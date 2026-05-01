package main

import (
	"config"
	"encoding/json"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("SETTINGS", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Settings Example")
	fmt.Println("  Tests: 39-40 (SettingGet, SettingSet, SettingDelete)")
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

	section := "web"
	key := "gotest.json"

	// =============================================
	// Step 1: SETTING_SET - Write a test settings file
	// =============================================
	fmt.Println("--- 39. SettingSet (write test file) ---")
	initialContent := `{"AppName":"GoAPI","Version":"1.0","Debug":true}`
	ret = api.SettingSet(section, key, initialContent)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  SettingSet OK: section=%s key=%s\n", section, key)
		fmt.Printf("  Content: %s\n", initialContent)
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// Step 2: SETTING_GET - Read back the file
	// =============================================
	fmt.Println("--- 40. SettingGet (read back test file) ---")
	result, ret := api.SettingGet(section, key)
	if ret == mt5webapi.MT_RET_OK && result != nil {
		resultJSON, _ := json.Marshal(result)
		fmt.Printf("  SettingGet OK: %s\n", string(resultJSON))
		// verify content
		if appName, ok := result["AppName"]; ok && fmt.Sprintf("%v", appName) == "GoAPI" {
			fmt.Println("  Verification: AppName matches!")
		} else {
			fmt.Println("  Verification: AppName mismatch!")
		}
		if version, ok := result["Version"]; ok && fmt.Sprintf("%v", version) == "1.0" {
			fmt.Println("  Verification: Version matches!")
		} else {
			fmt.Println("  Verification: Version mismatch!")
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// Step 3: SETTING_SET - Update the file with new content
	// =============================================
	fmt.Println("--- 39b. SettingSet (update test file) ---")
	updatedContent := `{"AppName":"GoAPI","Version":"2.0","Debug":false,"NewField":"added"}`
	ret = api.SettingSet(section, key, updatedContent)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  SettingSet OK (updated): section=%s key=%s\n", section, key)
		fmt.Printf("  Content: %s\n", updatedContent)
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// Step 4: SETTING_GET - Verify the update
	// =============================================
	fmt.Println("--- 40b. SettingGet (verify update) ---")
	result, ret = api.SettingGet(section, key)
	if ret == mt5webapi.MT_RET_OK && result != nil {
		resultJSON, _ := json.Marshal(result)
		fmt.Printf("  SettingGet OK: %s\n", string(resultJSON))
		// verify updated fields
		if version, ok := result["Version"]; ok && fmt.Sprintf("%v", version) == "2.0" {
			fmt.Println("  Verification: Version updated to 2.0!")
		} else {
			fmt.Println("  Verification: Version update failed!")
		}
		if newField, ok := result["NewField"]; ok && fmt.Sprintf("%v", newField) == "added" {
			fmt.Println("  Verification: NewField present!")
		} else {
			fmt.Println("  Verification: NewField missing!")
		}
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// Step 5: SETTING_DELETE - Delete the test file
	// =============================================
	fmt.Println("--- SettingDelete (delete test file) ---")
	ret = api.SettingDelete(section, key)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  SettingDelete OK: section=%s key=%s\n", section, key)
	} else {
		fmt.Printf("  FAILED: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// Step 6: SETTING_GET - Verify deleted (should return Not found)
	// =============================================
	fmt.Println("--- 40c. SettingGet (verify deleted) ---")
	result, ret = api.SettingGet(section, key)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("  SettingGet returned: [%d] %s (expected - file was deleted)\n", ret, mt5webapi.GetError(ret))
	} else {
		fmt.Println("  WARNING: SettingGet succeeded after delete - file may still exist")
		resultJSON, _ := json.Marshal(result)
		fmt.Printf("  Content: %s\n", string(resultJSON))
	}
	fmt.Println()

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
	fmt.Println("  All settings tests completed!")
	fmt.Println("================================================")
}
