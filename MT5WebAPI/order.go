package mt5webapi

// EnOrderType - order types.
type EnOrderType uint32

const (
	OpBuy          EnOrderType = 0
	OpSell         EnOrderType = 1
	OpBuyLimit     EnOrderType = 2
	OpSellLimit    EnOrderType = 3
	OpBuyStop      EnOrderType = 4
	OpSellStop     EnOrderType = 5
	OpBuyStopLimit EnOrderType = 6
	OpSellStopLimit EnOrderType = 7
	OpCloseBy      EnOrderType = 8

	OpFirst EnOrderType = OpBuy
	OpLast  EnOrderType = OpCloseBy
)

// EnOrderFilling - order filling types.
type EnOrderFilling uint32

const (
	OrderFillFOK    EnOrderFilling = 0
	OrderFillIOC    EnOrderFilling = 1
	OrderFillReturn EnOrderFilling = 2

	OrderFillFirst EnOrderFilling = OrderFillFOK
	OrderFillLast  EnOrderFilling = OrderFillReturn
)

// EnOrderTime - order expiration types.
type EnOrderTime uint32

const (
	OrderTimeGTC          EnOrderTime = 0
	OrderTimeDay          EnOrderTime = 1
	OrderTimeSpecified    EnOrderTime = 2
	OrderTimeSpecifiedDay EnOrderTime = 3

	OrderTimeFirst EnOrderTime = OrderTimeGTC
	OrderTimeLast  EnOrderTime = OrderTimeSpecifiedDay
)

// EnOrderState - order state.
type EnOrderState uint32

const (
	OrderStateStarted       EnOrderState = 0
	OrderStatePlaced        EnOrderState = 1
	OrderStateCanceled      EnOrderState = 2
	OrderStatePartial       EnOrderState = 3
	OrderStateFilled        EnOrderState = 4
	OrderStateRejected      EnOrderState = 5
	OrderStateExpired       EnOrderState = 6
	OrderStateRequestAdd    EnOrderState = 7
	OrderStateRequestModify EnOrderState = 8
	OrderStateRequestCancel EnOrderState = 9

	OrderStateFirst EnOrderState = OrderStateStarted
	OrderStateLast  EnOrderState = OrderStateRequestCancel
)

// EnOrderActivation - order activation state.
type EnOrderActivation uint32

const (
	OrderActivationNone       EnOrderActivation = 0
	OrderActivationPending    EnOrderActivation = 1
	OrderActivationStopLimit  EnOrderActivation = 2
	OrderActivationExpiration EnOrderActivation = 3
	OrderActivationStopOut    EnOrderActivation = 4

	OrderActivationFirst EnOrderActivation = OrderActivationNone
	OrderActivationLast  EnOrderActivation = OrderActivationStopOut
)

// EnOrderReason - order creation reasons.
type EnOrderReason uint32

const (
	OrderReasonClient          EnOrderReason = 0
	OrderReasonExpert          EnOrderReason = 1
	OrderReasonDealer          EnOrderReason = 2
	OrderReasonSL              EnOrderReason = 3
	OrderReasonTP              EnOrderReason = 4
	OrderReasonSO              EnOrderReason = 5
	OrderReasonRollover        EnOrderReason = 6
	OrderReasonExternal        EnOrderReason = 7
	OrderReasonVMargin         EnOrderReason = 8
	OrderReasonGateway         EnOrderReason = 9
	OrderReasonSignal          EnOrderReason = 10
	OrderReasonSettlement      EnOrderReason = 11
	OrderReasonTransfer        EnOrderReason = 12
	OrderReasonSync            EnOrderReason = 13
	OrderReasonExternalService EnOrderReason = 14
	OrderReasonMigration       EnOrderReason = 15
	OrderReasonMobile          EnOrderReason = 16
	OrderReasonWeb             EnOrderReason = 17
	OrderReasonSplit           EnOrderReason = 18

	OrderReasonFirst EnOrderReason = OrderReasonClient
	OrderReasonLast  EnOrderReason = OrderReasonSplit
)

// EnOrderTradeActivationFlags - order activation flags.
type EnOrderTradeActivationFlags uint32

const (
	OrderActivFlagsNoLimit      EnOrderTradeActivationFlags = 0x01
	OrderActivFlagsNoStop       EnOrderTradeActivationFlags = 0x02
	OrderActivFlagsNoSLimit     EnOrderTradeActivationFlags = 0x04
	OrderActivFlagsNoSL         EnOrderTradeActivationFlags = 0x08
	OrderActivFlagsNoTP         EnOrderTradeActivationFlags = 0x10
	OrderActivFlagsNoSO         EnOrderTradeActivationFlags = 0x20
	OrderActivFlagsNoExpiration EnOrderTradeActivationFlags = 0x40

	OrderActivFlagsNone EnOrderTradeActivationFlags = 0x00
	OrderActivFlagsAll  EnOrderTradeActivationFlags = OrderActivFlagsNoLimit | OrderActivFlagsNoStop | OrderActivFlagsNoSLimit |
		OrderActivFlagsNoSL | OrderActivFlagsNoTP | OrderActivFlagsNoSO | OrderActivFlagsNoExpiration
)

// EnOrderTradeModifyFlags - modification flags.
type EnOrderTradeModifyFlags uint32

const (
	OrderModifyFlagsAdmin        EnOrderTradeModifyFlags = 0x001
	OrderModifyFlagsManager      EnOrderTradeModifyFlags = 0x002
	OrderModifyFlagsPosition     EnOrderTradeModifyFlags = 0x004
	OrderModifyFlagsRestore      EnOrderTradeModifyFlags = 0x008
	OrderModifyFlagsAPIAdmin     EnOrderTradeModifyFlags = 0x010
	OrderModifyFlagsAPIManager   EnOrderTradeModifyFlags = 0x020
	OrderModifyFlagsAPIServer    EnOrderTradeModifyFlags = 0x040
	OrderModifyFlagsAPIGateway   EnOrderTradeModifyFlags = 0x080
	OrderModifyFlagsAPIServerAdd EnOrderTradeModifyFlags = 0x100

	OrderModifyFlagsNone EnOrderTradeModifyFlags = 0x000
	OrderModifyFlagsAll  EnOrderTradeModifyFlags = OrderModifyFlagsAdmin | OrderModifyFlagsManager | OrderModifyFlagsPosition |
		OrderModifyFlagsRestore | OrderModifyFlagsAPIAdmin | OrderModifyFlagsAPIManager |
		OrderModifyFlagsAPIServer | OrderModifyFlagsAPIGateway | OrderModifyFlagsAPIServerAdd
)

// MTOrder represents a trade order.
type MTOrder struct {
	Order           uint64                      `json:"Order"`
	ExternalID      string                      `json:"ExternalID"`
	Login           uint64                      `json:"Login"`
	Dealer          uint64                      `json:"Dealer"`
	Symbol          string                      `json:"Symbol"`
	Digits          uint32                      `json:"Digits"`
	DigitsCurrency  uint32                      `json:"DigitsCurrency"`
	ContractSize    float64                     `json:"ContractSize"`
	State           EnOrderState                `json:"State"`
	Reason          EnOrderReason               `json:"Reason"`
	TimeSetup       int64                       `json:"TimeSetup"`
	TimeExpiration  int64                       `json:"TimeExpiration"`
	TimeDone        int64                       `json:"TimeDone"`
	TimeSetupMsc    int64                       `json:"TimeSetupMsc"`
	TimeDoneMsc     int64                       `json:"TimeDoneMsc"`
	ModifyFlags     EnOrderTradeModifyFlags     `json:"ModifyFlags"`
	Type            EnOrderType                 `json:"Type"`
	TypeFill        EnOrderFilling              `json:"TypeFill"`
	TypeTime        EnOrderTime                 `json:"TypeTime"`
	PriceOrder      float64                     `json:"PriceOrder"`
	PriceTrigger    float64                     `json:"PriceTrigger"`
	PriceCurrent    float64                     `json:"PriceCurrent"`
	PriceSL         float64                     `json:"PriceSL"`
	PriceTP         float64                     `json:"PriceTP"`
	ExpertID        uint64                      `json:"ExpertID"`
	ExpertPositionID uint64                     `json:"ExpertPositionID"`
	PositionByID    uint64                      `json:"PositionByID"`
	Comment         string                      `json:"Comment"`
	ActivationMode  EnOrderActivation           `json:"ActivationMode"`
	ActivationTime  int64                       `json:"ActivationTime"`
	ActivationPrice float64                     `json:"ActivationPrice"`
	ActivationFlags EnOrderTradeActivationFlags `json:"ActivationFlags"`

	// private volume fields - use getter/setter methods
	volumeInitial uint64
	volumeCurrent uint64
}

// VolumeInitial returns the order initial volume (old 4-digit format).
func (o *MTOrder) VolumeInitial() uint64 {
	return ConvertToOldVolume(o.volumeInitial)
}

// SetVolumeInitial sets the order initial volume from old 4-digit format.
func (o *MTOrder) SetVolumeInitial(v uint64) {
	o.volumeInitial = ConvertToNewVolume(v)
}

// VolumeInitialExt returns the order initial volume in extended 8-digit format.
func (o *MTOrder) VolumeInitialExt() uint64 {
	return o.volumeInitial
}

// SetVolumeInitialExt sets the order initial volume in extended 8-digit format.
func (o *MTOrder) SetVolumeInitialExt(v uint64) {
	o.volumeInitial = v
}

// VolumeCurrent returns the order current volume (old 4-digit format).
func (o *MTOrder) VolumeCurrent() uint64 {
	return ConvertToOldVolume(o.volumeCurrent)
}

// SetVolumeCurrent sets the order current volume from old 4-digit format.
func (o *MTOrder) SetVolumeCurrent(v uint64) {
	o.volumeCurrent = ConvertToNewVolume(v)
}

// VolumeCurrentExt returns the order current volume in extended 8-digit format.
func (o *MTOrder) VolumeCurrentExt() uint64 {
	return o.volumeCurrent
}

// SetVolumeCurrentExt sets the order current volume in extended 8-digit format.
func (o *MTOrder) SetVolumeCurrentExt(v uint64) {
	o.volumeCurrent = v
}
