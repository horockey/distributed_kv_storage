package usecase

import (
	"fmt"

	instance_manager "github.com/horockey/consul_instance_manager"
	"github.com/horockey/distributed_kv_storage/internal/adapter/gateway/remote_storage"
	"github.com/horockey/distributed_kv_storage/internal/adapter/repository/local_storage"
	"github.com/rs/zerolog"
)

type KVManagement struct {
	localHostName string

	localKV local_storage.Repository

	remoteKV remote_storage.Gateway
	iman     *instance_manager.Client

	logger zerolog.Logger
}

func New(
	localHostName string,
	localKV local_storage.Repository,
	remoteKV remote_storage.Gateway,
	iman *instance_manager.Client,
	logger zerolog.Logger,
) *KVManagement {
	return &KVManagement{
		localHostName: localHostName,
		localKV:       localKV,
		remoteKV:      remoteKV,
		iman:          iman,
		logger:        logger,
	}
}

func (uc *KVManagement) Get(key string) (map[string]any, error) {
	dataHolder, err := uc.iman.GetDataHolder(key)
	if err != nil {
		err = fmt.Errorf("getting data owner for key %s from iman: %w", key, err)
		uc.logger.Error().Err(err).Send()
		return nil, err
	}

	var val map[string]any

	switch dataHolder.Name() {
	case uc.localHostName:
		val, err = uc.localKV.Get(key)
		if err != nil {
			err = fmt.Errorf("getting data from local storage: %w", err)
			uc.logger.Error().Err(err).Send()
			return nil, err
		}
	default:
		val, err = uc.remoteKV.Get(key, remote_storage.AppNode{
			Name:    dataHolder.Name(),
			Address: dataHolder.Address(),
		})
		if err != nil {
			err = fmt.Errorf("getting data from remote storage (%s): %w", dataHolder.Name(), err)
			uc.logger.Error().Err(err).Send()
			return nil, err
		}
	}

	return val, nil
}

func (uc *KVManagement) Set(key string, val map[string]any) error {
	dataHolder, err := uc.iman.GetDataHolder(key)
	if err != nil {
		err = fmt.Errorf("getting data owner for key %s from iman: %w", key, err)
		uc.logger.Error().Err(err).Send()
		return err
	}

	switch dataHolder.Name() {
	case uc.localHostName:
		err = uc.localKV.Set(key, val)
		if err != nil {
			err = fmt.Errorf("setting data to local storage: %w", err)
			uc.logger.Error().Err(err).Send()
			return err
		}
	default:
		err = uc.remoteKV.Set(key, val, remote_storage.AppNode{
			Name:    dataHolder.Name(),
			Address: dataHolder.Address(),
		})
		if err != nil {
			err = fmt.Errorf("setting data to remote storage (%s): %w", dataHolder.Name(), err)
			uc.logger.Error().Err(err).Send()
			return err
		}
	}

	return nil
}
