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
	"time"

	"github.com/fatih/structs"
)

/*
Tile numbers will be:
0 1 2
3 4 5
6 7 8
*/

var maxPPrintWidth int = 30
var printDetails = 0
var width = 3
var height = 3
var firstTileRotationNumber int = 0

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

var tiles2 = []Tile{
	{
		East:  "LGoat",
		South: "LBeetle",
		West:  "LAnt",
		North: "LGoat",
		ID:    1,
	}, {
		North: "LGoat",
		East:  "LGoat",
		South: "RButterfly",
		West:  "RBeetle",
		ID:    2,
	}, {
		North: "RAnt",
		East:  "RGrasshopper",
		South: "LGoat",
		West:  "LGoat",
		ID:    3,
	}, {
		North: "LButterfly",
		East:  "LGoat",
		South: "LGoat",
		West:  "LGrasshopper",
		ID:    4,
	},
}

var tiles = []Tile{
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
	// I know that these are not valid start tiles, so i moved them to the end.
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
		North: "LTree",
		East:  "LGoat",
		South: "LHouse",
		West:  "RMouse",
		ID:    6,
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

	x := 1
	y := 1
	i := 0
	totalTiles := len(tiles)
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
			if i >= totalTiles {
				break
			}
		}
		fmt.Printf(" %v\n", strings.Repeat("-", maxPPrintWidth*width+width))
		fmt.Printf("%v\n", topRow)
		fmt.Printf("%v\n", middleRow)
		fmt.Printf("%v\n", bottomRow)
		fmt.Printf(" %v\n", strings.Repeat("-", maxPPrintWidth*width+width))
		x = 1
		y += 1
		if i >= totalTiles {
			break
		}
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
	print(fmt.Sprintf("\t\t\t\tComparing %v and %v\n", current, test), 4) // really dont need this noise
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
	print(fmt.Sprintf("\t\tRotation Number: %v\n", rotationNumber), 3)
	if rotationNumber > 4 {
		return currentTile, errors.New("Invalid Tile")
	}
	// first tile doesnt match on anything
	if position == 0 {
		fmt.Printf("\n\tPlacing first tile. \n\n")
		return currentTile, nil
	}
	sidesToMatch := getSidesToMatch(position, width)
	print(fmt.Sprintf("\t\tSidesToMatch: %v\n", sidesToMatch), 4)
	tile := structs.Map(currentTile)
	isTileMatch := true

	for _, sideToMatch := range sidesToMatch {
		print(fmt.Sprintf("\t\t\tSideToMatch: %v\n", sideToMatch), 4)
		currentTileSide := fmt.Sprintf("%v", tile[sideToMatch.sideToMatchOnTile])
		tileToMatch := placedTiles[sideToMatch.tileToMatch]
		testTile := structs.Map(tileToMatch)
		testTileSide := fmt.Sprintf("%v", testTile[sideToMatch.sideToMatchOnMatchedTile])

		isMatch := checkEdgeMatch(currentTileSide, testTileSide)
		isTileMatch = isTileMatch && isMatch
		print(fmt.Sprintf("\t\t\tisMatch: %v \tisTileMatch: %v\n", isMatch, isTileMatch), 4)
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

func removeAvailableTile(availableTileIDs []int, tileID int) []int {
	out := []int{}
	for _, id := range availableTileIDs {
		if id != tileID {
			out = append(out, id)
		}
	}
	return out
}

func print(printText string, printLevel int) {
	/*
		levels: 0 - summary
				1 - basic
				2 - detail
	*/
	if printLevel <= printDetails {
		fmt.Printf(printText)
	}
}

func placeTile(placedTiles []Tile, tile Tile, availableTilesByPosition [][]int, positionsPossibleTiles []int, position int, isRetry bool) ([]Tile, [][]int) {
	print(fmt.Sprintf("\tPlacing tileID %v into position %v\n", tile.ID, position), 2)
	// place the tile into the placedTiles array
	placedTiles = append(placedTiles, tile)

	// Remove placed tile
	positionsPossibleTiles = removeAvailableTile(positionsPossibleTiles, tile.ID)
	// add or update the array of all available tiles
	if len(availableTilesByPosition) < position+1 || len(availableTilesByPosition) == 0 {
		availableTilesByPosition = append(availableTilesByPosition, positionsPossibleTiles)
	} else {
		availableTilesByPosition[position] = positionsPossibleTiles
	}

	print(fmt.Sprintf("\tpositionPossibleTiles: %v\n", positionsPossibleTiles), 2)
	print(fmt.Sprintf("\t--->Placed TilesA: %v<---\n", getTileIDs(placedTiles)), 0)

	return placedTiles, availableTilesByPosition
}

func rotateFirstTile(tile Tile) (placedTiles []Tile) {
	print("\tRotating first tile..\n", 1)
	print(fmt.Sprintf("\t\tBefore rotation: %v", tile), 2)
	tile = rotateTile(tile)
	placedTiles = append(placedTiles, tile)
	print(fmt.Sprintf("\t\tAfter rotation: %v", tile), 2)
	return placedTiles
}

func checkTilesForPosition(position int, placedTiles []Tile, availableTilesForPosition []int, availableTilesByPosition [][]int, isRetry bool) ([]Tile, [][]int, bool) {
	print(fmt.Sprintf("\tAvailableTiles before pick: %v\n", availableTilesForPosition), 3)
	positionsPossibleTiles := availableTilesForPosition
	for testTileNumber, testTile := range availableTilesForPosition { // try each available tile until it finds one that matches
		currentTile := getTileByID(testTile)
		print(fmt.Sprintf("\tTesting Tile: %v\n", currentTile), 3)
		currentTile, err := checkTileMatch(currentTile, position, placedTiles, 1)
		if err == nil { // tile worked, move on to the next position
			// place the tile
			placedTiles, availableTilesByPosition = placeTile(placedTiles, currentTile, availableTilesByPosition, positionsPossibleTiles, position, isRetry)
			availableTilesForPosition = availableTilesByPosition[position]
			isRetry = false
			break // break out of the loop of position's available tiles
		} else { // tile didnt work, so remove it from position's available tiles
			print(fmt.Sprintf("\tTileID %v did not work, removing..\n", currentTile.ID), 3)

			positionsPossibleTiles = removeAvailableTile(positionsPossibleTiles, currentTile.ID)

			// add or update the array of all available tiles
			if len(availableTilesByPosition) < position+1 || len(availableTilesByPosition) == 0 {
				availableTilesByPosition = append(availableTilesByPosition, positionsPossibleTiles)
			} else {
				availableTilesByPosition[position] = positionsPossibleTiles
			}
			print(fmt.Sprintf("\t\tPosition's available tiles after removal: %v\n", positionsPossibleTiles), 2)
		}
		if testTileNumber == len(availableTilesForPosition)-1 && position > 0 {
			print(fmt.Sprintf("\tNo tiles worked.  Retrying previous position. \n"), 3)
			// none of the available tiles for the position worked, so we need to go back one position and try a new tile
			isRetry = true
		}
	}

	return placedTiles, availableTilesByPosition, isRetry
}

func main() {
	startTime := time.Now()
	// Get commandline arguments
	args := os.Args[1:]
	printDetails, _ = strconv.Atoi(args[0])
	maxAttempts, _ := strconv.Atoi(args[1])
	firstTileID, _ := strconv.Atoi(args[2])
	// get a list of all tile ids
	allTileIDs := []int{}
	for _, tile := range tiles {
		allTileIDs = append(allTileIDs, tile.ID)
	}
	// Set Up Initial Variables
	attemptNumber := 0                    // tracks.. number of attempts
	isRetry := false                      // flag to tell if we are retrying a specific position
	placedTiles := []Tile{}               // array of tiles (in order, can calculate coorindates later)
	availableTilesByPosition := [][]int{} // array of arrays that holds available tiles for specific position
	poolOfAvailableTiles := []int{}       // tiles for this position
	checkTiles := true                    // flag to tell if you should check tiles (for instance, if there are no more tiles, dont try to check)
	firstTileRotationNumber := 0          // tracks which rotation number the first tile is on
	firstTileWasRotated := false
	maxPositionNumber := 0 // just bc im interested..
	// rotationsByPosition := []int{}

	// place the first tile
	placedTiles = append(placedTiles, getTileByID(firstTileID))
	availableTilesByPosition = append(availableTilesByPosition, getAvailableTiles(placedTiles, allTileIDs))
	position := 1 // position of tile (will calculate coordinates from this later)

	fmt.Printf("# Attempts: %v\n", maxAttempts)
	for position < len(allTileIDs) {
		poolOfAvailableTiles = []int{}
		attemptNumber += 1

		// if attemptNumber > 1 && position == 1 && placedTiles[1].ID == 1 {
		// 	break
		// }

		print(fmt.Sprintf("\nPosition: %v\tAttempt: %v/%v\tisRetry: %v\tfirstTileRotationNumber: %v\tnumAvailableFirstTiles: %v\n", position, attemptNumber, maxAttempts, isRetry, firstTileRotationNumber, len(availableTilesByPosition[0])), 0)
		print(fmt.Sprintf("\tPlacedTiles: %v\n", placedTiles), 1)
		if attemptNumber > maxAttempts {
			fmt.Printf("\n-----------")
			fmt.Printf("\nBreaking out due max attempts")
			break
		}

		if !isRetry || len(availableTilesByPosition) < position || firstTileWasRotated {
			/*	Recalc available tiles when:
				it is not a retry (so we moved forward in the position)
				it is the first time trying this position
			*/
			poolOfAvailableTiles = allTileIDs
		} else if isRetry {
			if len(placedTiles)-1 == position {
				// remove current tile from the placed tiles
				placedTiles = removePlacedTile(placedTiles, position)
				print(fmt.Sprintf("\tPlacedTiles after Removal: %v\n", placedTiles), 1)
			}
			/* use the already calculated tiles when:
			you are retrying and there are tiles available
			*/
			poolOfAvailableTiles = availableTilesByPosition[position]
		}
		print(fmt.Sprintf("\tpoolOfAvailableTileIDs: %v\tcheckTiles:%v\n", poolOfAvailableTiles, checkTiles), 1)
		// reset variables
		firstTileWasRotated = false
		checkTiles = true
		isRetry = false
		if len(poolOfAvailableTiles) == 0 {
			/* there are no more tiles available for this position, so
			do not check any tiles (bc there are none)
			set retry flag so we can go back
			*/
			checkTiles = false
			isRetry = true
		}
		if checkTiles {
			// Compare pool of available tiles to placed tiles to get the updated available tiles
			availableTilesForCurrentPosition := getAvailableTiles(placedTiles, poolOfAvailableTiles)
			print(fmt.Sprintf("\tavailableTilesForCurrentPosition: %v\n", availableTilesForCurrentPosition), 1)
			placedTiles, availableTilesByPosition, isRetry = checkTilesForPosition(position, placedTiles, availableTilesForCurrentPosition, availableTilesByPosition, isRetry)
		}
		print(fmt.Sprintf("\n\tisRetry: %v, numPlacedTiles: %v\n", isRetry, len(placedTiles)), 1)
		/*
			if the first tile got rotated, then the position needs to stay at position 1
			otherwise, if it is not a retry, move forward
			if it is a retry, move back
		*/
		if !isRetry {
			position += 1
		} else if position == 1 {
			if firstTileRotationNumber <= 3 { // rotate the first tile
				isRetry = true
				firstTileRotationNumber += 1
				placedTiles = rotateFirstTile(placedTiles[0])
				firstTileWasRotated = true
			} else {
				position -= 1
			}
		} else {
			position -= 1
		}
	}

	totalTime := time.Now().Sub(startTime)

	fmt.Printf("\n-----------\n")
	fmt.Printf("Total Time: %v sec (%v hours)\n", totalTime.Seconds(), totalTime.Hours())
	if totalTime.Seconds() > 0.5 {
		fmt.Printf("Averaging %v tiles per second.\n", attemptNumber/int(totalTime.Seconds()))
	}
	if position == len(tiles) {
		fmt.Printf("Found Solution in %v attempts.\n", attemptNumber)
		fmt.Printf("\nPlaced TileIDs: %v\n", getTileIDs(placedTiles))

		fmt.Printf("Placed Tiles: \n")
		pprintTiles(placedTiles)

	} else {
		fmt.Printf("Max position out of %v tries: %v\n", attemptNumber, maxPositionNumber)
		fmt.Print("Current attempt: \n")
		fmt.Println(getTileIDs(placedTiles))
		pprintTiles(placedTiles)
	}
}
