package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) emitProofRequest(ctx sdk.Context, pays, spends []byte, paysValue uint64, id types.RequestID) {
	ctx.EventManager().EmitEvent(types.NewProofRequestEvent(pays, spends, paysValue, id))
}

func (k Keeper) getRequestStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.RequestStorePrefix)
}

func (k Keeper) hasRequest(ctx sdk.Context, id types.RequestID) bool {
	store := k.getRequestStore(ctx)
	return store.Has(id[:])
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
	id, err := k.getNextID(ctx)
	if err != nil {
		return err
	}

	buf, marshalErr := json.Marshal(request)
	if marshalErr != nil {
		return types.ErrMarshalJSON(types.DefaultCodespace)
	}
	store.Set(id[:], buf)

	// Increment the ID
	incrementErr := k.incrementID(ctx)
	if incrementErr != nil {
		return incrementErr
	}

	// Emit Proof Request event
	k.emitProofRequest(ctx, pays, spends, request.PaysValue, id)
	return nil
}

func (k Keeper) setRequestState(ctx sdk.Context, requestID types.RequestID, active bool) sdk.Error {
	store := k.getRequestStore(ctx)
	request, err := k.getRequest(ctx, requestID)
	if err != nil {
		return err
	}

	request.ActiveState = active

	buf, marshalErr := json.Marshal(request)
	if marshalErr != nil {
		return types.ErrMarshalJSON(types.DefaultCodespace)
	}
	store.Set(requestID[:], buf)
	return nil
}

func (k Keeper) getRequest(ctx sdk.Context, id types.RequestID) (types.ProofRequest, sdk.Error) {
	store := k.getRequestStore(ctx)

	hasRequest := k.hasRequest(ctx, id)
	if !hasRequest {
		return types.ProofRequest{}, types.ErrUnknownRequest(types.DefaultCodespace)
	}

	buf := store.Get(id[:])

	var request types.ProofRequest
	jsonErr := json.Unmarshal(buf, &request)
	if jsonErr != nil {
		return types.ProofRequest{}, types.ErrExternal(types.DefaultCodespace, jsonErr)
	}
	return request, nil
}

func (k Keeper) incrementID(ctx sdk.Context) sdk.Error {
	store := k.getRequestStore(ctx)
	// get id
	id, err := k.getNextID(ctx)
	if err != nil {
		return err
	}
	// convert id to uint64 and add 1
	newID := binary.BigEndian.Uint64(id[:]) + 1
	// convert back to bytes and store
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, newID)
	store.Set([]byte(types.RequestIDTag), b)
	// if no errors, return nil
	return nil
}

// getNextID retrieves the ID.  The ID is incremented after storing a request,
// so this returns the next ID to be used.
func (k Keeper) getNextID(ctx sdk.Context) (types.RequestID, sdk.Error) {
	store := k.getRequestStore(ctx)
	idTag := []byte(types.RequestIDTag)
	if !store.Has(idTag) {
		store.Set(idTag, bytes.Repeat([]byte{0}, 8))
	}
	id := store.Get(idTag)
	newID, err := types.NewRequestID(id)
	if err != nil {
		return types.RequestID{}, err
	}
	return newID, nil
}

func (k Keeper) checkRequests(ctx sdk.Context, inputIndex, outputIndex uint8, vin []byte, vout []byte, requestID types.RequestID) sdk.Error {
	if !btcspv.ValidateVin(vin) {
		return types.ErrInvalidVin(types.DefaultCodespace)
	}
	if !btcspv.ValidateVout(vout) {
		return types.ErrInvalidVout(types.DefaultCodespace)
	}

	req, reqErr := k.getRequest(ctx, requestID)
	if reqErr != nil {
		return reqErr
	}
	if !req.ActiveState {
		return types.ErrClosedRequest(types.DefaultCodespace)
	}

	hasPays := req.Pays != btcspv.Hash256([]byte{0})
	if hasPays {
		// We can ignore this error because we know that ValidateVout passed
		out, _ := btcspv.ExtractOutputAtIndex(vout, outputIndex)
		// hash the output script (out[8:])
		outDigest := btcspv.Hash256(out[8:])
		if outDigest != req.Pays {
			return types.ErrRequestPays(types.DefaultCodespace)
		}
		paysValue := req.PaysValue
		if paysValue != 0 && uint64(btcspv.ExtractValue(out)) < paysValue {
			return types.ErrRequestValue(types.DefaultCodespace)
		}
	}

	hasSpends := req.Spends != btcspv.Hash256([]byte{0})
	if hasSpends {
		in := btcspv.ExtractInputAtIndex(vin, inputIndex)
		outpoint := btcspv.ExtractOutpoint(in)
		inDigest := btcspv.Hash256(outpoint)
		if hasSpends && inDigest != req.Spends {
			return types.ErrRequestSpends(types.DefaultCodespace)
		}
	}
	return nil
}
