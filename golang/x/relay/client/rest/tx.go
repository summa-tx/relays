package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

// type setLinkReq struct {
// 	BaseReq rest.BaseReq   `json:"base_req"`
// 	Header  string         `json:"header"`
// 	Sender  sdk.AccAddress `json:"sender"`
// }
//
// func setLinkHandler(cliCtx context.CLIContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req setLinkReq
//
// 		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
// 			return
// 		}
//
// 		baseReq := req.BaseReq.Sanitize()
// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}
//
// 		addr := sdk.AccAddress(req.Sender)
//
// 		// create the message
// 		msg := types.NewMsgSetLink(req.Sender, req.Header, addr)
// 		err := msg.ValidateBasic()
//
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}
//
// 		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
// 	}
// }

type IngestHeaderChainReq struct {
	BaseReq rest.BaseReq          `json:"base_req"`
	Headers []types.BitcoinHeader `json:"headers"`
	Sender  string                `json:"sender"`
}

func ingestHeaderChainHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IngestHeaderChainReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Sender)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgIngestHeaderChain(addr, req.Headers)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type NewRequestReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Spends    []byte       `json:"spends"`
	Pays      []byte       `json:"pays"`
	PaysValue uint64       `json:"paysValue"`
	NumConfs  uint8        `json:"numConfs"`
	Sender    string       `json:"sender"`
}

func newRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req NewRequestReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Sender)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgNewRequest(addr, req.Spends, req.Pays, req.PaysValue, req.NumConfs)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
