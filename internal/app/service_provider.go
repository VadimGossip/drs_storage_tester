package app

import (
	"context"
	"log"

	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/VadimGossip/platform_common/pkg/db/keydb"
	"github.com/VadimGossip/platform_common/pkg/db/keydb/kdb"
	"github.com/VadimGossip/platform_common/pkg/db/oracle"
	"github.com/VadimGossip/platform_common/pkg/db/oracle/odb"
	"github.com/VadimGossip/platform_common/pkg/db/tarantool"
	"github.com/VadimGossip/platform_common/pkg/db/tarantool/tdb"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/drs_storage_tester/internal/client/grpc"
	rateGrpc "github.com/VadimGossip/drs_storage_tester/internal/client/grpc/rate"
	"github.com/VadimGossip/drs_storage_tester/internal/config"
	clientCfg "github.com/VadimGossip/drs_storage_tester/internal/config/client"
	dbCfg "github.com/VadimGossip/drs_storage_tester/internal/config/db"
	serviceCfg "github.com/VadimGossip/drs_storage_tester/internal/config/service"
	"github.com/VadimGossip/drs_storage_tester/internal/repository"
	kdbRateRepo "github.com/VadimGossip/drs_storage_tester/internal/repository/rate/kdb"
	tarantoolRateRepo "github.com/VadimGossip/drs_storage_tester/internal/repository/rate/tarantool"
	requestRepo "github.com/VadimGossip/drs_storage_tester/internal/repository/request"
	"github.com/VadimGossip/drs_storage_tester/internal/service"
	dataService "github.com/VadimGossip/drs_storage_tester/internal/service/data"
	eventService "github.com/VadimGossip/drs_storage_tester/internal/service/event"
	imitatorService "github.com/VadimGossip/drs_storage_tester/internal/service/imitator"
	rateService "github.com/VadimGossip/drs_storage_tester/internal/service/rate"
	requestService "github.com/VadimGossip/drs_storage_tester/internal/service/request"
)

type serviceProvider struct {
	oracleConfig          config.OracleConfig
	kdbConfig             config.KdbConfig
	tarantoolConfig       config.TarantoolConfig
	rateGrpcConfig        config.RateGrpcConfig
	serviceProviderConfig config.ServiceProviderConfig
	imitatorConfig        config.ImitatorConfig

	odbClient       oracle.Client
	txManager       oracle.TxManager
	kdbClient       keydb.Client
	tarantoolClient tarantool.Client
	rateGrpcClient  grpc.RateClient

	rateRepo    repository.RateRepository
	requestRepo repository.RequestRepository

	requestService  service.RequestService
	dataService     service.DataService
	eventService    service.EventService
	rateService     service.RateService
	imitatorService service.ImitatorService
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) OracleConfig() config.OracleConfig {
	if s.oracleConfig == nil {
		cfg, err := dbCfg.NewOracleConfig()
		if err != nil {
			log.Fatalf("failed to get oracleConfig: %s", err)
		}

		s.oracleConfig = cfg
	}

	return s.oracleConfig
}

func (s *serviceProvider) KdbConfig() config.KdbConfig {
	if s.kdbConfig == nil {
		cfg, err := dbCfg.NewKdbConfig()
		if err != nil {
			log.Fatalf("failed to get kdbConfig: %s", err)
		}

		s.kdbConfig = cfg
	}

	return s.kdbConfig
}

func (s *serviceProvider) TarantoolConfig() config.TarantoolConfig {
	if s.tarantoolConfig == nil {
		cfg, err := dbCfg.NewTarantoolConfig()
		if err != nil {
			log.Fatalf("failed to get tarantoolConfig: %s", err)
		}

		s.tarantoolConfig = cfg
	}

	return s.tarantoolConfig
}

func (s *serviceProvider) RateGRPCClientConfig() config.RateGrpcConfig {
	if s.rateGrpcConfig == nil {
		cfg, err := clientCfg.NewRateGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get rateGrpcConfig: %s", err)
		}

		s.rateGrpcConfig = cfg
	}

	return s.rateGrpcConfig
}

func (s *serviceProvider) ServiceProviderConfig() config.ServiceProviderConfig {
	if s.serviceProviderConfig == nil {
		cfg, err := serviceCfg.NewServiceProviderConfig()
		if err != nil {
			log.Fatalf("failed to get serviceProviderConfig: %s", err)
		}

		s.serviceProviderConfig = cfg
	}

	return s.serviceProviderConfig
}

func (s *serviceProvider) ImitatorConfig() config.ImitatorConfig {
	if s.imitatorConfig == nil {
		cfg, err := serviceCfg.NewImitatorConfig()
		if err != nil {
			log.Fatalf("failed to get imitatorConfig: %s", err)
		}

		s.imitatorConfig = cfg
	}

	return s.imitatorConfig
}

func (s *serviceProvider) OdbClient(ctx context.Context) oracle.Client {
	if s.odbClient == nil {
		cl, err := odb.New(s.OracleConfig().DSN())
		if err != nil {
			logrus.Fatalf("failed to create odb client: %s", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatalf("odb ping error: %s", err)
		}
		closer.Add(cl.Close)
		s.odbClient = cl
	}

	return s.odbClient
}

func (s *serviceProvider) KeyDbClient(ctx context.Context) keydb.Client {
	if s.kdbClient == nil {
		cl := kdb.New(kdb.ClientOptions{
			Addr:         s.KdbConfig().Address(),
			Username:     s.KdbConfig().Username(),
			Password:     s.KdbConfig().Password(),
			DB:           s.KdbConfig().DB(),
			ReadTimeout:  s.KdbConfig().ReadTimeoutSec(),
			WriteTimeout: s.KdbConfig().ReadTimeoutSec(),
		})

		if err := cl.DB().Ping(ctx); err != nil {
			log.Fatalf("kdb ping error: %s", err)
		}

		closer.Add(cl.Close)
		s.kdbClient = cl
	}

	return s.kdbClient
}

func (s *serviceProvider) TarantoolClient(ctx context.Context) tarantool.Client {
	if s.tarantoolClient == nil {
		cl, err := tdb.New(ctx,
			s.TarantoolConfig().Address(),
			s.TarantoolConfig().Username(),
			s.TarantoolConfig().Password(),
		)
		if err != nil {
			logrus.Fatalf("failed to create tarantool client: %s", err)
		}
		closer.Add(cl.Close)
		s.tarantoolClient = cl
	}

	return s.tarantoolClient
}

func (s *serviceProvider) RateGRPCClient() grpc.RateClient {
	if s.rateGrpcClient == nil {
		grpcAuthClient, err := rateGrpc.NewClient(s.RateGRPCClientConfig())
		if err != nil {
			logrus.Fatalf("failed to create rateGrpcClient: %s", err)
		}
		s.rateGrpcClient = grpcAuthClient
	}
	return s.rateGrpcClient
}

func (s *serviceProvider) RateRepo(ctx context.Context) repository.RateRepository {
	if s.rateRepo == nil {
		testDb := s.ServiceProviderConfig().TestDB()
		if testDb == s.ServiceProviderConfig().TarantoolTestDB() {
			s.rateRepo = tarantoolRateRepo.NewRepository(s.TarantoolClient(ctx))
		} else if testDb == s.ServiceProviderConfig().KdbTestDB() {
			s.rateRepo = kdbRateRepo.NewRepository(s.KeyDbClient(ctx))
		} else {
			s.rateRepo = s.RateGRPCClient()
		}

	}
	return s.rateRepo
}

func (s *serviceProvider) RequestRepo(ctx context.Context) repository.RequestRepository {
	if s.requestRepo == nil {
		s.requestRepo = requestRepo.NewRepository(s.OdbClient(ctx))
	}
	return s.requestRepo
}

func (s *serviceProvider) RequestService(ctx context.Context) service.RequestService {
	if s.requestService == nil {
		s.requestService = requestService.NewService(s.RequestRepo(ctx))
	}
	return s.requestService
}

func (s *serviceProvider) DataService(ctx context.Context) service.DataService {
	if s.dataService == nil {
		s.dataService = dataService.NewService(s.RequestService(ctx))
	}
	return s.dataService
}

func (s *serviceProvider) EventService() service.EventService {
	if s.eventService == nil {
		s.eventService = eventService.NewService()
	}
	return s.eventService
}

func (s *serviceProvider) RateService(ctx context.Context) service.RateService {
	if s.rateService == nil {
		s.rateService = rateService.NewService(s.RateRepo(ctx))
	}
	return s.rateService
}

func (s *serviceProvider) ImitatorService(ctx context.Context) service.ImitatorService {
	if s.imitatorService == nil {
		s.imitatorService = imitatorService.NewService(s.ImitatorConfig(), s.RateService(ctx), s.DataService(ctx))
	}
	return s.imitatorService
}
