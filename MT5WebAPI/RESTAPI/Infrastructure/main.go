package main

import (
	"config"
	"encoding/json"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("INFRASTRUCTURE", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Infrastructure Example")
	fmt.Println("  Tests: Email, Messenger, Gateway, Feeder,")
	fmt.Println("         Report, Plugin")
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
	// 1. Email
	// =============================================
	fmt.Println("--- 1. Email ---")
	total, ret := api.EmailTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.EmailNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  EmailNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  EmailTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 2. Messenger
	// =============================================
	fmt.Println("--- 2. Messenger ---")
	total, ret = api.MessengerTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.MessengerNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  MessengerNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  MessengerTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 3. Gateway
	// =============================================
	fmt.Println("--- 3. Gateway ---")
	total, ret = api.GatewayTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.GatewayNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  GatewayNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  GatewayTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	moduleTotal, ret := api.GatewayModuleTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Module Total: %d\n", moduleTotal)
	} else {
		fmt.Printf("  GatewayModuleTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 4. Feeder
	// =============================================
	fmt.Println("--- 4. Feeder ---")
	total, ret = api.FeederTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.FeederNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  FeederNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  FeederTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	moduleTotal, ret = api.FeederModuleTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Module Total: %d\n", moduleTotal)
	} else {
		fmt.Printf("  FeederModuleTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 5. Report
	// =============================================
	fmt.Println("--- 5. Report ---")
	total, ret = api.ReportTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.ReportNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  ReportNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  ReportTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	moduleTotal, ret = api.ReportModuleTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Module Total: %d\n", moduleTotal)
	} else {
		fmt.Printf("  ReportModuleTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 6. Plugin
	// =============================================
	fmt.Println("--- 6. Plugin ---")
	total, ret = api.PluginTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total: %d\n", total)
		if total > 0 {
			item, ret := api.PluginNext(0)
			if ret == mt5webapi.MT_RET_OK {
				printMap("  First", item)
			} else {
				fmt.Printf("  PluginNext(0) FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("  PluginTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	moduleTotal, ret = api.PluginModuleTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Module Total: %d\n", moduleTotal)
	} else {
		fmt.Printf("  PluginModuleTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All Infrastructure tests completed!")
	fmt.Println("================================================")
}

func printMap(label string, m map[string]interface{}) {
	data, _ := json.MarshalIndent(m, "  ", "  ")
	fmt.Printf("%s: %s\n", label, string(data))
}
