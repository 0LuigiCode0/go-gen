package tmp

const HandlerTCPTmp = `package {{printf "%v_handler" (index . 0)}}

import (
	"encoding/json"
	"net/http"
	{{printf "\"%v/handlers/%v_handler/%v_helper\"" (index . 1) (index . 0) (index . 0)}}
	{{printf "\"%v/helper\"" (index . 1)}}
	{{printf "\"%v/hub/hub_helper\"" (index . 1)}}
)

type handler struct {
	hub_helper.HelperForHandler
}

func InitHandler(hub hub_helper.HelperForHandler, conf *helper.HandlerConfig) (H {{print (index . 0)}}_helper.Handler, err error) {
	h := &handler{HelperForHandler: hub}
	H = h

	h.Router().Use(h.middleware)
	h.SetHandler(applyCORS(h.Router()))
	return
}

func (h *handler) respOk(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := &helper.ResponseModel{
		Success: true,
		Result:  data,
	}
	buf, err := json.Marshal(resp)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorParse+": josn: %v", err)
		h.respError(w, helper.ErrorParse, helper.KeyErrorParse+": josn")
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorWrite+": response: %v", err)
		h.respError(w, helper.ErrorWrite, helper.KeyErrorWrite+": response")
		return
	}
}

func (h *handler) respError(w http.ResponseWriter, code helper.ErrCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	resp := &helper.ResponseModel{
		Success: false,
		Result: &helper.ResponseError{
			Code: code,
			Msg:  msg,
		},
	}
	buf, _ := json.Marshal(resp)
	_, err := w.Write(buf)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorWrite+": response: %v", err)
	}
}`
