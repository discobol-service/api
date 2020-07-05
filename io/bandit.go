package io


import (
	"sync"
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
	"fmt"

	"log"
)

type Bandit struct {
	Address string
}

type BanditStat struct {
	Arm    string  `json:"arm"`
	Scores float64 `json:"scores"`
}

var once2 sync.Once
var bandit *Bandit


func GetBandit(address string) *Bandit {
	once2.Do(func() {
		bandit = &Bandit{address}
	})
	return bandit
}

// Взять веса урлов у бандита
func (b *Bandit) GetStat(domain string, requestJSON []byte) ([]BanditStat, error) {
	stat := []BanditStat{}

	if domain == "" {
		domain = "default"
	}

	requestURL := fmt.Sprintf("%s/stat/list/%s", b.Address, domain)

	log.Println("requestURL = " + requestURL)

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return stat, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return stat, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println("response Body:", string(body))


	if err = json.Unmarshal(body, &stat); err != nil {
		return stat, err
	}

	return stat, nil
}