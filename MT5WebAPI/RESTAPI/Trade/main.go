package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("TRADE_CALC", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Trade Calculation Example (Go)")
	fmt.Println("  Tests: 27-30 (4 trade calc functions)")
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

	// =============================================
	// 27. CalcRateBuy - Calculate buy conversion rates
	// =============================================
	fmt.Println("--- 27. CalcRateBuy ---")

	// EUR -> USD buy rate (should be ~1.17)
	rate, ret := api.CalcRateBuy("EUR", "USD")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  EUR/USD buy rate: %.6f\n", rate)
	} else {
		fmt.Printf("  EUR/USD buy rate FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	// GBP -> USD buy rate (should be ~1.35)
	rate, ret = api.CalcRateBuy("GBP", "USD")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  GBP/USD buy rate: %.6f\n", rate)
	} else {
		fmt.Printf("  GBP/USD buy rate FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	// USD -> JPY buy rate (should be ~156)
	rate, ret = api.CalcRateBuy("USD", "JPY")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  USD/JPY buy rate: %.6f\n", rate)
	} else {
		fmt.Printf("  USD/JPY buy rate FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 28. CalcRateSell - Calculate sell conversion rates
	// =============================================
	fmt.Println("--- 28. CalcRateSell ---")

	// EUR -> USD sell rate
	rate, ret = api.CalcRateSell("EUR", "USD")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  EUR/USD sell rate: %.6f\n", rate)
	} else {
		fmt.Printf("  EUR/USD sell rate FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	// GBP -> USD sell rate
	rate, ret = api.CalcRateSell("GBP", "USD")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  GBP/USD sell rate: %.6f\n", rate)
	} else {
		fmt.Printf("  GBP/USD sell rate FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 29. CheckMargin - Check margin for orders
	// =============================================
	fmt.Println("--- 29. CheckMargin ---")

	// Get current prices for EURUSD
	ticks, ret := api.TickLast("EURUSD")
	if ret != mt5webapi.MT_RET_OK || len(ticks) == 0 {
		fmt.Printf("  Failed to get EURUSD ticks: [%d] %s\n", ret, mt5webapi.GetError(ret))
		fmt.Println("  Skipping margin check tests")
	} else {
		ask := ticks[0].Ask
		bid := ticks[0].Bid
		fmt.Printf("  EURUSD: Bid=%.5f  Ask=%.5f\n", bid, ask)

		// 0.1 lot buy at ask price (login 5000)
		marginResult, ret := api.CheckMargin(5000, "EURUSD", 0, 10000, ask)
		if ret == mt5webapi.MT_RET_OK && marginResult != nil {
			fmt.Println("  CheckMargin 0.1 lot BUY:")
			if newState, ok := marginResult["new"]; ok {
				fmt.Printf("    new state: %v\n", newState)
			}
			if curState, ok := marginResult["current"]; ok {
				fmt.Printf("    current state: %v\n", curState)
			}
		} else {
			fmt.Printf("  CheckMargin 0.1 lot BUY FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}

		// 1.0 lot sell at bid price (login 5000)
		marginResult, ret = api.CheckMargin(5000, "EURUSD", 1, 100000, bid)
		if ret == mt5webapi.MT_RET_OK && marginResult != nil {
			fmt.Println("  CheckMargin 1.0 lot SELL:")
			if newState, ok := marginResult["new"]; ok {
				fmt.Printf("    new state: %v\n", newState)
			}
			if curState, ok := marginResult["current"]; ok {
				fmt.Printf("    current state: %v\n", curState)
			}
		} else {
			fmt.Printf("  CheckMargin 1.0 lot SELL FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	}
	fmt.Println()

	// =============================================
	// 30. CalcProfit - Calculate profit
	// =============================================
	fmt.Println("--- 30. CalcProfit ---")

	// Buy profit: open at 1.16000, close at 1.17000 (0.1 lot)
	profit, ret := api.CalcProfit("demo\\forex", "EURUSD", 0, 10000, 1.16000, 1.17000)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Buy profit (0.1 lot, 1.16000->1.17000): %.2f\n", profit)
	} else {
		fmt.Printf("  Buy profit FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	// Sell profit: open at 1.17000, close at 1.16000 (0.1 lot)
	profit, ret = api.CalcProfit("demo\\forex", "EURUSD", 1, 10000, 1.17000, 1.16000)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Sell profit (0.1 lot, 1.17000->1.16000): %.2f\n", profit)
	} else {
		fmt.Printf("  Sell profit FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All trade calculation tests completed!")
	fmt.Println("================================================")
}
