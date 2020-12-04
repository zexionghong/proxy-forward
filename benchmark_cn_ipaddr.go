package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}
type Data struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func get_ping_request(data *Data) error {
	request, _ := http.NewRequest("GET", "http://myip.ipip.net", nil)
}

func send_apply_request(geo string, wg *sync.WaitGroup) (*Data, error) {
	defer wg.Done()
	result := Result{}
	defer func() {
		log.Println(result)
	}()
	request, _ := http.NewRequest("GET", "http://test-api.proxy302.com/api/v1/create_proxy/area?country="+geo, nil)
	request.SetBasicAuth("test", "oOZESDk7Rz3lhuWm")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func send_apply_request_v1(geo string, wg *sync.WaitGroup, ipProxys []*Data) error {
	defer wg.Done()
	result := Result{}
	defer func() {
		log.Println(result)
	}()
	request, _ := http.NewRequest("GET", "http://test-api.proxy302.com/api/v1/create_proxy/area?country="+geo, nil)
	request.SetBasicAuth("test", "oOZESDk7Rz3lhuWm")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &result); err != nil {
		return err
	}
	ipProxys = append(ipProxys, &result.Data)
	return nil
}

func benchmark_ping_request(data *Data) (int, int) {
	var (
		success int
		fail    int
	)
	for i := 0; i < num; i++ {
	}
}

// benchmark apply ip
func benchmark_apply_ip_addr(num int) {
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go send_apply_request("CN", &wg)
	}
	wg.Wait()
}

// benchmark ip 并发
func benchmark_used_ip_addr(num int) {
	ipProxys := []*Data{}
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go send_apply_request_v1("CN", &wg, ipProxys)
	}
	wg.Wait()
}

func main() {
	benchmark_apply_ip_addr(50)
}
