package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("GROUP", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Group Example (Go)")
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
	// 32. GROUP_TOTAL - Grup sayısı
	// =============================================
	fmt.Println("--- 32. GroupTotal ---")
	total, ret := api.GroupTotal()
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("Total groups: %d\n", total)
	} else {
		fmt.Printf("GroupTotal failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 33. GROUP_NEXT - Sonraki grubu al
	// =============================================
	fmt.Println("--- 33. GroupNext ---")
	if total > 0 {
		count := total
		if count > 10 {
			count = 10
		}
		for i := 0; i < count; i++ {
			group, ret := api.GroupNext(uint32(i))
			if ret == mt5webapi.MT_RET_OK && group != nil {
				fmt.Printf("  [%d] %-30s Currency: %-5s Leverage: %d  MarginCall: %.0f%%  StopOut: %.0f%%\n",
					i, group.Group, group.Currency, group.DemoLeverage, group.MarginCall, group.MarginStopOut)
			} else {
				fmt.Printf("  [%d] failed: %s\n", i, mt5webapi.GetError(ret))
			}
		}
		if total > 10 {
			fmt.Printf("  ... and %d more groups\n", total-10)
		}
	}
	fmt.Println()

	// =============================================
	// 34. GROUP_GET - Grup bilgisi al
	// =============================================
	fmt.Println("--- 34. GroupGet ---")
	groupName := "demo\\forex"
	group, ret := api.GroupGet(groupName)
	if ret == mt5webapi.MT_RET_OK && group != nil {
		fmt.Printf("Group:              %s\n", group.Group)
		fmt.Printf("Server:             %d\n", group.Server)
		fmt.Printf("Currency:           %s\n", group.Currency)
		fmt.Printf("CurrencyDigits:     %d\n", group.CurrencyDigits)
		fmt.Printf("Company:            %s\n", group.Company)
		fmt.Printf("CompanyPage:        %s\n", group.CompanyPage)
		fmt.Printf("CompanyEmail:       %s\n", group.CompanyEmail)
		fmt.Printf("AuthMode:           %d\n", group.AuthMode)
		fmt.Printf("AuthPasswordMin:    %d\n", group.AuthPasswordMin)
		fmt.Printf("PermissionsFlags:   %d\n", group.PermissionsFlags)
		fmt.Printf("MarginMode:         %d\n", group.MarginMode)
		fmt.Printf("MarginCall:         %.2f\n", group.MarginCall)
		fmt.Printf("MarginStopOut:      %.2f\n", group.MarginStopOut)
		fmt.Printf("MarginSOMode:       %d\n", group.MarginSOMode)
		fmt.Printf("MarginFreeMode:     %d\n", group.MarginFreeMode)
		fmt.Printf("TradeFlags:         %d\n", group.TradeFlags)
		fmt.Printf("TradeInterestRate:  %.2f\n", group.TradeInterestrate)
		fmt.Printf("TradeVirtualCredit: %.2f\n", group.TradeVirtualCredit)
		fmt.Printf("DemoLeverage:       %d\n", group.DemoLeverage)
		fmt.Printf("DemoDeposit:        %.2f\n", group.DemoDeposit)
		fmt.Printf("LimitHistory:       %d\n", group.LimitHistory)
		fmt.Printf("LimitOrders:        %d\n", group.LimitOrders)
		fmt.Printf("LimitSymbols:       %d\n", group.LimitSymbols)
		fmt.Printf("LimitPositions:     %d\n", group.LimitPositions)
		fmt.Printf("NewsMode:           %d\n", group.NewsMode)
		fmt.Printf("MailMode:           %d\n", group.MailMode)
		fmt.Printf("ReportsMode:        %d\n", group.ReportsMode)
	} else {
		fmt.Printf("GroupGet '%s' failed: [%d] %s\n", groupName, ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 35. GROUP_ADD - Yeni grup oluştur
	// =============================================
	fmt.Println("--- 35. GroupAdd ---")
	// clone an existing group and change the name
	templateGroup, ret := api.GroupGet("demo\\forex")
	if ret != mt5webapi.MT_RET_OK || templateGroup == nil {
		fmt.Printf("Cannot get template group: %s\n", mt5webapi.GetError(ret))
		api.Disconnect()
		return
	}
	// clone it with a new name
	newGroup := templateGroup
	newGroup.Group = "demo\\gotest"

	createdGroup, ret := api.GroupAdd(newGroup)
	if ret == mt5webapi.MT_RET_OK && createdGroup != nil {
		fmt.Printf("Group created!\n")
		fmt.Printf("Name:        %s\n", createdGroup.Group)
		fmt.Printf("Currency:    %s\n", createdGroup.Currency)
		fmt.Printf("Company:     %s\n", createdGroup.Company)
		fmt.Printf("MarginCall:  %.2f\n", createdGroup.MarginCall)
		fmt.Printf("StopOut:     %.2f\n", createdGroup.MarginStopOut)
		fmt.Printf("Leverage:    %d\n", createdGroup.DemoLeverage)
	} else {
		fmt.Printf("GroupAdd failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// verify group exists
	fmt.Println("Verifying group exists...")
	verifyGroup, ret := api.GroupGet("demo\\gotest")
	if ret == mt5webapi.MT_RET_OK && verifyGroup != nil {
		fmt.Printf("Verified: %s (Currency: %s, Company: %s)\n", verifyGroup.Group, verifyGroup.Currency, verifyGroup.Company)
	} else {
		fmt.Printf("Verify failed: %s\n", mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 36. GROUP_DELETE - Grup sil
	// =============================================
	fmt.Println("--- 36. GroupDelete ---")
	ret = api.GroupDelete("demo\\gotest")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Println("Group 'demo\\gotest' deleted successfully!")
	} else {
		fmt.Printf("GroupDelete failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	// verify deleted
	_, ret = api.GroupGet("demo\\gotest")
	fmt.Printf("GroupGet after delete: [%d] %s\n", ret, mt5webapi.GetError(ret))
	fmt.Println()

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All group tests completed!")
	fmt.Println("========================================")
}
