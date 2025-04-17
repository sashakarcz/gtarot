package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"gopkg.in/yaml.v3"
)

//go:embed cards/*.png
var cardFS embed.FS

// Config defines the YAML input file structure
type Config struct {
	Cards  []string `yaml:"cards"`  // List of card names (may include "!" for reversed)
	Output string   `yaml:"output"` // Output file
}

// Card stores the name of a card and whether it should be reversed
type Card struct {
	Name     string
	Reversed bool
}

// parseYAML reads YAML into a Config struct
func parseYAML(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// cardToFilename maps card names to their embedded image filenames
func cardToFilename(card string) string {
	card = strings.ToLower(card)
	card = strings.ReplaceAll(card, " ", "_")
	card = strings.ReplaceAll(card, "-", "_")
	if strings.HasPrefix(card, "major_") || strings.HasPrefix(card, "minor_") {
		return card + ".png"
	}
	if strings.Contains(card, "_of_") {
		parts := strings.Split(card, "_of_")
		return fmt.Sprintf("minor_arcana_%s_%s.png", parts[1], parts[0])
	}
	return "major_arcana_" + card + ".png"
}

// loadCardImage reads a PNG image from the embedded filesystem
func loadCardImage(card string) (image.Image, error) {
	filename := "cards/" + cardToFilename(card)

	data, err := cardFS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read embedded card %s: %w", filename, err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("could not decode PNG %s: %w", filename, err)
	}

	return img, nil
}

// drawSpread renders the given cards into a horizontal spread
func drawSpread(cards []Card, output string) error {
	if len(cards) == 0 {
		return errors.New("no cards provided")
	}

	cardImages := make([]image.Image, 0, len(cards))
	for _, card := range cards {
		img, err := loadCardImage(card.Name)
		if err != nil {
			return fmt.Errorf("failed to load card '%s': %w", card.Name, err)
		}
		// If reversed, rotate the image 180 degrees
		if card.Reversed {
			w := img.Bounds().Dx()
			h := img.Bounds().Dy()
			reversed := gg.NewContext(w, h)
			reversed.RotateAbout(gg.Radians(180), float64(w)/2, float64(h)/2)
			reversed.DrawImage(img, 0, 0)
			img = reversed.Image()
		}
		cardImages = append(cardImages, img)
	}

	cardWidth := cardImages[0].Bounds().Dx()
	cardHeight := cardImages[0].Bounds().Dy()
	spacing := 20
	totalWidth := len(cardImages)*cardWidth + (len(cardImages)-1)*spacing
	totalHeight := cardHeight

	dc := gg.NewContext(totalWidth, totalHeight)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	x := 0
	for _, img := range cardImages {
		dc.DrawImageAnchored(img, x+cardWidth/2, cardHeight/2, 0.5, 0.5)
		x += cardWidth + spacing
	}

	outfile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outfile.Close()

	return png.Encode(outfile, dc.Image())
}

// listAvailableCards prints all embedded card names in user-friendly format
func listAvailableCards() error {
	entries, err := cardFS.ReadDir("cards")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		base := strings.TrimSuffix(entry.Name(), ".png")
		parts := strings.Split(base, "_")
		if strings.HasPrefix(base, "major_arcana_") {
			fmt.Println(strings.TrimPrefix(base, "major_arcana_"))
		} else if strings.HasPrefix(base, "minor_arcana_") && len(parts) >= 4 {
			fmt.Printf("%s_of_%s\n", strings.Join(parts[3:], "_"), parts[2])
		}
	}

	return nil
}

func main() {
	var yamlPath string
	var output string
	var listCards bool
	var cardsCSV string

	flag.StringVar(&yamlPath, "yaml", "", "Path to YAML file describing cards and output")
	flag.StringVar(&output, "o", "spread.png", "Output PNG filename")
	flag.StringVar(&cardsCSV, "c", "", "Comma-separated list of cards (e.g. strength,!hermit,5_of_swords)")
	flag.BoolVar(&listCards, "list", false, "List all available card names")
	flag.Parse()

	// Handle -list flag
	if listCards {
		if err := listAvailableCards(); err != nil {
			log.Fatalf("Error listing cards: %v", err)
		}
		return
	}

	var cards []Card

	// Load from YAML if provided
	if yamlPath != "" {
		cfg, err := parseYAML(yamlPath)
		if err != nil {
			log.Fatalf("Failed to read YAML: %v", err)
		}
		for _, name := range cfg.Cards {
			reversed := strings.HasPrefix(name, "!")
			name = strings.TrimPrefix(name, "!")
			cards = append(cards, Card{Name: name, Reversed: reversed})
		}
		output = cfg.Output
	} else if cardsCSV != "" {
		// Load from comma-separated flag
		for _, item := range strings.Split(cardsCSV, ",") {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			reversed := strings.HasPrefix(item, "!")
			name := strings.TrimPrefix(item, "!")
			cards = append(cards, Card{Name: name, Reversed: reversed})
		}
	}

	// If nothing was parsed, show usage
	if len(cards) == 0 {
		fmt.Println("Usage:")
		fmt.Println("  gtarot -c strength,!hermit,5_of_swords -o spread.png")
		fmt.Println("  gtarot -yaml spread.yaml")
		fmt.Println("  gtarot -list")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := drawSpread(cards, output)
	if err != nil {
		log.Fatalf("Failed to draw spread: %v", err)
	}
	fmt.Printf("Tarot spread saved to %s\n", output)
}
