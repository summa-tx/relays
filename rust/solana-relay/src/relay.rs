use bitcoin_spv::types::Hash256Digest;

use solana_sdk::{
    account_info::{AccountInfo, next_account_info},
    entrypoint::ProgramResult,
    info,
    // program_error::ProgramError,
    pubkey::Pubkey,
};

use crate::{
    errors::*,
    instructions::*,
};

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
#[derive(Clone, Copy, Debug, PartialEq, serde::Serialize, serde::Deserialize)]
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
    relay_genesis: Hash256Digest,
    current_best_index: u64,
    best_known_digest: Hash256Digest,
    last_reorg_lca: Hash256Digest,
    // TODO: reap things older than 4032?
    header_store: Vec<HeaderInfo>,
}

impl State {
    /// Process the `Initialize` instruction
    pub fn process_initialize(
        _genesis_header: Vec<u8>, // always 80 bytes,
        _genesis_height: u64,
        _epoch_start: [u8; 32],
        accounts: &[AccountInfo],
    ) -> ProgramResult {
        let account_info_iter = &mut accounts.iter();
        let relay_info = next_account_info(account_info_iter)?;

        // lol json
        match serde_json::from_slice::<State>(&relay_info.data.borrow()) {
            Ok(State::Uninitialized) => {},
            _ => { return Err(RelayError::AlreadyInit.into()); },
        }

        unimplemented!()
    }

    /// Process the `AddHeaders` instruction
    pub fn process_add_headers(
        _anchor_index: u64,
        _headers: Vec<u8>, // should be a vec of [u8; 80]
    ) -> ProgramResult {
        unimplemented!()
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
        // lol json
        let instruction = serde_json::from_slice(input).expect("Invalid instruction");
        match instruction {
            RelayInstruction::Initialize {
                genesis_header,
                genesis_height,
                epoch_start,
            } => {
                info!("Instruction: Initialize");
                Self::process_initialize(
                    genesis_header,
                    genesis_height,
                    epoch_start,
                    accounts,
                )
            },
            RelayInstruction::AddHeaders {
                anchor_index,
                headers,
            } => {
                info!("Instruction: AddHeaders");
                Self::process_add_headers(
                    anchor_index,
                    headers,
                )
            }
            RelayInstruction::AddDifficultyChange {
                old_period_end_index,
                headers,
            } => {
                info!("Instruction: AddDifficultyChange");
                Self::process_add_difficulty_change(
                    old_period_end_index,
                    headers,
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
                )
            }
        }
    }
}
