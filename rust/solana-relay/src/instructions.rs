// NB: Vec<u8> that are always 80 bytes should be a [u8; 80],
//   but we can't derive serde that way

/// Instructions for the relay
///
/// All instructions require only this account. No outside interaction required
#[repr(C)]
#[derive(Clone, serde::Serialize, serde::Deserialize)]
pub enum RelayInstruction {
    /// Initialize the relay
    ///   0. `[writable, signer]` New Relay to create.
    Initialize {
        /// The genesis header with which to initialize the relay
        genesis_header: Vec<u8>, // always 80 bytes,
        /// The height of the genesis header
        genesis_height: u64,
        /// The LE digest of the Bitcoin block that starts the epoch containing the genesis header
        epoch_start: [u8; 32],
    },
    /// AddHeaders command
    ///
    ///   0. `[writable]` Existing Relay to update
    AddHeaders {
        /// The index of the anchor in the state vector
        anchor_index: u64,
        /// The raw anchor header
        anchor_bytes: Vec<u8>, // always 80 bytes
        /// The tightly-packed raw headers
        headers: Vec<u8>, // should be a vec of [u8; 80]
    },

    /// AddDifficultyChange command
    ///
    ///   0. `[writable]` Existing Relay to update
    AddDifficultyChange {
        /// The index of the old period end header in the state vector
        old_period_end_index: u64,
        /// The tightly-packed raw headers
        headers: Vec<u8>,
    },

    /// MarkNewHeaviest command
    ///
    ///   0. `[writable]` Existing Relay to update
    MarkNewHeaviest {
        /// The index of the latest common ancestor header in the state vector
        lca_index: u64,
        /// The current best header
        current_best: Vec<u8>, // always 80 bytes
        /// The index of the new best header in the state vector
        new_best_index: u64,
        /// The new best header
        new_best: Vec<u8>, // always 80 bytes
    },
}
