package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"seitaikei/utils"
	"time"
)

type Status struct {
	Ip   string `json:"ip"`
	Code string `json:"code"`
}
type ss struct {
	Type  string   `json:"type"`
	Stats []Status `json:"stats"`
}

var (
	servers = []string{
		"https://192.168.0.6/",
		"http://192.168.0.31:8096/",
	}
	stats = ss{Type: "misc"}
)

func main() {
	utils.Banner()
	log.Println("Starting Misc Service ....")

	go func() {
		for {
			for _, server := range servers {
				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true}},
				}
				req, err := http.NewRequest("GET", server, nil)
				if err != nil {
					log.Println(err)
					return
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("%s status: %s", server, resp.Status)
				stats.Stats = append(stats.Stats, Status{
					Ip:   server,
					Code: resp.Status})

				data, _ := json.Marshal(stats)
				r, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/misc", bytes.NewBuffer(data))
				r.Header.Set("Content-Type", "application/json")
				client.Do(r)
			}
			<-time.After(time.Minute * 10)
		}
	}()

	fmt.Scanln()
}
