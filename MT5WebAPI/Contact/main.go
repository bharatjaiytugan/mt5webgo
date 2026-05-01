package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("CONTACT", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - Contact Example (Go)")
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

	// create test user for mail
	testLogin := uint64(99999)
	fmt.Printf("Creating test user %d...\n", testLogin)
	newUser := mt5webapi.CreateDefault()
	newUser.Login = testLogin
	newUser.Name = "Contact Test User"
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

	// =============================================
	// 44. MAIL_SEND - Kullanıcıya mesaj gönder
	// =============================================
	fmt.Println("--- 44. MailSend ---")

	// send to specific login
	fmt.Print("1) Mail to login 99999... ")
	ret = api.MailSend("99999", "Test Message", "Hello from Go MT5 Web API!")
	fmt.Printf("[%d] %s\n", ret, mt5webapi.GetError(ret))

	// send with HTML content
	fmt.Print("2) Mail with HTML content... ")
	htmlBody := "<h1>Go MT5 Web API</h1><p>This is a <b>test</b> message with HTML content.</p><ul><li>Item 1</li><li>Item 2</li></ul>"
	ret = api.MailSend("99999", "HTML Test Message", htmlBody)
	fmt.Printf("[%d] %s\n", ret, mt5webapi.GetError(ret))

	// send to group mask
	fmt.Print("3) Mail to group demo\\*... ")
	ret = api.MailSend("demo\\*", "Broadcast Message", "This message is sent to all demo users via Go API.")
	fmt.Printf("[%d] %s\n", ret, mt5webapi.GetError(ret))
	fmt.Println()

	// =============================================
	// 45. NEWS_SEND - Haber yayınla
	// =============================================
	fmt.Println("--- 45. NewsSend ---")

	// send news
	fmt.Print("1) News (English, normal priority)... ")
	ret = api.NewsSend(
		"Go API Test News",          // subject
		"General",                   // category
		9,                           // language (9 = English LANGID)
		0,                           // priority (0 = normal)
		"This is a test news article published from Go MT5 Web API.",
	)
	fmt.Printf("[%d] %s\n", ret, mt5webapi.GetError(ret))

	// send news with HTML
	fmt.Print("2) News with HTML content... ")
	newsHTML := "<h2>Market Update</h2><p>EURUSD reached <b>1.1678</b> today.</p><p>Published via Go MT5 Web API.</p>"
	ret = api.NewsSend(
		"Market Update - Go API",
		"Market",
		9,
		1, // priority 1 = high
		newsHTML,
	)
	fmt.Printf("[%d] %s\n", ret, mt5webapi.GetError(ret))

	// send news in Turkish
	fmt.Print("3) News (Turkish)... ")
	ret = api.NewsSend(
		"Go API Test Haberi",
		"Genel",
		31, // language (31 = Turkish LANGID)
		0,
		"Bu haber Go MT5 Web API ile yayinlanmistir.",
	)
	fmt.Printf("[%d] %s\n", ret, mt5webapi.GetError(ret))
	fmt.Println()

	// cleanup
	fmt.Print("Cleaning up... ")
	api.UserDelete(testLogin)
	fmt.Printf("User %d deleted.\n", testLogin)

	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All contact tests completed!")
	fmt.Println("========================================")
}
