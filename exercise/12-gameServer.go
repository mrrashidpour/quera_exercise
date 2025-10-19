package main

import (
	"errors"
	"strings"
	"sync"
)

// ---------------------------
// Types
// ---------------------------

type Player struct {
	name        string // lowercase key
	displayName string // capitalized for display
	ch          chan string
	currentMap  *Map
	mu          sync.RWMutex
}

type Map struct {
	id        int
	players   map[string]*Player
	broadcast chan mapMessage
	mu        sync.RWMutex
}

type Game struct {
	maps    map[int]*Map
	players map[string]*Player
	mu      sync.RWMutex
}

type mapMessage struct {
	fromLower   string
	fromDisplay string
	text        string
}

// ---------------------------
// Helpers
// ---------------------------

func capitalize(s string) string {
	if s == "" {
		return s
	}
	s = strings.ToLower(s)
	r := []rune(s)
	r[0] = []rune(strings.ToUpper(string(r[0])))[0]
	return string(r)
}

// ---------------------------
// Game functions / methods
// ---------------------------

func NewGame(mapIds []int) (*Game, error) {
	g := &Game{
		maps:    make(map[int]*Map),
		players: make(map[string]*Player),
	}

	for _, id := range mapIds {
		if id <= 0 {
			return nil, errors.New("invalid map id")
		}
		// avoid duplicate map id: if exists, skip creating again
		if _, exists := g.maps[id]; exists {
			continue
		}
		m := &Map{
			id:        id,
			players:   make(map[string]*Player),
			broadcast: make(chan mapMessage, 200),
		}
		g.maps[id] = m

		// start fan-out goroutine for this map
		go m.FanOutMessages()
	}

	return g, nil
}

func (g *Game) ConnectPlayer(name string) error {
	if name == "" {
		return errors.New("invalid player name")
	}
	lower := strings.ToLower(name)

	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.players[lower]; exists {
		return errors.New("player already connected")
	}

	p := &Player{
		name:        lower,
		displayName: capitalize(name),
		ch:          make(chan string, 100),
		currentMap:  nil,
	}
	g.players[lower] = p
	return nil
}

func (g *Game) SwitchPlayerMap(name string, mapId int) error {
	if mapId <= 0 {
		return errors.New("invalid map id")
	}
	lower := strings.ToLower(name)

	// locate player
	g.mu.RLock()
	p, pExists := g.players[lower]
	g.mu.RUnlock()
	if !pExists {
		return errors.New("player not found")
	}

	// locate map
	g.mu.RLock()
	m, mExists := g.maps[mapId]
	g.mu.RUnlock()
	if !mExists {
		return errors.New("map not found")
	}

	// check if already in that map
	p.mu.RLock()
	if p.currentMap != nil && p.currentMap.id == mapId {
		p.mu.RUnlock()
		return errors.New("player already in this map")
	}
	p.mu.RUnlock()

	// remove from old map if any
	if p.currentMap != nil {
		old := p.currentMap
		old.mu.Lock()
		delete(old.players, lower)
		old.mu.Unlock()
	}

	// add to new map
	m.mu.Lock()
	m.players[lower] = p
	m.mu.Unlock()

	// set player's currentMap
	p.mu.Lock()
	p.currentMap = m
	p.mu.Unlock()

	return nil
}

func (g *Game) GetPlayer(name string) (*Player, error) {
	lower := strings.ToLower(name)
	g.mu.RLock()
	defer g.mu.RUnlock()
	if p, ok := g.players[lower]; ok {
		return p, nil
	}
	return nil, errors.New("player not found")
}

func (g *Game) GetMap(mapId int) (*Map, error) {
	if mapId <= 0 {
		return nil, errors.New("invalid map id")
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	if m, ok := g.maps[mapId]; ok {
		return m, nil
	}
	return nil, errors.New("map not found")
}

// ---------------------------
// Map methods
// ---------------------------

func (m *Map) FanOutMessages() {
	for msg := range m.broadcast {
		// snapshot players to avoid holding lock while sending
		m.mu.RLock()
		players := make([]*Player, 0, len(m.players))
		for _, p := range m.players {
			players = append(players, p)
		}
		m.mu.RUnlock()

		for _, p := range players {
			// do not deliver to sender
			if p == nil {
				continue
			}
			if p.name == msg.fromLower {
				continue
			}
			// send (blocking). channels are buffered (100) and problem guarantees max 100 concurrent messages.
			p.ch <- (msg.fromDisplay + " says: " + msg.text)
		}
	}
}

// ---------------------------
// Player methods
// ---------------------------

func (p *Player) GetChannel() <-chan string {
	return p.ch
}

func (p *Player) SendMessage(msg string) error {
	if msg == "" {
		return errors.New("message is empty")
	}
	p.mu.RLock()
	m := p.currentMap
	p.mu.RUnlock()
	if m == nil {
		return errors.New("player is not in a map")
	}

	// prepare message
	mMsg := mapMessage{
		fromLower:   p.name,
		fromDisplay: p.displayName,
		text:        msg,
	}

	// send to map's broadcast
	m.broadcast <- mMsg
	return nil
}

func (p *Player) GetName() string {
	return p.displayName
}
