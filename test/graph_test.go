package test

import (
	"fmt"
	"github.com/xfali/goutils/v2/container/mapGraph"
	"testing"
)

func TestGraphBFS(t *testing.T) {
	testGraph := mapGraph.New()
	testGraph.AddEdge(1, 2)
	testGraph.AddEdge(1, 3)
	testGraph.AddEdge(1, 5)
	testGraph.AddEdge(2, 4)
	testGraph.AddEdge(3, 4)
	testGraph.AddEdge(3, 5)

	testGraph.BFS(1, func(i interface{}) {
		fmt.Printf("BFS get value: %d\n", i.(int))
	})
}

func TestGraphBFS2(t *testing.T) {
	testGraph := mapGraph.New()
	testGraph.AddEdge(1, 2)
	testGraph.AddEdge(1, 3)
	testGraph.AddEdge(1, 5)
	testGraph.AddEdge(2, 4)
	testGraph.AddEdge(3, 4)
	testGraph.AddEdge(3, 5)
	testGraph.AddEdge(4, 1)

	testGraph.BFS(1, func(i interface{}) {
		fmt.Printf("BFS get value: %d\n", i.(int))
	})
}

func TestGraphDFS(t *testing.T) {
	testGraph := mapGraph.New()
	testGraph.AddEdge(1, 2)
	testGraph.AddEdge(1, 3)
	testGraph.AddEdge(1, 5)
	testGraph.AddEdge(2, 4)
	testGraph.AddEdge(3, 4)
	testGraph.AddEdge(3, 5)

	testGraph.DFS(1, func(i interface{}) {
		fmt.Printf("DFS get value: %d\n", i.(int))
	})
}
