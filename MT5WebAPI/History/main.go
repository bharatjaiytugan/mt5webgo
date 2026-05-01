package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("HISTORY", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - History Example (Go)")
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

	// check multiple logins for history orders
	loginsToCheck := []uint64{5000, 2000, 5005}
	demoLogins, _ := api.UserLogins("demo\\*")
	loginsToCheck = append(loginsToCheck, demoLogins...)

	// time range: last year
	from := time.Now().AddDate(-1, 0, 0).Unix()
	to := time.Now().Unix() + 86400

	// find a login with history
	var testLogin uint64

	// =============================================
	// 27. HISTORY_GET_TOTAL - Geçmiş emir sayısı
	// =============================================
	fmt.Println("--- 27. HistoryGetTotal ---")
	var total uint32
	for _, login := range loginsToCheck {
		t, ret := api.HistoryGetTotal(login, from, to)
		if ret == mt5webapi.MT_RET_OK {
			fmt.Printf("Login %d: %d history orders\n", login, t)
			if t > 0 && testLogin == 0 {
				testLogin = login
				total = t
			}
		} else {
			fmt.Printf("Login %d: failed - %s\n", login, mt5webapi.GetError(ret))
		}
	}
	if testLogin == 0 {
		fmt.Println("\nNo history orders found on any login.")
		fmt.Println("Usage:")
		fmt.Println("  total, ret := api.HistoryGetTotal(login, from, to)")
		fmt.Println("  orders, ret := api.HistoryGetPage(login, from, to, 0, total)")
		fmt.Println("  order, ret := api.HistoryGet(ticket)")
		api.Disconnect()
		fmt.Println("\nDisconnected.")
		fmt.Println("========================================")
		return
	}
	fmt.Printf("\nUsing login %d with %d history orders\n", testLogin, total)
	fmt.Println()

	// =============================================
	// 28. HISTORY_GET_PAGE - Geçmiş emir listesi
	// =============================================
	fmt.Println("--- 28. HistoryGetPage ---")
	if total > 0 {
		pageSize := uint32(10)
		if total < pageSize {
			pageSize = total
		}
		orders, ret := api.HistoryGetPage(testLogin, from, to, 0, pageSize)
		if ret == mt5webapi.MT_RET_OK && orders != nil {
			fmt.Printf("Got %d history orders (of %d total):\n", len(orders), total)
			for i, o := range orders {
				fmt.Printf("\n  History Order #%d:\n", i+1)
				fmt.Printf("    Ticket:      %d\n", o.Order)
				fmt.Printf("    Symbol:      %s\n", o.Symbol)
				fmt.Printf("    Type:        %s\n", orderTypeStr(o.Type))
				fmt.Printf("    State:       %s\n", orderStateStr(o.State))
				fmt.Printf("    Volume:      %d\n", o.VolumeInitialExt)
				fmt.Printf("    Price:       %.5f\n", o.PriceOrder)
				fmt.Printf("    Comment:     %s\n", o.Comment)
				fmt.Printf("    TimeSetup:   %s\n", fmtTime(o.TimeSetup))
				fmt.Printf("    TimeDone:    %s\n", fmtTime(o.TimeDone))

				// =============================================
				// 26. HISTORY_GET - Geçmiş emir detayı (ilk order)
				// =============================================
				if i == 0 {
					fmt.Printf("\n--- 26. HistoryGet (ticket=%d) ---\n", o.Order)
					histOrder, ret := api.HistoryGet(o.Order)
					if ret == mt5webapi.MT_RET_OK && histOrder != nil {
						fmt.Printf("    Ticket:       %d\n", histOrder.Order)
						fmt.Printf("    ExternalID:   %s\n", histOrder.ExternalID)
						fmt.Printf("    Login:        %d\n", histOrder.Login)
						fmt.Printf("    Dealer:       %d\n", histOrder.Dealer)
						fmt.Printf("    Symbol:       %s\n", histOrder.Symbol)
						fmt.Printf("    Digits:       %d\n", histOrder.Digits)
						fmt.Printf("    Type:         %s\n", orderTypeStr(histOrder.Type))
						fmt.Printf("    TypeFill:     %d\n", histOrder.TypeFill)
						fmt.Printf("    TypeTime:     %d\n", histOrder.TypeTime)
						fmt.Printf("    State:        %s\n", orderStateStr(histOrder.State))
						fmt.Printf("    Reason:       %d\n", histOrder.Reason)
						fmt.Printf("    PriceOrder:   %.5f\n", histOrder.PriceOrder)
						fmt.Printf("    PriceCurrent: %.5f\n", histOrder.PriceCurrent)
						fmt.Printf("    PriceSL:      %.5f\n", histOrder.PriceSL)
						fmt.Printf("    PriceTP:      %.5f\n", histOrder.PriceTP)
						fmt.Printf("    VolumeInit:   %d\n", histOrder.VolumeInitialExt)
						fmt.Printf("    VolumeCurr:   %d\n", histOrder.VolumeCurrentExt)
						fmt.Printf("    ExpertID:     %d\n", histOrder.ExpertID)
						fmt.Printf("    ContractSize: %.2f\n", histOrder.ContractSize)
						fmt.Printf("    Comment:      %s\n", histOrder.Comment)
						fmt.Printf("    TimeSetup:    %s\n", fmtTime(histOrder.TimeSetup))
						fmt.Printf("    TimeDone:     %s\n", fmtTime(histOrder.TimeDone))
					} else {
						fmt.Printf("    HistoryGet failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
					}
				}
			}
		} else {
			fmt.Printf("HistoryGetPage failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	} else {
		fmt.Println("No history orders found.")
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All history tests completed!")
	fmt.Println("========================================")
}

func fmtTime(unix int64) string {
	if unix == 0 {
		return "N/A"
	}
	return time.Unix(unix, 0).UTC().Format("2006-01-02 15:04:05")
}

func orderTypeStr(t mt5webapi.EnOrderType) string {
	switch t {
	case mt5webapi.OpBuy:
		return "Buy"
	case mt5webapi.OpSell:
		return "Sell"
	case mt5webapi.OpBuyLimit:
		return "Buy Limit"
	case mt5webapi.OpSellLimit:
		return "Sell Limit"
	case mt5webapi.OpBuyStop:
		return "Buy Stop"
	case mt5webapi.OpSellStop:
		return "Sell Stop"
	case mt5webapi.OpBuyStopLimit:
		return "Buy Stop Limit"
	case mt5webapi.OpSellStopLimit:
		return "Sell Stop Limit"
	case mt5webapi.OpCloseBy:
		return "Close By"
	}
	return fmt.Sprintf("Type(%d)", t)
}

func orderStateStr(s mt5webapi.EnOrderState) string {
	switch s {
	case mt5webapi.OrderStateStarted:
		return "Started"
	case mt5webapi.OrderStatePlaced:
		return "Placed"
	case mt5webapi.OrderStateCanceled:
		return "Canceled"
	case mt5webapi.OrderStatePartial:
		return "Partial"
	case mt5webapi.OrderStateFilled:
		return "Filled"
	case mt5webapi.OrderStateRejected:
		return "Rejected"
	case mt5webapi.OrderStateExpired:
		return "Expired"
	}
	return fmt.Sprintf("State(%d)", s)
}
