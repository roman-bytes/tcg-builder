package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"github.com/gofiber/fiber/v2"
)

// Helper function to perform HTTP requests in tests
func performRequest(app *fiber.App, method, url string, body []byte) (*http.Response, []byte) {
	req := bytes.NewReader(body)
	reqTest := httptest.NewRequest(method, url, req)
	reqTest.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(reqTest, -1) // Disable timeout for the request
	bodyRes, _ := io.ReadAll(res.Body)
	return res, bodyRes
}

// Test for the random-card route
func TestGetRandomCard(t *testing.T) {
	app := setupApp()

	// Perform GET request
	res, body := performRequest(app, "GET", "/random-card", nil)

	// Assert response
	assert.Equal(t, http.StatusOK, res.StatusCode, "Status should be 200 OK")

	// Check if the response has expected fields
	var cards []Card
	err := json.Unmarshal(body, &cards)
	assert.NoError(t, err, "Response should be a valid JSON")
	assert.NotEmpty(t, cards, "Response should contain at least one card")
}

// Test for storing a card
func TestStoreCard(t *testing.T) {
	app := setupApp()

	// Mock card to store
	card := Card{
		ID:    "xy7-54",
		Name:  "Mock Card",
		Image: "http://example.com/image.png",
	}
	cardJSON, _ := json.Marshal(card)

	// Perform POST request
	res, _ := performRequest(app, "POST", "/store", cardJSON)

	// Assert response
	assert.Equal(t, http.StatusOK, res.StatusCode, "Status should be 200 OK")
}

// Test for getting stored cards
func TestGetStoredCards(t *testing.T) {
	app := setupApp()

	// Mock card to store
	card := Card{
		ID:    "xy7-54",
		Name:  "Mock Card",
		Image: "http://example.com/image.png",
	}
	cardJSON, _ := json.Marshal(card)

	// First, perform POST request to store the card
	req := httptest.NewRequest("POST", "/store", bytes.NewReader(cardJSON))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1) // Disable timeout
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Retrieve the cookie set by the /store request
	cookie := res.Header.Get("Set-Cookie")
	assert.NotEmpty(t, cookie, "Cookie should be set when a card is stored")

	// Use the retrieved cookie in the GET request to /stored
	getReq := httptest.NewRequest("GET", "/stored", nil)
	getReq.Header.Set("Cookie", cookie)
	getRes, err := app.Test(getReq, -1)
	assert.NoError(t, err)

	// Read and unmarshal the response body
	body, err := io.ReadAll(getRes.Body)
	assert.NoError(t, err)

	var storedCards []Card
	err = json.Unmarshal(body, &storedCards)
	assert.NoError(t, err)

	// Validate the response
	assert.Equal(t, http.StatusOK, getRes.StatusCode)
	assert.Len(t, storedCards, 1, "There should be exactly one stored card")
	assert.Equal(t, card.ID, storedCards[0].ID, "Stored card ID should match")
}

// Test for storing a duplicate card
func TestStoreDuplicateCard(t *testing.T) {
	app := setupApp()

	// Mock card to store
	card := Card{
		ID:    "xy7-54",
		Name:  "Mock Card",
		Image: "http://example.com/image.png",
	}
	cardJSON, _ := json.Marshal(card)

	// First, perform POST request to store the card
	req1 := httptest.NewRequest("POST", "/store", bytes.NewReader(cardJSON))
	req1.Header.Set("Content-Type", "application/json")
	res1, err := app.Test(req1, -1) // Disable timeout
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res1.StatusCode)

	// Retrieve the cookie set by the first request
	cookie := res1.Header.Get("Set-Cookie")
	assert.NotEmpty(t, cookie, "Cookie should be set when a card is stored")

	// Second, perform POST request to store the same card (simulate duplicate)
	req2 := httptest.NewRequest("POST", "/store", bytes.NewReader(cardJSON))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Cookie", cookie) // Pass the previously set cookie
	res2, err := app.Test(req2, -1)
	assert.NoError(t, err)

	// Validate the response
	body, err := io.ReadAll(res2.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res1.StatusCode)
	assert.Equal(t, "Card already stored", string(body), "Duplicate card should not be stored again")
}

// Test for storing beyond the limit of 6 cards
func TestStoreCardLimit(t *testing.T) {
	app := setupApp()

	// Mock cards to store
	mockCards := []Card{
		{ID: "card1", Name: "Mock Card 1", Image: "http://example.com/image1.png"},
		{ID: "card2", Name: "Mock Card 2", Image: "http://example.com/image2.png"},
		{ID: "card3", Name: "Mock Card 3", Image: "http://example.com/image3.png"},
		{ID: "card4", Name: "Mock Card 4", Image: "http://example.com/image4.png"},
		{ID: "card5", Name: "Mock Card 5", Image: "http://example.com/image5.png"},
		{ID: "card6", Name: "Mock Card 6", Image: "http://example.com/image6.png"},
		{ID: "card7", Name: "Mock Card 7", Image: "http://example.com/image7.png"}, // Exceeds limit
	}

	var cookie string

	// Store up to 6 cards
	for _, card := range mockCards[:6] {
		cardJSON, _ := json.Marshal(card)
		req := httptest.NewRequest("POST", "/store", bytes.NewReader(cardJSON))
		req.Header.Set("Content-Type", "application/json")
		if cookie != "" {
			req.Header.Set("Cookie", cookie)
		}

		res, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// Retrieve cookie after each request
		cookie = res.Header.Get("Set-Cookie")
		assert.NotEmpty(t, cookie, "Cookie should be set after storing a card")

		// Validate response for each successful storage
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Card stored successfully", string(body), "Response should confirm card storage")
	}

	// Attempt to store the 7th card (exceeding the limit)
	card7JSON, _ := json.Marshal(mockCards[6])
	req := httptest.NewRequest("POST", "/store", bytes.NewReader(card7JSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie) // Pass the updated cookie

	res, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Validate that the limit is enforced
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "Limit of 6 cards reached", string(body), "Response should indicate card limit reached")
}
