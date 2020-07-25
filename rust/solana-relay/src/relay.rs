use solana_sdk::{
    account_info::{next_account_info, AccountInfo},
    entrypoint::ProgramResult,
    info,
    program_error::ProgramError,
    pubkey::Pubkey,
};

use relay_core::Relay;

use crate::{error::*, instructions::*};

#[repr(C)]
#[derive(Clone, Debug, PartialEq, serde::Serialize, serde::Deserialize)]
#[allow(clippy::large_enum_variant)]
/// The state of this account
pub enum State {
    /// Uninitialized
    Uninitialized,
    /// Actively running
    // Active variant is 180376 bytes
    Active(Relay),
}

impl State {
    fn get_relay(relay_state: &AccountInfo) -> Result<Option<Relay>, ProgramError> {
        match serde_cbor::from_slice(&relay_state.try_borrow_data()?) {
            Ok(State::Uninitialized) => Ok(None),
            Ok(State::Active(relay)) => Ok(Some(relay)),
            _ => Err(SolanaRelayError::UnknownError.into()),
        }
    }

    fn commit_relay(relay: Relay, relay_state: &AccountInfo) -> Result<(), ProgramError> {
        let serialized = serde_cbor::to_vec(&relay).expect("No serialization failure");
        let dest: &mut [u8] = &mut relay_state.data.borrow_mut();
        if dest.len() < serialized.len() {
            return Err(SolanaRelayError::InsufficientStateSpace.into());
        }
        dest[..serialized.len()].copy_from_slice(&serialized);
        dest[serialized.len()..].iter_mut().for_each(|i| *i = 0);

        Ok(())
    }

    /// Process the `Initialize` instruction
    pub fn process_initialize(
        genesis_header: [u8; 80],
        genesis_height: u32,
        epoch_start_digest: [u8; 32],
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let iter = &mut accounts.iter();
        let relay_state = next_account_info(iter)?;

        if Self::get_relay(relay_state)?.is_some() {
            return Err(SolanaRelayError::AlreadyInit.into());
        }

        let relay = Relay::new(genesis_header, genesis_height, epoch_start_digest)
            .map_err(Into::<SolanaRelayError>::into)?;

        Self::commit_relay(relay, relay_state)?;
        Ok(())
    }

    /// Process the `AddHeaders` instruction
    pub fn process_add_headers(
        anchor_index: u32,
        anchor_bytes: [u8; 80],
        header_bytes: Vec<u8>,
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let iter = &mut accounts.iter();
        let relay_state = next_account_info(iter)?;
        let mut relay = Self::get_relay(relay_state)?.ok_or(SolanaRelayError::NotYetInit)?;

        relay
            .add_headers(anchor_index, anchor_bytes, header_bytes, false)
            .map_err(Into::<SolanaRelayError>::into)?;

        // Commit and return
        Self::commit_relay(relay, relay_state)?;
        Ok(())
    }

    /// Process the `AddDifficultyChange` instruction
    pub fn process_add_difficulty_change(
        old_period_start_bytes: [u8; 80],
        old_period_end_index: u32,
        old_period_end_bytes: [u8; 80],
        header_bytes: Vec<u8>, // should be a vec of [u8; 80]
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let iter = &mut accounts.iter();
        let relay_state = next_account_info(iter)?;
        let mut relay = Self::get_relay(relay_state)?.ok_or(SolanaRelayError::NotYetInit)?;

        relay
            .add_difficulty_change(
                old_period_start_bytes,
                old_period_end_index,
                old_period_end_bytes,
                header_bytes,
            )
            .map_err(Into::<SolanaRelayError>::into)?;

        Self::commit_relay(relay, relay_state)?;
        Ok(())
    }

    /// Process the `MarkNewHeaviest` instruction
    pub fn process_mark_new_heaviest(
        lca_index: u32,
        current_best: [u8; 80],
        new_best_index: u32,
        new_best: [u8; 80],
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let iter = &mut accounts.iter();
        let relay_state = next_account_info(iter)?;
        let mut relay = Self::get_relay(relay_state)?.ok_or(SolanaRelayError::NotYetInit)?;

        relay
            .mark_new_heaviest(lca_index, current_best, new_best_index, new_best)
            .map_err(Into::<SolanaRelayError>::into)?;

        Self::commit_relay(relay, relay_state)?;
        Ok(())
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
                Self::process_add_headers(anchor_index, anchor_bytes, headers, accounts)
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
                Self::process_mark_new_heaviest(
                    lca_index,
                    current_best,
                    new_best_index,
                    new_best,
                    accounts,
                )
            }
        }
    }
}
