package usecase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/horockey/distributed_kv_storage/internal/usecase"
	"github.com/horockey/distributed_kv_storage/internal/usecase/mocks"
	consul_iman "github.com/horockey/go-consul-instance-manager"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

//go:generate mockery --dir ../adapter/repository/local_storage --name Repository --filename local_storage.go --structname LocalStorage
//go:generate mockery --dir ../adapter/gateway/remote_storage --name Gateway --filename remote_storage.go --structname RemoteStorage

const (
	hostname    = "host1"
	serviceName = "test_service"
	consulAddr  = "localhost:8500"
	addr        = "http://host1:8080"
)

var logger = zerolog.New(zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: time.RFC3339,
}).With().Timestamp().Logger()

type UsecaseTestSuite struct {
	suite.Suite

	iman *consul_iman.Client
	uc   *usecase.KVManagement

	localStorage  *mocks.LocalStorage
	remoteStorage *mocks.RemoteStorage

	consul testcontainers.Container
}

func (s *UsecaseTestSuite) SetupTest() {
	t := s.T()
	ctx := context.TODO()

	s.localStorage = mocks.NewLocalStorage(t)
	s.remoteStorage = mocks.NewRemoteStorage(t)

	var err error
	s.consul, err = setupConsul()
	require.NoError(t, err)

	err = s.consul.Start(ctx)
	require.NoError(t, err)
	time.Sleep(time.Second)

	consulCfg := api.DefaultConfig()
	consulCfg.Address = consulAddr
	consulClient, err := api.NewClient(consulCfg)
	require.NoError(t, err)

	s.iman, err = consul_iman.NewClient(
		serviceName,
		consul_iman.WithConsulClient(consulClient),
		consul_iman.WithPollInterval(time.Millisecond*500),
		consul_iman.WithDownHoldDuration(time.Second*3),
	)
	require.NoError(t, err)

	err = s.iman.Register(hostname, addr)
	require.NoError(t, err)

	s.uc = usecase.New(hostname, s.localStorage, s.remoteStorage, s.iman, logger)
}

func (s *UsecaseTestSuite) TearDownTest() {
	t := s.T()
	ctx := context.TODO()

	err := s.consul.Terminate(ctx)
	require.NoError(t, err)
}

func TestUsecaseSuite(t *testing.T) {
	suite.Run(t, &UsecaseTestSuite{})
}

func (s *UsecaseTestSuite) TestGet() {
	t := s.T()
	ctx := context.TODO()

	expectedVal := map[string]any{"k2": "v2"}

	s.localStorage.On("Get", "k1").Return(expectedVal, nil)

	go s.iman.Start(ctx)
	time.Sleep(time.Millisecond * 600)

	val, err := s.uc.Get("k1")
	require.NoError(t, err)
	require.Equal(t, expectedVal, val)
}

func (s *UsecaseTestSuite) TestSet() {
	t := s.T()
	ctx := context.TODO()

	val := map[string]any{"k2": "v2"}

	s.localStorage.On("Set", "k1", val).Return(nil)

	go s.iman.Start(ctx)
	time.Sleep(time.Millisecond * 600)

	err := s.uc.Set("k1", val)
	require.NoError(t, err)
}

func setupConsul() (testcontainers.Container, error) {
	return testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "hashicorp/consul:latest",
			ExposedPorts: []string{"8500:8500"},
			Name:         "consul",
			Hostname:     "consul",
			Networks:     []string{"testnet"},
		},
	})
}
