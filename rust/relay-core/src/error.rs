use bitcoin_spv::types::SPVError;
use thiserror::Error;

/// Errors that may be returned by the TokenSwap program.
///
/// Most of these are copied directly from bitcoin spv
#[derive(Clone, Debug, Eq, Error, PartialEq)]
pub enum RelayError {
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

impl From<SPVError> for RelayError {
    fn from(e: SPVError) -> RelayError {
        match e {
            SPVError::ReadOverrun => RelayError::ReadOverrun,
            SPVError::BadCompactInt => RelayError::BadCompactInt,
            SPVError::MalformattedOpReturnOutput => RelayError::MalformattedOpReturnOutput,
            SPVError::MalformattedP2SHOutput => RelayError::MalformattedP2SHOutput,
            SPVError::MalformattedP2PKHOutput => RelayError::MalformattedP2PKHOutput,
            SPVError::MalformattedWitnessOutput => RelayError::MalformattedWitnessOutput,
            SPVError::MalformattedOutput => RelayError::MalformattedOutput,
            SPVError::WrongLengthHeader => RelayError::WrongLengthHeader,
            SPVError::UnexpectedDifficultyChange => RelayError::UnexpectedDifficultyChange,
            SPVError::InsufficientWork => RelayError::InsufficientWork,
            SPVError::InvalidChain => RelayError::InvalidChain,
            SPVError::WrongDigest => RelayError::WrongDigest,
            SPVError::WrongMerkleRoot => RelayError::WrongMerkleRoot,
            SPVError::WrongPrevHash => RelayError::WrongPrevHash,
            SPVError::InvalidVin => RelayError::InvalidVin,
            SPVError::InvalidVout => RelayError::InvalidVout,
            SPVError::WrongTxID => RelayError::WrongTxID,
            SPVError::BadMerkleProof => RelayError::BadMerkleProof,
            SPVError::OutputLengthMismatch => RelayError::OutputLengthMismatch,
            SPVError::UnknownError => RelayError::UnknownError,
        }
    }
}
