package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func searchRestaurant(name string, lat, lon float32) ([]*Restaurant, error) {
	searchResp := searchResponse{}
	result := []*Restaurant{}
	//call backend
	searchURL := fmt.Sprintf("https://de.fd-api.com/api/v5/vendors?search_term=%s&latitude=%f&longitude=%f&opening_type=delivery", name, lat, lon)
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("x-fp-api-key", "volo")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != 200 {
		return result, fmt.Errorf("status code is %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("can not read response body: %s", err)
	}
	// fmt.Println("search rs is: ", string(body))
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		return result, fmt.Errorf("can unmarshal body: %s", err)
	}
	for _, rs := range searchResp.Data.Items {
		result = append(result, rs)
	}
	return result, nil
}

type Dish struct {
	Name       string      `json:"name"`
	ID         int         `json:"id"`
	Restaurant *Restaurant `json:"restaurant"`
	Variations []struct {
		Name  string  `json:"name"`
		ID    int     `json:"id"`
		Price float32 `json:"price"`
	} `json:"variations"`
}

func getDish(rstID int, name string) (*Dish, error) {
	return &Dish{}, nil
}
