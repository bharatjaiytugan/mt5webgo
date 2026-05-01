package main

import (
	"config"
	"fmt"
	"math"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("TRADING", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Trading Example (Go)")
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

	// use login 5000 which already has trades
	testLogin := uint64(5000)

	// get prices
	symbol := "EURUSD"
	ticks, _ := api.TickLast(symbol)
	sym, _ := api.SymbolGet(symbol)
	if ticks == nil || len(ticks) == 0 || sym == nil {
		fmt.Println("Failed to get prices")
		api.Disconnect()
		return
	}
	bid := ticks[0].Bid
	ask := ticks[0].Ask
	digits := sym.Digits
	point := sym.Point
	fmt.Printf("EURUSD: Bid=%.5f  Ask=%.5f  Digits=%d\n\n", bid, ask, digits)

	// check initial state
	posTotal, _ := api.PositionGetTotal(testLogin)
	ordTotal, _ := api.OrderGetTotal(testLogin)
	fmt.Printf("Initial: %d positions, %d pending orders\n\n", posTotal, ordTotal)

	volume := uint64(100) // 0.01 lot

	// =============================================
	// 1) Market Buy
	// =============================================
	fmt.Println("--- 1) Market Buy EURUSD 0.01 ---")
	r, ret := api.DealerBuyMarket(testLogin, symbol, volume, "Go Market Buy")
	printResult("Market Buy", r, ret)
	time.Sleep(2 * time.Second)

	// =============================================
	// 2) Market Sell
	// =============================================
	fmt.Println("--- 2) Market Sell EURUSD 0.01 ---")
	r, ret = api.DealerSellMarket(testLogin, symbol, volume, "Go Market Sell")
	printResult("Market Sell", r, ret)
	time.Sleep(2 * time.Second)

	// check positions
	posTotal2, _ := api.PositionGetTotal(testLogin)
	fmt.Printf("After market orders: %d positions (was %d)\n\n", posTotal2, posTotal)

	// =============================================
	// 3) Buy Limit (1000 points below bid)
	// =============================================
	blPrice := roundPrice(bid-1000*point, digits)
	fmt.Printf("--- 3) Buy Limit @ %.5f ---\n", blPrice)
	r, ret = api.DealerPendingOrder(testLogin, symbol, mt5webapi.OpBuyLimit, volume, blPrice, 0, 0, "Go Buy Limit")
	printResult("Buy Limit", r, ret)
	time.Sleep(2 * time.Second)

	// =============================================
	// 4) Sell Limit (1000 points above ask)
	// =============================================
	slPrice := roundPrice(ask+1000*point, digits)
	fmt.Printf("--- 4) Sell Limit @ %.5f ---\n", slPrice)
	r, ret = api.DealerPendingOrder(testLogin, symbol, mt5webapi.OpSellLimit, volume, slPrice, 0, 0, "Go Sell Limit")
	printResult("Sell Limit", r, ret)
	time.Sleep(2 * time.Second)

	// =============================================
	// 5) Buy Stop (1000 points above ask)
	// =============================================
	bsPrice := roundPrice(ask+1000*point, digits)
	fmt.Printf("--- 5) Buy Stop @ %.5f ---\n", bsPrice)
	r, ret = api.DealerPendingOrder(testLogin, symbol, mt5webapi.OpBuyStop, volume, bsPrice, 0, 0, "Go Buy Stop")
	printResult("Buy Stop", r, ret)
	time.Sleep(2 * time.Second)

	// =============================================
	// 6) Sell Stop (1000 points below bid)
	// =============================================
	ssPrice := roundPrice(bid-1000*point, digits)
	fmt.Printf("--- 6) Sell Stop @ %.5f ---\n", ssPrice)
	r, ret = api.DealerPendingOrder(testLogin, symbol, mt5webapi.OpSellStop, volume, ssPrice, 0, 0, "Go Sell Stop")
	printResult("Sell Stop", r, ret)
	time.Sleep(3 * time.Second)

	// check pending orders
	ordTotal2, _ := api.OrderGetTotal(testLogin)
	fmt.Printf("After pending orders: %d pending (was %d)\n", ordTotal2, ordTotal)

	// list pending orders
	var buyLimitTicket uint64
	if ordTotal2 > 0 {
		orders, ret := api.OrderGetPage(testLogin, 0, uint32(ordTotal2))
		if ret == mt5webapi.MT_RET_OK {
			for _, o := range orders {
				fmt.Printf("  ticket=%d type=%d symbol=%s price=%.5f comment=%s\n",
					o.Order, o.Type, o.Symbol, o.PriceOrder, o.Comment)
				if o.Type == mt5webapi.OpBuyLimit && buyLimitTicket == 0 {
					buyLimitTicket = o.Order
				}
			}
		}
	}
	fmt.Println()

	// =============================================
	// 7) Modify Buy Limit
	// =============================================
	if buyLimitTicket > 0 {
		newPrice := roundPrice(blPrice-500*point, digits)
		sl := roundPrice(newPrice-2000*point, digits)
		tp := roundPrice(newPrice+2000*point, digits)
		fmt.Printf("--- 7) Modify Order %d ---\n", buyLimitTicket)
		fmt.Printf("  Price=%.5f  SL=%.5f  TP=%.5f\n", newPrice, sl, tp)
		r, ret = api.DealerModifyOrder(testLogin, buyLimitTicket, newPrice, sl, tp)
		printResult("Modify", r, ret)
		time.Sleep(2 * time.Second)

		// =============================================
		// 8) Cancel Buy Limit
		// =============================================
		fmt.Printf("--- 8) Cancel Order %d ---\n", buyLimitTicket)
		r, ret = api.DealerCancelOrder(testLogin, buyLimitTicket)
		printResult("Cancel", r, ret)
		time.Sleep(2 * time.Second)
	} else {
		fmt.Println("--- 7) Modify --- SKIPPED (no pending)")
		fmt.Println("--- 8) Cancel --- SKIPPED (no pending)")
	}

	// final state
	fmt.Println()
	posTotal3, _ := api.PositionGetTotal(testLogin)
	ordTotal3, _ := api.OrderGetTotal(testLogin)
	fmt.Printf("Final: %d positions, %d pending orders\n", posTotal3, ordTotal3)

	acc, _ := api.UserAccountGet(testLogin)
	if acc != nil {
		fmt.Printf("Balance=%.2f  Equity=%.2f  Margin=%.2f\n", acc.Balance, acc.Equity, acc.Margin)
	}

	api.Disconnect()
	fmt.Println("\nDisconnected.")
	fmt.Println("========================================")
}

func printResult(label string, r *mt5webapi.TradeResult, ret mt5webapi.MTRetCode) {
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("  %s FAILED: [%d] %s\n\n", label, ret, mt5webapi.GetError(ret))
	} else if r != nil {
		fmt.Printf("  %s OK: ID=%d\n\n", label, r.ID)
	}
}

func roundPrice(price float64, digits uint32) float64 {
	mult := math.Pow(10, float64(digits))
	return math.Round(price*mult) / mult
}
