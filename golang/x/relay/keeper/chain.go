package keeper

import (
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

func (k Keeper) getChainStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.ChainStorePrefix)
}

func (k Keeper) emitReorg(ctx sdk.Context, prev, new, lca types.Hash256Digest) {
	ctx.EventManager().EmitEvent(types.NewReorgEvent(prev, new, lca))
}

func (k Keeper) getDigestByStoreKey(ctx sdk.Context, key string) (types.Hash256Digest, sdk.Error) {
	store := k.getChainStore(ctx)
	result := store.Get([]byte(key))

	digest, err := btcspv.NewHash256Digest(result)
	if err != nil {
		return types.Hash256Digest{}, types.ErrBadHash256Digest(types.DefaultCodespace)
	}
	return digest, nil
}

func (k Keeper) setDigestByStoreKey(ctx sdk.Context, key string, digest types.Hash256Digest) {
	// TODO: Remove this in favor of Genesis state
	store := k.getChainStore(ctx)
	store.Set([]byte(key), digest[:])
}

// setBestKnownDigest sets the best known chain tip
func (k Keeper) setBestKnownDigest(ctx sdk.Context, bestKnown types.Hash256Digest) {
	k.setDigestByStoreKey(ctx, types.BestKnownDigestStorage, bestKnown)
}

// GetBestKnownDigest returns the best known digest in the relay
func (k Keeper) GetBestKnownDigest(ctx sdk.Context) (types.Hash256Digest, sdk.Error) {
	return k.getDigestByStoreKey(ctx, types.BestKnownDigestStorage)
}

// setLastReorgLCA sets the latest common ancestor of the last reorg
func (k Keeper) setLastReorgLCA(ctx sdk.Context, lca types.Hash256Digest) {
	k.setDigestByStoreKey(ctx, types.LastReorgLCAStorage, lca)
}

// GetLastReorgLCA returns the best known digest in the relay
func (k Keeper) GetLastReorgLCA(ctx sdk.Context) (types.Hash256Digest, sdk.Error) {
	return k.getDigestByStoreKey(ctx, types.LastReorgLCAStorage)
}

// IsMostRecentCommonAncestor checks if a proposed ancestor is the LCA of two digests
func (k Keeper) IsMostRecentCommonAncestor(ctx sdk.Context, ancestor, left, right types.Hash256Digest, limit uint32) bool {
	if ancestor == left && ancestor == right {
		return true
	}

	leftCurrent := left
	leftPrev := left

	rightCurrent := right
	rightPrev := right

	for i := uint32(0); i < limit; i++ {
		if leftPrev != ancestor {
			leftCurrent = leftPrev
			leftPrev = k.getLink(ctx, leftPrev)
		}
		if rightPrev != ancestor {
			rightCurrent = rightPrev
			rightPrev = k.getLink(ctx, rightPrev)
		}
		if leftPrev == rightPrev {
			break
		}
	}

	if leftCurrent == rightCurrent {
		return false
	}

	if leftPrev != rightPrev {
		return false
	}

	return true
}

// HeaviestFromAncestor determines the heavier descendant of a common ancestor
func (k Keeper) HeaviestFromAncestor(ctx sdk.Context, ancestor, currentBest, newBest types.Hash256Digest, limit uint32) (types.Hash256Digest, sdk.Error) {
	ancestorBlock, err := k.GetHeader(ctx, ancestor)
	if err != nil {
		return types.Hash256Digest{}, types.ErrUnknownBlock(types.DefaultCodespace)
	}
	leftBlock, err := k.GetHeader(ctx, currentBest)
	if err != nil {
		return types.Hash256Digest{}, types.ErrUnknownBlock(types.DefaultCodespace)
	}
	rightBlock, err := k.GetHeader(ctx, newBest)
	if err != nil {
		return types.Hash256Digest{}, types.ErrUnknownBlock(types.DefaultCodespace)
	}

	if leftBlock.Height < ancestorBlock.Height || rightBlock.Height < ancestorBlock.Height {
		return types.Hash256Digest{}, types.ErrBadHeight(types.DefaultCodespace)
	}

	nextPeriodStartHeight := ancestorBlock.Height + 2016 - (ancestorBlock.Height % 2016)
	leftInPeriod := leftBlock.Height < nextPeriodStartHeight
	rightInPeriod := rightBlock.Height < nextPeriodStartHeight

	/*
		NB:
		1. Left is in a new window, right is in the old window. Left is heavier
		2. Right is in a new window, left is in the old window. Right is heavier
		3. Both are in the same window, choose the higher one
		4. They're in different new windows. Choose the heavier one
	*/

	if !leftInPeriod && rightInPeriod {
		return leftBlock.HashLE, nil
	}
	if leftInPeriod && !rightInPeriod {
		return rightBlock.HashLE, nil
	}
	if leftInPeriod && rightInPeriod {
		if leftBlock.Height >= rightBlock.Height {
			return leftBlock.HashLE, nil
		}
		return rightBlock.HashLE, nil
	}

	// if !leftInPeriod && !rightInPeriod
	leftDiff := btcspv.ExtractDifficulty(leftBlock.Raw)
	leftAccDiff := leftDiff.Mul(sdk.NewUint(uint64(leftBlock.Height % 2016)))

	rightDiff := btcspv.ExtractDifficulty(rightBlock.Raw)
	rightAccDiff := rightDiff.Mul(sdk.NewUint(uint64(rightBlock.Height % 2016)))

	if leftAccDiff.GTE(rightAccDiff) {
		return leftBlock.HashLE, nil
	}
	return rightBlock.HashLE, nil
}

// MarkNewHeaviest updates the best known digest and LCA
func (k Keeper) MarkNewHeaviest(ctx sdk.Context, ancestor types.Hash256Digest, currentBest, newBest types.RawHeader, limit uint32) sdk.Error {
	newBestDigest := btcspv.Hash256(newBest[:])
	currentBestDigest := btcspv.Hash256(currentBest[:])

	knownBestDigest, err := k.GetBestKnownDigest(ctx)
	if err != nil || currentBestDigest != knownBestDigest {
		return types.ErrNotBestKnown(types.DefaultCodespace)
	}

	if !k.HasHeader(ctx, newBestDigest) {
		return types.ErrUnknownBlock(types.DefaultCodespace)
	}

	if !k.IsMostRecentCommonAncestor(ctx, ancestor, knownBestDigest, newBestDigest, limit) {
		return types.ErrNotHeaviestAncestor(types.DefaultCodespace)
	}

	better, err := k.HeaviestFromAncestor(ctx, ancestor, knownBestDigest, newBestDigest, limit)
	if err != nil {
		return err
	}

	if newBestDigest != better {
		return types.ErrNotHeavier(types.DefaultCodespace)
	}

	k.setLastReorgLCA(ctx, ancestor)
	k.setBestKnownDigest(ctx, newBestDigest)
	k.emitReorg(ctx, knownBestDigest, newBestDigest, ancestor)

	return nil
}
