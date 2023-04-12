package tibiadata

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Vocation string

type player struct {
	Name string `json:"name"`
}

type world struct {
	Name          string   `json:"name"`
	OnlinePlayers []player `json:"online_players"`
}

type worlds struct {
	World world `json:"world"`
}

type worldResponse struct {
	Worlds worlds `json:"worlds"`
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close body")
		}
	}(resp.Body)

	return io.ReadAll(resp.Body)
}

func FetchWorldOnlinePlayers(world string) ([]string, error) {
	url := "https://api.tibiadata.com/v3/world/" + world
	body, err := fetch(url)
	if err != nil {
		return nil, err
	}
	var wr worldResponse
	err = json.Unmarshal(body, &wr)
	if err != nil {
		return nil, err
	}
	var playerNames []string
	for _, onlinePlayer := range wr.Worlds.World.OnlinePlayers {
		playerNames = append(playerNames, onlinePlayer.Name)
	}

	return playerNames, nil
}
