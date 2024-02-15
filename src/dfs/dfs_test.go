package dfs

import (
	"fmt"
	"testing"
)

func TestDFS(t *testing.T) {
	g := NewGraph()
	g.AddNode("salt", 160.0)
	g.AddNode("butter", 2.0)
	g.AddNode("apple pie", 1.0)
	g.AddNode("cream", 1.0)
	g.AddNode("unsalted butter", 2.0)
	g.AddNode("pie crust", 1.0)

	g.AddEdge("salt", "butter")
	g.AddEdge("butter", "apple pie")
	g.AddEdge("pie crust", "apple pie")
	g.AddEdge("cream", "unsalted butter")
	g.AddEdge("cream", "butter")
	g.AddEdge("unsalted butter", "pie crust")
	g.AddEdge("salt", "pie crust")

	fmt.Println(g.LongestOrdering())
}
