package mt5webapi

// EnPositionAction - position types.
type EnPositionAction uint32

const (
	PositionBuy  EnPositionAction = 0
	PositionSell EnPositionAction = 1

	PositionFirst EnPositionAction = PositionBuy
	PositionLast  EnPositionAction = PositionSell
)

// EnPositionActivation - activation modes.
type EnPositionActivation uint32

const (
	PositionActivationNone    EnPositionActivation = 0
	PositionActivationSL      EnPositionActivation = 1
	PositionActivationTP      EnPositionActivation = 2
	PositionActivationStopOut EnPositionActivation = 3

	PositionActivationFirst EnPositionActivation = PositionActivationNone
	PositionActivationLast  EnPositionActivation = PositionActivationStopOut
)

// EnPositionTradeActivationFlags - position activation flags.
type EnPositionTradeActivationFlags uint32

const (
	PositionActivFlagsNoLimit      EnPositionTradeActivationFlags = 0x01
	PositionActivFlagsNoStop       EnPositionTradeActivationFlags = 0x02
	PositionActivFlagsNoSLimit     EnPositionTradeActivationFlags = 0x04
	PositionActivFlagsNoSL         EnPositionTradeActivationFlags = 0x08
	PositionActivFlagsNoTP         EnPositionTradeActivationFlags = 0x10
	PositionActivFlagsNoSO         EnPositionTradeActivationFlags = 0x20
	PositionActivFlagsNoExpiration EnPositionTradeActivationFlags = 0x40

	PositionActivFlagsNone EnPositionTradeActivationFlags = 0x00
	PositionActivFlagsAll  EnPositionTradeActivationFlags = PositionActivFlagsNoLimit | PositionActivFlagsNoStop |
		PositionActivFlagsNoSLimit | PositionActivFlagsNoSL | PositionActivFlagsNoTP |
		PositionActivFlagsNoSO | PositionActivFlagsNoExpiration
)

// EnPositionReason - position creation reasons.
type EnPositionReason uint32

const (
	PositionReasonClient          EnPositionReason = 0
	PositionReasonExpert          EnPositionReason = 1
	PositionReasonDealer          EnPositionReason = 2
	PositionReasonSL              EnPositionReason = 3
	PositionReasonTP              EnPositionReason = 4
	PositionReasonSO              EnPositionReason = 5
	PositionReasonRollover        EnPositionReason = 6
	PositionReasonExternalClient  EnPositionReason = 7
	PositionReasonVMargin         EnPositionReason = 8
	PositionReasonGateway         EnPositionReason = 9
	PositionReasonSignal          EnPositionReason = 10
	PositionReasonSettlement      EnPositionReason = 11
	PositionReasonTransfer        EnPositionReason = 12
	PositionReasonSync            EnPositionReason = 13
	PositionReasonExternalService EnPositionReason = 14
	PositionReasonMigration       EnPositionReason = 15
	PositionReasonMobile          EnPositionReason = 16
	PositionReasonWeb             EnPositionReason = 17
	PositionReasonSplit           EnPositionReason = 18

	PositionReasonFirst EnPositionReason = PositionReasonClient
	PositionReasonLast  EnPositionReason = PositionReasonSplit
)

// EnPositionTradeModifyFlags - modification flags.
type EnPositionTradeModifyFlags uint32

const (
	PositionModifyFlagsAdmin      EnPositionTradeModifyFlags = 0x01
	PositionModifyFlagsManager    EnPositionTradeModifyFlags = 0x02
	PositionModifyFlagsPosition   EnPositionTradeModifyFlags = 0x04
	PositionModifyFlagsRestore    EnPositionTradeModifyFlags = 0x08
	PositionModifyFlagsAPIAdmin   EnPositionTradeModifyFlags = 0x10
	PositionModifyFlagsAPIManager EnPositionTradeModifyFlags = 0x20
	PositionModifyFlagsAPIServer  EnPositionTradeModifyFlags = 0x40
	PositionModifyFlagsAPIGateway EnPositionTradeModifyFlags = 0x80

	PositionModifyFlagsNone EnPositionTradeModifyFlags = 0x00
	PositionModifyFlagsAll  EnPositionTradeModifyFlags = PositionModifyFlagsAdmin | PositionModifyFlagsManager |
		PositionModifyFlagsPosition | PositionModifyFlagsRestore | PositionModifyFlagsAPIAdmin |
		PositionModifyFlagsAPIManager | PositionModifyFlagsAPIServer | PositionModifyFlagsAPIGateway
)

// MTPosition represents a position.
type MTPosition struct {
	Position        uint64                         `json:"Position"`
	ExternalID      string                         `json:"ExternalID"`
	Login           uint64                         `json:"Login"`
	Dealer          uint64                         `json:"Dealer"`
	Symbol          string                         `json:"Symbol"`
	Action          EnPositionAction               `json:"Action"`
	Digits          uint32                         `json:"Digits"`
	DigitsCurrency  uint32                         `json:"DigitsCurrency"`
	Reason          EnPositionReason               `json:"Reason"`
	ContractSize    float64                        `json:"ContractSize"`
	TimeCreate      int64                          `json:"TimeCreate"`
	TimeUpdate      int64                          `json:"TimeUpdate"`
	ModifyFlags     EnPositionTradeModifyFlags     `json:"ModifyFlags"`
	PriceOpen       float64                        `json:"PriceOpen"`
	PriceCurrent    float64                        `json:"PriceCurrent"`
	PriceSL         float64                        `json:"PriceSL"`
	PriceTP         float64                        `json:"PriceTP"`
	Profit          float64                        `json:"Profit"`
	Storage         float64                        `json:"Storage"`
	Commission      float64                        `json:"Commission"`
	RateProfit      float64                        `json:"RateProfit"`
	RateMargin      float64                        `json:"RateMargin"`
	ExpertID        uint64                         `json:"ExpertID"`
	ExpertPositionID uint64                        `json:"ExpertPositionID"`
	Comment         string                         `json:"Comment"`
	ActivationMode  EnPositionActivation           `json:"ActivationMode"`
	ActivationTime  int64                          `json:"ActivationTime"`
	ActivationPrice float64                        `json:"ActivationPrice"`
	ActivationFlags EnPositionTradeActivationFlags `json:"ActivationFlags"`

	// private volume field - use getter/setter methods
	volume uint64
}

// Volume returns the position volume (old 4-digit format).
// Note: The C# original uses ConvertToOldVolume in both getter and setter.
func (p *MTPosition) Volume() uint64 {
	return ConvertToOldVolume(p.volume)
}

// SetVolume sets the position volume from old 4-digit format.
// Note: The C# original uses ConvertToOldVolume (not ConvertToNewVolume) in the setter.
func (p *MTPosition) SetVolume(v uint64) {
	p.volume = ConvertToOldVolume(v)
}

// VolumeExt returns the position volume in extended 8-digit format.
func (p *MTPosition) VolumeExt() uint64 {
	return p.volume
}

// SetVolumeExt sets the position volume in extended 8-digit format.
func (p *MTPosition) SetVolumeExt(v uint64) {
	p.volume = v
}
