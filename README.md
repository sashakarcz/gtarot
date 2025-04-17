# gtarot

gtarot is a command-line tool written in Go that creates PNG images of tarot card spreads.

- Embed your full tarot deck directly in the binary

- Supports reversed cards (e.g., `!hermit`)

- Use CLI or YAML input

- Outputs beautiful side-by-side layouts

## Features

Embed all tarot card images (*.png) into the Go binary using embed.FS


Specify cards directly via `-c` or via `YAML` file

Reversed cards (rotated 180°) supported using !cardname

Generate a horizontal PNG layout of your spread

List all available card names with `-list`

## Requirements

Go 1.16+ (for embed)

PNG tarot card images stored in cards/ directory

## Installation

```bash

git clone https://github.com/sashakarcz/gtarot
cd gtarot
go mod tidy
go build -o gtarot
```

Ifyou wish to add custom cards, make sure your cards/ folder contains images named like:

`major_arcana_fool.png`

`minor_arcana_swords_5.png`

`minor_arcana_cups_queen.png`

## Usage

### List available cards

```
./gtarot -list
```

### Generate a spread with CLI input

```
./gtarot -c strength,'!hermit',5_of_swords -o spread.png

```

NOTE: `!hermit` → reversed (rotated 180°)

### Generate a spread using a YAML file

```
cat spread.yaml

cards:
  - strength
  - '!hermit'
  - 5_of_swords
output: spread.png
```

Then run:

```
./gtarot -yaml spread.yaml
```

### Command Line Flags

| Flag             | Description                                                           | Example Usage                                  |
|------------------|-----------------------------------------------------------------------|------------------------------------------------|
| `-c`             | Comma-separated list of card names. Use `!` for reversed cards.       | `-c strength,!hermit,5_of_swords`              |
| `-yaml`          | Path to a YAML file specifying cards and output filename.             | `-yaml spread.yaml`                            |
| `-o`             | Output PNG filename (default: `spread.png`).                          | `-o custom.png`                                |
| `-list`          | List all available embedded tarot card names.                         | `-list`                                        |
| `-random`        | Creates a random draw of cards with number of cards entered.          | `-random 5`                                    |
| `-allow-reverse` | Used with `-random`; allows for cards to be reversed. Default: False. | `-random 5 -allow-reverse`                     |


### Example Output

![Example tarot spread](output.png)


## License

MIT – do as thou wilt!

## Roadmap Ideas



PRs welcome!

## Sources

Thank you to the [Internet Archive](https://archive.org/details/rider-waite-tarot) for the public domain images of the Rider-Waite deck.

