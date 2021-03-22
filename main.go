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
var allTileIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
var printDetails = true

type TileSide struct {
	Direction   string
	Description string
}

type Tile struct {
	North string
	East  string
	South string
	West  string
	ID    int
}

var tiles = []Tile{
	{
		North: "RBeetle",
		East:  "LDragon",
		South: "RDragon",
		West:  "RGrasshopper",
		ID:    1,
	},
	{
		North: "LBeetle",
		East:  "RGrasshopper",
		South: "LAnt",
		West:  "RDragon",
		ID:    2,
	},
	{
		North: "RAnt",
		East:  "RBeetle",
		South: "LBeetle",
		West:  "LGrasshopper",
		ID:    3,
	},
	{
		North: "LDragon",
		East:  "RAnt",
		South: "RGrasshopper",
		West:  "LAnt",
		ID:    4,
	},
	{
		North: "RAnt",
		East:  "LDragon",
		South: "RGrasshopper",
		West:  "LBeetle",
		ID:    5,
	},
	{
		North: "RBeetle",
		East:  "LAnt",
		South: "RDragon",
		West:  "LGrasshopper",
		ID:    6,
	},
	{
		North: "LAnt",
		East:  "RBeetle",
		South: "LGrasshopper",
		West:  "LDragon",
		ID:    7,
	},
	{
		North: "LBeetle",
		East:  "RGrasshopper",
		South: "LDragon",
		West:  "RAnt",
		ID:    8,
	},
	{
		North: "RGrasshopper",
		East:  "LDragon",
		South: "RAnt",
		West:  "RBeetle",
		ID:    9,
	},
}

// when you move on to a new position, you can add back in all of the tiles that are currently available to this.
var potentialTilesByPosition = [][]int{}

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
		ID:    tile.ID,
	}
	return newTile
}

func getAvailableTiles(placedTiles []int, availableTileIDs []int) (diff []int) {
	print(fmt.Sprintf("\t\tplacedTiles: %v | available tile ids: %v\n", placedTiles, availableTileIDs))

	m := make(map[int]bool)
	for _, item := range placedTiles {
		m[item] = true
	}
	for _, item := range availableTileIDs {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	return diff
}

func splitTileSide(side string) TileSide {
	tileSide := TileSide{
		Direction:   side[0:1],
		Description: side[1:],
	}
	return tileSide
}

func getTileByID(tileID int) (out Tile) {
	for _, tile := range tiles {
		if tile.ID == tileID {
			out = tile
		}
	}
	return out
}

// checks an individual edge
func checkEdgeMatch(currentTileSide string, testTileSide string) bool {
	current := splitTileSide(currentTileSide)
	test := splitTileSide(testTileSide)
	// could do this with ORs and probably other ways, but this is easiest for now
	if current.Description == test.Description {
		if current.Direction == "R" && test.Direction == "L" {
			return true
		} else if current.Direction == "L" && test.Direction == "R" {
			return true
		}
	}
	return false
}

// checks all edges of test tile
func checkTileMatch(currentTile Tile, position int, placedTiles []int, rotationNumber int) (Tile, error) {
	print(fmt.Sprintf("\t\tRotation Number: %v\n", rotationNumber))
	if rotationNumber > 3 {
		return currentTile, errors.New("Invalid Tile")
	}
	// first tile doesnt match on anything
	if position == 0 {
		return currentTile, nil
	}
	sidesToMatch := sidesToMatchArray[position]
	tile := structs.Map(currentTile)
	isTileMatch := true

	for _, sideToMatch := range sidesToMatch {
		currentTileSide := fmt.Sprintf("%v", tile[sideToMatch.sideToMatchOnTile])
		tileToMatch := getTileByID(placedTiles[sideToMatch.tileToMatch])
		testTile := structs.Map(tileToMatch)
		testTileSide := fmt.Sprintf("%v", testTile[sideToMatch.sideToMatchOnMatchedTile])

		isMatch := checkEdgeMatch(currentTileSide, testTileSide)
		isTileMatch = isTileMatch && isMatch
	}

	if !isTileMatch {
		currentTile = rotateTile(currentTile)
		currentTileOut, err := checkTileMatch(currentTile, position, placedTiles, rotationNumber+1)
		if err == nil {
			return currentTileOut, nil
		}
	} else {
		return currentTile, nil
	}

	return currentTile, errors.New("Tile does not work..")
}

func removePlacedTile(placedTiles []int, index int) []int {
	return append(placedTiles[:index], placedTiles[index+1:]...)
}

func print(printText string) {
	if printDetails {
		fmt.Printf(printText)
	}
}

func main() {
	pprintTile(tiles[0])

	placedTiles := []int{}
	position := 0
	firstTileIndex := 0
	isPrevious := false

	availableTilesByPosition := [][]int{}
	availableTiles := []int{}

	attemptNumber := 0
	maxAttempts := 78

	placedTiles = append(placedTiles, tiles[firstTileIndex].ID)
	availableTiles = getAvailableTiles(placedTiles, allTileIDs)
	availableTilesByPosition = append(availableTilesByPosition, availableTiles)
	print(fmt.Sprintf("\tcurrentTile: %v\n", tiles[firstTileIndex]))

	print(fmt.Sprintf("\tavailableTiles: %v\n", availableTiles))

	position += 1

	for position < 9 {
		attemptNumber += 1
		runTileCheck := true
		if attemptNumber > maxAttempts {
			break
		}
		fmt.Printf("Position: %v\tAttempt: %v\n", position, attemptNumber)
		fmt.Printf("\tisPrevious: %v\n", isPrevious)
		{
			if !isPrevious {
				availableTiles = getAvailableTiles(placedTiles, allTileIDs)
				if len(availableTiles) == 0 {
					// none of the available tiles worked
					// need to use a different tile in the previous position (so, position=position-2)
					position = position - 2
					// need to also remove the placed tile
					placedTiles = removePlacedTile(placedTiles, len(placedTiles)-1)
					print(fmt.Sprintf("\t--->Placed Tiles: %v<---\n", placedTiles))
					// need to mark that this is a retry
					isPrevious = true
					runTileCheck = false
				}
			} else {

				print(fmt.Sprintf("\tTotal tiles available in this position currently: %v\n", len(availableTilesByPosition[position])))
				if len(availableTilesByPosition[position]) == 0 {
					// none of the available tiles worked
					// need to use a different tile in the previous position (so, position=position-2)
					position = position - 2
					// need to also remove the placed tile
					placedTiles = removePlacedTile(placedTiles, len(placedTiles)-1)
					print(fmt.Sprintf("\t--->Placed Tiles After Removal: %v<---\n", placedTiles))
					// need to mark that this is a retry
					isPrevious = true
					runTileCheck = false
				} else {
					availableTilesByPosition[position] = removePlacedTile(availableTilesByPosition[position], 0)
				}
				if len(availableTilesByPosition[position]) == 0 {
					if !runTileCheck {
						position = position - 1
					} else {
						position = position - 2
					}
					print(fmt.Sprintf("\tPosition: %v\n", position))

					// none of the available tiles worked
					// need to use a different tile in the previous position (so, position=position-2)
					// position = position - 2
					// need to also remove the placed tile
					placedTiles = removePlacedTile(placedTiles, len(placedTiles)-1)
					fmt.Printf("\t--->Placed Tiles: %v<---\n", placedTiles)
					// need to mark that this is a retry
					isPrevious = true
					runTileCheck = false
				} else {
					availableTiles = getAvailableTiles(placedTiles, availableTilesByPosition[position])
				}

			}
			if runTileCheck {

				print(fmt.Sprintf("\tavailableTiles before pick: %v\n", availableTiles))
				for testTileNumber, testTile := range availableTiles { // try each available tile until it finds one that matches
					currentTile := getTileByID(testTile)
					print(fmt.Sprintf("\tTesting Tile: %v\n", currentTile))
					currentTile, err := checkTileMatch(currentTile, position, placedTiles, 0)
					if err == nil { // tile worked, move on to next position
						print(fmt.Sprintf("\tchosen tile: %v\n", currentTile))
						placedTiles = append(placedTiles, currentTile.ID)
						availableTiles = getAvailableTiles(placedTiles, availableTiles)
						print(fmt.Sprintf("\tAvailableTilesByPosition: %v | position: %v\n", len(availableTilesByPosition), position))
						if len(availableTilesByPosition) <= position {
							availableTilesByPosition = append(availableTilesByPosition, availableTiles)
						} else {
							availableTilesByPosition[position] = availableTiles
						}

						print(fmt.Sprintf("\tavailableTiles after pick: %v\n", availableTiles))
						print(fmt.Sprintf("\tPosition's available tiles: %v\n", availableTilesByPosition[position]))
						isPrevious = false
						fmt.Printf("\t--->Placed Tiles: %v<---\n", placedTiles)
						break
					} else { // tile did not work
						print(fmt.Sprintf("\t---------\n\tERROR: %v\n\t---------\n", err))
					}
					if testTileNumber == len(availableTiles)-1 {
						// none of the available tiles worked
						// need to use a different tile in the previous position (so, position=position-2)
						position = position - 2
						// need to also remove the placed tile
						placedTiles = removePlacedTile(placedTiles, len(placedTiles)-1)
						fmt.Printf("\t--->Placed Tiles: %v<---\n", placedTiles)
						// need to mark that this is a retry
						isPrevious = true
					}
				}
			}

		}
		position += 1
	}

	fmt.Printf("\n\n-----------------")
	fmt.Printf("\nPlacedTiles: %v", placedTiles)
	fmt.Printf("\n\n-----------------")

}
