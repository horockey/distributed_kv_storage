package http_controller

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/horockey/distributed_kv_storage/internal/usecase"
	"github.com/rs/zerolog"
)

const shutdownTimeout = time.Second

//go:embed docs/openapi.html
var docsHtml []byte

type httpController struct {
	router *mux.Router
	serv   *http.Server

	uc *usecase.KVManagement

	logger zerolog.Logger
}

func New(
	adress string,
	uc *usecase.KVManagement,
	logger zerolog.Logger,
) *httpController {
	ctrl := httpController{
		uc:     uc,
		logger: logger,
	}

	ctrl.initRouter()

	ctrl.serv = &http.Server{
		Addr:    adress,
		Handler: ctrl.router,
	}

	return &ctrl
}

func (ctrl *httpController) Start(ctx context.Context) error {
	errs := make(chan error)
	go func() {
		defer close(errs)
		if err := ctrl.serv.ListenAndServe(); err != nil {
			errs <- err
		}
	}()
	ctrl.logger.Info().Str("addr", ctrl.serv.Addr).Msg("started")

	select {
	case err := <-errs:
		return fmt.Errorf("running http server: %w", err)
	case <-ctx.Done():
		resErr := ctx.Err()
		if errors.Is(resErr, context.Canceled) {
			resErr = nil
		}

		sdCtx, cancel := context.WithTimeout(context.TODO(), shutdownTimeout)
		defer cancel()

		if err := ctrl.serv.Shutdown(sdCtx); err != nil {
			resErr = errors.Join(resErr, fmt.Errorf("shutting down http server: %w", err))
		}
		return resErr
	}
}
