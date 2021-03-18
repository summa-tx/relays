use crate::{error::SolanaRelayError, relay::State};
use solana_sdk::{
    account_info::AccountInfo, entrypoint, entrypoint::ProgramResult,
    program_error::PrintProgramError, pubkey::Pubkey,
};

entrypoint!(process_instruction);

#[allow(dead_code)]
fn process_instruction<'a>(
    program_id: &Pubkey,
    accounts: &'a [AccountInfo<'a>],
    instruction_data: &[u8],
) -> ProgramResult {
    match State::process(program_id, accounts, instruction_data) {
        Err(error) => {
            // catch the error so we can print it
            error.print::<SolanaRelayError>();
            Err(error)
        }
        Ok(_) => Ok(()),
    }
}
