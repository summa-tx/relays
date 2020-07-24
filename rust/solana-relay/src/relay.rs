use bitcoin_spv::{btcspv::retarget_algorithm, types::{Hash256Digest, HeaderArray, RawHeader}};

use solana_sdk::{
    account_info::{next_account_info, AccountInfo},
    entrypoint::ProgramResult,
    info,
    program_error::ProgramError,
    // program_error::ProgramError,
    pubkey::Pubkey,
};

use crate::{errors::*, instructions::*};

#[repr(C)]
#[derive(Clone, Debug, PartialEq, serde::Serialize, serde::Deserialize)]
/// The state of this account
pub enum State {
    /// Uninitialized
    Uninitialized,
    /// Actively running
    Active(Relay),
}

#[repr(C)]
#[derive(Clone, Copy, Default, Debug, PartialEq, serde::Serialize, serde::Deserialize)]
/// Information about a header
pub struct HeaderInfo {
    digest: Hash256Digest,
    parent_index: u32,
    epoch_start_index: u32,
    height: u32,
}

/// A Raw header with its associated info
#[derive(Debug, Clone, Copy)]
pub struct RawWithInfo<'a> {
    raw: RawHeader,
    info: &'a HeaderInfo,
}

#[repr(C)]
#[derive(Clone, Debug, Default, PartialEq, serde::Serialize, serde::Deserialize)]
/// A Bitcoin relay
pub struct Relay {
    relay_genesis: HeaderInfo,
    current_best_index: u32,
    epoch_start: Hash256Digest,
    best_known_digest: Hash256Digest,
    last_reorg_lca: Hash256Digest,
    // TODO: reap things older than 4032?
    header_store: Vec<HeaderInfo>,
}

impl Relay {
    fn read_header_store(&self, index: u32) -> &HeaderInfo {
        &self.header_store[index as usize]
    }

    fn attach_metadata(&self, index: u32, raw: RawHeader) -> Result<RawWithInfo, RelayError> {
        let info = self.read_header_store(index);
        if info.digest != raw.digest() {
            return Err(RelayError::WrongDigest);
        }
        Ok(RawWithInfo {
            raw,
            info,
        })
    }

    // Load a header using its index and 80-bytes form
    fn load_header<T: AsRef<[u8]>>(&self, index: u32, raw: &T) -> Result<RawWithInfo, RelayError> {
        let header = RawHeader::new(raw).map_err(Into::<RelayError>::into)?;
        self.attach_metadata(index, header)
    }


    // Call only after validating the chain
    fn attach_to_store(&mut self, anchor_index: u32, headers: &HeaderArray) {
        let anchor = self.read_header_store(anchor_index);
        // sanity check
        assert!(
            anchor.digest == headers.index(0).parent(),
            "Attempted to attach to wrong anchor"
        );

        // We'll update each of these for each header
        let mut parent_index = anchor_index;
        let mut height = anchor.height;
        let mut epoch_start_index = anchor.epoch_start_index;

        for header in headers.iter() {
            self.header_store.push(HeaderInfo {
                digest: header.digest(),
                parent_index,
                epoch_start_index,
                height,
            });
            parent_index += 1;
            height += 1;

            // Consider: the way this works, a header at 4032 will point to 2016, not to itself.
            if height % 2016 == 1 {
                epoch_start_index = parent_index;
            }
        }
    }
}

impl State {
    fn get_relay(relay_state: &AccountInfo) -> Result<Relay, ProgramError> {
        match serde_cbor::from_slice(&relay_state.try_borrow_data()?) {
            Ok(State::Uninitialized) => Ok(Relay::default()),
            Ok(State::Active(relay)) => Ok(relay),
            _ => Err(RelayError::AlreadyInit.into()),
        }
    }

    fn commit_relay(relay: Relay, relay_state: &AccountInfo) -> Result<(), ProgramError> {
        let serialized = serde_cbor::to_vec(&relay).expect("No serialization failure");
        let dest: &mut [u8] = &mut relay_state.data.borrow_mut();
        if dest.len() < serialized.len() {
            return Err(RelayError::InsufficientStateSpace.into());
        }
        dest[..serialized.len()].copy_from_slice(&serialized);
        dest[serialized.len()..].iter_mut().for_each(|i| *i = 0);

        Ok(())
    }

    /// Process the `Initialize` instruction
    pub fn process_initialize(
        genesis_header: Vec<u8>, // always 80 bytes,
        genesis_height: u32,
        epoch_start: [u8; 32],
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let iter = &mut accounts.iter();
        let relay_state = next_account_info(iter)?;
        let mut relay = Self::get_relay(relay_state)?;

        let mut genesis = RawHeader::default();

        if genesis_header.len() != 80 {
            return Err(RelayError::WrongLengthHeader.into());
        }

        genesis.as_mut().copy_from_slice(&genesis_header[..80]);
        let genesis_digest = genesis.digest();

        if genesis_digest.as_ref()[28..32] != [0, 0, 0, 0] {
            return Err(RelayError::InsufficientWork.into());
        }

        let genesis_info = HeaderInfo {
            digest: genesis_digest,
            parent_index: u32::MAX, // will panic when indexing the vec
            epoch_start_index: genesis_height - (genesis_height % 2016),
            height: genesis_height,
        };

        relay.relay_genesis = genesis_info;
        relay.epoch_start = epoch_start.into();
        relay.current_best_index = genesis_height;
        relay.best_known_digest = genesis_digest;
        relay.last_reorg_lca = genesis_digest;
        relay.header_store.push(genesis_info);

        Self::commit_relay(relay, relay_state)?;
        Ok(())
    }

    fn add_headers(
        relay: &mut Relay,
        anchor_index: u32,
        anchor_bytes: Vec<u8>,
        header_bytes: Vec<u8>,
        internal: bool,
    ) -> ProgramResult {
        let anchor = relay.load_header(anchor_index, &anchor_bytes)?;
        let headers = HeaderArray::new(&header_bytes).map_err(Into::<RelayError>::into)?;

        let first_new = headers.index(0);
        if !internal && first_new.target() != anchor.raw.target() {
            return Err(RelayError::UnexpectedDifficultyChange.into());
        }

        headers
            .valid_difficulty(true)
            .map_err(Into::<RelayError>::into)?;

        // If we haven't errored yet, they're good.
        relay.attach_to_store(anchor_index, &headers);
        Ok(())
    }

    /// Process the `AddHeaders` instruction
    pub fn process_add_headers(
        anchor_index: u32,
        anchor_bytes: Vec<u8>,
        header_bytes: Vec<u8>,
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let iter = &mut accounts.iter();
        let relay_state = next_account_info(iter)?;
        let mut relay = Self::get_relay(relay_state)?;

        Self::add_headers(
            &mut relay,
            anchor_index,
            anchor_bytes,
            header_bytes,
            false,
        )?;

        // Commit and return
        Self::commit_relay(relay, relay_state)?;
        Ok(())
    }

    /// Process the `AddDifficultyChange` instruction
    pub fn process_add_difficulty_change(
        old_period_start_bytes: Vec<u8>,
        old_period_end_index: u32,
        old_period_end_bytes: Vec<u8>,
        header_bytes: Vec<u8>, // should be a vec of [u8; 80]
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let iter = &mut accounts.iter();
        let relay_state = next_account_info(iter)?;
        let mut relay = Self::get_relay(relay_state)?;

        let headers = HeaderArray::new(&header_bytes).map_err(Into::<RelayError>::into)?;
        let old_period_end = relay.load_header(old_period_end_index, &old_period_end_bytes)?;
        let old_period_start = relay.load_header(old_period_end.info.epoch_start_index, &old_period_start_bytes)?;

        // Ensure a change is allowed
        if old_period_end.info.height % 2016 != 2015 {
            return Err(RelayError::UnexpectedDifficultyChange.into());
        }

        // sanity checks. These should only fail if the store is corrupted
        assert!(old_period_start.info.height == old_period_end.info.height - 2015);
        assert!(old_period_start.raw.target() == old_period_end.raw.target());

        // Validate the difficulty change
        let new_target = headers.index(0).target();
        let expected_target = &retarget_algorithm(
            &old_period_start.raw.target(),
            old_period_start.raw.timestamp(),
            old_period_end.raw.timestamp()
        );

        // NB:
        // This comparison looks weird because header nBits encoding truncates targets
        // The target is encoded as a 3-byte LE significand with a 1-byte mantissa in base 256.
        // It is expanded into a u256 for PoW checks.
        // But the new target is generated as a full-precision u256.
        if (new_target & expected_target) != *expected_target {
            return Err(RelayError::IncorrectDifficultyChange.into());
        }

        // Proceed to add the headers
        Self::add_headers(
            &mut relay,
            old_period_end_index,
            old_period_end_bytes,
            header_bytes,
            true,
        )?;

        Self::commit_relay(relay, relay_state)?;
        Ok(())
    }

    /// Process the `MarkNewHeaviest` instruction
    pub fn process_mark_new_heaviest(
        _lca_index: u32,
        _current_best: Vec<u8>, // always 80 bytes
        _new_best_index: u32,
        _new_best: Vec<u8>, // always 80 bytes
    ) -> ProgramResult {
        unimplemented!()
    }

    /// Processes an [Instruction](enum.Instruction.html).
    pub fn process(_program_id: &Pubkey, accounts: &[AccountInfo], input: &[u8]) -> ProgramResult {
        let instruction = serde_cbor::from_slice(input).expect("Invalid instruction");
        match instruction {
            RelayInstruction::Initialize {
                genesis_header,
                genesis_height,
                epoch_start,
            } => {
                info!("Instruction: Initialize");
                Self::process_initialize(genesis_header, genesis_height, epoch_start, accounts)
            }
            RelayInstruction::AddHeaders {
                anchor_index,
                anchor_bytes,
                headers,
            } => {
                info!("Instruction: AddHeaders");
                Self::process_add_headers(
                    anchor_index,
                    anchor_bytes,
                    headers,
                    accounts,
                )
            }
            RelayInstruction::AddDifficultyChange {
                old_period_start_bytes,
                old_period_end_index,
                old_period_end_bytes,
                headers,
            } => {
                info!("Instruction: AddDifficultyChange");
                Self::process_add_difficulty_change(
                    old_period_start_bytes,
                    old_period_end_index,
                    old_period_end_bytes,
                    headers,
                    accounts,
                )
            }
            RelayInstruction::MarkNewHeaviest {
                lca_index,
                current_best,
                new_best_index,
                new_best,
            } => {
                info!("Instruction: MarkNewHeaviest");
                Self::process_mark_new_heaviest(lca_index, current_best, new_best_index, new_best)
            }
        }
    }
}
