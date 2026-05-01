package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("DEAL_BATCH", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Deal Batch Operations Example")
	fmt.Println("  Tests: 23-25 (3 deal batch functions)")
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
	// 23. DEAL_GET_BATCH - Get deals for login 5000 (last 30 days)
	// =============================================
	fmt.Println("--- 23. DealGetBatch (login 5000) ---")
	now := time.Now().Unix() + 86400
	from := time.Now().AddDate(0, -1, 0).Unix()
	deals, ret := api.DealGetBatch("5000", "", "", from, now)
	if ret == mt5webapi.MT_RET_OK && deals != nil {
		fmt.Printf("  Found %d deals for login 5000:\n", len(deals))
		limit := 5
		if len(deals) < limit {
			limit = len(deals)
		}
		for i := 0; i < limit; i++ {
			d := deals[i]
			fmt.Printf("    deal=%d action=%d symbol=%s price=%.5f profit=%.2f comment=%s\n",
				d.Deal, d.Action, d.Symbol, d.Price, d.Profit, d.Comment)
		}
		if len(deals) > 5 {
			fmt.Printf("    ... +%d more\n", len(deals)-5)
		}
	} else {
		fmt.Printf("  FAILED or no data: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// --- Setup test user for deal tests ---
	testLogin := uint64(99998)
	fmt.Println("--- Setup: creating test user 99998 ---")
	newUser := mt5webapi.CreateDefault()
	newUser.Login = testLogin
	newUser.Name = "Deal Batch Test User"
	newUser.Group = "demo\\forex"
	newUser.Leverage = 100
	newUser.MainPassword = "Test1234!"
	newUser.InvestPassword = "Invest1234!"
	api.UserDelete(testLogin) // cleanup any leftover
	time.Sleep(500 * time.Millisecond)
	_, ret = api.UserAdd(newUser)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("  UserAdd failed: %s - using login 5000 for remaining tests\n", mt5webapi.GetError(ret))
		testLogin = 5000
	} else {
		fmt.Printf("  User %d created\n", testLogin)
		// deposit to generate balance deals
		ret = api.UserDepositChange(testLogin, 50000, "initial deposit", mt5webapi.DealBalance)
		if ret == mt5webapi.MT_RET_OK {
			fmt.Println("  Deposited 50000 (should generate a balance deal)")
		} else {
			fmt.Printf("  Deposit failed: %s\n", mt5webapi.GetError(ret))
		}
		time.Sleep(1 * time.Second)
		// second deposit for more deals
		ret = api.UserDepositChange(testLogin, 10000, "bonus deposit", mt5webapi.DealBalance)
		if ret == mt5webapi.MT_RET_OK {
			fmt.Println("  Deposited 10000 (second balance deal)")
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println()

	// =============================================
	// 23b. DEAL_GET_BATCH - Get deals for test user (should see deposit deals)
	// =============================================
	fmt.Println("--- 23b. DealGetBatch (test user) ---")
	loginStr := fmt.Sprintf("%d", testLogin)
	testDeals, ret := api.DealGetBatch(loginStr, "", "", from, now)
	if ret == mt5webapi.MT_RET_OK && testDeals != nil {
		fmt.Printf("  Found %d deals for login %d:\n", len(testDeals), testLogin)
		for _, d := range testDeals {
			fmt.Printf("    deal=%d action=%d symbol=%s profit=%.2f comment=%s\n",
				d.Deal, d.Action, d.Symbol, d.Profit, d.Comment)
		}
	} else {
		fmt.Printf("  FAILED or no data: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 24. DEAL_UPDATE - Update comment on a deal
	// =============================================
	fmt.Println("--- 24. DealUpdateDirect ---")
	if testDeals != nil && len(testDeals) > 0 {
		ticket := testDeals[0].Deal
		newComment := "Updated by DealBatch"
		updateJSON := fmt.Sprintf(`{"Deal":"%d","Comment":"%s"}`, ticket, newComment)
		updatedDeal, ret := api.DealUpdateDirect(updateJSON)
		if ret == mt5webapi.MT_RET_OK && updatedDeal != nil {
			fmt.Printf("  Updated deal %d: comment=%s\n", updatedDeal.Deal, updatedDeal.Comment)
		} else {
			fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	} else {
		fmt.Println("  SKIPPED - no deals to update")
	}
	fmt.Println()

	// =============================================
	// 25. DEAL_DELETE - Delete deals (SKIPPED - dangerous)
	// =============================================
	fmt.Println("--- 25. DealDeleteDirect ---")
	fmt.Println("  SKIPPED - dangerous operation, would permanently delete deal records")
	fmt.Println("  Usage: ret := api.DealDeleteDirect(\"ticket1,ticket2\")")
	fmt.Println()

	// =============================================
	// Cleanup
	// =============================================
	fmt.Println("--- Cleanup ---")
	if testLogin == 99998 {
		api.UserDelete(99998)
		fmt.Println("  User 99998 cleaned up")
	}

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All deal batch tests completed!")
	fmt.Println("================================================")
}
