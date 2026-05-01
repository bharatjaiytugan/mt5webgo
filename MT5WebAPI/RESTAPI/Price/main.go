package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("PRICE", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Live Price Streaming Demo (Go)")
	fmt.Println("================================================")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		return
	}

	// --- 1. Connect ---
	fmt.Print("Connecting... ")
	ret := api.ConnectDefault(cfg.Server, cfg.Port, cfg.Login, cfg.Password, mt5webapi.PUMP_MODE_NONE)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("FAILED! [%d] %s\n", ret, mt5webapi.GetError(ret))
		return
	}
	fmt.Println("OK")
	fmt.Println()

	symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}

	// --- 2. One-shot TickLast ---
	fmt.Println("--- TickLast (one-shot) ---")
	for _, sym := range symbols {
		ticks, ret := api.TickLast(sym)
		if ret == mt5webapi.MT_RET_OK && len(ticks) > 0 {
			t := ticks[0]
			fmt.Printf("  %-8s Bid=%.5f  Ask=%.5f  Last=%.5f\n", t.Symbol, t.Bid, t.Ask, t.Last)
		} else {
			fmt.Printf("  %-8s no data [%d]\n", sym, ret)
		}
	}
	fmt.Println()

	// --- 3. One-shot TickStat ---
	fmt.Println("--- TickStat EURUSD (one-shot) ---")
	stats, ret := api.TickStat("EURUSD")
	if ret == mt5webapi.MT_RET_OK && len(stats) > 0 {
		s := stats[0]
		fmt.Printf("  Bid: %.5f (Low=%.5f High=%.5f)\n", s.Bid, s.BidLow, s.BidHigh)
		fmt.Printf("  Ask: %.5f (Low=%.5f High=%.5f)\n", s.Ask, s.AskLow, s.AskHigh)
		fmt.Printf("  Open=%.5f  Close=%.5f\n", s.PriceOpen, s.PriceClose)
	} else {
		fmt.Printf("  no data [%d]\n", ret)
	}
	fmt.Println()

	// --- 4. Start PriceStreamer ---
	fmt.Println("--- Starting PriceStreamer ---")
	fmt.Printf("  Symbols: %v\n", symbols)
	fmt.Println("  Interval: 1ms")
	fmt.Println()

	streamer := api.NewPriceStreamer(func(update *mt5webapi.PriceUpdate) {
		ts := time.Unix(int64(update.Datetime), 0)
		fmt.Printf("[%s] %-8s Bid=%.5f  Ask=%.5f  Spread=%.5f\n",
			ts.Format("15:04:05.000"), update.Symbol, update.Bid, update.Ask, update.Spread)
	})
	streamer.AddSymbols(symbols...)
	streamer.SetInterval(1 * time.Millisecond)
	streamer.Start()

	// --- 5. Stream for 10 seconds ---
	fmt.Println("--- Streaming prices for 10 seconds ---")
	time.Sleep(10 * time.Second)
	fmt.Println()

	// --- 6. Show GetLastPrice for each symbol ---
	fmt.Println("--- Last cached prices ---")
	for _, sym := range symbols {
		p := streamer.GetLastPrice(sym)
		if p != nil {
			fmt.Printf("  %-8s Bid=%.5f  Ask=%.5f  Spread=%.5f\n", p.Symbol, p.Bid, p.Ask, p.Spread)
		} else {
			fmt.Printf("  %-8s no cached price\n", sym)
		}
	}
	fmt.Println()

	// --- 7. Start BarBuilder for EURUSD ---
	fmt.Println("--- Starting BarBuilder for EURUSD (1-minute bars) ---")
	barBuilder := mt5webapi.NewBarBuilder("EURUSD", 1*time.Minute, func(bar *mt5webapi.MTBar) {
		bt := time.Unix(bar.Datetime, 0)
		fmt.Printf("[BAR] EURUSD %s  O=%.5f H=%.5f L=%.5f C=%.5f  Ticks=%d\n",
			bt.Format("15:04"), bar.Open, bar.High, bar.Low, bar.Close, bar.TickVol)
	})

	// Replace the streamer callback to also feed the bar builder
	streamer.Stop()
	streamer = api.NewPriceStreamer(func(update *mt5webapi.PriceUpdate) {
		ts := time.Unix(int64(update.Datetime), 0)
		fmt.Printf("[%s] %-8s Bid=%.5f  Ask=%.5f  Spread=%.5f\n",
			ts.Format("15:04:05.000"), update.Symbol, update.Bid, update.Ask, update.Spread)
		// Feed EURUSD ticks to the bar builder
		if update.Symbol == "EURUSD" {
			barBuilder.OnTick(update.Bid, ts)
		}
	})
	streamer.AddSymbols(symbols...)
	streamer.SetInterval(1 * time.Millisecond)
	streamer.Start()

	// --- 8. Stream for another 10 seconds with bar building ---
	fmt.Println("--- Streaming with bar building for 10 seconds ---")
	time.Sleep(10 * time.Second)

	// Show current incomplete bar
	bar := barBuilder.CurrentBar()
	if bar != nil {
		bt := time.Unix(bar.Datetime, 0)
		fmt.Printf("\n[BAR-CURRENT] EURUSD %s  O=%.5f H=%.5f L=%.5f C=%.5f  Ticks=%d\n",
			bt.Format("15:04"), bar.Open, bar.High, bar.Low, bar.Close, bar.TickVol)
	}
	fmt.Println()

	// --- 9. Stop streamer ---
	fmt.Println("--- Stopping PriceStreamer ---")
	streamer.Stop()
	fmt.Printf("  Running: %v\n", streamer.IsRunning())
	fmt.Println()

	// --- 10. Disconnect ---
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
}
