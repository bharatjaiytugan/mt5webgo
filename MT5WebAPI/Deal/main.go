package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("DEAL", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Deal Example (Go)")
	fmt.Println("========================================")
	fmt.Println()

	// load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		return
	}

	// connect
	fmt.Print("Connecting... ")
	ret := api.ConnectDefault(cfg.Server, cfg.Port, cfg.Login, cfg.Password, mt5webapi.PUMP_MODE_NONE)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("FAILED! [%d] %s\n", ret, mt5webapi.GetError(ret))
		return
	}
	fmt.Println("OK")
	fmt.Println()

	// create a test user and make a deposit to generate a deal
	fmt.Println("Creating test user 99999 with deposit to generate deals...")
	testLogin := uint64(99999)
	newUser := mt5webapi.CreateDefault()
	newUser.Login = testLogin
	newUser.Name = "Deal Test User"
	newUser.Group = "demo\\forex"
	newUser.Leverage = 100
	newUser.MainPassword = "Test1234!"
	newUser.InvestPassword = "Invest1234!"

	_, ret = api.UserAdd(newUser)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("UserAdd failed: %s\n", mt5webapi.GetError(ret))
		api.Disconnect()
		return
	}

	// deposit to create deals
	api.UserDepositChange(testLogin, 5000.00, "Deposit #1", mt5webapi.DealBalance)
	api.UserDepositChange(testLogin, 2000.00, "Deposit #2", mt5webapi.DealBalance)
	api.UserDepositChange(testLogin, -500.00, "Withdrawal", mt5webapi.DealBalance)
	api.UserDepositChange(testLogin, 100.00, "Credit bonus", mt5webapi.DealCredit)
	fmt.Println("4 balance operations done.")
	fmt.Println()

	// time range: last 30 days to now
	from := time.Now().AddDate(0, 0, -30).Unix()
	to := time.Now().Unix() + 86400

	// =============================================
	// 24. DEAL_GET_TOTAL - Tarih aralığında işlem sayısı
	// =============================================
	fmt.Println("--- 24. DealGetTotal ---")
	total, ret := api.DealGetTotal(testLogin, from, to)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("Login %d: %d deals (all time)\n", testLogin, total)
	} else {
		fmt.Printf("DealGetTotal failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 25. DEAL_GET_PAGE - İşlem listesi (sayfalı)
	// =============================================
	fmt.Println("--- 25. DealGetPage ---")
	if total > 0 {
		pageSize := uint32(10)
		if total < pageSize {
			pageSize = total
		}
		deals, ret := api.DealGetPage(testLogin, from, to, 0, pageSize)
		if ret == mt5webapi.MT_RET_OK && deals != nil {
			fmt.Printf("Got %d deals (of %d total):\n", len(deals), total)
			for i, d := range deals {
				fmt.Printf("\n  Deal #%d:\n", i+1)
				fmt.Printf("    Ticket:     %d\n", d.Deal)
				fmt.Printf("    Login:      %d\n", d.Login)
				fmt.Printf("    Action:     %s\n", dealActionStr(d.Action))
				fmt.Printf("    Entry:      %s\n", dealEntryStr(d.Entry))
				fmt.Printf("    Symbol:     %s\n", d.Symbol)
				fmt.Printf("    Volume:     %d\n", d.VolumeExt)
				fmt.Printf("    Price:      %.5f\n", d.Price)
				fmt.Printf("    Profit:     %.2f\n", d.Profit)
				fmt.Printf("    Commission: %.2f\n", d.Commission)
				fmt.Printf("    Storage:    %.2f\n", d.Storage)
				fmt.Printf("    Comment:    %s\n", d.Comment)
				fmt.Printf("    Time:       %s\n", time.Unix(d.Time, 0).UTC().Format("2006-01-02 15:04:05"))

				// =============================================
				// 23. DEAL_GET - İşlem detayı (ilk deal için)
				// =============================================
				if i == 0 {
					fmt.Printf("\n--- 23. DealGet (ticket=%d) ---\n", d.Deal)
					deal, ret := api.DealGet(d.Deal)
					if ret == mt5webapi.MT_RET_OK && deal != nil {
						fmt.Printf("    Ticket:       %d\n", deal.Deal)
						fmt.Printf("    ExternalID:   %s\n", deal.ExternalID)
						fmt.Printf("    Login:        %d\n", deal.Login)
						fmt.Printf("    Dealer:       %d\n", deal.Dealer)
						fmt.Printf("    Order:        %d\n", deal.Order)
						fmt.Printf("    Action:       %s\n", dealActionStr(deal.Action))
						fmt.Printf("    Entry:        %s\n", dealEntryStr(deal.Entry))
						fmt.Printf("    Reason:       %d\n", deal.Reason)
						fmt.Printf("    Digits:       %d\n", deal.Digits)
						fmt.Printf("    Symbol:       %s\n", deal.Symbol)
						fmt.Printf("    Price:        %.5f\n", deal.Price)
						fmt.Printf("    Volume:       %d\n", deal.VolumeExt)
						fmt.Printf("    Profit:       %.2f\n", deal.Profit)
						fmt.Printf("    Storage:      %.2f\n", deal.Storage)
						fmt.Printf("    Commission:   %.2f\n", deal.Commission)
						fmt.Printf("    RateProfit:   %.5f\n", deal.RateProfit)
						fmt.Printf("    RateMargin:   %.5f\n", deal.RateMargin)
						fmt.Printf("    ExpertID:     %d\n", deal.ExpertID)
						fmt.Printf("    PositionID:   %d\n", deal.PositionID)
						fmt.Printf("    Comment:      %s\n", deal.Comment)
						fmt.Printf("    PriceSL:      %.5f\n", deal.PriceSL)
						fmt.Printf("    PriceTP:      %.5f\n", deal.PriceTP)
						fmt.Printf("    Flags:        %d\n", deal.Flags)
						fmt.Printf("    Gateway:      %s\n", deal.Gateway)
						fmt.Printf("    Time:         %s\n", time.Unix(deal.Time, 0).UTC().Format("2006-01-02 15:04:05"))
					} else {
						fmt.Printf("    DealGet failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
					}
				}
			}
		} else {
			fmt.Printf("DealGetPage failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	} else {
		fmt.Println("No deals found.")
	}
	fmt.Println()

	// cleanup: delete test user
	fmt.Print("Cleaning up... ")
	ret = api.UserDelete(testLogin)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("User %d deleted.\n", testLogin)
	} else {
		fmt.Printf("Delete failed: %s\n", mt5webapi.GetError(ret))
	}

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All deal tests completed!")
	fmt.Println("========================================")
}

func dealActionStr(a mt5webapi.EnDealAction) string {
	switch a {
	case mt5webapi.DealBuy:
		return "Buy"
	case mt5webapi.DealSell:
		return "Sell"
	case mt5webapi.DealBalance:
		return "Balance"
	case mt5webapi.DealCredit:
		return "Credit"
	case mt5webapi.DealCharge:
		return "Charge"
	case mt5webapi.DealCorrection:
		return "Correction"
	case mt5webapi.DealBonus:
		return "Bonus"
	case mt5webapi.DealCommission:
		return "Commission"
	case mt5webapi.DealCommissionDaily:
		return "Commission Daily"
	case mt5webapi.DealCommissionMonthly:
		return "Commission Monthly"
	case mt5webapi.DealInterestRate:
		return "Interest Rate"
	case mt5webapi.DealDividend:
		return "Dividend"
	case mt5webapi.DealTax:
		return "Tax"
	}
	return fmt.Sprintf("Action(%d)", a)
}

func dealEntryStr(e mt5webapi.EnEntryFlags) string {
	switch e {
	case mt5webapi.EntryIn:
		return "In"
	case mt5webapi.EntryOut:
		return "Out"
	case mt5webapi.EntryInOut:
		return "InOut"
	case mt5webapi.EntryOutBy:
		return "OutBy"
	}
	return fmt.Sprintf("Entry(%d)", e)
}
