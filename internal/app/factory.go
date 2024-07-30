package app

import (
	"github.com/VadimGossip/tj-drs-storage/internal/imitator"
	"github.com/VadimGossip/tj-drs-storage/internal/rate"
)

type Factory struct {
	dbAdapter *DBAdapter

	rate     rate.Service
	imitator imitator.Service
}

var factory *Factory

//func newFactory(dbAdapter *DBAdapter, cfg *domain.Config) (*Factory, error) {
//	factory = &Factory{dbAdapter: dbAdapter}
//
//	factory.prometheusService = prometheus.NewService(cfg.Prometheus)
//	factory.protocolService = protocol.NewService()
//	othersStorages := make([]storage.Service, 0, cfg.OthersStorageManager.NumOfInstances)
//	for key := 1; key <= cfg.OthersStorageManager.NumOfInstances; key++ {
//		destStorageService := storageDest.NewService(dbAdapter.destRepo)
//		gatewayStorageService := storageGateway.NewService(dbAdapter.gatewayRepo)
//		ruleStorageService := storageRule.NewService(dbAdapter.ruleRepo, cfg.RouterLogic.SupPortionPrecision)
//		othersStorages = append(othersStorages, storage.NewService(int64(key), domain.OthersStorage, dbAdapter.storageRepo, dbAdapter.dstStorageRepo, destStorageService, gatewayStorageService, nil, ruleStorageService))
//	}
//	factory.othersStorageManager = smanager.NewService(cfg.OthersStorageManager, "others", domain.OthersStorage, othersStorages, factory.prometheusService)
//
//	ratesStorages := make([]storage.Service, 0, cfg.RatesStorageManager.NumOfInstances)
//	for key := 1; key <= cfg.RatesStorageManager.NumOfInstances; key++ {
//		ratesStorageService := storageRate.NewService(dbAdapter.rateRepo, dbAdapter.dstRateRepo)
//		ratesStorages = append(ratesStorages, storage.NewService(int64(key), domain.RatesStorage, dbAdapter.storageRepo, dbAdapter.dstStorageRepo, nil, nil, ratesStorageService, nil))
//	}
//	factory.ratesStorageManager = smanager.NewService(cfg.RatesStorageManager, "rates", domain.RatesStorage, ratesStorages, factory.prometheusService)
//
//	factory.destService = destination.NewService()
//	factory.rateService = rate.NewService()
//	factory.clientService = client.NewService(factory.destService, factory.rateService)
//	factory.tagService = tag.NewService()
//	factory.supplierService = supplier.NewService(factory.rateService, cfg.RouterLogic.MaxAnswerSuppliers, cfg.RouterLogic.SupPortionPrecision, cfg.RouterLogic.SupMinUIPortion)
//	factory.ruleService = rule.NewService(factory.supplierService, cfg.RouterLogic.MaxTraversLoops)
//	factory.routerService = router.NewService(factory.tagService, factory.clientService, factory.ruleService, factory.othersStorageManager, factory.ratesStorageManager, cfg.RouterLogic.DefaultScope)
//
//	return factory, nil
//}

func newFactory(dbAdapter *DBAdapter) *Factory {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.rate = rate.NewService(factory.dbAdapter.rateRepo)
	factory.imitator = imitator.NewService(factory.rate)

	return factory
}
