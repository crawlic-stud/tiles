package game

import (
	"container/heap"
	"tiles/pkg/models"
)

type Node struct {
	Position models.Position

	G int // cost from start
	H int // heuristic to target
	F int // G + H

	Index int
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	node := x.(*Node)
	node.Index = len(*pq)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.Index = -1
	*pq = old[:n-1]
	return node
}

// Manhattan distance
func Heuristic(a, b models.Position) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func reconstructPath(
	cameFrom map[models.Position]models.Position,
	current models.Position,
) []models.Position {
	path := []models.Position{current}
	for {
		prev, ok := cameFrom[current]
		if !ok {
			break
		}
		current = prev
		path = append([]models.Position{current}, path...)
	}
	return path
}

func findPathAStar(
	grid models.Grid,
	start models.Position,
	target models.Position,
) []models.Position {

	directions := []models.Position{
		{X: 0, Y: -1},
		{X: 0, Y: 1},
		{X: -1, Y: 0},
		{X: 1, Y: 0},
	}

	openSet := &PriorityQueue{}
	heap.Init(openSet)

	startNode := &Node{
		Position: start,
		G:        0,
		H:        Heuristic(start, target),
	}
	startNode.F = startNode.G + startNode.H

	heap.Push(openSet, startNode)
	cameFrom := map[models.Position]models.Position{}
	gScore := map[models.Position]int{
		start: 0,
	}
	closedSet := map[models.Position]bool{}

	for openSet.Len() > 0 {
		currentNode := heap.Pop(openSet).(*Node)
		current := currentNode.Position

		// Skip outdated nodes
		if currentNode.G > gScore[current] {
			continue
		}

		// Skip already processed nodes
		if closedSet[current] {
			continue
		}
		closedSet[current] = true

		if current == target {
			return reconstructPath(cameFrom, current)
		}

		for _, dir := range directions {
			neighbor := models.Position{
				X: current.X + dir.X,
				Y: current.Y + dir.Y,
			}

			// bounds check
			if neighbor.X < 0 ||
				neighbor.Y < 0 ||
				neighbor.X >= grid.Width() ||
				neighbor.Y >= grid.Height() {
				continue
			}

			if closedSet[neighbor] {
				continue
			}

			tile := grid.GetTile(neighbor.X, neighbor.Y)
			if !tile.IsWalkable() {
				continue
			}

			moveCost := tile.MovementCost()
			tentativeG := gScore[current] + moveCost

			oldG, exists := gScore[neighbor]

			if !exists || tentativeG < oldG {

				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeG

				h := Heuristic(neighbor, target)
				node := &Node{
					Position: neighbor,
					G:        tentativeG,
					H:        h,
					F:        tentativeG + h,
				}
				heap.Push(openSet, node)
			}
		}
	}

	return nil
}
func FindPath(
	grid models.Grid,
	character *models.Character,
	target models.Position,
) []models.Position {
	path := findPathAStar(
		grid,
		character.Position,
		target,
	)
	return path
}
