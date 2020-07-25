#![warn(missing_docs)]

//! A Bitcoin Relay for Solana

/// A generic relay, parameterized by the chain KVStore
pub mod relay;

/// Utils
pub mod utils;

/// State transitions
pub mod instructions;

/// Errors
pub mod error;

/// Entrypoint
pub mod entry;
