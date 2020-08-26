package main

// This is where we learn how goroutines work.
//
// I generate a random graph, then spawn a bunch of routines to traverse it via
// DFS. There's synchronization, since we don't want multiple routines touching
// more nodes than they need to, and there's threading, since each one can
// traverse independently.
//
// Also there's structs, to represent the graph itself! Yay learning.
import (
	"container/list"
	"fmt"
	"math/rand"
)

type Node struct {
	neighbors []*Node
	value     int
}

func createRandomDigraph() []Node {
	// Create an arbitrary number of nodes.
	nodes := make([]Node, 500+rand.Intn(1000))

	// Connect them randomly by picking pairs.
	for i, node := range nodes {
		node.value = i
		edges := 500 + rand.Intn(len(nodes)-500)

		// Create list of possible endpoints
		candidates := list.New()
		for j, _ := range nodes {
			candidates.PushBack(j)
		}

		for j := 0; j < edges; j++ {
			// Choose random endpoint
			edge := rand.Intn(candidates.Len())

			// Find the actual neighbor it refers to
			head := candidates.Front()
			k := 0
			for k != edge {
				k++
				head.Next()
			}

			// Connect this node to the chosen node.
			nodes[i].neighbors = append(node.neighbors, &nodes[head.Value.(int)])

			// Remove this as a possible edge for next time.
			candidates.Remove(head)
		}
	}

	return nodes
}

func main() {
	graph := createRandomDigraph()
	fmt.Printf("%s\n", graph)
	// traverse(graph)
}
