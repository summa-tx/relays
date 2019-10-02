package keeper

import (
	"bytes"
	"encoding/hex"

	"github.com/summa-tx/relays/golang/x/relay/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryGetParent is a query string tag for get parent
const QueryGetParent = "getparent"

// NewQuerier makes a query routing function
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGetParent:
			return queryGetParent(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryGetParent(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	digest, err := hex.DecodeString(path[0])
	if err != nil || len(digest) != 32 {
		panic("could not unhex argument string")
	}

	var buf [32]byte
	copy(buf[:], digest[:32])
	parent := keeper.GetLink(ctx, buf)

	if bytes.Equal(parent[:], []byte{}) {
		return []byte{}, types.ErrUnknownBlock(types.DefaultCodespace)
	}

	response := types.QueryResGetParent{
		Digest: hex.EncodeToString(digest),
		Parent: hex.EncodeToString(parent[:])}
	res, err := codec.MarshalJSONIndent(keeper.cdc, response)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return res, nil
}
