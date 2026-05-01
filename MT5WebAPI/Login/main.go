package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	// create API with logging
	api := mt5webapi.NewMT5WebAPI("LOGIN", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Login Example (Go)")
	fmt.Println("========================================")
	fmt.Println()

	// load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		return
	}

	fmt.Printf("Server:   %s:%d\n", cfg.Server, cfg.Port)
	fmt.Printf("Login:    %d\n", cfg.Login)
	fmt.Println()

	// =============================================
	// 1. CONNECT - Sunucuya bağlan
	// =============================================
	fmt.Println("--- 1. Connect ---")
	fmt.Print("Connecting... ")
	ret := api.ConnectDefault(cfg.Server, cfg.Port, cfg.Login, cfg.Password, mt5webapi.PUMP_MODE_NONE)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("FAILED! [%d] %s\n", ret, mt5webapi.GetError(ret))
		return
	}
	fmt.Println("OK")
	fmt.Printf("IsConnected: %v\n", api.IsConnected())
	fmt.Println()

	// server time
	serverTime := api.TimeServer()
	t := time.Unix(serverTime, 0).UTC()
	fmt.Printf("Server Time: %s (unix: %d)\n", t.Format("2006-01-02 15:04:05 UTC"), serverTime)
	fmt.Println()

	// server info
	common, ret := api.CommonGet()
	if ret == mt5webapi.MT_RET_OK && common != nil {
		fmt.Printf("Server Name:  %s\n", common.Name)
		fmt.Printf("Owner:        %s\n", common.Owner)
		fmt.Printf("Product:      %s\n", common.Product)
		fmt.Printf("Total Users:  %d\n", common.TotalUsers)
		fmt.Printf("Total Deals:  %d\n", common.TotalDeals)
		fmt.Printf("Positions:    %d\n", common.TotalPositions)
	}
	fmt.Println()

	// user info
	user, ret := api.UserGet(cfg.Login)
	if ret == mt5webapi.MT_RET_OK && user != nil {
		fmt.Printf("Login:    %d\n", user.Login)
		fmt.Printf("Name:     %s\n", user.Name)
		fmt.Printf("Group:    %s\n", user.Group)
		fmt.Printf("Country:  %s\n", user.Country)
		fmt.Printf("Leverage: %d\n", user.Leverage)
		fmt.Printf("Balance:  %.2f\n", user.Balance)
	}
	fmt.Println()

	// =============================================
	// 3. PING - Keepalive gönder
	// =============================================
	fmt.Println("--- 3. Ping ---")
	for i := 1; i <= 3; i++ {
		ret = api.Ping()
		fmt.Printf("Ping #%d: %s\n", i, mt5webapi.GetError(ret))
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Printf("IsConnected after pings: %v\n", api.IsConnected())
	fmt.Println()

	// =============================================
	// 2. DISCONNECT - Bağlantıyı kes
	// =============================================
	fmt.Println("--- 2. Disconnect ---")
	api.Disconnect()
	fmt.Printf("IsConnected: %v\n", api.IsConnected())
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("  All connection tests completed!")
	fmt.Println("========================================")
}
