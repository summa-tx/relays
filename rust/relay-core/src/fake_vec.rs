
use generic_array::{ArrayLength, GenericArray};

/// A simple Fake Vector. It has a fixed-sized allocation and overwrites things at the front when
/// it runs out of space.
///
/// It can be used at most `usize::MAX` times. It panics after that. This shouldn't be an issue.
#[repr(C)]
#[derive(Debug, Clone, Default, PartialEq, Eq, serde::Serialize, serde::Deserialize)]
#[serde(bound = "N: ArrayLength<T>")]
pub struct FakeVec<T, N>
where
    T: serde::Serialize + for<'d> serde::Deserialize<'d> + Default,
    N: ArrayLength<T>,
{
    latest: usize,
    internal: GenericArray<T, N>,
}

impl<T, N> FakeVec<T, N>
where
    T: serde::Serialize + for<'d> serde::Deserialize<'d> + Default,
    N: ArrayLength<T>,
{
    // Get the capacity of the underlying array
    fn capacity(&self) -> usize {
        self.internal.len()
    }

    // Get the actual index within the array
    fn normalize_index(&self, index: usize) -> usize {
        index % self.capacity()
    }

    // Check if we have the item referenced by an index
    fn valid(&self, index: usize) -> bool {
        index <= self.latest && index > self.latest - self.capacity()
    }

    /// Push an item to the Buffer
    pub fn push(&mut self, item: T) -> usize {
        self.latest += 1;
        let index = self.normalize_index(self.latest);
        self.internal[index] = item;
        self.latest
    }

    /// Get an item at an index, if we still have it
    pub fn get(&self, index: usize) -> Option<&T> {
        if !self.valid(index) {
            None
        } else {
            let index = self.normalize_index(index);
            Some(&self.internal[index])
        }
    }
}
