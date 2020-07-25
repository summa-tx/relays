#![warn(missing_docs)]
//! A bitcoin Relay

/// Error type
pub mod error;
pub use error::*;

/// FakeVec
pub mod fake_vec;
pub use fake_vec::*;

use bitcoin_spv::{
    btcspv::retarget_algorithm,
    std_types::SPVProof,
    types::{Hash256Digest, HeaderArray, RawHeader},
};

#[repr(C)]
#[derive(Clone, Copy, Default, Debug, PartialEq, serde::Serialize, serde::Deserialize)]
/// Information about a header
pub struct HeaderInfo {
    /// The digest of the header
    pub digest: Hash256Digest,
    /// The index of the parent in the store
    pub parent_index: u32,
    /// The index of the epoch start in the store
    pub epoch_start_index: u32,
    /// The height of the header
    pub height: u32,
}

/// A Raw header with its associated info.
/// Convenience struct to avoid declaring twice as many let bindings
#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct RawWithInfo<'a> {
    /// The raw header
    pub raw: RawHeader,
    /// A reference to the header info in the store
    pub info: &'a HeaderInfo,
}


#[repr(C)]
#[derive(Clone, Debug, Default, PartialEq, serde::Serialize, serde::Deserialize)]
/// A Bitcoin relay
pub struct Relay {
    relay_genesis: HeaderInfo,
    pre_genesis_epoch_start: Hash256Digest,
    /// The index of the current best-known header in the store
    pub current_best_index: u32,
    /// The digest of the chaintip
    pub best_known_digest: Hash256Digest,
    /// The LCA of the most recent reorg or extension
    pub last_reorg_lca: Hash256Digest,
    // TODO: generalize
    header_store: FakeVec<HeaderInfo, generic_array::typenum::U4096>,
}

impl Relay {
    /// Instantiate a new relay
    pub fn new(
        genesis_header: [u8; 80],
        genesis_height: u32,
        epoch_start_digest: [u8; 32],
    ) -> Result<Self, RelayError> {
        let mut genesis = RawHeader::default();

        if genesis_header.len() != 80 {
            return Err(RelayError::WrongLengthHeader);
        }

        genesis.as_mut().copy_from_slice(&genesis_header[..80]);
        let genesis_digest = genesis.digest();

        if genesis_digest.as_ref()[28..32] != [0, 0, 0, 0] {
            return Err(RelayError::InsufficientWork);
        }

        let epoch_start = HeaderInfo {
            digest: epoch_start_digest.into(),
            parent_index: u32::MAX, // will panic when indexing the vec
            epoch_start_index: 0,
            height: genesis_height - (genesis_height % 2016),
        };

        let genesis_info = HeaderInfo {
            digest: genesis_digest,
            parent_index: u32::MAX, // will panic when indexing the vec
            epoch_start_index: 0,
            height: genesis_height,
        };

        let mut relay = Self {
            relay_genesis: genesis_info,
            pre_genesis_epoch_start: epoch_start_digest.into(),
            current_best_index: 1,
            best_known_digest: genesis_digest,
            last_reorg_lca: genesis_digest,
            ..Default::default()
        };
        relay.header_store.push(epoch_start); // index 0
        relay.header_store.push(genesis_info); // index 1
        Ok(relay)
    }

    /// Read the info store at a specified index
    pub fn read_info_store(&self, index: u32) -> &HeaderInfo {
        self.header_store.get(index as usize).unwrap() // panic if not found
    }

    /// Read the parent of a header
    pub fn parent_of(&self, info: &HeaderInfo) -> &HeaderInfo {
        let parent = self.read_info_store(info.parent_index);
        assert!(parent.height == info.height - 1);
        parent
    }

    /// Read the nth ancestor of a header. Linear look ups in the distance between them.
    pub fn ancestor_of<'a>(&'a self, info: &'a HeaderInfo, depth: u32) -> &'a HeaderInfo {
        let mut current = info;
        for _ in 0..depth {
            current = self.parent_of(info);
        }
        current
    }

    fn attach_metadata(&self, index: u32, raw: RawHeader) -> Result<RawWithInfo, RelayError> {
        let info = self.read_info_store(index);
        if info.digest != raw.digest() {
            return Err(RelayError::WrongDigest);
        }
        Ok(RawWithInfo { raw, info })
    }

    /// Load a header using its index and 80-bytes raw form
    pub fn load_header(&self, index: u32, raw: [u8; 80]) -> Result<RawWithInfo, RelayError> {
        let header = raw.into();
        self.attach_metadata(index, header)
    }

    // Call only after validating the chain
    fn attach_to_store(&mut self, anchor_index: u32, headers: &HeaderArray) {
        let anchor = self.read_info_store(anchor_index);
        // sanity check
        assert!(
            anchor.digest == headers.index(0).parent(),
            "Attempted to attach to wrong anchor"
        );

        // We'll update each of these for each header
        let mut parent_index = anchor_index;
        let mut height = anchor.height;
        let mut epoch_start_index = anchor.epoch_start_index;

        for header in headers.iter() {
            self.header_store.push(HeaderInfo {
                digest: header.digest(),
                parent_index,
                epoch_start_index,
                height,
            });
            parent_index += 1;
            height += 1;

            // Consider: the way this works, a header at 4032 will point to 2016, not to itself.
            if height % 2016 == 1 {
                epoch_start_index = parent_index;
            }
        }
    }

    /// Validate an SPV Proof
    pub fn validate_proof(
        &self,
        confirming_header_index: u32,
        proof: &SPVProof,
    ) -> Result<u32, RelayError> {
        proof.validate()?;
        let header_info = self.read_info_store(confirming_header_index);
        let chaintip = self.read_info_store(self.current_best_index);

        if chaintip.height - header_info.height > 2016 {
            return Err(RelayError::TooDeep);
        }

        let ancestor = self.ancestor_of(chaintip, chaintip.height - header_info.height);
        if header_info.digest != proof.confirming_header.hash {
            return Err(RelayError::WrongDigest);
        }
        if ancestor.digest != header_info.digest {
            return Err(RelayError::NotInBestChain);
        }
        Ok(chaintip.height - header_info.height + 1)
    }

    /// Add headers to the relay
    pub fn add_headers(
        &mut self,
        anchor_index: u32,
        anchor_bytes: [u8; 80],
        header_bytes: Vec<u8>,
        internal: bool,
    ) -> Result<(), RelayError> {
        let anchor = self.load_header(anchor_index, anchor_bytes)?;
        let headers = HeaderArray::new(&header_bytes).map_err(Into::<RelayError>::into)?;

        let first_new = headers.index(0);
        if !internal && first_new.target() != anchor.raw.target() {
            return Err(RelayError::UnexpectedDifficultyChange);
        }

        headers
            .valid_difficulty(true)
            .map_err(Into::<RelayError>::into)?;

        // If we haven't errored yet, they're good.
        self.attach_to_store(anchor_index, &headers);
        Ok(())
    }

    /// Add a difficulty change to the relay
    pub fn add_difficulty_change(
        &mut self,
        old_period_start_bytes: [u8; 80],
        old_period_end_index: u32,
        old_period_end_bytes: [u8; 80],
        header_bytes: Vec<u8>, // should be a vec of [u8; 80]
    ) -> Result<(), RelayError> {
        let headers = HeaderArray::new(&header_bytes)?;
        let old_period_end = self.load_header(old_period_end_index, old_period_end_bytes)?;
        let old_period_start = self.load_header(
            old_period_end.info.epoch_start_index,
            old_period_start_bytes,
        )?;

        // Ensure a change is allowed
        if old_period_end.info.height % 2016 != 2015 {
            return Err(RelayError::UnexpectedDifficultyChange);
        }

        // sanity checks. These should only fail if the store is corrupted
        assert!(old_period_start.info.height == old_period_end.info.height - 2015);
        assert!(old_period_start.raw.target() == old_period_end.raw.target());

        // Validate the difficulty change
        let new_target = headers.index(0).target();
        let expected_target = &retarget_algorithm(
            &old_period_start.raw.target(),
            old_period_start.raw.timestamp(),
            old_period_end.raw.timestamp(),
        );

        // NB:
        // This comparison looks weird because header nBits encoding truncates targets
        // The target is encoded as a 3-byte LE significand with a 1-byte mantissa in base 256.
        // It is expanded into a u256 for PoW checks.
        // But the new target is generated as a full-precision u256.
        if (new_target & expected_target) != *expected_target {
            return Err(RelayError::IncorrectDifficultyChange);
        }

        // Proceed to add the headers
        self.add_headers(
            old_period_end_index,
            old_period_end_bytes,
            header_bytes,
            true,
        )?;
        Ok(())
    }

    // true for right (should update), false for left
    fn verify_better_descendant<'a>(
        &self,
        ancestor: &'a HeaderInfo,
        left: &'a RawWithInfo,
        right: &'a RawWithInfo,
    ) -> Result<(), RelayError> {
        if ancestor.digest == left.info.digest && ancestor.digest == right.info.digest {
            return Err(RelayError::NotHeavier); // Everything is the same
        }

        // First we have to check that the ancestor is the latest common ancestor
        {
            let mut left_current = left.info;
            let mut right_current = right.info;
            let mut left_prev = left.info;
            let mut right_prev = right.info;

            // TODO: limit this better
            for _ in 0..200 {
                if left_prev != ancestor {
                    left_current = left_prev;
                    left_prev = self.read_info_store(left_prev.parent_index);
                }
                if right_prev != ancestor {
                    right_current = right_prev;
                    right_prev = self.read_info_store(right_prev.parent_index);
                }
            }
            if left_current == right_current {
                return Err(RelayError::NotLatestAncestor);
            } // newer ancestor exists.
            if left_prev != right_prev {
                return Err(RelayError::NotLatestAncestor);
            } // Common ancestor not found
        }

        // NB:
        // 1. Left is in a new window, right is in the old window. Left is heavier
        // 2. Right is in a new window, left is in the old window. Right is heavier
        // 3. Both are in the same window, choose the higher one
        // 4. They're in different new windows. Choose the heavier one
        let next_period_start_height = (ancestor.height + 2016) - (ancestor.height % 2016);
        let left_in_period = left.info.height < next_period_start_height;
        let right_in_period = right.info.height < next_period_start_height;

        if !left_in_period && right_in_period {
            return Err(RelayError::NotHeavier);
        }
        if left_in_period && !right_in_period {
            return Ok(());
        }
        if left_in_period && right_in_period {
            if right.info.height > left.info.height {
                return Ok(());
            } else {
                return Err(RelayError::NotHeavier);
            }
        }

        let left_acc_diff = (left.info.height % 2016) * left.raw.difficulty();
        let right_acc_diff = (right.info.height % 2016) * right.raw.difficulty();

        if right_acc_diff > left_acc_diff {
            Ok(())
        } else {
            Err(RelayError::NotHeavier)
        }
    }

    /// Mark a new heaviest digest
    pub fn mark_new_heaviest(
        &mut self,
        lca_index: u32,
        current_best: [u8; 80],
        new_best_index: u32,
        new_best: [u8; 80],
    ) -> Result<(), RelayError> {
        // Isolate all the borrows
        let (new_best, ancestor) = {
            let new_best = self.load_header(new_best_index, new_best)?;
            let current_best = self.load_header(self.current_best_index, current_best)?;
            let ancestor = self.read_info_store(lca_index);
            Self::verify_better_descendant(&self, &ancestor, &current_best, &new_best)?;
            (new_best.info.digest, ancestor.digest)
        };

        self.last_reorg_lca = ancestor;
        self.best_known_digest = new_best;
        self.current_best_index = new_best_index;
        Ok(())
    }
}
