package main

import (
	"log"
	"sync"
	"tibia-death-tracker/tibiadata"
	"time"
)

func main() {
	world := "Antica"
	workers := 8

	workInput := make(chan string, 2000)
	minDeathTime := time.Now()
	log.Printf("Starting %v", minDeathTime)

	lastPlayersDeath := make(map[string]time.Time)
	var m sync.RWMutex

	for i := 0; i < workers; i++ {
		go func() {
			for playerName := range workInput {
				deaths, err := tibiadata.FetchPlayersDeaths(playerName)
				if err != nil {
					log.Printf("Failed to fetch deaths of %v %v", playerName, err)
				}
				for i := len(deaths) - 1; i >= 0; i-- {
					death := deaths[i]
					if death.Time.After(minDeathTime) {
						m.RLock()
						last, exists := lastPlayersDeath[playerName]
						m.RUnlock()

						if !exists || death.Time.After(last) {
							m.Lock()
							lastPlayersDeath[playerName] = death.Time
							m.Unlock()
							log.Printf("%v %v %v", playerName, death.Time, death.Reason)
						}

					}
				}
			}
		}()
	}

	ticker := time.NewTicker(65 * time.Second)
	lastSeenTime := make(map[string]time.Time)
	for range ticker.C {
		players, err := tibiadata.FetchWorldOnlinePlayers(world)
		if err != nil {
			log.Printf("Failed to fetch players %v", err)
		} else {
			for _, playerName := range players {
				lastSeenTime[playerName] = time.Now()
			}
		}

		limit := time.Now().Add(-15 * time.Minute)
		for playerName, seenTime := range lastSeenTime {
			if seenTime.After(limit) {
				workInput <- playerName
			}
		}
	}
}
