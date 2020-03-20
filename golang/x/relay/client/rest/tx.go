package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

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

type IngestDifficultyChangeReq struct {
	BaseReq rest.BaseReq          `json:"base_req"`
	Start   types.Hash256Digest   `json:"prevEpochStart"`
	Headers []types.BitcoinHeader `json:"headers"`
	Sender  string                `json:"sender"`
}

func ingestDifficultyChangeHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IngestDifficultyChangeReq

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

		msg := types.NewMsgIngestDifficultyChange(addr, req.Start, req.Headers)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type MarkNewHeaviestReq struct {
	BaseReq     rest.BaseReq        `json:"base_req"`
	Ancestor    types.Hash256Digest `json:"ancestor"`
	CurrentBest types.RawHeader     `json:"currentBest"`
	NewBest     types.RawHeader     `json:"newBest"`
	Limit       uint32              `json:"limit"`
	Sender      string              `json:"sender"`
}

func markNewHeaviestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MarkNewHeaviestReq

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

		msg := types.NewMsgMarkNewHeaviest(addr, req.Ancestor, req.CurrentBest, req.NewBest, req.Limit)
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

		msg := types.NewMsgNewRequest(addr, req.Spends, req.Pays, req.PaysValue, req.NumConfs, types.Local, nil)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type ProvideProofReq struct {
	BaseReq  rest.BaseReq              `json:"base_req"`
	Proof    types.SPVProof            `json:"proof"`
	Requests []types.FilledRequestInfo `json:"filled_requests"`
	Sender   string                    `json:"sender"`
}

func provideProofHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProvideProofReq

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

		filledRequests := types.NewFilledRequests(req.Proof, req.Requests)

		msg := types.NewMsgProvideProof(addr, filledRequests)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
