package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("POSITION", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Position Example (Go)")
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

	// CommonGet says TotalPositions=1, find which login has it
	loginsToCheck := []uint64{2000, 5005}
	demoLogins, ret := api.UserLogins("demo\\*")
	if ret == mt5webapi.MT_RET_OK {
		loginsToCheck = append(loginsToCheck, demoLogins...)
	}

	// =============================================
	// 21. POSITION_GET_TOTAL - Toplam pozisyon sayısı
	// =============================================
	fmt.Println("--- 21. PositionGetTotal ---")
	var loginWithPositions uint64
	var totalPositions uint32

	for _, login := range loginsToCheck {
		total, ret := api.PositionGetTotal(login)
		if ret == mt5webapi.MT_RET_OK {
			fmt.Printf("Login %d: %d open positions\n", login, total)
			if total > 0 && loginWithPositions == 0 {
				loginWithPositions = login
				totalPositions = total
			}
		} else {
			fmt.Printf("Login %d: failed - %s\n", login, mt5webapi.GetError(ret))
		}
	}
	fmt.Println()

	if loginWithPositions == 0 {
		fmt.Println("No open positions found on any account.")
		fmt.Println("Usage examples:")
		fmt.Println("  total, ret := api.PositionGetTotal(login)")
		fmt.Println("  positions, ret := api.PositionGetPage(login, 0, total)")
		fmt.Println("  position, ret := api.PositionGet(login, \"EURUSD\")")
		fmt.Println()
		api.Disconnect()
		fmt.Println("Disconnected.")
		fmt.Println("========================================")
		return
	}

	fmt.Printf("Using login %d with %d open positions\n", loginWithPositions, totalPositions)
	fmt.Println()

	// =============================================
	// 22. POSITION_GET_PAGE - Pozisyon listesi
	// =============================================
	fmt.Println("--- 22. PositionGetPage ---")
	pageSize := uint32(10)
	if totalPositions < pageSize {
		pageSize = totalPositions
	}
	positions, ret := api.PositionGetPage(loginWithPositions, 0, pageSize)
	if ret == mt5webapi.MT_RET_OK && positions != nil {
		fmt.Printf("Got %d positions (of %d total):\n", len(positions), totalPositions)
		for i, p := range positions {
			fmt.Printf("\n  Position #%d:\n", i+1)
			fmt.Printf("    Ticket:       %d\n", p.Position)
			fmt.Printf("    Symbol:       %s\n", p.Symbol)
			fmt.Printf("    Action:       %s\n", posActionStr(p.Action))
			fmt.Printf("    Volume:       %d\n", p.VolumeExt)
			fmt.Printf("    PriceOpen:    %.5f\n", p.PriceOpen)
			fmt.Printf("    PriceCurrent: %.5f\n", p.PriceCurrent)
			fmt.Printf("    PriceSL:      %.5f\n", p.PriceSL)
			fmt.Printf("    PriceTP:      %.5f\n", p.PriceTP)
			fmt.Printf("    Profit:       %.2f\n", p.Profit)
			fmt.Printf("    Storage:      %.2f\n", p.Storage)
			fmt.Printf("    Commission:   %.2f\n", p.Commission)
			fmt.Printf("    Comment:      %s\n", p.Comment)

			// =============================================
			// 20. POSITION_GET - Pozisyon detayı
			// =============================================
			if i == 0 {
				fmt.Printf("\n--- 20. PositionGet (login=%d, symbol=%s) ---\n", loginWithPositions, p.Symbol)
				pos, ret := api.PositionGet(loginWithPositions, p.Symbol)
				if ret == mt5webapi.MT_RET_OK && pos != nil {
					fmt.Printf("    Ticket:         %d\n", pos.Position)
					fmt.Printf("    ExternalID:     %s\n", pos.ExternalID)
					fmt.Printf("    Login:          %d\n", pos.Login)
					fmt.Printf("    Dealer:         %d\n", pos.Dealer)
					fmt.Printf("    Symbol:         %s\n", pos.Symbol)
					fmt.Printf("    Action:         %s\n", posActionStr(pos.Action))
					fmt.Printf("    Digits:         %d\n", pos.Digits)
					fmt.Printf("    Reason:         %d\n", pos.Reason)
					fmt.Printf("    ContractSize:   %.2f\n", pos.ContractSize)
					fmt.Printf("    TimeCreate:     %d\n", pos.TimeCreate)
					fmt.Printf("    TimeUpdate:     %d\n", pos.TimeUpdate)
					fmt.Printf("    PriceOpen:      %.5f\n", pos.PriceOpen)
					fmt.Printf("    PriceCurrent:   %.5f\n", pos.PriceCurrent)
					fmt.Printf("    PriceSL:        %.5f\n", pos.PriceSL)
					fmt.Printf("    PriceTP:        %.5f\n", pos.PriceTP)
					fmt.Printf("    Volume:         %d\n", pos.VolumeExt)
					fmt.Printf("    Profit:         %.2f\n", pos.Profit)
					fmt.Printf("    Storage:        %.2f\n", pos.Storage)
					fmt.Printf("    Commission:     %.2f\n", pos.Commission)
					fmt.Printf("    RateProfit:     %.5f\n", pos.RateProfit)
					fmt.Printf("    RateMargin:     %.5f\n", pos.RateMargin)
					fmt.Printf("    ExpertID:       %d\n", pos.ExpertID)
					fmt.Printf("    ExpertPositionID:%d\n", pos.ExpertPositionID)
					fmt.Printf("    Comment:        %s\n", pos.Comment)
					fmt.Printf("    ActivationMode: %d\n", pos.ActivationMode)
				} else {
					fmt.Printf("    PositionGet failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
				}
			}
		}
	} else {
		fmt.Printf("PositionGetPage failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All position tests completed!")
	fmt.Println("========================================")
}

func posActionStr(a mt5webapi.EnPositionAction) string {
	switch a {
	case mt5webapi.PositionBuy:
		return "Buy"
	case mt5webapi.PositionSell:
		return "Sell"
	}
	return fmt.Sprintf("Unknown(%d)", a)
}
