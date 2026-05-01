package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("SERVER_CONFIG", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Server Config Example")
	fmt.Println("========================================")
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
	// 1. CommonGet / CommonSet
	// =============================================
	fmt.Println("--- 1. CommonGet (before) ---")
	common, ret := api.CommonGet()
	if ret == mt5webapi.MT_RET_OK && common != nil {
		fmt.Printf("Server Name:   %s\n", common.Name)
		fmt.Printf("Owner:         %s\n", common.Owner)
		fmt.Printf("AccountAuto:   %d\n", common.AccountAuto)
		fmt.Printf("AccountURL:    %s\n", common.AccountURL)
		fmt.Println()

		// toggle AccountAuto: 0 -> 1, nonzero -> 0
		originalAuto := common.AccountAuto
		if common.AccountAuto == 0 {
			common.AccountAuto = 1
		} else {
			common.AccountAuto = 0
		}

		fmt.Printf("--- 1b. CommonSet (AccountAuto: %d -> %d) ---\n", originalAuto, common.AccountAuto)
		updated, ret := api.CommonSet(common)
		if ret == mt5webapi.MT_RET_OK && updated != nil {
			fmt.Printf("CommonSet OK:  AccountAuto=%d\n", updated.AccountAuto)
		} else {
			fmt.Printf("CommonSet FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
		fmt.Println()

		// verify with CommonGet
		fmt.Println("--- 1c. CommonGet (verify) ---")
		verify, ret := api.CommonGet()
		if ret == mt5webapi.MT_RET_OK && verify != nil {
			fmt.Printf("AccountAuto:   %d\n", verify.AccountAuto)
		} else {
			fmt.Printf("CommonGet FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
		fmt.Println()

		// restore original value
		common.AccountAuto = originalAuto
		fmt.Printf("--- 1d. CommonSet (restore AccountAuto to %d) ---\n", originalAuto)
		restored, ret := api.CommonSet(common)
		if ret == mt5webapi.MT_RET_OK && restored != nil {
			fmt.Printf("Restored OK:   AccountAuto=%d\n", restored.AccountAuto)
		} else {
			fmt.Printf("Restore FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	} else {
		fmt.Printf("CommonGet FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 2. AccessServerTotal / AccessServerNext
	// =============================================
	fmt.Println("--- 2. AccessServerTotal ---")
	total, ret := api.AccessServerTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("Total access servers: %d\n", total)

		// list all servers
		limit := total
		for i := 0; i < limit; i++ {
			fmt.Printf("\n--- 2b. AccessServerNext(%d) ---\n", i)
			srv, ret := api.AccessServerNext(i)
			if ret == mt5webapi.MT_RET_OK && srv != nil {
				for k, v := range srv {
					fmt.Printf("  %s: %v\n", k, v)
				}
			} else {
				fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Printf("AccessServerTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 3. TLSCertTotal
	// =============================================
	fmt.Println("--- 3. TLSCertTotal ---")
	tlsTotal, ret := api.TLSCertTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("Total TLS certificates: %d\n", tlsTotal)
	} else {
		fmt.Printf("TLSCertTotal FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 4. TimeGet (show current) / TimeSet (skip)
	// =============================================
	fmt.Println("--- 4. TimeGet ---")
	conTime, ret := api.TimeGet()
	if ret == mt5webapi.MT_RET_OK && conTime != nil {
		fmt.Printf("TimeZone:       %d seconds (GMT%+.1f)\n", conTime.TimeZone, float64(conTime.TimeZone)/3600.0)
		fmt.Printf("Daylight:       %d\n", conTime.Daylight)
		fmt.Printf("DaylightState:  %d\n", conTime.DaylightState)
		fmt.Printf("TimeServer:     %s\n", conTime.TimeServer)

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
		fmt.Printf("TimeGet FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	fmt.Println("--- 4b. TimeSet ---")
	fmt.Println("SKIPPED: TimeSet not called to avoid disrupting server time config.")
	fmt.Println("Usage:")
	fmt.Println("  updated, ret := api.TimeSet(conTime)")
	fmt.Println()

	// server time
	fmt.Println("--- 5. TimeServer ---")
	serverTime := api.TimeServer()
	if serverTime > 0 {
		t := time.Unix(serverTime, 0).UTC()
		fmt.Printf("Server Time: %s (unix: %d)\n", t.Format("2006-01-02 15:04:05 UTC"), serverTime)
	} else {
		fmt.Println("TimeServer failed!")
	}
	fmt.Println()

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All server config tests completed!")
	fmt.Println("========================================")
}
