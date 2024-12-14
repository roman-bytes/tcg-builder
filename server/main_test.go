package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestStoreCard(t *testing.T) {
	app := setupApp()

    card := Card{Name: "Pikachu", Image: "https://example.com/pikachu.png", ID: "1"}
    jsonCard, _ := json.Marshal(card)

    req := httptest.NewRequest("POST", "/store", bytes.NewReader(jsonCard))
    req.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(req)

    if err != nil {
        t.Fatalf("Failed to store card: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    // Test limit of stored cards
    for i := 0; i < 6; i++ {
        req := httptest.NewRequest("POST", "/store", bytes.NewReader(jsonCard))
        req.Header.Set("Content-Type", "application/json")
        app.Test(req)
    }

    // Attempt to store a seventh card
    req = httptest.NewRequest("POST", "/store", bytes.NewReader(jsonCard))
    req.Header.Set("Content-Type", "application/json")
    resp, err = app.Test(req)

    if resp.StatusCode != http.StatusBadRequest {
        t.Errorf("Expected status BadRequest, got %v", resp.Status)
    }
}

func TestGetStoredCards(t *testing.T) {
	app := setupApp()

    card := Card{Name: "Pikachu", Image: "https://example.com/pikachu.png", ID: "1"}
    jsonCard, _ := json.Marshal(card)

    req := httptest.NewRequest("POST", "/store", bytes.NewReader(jsonCard))
    req.Header.Set("Content-Type", "application/json")
    app.Test(req)

    req = httptest.NewRequest("GET", "/stored", nil)
    resp, err := app.Test(req)

    if err != nil {
        t.Fatalf("Failed to get stored cards: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    var cards []Card
    if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if len(cards) == 0 {
        t.Errorf("Expected at least one stored card, got none")
    }
}

func TestGetRandomCard(t *testing.T) {
	app := setupApp()
    // Mock the environment variable for testing
    os.Setenv("POKEMON_API_KEY", "test_api_key")

    req := httptest.NewRequest("GET", "/random-card", nil)
    resp, err := app.Test(req)

    if err != nil {
        t.Fatalf("Failed to get random card: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    var card Card
    if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if card.Name == "" {
        t.Errorf("Expected a card name, got empty")
    }
}