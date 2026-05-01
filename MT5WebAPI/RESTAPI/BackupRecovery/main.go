package main

import (
	"config"
	"fmt"
	"math"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("BACKUP", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Backup/Recovery Example (Go)")
	fmt.Println("  Tests: 19-22")
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
	// 19. ORDER_BACKUP_LIST
	// =============================================
	fmt.Println("--- 19. OrderBackupList ---")
	now := time.Now().Unix() + 86400
	from30 := time.Now().AddDate(0, -1, 0).Unix()
	backupDates, ret := api.OrderBackupList(from30, now)
	if ret == mt5webapi.MT_RET_OK && backupDates != nil {
		fmt.Printf("  Found %d backup dates:\n", len(backupDates))
		for i, d := range backupDates {
			t := time.Unix(d, 0)
			fmt.Printf("    [%d] %d (%s)\n", i, d, t.Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 20. ORDER_BACKUP_GET
	// =============================================
	fmt.Println("--- 20. OrderBackupGet ---")
	if backupDates != nil && len(backupDates) > 0 {
		// find a non-zero backup date
		var backupDate int64
		for i := len(backupDates) - 1; i >= 0; i-- {
			if backupDates[i] > 0 {
				backupDate = backupDates[i]
				break
			}
		}
		if backupDate > 0 {
			t := time.Unix(backupDate, 0)
			fmt.Printf("  Using backup: %d (%s)\n", backupDate, t.Format("2006-01-02 15:04:05"))

			// try with login 5000
			fmt.Println("  Querying for login 5000...")
			api.EnableWriteDebugLog(true)
			orders, ret := api.OrderBackupGet(backupDate, 5000)
			api.EnableWriteDebugLog(false)
			if ret == mt5webapi.MT_RET_OK && orders != nil {
				fmt.Printf("  Found %d orders in backup:\n", len(orders))
				for i, o := range orders {
					if i >= 5 {
						fmt.Printf("    ... +%d more\n", len(orders)-5)
						break
					}
					fmt.Printf("    ticket=%d type=%d symbol=%s login=%d\n",
						o.Order, o.Type, o.Symbol, o.Login)
				}
			} else {
				fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		} else {
			fmt.Println("  No valid backup dates found")
		}
	} else {
		fmt.Println("  SKIPPED - no backup dates")
	}
	fmt.Println()

	// =============================================
	// Setup: Create user, deposit, place pending, cancel it
	// =============================================
	fmt.Println("--- Setup for OrderReopen test ---")
	testLogin := uint64(10404)

	// cleanup any leftover first - cancel pending orders then delete
	existingOrders, _ := api.OrderGetBatch(fmt.Sprintf("%d", testLogin), "", "")
	if existingOrders != nil {
		for _, o := range existingOrders {
			api.DealerCancelOrder(testLogin, o.Order)
			time.Sleep(500 * time.Millisecond)
		}
	}
	api.UserDelete(testLogin)
	time.Sleep(1 * time.Second)

	newUser := mt5webapi.CreateDefault()
	newUser.Login = testLogin
	newUser.Name = "Backup Test"
	newUser.Group = "demo\\forex"
	newUser.Leverage = 100
	newUser.MainPassword = "Test1234!"
	newUser.InvestPassword = "Invest1234!"
	_, ret = api.UserAdd(newUser)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("  UserAdd: %s (using existing)\n", mt5webapi.GetError(ret))
	}
	api.UserDepositChange(testLogin, 100000, "deposit", mt5webapi.DealBalance)
	fmt.Printf("  User %d created with deposit\n", testLogin)

	// get EURUSD price
	ticks, _ := api.TickLast("EURUSD")
	sym, _ := api.SymbolGet("EURUSD")
	if ticks == nil || len(ticks) == 0 || sym == nil {
		fmt.Println("  Cannot get EURUSD price")
		api.UserDelete(testLogin)
		api.Disconnect()
		return
	}
	bid := ticks[0].Bid
	point := sym.Point
	digits := sym.Digits

	// place a Buy Limit
	blPrice := roundPrice(bid-1000*point, digits)
	fmt.Printf("  Placing Buy Limit @ %.5f...\n", blPrice)
	r, ret := api.DealerPendingOrder(testLogin, "EURUSD", mt5webapi.OpBuyLimit, 100, blPrice, 0, 0, "Reopen Test")
	if ret != mt5webapi.MT_RET_OK || r == nil {
		fmt.Printf("  DealerPendingOrder failed: %s\n", mt5webapi.GetError(ret))
		api.UserDelete(testLogin)
		api.Disconnect()
		return
	}
	fmt.Printf("  Dealer ID=%d\n", r.ID)
	time.Sleep(3 * time.Second)

	// get the order ticket
	ordTotal, _ := api.OrderGetTotal(testLogin)
	fmt.Printf("  Pending orders: %d\n", ordTotal)
	var cancelTicket uint64
	if ordTotal > 0 {
		orders, ret := api.OrderGetPage(testLogin, 0, uint32(ordTotal))
		if ret == mt5webapi.MT_RET_OK && orders != nil {
			for _, o := range orders {
				fmt.Printf("  Order: ticket=%d type=%d price=%.5f\n", o.Order, o.Type, o.PriceOrder)
				if cancelTicket == 0 {
					cancelTicket = o.Order
				}
			}
		}
	}

	if cancelTicket == 0 {
		fmt.Println("  No pending order found for cancel/reopen test")
		api.UserDelete(testLogin)
		api.Disconnect()
		return
	}

	// cancel it (move to history) using OrderCancel
	fmt.Printf("  Canceling order %d via OrderCancel...\n", cancelTicket)
	ret = api.OrderCancel(fmt.Sprintf("%d", cancelTicket))
	fmt.Printf("  OrderCancel result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	time.Sleep(3 * time.Second)

	// verify it's in history
	ordTotal2, _ := api.OrderGetTotal(testLogin)
	fmt.Printf("  Pending orders after cancel: %d\n", ordTotal2)
	fmt.Println()

	// =============================================
	// 22. ORDER_REOPEN
	// =============================================
	fmt.Println("--- 22. OrderReopen ---")
	fmt.Printf("  Reopening order ticket=%d...\n", cancelTicket)
	reopened, ret := api.OrderReopen(cancelTicket)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Println("  OrderReopen OK!")
		if reopened != nil {
			fmt.Printf("  Order: ticket=%d type=%d symbol=%s price=%.5f\n",
				reopened.Order, reopened.Type, reopened.Symbol, reopened.PriceOrder)
		}
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	// verify it's back as pending
	time.Sleep(1 * time.Second)
	ordTotal3, _ := api.OrderGetTotal(testLogin)
	fmt.Printf("  Pending orders after reopen: %d\n", ordTotal3)
	fmt.Println()

	// =============================================
	// 21. ORDER_BACKUP_RESTORE
	// =============================================
	fmt.Println("--- 21. OrderBackupRestore ---")
	fmt.Println("  SKIPPED - requires specific backup order JSON")
	fmt.Println("  Usage:")
	fmt.Println("    orders, _ := api.OrderBackupGet(backupDate, login)")
	fmt.Println("    restored, ret := api.OrderBackupRestore(orderJSON)")
	fmt.Println()

	// cleanup
	fmt.Println("--- Cleanup ---")
	// cancel remaining orders
	remaining, _ := api.OrderGetBatch(fmt.Sprintf("%d", testLogin), "", "")
	if remaining != nil && len(remaining) > 0 {
		for _, o := range remaining {
			api.DealerCancelOrder(testLogin, o.Order)
			time.Sleep(500 * time.Millisecond)
		}
	}
	time.Sleep(1 * time.Second)
	ret = api.UserDelete(testLogin)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  User %d deleted\n", testLogin)
	} else {
		fmt.Printf("  User delete: %s\n", mt5webapi.GetError(ret))
	}

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
}

func roundPrice(price float64, digits uint32) float64 {
	mult := math.Pow(10, float64(digits))
	return math.Round(price*mult) / mult
}
