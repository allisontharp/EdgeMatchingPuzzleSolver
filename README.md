# EdgeMatchingPuzzleSolver
This is an edge matching puzzle solver written in Go. 

**Edge Matching Puzzles**
An edge matching puzzle is a type of tiling puzzle involving tiling an area with (typically regular) polygons who edges are distinguished with colors or patterns in such a way that the edges match adjacent puzzles. [1](https://en.wikipedia.org/wiki/Edge-matching_puzzle)  

<img src="https://raw.githubusercontent.com/allisontharp/EdgeMatchingPuzzleSolver/SinglePlacementAttempt/images/ExamplePuzzle.png" height="500"/>

## Input

Update the inputPuzzle.json file with the details of the puzzle. 

- Each tile is represented by a distinct tile ID
- The edge descriptions are broken into 2 components, L/R and description
    - An edge matches if the description exactly matches and one side starts with L and the other starts with R (for example, LTree matches RTree).  **It is very important that these settings are exactly right!**

## Usage

Build 
    `go build main.go`
Run
    `go run main.go <printLevel> <maxAttempts> <startTileID>`

- `<printLevel>` 
    - The level of details to print (0 is minimal -> 4 very detailed)
- `<maxAttempts>` 
    - The max number of attempts it will try before stopping (it will stop when it finds a solution if it is under this number)
- `<startTileID>` 
    - ID of the tile to start with

Example
    `go run main.go 0 20000 1`

## Output
If it was able to find a solution (for a 3x3, it is generally between 10k and 20k attempts), it will output as follows.  It describes the tile ids in order as well as the orientation (1 rotation means it is in the starting position, 4 rotations means it was rotated once counter clockwise).

![example-output](https://raw.githubusercontent.com/allisontharp/EdgeMatchingPuzzleSolver/SinglePlacementAttempt/images/ExampleOutput.png)

## To Do
- Allow for other sizes (not just 3x3)
- Allow for trying to find multiple solutions (besides the 4 rotations of the same solution)
- Make command line arguments better (flags)
