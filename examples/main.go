package main

import (
	"context"
	"log"
	"os"

	gopastemyst "github.com/Sammie156/go-pastemyst"
	"github.com/joho/godotenv"
)

// TODO: Work on making better examples with understandable code and comments
func main() {
	err := godotenv.Load("examples/.env")

	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	token := os.Getenv("PASTEMYST_TOKEN")

	client := gopastemyst.NewClient(token)
	ctx := context.Background()

	client.PrivatePaste(ctx, "9xx98zrf")
}
