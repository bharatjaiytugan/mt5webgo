package mt5webapi

import (
	"fmt"
	"strconv"
	"strings"
)

// EnUsersRights - client rights bit flags.
type EnUsersRights uint64

const (
	UserRightNone             EnUsersRights = 0x0000000000000000
	UserRightEnabled          EnUsersRights = 0x0000000000000001
	UserRightPassword         EnUsersRights = 0x0000000000000002
	UserRightTradeDisabled    EnUsersRights = 0x0000000000000004
	UserRightInvestor         EnUsersRights = 0x0000000000000008
	UserRightConfirmed        EnUsersRights = 0x0000000000000010
	UserRightTrailing         EnUsersRights = 0x0000000000000020
	UserRightExpert           EnUsersRights = 0x0000000000000040
	UserRightObsolete         EnUsersRights = 0x0000000000000080
	UserRightReports          EnUsersRights = 0x0000000000000100
	UserRightReadonly         EnUsersRights = 0x0000000000000200
	UserRightResetPass        EnUsersRights = 0x0000000000000400
	UserRightOTPEnabled       EnUsersRights = 0x0000000000000800
	UserRightClientConfirmed  EnUsersRights = 0x0000000000001000
	UserRightSponsoredHosting EnUsersRights = 0x0000000000002000
	UserRightAPIEnabled       EnUsersRights = 0x0000000000004000

	UserRightDefault EnUsersRights = UserRightEnabled | UserRightPassword | UserRightTrailing | UserRightExpert | UserRightReports
	UserRightAll     EnUsersRights = UserRightEnabled | UserRightPassword | UserRightTradeDisabled |
		UserRightInvestor | UserRightConfirmed | UserRightTrailing |
		UserRightExpert | UserRightReports |
		UserRightReadonly | UserRightResetPass | UserRightOTPEnabled | UserRightClientConfirmed | UserRightSponsoredHosting
)

// EnUsersPasswords - password types.
type EnUsersPasswords uint32

const (
	UserPassMain     EnUsersPasswords = 0
	UserPassInvestor EnUsersPasswords = 1
	UserPassAPI      EnUsersPasswords = 2

	UserPassFirst EnUsersPasswords = UserPassMain
	UserPassLast  EnUsersPasswords = UserPassAPI
)

// EnUsersConnectionTypes - connection types.
type EnUsersConnectionTypes uint32

const (
	// client types
	UserTypeClient      EnUsersConnectionTypes = 0
	UserTypeClientWinMobile EnUsersConnectionTypes = 1
	UserTypeClientWinPhone  EnUsersConnectionTypes = 2
	UserTypeClientIPhone    EnUsersConnectionTypes = 4
	UserTypeClientAndroid   EnUsersConnectionTypes = 5
	UserTypeClientBlackBerry EnUsersConnectionTypes = 6
	UserTypeClientWeb       EnUsersConnectionTypes = 11
	// manager types
	UserTypeAdmin      EnUsersConnectionTypes = 32
	UserTypeManager    EnUsersConnectionTypes = 33
	UserTypeManagerAPI EnUsersConnectionTypes = 34
	UserTypeAdminAPI   EnUsersConnectionTypes = 36

	UserTypeFirst EnUsersConnectionTypes = UserTypeClient
	UserTypeLast  EnUsersConnectionTypes = UserTypeAdminAPI
)

const (
	externalIDMaxLen = 32
	externalIDLimit  = 128
)

// MTUser represents user data.
type MTUser struct {
	Login                 uint64        `json:"Login"`
	Group                 string        `json:"Group"`
	CertSerialNumber      uint64        `json:"CertSerialNumber"`
	Rights                EnUsersRights `json:"Rights"`
	MQID                  string        `json:"MQID"`
	Registration          int64         `json:"Registration"`
	LastAccess            int64         `json:"LastAccess"`
	LastPassChange        int64         `json:"LastPassChange"`
	LastIP                string        `json:"LastIP"`
	Name                  string        `json:"Name"`
	Company               string        `json:"Company"`
	Account               string        `json:"Account"`
	Country               string        `json:"Country"`
	Language              uint32        `json:"Language"`
	ClientID              uint64        `json:"ClientID"`
	City                  string        `json:"City"`
	State                 string        `json:"State"`
	ZIPCode               string        `json:"ZIPCode"`
	Address               string        `json:"Address"`
	Phone                 string        `json:"Phone"`
	Email                 string        `json:"Email"`
	ID                    string        `json:"ID"`
	Status                string        `json:"Status"`
	Comment               string        `json:"Comment"`
	Color                 uint32        `json:"Color"`
	MainPassword          string        `json:"MainPassword"`
	InvestPassword        string        `json:"InvestPassword"`
	PhonePassword         string        `json:"PhonePassword"`
	Leverage              uint32        `json:"Leverage"`
	Agent                 uint64        `json:"Agent"`
	Balance               float64       `json:"Balance"`
	Credit                float64       `json:"Credit"`
	InterestRate          float64       `json:"InterestRate"`
	CommissionDaily       float64       `json:"CommissionDaily"`
	CommissionMonthly     float64       `json:"CommissionMonthly"`
	CommissionAgentDaily  float64       `json:"CommissionAgentDaily"`
	CommissionAgentMonthly float64      `json:"CommissionAgentMonthly"`
	BalancePrevDay        float64       `json:"BalancePrevDay"`
	BalancePrevMonth      float64       `json:"BalancePrevMonth"`
	EquityPrevDay         float64       `json:"EquityPrevDay"`
	EquityPrevMonth       float64       `json:"EquityPrevMonth"`
	TradeAccounts         string        `json:"TradeAccounts"`
	LeadCampaign          string        `json:"LeadCampaign"`
	LeadSource            string        `json:"LeadSource"`
}

// CreateDefault creates a new MTUser with default values.
func CreateDefault() *MTUser {
	return &MTUser{
		Rights:   UserRightDefault,
		Leverage: 100,
		Color:    0xFF000000,
	}
}

// ExternalAccountAdd adds an external account.
func (u *MTUser) ExternalAccountAdd(gatewayID uint64, account string) MTRetCode {
	if len(account) >= externalIDMaxLen {
		return MT_RET_ERR_DATA
	}
	result := fmt.Sprintf("%s%d=%s|", u.TradeAccounts, gatewayID, account)
	if externalIDLimit <= len(result) {
		return MT_RET_ERR_DATA
	}
	u.TradeAccounts = result
	return MT_RET_OK
}

// ExternalAccountUpdate updates an external account at the given position.
func (u *MTUser) ExternalAccountUpdate(pos uint32, gatewayID uint64, account string) MTRetCode {
	if len(account) >= externalIDMaxLen {
		return MT_RET_ERR_DATA
	}
	tokens := strings.Split(u.TradeAccounts, "|")
	var result strings.Builder
	var count uint32
	for _, token := range tokens {
		if len(token) < 1 {
			continue
		}
		if pos == count {
			fmt.Fprintf(&result, "%d=%s|", gatewayID, account)
		} else {
			fmt.Fprintf(&result, "%s|", token)
		}
		count++
	}
	if pos >= count {
		return MT_RET_ERR_PARAMS
	}
	if externalIDLimit <= result.Len() {
		return MT_RET_ERR_DATA
	}
	u.TradeAccounts = result.String()
	return MT_RET_OK
}

// ExternalAccountDelete deletes an external account at the given position.
func (u *MTUser) ExternalAccountDelete(pos uint32) MTRetCode {
	tokens := strings.Split(u.TradeAccounts, "|")
	var result strings.Builder
	var count uint32
	for _, token := range tokens {
		if len(token) < 1 {
			continue
		}
		if pos != count {
			fmt.Fprintf(&result, "%s|", token)
		}
		count++
	}
	if pos >= count {
		return MT_RET_ERR_PARAMS
	}
	if externalIDLimit <= result.Len() {
		return MT_RET_ERR_DATA
	}
	u.TradeAccounts = result.String()
	return MT_RET_OK
}

// ExternalAccountClear clears all external accounts.
func (u *MTUser) ExternalAccountClear() MTRetCode {
	u.TradeAccounts = ""
	return MT_RET_OK
}

// ExternalAccountTotal returns the total count of external accounts.
func (u *MTUser) ExternalAccountTotal() uint32 {
	tokens := strings.Split(u.TradeAccounts, "|")
	var count uint32
	for _, token := range tokens {
		if len(token) < 1 {
			continue
		}
		count++
	}
	return count
}

// ExternalAccountNext gets the external account at the given position.
func (u *MTUser) ExternalAccountNext(pos uint32) (gatewayID uint64, account string, ret MTRetCode) {
	tokens := strings.Split(u.TradeAccounts, "|")
	var count uint32
	for _, token := range tokens {
		if len(token) < 1 {
			continue
		}
		if pos == count {
			items := strings.Split(token, "=")
			if len(items) == 2 {
				id, err := strconv.ParseUint(items[0], 10, 64)
				if err != nil {
					id = 0
				}
				return id, items[1], MT_RET_OK
			}
			break
		}
		count++
	}
	return 0, "", MT_RET_ERR_PARAMS
}

// ExternalAccountGet finds an external account for the given gateway.
func (u *MTUser) ExternalAccountGet(gatewayID uint64) (account string, ret MTRetCode) {
	tokens := strings.Split(u.TradeAccounts, "|")
	for _, token := range tokens {
		if len(token) < 1 {
			continue
		}
		items := strings.Split(token, "=")
		if len(items) == 2 {
			id, err := strconv.ParseUint(items[0], 10, 64)
			if err == nil && id == gatewayID {
				return items[1], MT_RET_OK
			}
		}
	}
	return "", MT_RET_ERR_PARAMS
}
