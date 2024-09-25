package model

const (
	SUPGWObjectKey  string = "SUPGW"
	RAObjectKey     string = "RA"
	RBObjectKey     string = "RB"
	RTSObjectKey    string = "RTS"
	RVObjectKey     string = "RV"
	CURRTSObjectKey string = "CURRTS"
)

type ARmsgKey struct {
	GwgrId    int64
	Direction uint8
	BRmsgId   int64
	Code      string
}

type BRmsgKey struct {
	GwgrId    int64
	Direction uint8
	Code      string
}

type IdHistItem struct {
	Id     int64
	DBegin int64
	DEnd   int64
}

type RateKey struct {
	GwgrId    int64
	Direction uint8
	ARmsgId   int64
	BRmsgId   int64
}

type RmsRateHistItem struct {
	RmsrId int64
	RmsvId int64
	DBegin int64
	DEnd   int64
}

type Rate struct {
	Price      float64
	CurrencyId int64
}

type CurrencyRateHist struct {
	CurrencyRate float64
	DBegin       int64
	DEnd         int64
}

type RateBase struct {
	RmsrId    int64
	PriceBase float64
}
