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

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Position Batch Example (Go)")
	fmt.Println("  Test: 26. PositionGetBatch")
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
	// 26. POSITION_GET_BATCH - by login
	// =============================================
	fmt.Println("--- 26a. PositionGetBatch (login=5000) ---")
	positions, ret := api.PositionGetBatch("5000", "", "")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Found %d positions:\n", len(positions))
		for _, p := range positions {
			action := "Buy"
			if p.Action == mt5webapi.PositionSell {
				action = "Sell"
			}
			fmt.Printf("    ticket=%d %s %s vol=%d open=%.5f current=%.5f profit=%.2f\n",
				p.Position, action, p.Symbol, p.VolumeExt, p.PriceOpen, p.PriceCurrent, p.Profit)
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 26b. POSITION_GET_BATCH - by group
	// =============================================
	fmt.Println("--- 26b. PositionGetBatch (group=demo\\*) ---")
	positions2, ret := api.PositionGetBatch("", "demo\\*", "")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Found %d positions in demo groups:\n", len(positions2))
		for i, p := range positions2 {
			if i >= 10 {
				fmt.Printf("    ... +%d more\n", len(positions2)-10)
				break
			}
			action := "Buy"
			if p.Action == mt5webapi.PositionSell {
				action = "Sell"
			}
			fmt.Printf("    login=%d ticket=%d %s %s profit=%.2f\n",
				p.Login, p.Position, action, p.Symbol, p.Profit)
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 26c. POSITION_GET_BATCH - all (managers\*)
	// =============================================
	fmt.Println("--- 26c. PositionGetBatch (group=managers\\*) ---")
	positions3, ret := api.PositionGetBatch("", "managers\\*", "")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Found %d positions in managers groups:\n", len(positions3))
		for _, p := range positions3 {
			action := "Buy"
			if p.Action == mt5webapi.PositionSell {
				action = "Sell"
			}
			fmt.Printf("    login=%d ticket=%d %s %s open=%.5f profit=%.2f\n",
				p.Login, p.Position, action, p.Symbol, p.PriceOpen, p.Profit)
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
}
