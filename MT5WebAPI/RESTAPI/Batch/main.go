package main

import (
	"config"
	"fmt"
	"math"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("BATCH", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Batch Operations Example (Go)")
	fmt.Println("  Tests: 7-18 (12 batch functions)")
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
	// 9. GROUP_LIST - Grup listesi (mask ile)
	// =============================================
	fmt.Println("--- 9. GroupList ---")
	groups, ret := api.GroupList("*")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  GroupList(\"*\"): %d groups found\n", len(groups))
		// show first 5
		for i := 0; i < len(groups) && i < 5; i++ {
			fmt.Printf("    [%d] %s\n", i, groups[i])
		}
		if len(groups) > 5 {
			fmt.Printf("    ... +%d more\n", len(groups)-5)
		}
	} else {
		fmt.Printf("  FAILED: %s\n", mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 7. GROUP_ADD_BATCH - Toplu grup ekleme
	// =============================================
	fmt.Println("--- 7. GroupAddBatch ---")
	// clone existing group for template
	templateGroup, ret := api.GroupGet("demo\\forex")
	if ret == mt5webapi.MT_RET_OK && templateGroup != nil {
		templateGroup.Group = "demo\\batchtest1"
		json1 := groupMinJSON("demo\\batchtest1", templateGroup)
		templateGroup.Group = "demo\\batchtest2"
		json2 := groupMinJSON("demo\\batchtest2", templateGroup)
		batchJSON := "[" + json1 + "," + json2 + "]"

		ret = api.GroupAddBatch(batchJSON)
		fmt.Printf("  GroupAddBatch: [%d] %s\n", ret, mt5webapi.GetError(ret))

		// verify
		g1, r1 := api.GroupGet("demo\\batchtest1")
		g2, r2 := api.GroupGet("demo\\batchtest2")
		if r1 == mt5webapi.MT_RET_OK && g1 != nil {
			fmt.Printf("  Verified: %s exists\n", g1.Group)
		}
		if r2 == mt5webapi.MT_RET_OK && g2 != nil {
			fmt.Printf("  Verified: %s exists\n", g2.Group)
		}
	} else {
		fmt.Printf("  Cannot get template group: %s\n", mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 8. GROUP_DELETE_BATCH - Toplu grup silme
	// =============================================
	fmt.Println("--- 8. GroupDeleteBatch ---")
	ret = api.GroupDeleteBatch("demo\\batchtest1,demo\\batchtest2")
	fmt.Printf("  GroupDeleteBatch: [%d] %s\n", ret, mt5webapi.GetError(ret))
	// verify deleted
	_, r1 := api.GroupGet("demo\\batchtest1")
	_, r2 := api.GroupGet("demo\\batchtest2")
	fmt.Printf("  batchtest1 after delete: [%d] %s\n", r1, mt5webapi.GetError(r1))
	fmt.Printf("  batchtest2 after delete: [%d] %s\n", r2, mt5webapi.GetError(r2))
	fmt.Println()

	// =============================================
	// 10. SYMBOL_ADD_BATCH - Toplu sembol ekleme
	// =============================================
	fmt.Println("--- 10. SymbolAddBatch ---")
	// clone EURUSD as template
	templateSym, ret := api.SymbolGet("EURUSD")
	if ret == mt5webapi.MT_RET_OK && templateSym != nil {
		templateSym.Symbol = "GOTEST1"
		templateSym.Path = "Forex\\GOTEST1"
		templateSym.Description = "Batch Test 1"
		j1 := mt5webapi.SymbolToJSON(templateSym)
		templateSym.Symbol = "GOTEST2"
		templateSym.Path = "Forex\\GOTEST2"
		templateSym.Description = "Batch Test 2"
		j2 := mt5webapi.SymbolToJSON(templateSym)
		batchJSON := "[" + j1 + "," + j2 + "]"

		ret = api.SymbolAddBatch(batchJSON)
		fmt.Printf("  SymbolAddBatch: [%d] %s\n", ret, mt5webapi.GetError(ret))

		// verify
		s1, r1 := api.SymbolGet("GOTEST1")
		s2, r2 := api.SymbolGet("GOTEST2")
		if r1 == mt5webapi.MT_RET_OK && s1 != nil {
			fmt.Printf("  Verified: %s exists\n", s1.Symbol)
		}
		if r2 == mt5webapi.MT_RET_OK && s2 != nil {
			fmt.Printf("  Verified: %s exists\n", s2.Symbol)
		}
	} else {
		fmt.Printf("  Cannot get template symbol: %s\n", mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 11. SYMBOL_DELETE_BATCH - Toplu sembol silme
	// =============================================
	fmt.Println("--- 11. SymbolDeleteBatch ---")
	ret = api.SymbolDeleteBatch("GOTEST1,GOTEST2")
	fmt.Printf("  SymbolDeleteBatch: [%d] %s\n", ret, mt5webapi.GetError(ret))
	_, r1 = api.SymbolGet("GOTEST1")
	_, r2 = api.SymbolGet("GOTEST2")
	fmt.Printf("  GOTEST1 after delete: [%d] %s\n", r1, mt5webapi.GetError(r1))
	fmt.Printf("  GOTEST2 after delete: [%d] %s\n", r2, mt5webapi.GetError(r2))
	fmt.Println()

	// --- Setup test user for order/history tests ---
	testLogin := uint64(99999)
	fmt.Println("--- Setup: creating test user 99999 ---")
	newUser := mt5webapi.CreateDefault()
	newUser.Login = testLogin
	newUser.Name = "Batch Test User"
	newUser.Group = "demo\\forex"
	newUser.Leverage = 100
	newUser.MainPassword = "Test1234!"
	newUser.InvestPassword = "Invest1234!"
	api.UserDelete(testLogin) // cleanup any leftover
	time.Sleep(500 * time.Millisecond)
	_, ret = api.UserAdd(newUser)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("  UserAdd failed: %s - using login 5000 instead\n", mt5webapi.GetError(ret))
		testLogin = 5000
	} else {
		api.UserDepositChange(testLogin, 100000, "deposit", mt5webapi.DealBalance)
		fmt.Printf("  User %d created with 100000 deposit\n", testLogin)
	}

	// place pending orders for order tests
	symbol := "EURUSD"
	ticks, _ := api.TickLast(symbol)
	sym, _ := api.SymbolGet(symbol)
	if ticks != nil && len(ticks) > 0 && sym != nil {
		bid := ticks[0].Bid
		ask := ticks[0].Ask
		point := sym.Point
		digits := sym.Digits

		if testLogin == 99999 {
			blPrice := roundPrice(bid-1000*point, digits)
			slPrice := roundPrice(ask+1000*point, digits)
			api.DealerPendingOrder(testLogin, symbol, mt5webapi.OpBuyLimit, 100, blPrice, 0, 0, "Batch BL")
			time.Sleep(2 * time.Second)
			api.DealerPendingOrder(testLogin, symbol, mt5webapi.OpSellLimit, 100, slPrice, 0, 0, "Batch SL")
			time.Sleep(2 * time.Second)
			fmt.Println("  2 pending orders placed")
		}
	}
	fmt.Println()

	// =============================================
	// 12. ORDER_GET_BATCH - Toplu emir sorgulama
	// =============================================
	fmt.Println("--- 12. OrderGetBatch ---")
	loginStr := fmt.Sprintf("%d", testLogin)
	orders, ret := api.OrderGetBatch(loginStr, "", "")
	if ret == mt5webapi.MT_RET_OK && orders != nil {
		fmt.Printf("  Found %d orders for login %d:\n", len(orders), testLogin)
		for _, o := range orders {
			fmt.Printf("    ticket=%d type=%d symbol=%s price=%.5f comment=%s\n",
				o.Order, o.Type, o.Symbol, o.PriceOrder, o.Comment)
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 13. ORDER_UPDATE - Açık emri doğrudan değiştir
	// =============================================
	fmt.Println("--- 13. OrderUpdateDirect ---")
	if orders != nil && len(orders) > 0 {
		ticket := orders[0].Order
		newComment := "Updated by Batch"
		updateJSON := fmt.Sprintf(`{"Order":"%d","Comment":"%s"}`, ticket, newComment)
		updatedOrder, ret := api.OrderUpdateDirect(updateJSON)
		if ret == mt5webapi.MT_RET_OK && updatedOrder != nil {
			fmt.Printf("  Updated order %d: comment=%s\n", updatedOrder.Order, updatedOrder.Comment)
		} else {
			fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	} else {
		fmt.Println("  SKIPPED - no orders to update")
	}
	fmt.Println()

	// =============================================
	// 15. ORDER_CANCEL - Emri history'ye taşı
	// =============================================
	fmt.Println("--- 15. OrderCancel ---")
	if orders != nil && len(orders) > 0 {
		ticket := orders[0].Order
		fmt.Printf("  Canceling order %d (move to history)...\n", ticket)
		ret = api.OrderCancel(fmt.Sprintf("%d", ticket))
		fmt.Printf("  OrderCancel: [%d] %s\n", ret, mt5webapi.GetError(ret))
		time.Sleep(1 * time.Second)
	} else {
		fmt.Println("  SKIPPED - no orders")
	}
	fmt.Println()

	// =============================================
	// 14. ORDER_DELETE - Açık emri sil
	// =============================================
	fmt.Println("--- 14. OrderDeleteDirect ---")
	// refresh orders
	orders2, ret := api.OrderGetBatch(loginStr, "", "")
	if ret == mt5webapi.MT_RET_OK && orders2 != nil && len(orders2) > 0 {
		ticket := orders2[0].Order
		fmt.Printf("  Deleting order %d...\n", ticket)
		ret = api.OrderDeleteDirect(fmt.Sprintf("%d", ticket))
		fmt.Printf("  OrderDeleteDirect: [%d] %s\n", ret, mt5webapi.GetError(ret))
	} else {
		fmt.Println("  SKIPPED - no orders to delete")
	}
	fmt.Println()

	// =============================================
	// 16. HISTORY_GET_BATCH - Toplu geçmiş sorgulama
	// =============================================
	fmt.Println("--- 16. HistoryGetBatch ---")
	now := time.Now().Unix() + 86400
	from := time.Now().AddDate(0, -1, 0).Unix()
	histOrders, ret := api.HistoryGetBatch("5000", "", "", from, now)
	if ret == mt5webapi.MT_RET_OK && histOrders != nil {
		fmt.Printf("  Found %d history orders for login 5000:\n", len(histOrders))
		limit := 5
		if len(histOrders) < limit {
			limit = len(histOrders)
		}
		for i := 0; i < limit; i++ {
			o := histOrders[i]
			fmt.Printf("    ticket=%d type=%d symbol=%s state=%d price=%.5f\n",
				o.Order, o.Type, o.Symbol, o.State, o.PriceOrder)
		}
		if len(histOrders) > 5 {
			fmt.Printf("    ... +%d more\n", len(histOrders)-5)
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 17. HISTORY_UPDATE - Geçmiş emri değiştir
	// =============================================
	fmt.Println("--- 17. HistoryUpdateDirect ---")
	if histOrders != nil && len(histOrders) > 0 {
		ticket := histOrders[0].Order
		updateJSON := fmt.Sprintf(`{"Order":"%d","Comment":"Updated history"}`, ticket)
		updatedHist, ret := api.HistoryUpdateDirect(updateJSON)
		if ret == mt5webapi.MT_RET_OK && updatedHist != nil {
			fmt.Printf("  Updated history order %d: comment=%s\n", updatedHist.Order, updatedHist.Comment)
		} else {
			fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	} else {
		fmt.Println("  SKIPPED - no history orders")
	}
	fmt.Println()

	// =============================================
	// 18. HISTORY_DELETE - Geçmiş emri sil
	// =============================================
	fmt.Println("--- 18. HistoryDeleteDirect ---")
	fmt.Println("  SKIPPED - dangerous operation, would permanently delete history")
	fmt.Println("  Usage: ret := api.HistoryDeleteDirect(\"ticket1,ticket2\")")
	fmt.Println()

	// cleanup
	fmt.Println("--- Cleanup ---")
	if testLogin == 99999 {
		// cancel remaining orders
		remaining, _ := api.OrderGetBatch("99999", "", "")
		if remaining != nil && len(remaining) > 0 {
			tickets := ""
			for i, o := range remaining {
				if i > 0 {
					tickets += ","
				}
				tickets += fmt.Sprintf("%d", o.Order)
			}
			api.OrderCancel(tickets)
			time.Sleep(1 * time.Second)
		}
		api.UserDelete(99999)
		fmt.Println("  User 99999 cleaned up")
	}

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All 12 batch tests completed!")
	fmt.Println("================================================")
}

func groupMinJSON(name string, g *mt5webapi.MTConGroup) string {
	w := &mt5webapi.JSONWriter{}
	w.WriteBeginObject()
	w.WriteAttributeString("Group", name)
	w.WriteAttributeUint64("Server", g.Server)
	w.WriteAttributeString("Currency", g.Currency)
	w.WriteAttributeUint32("CurrencyDigits", g.CurrencyDigits)
	w.WriteAttributeUint64("PermissionsFlags", uint64(g.PermissionsFlags))
	w.WriteAttributeUint32("MarginMode", uint32(g.MarginMode))
	w.WriteAttributeFloat64("MarginCall", g.MarginCall)
	w.WriteAttributeFloat64("MarginStopOut", g.MarginStopOut)
	w.WriteEndObject()
	return w.String()
}

func roundPrice(price float64, digits uint32) float64 {
	mult := math.Pow(10, float64(digits))
	return math.Round(price*mult) / mult
}
