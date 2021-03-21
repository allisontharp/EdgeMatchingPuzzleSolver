package main

import (
	"fmt"
	"strings"
)

type Tile struct {
	North string
	East  string
	South string
	West  string
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

func pprintTile(tile Tile, maxWidth int) {
	topSpaces := strings.Repeat(" ", (maxWidth-1-len(tile.North))/2)
	bottomSpaces := strings.Repeat(" ", (maxWidth-1-len(tile.South))/2)
	middleSpaces := strings.Repeat(" ", (maxWidth - 1 - len(tile.East) - len(tile.West)))
	fmt.Printf("%v\n", strings.Repeat("-", maxWidth))
	fmt.Printf("|%v%v%v|\n", topSpaces, tile.North, topSpaces)
	fmt.Printf("|%v%v%v|\n", tile.West, middleSpaces, tile.East)
	fmt.Printf("|%v%v%v|\n", bottomSpaces, tile.South, bottomSpaces)
	fmt.Printf("%v\n", strings.Repeat("-", maxWidth))
}

func main() {
	t := tiles[0]
	pprintTile(t, 60)

}
