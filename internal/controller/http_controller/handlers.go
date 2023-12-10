package http_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/horockey/distributed_kv_storage/internal/controller/http_controller/dto"
	"github.com/horockey/go-toolbox/http_helpers"
)

func (ctrl *httpController) kvGetKey(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key, found := vars["key"]
	if !found {
		http_helpers.RespondWithErr(
			w,
			http.StatusBadRequest,
			errors.New("missing key"),
		)
		return
	}

	val, err := ctrl.uc.Get(key)
	if err != nil {
		http_helpers.RespondWithErr(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("getting KV from usecase: %w", err),
		)
		return
	}

	http_helpers.RespondOK(w, dto.KV{Key: key, Value: val})
}

func (ctrl *httpController) kvPost(w http.ResponseWriter, req *http.Request) {
	kv := dto.KV{}

	if err := json.NewDecoder(req.Body).Decode(&kv); err != nil {
		http_helpers.RespondWithErr(
			w,
			http.StatusBadRequest,
			fmt.Errorf("decoding body json: %w", err),
		)
		return
	}

	if err := ctrl.uc.Set(kv.Key, kv.Value); err != nil {
		http_helpers.RespondWithErr(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("setting KV to usecase: %w", err),
		)
		return
	}

	http_helpers.RespondOK(w, kv)
}
