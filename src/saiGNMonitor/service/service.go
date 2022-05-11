package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/webmakom-com/saiGNMonitor/config"
)

type Service struct {
	config config.Configuration
}

type Result struct {
	Time int `json:"t"`
	Value float64 `json:"v"`
}

type Response []Result

const (
	MinerRevenue = "revenue_sum"
	HashRate     = "hash_rate_mean"
)

func NewService() Service {
	return Service{
		config: config.Get(),
	}
}

func (s Service) Start() {
	fmt.Println("Start")
	for range time.Tick(time.Second * time.Duration(s.config.GlassNode.Period)) {
		fmt.Println("Tick")
		s.Execute()
	}
}

func (s Service) Execute() {
	attempts := 1
	for range time.Tick(time.Second * time.Duration(s.config.GlassNode.Interval)) {
		fmt.Println("Attempts:", attempts)
		if attempts > s.config.GlassNode.Retries {
			break
		}

		mr, err := s.sendRequest(MinerRevenue)
		if err != nil {
			fmt.Println(err)
			attempts++
			return
		}

		hr, err := s.sendRequest(HashRate)
		if err != nil {
			fmt.Println(err)
			attempts++
			return
		}

		var mrt Response
		var hrt Response

		mrErr := json.Unmarshal(mr, &mrt)
		if mrErr != nil {
			fmt.Println(mrErr)
			attempts++
			return
		}

		hrErr := json.Unmarshal(hr, &hrt)
		if hrErr != nil {
			fmt.Println(hrErr)
			attempts++
			return
		}

		profit := mrt[len(mrt)-1].Value / hrt[len(hrt)-1].Value * math.Pow(10, 9)

		fmt.Println("Profit:", profit)
		//Execute contract call

		break
	}
}

func (s Service) sendRequest(path string) ([]byte, error) {
	url := s.config.GlassNode.Url + path + "?a=BTC"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Service error: ", err)
		return []byte(""), err
	}

	req.Header.Set("X-Api-Key", s.config.GlassNode.Key)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Service error: ", err)
		return []byte(""), err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}