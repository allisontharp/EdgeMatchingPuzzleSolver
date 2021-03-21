package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/structs"
)

/*
Tile numbers will be:
0 1 2
3 4 5
6 7 8
*/

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
		South: "LBeetle",
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
		West:  "LBeetle",
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

// There must be a better way to do this...
type sideToMatch struct {
	tileToMatch              int
	sideToMatchOnTile        string // side to match on the current tile
	sideToMatchOnMatchedTile string // side to match on the tile you are matching to
}

var sidesToMatchArray = [][]sideToMatch{
	// Tile 0 has no matches
	{{}},
	// Tile 1 matches tile 0 on the W side
	{{tileToMatch: 0, sideToMatchOnTile: "West", sideToMatchOnMatchedTile: "East"}},
	// Tile 2 Matches tile 1 on the W side
	{{tileToMatch: 1, sideToMatchOnTile: "West", sideToMatchOnMatchedTile: "East"}},
	// Tile 3 matches tile 0 on the N side
	{{tileToMatch: 0, sideToMatchOnTile: "North", sideToMatchOnMatchedTile: "South"}},
	// tile 4 matches tile 3 on the W side and tile 1 on the N side
	{{tileToMatch: 3, sideToMatchOnTile: "West", sideToMatchOnMatchedTile: "East"},
		{tileToMatch: 1, sideToMatchOnTile: "North", sideToMatchOnMatchedTile: "South"}},
	// Tile 5 matches tile 4 on W side and tile 2 on N side
	{{tileToMatch: 4, sideToMatchOnTile: "West", sideToMatchOnMatchedTile: "East"},
		{tileToMatch: 2, sideToMatchOnTile: "North", sideToMatchOnMatchedTile: "South"}},
	// Tile 6 matches tile 3 on the N side
	{{tileToMatch: 3, sideToMatchOnTile: "North", sideToMatchOnMatchedTile: "South"}},
	// Tile 7 matches tile 4 on N sideand tile 6 on the W side
	{{tileToMatch: 6, sideToMatchOnTile: "West", sideToMatchOnMatchedTile: "East"},
		{tileToMatch: 4, sideToMatchOnTile: "North", sideToMatchOnMatchedTile: "South"}},
	// Tile 8 matches tile 7 on W side and tile 4 on N side
	{{tileToMatch: 7, sideToMatchOnTile: "West", sideToMatchOnMatchedTile: "East"},
		{tileToMatch: 5, sideToMatchOnTile: "North", sideToMatchOnMatchedTile: "South"}},
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

func checkForEdgeMatch(currentTileSide string, testTileSide string) bool {
	current := splitTileSide(currentTileSide)
	test := splitTileSide(testTileSide)

	fmt.Printf("\tTile1: %v,\t\tTile2: %v\n", currentTileSide, testTileSide)
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

func checkForTileMatch(currentTile Tile, tileNumber int, tileArray []Tile, rotationNumber int) (Tile, error) {
	if rotationNumber > 3 {
		return currentTile, errors.New("Invalid Tile")
	}

	sidesToMatch := sidesToMatchArray[tileNumber]
	if tileNumber == 3 {
		fmt.Println(sidesToMatch)
	}

	tile := structs.Map(currentTile)

	isTileMatch := true

	fmt.Println(sidesToMatch)
	for _, sideToMatch := range sidesToMatch {
		currentTileSide := fmt.Sprintf("%v", tile[sideToMatch.sideToMatchOnTile])
		testTile := structs.Map(tileArray[sideToMatch.tileToMatch])
		testTileSide := fmt.Sprintf("%v", testTile[sideToMatch.sideToMatchOnMatchedTile])

		isMatch := checkForEdgeMatch(currentTileSide, testTileSide)
		isTileMatch = isTileMatch && isMatch
	}

	if !isTileMatch {
		currentTile = rotateTile(currentTile)
		checkForTileMatch(currentTile, tileNumber, tileArray, rotationNumber+1)
	} else {
		return currentTile, nil
	}

	return currentTile, errors.New("Invalid Tile")
}

func removeTile(tiles []Tile, index int) []Tile {
	return append(tiles[:index], tiles[index+1:]...)
}

func iterateOverArray() {

}

func main() {
	pprintTile(tiles[0])
	t := rotateTile(tiles[0])
	pprintTile(t)

	availableTiles := tiles

	// This holds each of the 9 tiles in order
	tileArray := []Tile{}

	// Add the first tile to the array and remove it from the available tiles
	tileArray = append(tileArray, tiles[0])
	availableTiles = removeTile(availableTiles, 0)

	// Iterate over each tile placement.  We really start with 1 bc tile 0 is a little irrelevant
	for i := 1; i <= 8; i++ {
		for indexOfTileToCheck := range availableTiles {
			numAvailableTilesLeft := len(availableTiles)
			fmt.Printf("\n\nPosition %v Test Tile %v\n", i, indexOfTileToCheck)
			tile, err := checkForTileMatch(availableTiles[indexOfTileToCheck], i, tileArray, 0)
			if err == nil {
				tileArray = append(tileArray, tile)
				availableTiles = removeTile(availableTiles, indexOfTileToCheck)
				fmt.Printf("\tChosen Tile: %v", tile)
				break
			} else if indexOfTileToCheck == numAvailableTilesLeft-1 {
				fmt.Printf("no tiles match for slot %v (current tiles: %v)\n\n", i, tileArray)
				return
			}
		}
	}

	fmt.Println(len(availableTiles))
}
