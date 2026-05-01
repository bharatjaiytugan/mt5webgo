package main

import (
	"config"
	"encoding/json"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("SUBSCRIPTION", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Subscription Example")
	fmt.Println("  Tests: 41-42 (Subscription Config & Management)")
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
	// 41. Subscription Config
	// =============================================

	// Step 1: Get total subscription configs
	fmt.Println("--- 41a. SubscriptionCfgTotal ---")
	total, ret := api.SubscriptionCfgTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total subscription configs: %d\n", total)
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// Step 2: List first config if any exist
	var firstConfigID int
	hasConfigs := false
	if total > 0 {
		fmt.Println("--- 41b. SubscriptionCfgNext(0) ---")
		cfgMap, ret := api.SubscriptionCfgNext(0)
		if ret == mt5webapi.MT_RET_OK && cfgMap != nil {
			cfgJSON, _ := json.MarshalIndent(cfgMap, "  ", "  ")
			fmt.Printf("  Config[0]: %s\n", string(cfgJSON))
			// try to extract ID for later use
			if idVal, ok := cfgMap["ID"]; ok {
				switch v := idVal.(type) {
				case float64:
					firstConfigID = int(v)
					hasConfigs = true
				}
			}
		} else {
			fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
		fmt.Println()

		// Step 3: Get config by ID if we have one
		if hasConfigs {
			fmt.Printf("--- 41c. SubscriptionCfgGet(ID=%d) ---\n", firstConfigID)
			cfgMap, ret := api.SubscriptionCfgGet(firstConfigID)
			if ret == mt5webapi.MT_RET_OK && cfgMap != nil {
				cfgJSON, _ := json.MarshalIndent(cfgMap, "  ", "  ")
				fmt.Printf("  Config: %s\n", string(cfgJSON))
			} else {
				fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
			fmt.Println()
		}
	}

	// =============================================
	// 42. Subscription Management
	// =============================================

	testLogin := uint64(5000)

	// Step 4: Get subscriptions for login
	fmt.Printf("--- 42a. SubscriptionGet(login=%d) ---\n", testLogin)
	subMap, ret := api.SubscriptionGet(testLogin)
	if ret == mt5webapi.MT_RET_OK && subMap != nil {
		subJSON, _ := json.MarshalIndent(subMap, "  ", "  ")
		fmt.Printf("  Subscriptions: %s\n", string(subJSON))
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// Step 5: If configs exist, test exist/join/cancel
	if hasConfigs {
		fmt.Printf("--- 42b. SubscriptionExist(login=%d, ID=%d) ---\n", testLogin, firstConfigID)
		exists, ret := api.SubscriptionExist(testLogin, firstConfigID)
		if ret == mt5webapi.MT_RET_OK {
			fmt.Printf("  Exists: %v\n", exists)
		} else {
			fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
		fmt.Println()

		fmt.Printf("--- 42c. SubscriptionJoin(login=%d, ID=%d) ---\n", testLogin, firstConfigID)
		ret = api.SubscriptionJoin(testLogin, firstConfigID)
		if ret == mt5webapi.MT_RET_OK {
			fmt.Println("  Join: OK")
		} else {
			fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
		fmt.Println()

		fmt.Printf("--- 42d. SubscriptionCancel(login=%d, ID=%d) ---\n", testLogin, firstConfigID)
		ret = api.SubscriptionCancel(testLogin, firstConfigID)
		if ret == mt5webapi.MT_RET_OK {
			fmt.Println("  Cancel: OK")
		} else {
			fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
		fmt.Println()
	}

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All subscription tests completed!")
	fmt.Println("================================================")
}
