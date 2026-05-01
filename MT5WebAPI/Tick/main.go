package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("TICK", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Tick Example (Go)")
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

	symbols := []string{"EURUSD", "GBPUSD", "USDJPY", "XAUUSD", "WTI"}

	// =============================================
	// 29. TICK_LAST - Son fiyat (bid/ask/last)
	// =============================================
	fmt.Println("--- 29. TickLast ---")
	for _, symbol := range symbols {
		ticks, ret := api.TickLast(symbol)
		if ret == mt5webapi.MT_RET_OK && ticks != nil && len(ticks) > 0 {
			t := ticks[0]
			ts := time.Unix(int64(t.Datetime), 0).UTC().Format("15:04:05")
			fmt.Printf("  %-8s  Bid: %-10s  Ask: %-10s  Last: %-10s  Vol: %d  Time: %s\n",
				t.Symbol,
				fmtPrice(t.Bid, t.Digits),
				fmtPrice(t.Ask, t.Digits),
				fmtPrice(t.Last, t.Digits),
				t.Volume,
				ts)
		} else {
			fmt.Printf("  %-8s  failed: %s\n", symbol, mt5webapi.GetError(ret))
		}
	}
	fmt.Println()

	// =============================================
	// 30. TICK_LAST_GROUP - Grup bazlı son fiyat
	// =============================================
	fmt.Println("--- 30. TickLastGroup ---")
	groups := []string{"demo\\forex", "managers\\*"}
	for _, group := range groups {
		fmt.Printf("  Group: %s\n", group)
		for _, symbol := range symbols[:3] {
			ticks, ret := api.TickLastGroup(symbol, group)
			if ret == mt5webapi.MT_RET_OK && ticks != nil && len(ticks) > 0 {
				t := ticks[0]
				fmt.Printf("    %-8s  Bid: %-10s  Ask: %-10s\n",
					t.Symbol,
					fmtPrice(t.Bid, t.Digits),
					fmtPrice(t.Ask, t.Digits))
			} else {
				fmt.Printf("    %-8s  failed: %s\n", symbol, mt5webapi.GetError(ret))
			}
		}
		fmt.Println()
	}

	// =============================================
	// 31. TICK_STAT - Fiyat istatistikleri
	// =============================================
	fmt.Println("--- 31. TickStat ---")
	for _, symbol := range symbols {
		stats, ret := api.TickStat(symbol)
		if ret == mt5webapi.MT_RET_OK && stats != nil && len(stats) > 0 {
			s := stats[0]
			fmt.Printf("  %s:\n", s.Symbol)
			fmt.Printf("    Bid:       %-10s  Low: %-10s  High: %-10s  Dir: %s\n",
				fmtPrice(s.Bid, s.Digits),
				fmtPrice(s.BidLow, s.Digits),
				fmtPrice(s.BidHigh, s.Digits),
				dirStr(s.BidDir))
			fmt.Printf("    Ask:       %-10s  Low: %-10s  High: %-10s  Dir: %s\n",
				fmtPrice(s.Ask, s.Digits),
				fmtPrice(s.AskLow, s.Digits),
				fmtPrice(s.AskHigh, s.Digits),
				dirStr(s.AskDir))
			fmt.Printf("    Last:      %-10s  Low: %-10s  High: %-10s  Dir: %s\n",
				fmtPrice(s.Last, s.Digits),
				fmtPrice(s.LastLow, s.Digits),
				fmtPrice(s.LastHigh, s.Digits),
				dirStr(s.LastDir))
			fmt.Printf("    Volume:    %d  VolumeReal: %.2f\n", s.Volume, s.VolumeReal)
			fmt.Printf("    VolLow:    %d  VolHigh: %d  Dir: %s\n", s.VolumeLow, s.VolumeHigh, dirStr(s.VolumeDir))
			fmt.Printf("    Deals:     %d\n", s.TradeDeals)
			fmt.Printf("    BuyOrd:    %d  BuyVol: %d\n", s.TradeBuyOrders, s.TradeBuyVolume)
			fmt.Printf("    SellOrd:   %d  SellVol: %d\n", s.TradeSellOrders, s.TradeSellVolume)
			fmt.Printf("    PriceOpen: %-10s  PriceClose: %s\n",
				fmtPrice(s.PriceOpen, s.Digits),
				fmtPrice(s.PriceClose, s.Digits))
			fmt.Printf("    Change:    %.2f%%  Volatility: %.2f%%\n", s.PriceChange, s.PriceVolatility)
			fmt.Println()
		} else {
			fmt.Printf("  %-8s  failed: %s\n", symbol, mt5webapi.GetError(ret))
			fmt.Println()
		}
	}

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All tick tests completed!")
	fmt.Println("========================================")
}

func fmtPrice(price float64, digits uint32) string {
	return fmt.Sprintf("%.*f", digits, price)
}

func dirStr(d mt5webapi.EnTickStatDirection) string {
	switch d {
	case mt5webapi.TickStatDirectionNone:
		return "-"
	case mt5webapi.TickStatDirectionUp:
		return "UP"
	case mt5webapi.TickStatDirectionDown:
		return "DOWN"
	}
	return "?"
}
