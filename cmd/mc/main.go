package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"seitaikei/utils"
	"time"

	"github.com/Cryptkeeper/go-minecraftping"
)

type ss struct {
	Type    string      `json:"type"`
	Players interface{} `json:"players"`
}

func main() {
	utils.Banner()
	log.Print("Minecraft Service Starting Up....")

	go func() {
		// Gets cursor
		fmt.Print("\033[s")
		for {
			<-time.After(time.Second * 5)

			resp, err := minecraftping.Ping("192.168.0.60", 25565, 760, time.Second*5)
			if err != nil {
				log.Fatal(err)
			}

			// Restors cursor pos and clears line
			fmt.Print("\033[u\033[K")
			fmt.Printf("%d/%d players are online.\n", resp.Players.Online, resp.Players.Max)

			for _, player := range resp.Players.Sample {
				fmt.Println(player.Name)
			}

			tmp := ss{
				Type:    "mc",
				Players: resp.Players,
			}
			data, _ := json.Marshal(tmp)
			r, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/mc", bytes.NewBuffer(data))
			r.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			client.Do(r)
		}
	}()
	fmt.Scanln()
}
