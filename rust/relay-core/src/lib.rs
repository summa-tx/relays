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
use generic_array::typenum::U4096;

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
/// A Bitcoin relay.
///
/// This struct tracks the state of the bitcoin mainnet and stores the most information about the
/// most recent 4096 headers. It is designed to be run on-chain.
///
/// Rather than querying the blockchain, this relay passively accepts proofs of chain updates.
/// These proofs are validated and then added to the state. See this repo's README for more
/// information.
pub struct Relay {
    mainnet: bool,
    relay_genesis: HeaderInfo,
    pre_genesis_epoch_start: Hash256Digest,
    /// The index of the current best-known header in the store
    pub current_best_index: u32,
    /// The digest of the chaintip
    pub best_known_digest: Hash256Digest,
    /// The LCA of the most recent reorg or extension
    pub last_reorg_lca: Hash256Digest,
    // TODO: generalize?
    header_store: FakeVec<HeaderInfo, U4096>,
}

impl Relay {
    /// Instantiate a new relay by inserting a trusted state.
    ///
    /// # Arguments
    ///
    /// - mainnet - True for Bitcoin mainnet, false for Bitcoin testnet. Setting to false
    ///    disables difficulty checks and RENDERS THE RELAY INSECURE.
    /// - genesis_header - an 80-byte bitcoin header that serves as the anchor of the relay.
    /// - genesis_height - the height of the genesis header in the Bitcoin chain
    /// - epoch_start_digest - the 32-byte LE hash of the header that began the difficulty epoch
    ///    containing the genesis header. its height will always be `0 % 2016`
    pub fn new<H: Into<RawHeader>, D: Into<Hash256Digest>>(
        mainnet: bool,
        genesis_header: H,
        genesis_height: u32,
        epoch_start_digest: D,
    ) -> Result<Self, RelayError> {
        let epoch_start_digest: Hash256Digest = epoch_start_digest.into();
        let genesis: RawHeader = genesis_header.into();
        let genesis_digest = genesis.digest();
        // Sanity check
        if genesis_digest.as_ref()[28..32] != [0, 0, 0, 0] {
            return Err(RelayError::InsufficientWork);
        }

        let epoch_start = HeaderInfo {
            digest: epoch_start_digest,
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
        // Thanks, I hate it
        let header_store: FakeVec<HeaderInfo, U4096> = unsafe {
            let hs = [HeaderInfo::default(); 4096];
            let store = std::mem::transmute::<
                [HeaderInfo; 4096],
                generic_array::GenericArray<HeaderInfo, U4096>,
            >(hs);
            FakeVec {
                next: 0,
                internal: store,
            }
        };

        let mut relay = Self {
            mainnet,
            relay_genesis: genesis_info,
            pre_genesis_epoch_start: epoch_start_digest,
            current_best_index: 1,
            best_known_digest: genesis_digest,
            last_reorg_lca: genesis_digest,
            header_store,
        };
        assert_eq!(relay.header_store.push(epoch_start), 0);
        assert_eq!(relay.header_store.push(genesis_info), 1);
        Ok(relay)
    }

    /// Read the info store at a specified index.
    pub fn read_metadata_store(&self, index: u32) -> &HeaderInfo {
        self.header_store.get(index as usize).unwrap() // panic if not found
    }

    /// Read the stored metadata about the parent of a header.
    pub fn parent_of(&self, info: &HeaderInfo) -> &HeaderInfo {
        let parent = self.read_metadata_store(info.parent_index);
        assert!(parent.height == info.height - 1);
        parent
    }

    /// Read the stored metadata about the nth ancestor of a header.
    ///
    /// # Note
    ///
    /// This function makes linear state look ups in the distance between them.
    pub fn ancestor_of<'a>(&'a self, info: &'a HeaderInfo, depth: u32) -> &'a HeaderInfo {
        let mut current = info;
        for _ in 0..depth {
            current = self.parent_of(info);
        }
        current
    }

    // Convenience function for associateing a header with its metadata
    fn attach_metadata(&self, index: u32, raw: RawHeader) -> Result<RawWithInfo, RelayError> {
        let info = self.read_metadata_store(index);
        if info.digest != raw.digest() {
            return Err(RelayError::WrongDigest);
        }
        Ok(RawWithInfo { raw, info })
    }

    /// Load metadata about a header using its index and 80-bytes raw form. Returns a struct with
    /// the header and a reference to its metadata in the relay store
    pub fn load_header<R: Into<RawHeader>>(
        &self,
        index: u32,
        raw: R,
    ) -> Result<RawWithInfo, RelayError> {
        self.attach_metadata(index, raw.into())
    }

    // Attach a valid chain of headers to the store. Call only after validating the chain.
    fn attach_to_store(&mut self, anchor_index: u32, headers: &HeaderArray) {
        let anchor = self.read_metadata_store(anchor_index);
        // sanity check
        assert!(
            anchor.digest == headers.index(0).parent(),
            "Attempted to attach to wrong anchor"
        );

        // We'll update each of these for each header
        let mut parent_index = anchor_index;
        let mut height = anchor.height + 1;
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

    /// Validate an SPV Proof against the relay.
    ///
    /// # Arguments
    ///
    /// - confirming_header_index - the index of the confirming header in the relay's store.
    /// - proof - the SPVProof to be validated. See the `bitcoin-spv` crate for details.
    pub fn validate_proof(
        &self,
        confirming_header_index: u32,
        proof: &SPVProof,
    ) -> Result<u32, RelayError> {
        proof.validate()?;
        let header_info = self.read_metadata_store(confirming_header_index);
        let chaintip = self.read_metadata_store(self.current_best_index);

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

    fn internal_add_headers<R: Into<RawHeader>>(
        &mut self,
        anchor_index: u32,
        anchor_bytes: R,
        header_bytes: Vec<u8>,
        internal: bool,
    ) -> Result<(), RelayError> {
        let anchor = self.load_header(anchor_index, anchor_bytes)?;
        let headers = HeaderArray::new(&header_bytes).map_err(Into::<RelayError>::into)?;

        let first_new = headers.index(0);
        if !internal && self.mainnet && first_new.target() != anchor.raw.target() {
            return Err(RelayError::UnexpectedDifficultyChange);
        }

        headers
            .valid_difficulty(true)
            .map_err(Into::<RelayError>::into)?;

        // If we haven't errored yet, they're good.
        self.attach_to_store(anchor_index, &headers);
        Ok(())
    }

    /// Add headers to the relay.
    ///
    /// This function validates a chain of header and its connection to some previously-known
    /// header. It checks for adherence to Bitcoin consensus rules, under the SPV assumption.
    /// Once the chain has been validated, it attaches the header info to the metadata store.
    /// Chains are attached to an "anchor", which is an already-stored header. The anchor must
    /// immediately precede the first header in the chain to be ingested.
    ///
    ///
    /// # Arguments
    ///
    /// - anchor_index - the index of the anchor header metadata in the relay store
    /// - anchor_bytes - the raw 80-byte anchor header
    /// - header_bytes - an array of bytes containing any number of tightly-packed headers. Its
    ///     length MUST always be a multiple of 80.
    pub fn add_headers<R: Into<RawHeader>>(
        &mut self,
        anchor_index: u32,
        anchor_bytes: R,
        header_bytes: Vec<u8>,
    ) -> Result<(), RelayError> {
        self.internal_add_headers(anchor_index, anchor_bytes, header_bytes, false)
    }

    /// Add a difficulty change to the relay
    ///
    /// This function validates a chain of header and its connection to some previously-known
    /// header. It checks for adherence to Bitcoin consensus rules, under the SPV assumption.
    /// Once the chain has been validated, it attaches the header info to the metadata store.
    /// Chains are attached to an "anchor", which is an already-stored header. The anchor must
    /// immediately precede the first header in the chain to be ingested.
    ///
    /// Particularly, this function checks application of the retarget algorithm, so it must
    /// accept and validate information about the difficulty epoch that is closing. The anchor
    /// MUST be the last header of the closing difficulty period (`old_period_end`).
    ///
    /// # Arguments
    ///
    /// - old_period_start_bytes - the raw 80-byte header that opened the difficulty epoch
    /// - old_period_end_index - the index of the anchor header metadata in the relay store
    /// - old_period_end_bytes - the raw 80-byte anchor header
    /// - header_bytes - an array of bytes containing any number of tightly-packed headers. Its
    ///     length MUST always be a multiple of 80.
    pub fn add_difficulty_change<R1, R2>(
        &mut self,
        old_period_start_bytes: R1,
        old_period_end_index: u32,
        old_period_end_bytes: R2,
        header_bytes: Vec<u8>, // should be a vec of [u8; 80]
    ) -> Result<(), RelayError>
    where
        R1: Into<RawHeader>,
        R2: Into<RawHeader>,
    {
        let old_period_start_bytes: RawHeader = old_period_start_bytes.into();
        let old_period_end_bytes: RawHeader = old_period_end_bytes.into();

        let headers = HeaderArray::new(&header_bytes)?;
        let old_period_end = self.load_header(old_period_end_index, old_period_end_bytes)?;
        let old_period_start = self.load_header(
            old_period_end.info.epoch_start_index,
            old_period_start_bytes,
        )?;

        // Ensure a change is allowed
        if old_period_end.info.height % 2016 != 2015 {
            return Err(RelayError::WrongEnd);
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
        if (&new_target & expected_target) != new_target {
            return Err(RelayError::IncorrectDifficultyChange);
        }

        // Proceed to add the headers
        self.internal_add_headers(
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
                    left_prev = self.read_metadata_store(left_prev.parent_index);
                }
                if right_prev != ancestor {
                    right_current = right_prev;
                    right_prev = self.read_metadata_store(right_prev.parent_index);
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
    ///
    /// Rather than tracking chain tips as headers are ingested, the relay allows callers to mark
    /// them retrospectively. To do this, we store the current best tip, and then prove that a new
    /// candidate tip is better. The relay can verify this cheaply (arund `O(2n)` where n is the
    /// distance between tip and LCA).
    ///
    /// This function takes the Latest Common Ancestor as an argument. This is the latest header
    /// in the shared history of the new and old tips. Which is to say, the last header that both
    /// tips consider an ancestor.
    ///
    /// # Arguments
    ///
    /// - lca_index - the index of the latest common ancestor in the header store.
    /// - current_best - the 80-byte raw header corresponding to the relay's `best_known_digest`
    /// - new_best_idx - the index of the best known header in the header store
    /// - new_best - the 80-byte raw candidate chain tip
    pub fn mark_new_heaviest<R1, R2>(
        &mut self,
        lca_index: u32,
        current_best: R1,
        new_best_index: u32,
        new_best: R2,
    ) -> Result<(), RelayError>
    where
        R1: Into<RawHeader>,
        R2: Into<RawHeader>,
    {
        // Isolate all the borrows
        let (new_best, ancestor) = {
            let new_best = self.load_header(new_best_index, new_best)?;
            let current_best = self.load_header(self.current_best_index, current_best)?;
            let ancestor = self.read_metadata_store(lca_index);
            Self::verify_better_descendant(&self, &ancestor, &current_best, &new_best)?;
            (new_best.info.digest, ancestor.digest)
        };

        self.last_reorg_lca = ancestor;
        self.best_known_digest = new_best;
        self.current_best_index = new_best_index;
        Ok(())
    }

    /// Find the index of the first header_info that is equal to the value passed in.
    ///
    /// # Note
    ///
    /// This iterates over the internal array of size N. So may be expensive
    pub fn find<A: AsRef<HeaderInfo>>(&self, value: &A) -> Option<usize> {
        self.header_store.find(value)
    }

    /// Find the index of the first header_info that has the digest
    ///
    /// # Note
    ///
    /// This iterates over the internal array of size N. So may be expensive
    pub fn find_digest(&self, value: Hash256Digest) -> Option<usize> {
        self.header_store.position(|v| v.digest == value)
    }
}

#[cfg(test)]
mod test_utils;

#[cfg(test)]
mod test {
    use super::*;
    use crate::test_utils;

    macro_rules! check_err {
        ($err:expr, $expected:expr) => {
            if $err.is_err() {
                let err = $err.unwrap_err();
                let actual_code = test_utils::error_to_code(err);
                if actual_code != $expected {
                    dbg!($expected);
                    assert!(
                        false,
                        "\nExpected error {:?}. \nGot {:?}\n{:?}\n",
                        $expected, actual_code, err
                    );
                }
                continue;
            }
        };
    }

    #[test]
    fn it_ingests_headers() {
        let cases = &test_utils::TEST_VECTORS.header.header_chain_cases;
        for case in cases.iter() {
            let anchor = case.anchor;

            let instantiation_res =
                Relay::new(case.mainnet, anchor.raw, anchor.height, anchor.hash);
            check_err!(instantiation_res, case.output);

            let mut relay = instantiation_res.unwrap();
            let add_res = relay.internal_add_headers(
                relay.find_digest(anchor.hash).unwrap() as u32,
                anchor.raw,
                case.flat_raw_headers(),
                case.internal,
            );
            check_err!(add_res, case.output);
        }
    }

    #[test]
    fn it_ingests_diff_changes() {
        let cases = &test_utils::TEST_VECTORS.header.diff_change_cases;
        for case in cases.iter() {
            if case.rust_skip { continue; }

            let anchor = case.anchor;
            let prev_epoch_start = case.prev_epoch_start;
            let instantiation_res =
                Relay::new(true, anchor.raw, anchor.height, prev_epoch_start.hash);
            check_err!(instantiation_res, case.output);

            let mut relay = instantiation_res.unwrap();
            let add_res = relay.add_difficulty_change(
                prev_epoch_start.raw,
                relay.find_digest(anchor.hash).unwrap() as u32,
                anchor.raw,
                case.flat_raw_headers(),
            );
            check_err!(add_res, case.output);
        }
    }

    #[test]
    fn it_marks_new_heaviest_digests() {
        let chain_block = &test_utils::TEST_VECTORS.chain;
        let setup_info = &chain_block.is_most_recent_common_ancestor;
        let cases = &chain_block.mark_new_heaviest_cases;
        let genesis = &setup_info.genesis;
        let prev_epoch_start = &setup_info.prev_epoch_start;

        let old_epoch_end = setup_info.pre_retarget_chain.last().unwrap();

        let pre_retarget_chain = setup_info.flat_raw_pre();
        let post_retarget_chain = setup_info.flat_raw_post();

        // Creating a chain with an orphan we reorg away from
        let mut post_with_orphan = post_retarget_chain[..post_retarget_chain.len() - 160].to_vec();
        post_with_orphan.extend(&setup_info.orphan.raw[..]);

        let instantiation_res =
            Relay::new(true, genesis.raw, genesis.height, prev_epoch_start.hash);
        let mut relay = instantiation_res.unwrap();

        relay.internal_add_headers(
            relay.find_digest(genesis.hash).unwrap() as u32,
            genesis.raw,
            pre_retarget_chain,
            false,
        ).unwrap();
        relay.add_difficulty_change(
            prev_epoch_start.raw,
            relay.find_digest(old_epoch_end.hash).unwrap() as u32,
            old_epoch_end.raw,
            post_with_orphan,
        ).unwrap();
        relay.add_difficulty_change(
            prev_epoch_start.raw,
            relay.find_digest(old_epoch_end.hash).unwrap() as u32,
            old_epoch_end.raw,
            post_retarget_chain,
        ).unwrap();

        // Update heaviest to the first block we ingested
        let new_heaviest = setup_info.pre_retarget_chain[0];
        relay.mark_new_heaviest(
            relay.find_digest(genesis.hash).unwrap() as u32,
            genesis.raw,
            relay.find_digest(new_heaviest.hash).unwrap() as u32,
            new_heaviest.raw,
        ).unwrap();

        let old_heaviest = new_heaviest;
        let new_heaviest = setup_info.pre_retarget_chain[1];
        // NotHeaviestAncestor
        assert!(relay.mark_new_heaviest(
            relay.find_digest(genesis.hash).unwrap() as u32,
            old_heaviest.raw,
            relay.find_digest(new_heaviest.hash).unwrap() as u32,
            new_heaviest.raw,
        ).is_err());

        for case in cases.iter() {
            dbg!(case);

            let idx_opt = relay.find_digest(case.best_known_digest);
            if idx_opt.is_none() { continue; }

            relay.current_best_index = idx_opt.unwrap() as u32;
            relay.best_known_digest = case.best_known_digest;

            let new_best_idx_opt = relay.find_digest(case.new_best.digest());
            if new_best_idx_opt.is_none() { continue; }
            let new_best_idx = new_best_idx_opt.unwrap() as u32;

            let result = relay.mark_new_heaviest(
                relay.find_digest(case.ancestor).unwrap() as u32,
                case.current_best,
                new_best_idx,
                case.new_best,
            );
            check_err!(result, case.error);
        }
    }
}
