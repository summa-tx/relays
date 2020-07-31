use generic_array::{ArrayLength, GenericArray};

/// A simple Fake Vector. It has a fixed-sized allocation and overwrites things at the front when
/// it runs out of space.
///
/// It can be used at most `usize::MAX` times. It panics after that. This shouldn't be an issue.
#[repr(C)]
#[derive(Debug, Default, Clone, PartialEq, Eq, serde::Serialize, serde::Deserialize)]
#[serde(bound = "N: ArrayLength<T>")]
pub struct FakeVec<T, N>
where
    T: serde::Serialize + for<'d> serde::Deserialize<'d> + Default + Clone,
    N: ArrayLength<T>,
{
    pub(crate) next: usize,
    pub(crate) internal: GenericArray<T, N>,
}

impl<T, N> FakeVec<T, N>
where
    T: serde::Serialize + for<'d> serde::Deserialize<'d> + Default + Clone,
    N: ArrayLength<T>,
{
    // Panics if called when empty. TODO maybe fix that
    fn latest(&self) -> usize {
        self.next - 1
    }

    // Get the capacity of the underlying array
    fn capacity(&self) -> usize {
        self.internal.len()
    }

    // Get the actual index within the internal array
    fn to_internal_index(&self, index: usize) -> usize {
        index % self.capacity()
    }

    // Transform an internal array index to an external retrieval
    fn to_external_index(&self, index: usize) -> usize {
        assert!(index < self.capacity());

        // lazy algo. there's probably a better way
        let mut internal_latest = self.to_internal_index(self.latest());
        if index > internal_latest {
            internal_latest += self.capacity();
        }
        self.latest() - (internal_latest - index)
    }

    // Check if we have the item referenced by an index
    fn valid(&self, index: usize) -> bool {
        index <= self.latest()
            && (self.latest() < self.capacity() || index > self.latest() - self.capacity())
    }

    /// Push an item to the Buffer
    pub fn push(&mut self, item: T) -> usize {
        let index = self.to_internal_index(self.next);
        self.internal[index] = item;
        self.next += 1;
        index
    }

    /// Get an item at an index, if we still have it
    pub fn get(&self, index: usize) -> Option<&T> {
        if !self.valid(index) {
            None
        } else {
            let index = self.to_internal_index(index);
            Some(&self.internal[index])
        }
    }

    /// Find the index of the first item in the array that matches the filter.
    ///
    /// # Notes
    ///
    /// - This iterates over the internal array of size N. So may be expensive
    /// - This always returns a valid index. However, that index may become invalid as new things
    ///   are pushed to the array
    pub fn position<F>(&self, func: F) -> Option<usize>
    where
        F: FnMut(&T) -> bool,
    {
        self.internal
            .iter()
            .position(func)
            .map(|v| self.to_external_index(v))
    }

    /// Find the index of the first item in the array that is equal to the value passed in.
    ///
    /// # Note
    ///
    /// This iterates over the internal array of size N. So may be expensive
    pub fn find<A: AsRef<T>>(&self, value: &A) -> Option<usize>
    where
        T: PartialEq,
    {
        self.position(|v| v == value.as_ref())
    }
}

#[cfg(test)]
mod test {
    use super::*;

    // #[test]
    // fn it_defaults() {
    //     dbg!(std::mem::size_of::<GenericArray<[u64; 32], U816>>());
    //     dbg!(std::mem::size_of::<GenericArray<[u64; 32], U817>>());
    //
    //     // Works
    //     GenericArray::<[u64; 32], U816>::default();
    //
    //     // Stack Overflow
    //     GenericArray::<[u64; 32], U817>::default();
    // }

    #[test]
    fn it_generates_external_indices() {
        let mut store = FakeVec::<u8, generic_array::typenum::U8> {
            next: 0,
            internal: Default::default(),
        };
        store.next = 502;

        assert_eq!(store.to_external_index(5), store.latest());
        assert_eq!(store.to_external_index(6), store.latest() - 7);
        assert_eq!(store.to_external_index(7), store.latest() - 6);
        assert_eq!(store.to_external_index(0), store.latest() - 5);
        assert_eq!(store.to_external_index(1), store.latest() - 4);
        assert_eq!(store.to_external_index(2), store.latest() - 3);
        assert_eq!(store.to_external_index(3), store.latest() - 2);
        assert_eq!(store.to_external_index(4), store.latest() - 1);
    }
}
