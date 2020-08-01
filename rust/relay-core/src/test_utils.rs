use bitcoin_spv::types::*;

use std::{fs::File, io::Read};

use lazy_static::lazy_static;

use crate::error::RelayError;

lazy_static! {
    pub static ref TEST_VECTORS: TestJson = setup();
}

macro_rules! impl_has_header_chain {
    ($name:ty, $prop:ident) => {
        impl $name {
            pub fn flat_raw_headers(&self) -> Vec<u8> {
                self.$prop
                    .iter()
                    .map(|v| v.raw.as_ref().to_vec().into_iter())
                    .flatten()
                    .collect::<Vec<u8>>()
            }
        }
    };
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

impl_has_header_chain!(AddHeadersCase, headers);

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

impl_has_header_chain!(AddDiffChangeCase, headers);

#[derive(serde::Deserialize, Debug, Clone)]
pub struct HeaderBlock {
    #[serde(rename = "validateHeaderChain")]
    pub header_chain_cases: Vec<AddHeadersCase>,
    #[serde(rename = "validateDifficultyChange")]
    pub diff_change_cases: Vec<AddDiffChangeCase>,
    // incomplete
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct LCABlock {
    pub orphan: TestHeaderDetails,
    #[serde(rename = "oldPeriodStart")]
    pub prev_epoch_start: TestHeaderDetails,
    pub genesis: TestHeaderDetails,
    #[serde(rename = "preRetargetChain")]
    pub pre_retarget_chain: Vec<TestHeaderDetails>,
    #[serde(rename = "postRetargetChain")]
    pub post_retarget_chain: Vec<TestHeaderDetails>,
    // incomplete
    // testCases
}

impl LCABlock {
    pub fn flat_raw_pre(&self) -> Vec<u8> {
        self.pre_retarget_chain
            .iter()
            .map(|v| v.raw.as_ref().to_vec().into_iter())
            .flatten()
            .collect::<Vec<u8>>()
    }

    pub fn flat_raw_post(&self) -> Vec<u8> {
        self.post_retarget_chain
            .iter()
            .map(|v| v.raw.as_ref().to_vec().into_iter())
            .flatten()
            .collect::<Vec<u8>>()
    }
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct NewHeaviestBlockCase {
    #[serde(rename = "bestKnownDigest")]
    pub best_known_digest: Hash256Digest,
    pub ancestor: Hash256Digest,
    #[serde(rename = "currentBest")]
    pub current_best: RawHeader,
    #[serde(rename = "newBest")]
    pub new_best: RawHeader,
    pub limit: usize,
    pub error: usize,
    pub output: String,
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct ChainBlock {
    #[serde(rename = "isMostRecentCommonAncestor")]
    pub is_most_recent_common_ancestor: LCABlock,
    #[serde(rename = "markNewHeaviest")]
    pub mark_new_heaviest_cases: Vec<NewHeaviestBlockCase>,
    // incomplete
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct TestJson {
    pub header: HeaderBlock,
    pub chain: ChainBlock,
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
