package main

import (
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
  	tcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
  	"github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg/request"

	"os"
	"time"
	"encoding/json"
	"github.com/joho/godotenv"
)

type Card struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
}

var storedCards []Card

func clearCookies(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name: "storedCards",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	return c.Next()
}

/* Setup App */
func setupApp() *fiber.App {
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))


	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Clear stored cards cookie
	app.Use(clearCookies)

	// Deinfe Routes
	app.Get("/random-card", getRandomCard)
	app.Post("/store", storeCard)
	app.Get("/stored", getStoredCards)

	return app
}

/* Routes */
func storeCard(c *fiber.Ctx) error {
	var card Card
	if err := c.BodyParser(&card); err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	// Get the existing cookie
	existingCookie := c.Cookies("storedCards")
	var selectedCards []Card
		
	// Only unmarshal if the cookie is not empty
	if existingCookie != "" {
		if err := json.Unmarshal([]byte(existingCookie), &selectedCards); err != nil {
			return c.Status(500).SendString("Error unmarshalling cookie")
		}
	}

	// Add the new card if its not already in the cookie
	for _, storedCard := range selectedCards {
		if storedCard.ID == card.ID {
			return c.SendString("Card already stored")
		}
	}

	// Check if the limit of 6 cards is reached
	if len(selectedCards) >= 6 {
		return c.Status(400).SendString("Limit of 6 cards reached")
	}

	// Add the new card
	selectedCards = append(selectedCards, card)

	// Covert back to JSON and update the cookie
	updatedCookie, err := json.Marshal(selectedCards)
	if err != nil {
		return c.Status(500).SendString("Error marshalling cookie")
	}

	c.Cookie(&fiber.Cookie{
		Name: "storedCards",
		Value: string(updatedCookie),
		Expires: time.Now().Add(time.Hour * 24),
	})

	return c.Status(200).SendString("Card stored successfully")
}

func getStoredCards(c *fiber.Ctx) error {

	// Get cookie
	existingCookie := c.Cookies("storedCards")
	var selectedCards []Card
	
	if existingCookie == "" {
		return c.JSON([]Card{})
	}

	if err := json.Unmarshal([]byte(existingCookie), &selectedCards); err != nil {
		return c.Status(500).SendString("Error unmarshalling cookie")
	}

	return c.JSON(selectedCards)
}

func getRandomCard(c *fiber.Ctx) error {

	apiKey := os.Getenv("POKEMON_API_KEY")
	if apiKey == "" {
		log.Fatal("POKEMON_API_KEY is not set")
	}

	// Create tcg client to make API calls and use our API key
  	cardClient := tcg.NewClient(apiKey)
	tcgTotalCards := 18506
	randomPage := rand.Intn((tcgTotalCards / 250) + 1) 
	cards, err := cardClient.GetCards(
		request.Page(randomPage),
		request.PageSize(1), // Needs to stay 1 so that we only get one card back
	)

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Simplify what we return from the API.

	return c.JSON(cards);
}

func main() {
	app := setupApp()
	app.Listen(":3000")
}