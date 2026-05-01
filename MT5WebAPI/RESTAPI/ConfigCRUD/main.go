package main

import (
	"config"
	"encoding/json"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("CONFIG_CRUD", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Config CRUD Example")
	fmt.Println("  Tests: Leverage, Firewall, Holiday, Manager,")
	fmt.Println("         Route, Spread, HistorySync, SymbolGroup")
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
	// 1. Leverage
	// =============================================
	fmt.Println("--- 1. Leverage ---")
	total, ret := api.LeverageTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.LeverageNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  LeverageNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  LeverageTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 2. Firewall
	// =============================================
	fmt.Println("--- 2. Firewall ---")
	total, ret = api.FirewallTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.FirewallNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  FirewallNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  FirewallTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 3. Holiday
	// =============================================
	fmt.Println("--- 3. Holiday ---")
	total, ret = api.HolidayTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.HolidayNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  HolidayNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  HolidayTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 4. Manager
	// =============================================
	fmt.Println("--- 4. Manager ---")
	total, ret = api.ManagerTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.ManagerNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  ManagerNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  ManagerTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 5. Route
	// =============================================
	fmt.Println("--- 5. Route ---")
	total, ret = api.RouteTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.RouteNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  RouteNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  RouteTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 6. Spread
	// =============================================
	fmt.Println("--- 6. Spread ---")
	total, ret = api.SpreadTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.SpreadNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  SpreadNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  SpreadTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 7. HistorySync
	// =============================================
	fmt.Println("--- 7. HistorySync ---")
	total, ret = api.HistorySyncTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.HistorySyncNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  HistorySyncNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  HistorySyncTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 8. SymbolGroup
	// =============================================
	fmt.Println("--- 8. SymbolGroup ---")
	total, ret = api.SymbolGroupTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.SymbolGroupNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  SymbolGroupNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  SymbolGroupTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All Config CRUD tests completed!")
	fmt.Println("================================================")
}

func printMap(label string, m map[string]interface{}) {
	data, _ := json.MarshalIndent(m, "  ", "  ")
	fmt.Printf("%s: %s\n", label, string(data))
}
