package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	instance_manager "github.com/horockey/consul_instance_manager"
	"github.com/horockey/distributed_kv_storage/internal/adapter/gateway/remote_storage/http_remote_storage"
	"github.com/horockey/distributed_kv_storage/internal/adapter/repository/local_storage/inmemory_local_storage"
	"github.com/horockey/distributed_kv_storage/internal/config"
	"github.com/horockey/distributed_kv_storage/internal/controller/http_controller"
	"github.com/horockey/distributed_kv_storage/internal/usecase"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().
			Err(fmt.Errorf("creating config: %w", err)).
			Send()
	}

	localStorage := inmemory_local_storage.New()
	remoteStorage := http_remote_storage.New()

	iman, err := instance_manager.NewClient(
		config.AppName,
		instance_manager.WithLogger(logger.With().Str("layer", "instance manager").Logger()),
		instance_manager.WithDownHoldDuration(time.Duration(cfg.InstanceManager.DownHoldDirationMsec)*time.Millisecond),
		instance_manager.WithPollInterval(time.Duration(cfg.InstanceManager.PollIntervalMsec)*time.Millisecond),
	)
	if err != nil {
		logger.Fatal().
			Err(fmt.Errorf("creating instance manager: %w", err)).
			Send()
	}

	uc := usecase.New(
		cfg.Hostname,
		localStorage,
		remoteStorage,
		iman,
		logger.With().Str("layer", "usecase").Logger(),
	)

	serviceHttpAddr := fmt.Sprintf("%s:%d", cfg.Http.BindAddr, cfg.Http.Port)
	ctrl := http_controller.New(
		serviceHttpAddr,
		uc,
		logger.With().Str("layer", "http_controller").Logger(),
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGABRT,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := iman.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error().
				Err(fmt.Errorf("running instance manager: %w", err)).
				Send()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := ctrl.Start(ctx); err != nil {
			logger.Error().
				Err(fmt.Errorf("running http conmtroller: %w", err)).
				Send()
		}
	}()

	if err := registerInConsul(cfg.Hostname, serviceHttpAddr); err != nil {
		logger.Fatal().
			Err(fmt.Errorf("registering in consul")).
			Send()
	}
	logger.Debug().Msg("service registered in consul")

	defer func() {
		if err := deregisterFromConsul(cfg.Hostname); err != nil {
			logger.Error().
				Err(fmt.Errorf("deregistering from consul")).
				Send()
		}
		logger.Debug().Msg("service deregistered from consul")
	}()

	logger.Info().Msg("app started")
	wg.Wait()
	logger.Info().Msg("app stopped")
}

func registerInConsul(hostname string, address string) error {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return fmt.Errorf("creating consul client: %w", err)
	}

	if _, err := client.Catalog().Register(&consul.CatalogRegistration{
		ID:      uuid.NewString(),
		Node:    hostname,
		Address: address,
		Service: &consul.AgentService{
			ID:      config.AppName + "_" + hostname,
			Service: config.AppName,
		},
		Checks: consul.HealthChecks{
			{
				Node:    hostname,
				CheckID: uuid.NewString(),
				Status:  consul.HealthPassing,
			},
		},
	}, nil); err != nil {
		return fmt.Errorf("registering in consul: %w", err)
	}

	return nil
}

func deregisterFromConsul(hostname string) error {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return fmt.Errorf("creting consul client: %w", err)
	}

	_, err = client.Catalog().Deregister(&consul.CatalogDeregistration{
		Node:      hostname,
		ServiceID: config.AppName + "_" + hostname,
	}, nil)
	if err != nil {
		return fmt.Errorf("deregistering from consul: %w", err)
	}

	return nil
}
