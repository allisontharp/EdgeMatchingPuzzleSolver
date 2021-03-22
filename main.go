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
	}
	return newTile
}

func getTileIDs(tileList []Tile) (out []int) {
	for _, tile := range tileList {
		out = append(out, tile.ID)
	}
	return out
}

func getAvailableTiles(tileList []Tile) (diff []int) {
	usedTiles := getTileIDs(tileList)
	fmt.Printf("\tusedTiles:%v\n", usedTiles)

	allTiles := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	m := make(map[int]bool)
	for _, item := range usedTiles {
		m[item] = true
	}

	for _, item := range allTiles {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	fmt.Printf("diff:%v\n", diff)
	return diff
}

func getTileByID(tileID int) (out Tile) {
	for _, tile := range tiles {
		if tile.ID == tileID {
			out = tile
		}
	}
	return out
}

func getTilesByIDs(tileIDs []int) (out []Tile) {
	for id := range tileIDs {
		tile := getTileByID(id)
		if tile.ID != 0 {
			out = append(out, getTileByID(id))
		}
	}
	return out
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

	tile := structs.Map(currentTile)

	isTileMatch := true
	for _, sideToMatch := range sidesToMatch {
		currentTileSide := fmt.Sprintf("%v", tile[sideToMatch.sideToMatchOnTile])

		testTile := structs.Map(tileArray[sideToMatch.tileToMatch])
		testTileSide := fmt.Sprintf("%v", testTile[sideToMatch.sideToMatchOnMatchedTile])

		isMatch := checkForEdgeMatch(currentTileSide, testTileSide)
		isTileMatch = isTileMatch && isMatch
	}

	if !isTileMatch {
		currentTile = rotateTile(currentTile)
		currentTileOut, err := checkForTileMatch(currentTile, tileNumber, tileArray, rotationNumber+1)
		if err == nil {
			return currentTileOut, nil
		}
	} else {
		return currentTile, nil
	}

	return currentTile, errors.New("Tile does not work")
}

func removeTile(tiles []Tile, index int) []Tile {
	return append(tiles[:index], tiles[index+1:]...)
}

func removeItemFromIDList(ids []int, index int) []int {
	fmt.Printf("remove Item before:%v\n", ids)
	out := append(ids[:index], ids[index+1:]...)
	fmt.Printf("remove item after:%v\n", out)
	return out
}

func tryAllAvailableTiles(positionAvailableTileIDs []int, tileArray []Tile, positionOfTile int) ([]int, []Tile, int, error) {
	fmt.Printf("\ttryAllAvailableTiles:\n"+
		"\t\tpositionAvailableTileIds:%v\n"+
		"\t\tpositionOfTile:%v\n", positionAvailableTileIDs, positionOfTile)
	positionAvailableTiles := getTilesByIDs(positionAvailableTileIDs)
	for indexOfTileToCheck, tile := range positionAvailableTiles {
		fmt.Printf("\tPosition: %v\tTest Tile %v\n", positionOfTile, indexOfTileToCheck)
		tile, err := checkForTileMatch(tile, positionOfTile, tileArray, 0)
		fmt.Printf("positionAvailableTileIDs before: %v\n", positionAvailableTileIDs)
		positionAvailableTileIDs = removeItemFromIDList(positionAvailableTileIDs, 0)
		fmt.Printf("positionAvailableTileIDs after: %v\n", positionAvailableTileIDs)
		if err == nil {
			tileArray = append(tileArray, tile)
			fmt.Printf("\tChosenTileIndex:%v\n\tChosen TileID: %v\n", indexOfTileToCheck, tile.ID)
			fmt.Printf("positionAvailableTileIds: %v\n", positionAvailableTileIDs)
			return positionAvailableTileIDs, tileArray, indexOfTileToCheck, nil
		}
	}
	fmt.Printf("no tiles match for slot %v (current tiles: %v)\n", positionOfTile, tileArray)
	return positionAvailableTileIDs, tileArray, -1, errors.New("No valid tile for position.")
}

func main() {
	pprintTile(tiles[0])

	availableTiles := tiles

	// This holds each of the 9 tiles in order
	tileArray := []Tile{}
	position := 0
	firstTileIndex := 0
	isPrevious := false

	attempt := 0
	maxAttempt := 10

	positionAvailableTiles := []int{}
	for position < 9 {
		attempt += 1
		if attempt > maxAttempt {
			break
		}
		fmt.Printf("\n\nPosition: %v\n", position)
		fmt.Printf("\tisPrevious: %v\n", isPrevious)
		fmt.Printf("\tCurrent Tile IDs: %v\n", getTileIDs(tileArray))
		if position == 0 {
			potentialTilesByPosition = append(potentialTilesByPosition, getAvailableTiles(tileArray))
			tileArray = append(tileArray, availableTiles[firstTileIndex])
			availableTiles = removeTile(availableTiles, firstTileIndex)
			fmt.Printf("\ttileArray: %v\n", tileArray)
			fmt.Printf("\tpotentialTilesByPosition: %v\n", potentialTilesByPosition)
		} else {
			// Iterate over each tile placement.  We really start with 1 bc tile 0 is a little irrelevant
			if !isPrevious {
				positionAvailableTiles = getAvailableTiles(tileArray)
				fmt.Printf("\tpositionAvailableTiles:%v\n", positionAvailableTiles)
			} else {
				fmt.Printf("\tisPrevious mess:\n")
				fmt.Printf("\t\tpositionAvailableTiles length: %v\n"+
					"\t\tpositionAvailableTiles: %v\n", len(potentialTilesByPosition), potentialTilesByPosition[position])
				positionAvailableTiles = potentialTilesByPosition[position]
			}
			isPrevious = false
			fmt.Printf("\tAvailable Tile IDs: %v\n", positionAvailableTiles)
			positionAvailableTiles, tileArrayOut, indexOfChosenTile, err := tryAllAvailableTiles(positionAvailableTiles, tileArray, position)
			fmt.Printf("main loop positionAvailableTiles: %v\n", positionAvailableTiles)
			tileArray = tileArrayOut
			if err != nil {
				fmt.Printf("\tpositionAvailableTiles left: %v \n", len(positionAvailableTiles))
				availableTiles = append(availableTiles, tileArray[position-1])
				tileArray = removeTile(tileArray, position-1)
				position = position - 2
				isPrevious = true
			} else if indexOfChosenTile >= 0 {
				fmt.Printf("positionAvailableTiles being added:%v\n", positionAvailableTiles)
				potentialTilesByPosition = append(potentialTilesByPosition, positionAvailableTiles)
				availableTiles = removeTile(availableTiles, indexOfChosenTile)
			}
		}
		position += 1
	}

	fmt.Printf("\n\n\n# Availble: %v\nAvailableTiles: %v\n", len(availableTiles), availableTiles)
	fmt.Printf("TileArray: %v\n", tileArray)
}

// Add the first tile to the array and remove it from the available tiles
