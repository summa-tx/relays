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
        #[serde(with = "bad_header_ser")]
        genesis_header: [u8; 80],
        /// The height of the genesis header
        genesis_height: u32,
        /// The LE digest of the Bitcoin block that starts the epoch containing the genesis header
        epoch_start: [u8; 32],
    },
    /// AddHeaders command
    ///
    ///   0. `[writable]` Existing Relay to update
    AddHeaders {
        /// The index of the anchor in the state vector
        anchor_index: u32,
        /// The raw anchor header
        #[serde(with = "bad_header_ser")]
        anchor_bytes: [u8; 80],
        /// The tightly-packed raw headers
        headers: Vec<u8>, // should be a vec of [u8; 80]
    },

    /// AddDifficultyChange command
    ///
    ///   0. `[writable]` Existing Relay to update
    AddDifficultyChange {
        /// The raw old period start header
        #[serde(with = "bad_header_ser")]
        old_period_start_bytes: [u8; 80],
        /// The index of the old period end header in the state vector
        old_period_end_index: u32,
        /// The raw old period end header
        #[serde(with = "bad_header_ser")]
        old_period_end_bytes: [u8; 80],
        /// The tightly-packed raw headers
        headers: Vec<u8>,
    },

    /// MarkNewHeaviest command
    ///
    ///   0. `[writable]` Existing Relay to update
    MarkNewHeaviest {
        /// The index of the latest common ancestor header in the state vector
        lca_index: u32,
        /// The current best header
        #[serde(with = "bad_header_ser")]
        current_best: [u8; 80],
        /// The index of the new best header in the state vector
        new_best_index: u32,
        /// The new best header
        #[serde(with = "bad_header_ser")]
        new_best: [u8; 80],
    },
}

mod bad_header_ser {
    use std::fmt;

    use serde::{
        de::Error,
        de::{SeqAccess, Visitor},
        ser::SerializeTuple,
        Deserializer, Serializer,
    };

    pub fn deserialize<'de, D>(deserializer: D) -> Result<[u8; 80], D::Error>
    where
        D: Deserializer<'de>,
    {
        struct ArrayVisitor {}

        impl<'de> Visitor<'de> for ArrayVisitor {
            type Value = [u8; 80];

            fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
                formatter.write_str(concat!("an array of length 80"))
            }

            fn visit_seq<A>(self, mut seq: A) -> Result<[u8; 80], A::Error>
            where
                A: SeqAccess<'de>,
            {
                let mut arr = [0u8; 80];
                #[allow(clippy::needless_range_loop)]
                for i in 0..80 {
                    arr[i] = seq
                        .next_element()?
                        .ok_or_else(|| Error::invalid_length(i, &self))?;
                }
                Ok(arr)
            }
        }

        let visitor = ArrayVisitor {};
        deserializer.deserialize_tuple(80, visitor)
    }

    pub fn serialize<S>(buf: &[u8; 80], serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut seq = serializer.serialize_tuple(buf.len())?;
        for elem in &buf[..] {
            seq.serialize_element(elem)?;
        }
        seq.end()
    }
}
