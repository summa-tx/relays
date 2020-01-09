package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) emitProofRequest(ctx sdk.Context, pays, spends types.Hash256Digest, paysValue, id uint64) {
	ctx.EventManager().EmitEvent(types.NewProofRequestEvent(pays, spends, paysValue, id))
}

func (k Keeper) getRequestStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.RequestStorePrefix)
}

func (k Keeper) hasRequest(ctx sdk.Context, id []byte) bool {
	store := k.getRequestStore(ctx)
	return store.Has(id)
}

func (k Keeper) setRequest(ctx sdk.Context, spends []byte, pays []byte, paysValue uint64, numConfs uint8) sdk.Error {
	store := k.getRequestStore(ctx)

	valid := k.validateRequests(spends, pays)
	if !valid {
		return types.ErrInvalidRequest(types.DefaultCodespace)
	}

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
	id := k.getID(ctx)

	buf, err := json.Marshal(request)
	if err != nil {
		return types.ErrMarshalJSON(types.DefaultCodespace)
	}
	store.Set(id, buf)

	// Increment the ID
	k.incrementID(ctx)

	// Emit Proof Request event
	numID := binary.BigEndian.Uint64(id)
	k.emitProofRequest(ctx, request.Pays, request.Spends, request.PaysValue, numID)
	return nil
}

func (k Keeper) getRequest(ctx sdk.Context, id []byte) types.ProofRequest {
	store := k.getRequestStore(ctx)
	buf := store.Get(id)
	var request []types.ProofRequest
	json.Unmarshal(buf, &request)
	return request[0]
}

func (k Keeper) incrementID(ctx sdk.Context) {
	store := k.getRequestStore(ctx)
	// get id
	id := k.getID(ctx)
	// convert id to uint64 and add 1
	newID := binary.BigEndian.Uint64(id) + 1
	// convert back to bytes and store
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, newID)
	store.Set([]byte(types.RequestID), b)
}

func (k Keeper) getID(ctx sdk.Context) []byte {
	store := k.getRequestStore(ctx)
	id := []byte(types.RequestID)
	if !store.Has(id) {
		store.Set(id, []byte{0})
	}
	return store.Get(id)
}

func (k Keeper) validateRequests(spends []byte, pays []byte) bool {
	if len(spends) != 36 {
		return false
	}
	if len(pays) > 50 {
		return false
	}
	return true
}

func (k Keeper) checkRequests(ctx sdk.Context, inputIndex, outputIndex uint8, vin []byte, vout []byte, requestID []byte) (bool, sdk.Error) {
	// TODO: Add errors
	if !btcspv.ValidateVin(vin) {
		return false, types.ErrInvalidVin(types.DefaultCodespace)
	}
	if !btcspv.ValidateVout(vout) {
		return false, types.ErrInvalidVout(types.DefaultCodespace)
	}

	req := k.getRequest(ctx, requestID)
	if !req.ActiveState {
		return false, types.ErrClosedRequest(types.DefaultCodespace)
	}

	hasPays := !bytes.Equal(req.Pays[:], bytes.Repeat([]byte{0}, 32))
	if hasPays {
		out, _ := btcspv.ExtractOutputAtIndex(vout, outputIndex)
		len := btcspv.ExtractOutputScriptLen(out)
		if !bytes.Equal(out[8:len+1], req.Pays[:]) {
			return false, types.ErrRequestPays(types.DefaultCodespace)
		}
		paysValue := req.PaysValue
		if paysValue != 0 || uint64(btcspv.ExtractValue(out)) <= paysValue {
			return false, types.ErrRequestValue(types.DefaultCodespace)
		}
	}

	hasSpends := !bytes.Equal(req.Spends[:], bytes.Repeat([]byte{0}, 32))
	if hasSpends {
		in := btcspv.ExtractInputAtIndex(vin, inputIndex)
		if !hasSpends || !bytes.Equal(btcspv.ExtractOutpoint(in), req.Spends[:]) {
			return false, types.ErrRequestSpends(types.DefaultCodespace)
		}
	}
	return true, nil
}
