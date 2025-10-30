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

	// paste, err := client.GetPaste(context.Background(), pasteID)
	// if err != nil {
	// 	log.Fatalf("Failed to get paste: %v", err)
	// }

	// fmt.Printf("Got paste! Title: %s\n", paste.Title)

	// for _, pasty := range paste.Pasties {
	// 	fmt.Printf("Pasty ID: %s, Language: %s, Title: %s\n", pasty.ID, pasty.Language, pasty.Title)
	// 	// fmt.Printf("Pasty Content : %s\n", pasty.Content)
	// }

	// stats, err := client.GetPasteStats(context.Background(), pasteID)
	// if err != nil {
	// 	log.Fatalf("Failed to get paste stats: %v", err)
	// }

	// for pastyID, pasty := range stats.Pasties {
	// 	fmt.Printf("Pasty ID: %s, Lines: %d, Words: %d\n", pastyID, pasty.Lines, pasty.Words)
	// }

	pasty := gopastemyst.CreatePastyOptions{
		Title:    "API Test",
		Content:  "Testing API Calling",
		Language: "text",
	}

	options := gopastemyst.CreatePasteOptions{
		Title:     "My Private Go API Test",
		ExpiresIn: "1h",
		Private:   true,
		Pasties:   []gopastemyst.CreatePastyOptions{pasty},
	}

	log.Println("Attempting to create Paste")
	newPaste, err := client.CreatePaste(ctx, options)

	if err != nil {
		log.Fatalf("Failed to create paste: %v", err)
	}

	fmt.Printf("Created Paste! ID: %s\n", newPaste.ID)

}
