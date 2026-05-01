package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("SERVER", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Server Example (Go)")
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
	// 4. TIME_SERVER - Sunucu zamanını al
	// =============================================
	fmt.Println("--- 4. TimeServer ---")
	serverTime := api.TimeServer()
	if serverTime > 0 {
		t := time.Unix(serverTime, 0).UTC()
		fmt.Printf("Unix Timestamp: %d\n", serverTime)
		fmt.Printf("UTC:            %s\n", t.Format("2006-01-02 15:04:05"))
		fmt.Printf("Local:          %s\n", t.Local().Format("2006-01-02 15:04:05 MST"))
	} else {
		fmt.Println("TimeServer failed!")
	}
	fmt.Println()

	// =============================================
	// 5. TIME_GET - Zaman ayarlarını al
	// =============================================
	fmt.Println("--- 5. TimeGet ---")
	conTime, ret := api.TimeGet()
	if ret == mt5webapi.MT_RET_OK && conTime != nil {
		fmt.Printf("TimeZone:       %d minutes (GMT%+.1f)\n", conTime.TimeZone, float64(conTime.TimeZone)/3600.0)
		fmt.Printf("Daylight:       %d\n", conTime.Daylight)
		fmt.Printf("DaylightState:  %d\n", conTime.DaylightState)
		fmt.Printf("TimeServer:     %s\n", conTime.TimeServer)

		// days schedule
		if conTime.Days != nil {
			dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
			fmt.Println("Working Days:")
			for i, day := range conTime.Days {
				if i < len(dayNames) && day != nil {
					status := "disabled"
					for _, mode := range day {
						if mode == mt5webapi.TimeModeEnabled {
							status = "enabled"
							break
						}
					}
					fmt.Printf("  %-10s: %s\n", dayNames[i], status)
				}
			}
		}
	} else {
		fmt.Printf("TimeGet failed: %s\n", mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 6. COMMON_GET - Sunucu bilgilerini al
	// =============================================
	fmt.Println("--- 6. CommonGet ---")
	common, ret := api.CommonGet()
	if ret == mt5webapi.MT_RET_OK && common != nil {
		fmt.Printf("Server Name:      %s\n", common.Name)
		fmt.Printf("Owner:            %s\n", common.Owner)
		fmt.Printf("OwnerID:          %s\n", common.OwnerID)
		fmt.Printf("OwnerHost:        %s\n", common.OwnerHost)
		fmt.Printf("Product:          %s\n", common.Product)
		fmt.Println()

		// license info
		if common.ExpirationLicense > 0 {
			expLic := time.Unix(common.ExpirationLicense, 0).UTC()
			fmt.Printf("License Expires:  %s\n", expLic.Format("2006-01-02"))
		}
		if common.ExpirationSupport > 0 {
			expSup := time.Unix(common.ExpirationSupport, 0).UTC()
			fmt.Printf("Support Expires:  %s\n", expSup.Format("2006-01-02"))
		}
		fmt.Println()

		// limits
		fmt.Printf("Limit Trade Servers: %d\n", common.LimitTradeServers)
		fmt.Printf("Limit Web Servers:   %d\n", common.LimitWebServers)
		fmt.Printf("Limit Accounts:      %d\n", common.LimitAccounts)
		fmt.Printf("Limit Deals:         %d\n", common.LimitDeals)
		fmt.Printf("Limit Symbols:       %d\n", common.LimitSymbols)
		fmt.Printf("Limit Groups:        %d\n", common.LimitGroups)
		fmt.Println()

		// totals
		fmt.Printf("Total Users:         %d\n", common.TotalUsers)
		fmt.Printf("Total Users Real:    %d\n", common.TotalUsersReal)
		fmt.Printf("Total Deals:         %d\n", common.TotalDeals)
		fmt.Printf("Total Orders:        %d\n", common.TotalOrders)
		fmt.Printf("Total History:       %d\n", common.TotalOrdersHistory)
		fmt.Printf("Total Positions:     %d\n", common.TotalPositions)
		fmt.Println()

		fmt.Printf("LiveUpdate Mode:     %d\n", common.LiveUpdateMode)
		fmt.Printf("Account Auto:        %d\n", common.AccountAuto)
		fmt.Printf("Account URL:         %s\n", common.AccountURL)
	} else {
		fmt.Printf("CommonGet failed: %s\n", mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 7. SERVER_RESTART - Sunucuyu yeniden başlat
	// =============================================
	fmt.Println("--- 7. ServerRestart ---")
	fmt.Println("WARNING: ServerRestart will restart the MT5 server!")
	fmt.Println("Skipping to avoid disruption. Usage:")
	fmt.Println("  ret := api.ServerRestart()")
	fmt.Println("  // returns MT_RET_OK on success")
	fmt.Println()

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All server tests completed!")
	fmt.Println("========================================")
}
