use bitcoin_spv::types::SPVError;

use num_derive::FromPrimitive;
use num_traits::FromPrimitive;

use solana_sdk::{
    decode_error::DecodeError,
    info,
    program_error::{PrintProgramError, ProgramError},
};

use thiserror::Error;

use relay_core::RelayError;

/// Errors that may be returned by the TokenSwap program.
///
/// Most of these are copied directly from bitcoin spv
#[derive(Clone, Debug, Eq, Error, FromPrimitive, PartialEq)]
pub enum SolanaRelayError {
    /// Relay is not yet init
    #[error("NotYetInit")]
    NotYetInit,
    /// Relay is already init
    #[error("AlreadyInit")]
    AlreadyInit,
    /// State account has insufficient space to update the state
    #[error("InsufficientStateSpace")]
    InsufficientStateSpace,
    /// A difficulty change mismatch occured while adding headers
    #[error("IncorrectDifficultyChange")]
    IncorrectDifficultyChange,
    /// Failed to update to new heaviest header because new digest is not heavier
    #[error("NotHeavier")]
    NotHeavier,
    /// Failed to update to new heaviest header because a later ancestor was found
    #[error("NotLatestAncestor")]
    NotLatestAncestor,
    /// Confirming header of an SPV Proof is not confirmed by the chain tip
    #[error("NotInBestChain")]
    NotInBestChain,
    /// Confirming header of an SPV Proof is too deep to read
    #[error("TooDeep")]
    TooDeep,

    // Below copied from bitcoin-spv
    /// Overran a checked read on a slice
    #[error("ReadOverrun")]
    ReadOverrun,
    /// Attempted to parse a CompactInt without enough bytes
    #[error("BadCompactInt")]
    BadCompactInt,
    /// Called `extract_op_return_data` on an output without an op_return.
    #[error("MalformattedOpReturnOutput")]
    MalformattedOpReturnOutput,
    /// `extract_hash` identified a SH output prefix without a SH postfix.
    #[error("MalformattedP2SHOutput")]
    MalformattedP2SHOutput,
    /// `extract_hash` identified a PKH output prefix without a PKH postfix.
    #[error("MalformattedP2PKHOutput")]
    MalformattedP2PKHOutput,
    /// `extract_hash` identified a Witness output with a bad length tag.
    #[error("MalformattedWitnessOutput")]
    MalformattedWitnessOutput,
    /// `extract_hash` could not identify the output type.
    #[error("MalformattedOutput")]
    MalformattedOutput,
    /// Header not exactly 80 bytes.
    #[error("WrongLengthHeader")]
    WrongLengthHeader,
    /// Header chain changed difficulties unexpectedly
    #[error("UnexpectedDifficultyChange")]
    UnexpectedDifficultyChange,
    /// Header does not meet its own difficulty target.
    #[error("InsufficientWork")]
    InsufficientWork,
    /// Header in chain does not correctly reference parent header.
    #[error("InvalidChain")]
    InvalidChain,
    /// When validating a `BitcoinHeader`, the `hash` field is not the digest
    /// of the raw header.
    #[error("WrongDigest")]
    WrongDigest,
    /// When validating a `BitcoinHeader`, the `merkle_root` field does not
    /// match the root found in the raw header.
    #[error("WrongMerkleRoot")]
    WrongMerkleRoot,
    /// When validating a `BitcoinHeader`, the `prevhash` field does not
    /// match the parent hash found in the raw header.
    #[error("WrongPrevHash")]
    WrongPrevHash,
    /// A `vin` (transaction input vector) is malformatted.
    #[error("InvalidVin")]
    InvalidVin,
    /// A `vout` (transaction output vector) is malformatted.
    #[error("InvalidVout")]
    InvalidVout,
    /// When validating an `SPVProof`, the `tx_id` field is not the digest
    /// of the `version`, `vin`, `vout`, and `locktime`.
    #[error("WrongTxID")]
    WrongTxID,
    /// When validating an `SPVProof`, the `intermediate_nodes` is not a valid
    /// merkle proof connecting the `tx_id_le` to the `confirming_header`.
    #[error("BadMerkleProof")]
    BadMerkleProof,
    /// TxOut's reported length does not match passed-in byte slice's length
    #[error("OutputLengthMismatch")]
    OutputLengthMismatch,
    /// Any other error
    #[error("UnknownError")]
    UnknownError,
}

impl From<SPVError> for SolanaRelayError {
    fn from(e: SPVError) -> SolanaRelayError {
        RelayError::from(e).into()
    }
}

impl From<SolanaRelayError> for ProgramError {
    fn from(e: SolanaRelayError) -> Self {
        ProgramError::Custom(e as u32)
    }
}

impl<T> DecodeError<T> for SolanaRelayError {
    fn type_of() -> &'static str {
        "Relay Error"
    }
}

impl From<RelayError> for SolanaRelayError {
    fn from(e: RelayError) -> Self {
        match e {
            RelayError::IncorrectDifficultyChange => SolanaRelayError::IncorrectDifficultyChange,
            RelayError::NotHeavier => SolanaRelayError::NotHeavier,
            RelayError::NotLatestAncestor => SolanaRelayError::NotLatestAncestor,
            RelayError::NotInBestChain => SolanaRelayError::NotInBestChain,
            RelayError::TooDeep => SolanaRelayError::TooDeep,
            RelayError::ReadOverrun => SolanaRelayError::ReadOverrun,
            RelayError::BadCompactInt => SolanaRelayError::BadCompactInt,
            RelayError::MalformattedOpReturnOutput => SolanaRelayError::MalformattedOpReturnOutput,
            RelayError::MalformattedP2SHOutput => SolanaRelayError::MalformattedP2SHOutput,
            RelayError::MalformattedP2PKHOutput => SolanaRelayError::MalformattedP2PKHOutput,
            RelayError::MalformattedWitnessOutput => SolanaRelayError::MalformattedWitnessOutput,
            RelayError::MalformattedOutput => SolanaRelayError::MalformattedOutput,
            RelayError::WrongLengthHeader => SolanaRelayError::WrongLengthHeader,
            RelayError::UnexpectedDifficultyChange => SolanaRelayError::UnexpectedDifficultyChange,
            RelayError::InsufficientWork => SolanaRelayError::InsufficientWork,
            RelayError::InvalidChain => SolanaRelayError::InvalidChain,
            RelayError::WrongDigest => SolanaRelayError::WrongDigest,
            RelayError::WrongMerkleRoot => SolanaRelayError::WrongMerkleRoot,
            RelayError::WrongPrevHash => SolanaRelayError::WrongPrevHash,
            RelayError::InvalidVin => SolanaRelayError::InvalidVin,
            RelayError::InvalidVout => SolanaRelayError::InvalidVout,
            RelayError::WrongTxID => SolanaRelayError::WrongTxID,
            RelayError::BadMerkleProof => SolanaRelayError::BadMerkleProof,
            RelayError::OutputLengthMismatch => SolanaRelayError::OutputLengthMismatch,
            RelayError::UnknownError => SolanaRelayError::UnknownError,
        }
    }
}

impl PrintProgramError for SolanaRelayError {
    fn print<E>(&self)
    where
        E: 'static + std::error::Error + DecodeError<E> + PrintProgramError + FromPrimitive,
    {
        match self {
            SolanaRelayError::NotYetInit => info!("SolanaRelayError: NotYetInit"),
            SolanaRelayError::AlreadyInit => info!("SolanaRelayError: AlreadyInit"),
            SolanaRelayError::InsufficientStateSpace => {
                info!("SolanaRelayError: InsufficientStateSpace")
            }
            SolanaRelayError::IncorrectDifficultyChange => {
                info!("SolanaRelayError: IncorrectDifficultyChange")
            }
            SolanaRelayError::NotHeavier => info!("SolanaRelayError: NotHeavier"),
            SolanaRelayError::NotLatestAncestor => info!("SolanaRelayError: NotLatestAncestor"),
            SolanaRelayError::NotInBestChain => info!("SolanaRelayError: NotInBestChain"),
            SolanaRelayError::TooDeep => info!("SolanaRelayError: TooDeep"),
            SolanaRelayError::ReadOverrun => info!("SolanaRelayError: ReadOverrun"),
            SolanaRelayError::BadCompactInt => info!("SolanaRelayError: BadCompactInt"),
            SolanaRelayError::MalformattedOpReturnOutput => {
                info!("SolanaRelayError: MalformattedOpReturnOutput")
            }
            SolanaRelayError::MalformattedP2SHOutput => {
                info!("SolanaRelayError: MalformattedP2SHOutput")
            }
            SolanaRelayError::MalformattedP2PKHOutput => {
                info!("SolanaRelayError: MalformattedP2PKHOutput")
            }
            SolanaRelayError::MalformattedWitnessOutput => {
                info!("SolanaRelayError: MalformattedWitnessOutput")
            }
            SolanaRelayError::MalformattedOutput => info!("SolanaRelayError: MalformattedOutput"),
            SolanaRelayError::WrongLengthHeader => info!("SolanaRelayError: WrongLengthHeader"),
            SolanaRelayError::UnexpectedDifficultyChange => {
                info!("SolanaRelayError: UnexpectedDifficultyChange")
            }
            SolanaRelayError::InsufficientWork => info!("SolanaRelayError: InsufficientWork"),
            SolanaRelayError::InvalidChain => info!("SolanaRelayError: InvalidChain"),
            SolanaRelayError::WrongDigest => info!("SolanaRelayError: WrongDigest"),
            SolanaRelayError::WrongMerkleRoot => info!("SolanaRelayError: WrongMerkleRoot"),
            SolanaRelayError::WrongPrevHash => info!("SolanaRelayError: WrongPrevHash"),
            SolanaRelayError::InvalidVin => info!("SolanaRelayError: InvalidVin"),
            SolanaRelayError::InvalidVout => info!("SolanaRelayError: InvalidVout"),
            SolanaRelayError::WrongTxID => info!("SolanaRelayError: WrongTxID"),
            SolanaRelayError::BadMerkleProof => info!("SolanaRelayError: BadMerkleProof"),
            SolanaRelayError::OutputLengthMismatch => {
                info!("SolanaRelayError: OutputLengthMismatch")
            }
            SolanaRelayError::UnknownError => info!("SolanaRelayError: UnknownError"),
        }
    }
}
