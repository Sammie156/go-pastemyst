package main

import (
	"context"
	"fmt"
	"log"

	gopastemyst "github.com/Sammie156/go-pastemyst"
)

func main() {
	client := gopastemyst.NewClient("")

	paste, err := client.GetPaste(context.Background(), "17dex03t")
	if err != nil {
		log.Fatalf("Failed to get paste: %v", err)
	}

	fmt.Printf("Got paste! Title: %s\n", paste.Title)

	for _, pasty := range paste.Pasties {
		fmt.Printf("Pasty ID: %s, Language: %s, Title: %s\n", pasty.ID, pasty.Language, pasty.Title)
		fmt.Printf("Pasty Content : %s\n", pasty.Content)
	}
}
