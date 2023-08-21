package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContainsIP(t *testing.T) {
	ipList := []string{"192.168.1.1", "10.0.0.1", "127.0.0.1"}
	ipToFind := "10.0.0.1"
	ipNotInList := "8.8.8.8"

	if !containsIP(ipList, ipToFind) {
		t.Errorf("Se esperaba que containsIP devolviera true para la IP %s, pero devolvió false", ipToFind)
	}

	if containsIP(ipList, ipNotInList) {
		t.Errorf("Se esperaba que containsIP devolviera false para la IP %s, pero devolvió true", ipNotInList)
	}
}

func TestGetIpAddressInfo(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseJSON := `{"country_name": "Estados Unidos", "country_code": "US"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer mockServer.Close()

	ipAddress := "8.8.8.8"
	expectedCountryCode := "US"

	recorder := httptest.NewRecorder()
	getIpAddressInfo(ipAddress, recorder)

	if recorder.Code != http.StatusOK {
		t.Errorf("Se esperaba el código de estado %d, pero se obtuvo %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, expectedCountryCode) {
		t.Errorf("Se esperaba que el cuerpo de la respuesta contuviera el código de país %s, pero no lo hizo", expectedCountryCode)
	}
}

func TestGetCurrencyInfo(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseJSON := `{"currencies": {"USD": {"name": "Dólar estadounidense", "symbol": "$"}}}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer mockServer.Close()

	countryCode := "US"
	expectedCurrencyCode := "USD"

	recorder := httptest.NewRecorder()
	getCurrencyInfo(countryCode, recorder)

	if recorder.Code != http.StatusOK {
		t.Errorf("Se esperaba el código de estado %d, pero se obtuvo %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, expectedCurrencyCode) {
		t.Errorf("Se esperaba que el cuerpo de la respuesta contuviera el código de moneda %s, pero no lo hizo", expectedCurrencyCode)
	}
}

func TestGetExchangeInfo(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseJSON := `{"currency_pair": "USD_USD", "exchange_rate": 1.0}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer mockServer.Close()

	currencyCode := "USD"

	recorder := httptest.NewRecorder()
	getExchangeInfo(currencyCode, recorder)

	if recorder.Code != http.StatusOK {
		t.Errorf("Se esperaba el código de estado %d, pero se obtuvo %d", http.StatusOK, recorder.Code)
	}

	body := recorder.Body.String()
	expectedOutput := "Par de monedas: USD_USD\nTipo de cambio: 1.000000"
	if !strings.Contains(body, expectedOutput) {
		t.Errorf("Se esperaba que el cuerpo de la respuesta contuviera: %s\nCuerpo actual: %s", expectedOutput, body)
	}
}
