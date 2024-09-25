package converter

import (
	desc "github.com/VadimGossip/drs_data_loader/pkg/rate_v1"
	"github.com/VadimGossip/drs_storage_tester/internal/model"
)

func ToRateBaseFromFromDesc(rate *desc.RateBase) model.RateBase {
	return model.RateBase{
		RmsrId:    rate.RmsrId,
		PriceBase: rate.PriceBase,
	}
}

func ToSupRatesBaseFromDesc(rates []*desc.SupRateBase) map[int64]model.RateBase {
	result := make(map[int64]model.RateBase)
	for _, rate := range rates {
		result[rate.GwgrId] = ToRateBaseFromFromDesc(rate.Rate)
	}
	return result
}
