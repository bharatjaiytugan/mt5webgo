package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("SYMBOL", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Symbol Example (Go)")
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
	// 37. SYMBOL_TOTAL - Sembol sayısı
	// =============================================
	fmt.Println("--- 37. SymbolTotal ---")
	total, ret := api.SymbolTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("Total symbols: %d\n", total)
	} else {
		fmt.Printf("SymbolTotal failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 38. SYMBOL_NEXT - Sonraki sembolü al
	// =============================================
	fmt.Println("--- 38. SymbolNext ---")
	count := total
	if count > 15 {
		count = 15
	}
	for i := 0; i < count; i++ {
		sym, ret := api.SymbolNext(i)
		if ret == mt5webapi.MT_RET_OK && sym != nil {
			fmt.Printf("  [%2d] %-12s  Path: %-25s  Digits: %d  Spread: %d  Mode: %d\n",
				i, sym.Symbol, sym.Path, sym.Digits, sym.Spread, sym.TradeMode)
		} else {
			fmt.Printf("  [%2d] failed: %s\n", i, mt5webapi.GetError(ret))
		}
	}
	if total > 15 {
		fmt.Printf("  ... and %d more symbols\n", total-15)
	}
	fmt.Println()

	// =============================================
	// 39. SYMBOL_GET - Sembol bilgisi al
	// =============================================
	fmt.Println("--- 39. SymbolGet (EURUSD) ---")
	sym, ret := api.SymbolGet("EURUSD")
	if ret == mt5webapi.MT_RET_OK && sym != nil {
		printSymbol(sym)
	} else {
		fmt.Printf("SymbolGet failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 40. SYMBOL_GET_GROUP - Grup bazlı sembol
	// =============================================
	fmt.Println("--- 40. SymbolGetGroup (EURUSD, demo\\forex) ---")
	symGroup, ret := api.SymbolGetGroup("EURUSD", "demo\\forex")
	if ret == mt5webapi.MT_RET_OK && symGroup != nil {
		fmt.Printf("Symbol:       %s\n", symGroup.Symbol)
		fmt.Printf("Path:         %s\n", symGroup.Path)
		fmt.Printf("TradeMode:    %d\n", symGroup.TradeMode)
		fmt.Printf("ExecMode:     %d\n", symGroup.ExecMode)
		fmt.Printf("Spread:       %d\n", symGroup.Spread)
		fmt.Printf("ContractSize: %.0f\n", symGroup.ContractSize)
		fmt.Printf("StopsLevel:   %d\n", symGroup.StopsLevel)
		fmt.Printf("FreezeLevel:  %d\n", symGroup.FreezeLevel)
	} else {
		fmt.Printf("SymbolGetGroup failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 41. SYMBOL_ADD - Yeni sembol ekle
	// =============================================
	fmt.Println("--- 41. SymbolAdd ---")
	// clone EURUSD as template
	templateSym, ret := api.SymbolGet("EURUSD")
	if ret == mt5webapi.MT_RET_OK && templateSym != nil {
		templateSym.Symbol = "GOTEST"
		templateSym.Path = "Forex\\GOTEST"
		templateSym.Description = "Go API Test Symbol"

		newSym, ret := api.SymbolAdd(templateSym)
		if ret == mt5webapi.MT_RET_OK && newSym != nil {
			fmt.Printf("Symbol created!\n")
			fmt.Printf("Symbol:       %s\n", newSym.Symbol)
			fmt.Printf("Path:         %s\n", newSym.Path)
			fmt.Printf("Description:  %s\n", newSym.Description)
			fmt.Printf("Digits:       %d\n", newSym.Digits)

			// verify
			verifySym, ret := api.SymbolGet("GOTEST")
			if ret == mt5webapi.MT_RET_OK && verifySym != nil {
				fmt.Printf("Verified:     %s exists\n", verifySym.Symbol)
			}
			fmt.Println()

			// =============================================
			// 42. SYMBOL_DELETE - Sembol sil
			// =============================================
			fmt.Println("--- 42. SymbolDelete ---")
			ret = api.SymbolDelete("GOTEST")
			if ret == mt5webapi.MT_RET_OK {
				fmt.Println("Symbol 'GOTEST' deleted!")
			} else {
				fmt.Printf("SymbolDelete failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
			// verify deleted
			_, ret = api.SymbolGet("GOTEST")
			fmt.Printf("SymbolGet after delete: [%d] %s\n", ret, mt5webapi.GetError(ret))
		} else {
			fmt.Printf("SymbolAdd failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	} else {
		fmt.Printf("Cannot get template symbol: %s\n", mt5webapi.GetError(ret))
	}
	fmt.Println()

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All symbol tests completed!")
	fmt.Println("========================================")
}

func printSymbol(s *mt5webapi.MTConSymbol) {
	fmt.Printf("Symbol:         %s\n", s.Symbol)
	fmt.Printf("Path:           %s\n", s.Path)
	fmt.Printf("ISIN:           %s\n", s.ISIN)
	fmt.Printf("Description:    %s\n", s.Description)
	fmt.Printf("CurrencyBase:   %s\n", s.CurrencyBase)
	fmt.Printf("CurrencyProfit: %s\n", s.CurrencyProfit)
	fmt.Printf("CurrencyMargin: %s\n", s.CurrencyMargin)
	fmt.Printf("Digits:         %d\n", s.Digits)
	fmt.Printf("Point:          %f\n", s.Point)
	fmt.Printf("ContractSize:   %.0f\n", s.ContractSize)
	fmt.Printf("TradeMode:      %d\n", s.TradeMode)
	fmt.Printf("ExecMode:       %d\n", s.ExecMode)
	fmt.Printf("CalcMode:       %d\n", s.CalcMode)
	fmt.Printf("Spread:         %d\n", s.Spread)
	fmt.Printf("SpreadBalance:  %d\n", s.SpreadBalance)
	fmt.Printf("TickValue:      %.5f\n", s.TickValue)
	fmt.Printf("TickSize:       %.5f\n", s.TickSize)
	fmt.Printf("StopsLevel:     %d\n", s.StopsLevel)
	fmt.Printf("FreezeLevel:    %d\n", s.FreezeLevel)
	fmt.Printf("SwapMode:       %d\n", s.SwapMode)
	fmt.Printf("SwapLong:       %.2f\n", s.SwapLong)
	fmt.Printf("SwapShort:      %.2f\n", s.SwapShort)
	fmt.Printf("Swap3Day:       %d\n", s.Swap3Day)
	fmt.Printf("Sector:         %d\n", s.Sector)
	fmt.Printf("Industry:       %d\n", s.Industry)
	fmt.Printf("Category:       %s\n", s.Category)
	fmt.Printf("Exchange:       %s\n", s.Exchange)
}
