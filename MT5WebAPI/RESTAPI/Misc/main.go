package main

import (
	"config"
	"encoding/json"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("MISC", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Misc Handlers Example")
	fmt.Println("  Tests: Logger, UserTotal, UserGroup, Mail,")
	fmt.Println("         News, TestAccess, Backups, DailyGetLight")
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

	now := time.Now()

	// =============================================
	// 1. LoggerServerRequest (last 1 hour)
	// =============================================
	fmt.Println("--- 1. LoggerServerRequest ---")
	fromLog := now.Add(-1 * time.Hour).Unix()
	toLog := now.Unix()
	logs, ret := api.LoggerServerRequest("0", "0", fromLog, toLog, "")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Log entries: %d\n", len(logs))
		limit := 3
		if len(logs) < limit {
			limit = len(logs)
		}
		for i := 0; i < limit; i++ {
			printMap(fmt.Sprintf("  Entry[%d]", i), logs[i])
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 2. UserTotal
	// =============================================
	fmt.Println("--- 2. UserTotal ---")
	userTotal, ret := api.UserTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Total users: %d\n", userTotal)
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 3. UserGroup (get group of login 5000)
	// =============================================
	fmt.Println("--- 3. UserGroup ---")
	groupName, ret := api.UserGroup(5000)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Login 5000 group: %s\n", groupName)
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 4. MailGet (last message)
	// =============================================
	fmt.Println("--- 4. MailGet ---")
	mails, ret := api.MailGet(cfg.Login, 0, 1)
	if ret == mt5webapi.MT_RET_OK && len(mails) > 0 {
		printMap("  Last mail", mails[0])
	} else if ret == mt5webapi.MT_RET_OK {
		fmt.Println("  No mail messages found")
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 5. NewsGet - REST API only, not supported on WebSocket binary protocol
	// =============================================
	fmt.Println("--- 5. NewsGet ---")
	fmt.Println("  SKIPPED - NEWS_GET is REST API only (not in C# original SDK)")
	fmt.Println("  Usage: news, ret := api.NewsGet(\"\", 10)")
	fmt.Println()

	// =============================================
	// 6. TestAccess
	// =============================================
	fmt.Println("--- 6. TestAccess ---")
	ret = api.TestAccess()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Println("  TestAccess: OK")
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 7. PositionBackupList (last 30 days)
	// =============================================
	fmt.Println("--- 7. PositionBackupList ---")
	fromBackup := now.AddDate(0, 0, -30).Unix()
	toBackup := now.Unix()
	posDates, ret := api.PositionBackupList(fromBackup, toBackup)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Position backup dates: %d\n", len(posDates))
		limit := 5
		if len(posDates) < limit {
			limit = len(posDates)
		}
		for i := 0; i < limit; i++ {
			t := time.Unix(posDates[i], 0).UTC()
			fmt.Printf("  [%d] %s (%d)\n", i, t.Format("2006-01-02 15:04:05"), posDates[i])
		}
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 8. DealBackupList (last 30 days)
	// =============================================
	fmt.Println("--- 8. DealBackupList ---")
	dealDates, ret := api.DealBackupList(fromBackup, toBackup)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Deal backup dates: %d\n", len(dealDates))
		limit := 5
		if len(dealDates) < limit {
			limit = len(dealDates)
		}
		for i := 0; i < limit; i++ {
			t := time.Unix(dealDates[i], 0).UTC()
			fmt.Printf("  [%d] %s (%d)\n", i, t.Format("2006-01-02 15:04:05"), dealDates[i])
		}
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 9. DailyGetLight (login=5000, last 7 days)
	// =============================================
	fmt.Println("--- 9. DailyGetLight ---")
	fromDaily := now.AddDate(0, 0, -7).Unix()
	toDaily := now.Unix()
	dailies, ret := api.DailyGetLight(5000, fromDaily, toDaily)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  Daily reports for login 5000: %d\n", len(dailies))
		limit := 3
		if len(dailies) < limit {
			limit = len(dailies)
		}
		for i := 0; i < limit; i++ {
			d := dailies[i]
			fmt.Printf("  [%d] Login=%d Group=%s Balance=%.2f Equity=%.2f\n",
				i, d.Login, d.Group, d.Balance, d.Equity)
		}
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All Misc tests completed!")
	fmt.Println("================================================")
}

func printMap(label string, m map[string]interface{}) {
	data, _ := json.MarshalIndent(m, "  ", "  ")
	fmt.Printf("%s: %s\n", label, string(data))
}
