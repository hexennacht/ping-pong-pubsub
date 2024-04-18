package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hexennacht/ping-pong-pubsub/helper"
	"github.com/hexennacht/ping-pong-pubsub/module/entity"
	"github.com/hexennacht/ping-pong-pubsub/module/service"
)

type pingHandler struct {
	svc service.PingService
}

func RegisterPingHandler(r *chi.Mux, svc service.PingService) {
	ph := &pingHandler{
		svc: svc,
	}

	r.Post("/ping", ph.Ping)
}

func (ph *pingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	req, err := helper.ReadJsonBody[entity.PingRequest](r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		response, _ := json.Marshal(helper.ResponseBody{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Errors:  err.Error(),
		})

		w.Write(response)
		return
	}

	result, err := ph.svc.Ping(r.Context(), req)
	if err != nil {
		var responseError *helper.ResponseBody
		errors.As(err, &responseError)
		response, _ := json.Marshal(responseError)
		w.WriteHeader(responseError.Code)
		w.Write(response)
		return
	}

	response, _ := json.Marshal(helper.ResponseBody{
		Code:    http.StatusCreated,
		Message: "SUCCESS",
		Data:    result,
	})

	w.WriteHeader(http.StatusAccepted)
	w.Write(response)
}
