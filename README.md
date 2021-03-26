# EdgeMatchingPuzzleSolver
This is an edge matching puzzle solver written in Go. 

**Edge Matching Puzzles**
An edge matching puzzle is a type of tiling puzzle involving tiling an area with (typically regular) polygons who edges are distinguished with colors or patterns in such a way that the edges match adjacent puzzles. [1](https://en.wikipedia.org/wiki/Edge-matching_puzzle)  

## Input

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

