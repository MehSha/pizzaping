package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Restaurant struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Rating          float32 `json:"rating"`
	MinimumDelivery int     `json:"minimum_delivery_time"`
}

type searchResponse struct {
	Data struct {
		Items []*Restaurant `json:"items"`
	} `json:"data"`
}

func search(name string, lat, lon float32) ([]*Restaurant, error) {
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

type product struct {
	Name       string `json:"name"`
	ID         int    `json:"id"`
	Variations []struct {
		Name  string  `json:"name"`
		ID    int     `json:"id"`
		Price float32 `json:"price"`
	} `json:"product_variations"`
}
type category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []product `json:"products"`
}

type dishResponse struct {
	Data struct {
		Address string `json:"address"`
		Menus   []struct {
			ID         int        `json:"id"`
			Categories []category `json:"menu_categories"`
		} `json:"menus"`
	} `json:"data"`
}

func getDish(restaurantID int, name string, lat, lon float32) (*Dish, error) {
	dish := &Dish{}
	searchURL := fmt.Sprintf("https://de.fd-api.com/api/v5/vendors/%d?include=product_variations_normalized&latitude=%f&longitude=%f&opening_type=delivery",
		restaurantID, lat, lon)
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return dish, err
	}
	req.Header.Set("x-fp-api-key", "volo")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dish, err
	}
	if resp.StatusCode != 200 {
		return dish, fmt.Errorf("status code is %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dish, fmt.Errorf("can not read response body: %s", err)
	}
	// fmt.Printf("dish body is %s\n", body)
	dishRsp := &dishResponse{}
	err = json.Unmarshal(body, &dishRsp)
	if err != nil {
		return dish, fmt.Errorf("can dish unmarshal body: %s", err)
	}
	if len(dishRsp.Data.Menus) < 1 {
		return dish, fmt.Errorf("restaurant has no menu")
	}
	menu := dishRsp.Data.Menus[0]
	// fmt.Printf("dish response is %+v\n", dishRsp)
	fmt.Printf("restaurant id %+v\n", restaurantID)
	for _, cat := range menu.Categories {
		if strings.Contains(strings.ToLower(cat.Name), "pizza") {
			// fmt.Printf("found pizza category %+v\n", cat)
			//now search for actual product
			for _, prd := range cat.Products {
				fmt.Printf("checking product our search: %s, product name: %s, ok? %t\n", name, prd.Name, strings.Contains(strings.ToLower(prd.Name), name))
				if strings.Contains(strings.ToLower(prd.Name), name) {
					dish.Name = prd.Name
					dish.Variations = prd.Variations
					dish.ID = prd.ID
					return dish, nil
					// fmt.Printf("found dish %+v\n", prd)
				}
			}
		}
	}
	return dish, fmt.Errorf("dish was not found!")
}
