package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getRequestStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.RequestStorePrefix)
}

func (k Keeper) hasRequest(ctx sdk.Context, id []byte) bool {
	store := k.getRequestStore(ctx)
	return store.Has(id)
}

func (k Keeper) setRequest(ctx sdk.Context, request types.ProofRequest) {
	store := k.getRequestStore(ctx)
	buf, _ := json.Marshal(request)
	// When a new request comes in, get the id and use it to store request
	id := k.getID(ctx)
	store.Set(id, buf)
	// Increment the ID
	k.incrementID(ctx)
}

func (k Keeper) getRequest(ctx sdk.Context, id []byte) types.ProofRequest {
	store := k.getRequestStore(ctx)
	buf := store.Get(id)
	var request []types.ProofRequest
	json.Unmarshal(buf, &request)
	return request[0]
}

func (k Keeper) incrementID(ctx sdk.Context) {
	// TODO: Is there a better way of incrementing this? Have to store as bytes...
	store := k.getRequestStore(ctx)
	// get id
	id := k.getID(ctx)
	// convert id to uint64 and add 1
	newID := binary.BigEndian.Uint64(id) + 1
	// convert back to bytes and store
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, newID)
	store.Set([]byte("id"), b)
}

func (k Keeper) getID(ctx sdk.Context) []byte {
	store := k.getRequestStore(ctx)
	id := []byte("id")
	if !store.Has(id) {
		store.Set(id, []byte{0})
	}
	return store.Get(id)
}

func (k Keeper) checkRequests(ctx sdk.Context, reqIndices uint16, vin []byte, vout []byte, requestID []byte) bool {
	// TODO: Add errors
	if !btcspv.ValidateVin(vin) {
		return false
	}
	if !btcspv.ValidateVout(vout) {
		return false
	}

	inputIndex := uint8(reqIndices >> 8)
	outputIndex := uint8(reqIndices & 0xff)

	id := []byte{0, 0, 0}
	req := k.getRequest(ctx, id)

	if !req.ActiveState {
		return false
	}

	hasPays := !bytes.Equal(req.Pays[:], bytes.Repeat([]byte{0}, 32))
	if hasPays {
		out, _ := btcspv.ExtractOutputAtIndex(vout, outputIndex)
		len := btcspv.ExtractOutputScriptLen(out)
		if !bytes.Equal(out[8:len+1], req.Pays[:]) {
			return false
		}
		paysValue := req.PaysValue
		if paysValue != 0 || uint64(btcspv.ExtractValue(out)) <= paysValue {
			return false
		}
	}

	hasSpends := !bytes.Equal(req.Spends[:], bytes.Repeat([]byte{0}, 32))
	if hasSpends {
		in := btcspv.ExtractInputAtIndex(vin, inputIndex)
		if !hasSpends || !bytes.Equal(btcspv.ExtractOutpoint(in), req.Spends[:]) {
			return false
		}
	}
	return true
}
