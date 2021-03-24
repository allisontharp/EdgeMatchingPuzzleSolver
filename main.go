/*
	TODO: Handle cases where a single tile has multiple valid placements.
		in the availableTilesByPosition, maybe have each position be [[tileids], rotationNumber] ?
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

/*
Tile numbers will be:
0 1 2
3 4 5
6 7 8
*/

var maxPPrintWidth int = 30

// var allTileIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
var printDetails bool
var width = 3
var height = 3

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

// var tiles = []Tile{
// 	{
// 		North: "LGoat",
// 		East:  "LBeetle",
// 		South: "LAnt",
// 		West:  "LGoat",
// 		ID:    1,
// 	}, {
// 		North: "LGoat",
// 		East:  "LGoat",
// 		South: "RButterfly",
// 		West:  "RBeetle",
// 		ID:    2,
// 	}, {
// 		North: "RAnt",
// 		East:  "RGrasshopper",
// 		South: "LGoat",
// 		West:  "LGoat",
// 		ID:    3,
// 	}, {
// 		North: "LButterfly",
// 		East:  "LGoat",
// 		South: "LGoat",
// 		West:  "LGrasshopper",
// 		ID:    4,
// 	},
// }

var tiles = []Tile{
	{
		North: "LGoat",
		East:  "LHouse",
		South: "RMouse",
		West:  "LTree",
		ID:    1,
	},
	{
		North: "LTree",
		East:  "LHouse",
		South: "LMouse",
		West:  "RGoat",
		ID:    2,
	},
	{
		North: "RGoat",
		East:  "LHouse",
		South: "RTree",
		West:  "RMouse",
		ID:    3,
	},
	{
		North: "RHouse",
		East:  "LMouse",
		South: "RTree",
		West:  "RGoat",
		ID:    4,
	},
	{
		North: "LTree",
		East:  "LMouse",
		South: "LHouse",
		West:  "RGoat",
		ID:    5,
	},
	{
		North: "LTree",
		East:  "LGoat",
		South: "LHouse",
		West:  "RMouse",
		ID:    6,
	},
	{
		North: "RTree",
		East:  "RGoat",
		South: "RHouse",
		West:  "LMouse",
		ID:    7,
	},
	{
		North: "RHouse",
		East:  "RTree",
		South: "LGoat",
		West:  "RMouse",
		ID:    8,
	},
	{
		North: "LTree",
		East:  "LMouse",
		South: "LHouse",
		West:  "RMouse",
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

func pprintTile(tile Tile) {
	descLength := 4
	topSpaces := strings.Repeat(" ", (maxPPrintWidth-descLength+1)/2)
	bottomSpaces := strings.Repeat(" ", (maxPPrintWidth-descLength+1)/2)
	middleSpaces := strings.Repeat(" ", (maxPPrintWidth - (2 * descLength)))
	fmt.Printf(" %v\n", strings.Repeat("-", maxPPrintWidth))
	fmt.Printf("|%v%v%v|\n", topSpaces, tile.North[0:descLength], topSpaces)
	fmt.Printf("|%v%v%v|\n", tile.West[0:descLength], middleSpaces, tile.East[0:descLength])
	fmt.Printf("|%v%v%v|\n", bottomSpaces, tile.South[:descLength], bottomSpaces)
	fmt.Printf(" %v\n", strings.Repeat("-", maxPPrintWidth))
}

func pprintTiles(tiles []Tile) {
	descLength := 4
	topSpaces := strings.Repeat(" ", (maxPPrintWidth-descLength+1)/2)
	bottomSpaces := strings.Repeat(" ", (maxPPrintWidth-descLength+1)/2)
	middleSpaces := strings.Repeat(" ", (maxPPrintWidth - (2 * descLength)))

	// print the top bar
	x := 1
	y := 1
	i := 0
	for y <= height {
		topRow := ""
		middleRow := ""
		bottomRow := ""
		for x <= width {
			tile := tiles[i]
			topRow += fmt.Sprintf("|%v%v%v|", topSpaces, tile.North[0:descLength], topSpaces)
			middleRow += fmt.Sprintf("|%v%v%v|", tile.West[0:descLength], middleSpaces, tile.East[0:descLength])
			bottomRow += fmt.Sprintf("|%v%v%v|", bottomSpaces, tile.South[:descLength], bottomSpaces)
			x += 1
			i += 1
		}
		fmt.Printf(" %v\n", strings.Repeat("-", maxPPrintWidth*width+width))
		fmt.Printf("%v\n", topRow)
		fmt.Printf("%v\n", middleRow)
		fmt.Printf("%v\n", bottomRow)
		fmt.Printf(" %v\n", strings.Repeat("-", maxPPrintWidth*width+width))
		x = 1
		y += 1
	}
}

func rotateTile(tile Tile) Tile {
	newTile := Tile{
		North: tile.West,
		East:  tile.North,
		South: tile.East,
		West:  tile.South,
		ID:    tile.ID,
	}
	// save the rotation of the tile
	for index, t := range tiles {
		if t.ID == tile.ID {
			tiles[index] = newTile
		}
	}
	return newTile
}

func getAvailableTiles(placedTiles []Tile, availableTileIDs []int) (diff []int) {
	m := make(map[int]bool)
	for _, item := range placedTiles {
		m[item.ID] = true
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

func getTileIDs(tiles []Tile) (out []int) {
	for _, tile := range tiles {
		out = append(out, tile.ID)
	}
	return out
}

//Gets x,y position of tile by position number
func getTileCoordinates(position int, width int) (x int, y int) {
	return position % width, position / width
}

func getSidesToMatch(position int, width int) []sideToMatch {
	x, y := getTileCoordinates(position, width)
	out := []sideToMatch{}
	// tiles on left edge don't need to match anything in their row
	if x > 0 { // tile is not on left edge
		out = append(out, sideToMatch{
			tileToMatch:              position - 1,
			sideToMatchOnTile:        "West",
			sideToMatchOnMatchedTile: "East",
		})
	}

	// tiles on the top don't need to match anything above them
	if y > 0 { // tile is not on top row
		out = append(out, sideToMatch{
			tileToMatch:              position - width,
			sideToMatchOnTile:        "North",
			sideToMatchOnMatchedTile: "South",
		})
	}

	return out
}

// checks an individual edge
func checkEdgeMatch(currentTileSide string, testTileSide string) bool {
	current := splitTileSide(currentTileSide)
	test := splitTileSide(testTileSide)
	print(fmt.Sprintf("\t\tComparing %v and %v\n", current, test))
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
func checkTileMatch(currentTile Tile, position int, placedTiles []Tile, rotationNumber int) (Tile, error) {
	print(fmt.Sprintf("\t\tRotation Number: %v\n", rotationNumber))
	if rotationNumber > 3 {
		return currentTile, errors.New("Invalid Tile")
	}
	// first tile doesnt match on anything
	if position == 0 {
		fmt.Printf("\n\tPlacing first tile. \n\n")
		return currentTile, nil
	}
	sidesToMatch := getSidesToMatch(position, width)
	print(fmt.Sprintf("\t\tSidesToMatch: %v\n", sidesToMatch))
	tile := structs.Map(currentTile)
	isTileMatch := true

	for _, sideToMatch := range sidesToMatch {
		currentTileSide := fmt.Sprintf("%v", tile[sideToMatch.sideToMatchOnTile])
		tileToMatch := placedTiles[sideToMatch.tileToMatch]
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

func removePlacedTile(placedTiles []Tile, index int) []Tile {
	return append(placedTiles[:index], placedTiles[index+1:]...)
}

func removeAvailableTile(availableTileIDs []int, index int) []int {
	return append(availableTileIDs[:index], availableTileIDs[index+1:]...)
}

func print(printText string) {
	if printDetails {
		fmt.Printf(printText)
	}
}

func placeTile(placedTiles []Tile, tile Tile, availableTilesByPosition [][]int, positionsPossibleTiles []int, position int, isRetry bool) ([]Tile, [][]int) {
	print(fmt.Sprintf("\tPlacing tileID %v into position %v\n", tile.ID, position))
	// place the tile into the placedTiles array
	placedTiles = append(placedTiles, tile)

	// Remove placed tile
	positionsPossibleTiles = removeAvailableTile(positionsPossibleTiles, 0)
	// add or update the array of all available tiles
	if len(availableTilesByPosition) < position+1 || len(availableTilesByPosition) == 0 {
		availableTilesByPosition = append(availableTilesByPosition, positionsPossibleTiles)
	} else {
		availableTilesByPosition[position] = positionsPossibleTiles
	}

	fmt.Printf("\t--->Placed Tiles: %v<---\n", getTileIDs(placedTiles))

	return placedTiles, availableTilesByPosition
}

func main() {
	// Get commandline arguments
	args := os.Args[1:]
	printDetails, _ = strconv.ParseBool(args[0])
	maxAttempts, _ := strconv.Atoi(args[1])
	// firstTileIndex, _ := strconv.Atoi(args[2])
	// get a list of all tile ids
	allTileIDs := []int{}
	for _, tile := range tiles {
		allTileIDs = append(allTileIDs, tile.ID)
	}
	// Set Up Initial Variables
	position := 0                           // position of tile (will calculate coordinates from this later)
	attemptNumber := 0                      // tracks.. number of attempts
	isRetry := false                        // flag to tell if we are retrying a specific position
	placedTiles := []Tile{}                 // array of tiles (in order, can calculate coorindates later)
	availableTilesByPosition := [][][]int{} // array of arrays that holds available tiles for specific position
	poolOfAvailableTiles := []int{}         // tiles for this position
	checkTiles := true                      // flag to tell if you should check tiles (for instance, if there are no more tiles, dont try to check)

	fmt.Printf("# Attempts: %v\n", maxAttempts)
	for position < len(allTileIDs) {
		attemptNumber += 1
		checkTiles = true
		if attemptNumber > maxAttempts {
			return // hit max attempts, so stop
		}
		fmt.Printf("\nPosition: %v\tAttempt: %v/%v\tisRetry: %v\n", position, attemptNumber, maxAttempts, isRetry)
		// Get position's available tile IDs
		if len(availableTilesByPosition) < position || !isRetry { // available tiles will be from the pool of all tiles
			/*
				use pool of all tiles when:
					it is a retry for any position other than 0 (retries for position 0 need to be handled special)
					it is a position that hasnt been added yet to the available tiles (in other words, its the first time we are trying this position)
			*/
			poolOfAvailableTiles = allTileIDs

		} else {
			poolOfAvailableTiles = availableTilesByPosition[position]
		}
		availableTilesForPosition := getAvailableTiles(placedTiles, poolOfAvailableTiles)

		// Try available tiles in this position to try to find a match
		if len(availableTilesForPosition) == 0 {
			isRetry = true
			checkTiles = false
		}

		if checkTiles {
			print(fmt.Sprintf("\tAvailableTiles before pick: %v\n", availableTilesForPosition))
			fmt.Printf("\tAvailableTiles before pick: %v\n", availableTilesForPosition)
			for testTileNumber, testTile := range availableTilesForPosition { // try each available tile until it finds one that matches
				currentTile := getTileByID(testTile)
				print(fmt.Sprintf("\tTesting Tile: %v\n", currentTile))
				currentTile, err := checkTileMatch(currentTile, position, placedTiles, 0)
				if err == nil { // tile worked, move on to the next position
					// place the tile
					placedTiles, availableTilesByPosition = placeTile(placedTiles, currentTile, availableTilesByPosition, availableTilesForPosition, position, isRetry)
					availableTilesForPosition = availableTilesByPosition[position]
					print(fmt.Sprintf("\tPosition's available tiles: %v\n", availableTilesForPosition))
					fmt.Printf("\tPosition's available tiles: %v\n", availableTilesForPosition)
					isRetry = false
					break // break out of the loop of position's available tiles
				}
				if testTileNumber == len(availableTilesForPosition)-1 && position > 0 {
					print(fmt.Sprintf("\tNo tiles worked.  Retrying previous position.  But first..\n"))
					// none of the available tiles for the position worked, so we need to go back one position and try a new tile
					isRetry = true

					// need to remove the placed tile if it wasnt the first one
					if len(placedTiles) > 1 {
						placedTiles = removePlacedTile(placedTiles, position-1)
						fmt.Printf("\t--->Placed Tiles: %v<---\n", getTileIDs(placedTiles))
					}
				}
			}
		}

		if isRetry { // move back a tile to try a new tile
			position -= 1
			fmt.Printf("\tisRetry: %v\n", isRetry)
		} else { // go on to next position
			position += 1
			fmt.Printf("\tisRetry: %v\n", isRetry)
		}
	}

	fmt.Printf("\n-----------")
	fmt.Printf("\nPlaced TileIDs: %v\n", getTileIDs(placedTiles))

	fmt.Printf("Placed Tiles: \n")
	pprintTiles(placedTiles)
}
