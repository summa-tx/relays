package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) emitProofRequest(ctx sdk.Context, pays, spends []byte, paysValue, id uint64) {
	ctx.EventManager().EmitEvent(types.NewProofRequestEvent(pays, spends, paysValue, id))
}

func (k Keeper) getRequestStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.RequestStorePrefix)
}

func (k Keeper) hasRequest(ctx sdk.Context, id uint64) bool {
	// convert id to bytes
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, id)

	store := k.getRequestStore(ctx)
	return store.Has(idBytes)
}

func (k Keeper) setRequest(ctx sdk.Context, spends []byte, pays []byte, paysValue uint64, numConfs uint8) sdk.Error {
	store := k.getRequestStore(ctx)

	spendsDigest := btcspv.Hash256(spends)
	paysDigest := btcspv.Hash256(pays)

	request := types.ProofRequest{
		Spends:      spendsDigest,
		Pays:        paysDigest,
		PaysValue:   paysValue,
		ActiveState: true,
		NumConfs:    numConfs,
	}

	// When a new request comes in, get the id and use it to store request
	id := k.getNextID(ctx)

	buf, err := json.Marshal(request)
	if err != nil {
		return types.ErrMarshalJSON(types.DefaultCodespace)
	}
	store.Set(id, buf)

	// Increment the ID
	k.incrementID(ctx)

	// Emit Proof Request event
	numID := binary.BigEndian.Uint64(id)
	k.emitProofRequest(ctx, pays, spends, request.PaysValue, numID)
	return nil
}

func (k Keeper) getRequest(ctx sdk.Context, id uint64) (types.ProofRequest, sdk.Error) {
	store := k.getRequestStore(ctx)

	// convert id to bytes
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, id)

	hasRequest := k.hasRequest(ctx, id)
	if !hasRequest {
		return types.ProofRequest{}, types.ErrUnknownRequest(types.DefaultCodespace)
	}
	buf := store.Get(idBytes)
	var request types.ProofRequest
	jsonErr := json.Unmarshal(buf, &request)
	if jsonErr != nil {
		return types.ProofRequest{}, types.ErrExternal(types.DefaultCodespace, jsonErr)
	}
	return request, nil
}

func (k Keeper) incrementID(ctx sdk.Context) {
	store := k.getRequestStore(ctx)
	// get id
	id := k.getNextID(ctx)
	// convert id to uint64
	newID := binary.BigEndian.Uint64(id)
	// add 1, convert back to bytes, and store
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, newID+1)
	store.Set([]byte(types.RequestID), b)
}

// getNextID retrieves the ID.  The ID is incremented after storing a request,
// so this returns the next ID to be used.
func (k Keeper) getNextID(ctx sdk.Context) []byte {
	store := k.getRequestStore(ctx)
	id := []byte(types.RequestID)
	if !store.Has(id) {
		store.Set(id, bytes.Repeat([]byte{0}, 8))
	}
	return store.Get(id)
}

func (k Keeper) checkRequests(ctx sdk.Context, inputIndex, outputIndex uint8, vin, vout []byte, requestID uint64) (bool, sdk.Error) {
	if !btcspv.ValidateVin(vin) {
		return false, types.ErrInvalidVin(types.DefaultCodespace)
	}
	if !btcspv.ValidateVout(vout) {
		return false, types.ErrInvalidVout(types.DefaultCodespace)
	}

	req, reqErr := k.getRequest(ctx, requestID)
	if reqErr != nil {
		return false, reqErr
	}
	if !req.ActiveState {
		return false, types.ErrClosedRequest(types.DefaultCodespace)
	}

	hasPays := req.Pays != types.Hash256Digest{}
	if hasPays {
		// We can ignore this error because we know that ValidateVout passed
		out, _ := btcspv.ExtractOutputAtIndex(vout, outputIndex)
		outDigest := btcspv.Hash256(out[8:])
		if outDigest != req.Pays {
			return false, types.ErrRequestPays(types.DefaultCodespace)
		}
		paysValue := req.PaysValue
		if paysValue != 0 || uint64(btcspv.ExtractValue(out)) < paysValue {
			return false, types.ErrRequestValue(types.DefaultCodespace)
		}
	}

	hasSpends := req.Spends != types.Hash256Digest{}
	if hasSpends {
		in := btcspv.ExtractInputAtIndex(vin, inputIndex)
		inDigest := btcspv.Hash256(in)
		if !hasSpends || inDigest != req.Spends {
			return false, types.ErrRequestSpends(types.DefaultCodespace)
		}
	}
	return true, nil
}
