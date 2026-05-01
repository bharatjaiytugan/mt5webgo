package main

import (
	"config"
	"fmt"

	"mt5webapi"
)

func main() {
	api := mt5webapi.NewMT5WebAPI("CLIENT_CRUD", func(logType int, message string) {
		if logType == 1 {
			fmt.Printf("[ERROR] %s\n", message)
		}
	})

	fmt.Println("================================================")
	fmt.Println("  MT5 Web API - Client CRUD Example (Go)")
	fmt.Println("  Tests: 33, 34")
	fmt.Println("================================================")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		return
	}

	// 1. Connect
	fmt.Print("Connecting... ")
	ret := api.ConnectDefault(cfg.Server, cfg.Port, cfg.Login, cfg.Password, mt5webapi.PUMP_MODE_NONE)
	if ret != mt5webapi.MT_RET_OK {
		fmt.Printf("FAILED! [%d] %s\n", ret, mt5webapi.GetError(ret))
		return
	}
	fmt.Println("OK")
	fmt.Println()

	var clientID uint64

	// 2. ClientAdd - create a new client
	fmt.Println("--- 33. ClientAdd ---")
	jw := mt5webapi.NewJSONWriter()
	jw.WriteBeginObject()
	jw.WriteAttributeString("ClientType", "0")
	jw.WriteAttributeString("ClientStatus", "0")
	jw.WriteAttributeString("PersonName", "GoAPI")
	jw.WriteAttributeString("PersonLastName", "TestClient")
	jw.WriteAttributeString("ContactEmail", "goapi@test.com")
	jw.WriteAttributeString("AddressCountry", "Turkey")
	jw.WriteAttributeString("Comment", "Created by Go API")
	jw.WriteAttributeString("AssignedManager", fmt.Sprintf("%d", cfg.Login))
	jw.WriteEndObject()
	clientJSON := jw.String()
	fmt.Printf("  Request JSON: %s\n", clientJSON)


	client, ret := api.ClientAdd(clientJSON)
	if ret == mt5webapi.MT_RET_OK {
		if client != nil {
			clientID = client.RecordID
			fmt.Printf("  OK - RecordID=%d Name=%s\n", client.RecordID, client.Name)
		} else {
			fmt.Println("  OK - client created (no RecordID in response)")
			fmt.Println("  Finding client via ClientGetIDs...")
			ids, r := api.ClientGetIDs("")
			if r == mt5webapi.MT_RET_OK && len(ids) > 0 {
				clientID = ids[len(ids)-1]
				fmt.Printf("  All client IDs: %v\n", ids)
				fmt.Printf("  Using last ID: %d\n", clientID)
			} else {
				fmt.Printf("  ClientGetIDs failed: [%d] %s\n", r, mt5webapi.GetError(r))
			}
		}
	} else {
		fmt.Printf("  FAILED! [%d] %s\n", ret, mt5webapi.GetError(ret))
		api.Disconnect()
		return
	}
	fmt.Println()

	// 3. ClientGet - verify client was created
	fmt.Println("--- 34a. ClientGet ---")
	client, ret = api.ClientGet(clientID)
	if ret == mt5webapi.MT_RET_OK && client != nil {
		printClient(client)
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 4. ClientUpdate - update comment
	fmt.Println("--- 34b. ClientUpdate ---")
	jwUpd := mt5webapi.NewJSONWriter()
	jwUpd.WriteBeginObject()
	jwUpd.WriteAttributeUint64("RecordID", clientID)
	jwUpd.WriteAttributeString("Comment", "Updated by Go API")
	jwUpd.WriteEndObject()
	updateJSON := jwUpd.String()
	fmt.Printf("  Request JSON: %s\n", updateJSON)

	client, ret = api.ClientUpdate(updateJSON)
	if ret == mt5webapi.MT_RET_OK && client != nil {
		fmt.Printf("  OK - RecordID=%d Comment=%s\n", client.RecordID, client.Comment)
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 5. ClientUserAdd - link login 5005 to client
	fmt.Println("--- 34c. ClientUserAdd ---")
	ret = api.ClientUserAdd(clientID, 5005)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  OK - linked login 5005 to client %d\n", clientID)
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 6. ClientUserLogins - list users linked to client
	fmt.Println("--- 34d. ClientUserLogins ---")
	logins, ret := api.ClientUserLogins(clientID)
	if ret == mt5webapi.MT_RET_OK && logins != nil {
		fmt.Printf("  Found %d linked logins:", len(logins))
		for _, l := range logins {
			fmt.Printf(" %d", l)
		}
		fmt.Println()
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 7. ClientGetIDs - get all client IDs
	fmt.Println("--- 34e. ClientGetIDs ---")
	ids, ret := api.ClientGetIDs("")
	if ret == mt5webapi.MT_RET_OK && ids != nil {
		fmt.Printf("  Found %d client IDs:", len(ids))
		for _, id := range ids {
			fmt.Printf(" %d", id)
		}
		fmt.Println()
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 8. ClientUserDelete - unlink login 5005
	fmt.Println("--- 34f. ClientUserDelete ---")
	ret = api.ClientUserDelete(clientID, 5005)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  OK - unlinked login 5005 from client %d\n", clientID)
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 9. ClientDelete - delete client
	fmt.Println("--- 34g. ClientDelete ---")
	ret = api.ClientDelete(clientID)
	if ret == mt5webapi.MT_RET_OK {
		fmt.Printf("  OK - deleted client %d\n", clientID)
	} else {
		fmt.Printf("  Result: [%d] %s\n", ret, mt5webapi.GetError(ret))
	}
	fmt.Println()

	// 10. Disconnect
	api.Disconnect()
	fmt.Println("Disconnected.")
	fmt.Println("================================================")
}

func printClient(c *mt5webapi.MTClient) {
	fmt.Printf("  RecordID=%-8d Name=%-20s Email=%-25s Country=%-10s Comment=%s\n",
		c.RecordID, c.Name, c.Email, c.Country, c.Comment)
}
