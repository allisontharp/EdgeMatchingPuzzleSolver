/*
	TODO: Handle cases where a single tile has multiple valid placements.
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
var allTileIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
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

// var tiles = []Tile{
// 	{
// 		North: "RBeetle",
// 		East:  "LDragon",
// 		South: "RDragon",
// 		West:  "RGrasshopper",
// 		ID:    1,
// 	},
// 	{
// 		North: "LBeetle",
// 		East:  "RGrasshopper",
// 		South: "LAnt",
// 		West:  "RDragon",
// 		ID:    2,
// 	},
// 	{
// 		North: "RAnt",
// 		East:  "RBeetle",
// 		South: "LBeetle",
// 		West:  "LGrasshopper",
// 		ID:    3,
// 	},
// 	{
// 		North: "LDragon",
// 		East:  "RAnt",
// 		South: "RGrasshopper",
// 		West:  "LAnt",
// 		ID:    4,
// 	},
// 	{
// 		North: "RAnt",
// 		East:  "LDragon",
// 		South: "RGrasshopper",
// 		West:  "LBeetle",
// 		ID:    5,
// 	},
// 	{
// 		North: "RBeetle",
// 		East:  "LAnt",
// 		South: "RDragon",
// 		West:  "LGrasshopper",
// 		ID:    6,
// 	},
// 	{
// 		North: "LAnt",
// 		East:  "RBeetle",
// 		South: "LGrasshopper",
// 		West:  "LDragon",
// 		ID:    7,
// 	},
// 	{
// 		North: "LBeetle",
// 		East:  "RGrasshopper",
// 		South: "LDragon",
// 		West:  "RAnt",
// 		ID:    8,
// 	},
// 	{
// 		North: "RGrasshopper",
// 		East:  "LDragon",
// 		South: "RAnt",
// 		West:  "RBeetle",
// 		ID:    9,
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

func pprintTiles(tileIDs []int) {
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
			tile := getTileByID(tileIDs[i])
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
func checkTileMatch(currentTile Tile, position int, placedTiles []int, rotationNumber int) (Tile, error) {
	print(fmt.Sprintf("\t\tRotation Number: %v\n", rotationNumber))
	if rotationNumber > 3 {
		return currentTile, errors.New("Invalid Tile")
	}
	// first tile doesnt match on anything
	if position == 0 {
		fmt.Printf("FIRST TILE DOESNT MATCH DO SOMETHING HERE!\n")
		return currentTile, errors.New("First Tile Does Not Match")
	}
	sidesToMatch := getSidesToMatch(position, width)
	print(fmt.Sprintf("\t\tSidesToMatch: %v\n", sidesToMatch))
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

	args := os.Args[1:]
	printDetails, _ = strconv.ParseBool(args[0])
	maxAttempts, _ := strconv.Atoi(args[1])
	firstTileIndex, _ := strconv.Atoi(args[2])
	fmt.Printf("# Attempts: %v\n", maxAttempts)

	placedTiles := []int{}
	position := 0
	isPrevious := false

	availableTilesByPosition := [][]int{}
	availableTiles := []int{}

	attemptNumber := 0

	placedTiles = append(placedTiles, tiles[firstTileIndex].ID)
	availableTiles = getAvailableTiles(placedTiles, allTileIDs)
	availableTilesByPosition = append(availableTilesByPosition, availableTiles)
	print(fmt.Sprintf("\tcurrentTile: %v\n", tiles[firstTileIndex]))

	print(fmt.Sprintf("\tavailableTiles: %v\n", availableTiles))

	position += 1
	firstTileRotationNumber := 0

	for position < len(allTileIDs) {
		attemptNumber += 1
		runTileCheck := true
		if attemptNumber > maxAttempts {
			break
		}
		fmt.Printf("Position: %v\tAttempt: %v/%v\n", position, attemptNumber, maxAttempts)
		print(fmt.Sprintf("\tisPrevious: %v\n", isPrevious))
		print(fmt.Sprintf("\t--->Placed Tiles: %v<---\n", placedTiles))
		if !isPrevious {
			// the available tiles for the position reset to be whatever is available in the pool
			availableTiles = getAvailableTiles(placedTiles, allTileIDs)
			if len(availableTiles) == 0 {
				// none of the available tiles worked
				// need to use a different tile in the previous position (so, position=position-2)
				position = position - 2
				// need to also remove the placed tile
				placedTiles = removePlacedTile(placedTiles, len(placedTiles)-1)
				// need to mark that this is a retry
				isPrevious = true
				runTileCheck = false
			}
		} else if position > 0 {
			// the available tiles are the tiles that it hasn't yet tried
			// print(fmt.Sprintf("\tTotal tiles available in this position currently: %v\n", len(availableTilesByPosition[position])))
			// if len(availableTilesByPosition[position]) == 0 {
			// 	// none of the available tiles worked
			// 	// need to use a different tile in the previous position (so, position=position-2)
			// 	position = position - 2
			// 	// need to also remove the placed tile
			// 	placedTiles = removePlacedTile(placedTiles, len(placedTiles)-1)
			// 	print(fmt.Sprintf("\t--->Placed Tiles After Removal: %v<---\n", placedTiles))
			// 	// need to mark that this is a retry
			// 	isPrevious = true
			// 	runTileCheck = false
			// } else {
			// 	availableTilesByPosition[position] = removePlacedTile(availableTilesByPosition[position], 0)
			// 	print(fmt.Sprintf("\tTotal tiles available in this position after removal: %v\n", len(availableTilesByPosition[position])))
			// }
			if len(availableTilesByPosition[position]) == 0 {
				// there are no more tiles for this position
				originalPosition := position
				if !runTileCheck {
					position = position - 1
				} else {
					position = position - 2
				}
				if originalPosition == 1 { // the last tile for position 1 didnt work, so need to trigger a tile 0 rotate
					runTileCheck = false
					fmt.Printf("\nNumber of Rotations for First Tile Before: %v\n", firstTileRotationNumber)
					if firstTileRotationNumber <= 3 {
						firstTile := getTileByID(placedTiles[0])
						print("\tRotating first tile:\n")
						print(fmt.Sprintf("\t\tBefore: %v\n", firstTile))

						afterRotation := rotateTile(firstTile)
						print(fmt.Sprintf("\t\tAfter: %v\n", afterRotation))
						firstTileRotationNumber += 1
						fmt.Printf("\nNumber of Rotations for First Tile: %v\n", firstTileRotationNumber)
					} else {
						fmt.Printf("\t--->Placed tiles before replacement: %v<---\n", placedTiles)
						placedTiles[0] = availableTilesByPosition[0][0]
						availableTilesByPosition[0] = removePlacedTile(availableTilesByPosition[0], 0)
						fmt.Printf("\t--->Placed new Position 1 tile: %v<---\n", placedTiles)
						firstTileRotationNumber = 0
					}
				} else {
					// none of the available tiles worked
					// need to use a different tile in the previous position (so, position=position-2)
					// position = position - 2
					// need to also remove the placed tile
					placedTiles = removePlacedTile(placedTiles, len(placedTiles)-1)
					fmt.Printf("\t--->Placed Tiles: %v<---\n", placedTiles)
					// need to mark that this is a retry
					isPrevious = true
					runTileCheck = false
				}

			} else {
				availableTiles = getAvailableTiles(placedTiles, availableTilesByPosition[position])
			}

		} else { // redo first tile (it was already rotated)
			print(fmt.Sprintf("\tTotal tiles available in this position currently: %v\n", len(availableTilesByPosition[position])))
			runTileCheck = false
			isPrevious = false
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
					print(fmt.Sprintf("\t---------\n\tERROR: %v\n\t---------", err))
					print(fmt.Sprintf("\n\ttestTileNumber: %v | len(availableTiles): %v | position: %v\n", testTileNumber, len(availableTiles), position))
				}
				if testTileNumber == len(availableTiles)-1 && position > 0 {
					// none of the available tiles for the position worked, so we need to go back one position and try a new tile

					// go back a position and try again (-2 because we will add 1 to the position at the end of the loop)
					position = position - 2

					// need to also remove the placed tile if it wasn't the last one
					if len(placedTiles) > 1 {
						placedTiles = removePlacedTile(placedTiles, len(placedTiles)-1)
						fmt.Printf("\t--->Placed Tiles: %v<---\n", placedTiles)
					} else { // there is only one other tile
						if firstTileRotationNumber <= 3 {
							firstTile := getTileByID(placedTiles[0])
							print("\tRotating first tile:\n")
							print(fmt.Sprintf("\t\tBefore: %v\n", firstTile))

							afterRotation := rotateTile(firstTile)
							print(fmt.Sprintf("\t\tAfter: %v\n", afterRotation))
							firstTileRotationNumber += 1
							fmt.Printf("\nNumber of Rotations for First Tile: %v\n", firstTileRotationNumber)
						} else {
							fmt.Printf("\t--->Placed tiles before replacement: %v<---\n", placedTiles)
							fmt.Printf("\tavailableTilesByPosition[0]: %v\n", availableTilesByPosition[0])
							placedTiles[0] = availableTilesByPosition[0][0]
							availableTilesByPosition[0] = removePlacedTile(availableTilesByPosition[0], 0)
							fmt.Printf("\t--->Placed new Position 1 tile: %v<---\n", placedTiles)
							firstTileRotationNumber = 0
						}
					}
					// need to mark that this is a retry
					isPrevious = true
				}
			}
		}

		position += 1
	}

	// it found a match
	if len(placedTiles) == len(allTileIDs) {
		fmt.Printf("\n\nFound solution in %v attempts\n", attemptNumber)
		pprintTiles(placedTiles)
	}

	fmt.Println(getTileByID(7))

	fmt.Printf("\n\n-----------------\n")
	fmt.Printf("PlacedTiles: %v\n", placedTiles)
	fmt.Printf("-----------------\n")

}
