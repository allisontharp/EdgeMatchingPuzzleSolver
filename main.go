/*
	TODO: Handle cases where a single tile has multiple valid placements.
		in the availableTilesByPosition, maybe have each position be [[tileids], rotationNumber] ?
*/

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
Tile numbers will be:
0 1 2
3 4 5
6 7 8
*/

var maxPPrintWidth int = 30

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

var tiles2 = []Tile{
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

func pprintTiles(tiles []Tile, numTotalTiles int) {
	width := math.Sqrt(float64(numTotalTiles)) // dont want to do len(tiles) because this way, we can print partial puzzles
	height := width                            // not great, but for now this is fine (expecting puzzle to only be a square)

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

func splitTileSide(side String) TileSide{
	tileSide := TileSide{
		Direction: side[0:1],
		Description: [side1:],
	}
}
func checkEdgeMatch(currentTileSide string, testTileSide string) bool {
	current := splitTileSide(currentTileSide)
	test := splitTileSide(testTileSide)
}

func main() {
	// Get command line arguments
	args := os.Args[1:]
	maxAttempts, _ := strconv.Atoi(args[0])
	firstTileIndex, _ := strconv.Atoi(args[1])
	printLevel, _ := strconv.Atoi(args[2])

	// Set Variables
	placedTiles := []Tile{} // placed tiles (this also stores orientation)
	position := 0           // position of current tile attempt
	attemptNumber := 0      // current attempt number (used for stopping after X attempts)
	isRetry := false        // denotes if we are retrying a tile
	//// Set pool of all tiles
	allTileIDs := []int{}
	for _, tile := range tiles {
		allTileIDs = append(allTileIDs, tile.ID)
	}
	firstTileRotationNumber := 0 // first tile needs to rotate 3 times before replacing

	// loop through each position
	for position < len(allTileIDs) {
		// stop if we have reached max attempts
		attemptNumber += 1
		if attemptNumber > maxAttempts {
			break
		}
		fmt.Printf("Position: %v\tAttempt: %v/%v\tisRetry: %v\tfirstTileRotationNumber: %v\n", position, attemptNumber, maxAttempts, isRetry, firstTileRotationNumber)

	}

}
