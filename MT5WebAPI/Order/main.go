package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("ORDER", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Order Example (Go)")
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

	// use login 2000 which has trades
	testLogin := uint64(2000)

	// =============================================
	// 18. ORDER_GET_TOTAL - Toplam emir sayısı
	// =============================================
	fmt.Println("--- 18. OrderGetTotal ---")
	total, ret := api.OrderGetTotal(testLogin)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("Login %d: %d active (pending) orders\n", testLogin, total)
	} else {
		fmt.Printf("OrderGetTotal failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 19. ORDER_GET_PAGE - Emir listesi (sayfalı)
	// =============================================
	fmt.Println("--- 19. OrderGetPage ---")
	if total > 0 {
		pageSize := uint32(10)
		if total < pageSize {
			pageSize = total
		}
		orders, ret := api.OrderGetPage(testLogin, 0, pageSize)
		if ret == mt5webapi.MT_RET_OK && orders != nil {
			fmt.Printf("Got %d orders (of %d total):\n", len(orders), total)
			for i, o := range orders {
				fmt.Printf("\n  Order #%d:\n", i+1)
				printOrder(o)
			}
		} else {
			fmt.Printf("OrderGetPage failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	} else {
		fmt.Println("No active orders to list.")
		fmt.Println("(Active orders = pending orders like Buy Limit, Sell Stop, etc.)")
	}
	fmt.Println()

	// =============================================
	// 17. ORDER_GET - Emir detayı al (by ticket)
	// =============================================
	fmt.Println("--- 17. OrderGet ---")
	if total > 0 {
		orders, ret := api.OrderGetPage(testLogin, 0, 1)
		if ret == mt5webapi.MT_RET_OK && orders != nil && len(orders) > 0 {
			ticket := orders[0].Order
			order, ret := api.OrderGet(ticket)
			if ret == mt5webapi.MT_RET_OK && order != nil {
				fmt.Printf("OrderGet(ticket=%d):\n", ticket)
				printOrderFull(order)
			} else {
				fmt.Printf("OrderGet failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		}
	} else {
		fmt.Println("No active orders. Trying history order...")
		// try history as fallback
		histTotal, ret := api.HistoryGetTotal(testLogin, 0, time.Now().Unix())
		if ret == mt5webapi.MT_RET_OK && histTotal > 0 {
			fmt.Printf("Found %d history orders\n", histTotal)
			histOrders, ret := api.HistoryGetPage(testLogin, 0, time.Now().Unix(), 0, 3)
			if ret == mt5webapi.MT_RET_OK && histOrders != nil {
				for i, o := range histOrders {
					fmt.Printf("\n  History Order #%d:\n", i+1)
					printOrder(o)
				}

				// get single history order
				if len(histOrders) > 0 {
					ticket := histOrders[0].Order
					fmt.Printf("\n  HistoryGet(ticket=%d):\n", ticket)
					histOrder, ret := api.HistoryGet(ticket)
					if ret == mt5webapi.MT_RET_OK && histOrder != nil {
						printOrderFull(histOrder)
					} else {
						fmt.Printf("  HistoryGet failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
					}
				}
			} else {
				fmt.Printf("HistoryGetPage failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		} else {
			fmt.Printf("HistoryGetTotal failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	}
	fmt.Println()

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All order tests completed!")
	fmt.Println("========================================")
}

func printOrder(o *mt5webapi.MTOrder) {
	fmt.Printf("    Ticket:  %d\n", o.Order)
	fmt.Printf("    Symbol:  %s\n", o.Symbol)
	fmt.Printf("    Type:    %s\n", orderTypeStr(o.Type))
	fmt.Printf("    State:   %s\n", orderStateStr(o.State))
	fmt.Printf("    Volume:  %d\n", o.VolumeInitialExt)
	fmt.Printf("    Price:   %.5f\n", o.PriceOrder)
	fmt.Printf("    SL:      %.5f\n", o.PriceSL)
	fmt.Printf("    TP:      %.5f\n", o.PriceTP)
	fmt.Printf("    Comment: %s\n", o.Comment)
}

func printOrderFull(o *mt5webapi.MTOrder) {
	fmt.Printf("    Ticket:       %d\n", o.Order)
	fmt.Printf("    ExternalID:   %s\n", o.ExternalID)
	fmt.Printf("    Login:        %d\n", o.Login)
	fmt.Printf("    Dealer:       %d\n", o.Dealer)
	fmt.Printf("    Symbol:       %s\n", o.Symbol)
	fmt.Printf("    Digits:       %d\n", o.Digits)
	fmt.Printf("    Type:         %s\n", orderTypeStr(o.Type))
	fmt.Printf("    TypeFill:     %d\n", o.TypeFill)
	fmt.Printf("    TypeTime:     %d\n", o.TypeTime)
	fmt.Printf("    State:        %s\n", orderStateStr(o.State))
	fmt.Printf("    Reason:       %d\n", o.Reason)
	fmt.Printf("    PriceOrder:   %.5f\n", o.PriceOrder)
	fmt.Printf("    PriceTrigger: %.5f\n", o.PriceTrigger)
	fmt.Printf("    PriceCurrent: %.5f\n", o.PriceCurrent)
	fmt.Printf("    PriceSL:      %.5f\n", o.PriceSL)
	fmt.Printf("    PriceTP:      %.5f\n", o.PriceTP)
	fmt.Printf("    VolumeInit:   %d\n", o.VolumeInitialExt)
	fmt.Printf("    VolumeCurr:   %d\n", o.VolumeCurrentExt)
	fmt.Printf("    ExpertID:     %d\n", o.ExpertID)
	fmt.Printf("    PositionByID: %d\n", o.PositionByID)
	fmt.Printf("    Comment:      %s\n", o.Comment)
	fmt.Printf("    ContractSize: %.2f\n", o.ContractSize)
}

func orderTypeStr(t mt5webapi.EnOrderType) string {
	switch t {
	case mt5webapi.OpBuy:
		return "Buy"
	case mt5webapi.OpSell:
		return "Sell"
	case mt5webapi.OpBuyLimit:
		return "Buy Limit"
	case mt5webapi.OpSellLimit:
		return "Sell Limit"
	case mt5webapi.OpBuyStop:
		return "Buy Stop"
	case mt5webapi.OpSellStop:
		return "Sell Stop"
	case mt5webapi.OpBuyStopLimit:
		return "Buy Stop Limit"
	case mt5webapi.OpSellStopLimit:
		return "Sell Stop Limit"
	case mt5webapi.OpCloseBy:
		return "Close By"
	}
	return fmt.Sprintf("Unknown(%d)", t)
}

func orderStateStr(s mt5webapi.EnOrderState) string {
	switch s {
	case mt5webapi.OrderStateStarted:
		return "Started"
	case mt5webapi.OrderStatePlaced:
		return "Placed"
	case mt5webapi.OrderStateCanceled:
		return "Canceled"
	case mt5webapi.OrderStatePartial:
		return "Partial"
	case mt5webapi.OrderStateFilled:
		return "Filled"
	case mt5webapi.OrderStateRejected:
		return "Rejected"
	case mt5webapi.OrderStateExpired:
		return "Expired"
	}
	return fmt.Sprintf("Unknown(%d)", s)
}
