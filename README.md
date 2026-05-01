# MT5 Web API - Go Implementation

A comprehensive Go implementation of the MetaTrader 5 Web API. Fully ported from the original C# (.NET Framework 4.7.2) library and extended with all 284 REST API endpoints.

## Features

- **295 public API methods** covering the entire MT5 Manager Web API surface
- **WebSocket binary protocol** with AES-256 encrypted TCP communication
- **Zero external dependencies** - only Go standard library
- **Real-time price streaming** via WebSocket with configurable polling interval (down to 1ms)
- **Dealer trading** - send Market, Limit, Stop orders; modify and cancel
- **OHLC bar builder** - automatically construct bars from streaming tick data
- **Single native binary** - no DLL, no runtime, just compile and run

## Project Structure

```
Go/
├── mt5webapi/                    # Core library (57 files, 17,436 lines)
│   ├── mt5webapi.go              # Facade pattern - 295 public methods
│   ├── retcode.go                # Error codes (100+ MT5 return codes)
│   ├── format.go                 # Human-readable error messages
│   │
│   ├── # Data Models
│   ├── user.go                   # MTUser, EnUsersRights, EnUsersPasswords
│   ├── account.go                # MTAccount, EnSoActivation
│   ├── order.go                  # MTOrder, EnOrderType, EnOrderFilling
│   ├── deal.go                   # MTDeal, EnDealAction, EnEntryFlags
│   ├── position.go               # MTPosition, EnPositionAction
│   ├── consymbol.go              # MTConSymbol (1,174 lines, 100+ enums)
│   ├── congroup.go               # MTConGroup, trade flags, margin modes
│   ├── congroupsymbol.go         # MTConGroupSymbol
│   ├── concommon.go              # MTConCommon, server info
│   ├── concommission.go          # MTConCommission, MTConCommTier
│   ├── contime.go                # MTConTime, working hours
│   │
│   ├── # Protocol Layer
│   ├── protocol.go               # Constants, header parser, param parser
│   ├── connect.go                # Synchronous TCP connection
│   ├── async_connect.go          # Async WebSocket connection (goroutines)
│   ├── async_send.go             # Async send/receive with channels
│   ├── auth.go                   # Challenge-response authentication
│   ├── crypt.go                  # AES-256-OFB encryption/decryption
│   ├── crypt_aes.go              # Full AES implementation (S-box, Galois)
│   ├── api_base.go               # Base API class, JSON parsing helpers
│   │
│   ├── # Protocol Handlers (29 files)
│   ├── user_handler.go           # User CRUD operations
│   ├── order_handler.go          # Order queries
│   ├── deal_handler.go           # Deal queries
│   ├── position_handler.go       # Position queries
│   ├── group_handler.go          # Group management
│   ├── symbol_handler.go         # Symbol management
│   ├── time_handler.go           # Server time operations
│   ├── common_handler.go         # Common server configuration
│   ├── history_handler.go        # Order history
│   ├── tick_handler.go           # Tick/price data
│   ├── mail_handler.go           # Internal mail
│   ├── news_handler.go           # News publishing
│   ├── ping_handler.go           # Keepalive ping
│   ├── custom_handler.go         # Custom commands
│   ├── trade_handler.go          # Trade balance operations
│   ├── server_handler.go         # Server restart
│   ├── dealer_handler.go         # Dealer trading (order execution)
│   ├── server_config_handler.go  # Server configuration CRUD
│   ├── batch_handler.go          # Batch operations
│   ├── backup_handler.go         # Backup and recovery
│   ├── trade_calc_handler.go     # Rate/margin/profit calculations
│   ├── client_handler.go         # Client CRM operations
│   ├── chart_handler.go          # Chart and market depth data
│   ├── daily_handler.go          # Daily reports
│   ├── settings_handler.go       # Settings file management
│   ├── subscription_handler.go   # Subscription management
│   ├── config_crud_handler.go    # Config CRUD (8 entity types)
│   ├── infra_handler.go          # Infrastructure (6 entity types)
│   ├── misc_handler.go           # Miscellaneous operations (48 methods)
│   │
│   ├── # Real-Time Streaming
│   ├── streaming.go              # PriceStreamer, BarBuilder
│   │
│   ├── # Utilities
│   ├── log.go                    # Thread-safe logging
│   ├── utils.go                  # MD5, hex, UTF-16LE, buffer operations
│   ├── socket.go                 # TCP socket wrapper
│   ├── json_writer.go            # JSON serializer (matches C# JSONWriter)
│   └── handler_constructors.go   # Handler factory functions
│
├── MT5WebAPI/                    # Example applications
│   ├── .env                      # Connection credentials (add to .gitignore)
│   ├── .env.example              # Credential template
│   ├── config/                   # .env file reader package
│   │
│   ├── # Core API Examples
│   ├── Login/                    # Connect, Disconnect, Ping
│   ├── Server/                   # TimeServer, TimeGet, CommonGet
│   ├── User/                     # Full user lifecycle test
│   ├── Order/                    # Order queries
│   ├── Position/                 # Position queries
│   ├── Deal/                     # Deal queries
│   ├── History/                  # Order history
│   ├── Tick/                     # TickLast, TickLastGroup, TickStat
│   ├── Group/                    # Group CRUD
│   ├── Symbol/                   # Symbol CRUD
│   ├── Trade/                    # TradeBalance (8 operation types)
│   ├── Contact/                  # MailSend, NewsSend
│   ├── Custom/                   # CustomSend
│   │
│   └── RESTAPI/                  # Extended API Examples
│       ├── Trading/              # All 6 order types + modify + cancel
│       ├── Server/               # CommonSet, AccessServer, TLS, TimeSet
│       ├── Batch/                # Group/Symbol/Order/History batch ops
│       ├── BackupRecovery/       # Backup list/get/restore, OrderReopen
│       ├── Deal/                 # DealGetBatch, DealUpdate
│       ├── Position/             # PositionGetBatch
│       ├── Trade/                # CalcRateBuy/Sell, CheckMargin, CalcProfit
│       ├── User/                 # UserGetBatch, UserBackup
│       ├── Client/               # Client CRUD + user account binding
│       ├── Price/                # Live price streaming + BarBuilder demo
│       ├── Daily/                # Daily reports
│       ├── Settings/             # Settings file CRUD
│       ├── Subscription/         # Subscription configuration
│       ├── ConfigCRUD/           # Leverage/Firewall/Holiday/Manager/Route
│       ├── Infrastructure/       # Email/Messenger/Gateway/Feeder/Report/Plugin
│       └── Misc/                 # Logger, UserTotal, TestAccess, Backups
```

## Quick Start

### 1. Requirements

- Go 1.21 or later
- MetaTrader 5 trade server with Web API access
- Manager account with appropriate permissions

### 2. Installation

```bash
git clone <repo-url>
cd Go/mt5webapi
go build ./...
```

### 3. Configuration

Create a `.env` file in the `MT5WebAPI/` directory:

```env
MT5_SERVER=your_server_ip
MT5_PORT=443
MT5_LOGIN=your_manager_login
MT5_PASSWORD=your_password
```

### 4. Basic Usage

```go
package main

import (
    "fmt"
    "mt5webapi"
)

func main() {
    // Create API instance
    api := mt5webapi.NewMT5WebAPIDefault()

    // Connect to MT5 server
    ret := api.ConnectDefault("10.0.0.1", 443, 1000, "password", mt5webapi.PUMP_MODE_NONE)
    if ret != mt5webapi.MT_RET_OK {
        fmt.Printf("Connection failed: %s\n", mt5webapi.GetError(ret))
        return
    }
    defer api.Disconnect()

    // Get server time
    serverTime := api.TimeServer()
    fmt.Printf("Server time: %d\n", serverTime)

    // Get user info
    user, ret := api.UserGet(1000)
    if ret == mt5webapi.MT_RET_OK {
        fmt.Printf("Login: %d, Name: %s, Group: %s, Balance: %.2f\n",
            user.Login, user.Name, user.Group, user.Balance)
    }

    // Get live prices
    ticks, ret := api.TickLast("EURUSD")
    if ret == mt5webapi.MT_RET_OK && len(ticks) > 0 {
        fmt.Printf("EURUSD  Bid=%.5f  Ask=%.5f\n", ticks[0].Bid, ticks[0].Ask)
    }
}
```

## API Reference

### Connection

```go
// Create with logging callback
api := mt5webapi.NewMT5WebAPI("MyAgent", func(logType int, msg string) {
    fmt.Printf("[%d] %s\n", logType, msg)
})

// Create with defaults
api := mt5webapi.NewMT5WebAPIDefault()

// Connect
ret := api.Connect(server, port, login, password, pumpModes, crypt, timeout)
ret := api.ConnectDefault(server, port, login, password, pumpModes)

// Connection management
api.Disconnect()
api.IsConnected()  // bool
api.Ping()         // keepalive
```

### User Management

```go
// CRUD
user, ret := api.UserGet(login)
newUser, ret := api.UserAdd(user)
updatedUser, ret := api.UserUpdate(user)
ret := api.UserDelete(login)

// Queries
logins, ret := api.UserLogins("demo\\*")
users, ret := api.UserGetBatch("1000,1001", "")
account, ret := api.UserAccountGet(login)
total, ret := api.UserTotal()

// Authentication
ret := api.UserPasswordCheck(login, password, mt5webapi.UserPassMain)
ret := api.UserPasswordChange(login, newPassword, mt5webapi.UserPassMain)

// Balance
ret := api.UserDepositChange(login, 1000.0, "Deposit", mt5webapi.DealBalance)
```

### Dealer Trading (Order Execution)

```go
// Market orders
result, ret := api.DealerBuyMarket(login, "EURUSD", 100, "Buy")   // 0.01 lot
result, ret := api.DealerSellMarket(login, "EURUSD", 100, "Sell")

// Pending orders (Buy Limit, Sell Limit, Buy Stop, Sell Stop)
result, ret := api.DealerPendingOrder(login, "EURUSD",
    mt5webapi.OpBuyLimit,  // order type
    100,                   // volume: 0.01 lot
    1.15000,               // price
    1.14000,               // stop loss
    1.16000,               // take profit
    "Buy Limit")

// Modify pending order
result, ret := api.DealerModifyOrder(login, orderTicket, newPrice, sl, tp)

// Cancel pending order
result, ret := api.DealerCancelOrder(login, orderTicket)

// Balance operations
ret := api.TradeBalance(login, mt5webapi.DealBalance, 10000.0, "Deposit")
ticket, ret := api.TradeBalanceWithTicket(login, mt5webapi.DealCredit, 500.0, "Bonus")
```

**Volume format:** Values are in 1/10,000 of a lot:
- `100` = 0.01 lot
- `1000` = 0.1 lot
- `10000` = 1.0 lot
- `100000` = 10.0 lots

**Dealer action codes:**

| Code | Constant | Description |
|------|----------|-------------|
| 200 | `TA_DEALER_POS_EXECUTE` | Execute market order (buy/sell) |
| 201 | `TA_DEALER_ORD_PENDING` | Place pending order |
| 202 | `TA_DEALER_POS_MODIFY` | Modify position |
| 203 | `TA_DEALER_ORD_MODIFY` | Modify pending order |
| 204 | `TA_DEALER_ORD_REMOVE` | Remove pending order |
| 205 | `TA_DEALER_ORD_ACTIVATE` | Activate pending order |
| 206 | `TA_DEALER_BALANCE` | Balance operation |

### Real-Time Price Streaming

```go
// Create streamer with callback
streamer := api.NewPriceStreamer(func(update *mt5webapi.PriceUpdate) {
    fmt.Printf("[%s] %-8s Bid=%.5f Ask=%.5f Spread=%.5f\n",
        time.Now().Format("15:04:05.000"),
        update.Symbol, update.Bid, update.Ask, update.Spread)
})

// Add symbols to watch
streamer.AddSymbols("EURUSD", "GBPUSD", "USDJPY", "XAUUSD")

// Set polling interval (default: 200ms)
streamer.SetInterval(1 * time.Millisecond)

// Start streaming
streamer.Start()

// Get last cached price
price := streamer.GetLastPrice("EURUSD")
fmt.Printf("Last EURUSD: Bid=%.5f\n", price.Bid)

// Build OHLC bars from ticks
barBuilder := mt5webapi.NewBarBuilder("EURUSD", 1*time.Minute, func(bar *mt5webapi.MTBar) {
    fmt.Printf("[BAR] %s O=%.5f H=%.5f L=%.5f C=%.5f Ticks=%d\n",
        time.Unix(bar.Datetime, 0).Format("15:04"),
        bar.Open, bar.High, bar.Low, bar.Close, bar.TickVol)
})

// Feed ticks to bar builder in streamer callback
streamer = api.NewPriceStreamer(func(update *mt5webapi.PriceUpdate) {
    if update.Symbol == "EURUSD" {
        barBuilder.OnTick(update.Bid, time.Now())
    }
})
streamer.AddSymbol("EURUSD")
streamer.Start()

// Stop when done
streamer.Stop()
```

### Order / Position / Deal Queries

```go
// Orders (active pending orders)
order, ret := api.OrderGet(ticket)
total, ret := api.OrderGetTotal(login)
orders, ret := api.OrderGetPage(login, offset, count)
orders, ret := api.OrderGetBatch("1000,1001", "", "")  // by logins
ret := api.OrderCancel("12345,12346")                   // move to history
ret := api.OrderDeleteDirect("12345")                    // permanently delete
order, ret := api.OrderUpdateDirect(orderJSON)           // modify directly

// Positions (open trades)
position, ret := api.PositionGet(login, "EURUSD")
total, ret := api.PositionGetTotal(login)
positions, ret := api.PositionGetPage(login, offset, count)
positions, ret := api.PositionGetBatch("1000", "demo\\*", "")

// Deals (executed transactions)
deal, ret := api.DealGet(ticket)
total, ret := api.DealGetTotal(login, from, to)
deals, ret := api.DealGetPage(login, from, to, offset, count)
deals, ret := api.DealGetBatch("1000", "", "", from, to)

// History (closed orders)
history, ret := api.HistoryGet(ticket)
total, ret := api.HistoryGetTotal(login, from, to)
orders, ret := api.HistoryGetPage(login, from, to, offset, count)
```

### Trade Calculations

```go
// Currency conversion rates
buyRate, ret := api.CalcRateBuy("EUR", "USD")   // ~1.17
sellRate, ret := api.CalcRateSell("GBP", "USD") // ~1.36

// Margin check (what-if analysis)
result, ret := api.CheckMargin(login, "EURUSD", 0, 10000, askPrice)
// result["new"] = account state after order
// result["current"] = current account state

// Profit calculation
profit, ret := api.CalcProfit("demo\\forex", "EURUSD", 0, 10000, 1.16000, 1.17000)
```

### Group & Symbol Management

```go
// Groups
total, ret := api.GroupTotal()               // count
group, ret := api.GroupGet("demo\\forex")    // by name
groups, ret := api.GroupList("demo\\*")      // by mask
newGroup, ret := api.GroupAdd(group)          // create (clone existing)
ret := api.GroupDelete("demo\\test")         // delete

// Symbols
total, ret := api.SymbolTotal()              // count: 2183
symbol, ret := api.SymbolGet("EURUSD")       // by name
newSym, ret := api.SymbolAdd(symbol)         // create
ret := api.SymbolDelete("GOTEST")            // delete

// Batch operations
ret := api.GroupAddBatch(groupsJSON)
ret := api.SymbolAddBatch(symbolsJSON)
ret := api.GroupDeleteBatch("group1,group2")
ret := api.SymbolDeleteBatch("SYM1,SYM2")
```

### Client (CRM) Management

```go
// Client CRUD (request uses JSON array format)
client, ret := api.ClientAdd(`{"PersonName":"John","PersonLastName":"Doe","ClientType":"0"}`)
client, ret := api.ClientGet(recordID)
client, ret := api.ClientUpdate(`{"RecordID":1001,"Comment":"Updated"}`)
ret := api.ClientDelete(recordID)

// Bind trading accounts to clients
ret := api.ClientUserAdd(clientID, login)
ret := api.ClientUserDelete(clientID, login)
logins, ret := api.ClientUserLogins(clientID)
clientIDs, ret := api.ClientGetIDs("")  // all clients visible to manager
```

### Settings Files

```go
// Read
settings, ret := api.SettingGet("web", "config.json")

// Write (creates or overwrites)
ret := api.SettingSet("web", "config.json", `{"theme":"dark","lang":"en"}`)

// Delete
ret := api.SettingDelete("web", "config.json")
```

### Server Configuration

```go
// Common settings
common, ret := api.CommonGet()
updated, ret := api.CommonSet(common)

// Access servers
total, ret := api.AccessServerTotal()  // 4
server, ret := api.AccessServerNext(0) // Main Trade Server

// TLS certificates
total, ret := api.TLSCertTotal()

// Time configuration
conTime, ret := api.TimeGet()
updated, ret := api.TimeSet(conTime)
```

### Infrastructure Management

Each infrastructure entity follows the same CRUD pattern: `Total`, `Next`, `Get`, `GetAll`, `Add`, `Delete`, `Shift`.

```go
// Configuration entities
total, _ := api.LeverageTotal()     // Floating margin rules
total, _ := api.FirewallTotal()     // IP filtering rules
total, _ := api.HolidayTotal()      // Holiday schedules
total, _ := api.ManagerTotal()      // Manager accounts
total, _ := api.RouteTotal()        // Order routing rules
total, _ := api.SpreadTotal()       // Spread configurations
total, _ := api.HistorySyncTotal()  // History sync settings
total, _ := api.SymbolGroupTotal()  // Symbol groups

// Infrastructure with module support
total, _ := api.GatewayTotal()       // Gateways
modTotal, _ := api.GatewayModuleTotal()
total, _ := api.FeederTotal()        // Data feeds
total, _ := api.ReportTotal()        // Report modules
total, _ := api.PluginTotal()        // Server plugins
total, _ := api.EmailTotal()         // Email configurations
total, _ := api.MessengerTotal()     // Messenger configurations
```

### Backup & Recovery

```go
// List available backups
orderDates, ret := api.OrderBackupList(from, to)
positionDates, ret := api.PositionBackupList(from, to)
dealDates, ret := api.DealBackupList(from, to)

// Get data from backup
orders, ret := api.OrderBackupGet(backupDate, login)

// Restore from backup
restored, ret := api.OrderBackupRestore(orderJSON)

// Reopen a canceled pending order from history
reopened, ret := api.OrderReopen(ticket)
```

### Miscellaneous

```go
// Server journal logs
logs, ret := api.LoggerServerRequest("0", "0", from, to, "searchterm")

// User statistics
total, ret := api.UserTotal()

// Connection test
ret := api.TestAccess()

// Internal mail
ret := api.MailSend("1000", "Subject", "<b>HTML body</b>")

// News broadcasting
ret := api.NewsSend("Subject", "Category", 9, 0, "News body")

// Push notifications
ret := api.NotificationSend(logins, "Notification text")

// Raw custom command (any WebSocket command)
answer, ret := api.CustomSend("COMMAND_NAME", params, body)
```

## Architecture

### Protocol Flow

```
Your Go Application
         │
         ▼
MT5WebAPI (Facade - 295 methods)
         │
         ▼
Handler Layer (user_handler, dealer_handler, ...)
         │
         ▼
apiBase.send(command, params)
         │
         ▼
mtAsyncSend ──► Async Queue (channel-based)
         │
         ▼
mtAsyncConnect (3 goroutines)
  ├── sendThread:    queue → encrypt → TCP write
  ├── recvThread:    TCP read → decrypt → callback
  └── clearThread:   timeout cleanup
         │
         ▼
AES-256-OFB Encryption ↔ UTF-16LE Encoding ↔ TCP Socket
         │
         ▼
MetaTrader 5 Server (port 443)
```

### C# to Go Conversion Map

| C# (.NET) | Go |
|-----------|-----|
| `namespace` | `package` |
| `class` | `struct` |
| `property { get; set; }` | exported field |
| `delegate` | `func` type |
| `ManualResetEvent` | `chan struct{}` |
| `Thread` | goroutine |
| `ReaderWriterLockSlim` | `sync.RWMutex` |
| `Interlocked` | `sync/atomic` |
| `out` parameter | multiple return values |
| `Encoding.Unicode` | `unicode/utf16` + `encoding/binary` |
| `JavaScriptSerializer` | `encoding/json` + manual parsing |
| `.NET Framework 4.7.2 DLL` | Single native binary |
| `using` / IDisposable | `defer` |

### Key Design Decisions

1. **Single package** - All code in one `mt5webapi` package, avoiding circular imports that would occur with sub-packages
2. **Facade pattern** - `MT5WebAPI` struct provides a clean public API; handlers are internal
3. **Lazy initialization** - Handlers are created on first use, not at connection time
4. **Tolerant parsing** - Server may return command names in different formats (underscore vs slash); parsers accept any format
5. **UTF-16LE everywhere** - MT5 protocol uses UTF-16LE encoding; all string conversion handled transparently

## Test Results

All tests verified against a live MT5 trade server:

| Category | Operations Tested | Status |
|----------|-------------------|--------|
| **Connection** | Connect, Disconnect, Ping | ✅ |
| **Server** | TimeServer, TimeGet, CommonGet/Set, AccessServer, TLS | ✅ |
| **Users** | Add, Update, Delete, Get, Logins, Password, Deposit, Account, Batch | ✅ |
| **Orders** | Get, GetTotal, GetPage, GetBatch, Update, Delete, Cancel, Reopen | ✅ |
| **Positions** | Get, GetTotal, GetPage, GetBatch | ✅ |
| **Deals** | Get, GetTotal, GetPage, GetBatch, Update | ✅ |
| **History** | Get, GetTotal, GetPage, GetBatch, Update | ✅ |
| **Ticks** | Last, LastGroup, Stat | ✅ |
| **Groups** | Total, Next, Get, Add, Delete, List, AddBatch, DeleteBatch | ✅ |
| **Symbols** | Total, Next, Get, GetGroup, Add, Delete, AddBatch, DeleteBatch | ✅ |
| **Trading** | Market Buy/Sell, Buy/Sell Limit/Stop, Modify, Cancel | ✅ |
| **Calculations** | CalcRateBuy/Sell, CheckMargin, CalcProfit | ✅ |
| **Clients** | Add, Get, Update, Delete, UserAdd, UserDelete, UserLogins, GetIDs | ✅ |
| **Streaming** | Live price feed (1ms interval), OHLC BarBuilder | ✅ |
| **Daily** | DailyGet, DailyGetBatch | ✅ |
| **Settings** | Get, Set, Delete (full CRUD cycle verified) | ✅ |
| **Subscriptions** | CfgTotal (138 configs), CfgNext, CfgGet | ✅ |
| **Backup** | OrderBackupList, OrderBackupGet, OrderReopen | ✅ |
| **Config** | Leverage, Firewall, Holiday, Manager, Route, Spread, HistorySync, SymbolGroup | ✅ |
| **Infrastructure** | Email, Messenger, Gateway, Feeder (3), Report (84), Plugin | ✅ |
| **Misc** | Logger (1876 entries), UserTotal (17), MailGet, TestAccess, BackupList | ✅ |

## Important Notes

### WebSocket vs REST API

All communication uses the MT5 **WebSocket binary protocol** over TCP. No HTTP REST API calls are made. A few commands (`CHART_GET`, `BOOK_GET`, `NEWS_GET`) are only available via HTTPS REST API and are not supported on the binary WebSocket protocol. The library includes handler code for these, ready for use when an HTTPS transport is added.

### Volume Format

All volume values use the 1/10,000 lot format:

| Value | Lot Size |
|-------|----------|
| `100` | 0.01 lot |
| `1,000` | 0.1 lot |
| `10,000` | 1.0 lot |
| `100,000` | 10.0 lots |

### JSON Parsing

The MT5 server sends all numeric values as strings in JSON responses (e.g., `"Balance":"10000.00"`). The library automatically converts these to appropriate Go types (`uint64`, `float64`, etc.) using manual parsing rather than `json.Unmarshal` to handle this quirk.

### Thread Safety

- `PriceStreamer` is fully thread-safe (mutex-protected symbol list and price cache)
- `MT5WebAPI` instance should be used from a single goroutine for command operations
- The async connection layer runs 3 background goroutines (send, receive, cleanup)
- `Disconnect()` calls `Stop()` which properly signals all goroutines to exit

### Error Handling

All API methods return `MTRetCode`. Check against `MT_RET_OK` (0) for success:

```go
user, ret := api.UserGet(login)
if ret != mt5webapi.MT_RET_OK {
    fmt.Printf("Error: [%d] %s\n", ret, mt5webapi.GetError(ret))
    return
}
```

## Building Examples

```bash
# Build all examples
cd MT5WebAPI/Login && go build -o Login.exe .
cd MT5WebAPI/RESTAPI/Trading && go build -o Trading.exe .

# Run an example
./Login.exe
./Trading.exe
```

## License

This project is built upon the MetaTrader 5 Web API SDK provided by MetaQuotes Ltd. The original C# SDK license terms apply.

## Community

Join our Telegram group for discussions, support, and updates:

[![Telegram](https://img.shields.io/badge/Telegram-Join%20Group-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/+pfMRu5HgdFk3YjM0)

**https://t.me/+pfMRu5HgdFk3YjM0**

## Support & Contributing

For questions, bug reports, feature requests, and commercial licensing inquiries, contact directly:

[![Telegram](https://img.shields.io/badge/Telegram-@iamocap-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/iamocap)

**https://t.me/iamocap**
