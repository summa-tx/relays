use bitcoin_spv::types::*;

use std::{fs::File, io::Read};

use lazy_static::lazy_static;

use crate::error::RelayError;

lazy_static! {
    pub static ref TEST_VECTORS: TestJson = setup();
}

#[derive(serde::Deserialize, Debug, Clone, Copy)]
pub struct TestHeaderDetails {
    pub height: u32,
    pub hash: Hash256Digest,
    pub prevhash: Hash256Digest,
    pub merkle_root: Hash256Digest,
    pub raw: RawHeader,
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct AddHeadersCase {
    #[serde(rename = "rustSkip", default)]
    pub rust_skip: bool,
    pub comment: String,
    pub headers: Vec<TestHeaderDetails>,
    pub anchor: TestHeaderDetails,
    pub internal: bool,
    #[serde(rename = "isMainnet")]
    pub mainnet: bool,
    pub output: usize,
}

impl AddHeadersCase {
    pub fn flat_raw_headers(&self) -> Vec<u8> {
        self.headers
            .iter()
            .map(|v| v.raw.as_ref().to_vec().into_iter())
            .flatten()
            .collect::<Vec<u8>>()
    }
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct AddDiffChangeCase {
    #[serde(rename = "rustSkip", default)]
    pub rust_skip: bool,
    pub comment: String,
    pub headers: Vec<TestHeaderDetails>,
    #[serde(rename = "prevEpochStart")]
    pub prev_epoch_start: TestHeaderDetails,
    pub anchor: TestHeaderDetails,
    pub output: usize,
}

impl AddDiffChangeCase {
    pub fn flat_raw_headers(&self) -> Vec<u8> {
        self.headers
            .iter()
            .map(|v| v.raw.as_ref().to_vec().into_iter())
            .flatten()
            .collect::<Vec<u8>>()
    }
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct HeaderBlock {
    #[serde(rename = "validateHeaderChain")]
    pub header_chain_cases: Vec<AddHeadersCase>,
    #[serde(rename = "validateDifficultyChange")]
    pub diff_change_cases: Vec<AddDiffChangeCase>,
    // incomplete
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct TestJson {
    pub header: HeaderBlock,
    // incomplete
}

pub fn setup() -> TestJson {
    let mut file = File::open("../../testVectors.json").unwrap();
    let mut data = String::new();
    file.read_to_string(&mut data).unwrap();

    serde_json::from_str(&data).unwrap()
}

// Thanks I hate it
// TODO: make this not suck so bad
pub fn error_to_code(err: RelayError) -> usize {
    match err {
        RelayError::IncorrectDifficultyChange => 304, // BadRetarget
        RelayError::NotHeavier => 405,                // NotHeavier
        RelayError::NotLatestAncestor => 404,         // NotHeaviestAncestor
        RelayError::NotInBestChain => 701,
        RelayError::TooDeep => 701,
        RelayError::ReadOverrun => 701,
        RelayError::BadCompactInt => 701,
        RelayError::MalformattedOpReturnOutput => 108,
        RelayError::MalformattedP2SHOutput => 108,
        RelayError::MalformattedP2PKHOutput => 108,
        RelayError::MalformattedWitnessOutput => 108,
        RelayError::MalformattedOutput => 108,
        RelayError::WrongLengthHeader => 101, // BadHeaderLength
        RelayError::WrongEnd => 301,
        RelayError::UnexpectedDifficultyChange => 201, // UnexpectedRetarget
        RelayError::InsufficientWork => 108,
        RelayError::InvalidChain => 108,
        RelayError::WrongDigest => 108,
        RelayError::WrongMerkleRoot => 108,
        RelayError::WrongPrevHash => 108,
        RelayError::InvalidVin => 604,
        RelayError::InvalidVout => 605,
        RelayError::WrongTxID => 108,
        RelayError::BadMerkleProof => 108,
        RelayError::OutputLengthMismatch => 108,
        RelayError::UnknownError => 701,
    }
}
