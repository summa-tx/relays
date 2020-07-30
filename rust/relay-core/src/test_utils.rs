use bitcoin_spv::types::*;

use std::{
    fs::File,
    io::Read,
};

use lazy_static::lazy_static;

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
        self.headers.iter().map(|v| v.raw.as_ref().to_vec().into_iter()).flatten().collect::<Vec<u8>>()
    }
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct HeaderBlock {
    #[serde(rename = "validateHeaderChain")]
    pub header_chain_cases: Vec<AddHeadersCase>,
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
