package tibiadata

import (
	"encoding/json"
	"log"
	"time"
)

const isoDateTimeFormat = "2006-01-02T15:04:05Z"

type information struct {
	Timestamp string `json:"timestamp"`
}

type characterResponse struct {
	Characters  characters  `json:"characters"`
	Information information `json:"information"`
}

type characters struct {
	Deaths []death `json:"deaths"`
}

type death struct {
	Time   string `json:"time"`
	Level  int    `json:"level"`
	Reason string `json:"reason"`
}

type Death struct {
	Time   time.Time
	Level  int
	Reason string
}

func FetchPlayersDeaths(playerName string) ([]Death, error) {
	url := "https://api.tibiadata.com/v3/character/" + playerName
	body, err := fetch(url)
	if err != nil {
		return nil, err
	}
	var cr characterResponse
	err = json.Unmarshal(body, &cr)
	if err != nil {
		log.Printf("%v", string(body))
		return nil, err
	}

	var deaths []Death
	for _, death := range cr.Characters.Deaths {
		parsedTime, err := time.Parse(isoDateTimeFormat, death.Time)
		if err != nil {
			log.Printf("Failed to parse time %v", err)
			continue
		}
		deaths = append(deaths, Death{
			Level:  death.Level,
			Reason: death.Reason,
			Time:   parsedTime,
		})
	}

	return deaths, nil
}
