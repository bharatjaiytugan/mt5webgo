package mt5webapi

// EnSoActivation - stop-out activation method.
type EnSoActivation uint32

const (
	ActivationNone       EnSoActivation = 0
	ActivationMarginCall EnSoActivation = 1
	ActivationStopOut    EnSoActivation = 2

	ActivationFirst EnSoActivation = ActivationNone
	ActivationLast  EnSoActivation = ActivationStopOut
)

// MTAccount represents account data.
type MTAccount struct {
	Login              uint64         `json:"Login"`
	CurrencyDigits     uint32         `json:"CurrencyDigits"`
	Balance            float64        `json:"Balance"`
	Credit             float64        `json:"Credit"`
	Margin             float64        `json:"Margin"`
	MarginFree         float64        `json:"MarginFree"`
	MarginLevel        float64        `json:"MarginLevel"`
	MarginLeverage     uint32         `json:"MarginLeverage"`
	Profit             float64        `json:"Profit"`
	Storage            float64        `json:"Storage"`
	Commission         float64        `json:"Commission"`
	Floating           float64        `json:"Floating"`
	Equity             float64        `json:"Equity"`
	SOActivation       EnSoActivation `json:"SOActivation"`
	SOTime             int64          `json:"SOTime"`
	SOLevel            float64        `json:"SOLevel"`
	SOEquity           float64        `json:"SOEquity"`
	SOMargin           float64        `json:"SOMargin"`
	MarginInitial      float64        `json:"MarginInitial"`
	MarginMaintenance  float64        `json:"MarginMaintenance"`
	Assets             float64        `json:"Assets"`
	Liabilities        float64        `json:"Liabilities"`
	BlockedCommission  float64        `json:"BlockedCommission"`
	BlockedProfit      float64        `json:"BlockedProfit"`
}
