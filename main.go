package main

import (
	"fmt"
	"strings"
)

var maxWidth int = 30

type TileSide struct {
	Direction   string
	Description string
}

type Tile struct {
	North       string
	East        string
	South       string
	West        string
	Orientation string
}

var tiles = []Tile{
	{
		North: "RBeetle",
		East:  "LDragon",
		South: "RDragon",
		West:  "RGrasshopper",
	},
	{
		North: "LBeetle",
		East:  "RGrasshopper",
		South: "LAnt",
		West:  "RDragon",
	},
	{
		North: "RAnt",
		East:  "RBeetle",
		South: "LBettle",
		West:  "LGrasshopper",
	},
	{
		North: "LDragon",
		East:  "RAnt",
		South: "RGrasshopper",
		West:  "LAnt",
	},
	{
		North: "RAnt",
		East:  "LDragon",
		South: "RGrasshopper",
		West:  "LBettle",
	},
	{
		North: "RBeetle",
		East:  "LAnt",
		South: "RDragon",
		West:  "LGrasshopper",
	},
	{
		North: "LAnt",
		East:  "RBeetle",
		South: "LGrasshopper",
		West:  "LDragon",
	},
	{
		North: "LBeetle",
		East:  "RGrasshopper",
		South: "LDragon",
		West:  "RAnt",
	},
	{
		North: "RGrasshopper",
		East:  "LDragon",
		South: "RAnt",
		West:  "RBeetle",
	},
}

func pprintTile(tile Tile) {
	topSpaces := strings.Repeat(" ", (maxWidth-1-len(tile.North))/2)
	bottomSpaces := strings.Repeat(" ", (maxWidth-1-len(tile.South))/2)
	middleSpaces := strings.Repeat(" ", (maxWidth - 1 - len(tile.East) - len(tile.West)))
	fmt.Printf("%v\n", strings.Repeat("-", maxWidth))
	fmt.Printf("|%v%v%v|\n", topSpaces, tile.North, topSpaces)
	fmt.Printf("|%v%v%v|\n", tile.West, middleSpaces, tile.East)
	fmt.Printf("|%v%v%v|\n", bottomSpaces, tile.South, bottomSpaces)
	fmt.Printf("%v\n", strings.Repeat("-", maxWidth))
}

func rotateTile(tile Tile) Tile {
	newTile := Tile{
		North: tile.West,
		East:  tile.North,
		South: tile.East,
		West:  tile.South,
	}
	return newTile
}

func splitTileSide(side string) TileSide {
	tileSide := TileSide{
		Direction:   side[0:1],
		Description: side[1:],
	}
	return tileSide
}

func checkForMatch(currentTileSide string, testTileSide string) bool {
	current := splitTileSide(currentTileSide)
	test := splitTileSide(testTileSide)

	fmt.Println(currentTileSide, testTileSide)

	// could do this with ORs and probably other ways, but again, this is easiest for now.
	if current.Description == test.Description {
		if current.Direction == "R" && test.Direction == "L" {
			return true
		} else if current.Direction == "L" && test.Direction == "R" {
			return true
		}
	}

	return false
}

func main() {
	pprintTile(tiles[0])
	t := rotateTile(tiles[0])
	pprintTile(t)

	fmt.Println(checkForMatch(tiles[0].North, tiles[1].East))
}
