//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"sync"
	"time"
)

// PriceUpdate represents a real-time price update.
type PriceUpdate struct {
	Symbol     string
	Bid        float64
	Ask        float64
	Last       float64
	Volume     uint64
	VolumeReal float64
	Datetime   uint64
	Spread     float64 // Ask - Bid
}

// PriceCallback is called when a new price update is received.
type PriceCallback func(update *PriceUpdate)

// PriceStreamer manages real-time price streaming for multiple symbols
// by polling TickLast at a configurable interval.
type PriceStreamer struct {
	api        *MT5WebAPI
	symbols    []string
	interval   time.Duration
	callback   PriceCallback
	stopCh     chan struct{}
	running    bool
	mu         sync.Mutex
	lastPrices map[string]*PriceUpdate
}

// NewPriceStreamer creates a new price streamer with a default 200ms polling interval.
func NewPriceStreamer(api *MT5WebAPI, callback PriceCallback) *PriceStreamer {
	return &PriceStreamer{
		api:        api,
		interval:   200 * time.Millisecond,
		callback:   callback,
		lastPrices: make(map[string]*PriceUpdate),
	}
}

// NewPriceStreamer creates a new PriceStreamer bound to this API connection.
func (api *MT5WebAPI) NewPriceStreamer(callback PriceCallback) *PriceStreamer {
	return NewPriceStreamer(api, callback)
}

// AddSymbol adds a symbol to the streaming list.
func (ps *PriceStreamer) AddSymbol(symbol string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for _, s := range ps.symbols {
		if s == symbol {
			return
		}
	}
	ps.symbols = append(ps.symbols, symbol)
}

// AddSymbols adds multiple symbols to the streaming list.
func (ps *PriceStreamer) AddSymbols(symbols ...string) {
	for _, s := range symbols {
		ps.AddSymbol(s)
	}
}

// RemoveSymbol removes a symbol from the streaming list.
func (ps *PriceStreamer) RemoveSymbol(symbol string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for i, s := range ps.symbols {
		if s == symbol {
			ps.symbols = append(ps.symbols[:i], ps.symbols[i+1:]...)
			delete(ps.lastPrices, symbol)
			return
		}
	}
}

// SetInterval sets the polling interval. Default is 200ms.
func (ps *PriceStreamer) SetInterval(d time.Duration) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.interval = d
}

// Start begins streaming prices in a background goroutine.
// If already running, this is a no-op.
func (ps *PriceStreamer) Start() {
	ps.mu.Lock()
	if ps.running {
		ps.mu.Unlock()
		return
	}
	ps.running = true
	ps.stopCh = make(chan struct{})
	ps.mu.Unlock()

	go ps.pollLoop()
}

// Stop stops the streaming goroutine.
func (ps *PriceStreamer) Stop() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if !ps.running {
		return
	}
	close(ps.stopCh)
	ps.running = false
}

// IsRunning returns whether the streamer is actively polling.
func (ps *PriceStreamer) IsRunning() bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.running
}

// GetLastPrice returns the last known price for a symbol, or nil if none.
func (ps *PriceStreamer) GetLastPrice(symbol string) *PriceUpdate {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if p, ok := ps.lastPrices[symbol]; ok {
		cp := *p
		return &cp
	}
	return nil
}

// pollLoop is the main polling goroutine.
func (ps *PriceStreamer) pollLoop() {
	for {
		ps.mu.Lock()
		interval := ps.interval
		symbols := make([]string, len(ps.symbols))
		copy(symbols, ps.symbols)
		ps.mu.Unlock()

		for _, sym := range symbols {
			select {
			case <-ps.stopCh:
				return
			default:
			}
			ps.pollSymbol(sym)
		}

		select {
		case <-ps.stopCh:
			return
		case <-time.After(interval):
		}
	}
}

// pollSymbol fetches the latest tick for a single symbol and fires the callback if changed.
func (ps *PriceStreamer) pollSymbol(symbol string) {
	ticks, ret := ps.api.TickLast(symbol)
	if ret != MT_RET_OK || len(ticks) == 0 {
		return
	}
	t := ticks[0]
	update := &PriceUpdate{
		Symbol:     symbol,
		Bid:        t.Bid,
		Ask:        t.Ask,
		Last:       t.Last,
		Volume:     t.Volume,
		VolumeReal: t.VolumeReal,
		Datetime:   t.Datetime,
		Spread:     t.Ask - t.Bid,
	}

	ps.mu.Lock()
	prev, exists := ps.lastPrices[symbol]
	changed := !exists || prev.Bid != update.Bid || prev.Ask != update.Ask || prev.Last != update.Last
	if changed {
		ps.lastPrices[symbol] = update
	}
	ps.mu.Unlock()

	if changed && ps.callback != nil {
		ps.callback(update)
	}
}

// BarBuilder builds OHLC bars from streaming tick data.
type BarBuilder struct {
	symbol     string
	period     time.Duration
	callback   func(bar *MTBar)
	currentBar *MTBar
	barStart   time.Time
	mu         sync.Mutex
}

// NewBarBuilder creates a bar builder for a symbol with the given period.
// The callback is invoked when a completed bar is finalized.
func NewBarBuilder(symbol string, period time.Duration, callback func(bar *MTBar)) *BarBuilder {
	return &BarBuilder{
		symbol:   symbol,
		period:   period,
		callback: callback,
	}
}

// OnTick processes a new tick price and timestamp, building OHLC bars.
// When a bar's period expires, the completed bar is sent to the callback
// and a new bar is started.
func (bb *BarBuilder) OnTick(price float64, timestamp time.Time) {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	// Truncate timestamp to the period boundary
	barTime := timestamp.Truncate(bb.period)

	if bb.currentBar == nil || barTime.After(bb.barStart) {
		// Finalize the previous bar
		if bb.currentBar != nil && bb.callback != nil {
			completedBar := *bb.currentBar
			go bb.callback(&completedBar)
		}
		// Start a new bar
		bb.barStart = barTime
		bb.currentBar = &MTBar{
			Datetime: barTime.Unix(),
			Open:     price,
			High:     price,
			Low:      price,
			Close:    price,
			TickVol:  1,
		}
		return
	}

	// Update the current bar
	if price > bb.currentBar.High {
		bb.currentBar.High = price
	}
	if price < bb.currentBar.Low {
		bb.currentBar.Low = price
	}
	bb.currentBar.Close = price
	bb.currentBar.TickVol++
}

// CurrentBar returns a copy of the current incomplete bar, or nil if none.
func (bb *BarBuilder) CurrentBar() *MTBar {
	bb.mu.Lock()
	defer bb.mu.Unlock()
	if bb.currentBar == nil {
		return nil
	}
	cp := *bb.currentBar
	return &cp
}
