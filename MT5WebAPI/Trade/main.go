package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("TRADE", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Trade Example (Go)")
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

	// create test user
	testLogin := uint64(99999)
	fmt.Printf("Creating test user %d...\n", testLogin)
	newUser := mt5webapi.CreateDefault()
	newUser.Login = testLogin
	newUser.Name = "Trade Test User"
	newUser.Group = "demo\\forex"
	newUser.Leverage = 100
	newUser.MainPassword = "Test1234!"
	newUser.InvestPassword = "Invest1234!"
	_, ret = api.UserAdd(newUser)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("UserAdd failed: %s\n", mt5webapi.GetError(ret))
		api.Disconnect()
		return
	}
	fmt.Println("User created.")
	fmt.Println()

	// check initial balance
	printAccount(api, testLogin, "Initial")

	// =============================================
	// 43. TRADE_BALANCE - Bakiye işlemleri
	// =============================================
	fmt.Println("--- 43. TradeBalance ---")
	fmt.Println()

	// 1) Deposit (bakiye yatırma)
	fmt.Print("1) Deposit +10000 USD... ")
	ret = api.TradeBalance(testLogin, mt5webapi.DealBalance, 10000.00, "Deposit via Go API")
	fmt.Printf("[%d] %s\n", ret, mt5webapi.GetError(ret))
	printAccount(api, testLogin, "After deposit")

	// 2) Deposit with ticket (ticket döndüren versiyon)
	fmt.Print("2) Deposit +5000 USD (with ticket)... ")
	ticket, ret := api.TradeBalanceWithTicket(testLogin, mt5webapi.DealBalance, 5000.00, "Deposit #2 with ticket")
	fmt.Printf("[%d] %s, Ticket: %d\n", ret, mt5webapi.GetError(ret), ticket)
	printAccount(api, testLogin, "After deposit #2")

	// 3) Withdrawal (bakiye çekme)
	fmt.Print("3) Withdrawal -3000 USD... ")
	ticket, ret = api.TradeBalanceWithTicket(testLogin, mt5webapi.DealBalance, -3000.00, "Withdrawal via Go API")
	fmt.Printf("[%d] %s, Ticket: %d\n", ret, mt5webapi.GetError(ret), ticket)
	printAccount(api, testLogin, "After withdrawal")

	// 4) Credit (kredi)
	fmt.Print("4) Credit +1000 USD... ")
	ticket, ret = api.TradeBalanceWithTicket(testLogin, mt5webapi.DealCredit, 1000.00, "Credit bonus")
	fmt.Printf("[%d] %s, Ticket: %d\n", ret, mt5webapi.GetError(ret), ticket)
	printAccount(api, testLogin, "After credit")

	// 5) Charge (ek masraf)
	fmt.Print("5) Charge -200 USD... ")
	ticket, ret = api.TradeBalanceWithTicket(testLogin, mt5webapi.DealCharge, -200.00, "Service charge")
	fmt.Printf("[%d] %s, Ticket: %d\n", ret, mt5webapi.GetError(ret), ticket)
	printAccount(api, testLogin, "After charge")

	// 6) Bonus
	fmt.Print("6) Bonus +500 USD... ")
	ticket, ret = api.TradeBalanceWithTicket(testLogin, mt5webapi.DealBonus, 500.00, "Welcome bonus")
	fmt.Printf("[%d] %s, Ticket: %d\n", ret, mt5webapi.GetError(ret), ticket)
	printAccount(api, testLogin, "After bonus")

	// 7) Correction (düzeltme)
	fmt.Print("7) Correction -100 USD... ")
	ticket, ret = api.TradeBalanceWithTicket(testLogin, mt5webapi.DealCorrection, -100.00, "Balance correction")
	fmt.Printf("[%d] %s, Ticket: %d\n", ret, mt5webapi.GetError(ret), ticket)
	printAccount(api, testLogin, "After correction")

	// 8) TradeBalance with margin check
	fmt.Print("8) Deposit +2000 with margin check... ")
	ticket, ret = api.TradeBalanceWithMarginCheck(testLogin, mt5webapi.DealBalance, 2000.00, "Deposit with margin check", true)
	fmt.Printf("[%d] %s, Ticket: %d\n", ret, mt5webapi.GetError(ret), ticket)
	printAccount(api, testLogin, "Final")

	// cleanup
	fmt.Println()
	fmt.Print("Cleaning up... ")
	ret = api.UserDelete(testLogin)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("User %d deleted.\n", testLogin)
	} else {
		fmt.Printf("Delete failed: %s\n", mt5webapi.GetError(ret))
	}

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All trade tests completed!")
	fmt.Println("========================================")
}

func printAccount(api *mt5webapi.MT5WebAPI, login uint64, label string) {
	acc, ret := api.UserAccountGet(login)
	if ret == mt5webapi.MT_RET_OK && acc != nil {
		fmt.Printf("  [%s] Balance: %.2f  Credit: %.2f  Equity: %.2f\n", label, acc.Balance, acc.Credit, acc.Equity)
	}
	fmt.Println()
}
