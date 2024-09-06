package app

import (
	"context"
	"fmt"
	"log"

	"github.com/VadimGossip/tj-drs-storage/internal/client/db/tarantool"
	"github.com/VadimGossip/tj-drs-storage/internal/client/db/tarantool/tdb"
	"github.com/VadimGossip/tj-drs-storage/internal/repository"

	tRateRepo "github.com/VadimGossip/tj-drs-storage/internal/repository/rate/tarantool"
	"github.com/VadimGossip/tj-drs-storage/internal/service"
	tRateService "github.com/VadimGossip/tj-drs-storage/internal/service/rate"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/tj-drs-storage/internal/client/db/keydb"
	"github.com/VadimGossip/tj-drs-storage/internal/client/db/keydb/kdb"
	"github.com/VadimGossip/tj-drs-storage/internal/client/db/oracle"
	"github.com/VadimGossip/tj-drs-storage/internal/client/db/oracle/odb"
	"github.com/VadimGossip/tj-drs-storage/internal/closer"
	"github.com/VadimGossip/tj-drs-storage/internal/data"
	"github.com/VadimGossip/tj-drs-storage/internal/dbsource"
	"github.com/VadimGossip/tj-drs-storage/internal/domain"
	"github.com/VadimGossip/tj-drs-storage/internal/imitator"
	"github.com/VadimGossip/tj-drs-storage/internal/rate"
)

type serviceProvider struct {
	cfg          *domain.Config
	odbClient    oracle.Client
	txManager    oracle.TxManager
	dbSourceRepo dbsource.Repository

	kdbClient keydb.Client
	rateRepo  rate.Repository

	tarantoolClient   tarantool.Client
	tarantoolRateRepo repository.RateRepository

	dbSourceService dbsource.Service
	dataService     data.Service
	rateService     rate.Service
	tRateService    service.RateService
	imitatorService imitator.Service
}

func newServiceProvider(cfg *domain.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (s *serviceProvider) OdbClient(ctx context.Context) oracle.Client {
	if s.odbClient == nil {
		dsn := fmt.Sprintf(`user=%s password=%s connectString=%s:%d/%s`,
			s.cfg.DataSourceDb.Username,
			s.cfg.DataSourceDb.Password,
			s.cfg.DataSourceDb.Host,
			s.cfg.DataSourceDb.Port,
			s.cfg.DataSourceDb.Service)
		cl, err := odb.New(dsn)
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
		cl := kdb.New(s.cfg.TargetDb)

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
		cl, err := tdb.New(ctx, "todo_config")
		if err != nil {
			logrus.Fatalf("failed to create tarantool client: %s", err)
		}
		closer.Add(cl.Close)
		s.tarantoolClient = cl
	}

	return s.tarantoolClient
}

func (s *serviceProvider) DbSourceRepo(ctx context.Context) dbsource.Repository {
	if s.dbSourceRepo == nil {
		s.dbSourceRepo = dbsource.NewRepository(s.OdbClient(ctx))
	}
	return s.dbSourceRepo
}

func (s *serviceProvider) RateRepo(ctx context.Context) rate.Repository {
	if s.rateRepo == nil {
		s.rateRepo = rate.NewRepository(s.KeyDbClient(ctx))
	}
	return s.rateRepo
}

func (s *serviceProvider) TRateRepo(ctx context.Context) repository.RateRepository {
	if s.tarantoolRateRepo == nil {
		s.tarantoolRateRepo = tRateRepo.NewRepository(s.TarantoolClient(ctx))
	}
	return s.tarantoolRateRepo
}

func (s *serviceProvider) DbSourceService(ctx context.Context) dbsource.Service {
	if s.dbSourceService == nil {
		s.dbSourceService = dbsource.NewService(s.DbSourceRepo(ctx))
	}
	return s.dbSourceService
}

func (s *serviceProvider) DataService(ctx context.Context) data.Service {
	if s.dataService == nil {
		s.dataService = data.NewService(s.DbSourceService(ctx))
	}
	return s.dataService
}

func (s *serviceProvider) RateService(ctx context.Context) rate.Service {
	if s.rateService == nil {
		s.rateService = rate.NewService(s.RateRepo(ctx))
	}
	return s.rateService
}

func (s *serviceProvider) TRateService(ctx context.Context) service.RateService {
	if s.tRateService == nil {
		s.tRateService = tRateService.NewService(s.TRateRepo(ctx))
	}
	return s.rateService
}

func (s *serviceProvider) ImitatorService(ctx context.Context) imitator.Service {
	if s.imitatorService == nil {
		s.imitatorService = imitator.NewService(s.TRateService(ctx), s.DataService(ctx))
	}
	return s.imitatorService
}
