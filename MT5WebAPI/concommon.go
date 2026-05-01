package mt5webapi

// EnUpdateMode - LiveUpdate modes
type EnUpdateMode uint32

const (
	UpdateDisable    EnUpdateMode = 0
	UpdateEnable     EnUpdateMode = 1
	UpdateEnableBeta EnUpdateMode = 2
	UpdateFirst      EnUpdateMode = UpdateDisable
	UpdateLast       EnUpdateMode = UpdateEnableBeta
)

// MTConCommon represents common server configuration.
type MTConCommon struct {
	Name               string       `json:"Name"`
	Owner              string       `json:"Owner"`
	OwnerID            string       `json:"OwnerID"`
	OwnerHost          string       `json:"OwnerHost"`
	OwnerEmail         string       `json:"OwnerEmail"`
	Product            string       `json:"Product"`
	ExpirationLicense  int64        `json:"ExpirationLicense"`
	ExpirationSupport  int64        `json:"ExpirationSupport"`
	LimitTradeServers  uint32       `json:"LimitTradeServers"`
	LimitWebServers    uint32       `json:"LimitWebServers"`
	LimitAccounts      uint32       `json:"LimitAccounts"`
	LimitDeals         uint32       `json:"LimitDeals"`
	LimitSymbols       uint32       `json:"LimitSymbols"`
	LimitGroups        uint32       `json:"LimitGroups"`
	LiveUpdateMode     EnUpdateMode `json:"LiveUpdateMode"`
	TotalUsers         uint32       `json:"TotalUsers"`
	TotalUsersReal     uint32       `json:"TotalUsersReal"`
	TotalDeals         uint32       `json:"TotalDeals"`
	TotalOrders        uint32       `json:"TotalOrders"`
	TotalOrdersHistory uint32       `json:"TotalOrdersHistory"`
	TotalPositions     uint32       `json:"TotalPositions"`
	AccountURL         string       `json:"AccountURL"`
	AccountAuto        uint32       `json:"AccountAuto"`
}
