package main

import (
	"config"
	"fmt"
	"time"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("USER_BATCH", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - User Batch/Backup Example (Go)")
	fmt.Println("  Tests: 31, 32a, 32b")
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
	// 31. USER_GET_BATCH - by logins
	// =============================================
	fmt.Println("--- 31a. UserGetBatch (by logins) ---")
	users, ret := api.UserGetBatch("5000,5005", "")
	if ret == mt5webapi.MT_RET_OK && users != nil {
		fmt.Printf("  Found %d users:\n", len(users))
		for _, u := range users {
			printUser(u)
		}
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 31b. USER_GET_BATCH - by demo group
	// =============================================
	fmt.Println("--- 31b. UserGetBatch (demo group) ---")
	users, ret = api.UserGetBatch("", "demo\\*")
	if ret == mt5webapi.MT_RET_OK && users != nil {
		fmt.Printf("  Found %d users in demo\\*:\n", len(users))
		for i, u := range users {
			if i >= 10 {
				fmt.Printf("    ... +%d more\n", len(users)-10)
				break
			}
			printUser(u)
		}
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 31c. USER_GET_BATCH - by managers group
	// =============================================
	fmt.Println("--- 31c. UserGetBatch (managers group) ---")
	users, ret = api.UserGetBatch("", "manager\\*")
	if ret == mt5webapi.MT_RET_OK && users != nil {
		fmt.Printf("  Found %d users in manager\\*:\n", len(users))
		for i, u := range users {
			if i >= 10 {
				fmt.Printf("    ... +%d more\n", len(users)-10)
				break
			}
			printUser(u)
		}
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 32a. USER_BACKUP_LIST
	// =============================================
	fmt.Println("--- 32a. UserBackupList (last 30 days) ---")
	now := time.Now().Unix() + 86400
	from30 := time.Now().AddDate(0, -1, 0).Unix()
	backupDates, ret := api.UserBackupList(from30, now)
	if ret == mt5webapi.MT_RET_OK && backupDates != nil {
		fmt.Printf("  Found %d backup dates:\n", len(backupDates))
		for i, d := range backupDates {
			t := time.Unix(d, 0)
			fmt.Printf("    [%d] %d (%s)\n", i, d, t.Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 32b. USER_BACKUP_GET
	// =============================================
	fmt.Println("--- 32b. UserBackupGet ---")
	if backupDates != nil && len(backupDates) > 0 {
		// find the latest non-zero backup date
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
			fmt.Println("  Querying for login 5000...")
			users, ret := api.UserBackupGet(backupDate, 5000)
			if ret == mt5webapi.MT_RET_OK && users != nil {
				fmt.Printf("  Found %d users in backup:\n", len(users))
				for i, u := range users {
					if i >= 5 {
						fmt.Printf("    ... +%d more\n", len(users)-5)
						break
					}
					printUser(u)
				}
			} else {
				fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
			}
		} else {
			fmt.Println("  No valid backup dates found")
		}
	} else {
		fmt.Println("  SKIPPED - no backup dates available")
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
}

func printUser(u *mt5webapi.MTUser) {
	fmt.Printf("    Login=%-8d Name=%-20s Group=%-20s Leverage=%-5d Balance=%.2f\n",
		u.Login, u.Name, u.Group, u.Leverage, u.Balance)
}
