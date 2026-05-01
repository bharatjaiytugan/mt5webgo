package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("CUSTOM", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Custom Example (Go)")
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

	// =============================================
	// 46. CUSTOM_SEND - Özel komut gönder
	// =============================================
	fmt.Println("--- 46. CustomSend ---")
	fmt.Println()

	// 1) Use CustomSend to call TIME_SERVER (known command as test)
	fmt.Println("1) CustomSend: TIME_SERVER")
	answer, ret := api.CustomSend("TIME_SERVER", nil, "")
	if ret == mt5webapi.MT_RET_OK && answer != nil {
		answerStr := mt5webapi.GetString(answer)
		fmt.Printf("   Response: %s\n", answerStr)
	} else {
		fmt.Printf("   Failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 2) Use CustomSend to call COMMON_GET
	fmt.Println("2) CustomSend: COMMON_GET")
	answer, ret = api.CustomSend("COMMON_GET", nil, "")
	if ret == mt5webapi.MT_RET_OK && answer != nil {
		answerStr := mt5webapi.GetString(answer)
		// show first 200 chars
		if len(answerStr) > 200 {
			fmt.Printf("   Response (first 200 chars): %s...\n", answerStr[:200])
		} else {
			fmt.Printf("   Response: %s\n", answerStr)
		}
	} else {
		fmt.Printf("   Failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 3) Use CustomSend to call USER_GET with params
	fmt.Println("3) CustomSend: USER_GET (login=5005)")
	params := map[string]string{
		"LOGIN": "5005",
	}
	answer, ret = api.CustomSend("USER_GET", params, "")
	if ret == mt5webapi.MT_RET_OK && answer != nil {
		answerStr := mt5webapi.GetString(answer)
		if len(answerStr) > 300 {
			fmt.Printf("   Response (first 300 chars): %s...\n", answerStr[:300])
		} else {
			fmt.Printf("   Response: %s\n", answerStr)
		}
	} else {
		fmt.Printf("   Failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 4) Use CustomSend with GROUP_TOTAL
	fmt.Println("4) CustomSend: GROUP_TOTAL")
	answer, ret = api.CustomSend("GROUP_TOTAL", nil, "")
	if ret == mt5webapi.MT_RET_OK && answer != nil {
		answerStr := mt5webapi.GetString(answer)
		fmt.Printf("   Response: %s\n", answerStr)
	} else {
		fmt.Printf("   Failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	fmt.Println("CustomSend Usage:")
	fmt.Println("  answer, ret := api.CustomSend(command, params, body)")
	fmt.Println("  // command: any MT5 protocol command string")
	fmt.Println("  // params:  map[string]string of key=value pairs")
	fmt.Println("  // body:    optional body text (JSON, etc.)")
	fmt.Println("  // answer:  raw []byte response from server")
	fmt.Println()

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All custom tests completed!")
	fmt.Println("========================================")
}
