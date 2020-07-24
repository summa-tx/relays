use bitcoin_spv::types::{Hash256Digest, HeaderArray, RawHeader};

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
    parent_index: u64,
    epoch_start_index: u64,
    height: u64,
}

#[repr(C)]
#[derive(Clone, Debug, Default, PartialEq, serde::Serialize, serde::Deserialize)]
/// A Bitcoin relay
pub struct Relay {
    relay_genesis: HeaderInfo,
    current_best_index: u64,
    epoch_start: Hash256Digest,
    best_known_digest: Hash256Digest,
    last_reorg_lca: Hash256Digest,
    // TODO: reap things older than 4032?
    header_store: Vec<HeaderInfo>,
}

impl Relay {
    // Call only after validating the chain
    fn attach_to_store(&mut self, anchor_index: u64, headers: &HeaderArray) {
        let anchor = self.header_store[anchor_index as usize];
        assert!(
            anchor.digest == headers.index(0).parent(),
            "Attempted to attach to wrong anchor"
        );
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
        genesis_height: u64,
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
            parent_index: u64::MAX, // will panic when indexing the vec
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

    /// Process the `AddHeaders` instruction
    pub fn process_add_headers(
        anchor_index: u64,
        anchor_bytes: Vec<u8>,
        header_bytes: Vec<u8>,
        internal: bool,
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let iter = &mut accounts.iter();
        let relay_state = next_account_info(iter)?;
        let mut relay = Self::get_relay(relay_state)?;

        let anchor_info = relay.header_store[anchor_index as usize];
        let anchor = RawHeader::new(&anchor_bytes).map_err(Into::<RelayError>::into)?;
        if anchor.digest() != anchor_info.digest {
            return Err(RelayError::WrongPrevHash.into());
        }
        let headers = HeaderArray::new(&header_bytes).map_err(Into::<RelayError>::into)?;

        let first_new = headers.index(0);
        if !internal && first_new.target() != anchor.target() {
            return Err(RelayError::UnexpectedDifficultyChange.into());
        }

        let _acc_diff = headers
            .valid_difficulty(true)
            .map_err(Into::<RelayError>::into)?;

        // If we haven't errored yet, they're good.
        relay.attach_to_store(anchor_index, &headers);

        Self::commit_relay(relay, relay_state)?;
        Ok(())
    }

    /// Process the `AddDifficultyChange` instruction
    pub fn process_add_difficulty_change(
        _old_period_end_index: u64,
        _headers: Vec<u8>, // should be a vec of [u8; 80]
    ) -> ProgramResult {
        unimplemented!()
    }

    /// Process the `MarkNewHeaviest` instruction
    pub fn process_mark_new_heaviest(
        _lca_index: u64,
        _current_best: Vec<u8>, // always 80 bytes
        _new_best_index: u64,
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
                    false, // internal
                    accounts,
                )
            }
            RelayInstruction::AddDifficultyChange {
                old_period_end_index,
                headers,
            } => {
                info!("Instruction: AddDifficultyChange");
                Self::process_add_difficulty_change(old_period_end_index, headers)
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
