package usecase

import (
	"errors"
	"fmt"

	"github.com/horockey/distributed_kv_storage/internal/adapter/gateway/remote_storage"
	"github.com/horockey/distributed_kv_storage/internal/adapter/repository/local_storage"
	instance_manager "github.com/horockey/go-consul-instance-manager"
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

func (uc *KVManagement) Get(key string) (res map[string]any, resErr error) {
	dataHolders, err := uc.iman.GetDataHolders(key)
	if err != nil {
		err = fmt.Errorf("getting data owner for key %s from iman: %w", key, err)
		uc.logger.Error().Err(err).Send()
		return nil, err
	}

	var val map[string]any

	for _, dataHolder := range dataHolders {
		switch dataHolder.Name() {
		case uc.localHostName:
			val, err = uc.localKV.Get(key)
			if err != nil {
				err = fmt.Errorf("getting data from local storage: %w", err)
				resErr = errors.Join(resErr, err)
				uc.logger.Error().Err(err).Send()
				continue
			}
		default:
			val, err = uc.remoteKV.Get(key, remote_storage.AppNode{
				Name:    dataHolder.Name(),
				Address: dataHolder.Address(),
			})
			if err != nil {
				err = fmt.Errorf("getting data from remote storage (%s): %w", dataHolder.Name(), err)
				resErr = errors.Join(resErr, err)
				uc.logger.Error().Err(err).Send()
				continue
			}
		}
		return val, nil
	}

	return nil, resErr
}

func (uc *KVManagement) Set(key string, val map[string]any) (resErr error) {
	dataHolders, err := uc.iman.GetDataHolders(key)
	if err != nil {
		err = fmt.Errorf("getting data owners for key %s from iman: %w", key, err)
		uc.logger.Error().Err(err).Send()
		return err
	}

	for _, dataHolder := range dataHolders {
		switch dataHolder.Name() {
		case uc.localHostName:
			err = uc.localKV.Set(key, val)
			if err != nil {
				err = fmt.Errorf("setting data to local storage: %w", err)
				resErr = errors.Join(resErr, err)
				uc.logger.Error().Err(err).Send()
				continue
			}
		default:
			err = uc.remoteKV.Set(key, val, remote_storage.AppNode{
				Name:    dataHolder.Name(),
				Address: dataHolder.Address(),
			})
			if err != nil {
				err = fmt.Errorf("setting data to remote storage (%s): %w", dataHolder.Name(), err)
				resErr = errors.Join(resErr, err)
				uc.logger.Error().Err(err).Send()
				continue
			}
			continue
		}
		return nil
	}

	return resErr
}
