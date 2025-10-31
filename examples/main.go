package main

import (
	"context"
	"fmt"
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

	pasteDiff, err := client.GetDiffAtCertainEdit(ctx, "1shjmsu7", "l051n8gw")

	fmt.Printf("Old Paste: %s\n", pasteDiff.OldPaste.Pasties[2].Content)
	fmt.Printf("New Paste: %s\n", pasteDiff.NewPaste.Pasties[2].Content)
}
