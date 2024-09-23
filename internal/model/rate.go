package model

const (
	SUPGWObjectKey  string = "SUPGW"
	DSTRObjectKey   string = "DSTR"
	DSTBObjectKey   string = "DSTB"
	RLObjectKey     string = "RL"
	SCDRObjectKey   string = "SCDR"
	AADObjectKey    string = "AAD"
	ABDObjectKey    string = "ABD"
	ACObjectKey     string = "AC"
	RAObjectKey     string = "RA"
	RBObjectKey     string = "RB"
	RTSObjectKey    string = "RTS"
	RVObjectKey     string = "RV"
	CURRTSObjectKey string = "CURRTS"
	TAGSObjectKey   string = "TAGS"
	DRSObjectKey    string = "DRS"
	DRSGObjectKey   string = "DRSG"
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
