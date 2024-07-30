package rate

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

type aRmsgKey struct {
	gwgrId    int64
	direction uint8
	bRmsgId   int64
	code      uint64
}

type bRmsgKey struct {
	gwgrId    int64
	direction uint8
	code      uint64
}

type rateKey struct {
	gwgrId    int64
	direction uint8
	aRmsgId   int64
	bRmsgId   int64
}

type IdHistItem struct {
	Id     int64 `json:"id"`
	DBegin int64 `json:"dBegin"`
	DEnd   int64 `json:"dEnd"`
}

type RmsRateHistItem struct {
	RmsrId int64 `json:"rmsr_id"`
	RmsvId int64 `json:"rmsv_id"`
	DBegin int64 `json:"dBegin"`
	DEnd   int64 `json:"dEnd"`
}

type Rate struct {
	Price      float64 `json:"price"`
	CurrencyId int64   `json:"currencyId"`
}

type CurrencyRateHist struct {
	CurrencyRate float64 `json:"currency_rate"`
	DBegin       int64   `json:"dBegin"`
	DEnd         int64   `json:"dEnd"`
}
