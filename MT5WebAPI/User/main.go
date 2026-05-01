package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

const testLogin = uint64(99999)
const testPassword = "Test1234!"
const testInvestPassword = "Invest1234!"

func main() {
	api := mt5webapi.NewMT5WebAPI("USER", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("========================================")
	fmt.Println("  MT5 Web API - User Example (Go)")
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
	// 8. USER_GET - Kullanıcı bilgisi al
	// =============================================
	fmt.Println("--- 8. UserGet (manager 5005) ---")
	user, ret := api.UserGet(5005)
	if ret == mt5webapi.MT_RET_OK && user != nil {
		fmt.Printf("Login:       %d\n", user.Login)
		fmt.Printf("Name:        %s\n", user.Name)
		fmt.Printf("Group:       %s\n", user.Group)
		fmt.Printf("Country:     %s\n", user.Country)
		fmt.Printf("Leverage:    %d\n", user.Leverage)
		fmt.Printf("Balance:     %.2f\n", user.Balance)
		fmt.Printf("LastIP:      %s\n", user.LastIP)
	} else {
		fmt.Printf("Failed: %s\n", mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 9. USER_ADD - Login 99999 oluştur
	// =============================================
	fmt.Println("--- 9. UserAdd (login 99999) ---")
	newUser := mt5webapi.CreateDefault()
	newUser.Login = testLogin
	newUser.Name = "Go Test User"
	newUser.Group = "demo\\forex"
	newUser.Leverage = 200
	newUser.MainPassword = testPassword
	newUser.InvestPassword = testInvestPassword
	newUser.Email = "gotest@example.com"
	newUser.Country = "Turkey"
	newUser.City = "Istanbul"
	newUser.Phone = "+905551234567"
	newUser.Comment = "Created by Go MT5 Web API"

	createdUser, ret := api.UserAdd(newUser)
	if ret == mt5webapi.MT_RET_OK && createdUser != nil {
		fmt.Printf("User created!\n")
		fmt.Printf("Login:       %d\n", createdUser.Login)
		fmt.Printf("Name:        %s\n", createdUser.Name)
		fmt.Printf("Group:       %s\n", createdUser.Group)
		fmt.Printf("Leverage:    %d\n", createdUser.Leverage)
		fmt.Printf("Email:       %s\n", createdUser.Email)
		fmt.Printf("Country:     %s\n", createdUser.Country)
		fmt.Printf("City:        %s\n", createdUser.City)
	} else {
		fmt.Printf("UserAdd failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
		fmt.Println("Cannot continue tests without test user.")
		api.Disconnect()
		return
	}
	fmt.Println()

	// =============================================
	// 10. USER_UPDATE - Kullanıcı güncelle
	// =============================================
	fmt.Println("--- 10. UserUpdate (99999) ---")
	updateUser, ret := api.UserGet(testLogin)
	if ret == mt5webapi.MT_RET_OK && updateUser != nil {
		oldComment := updateUser.Comment
		oldCity := updateUser.City
		updateUser.Comment = "Updated by Go MT5 Web API"
		updateUser.City = "Ankara"
		updatedUser, ret := api.UserUpdate(updateUser)
		if ret == mt5webapi.MT_RET_OK && updatedUser != nil {
			fmt.Printf("User updated!\n")
			fmt.Printf("Login:       %d\n", updatedUser.Login)
			fmt.Printf("Comment:     %s -> %s\n", oldComment, updatedUser.Comment)
			fmt.Printf("City:        %s -> %s\n", oldCity, updatedUser.City)
		} else {
			fmt.Printf("UserUpdate failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
		}
	}
	fmt.Println()

	// =============================================
	// 12. USER_LOGINS - Grup bazlı login listesi
	// =============================================
	fmt.Println("--- 12. UserLogins (demo\\*) ---")
	logins, ret := api.UserLogins("demo\\*")
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("Found %d logins in group 'demo\\*':\n", len(logins))
		for i, login := range logins {
			fmt.Printf("  [%d] Login: %d\n", i+1, login)
		}
	} else {
		fmt.Printf("UserLogins failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 13. USER_PASSWORD_CHECK - Şifre doğrula
	// =============================================
	fmt.Println("--- 13. UserPasswordCheck (99999) ---")

	// main password - correct
	ret = api.UserPasswordCheckMain(testLogin, testPassword)
	fmt.Printf("Main password (correct):    [%d] %s\n", ret, mt5webapi.GetError(ret))

	// main password - wrong
	ret = api.UserPasswordCheckMain(testLogin, "WrongPassword!")
	fmt.Printf("Main password (wrong):      [%d] %s\n", ret, mt5webapi.GetError(ret))

	// investor password - correct
	ret = api.UserPasswordCheck(testLogin, testInvestPassword, mt5webapi.UserPassInvestor)
	fmt.Printf("Investor password (correct): [%d] %s\n", ret, mt5webapi.GetError(ret))
	fmt.Println()

	// =============================================
	// 14. USER_PASSWORD_CHANGE - Şifre değiştir
	// =============================================
	fmt.Println("--- 14. UserPasswordChange (99999) ---")

	// change main password
	newPass := "NewPass9999!"
	ret = api.UserPasswordChangeMain(testLogin, newPass)
	fmt.Printf("Change main password:       [%d] %s\n", ret, mt5webapi.GetError(ret))

	// verify new password works
	ret = api.UserPasswordCheckMain(testLogin, newPass)
	fmt.Printf("Check new password:         [%d] %s\n", ret, mt5webapi.GetError(ret))

	// verify old password no longer works
	ret = api.UserPasswordCheckMain(testLogin, testPassword)
	fmt.Printf("Check old password:         [%d] %s\n", ret, mt5webapi.GetError(ret))

	// change investor password
	newInvestPass := "NewInvest9999!"
	ret = api.UserPasswordChange(testLogin, newInvestPass, mt5webapi.UserPassInvestor)
	fmt.Printf("Change investor password:   [%d] %s\n", ret, mt5webapi.GetError(ret))
	fmt.Println()

	// =============================================
	// 15. USER_DEPOSIT_CHANGE - Bakiye değiştir
	// =============================================
	fmt.Println("--- 15. UserDepositChange (99999) ---")

	// deposit 10000
	ret = api.UserDepositChange(testLogin, 10000.00, "Initial deposit via Go API", mt5webapi.DealBalance)
	fmt.Printf("Deposit  +10000:  [%d] %s\n", ret, mt5webapi.GetError(ret))

	// check balance
	acc, ret := api.UserAccountGet(testLogin)
	if ret == mt5webapi.MT_RET_OK && acc != nil {
		fmt.Printf("Balance after deposit: %.2f\n", acc.Balance)
	}

	// credit 500
	ret = api.UserDepositChange(testLogin, 500.00, "Credit bonus via Go API", mt5webapi.DealCredit)
	fmt.Printf("Credit   +500:    [%d] %s\n", ret, mt5webapi.GetError(ret))

	// withdraw 2000
	ret = api.UserDepositChange(testLogin, -2000.00, "Withdrawal via Go API", mt5webapi.DealBalance)
	fmt.Printf("Withdraw -2000:   [%d] %s\n", ret, mt5webapi.GetError(ret))

	// final balance check
	acc, ret = api.UserAccountGet(testLogin)
	if ret == mt5webapi.MT_RET_OK && acc != nil {
		fmt.Printf("Final Balance:    %.2f\n", acc.Balance)
		fmt.Printf("Final Credit:     %.2f\n", acc.Credit)
		fmt.Printf("Final Equity:     %.2f\n", acc.Equity)
	}
	fmt.Println()

	// =============================================
	// 16. USER_ACCOUNT_GET - Hesap detayları
	// =============================================
	fmt.Println("--- 16. UserAccountGet (99999) ---")
	account, ret := api.UserAccountGet(testLogin)
	if ret == mt5webapi.MT_RET_OK && account != nil {
		fmt.Printf("Login:          %d\n", account.Login)
		fmt.Printf("Balance:        %.2f\n", account.Balance)
		fmt.Printf("Credit:         %.2f\n", account.Credit)
		fmt.Printf("Equity:         %.2f\n", account.Equity)
		fmt.Printf("Margin:         %.2f\n", account.Margin)
		fmt.Printf("MarginFree:     %.2f\n", account.MarginFree)
		fmt.Printf("MarginLevel:    %.2f%%\n", account.MarginLevel)
		fmt.Printf("MarginLeverage: %d\n", account.MarginLeverage)
		fmt.Printf("Profit:         %.2f\n", account.Profit)
		fmt.Printf("Floating:       %.2f\n", account.Floating)
	} else {
		fmt.Printf("UserAccountGet failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// =============================================
	// 11. USER_DELETE - Kullanıcı sil
	// =============================================
	fmt.Println("--- 11. UserDelete (99999) ---")
	ret = api.UserDelete(testLogin)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("User %d deleted successfully!\n", testLogin)
	} else {
		fmt.Printf("UserDelete failed: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}

	// verify deleted
	_, ret = api.UserGet(testLogin)
	fmt.Printf("UserGet after delete: [%d] %s\n", ret, mt5webapi.GetError(ret))
	fmt.Println()

	// disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("========================================")
	fmt.Println("  All user tests completed!")
	fmt.Println("========================================")
}
