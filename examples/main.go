package main

import (
	"context"
	"fmt"
	"log"

	gopastemyst "github.com/Sammie156/go-pastemyst"
)

func main() {
	pasteID := "54ntx8jj"
	client := gopastemyst.NewClient("")

	paste, err := client.GetPaste(context.Background(), pasteID)
	if err != nil {
		log.Fatalf("Failed to get paste: %v", err)
	}

	fmt.Printf("Got paste! Title: %s\n", paste.Title)

	for _, pasty := range paste.Pasties {
		fmt.Printf("Pasty ID: %s, Language: %s, Title: %s\n", pasty.ID, pasty.Language, pasty.Title)
		// fmt.Printf("Pasty Content : %s\n", pasty.Content)
	}

	stats, err := client.GetPasteStats(context.Background(), pasteID)
	if err != nil {
		log.Fatalf("Failed to get paste stats: %v", err)
	}

	for pastyID, pasty := range stats.Pasties {
		fmt.Printf("Pasty ID: %s, Lines: %d, Words: %d\n", pastyID, pasty.Lines, pasty.Words)
	}
}
