package mt5webapi

// EnDealFlags - common deal flags.
type EnDealFlags uint32

const (
	DealFlagsNone     EnDealFlags = 0x00000000
	DealFlagsRollover EnDealFlags = 0x00000001
	DealFlagsVMargin  EnDealFlags = 0x00000002

	DealFlagsAll EnDealFlags = DealFlagsRollover | DealFlagsVMargin
)

// EnDealAction - deal action types.
type EnDealAction uint32

const (
	DealBuy              EnDealAction = 0
	DealSell             EnDealAction = 1
	DealBalance          EnDealAction = 2
	DealCredit           EnDealAction = 3
	DealCharge           EnDealAction = 4
	DealCorrection       EnDealAction = 5
	DealBonus            EnDealAction = 6
	DealCommission       EnDealAction = 7
	DealCommissionDaily  EnDealAction = 8
	DealCommissionMonthly EnDealAction = 9
	DealAgentDaily       EnDealAction = 10
	DealAgentMonthly     EnDealAction = 11
	DealInterestRate     EnDealAction = 12
	DealBuyCanceled      EnDealAction = 13
	DealSellCanceled     EnDealAction = 14
	DealDividend         EnDealAction = 15
	DealDividendFranked  EnDealAction = 16
	DealTax              EnDealAction = 17
	DealAgent            EnDealAction = 18
	DealSOCompensation   EnDealAction = 19

	DealFirst EnDealAction = DealBuy
	DealLast  EnDealAction = DealSOCompensation
)

// EnEntryFlags - deal entry direction.
type EnEntryFlags uint32

const (
	EntryIn    EnEntryFlags = 0
	EntryOut   EnEntryFlags = 1
	EntryInOut EnEntryFlags = 2
	EntryOutBy EnEntryFlags = 3
	EntryState EnEntryFlags = 255

	EntryFirst EnEntryFlags = EntryIn
	EntryLast  EnEntryFlags = EntryState
)

// EnDealReason - deal creation reasons.
type EnDealReason uint32

const (
	DealReasonClient          EnDealReason = 0
	DealReasonExpert          EnDealReason = 1
	DealReasonDealer          EnDealReason = 2
	DealReasonSL              EnDealReason = 3
	DealReasonTP              EnDealReason = 4
	DealReasonSO              EnDealReason = 5
	DealReasonRollover        EnDealReason = 6
	DealReasonExternalClient  EnDealReason = 7
	DealReasonVMargin         EnDealReason = 8
	DealReasonGateway         EnDealReason = 9
	DealReasonSignal          EnDealReason = 10
	DealReasonSettlement      EnDealReason = 11
	DealReasonTransfer        EnDealReason = 12
	DealReasonSync            EnDealReason = 13
	DealReasonExternalService EnDealReason = 14
	DealReasonMigration       EnDealReason = 15
	DealReasonMobile          EnDealReason = 16
	DealReasonWeb             EnDealReason = 17
	DealReasonSplit           EnDealReason = 18

	DealReasonFirst EnDealReason = DealReasonClient
	DealReasonLast  EnDealReason = DealReasonSplit
)

// EnDealTradeModifyFlags - modification flags.
type EnDealTradeModifyFlags uint32

const (
	DealModifyFlagsAdmin        EnDealTradeModifyFlags = 0x001
	DealModifyFlagsManager      EnDealTradeModifyFlags = 0x002
	DealModifyFlagsPosition     EnDealTradeModifyFlags = 0x004
	DealModifyFlagsRestore      EnDealTradeModifyFlags = 0x008
	DealModifyFlagsAPIAdmin     EnDealTradeModifyFlags = 0x010
	DealModifyFlagsAPIManager   EnDealTradeModifyFlags = 0x020
	DealModifyFlagsAPIServer    EnDealTradeModifyFlags = 0x040
	DealModifyFlagsAPIGateway   EnDealTradeModifyFlags = 0x080
	DealModifyFlagsAPIServerAdd EnDealTradeModifyFlags = 0x100

	DealModifyFlagsNone EnDealTradeModifyFlags = 0x000
	DealModifyFlagsAll  EnDealTradeModifyFlags = DealModifyFlagsAdmin | DealModifyFlagsManager | DealModifyFlagsPosition |
		DealModifyFlagsRestore | DealModifyFlagsAPIAdmin | DealModifyFlagsAPIManager |
		DealModifyFlagsAPIServer | DealModifyFlagsAPIGateway | DealModifyFlagsAPIServerAdd
)

// MTDeal represents a deal.
type MTDeal struct {
	Deal            uint64                 `json:"Deal"`
	ExternalID      string                 `json:"ExternalID"`
	Login           uint64                 `json:"Login"`
	Dealer          uint64                 `json:"Dealer"`
	Order           uint64                 `json:"Order"`
	Action          EnDealAction           `json:"Action"`
	Entry           EnEntryFlags           `json:"Entry"`
	Reason          EnDealReason           `json:"Reason"`
	Digits          uint32                 `json:"Digits"`
	DigitsCurrency  uint32                 `json:"DigitsCurrency"`
	ContractSize    float64                `json:"ContractSize"`
	Time            int64                  `json:"Time"`
	TimeMsc         int64                  `json:"TimeMsc"`
	Symbol          string                 `json:"Symbol"`
	Price           float64                `json:"Price"`
	Profit          float64                `json:"Profit"`
	Storage         float64                `json:"Storage"`
	Commission      float64                `json:"Commission"`
	CommissionAgent float64                `json:"CommissionAgent"`
	RateProfit      float64                `json:"RateProfit"`
	RateMargin      float64                `json:"RateMargin"`
	ExpertID        uint64                 `json:"ExpertID"`
	PositionID      uint64                 `json:"PositionID"`
	Comment         string                 `json:"Comment"`
	ProfitRaw       float64                `json:"ProfitRaw"`
	PricePosition   float64                `json:"PricePosition"`
	TickValue       float64                `json:"TickValue"`
	TickSize        float64                `json:"TickSize"`
	Flags           uint64                 `json:"Flags"`
	Gateway         string                 `json:"Gateway"`
	PriceGateway    float64                `json:"PriceGateway"`
	ModifyFlags     EnDealTradeModifyFlags `json:"ModifyFlags"`
	PriceSL         float64                `json:"PriceSL"`
	PriceTP         float64                `json:"PriceTP"`

	// private volume fields - use getter/setter methods
	volume       uint64
	volumeClosed uint64
}

// Volume returns the deal volume (old 4-digit format).
func (d *MTDeal) Volume() uint64 {
	return ConvertToOldVolume(d.volume)
}

// SetVolume sets the deal volume from old 4-digit format.
func (d *MTDeal) SetVolume(v uint64) {
	d.volume = ConvertToNewVolume(v)
}

// VolumeExt returns the deal volume in extended 8-digit format.
func (d *MTDeal) VolumeExt() uint64 {
	return d.volume
}

// SetVolumeExt sets the deal volume in extended 8-digit format.
func (d *MTDeal) SetVolumeExt(v uint64) {
	d.volume = v
}

// VolumeClosed returns the closed volume (old 4-digit format).
func (d *MTDeal) VolumeClosed() uint64 {
	return ConvertToOldVolume(d.volumeClosed)
}

// SetVolumeClosed sets the closed volume from old 4-digit format.
func (d *MTDeal) SetVolumeClosed(v uint64) {
	d.volumeClosed = ConvertToNewVolume(v)
}

// VolumeClosedExt returns the closed volume in extended 8-digit format.
func (d *MTDeal) VolumeClosedExt() uint64 {
	return d.volumeClosed
}

// SetVolumeClosedExt sets the closed volume in extended 8-digit format.
func (d *MTDeal) SetVolumeClosedExt(v uint64) {
	d.volumeClosed = v
}
