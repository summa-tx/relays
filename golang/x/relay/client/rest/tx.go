package rest

//
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
