package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("DAILY", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Daily Report Example")
	fmt.Println("  Tests: 37-38 (DailyGet, DailyGetBatch)")
	fmt.Println("================================================")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		return
	}

	fmt.Print("Connecting... ")
	ret := api.ConnectDefault(cfg.Server, cfg.Port, cfg.Login, cfg.Password, mt5webapi.PUMP_MODE_NONE)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("FAILED! [%d] %s\n", ret, mt5webapi.GetError(ret))
		return
	}
	fmt.Println("OK")
	fmt.Println()

	now := time.Now().Unix() + 86400
	from := time.Now().AddDate(0, -6, 0).Unix() // last 6 months



	// =============================================
	// 37. DAILY_GET - Get daily reports for login 5000 (last 7 days)
	// =============================================
	fmt.Println("--- 37. DailyGet (login 5000, last 7 days) ---")
	dailies, ret := api.DailyGet(5000, from, now)
	if ret == mt5webapi.MT_RET_OK && dailies != nil {
		fmt.Printf("  Found %d daily reports for login 5000:\n", len(dailies))
		for _, d := range dailies {
			fmt.Printf("    Login=%d Name=%s Group=%s Balance=%.2f Equity=%.2f Profit=%.2f Margin=%.2f\n",
				d.Login, d.Name, d.Group, d.Balance, d.Equity, d.Profit, d.Margin)
		}
	} else {
		fmt.Printf("  FAILED or no data: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 37b. DAILY_GET - Get daily reports for login 5005 (last 7 days)
	// =============================================
	fmt.Println("--- 37b. DailyGet (login 5005, last 7 days) ---")
	dailies, ret = api.DailyGet(5005, from, now)
	if ret == mt5webapi.MT_RET_OK && dailies != nil {
		fmt.Printf("  Found %d daily reports for login 5005:\n", len(dailies))
		for _, d := range dailies {
			fmt.Printf("    Login=%d Name=%s Group=%s Balance=%.2f Equity=%.2f Profit=%.2f Margin=%.2f\n",
				d.Login, d.Name, d.Group, d.Balance, d.Equity, d.Profit, d.Margin)
		}
	} else {
		fmt.Printf("  FAILED or no data: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 38. DAILY_GET_BATCH - Get daily reports for logins 5000,5005 (last 7 days)
	// =============================================
	fmt.Println("--- 38. DailyGetBatch (logins 5000,5005) ---")
	dailies, ret = api.DailyGetBatch("5000,5005", from, now)
	if ret == mt5webapi.MT_RET_OK && dailies != nil {
		fmt.Printf("  Found %d daily reports:\n", len(dailies))
		for _, d := range dailies {
			fmt.Printf("    Login=%d Name=%s Group=%s Balance=%.2f Equity=%.2f Profit=%.2f Margin=%.2f\n",
				d.Login, d.Name, d.Group, d.Balance, d.Equity, d.Profit, d.Margin)
		}
	} else {
		fmt.Printf("  FAILED or no data: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All daily report tests completed!")
	fmt.Println("================================================")
}
