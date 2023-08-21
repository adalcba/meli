package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
}

var bannedIps []string

// @title MELI REST API
// @version 1.0
// @description Get Country and Currency info based on Ip Address.

// @host localhost:8080
// @BasePath /
func main() {
	/*http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"), // Replace with your correct JSON file path
	))*/

	// Create a router
	http.HandleFunc("/getAllInfo", ipInfoHandler)

	// Start the server
	port := ":8080"
	fmt.Printf("Server listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// @Summary Get IP information
// @Description Get information about an IP address
// @Param ip query string true "IP Address"
// @Success 200 {object} Response
// @Router /getAllInfo [get]
func ipInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		buttonName := r.FormValue("button_name")

		if buttonName == "getInfo" {
			if !containsIP(bannedIps, r.FormValue("ip_address")) {
				ipAddress := r.FormValue("ip_address")
				countryCode := getIpAddressInfo(ipAddress, w)
				currencyCode := getCurrencyInfo(countryCode, w)
				getExchangeInfo(currencyCode, w)
			} else {
				http.Error(w, "Access denied for this IP address.", http.StatusForbidden)
			}

		} else if buttonName == "customAction" {
			banIp(r)
		}
	} else {
		formHTML := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>IP Address Info</title>
		</head>
		<body>
			<h1>IP Address Info</h1>
			<form method="post">
				<label for="ip_address">Enter IP Address:</label>
				<input type="text" id="ip_address" name="ip_address">
				<button type="submit" name="button_name" value="getInfo">Get Info</button>
				<button type="submit" name="button_name" value="customAction">Custom Action</button>
			</form>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formHTML))
	}
}

func getIpAddressInfo(ipAddress string, w http.ResponseWriter) string {
	url := fmt.Sprintf("http://api.ipapi.com/api/%s?access_key=1cd1dce7bfa968fed7a39df86ceb70e9", ipAddress)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	name := data["country_name"].(string)
	code := data["country_code"].(string)

	fmt.Println("Country Name:", name)
	fmt.Println("Country Code:", code)

	// Create the formatted response
	response := fmt.Sprintf("Country Name: %s\nCountry Code: %s\n\n", name, code)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

	code = data["country_code"].(string)
	return code
}

func getCurrencyInfo(countryCode string, w http.ResponseWriter) string {
	url := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s?fields=currencies", countryCode)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var result map[string]map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Define a struct to match the JSON structure
	var data struct {
		Currencies map[string]struct {
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
		} `json:"currencies"`
	}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Extract the key
	var response string
	var _currencyCode string
	for currencyCode, currencyData := range data.Currencies {
		_currencyCode = currencyCode
		response += fmt.Sprintf("Currency Code: %s\n", currencyCode)
		response += fmt.Sprintf("Name: %s\n", currencyData.Name)
		response += fmt.Sprintf("Symbol: %s\n\n", currencyData.Symbol)
	}

	fmt.Println(response)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

	return _currencyCode

}

func getExchangeInfo(currencyCode string, w http.ResponseWriter) {

	type Data struct {
		CurrencyPair string  `json:"currency_pair"`
		ExchangeRate float64 `json:"exchange_rate"`
	}

	url := fmt.Sprintf("https://api.api-ninjas.com/v1/exchangerate?pair=%s_USD", currencyCode)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Format and print the desired output
	output := fmt.Sprintf("Currency Pair:  %s\nExchange Rate: %.6f",
		data.CurrencyPair, data.ExchangeRate)

	fmt.Println(output)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

func banIp(r *http.Request) {
	ip := r.FormValue("ip_address")
	bannedIps = append(bannedIps, ip)
}

func containsIP(slice []string, ip string) bool {
	for _, item := range slice {
		if item == ip {
			return true
		}
	}
	return false
}
