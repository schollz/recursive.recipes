package dfs

// depth-first search
type Node struct {
	name   string
	amount float64
}

type Edge struct {
	source *Node
	dest   *Node
	weight float64 // Add a weight field to represent the weight of the edge
}

type Graph struct {
	nodes []*Node
	edges []*Edge
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) AddNode(name string, amount float64) {
	g.nodes = append(g.nodes, &Node{name, amount})
}

// Modify AddEdge to include a weight parameter
func (g *Graph) AddEdge(s string, d string) {
	var source *Node
	var dest *Node
	for _, n := range g.nodes {
		if n.name == s {
			source = n
			break
		}
	}
	for _, n := range g.nodes {
		if n.name == d {
			dest = n
			break
		}
	}
	e := &Edge{source: source, dest: dest, weight: source.amount}
	g.edges = append(g.edges, e)
}

func (g *Graph) NodeList() []string {
	names := []string{}
	for _, n := range g.nodes {
		names = append(names, n.name)
	}
	return names
}

func (g *Graph) LongestOrdering() (ordering []string) {
	if len(g.NodeList()) < 1 {
		return
	}
	for {
		if len(g.NodeList()) == 1 {
			ordering = append(ordering, g.NodeList()[0])
			break
		}
		longestpath := g.LongestPath()
		if len(longestpath) == 0 {
			break
		}
		ordering = append(ordering, longestpath[0])
		g.RemoveNode(longestpath[0])
	}
	return
}

func (g *Graph) RemoveNode(name string) {
	// remove all edges that have n as a source or dest
	indices_to_remove := []int{}
	for i, e := range g.edges {
		if e.source.name == name || e.dest.name == name {
			indices_to_remove = append(indices_to_remove, i)
		}
	}
	new_edges := []*Edge{}
	for i, e := range g.edges {
		keep := true
		for _, j := range indices_to_remove {
			if i == j {
				keep = false
				break
			}
		}
		if keep {
			new_edges = append(new_edges, e)
		}
	}
	g.edges = new_edges

	// remove node
	for i, n := range g.nodes {
		if n.name == name {
			g.nodes = append(g.nodes[:i], g.nodes[i+1:]...)
		}
	}
}

func (g *Graph) DFS(s *Node, d *Node, visited map[*Node]bool, path []string, paths *[][]string) {
	visited[s] = true
	path = append(path, s.name)

	if s == d {
		*paths = append(*paths, path)
	} else {
		for _, e := range g.edges {
			if e.source == s && !visited[e.dest] {
				g.DFS(e.dest, d, visited, path, paths)
			}
		}
	}

	delete(visited, s)
	path = path[:len(path)-1]
}

func (g *Graph) AllPaths(s *Node, d *Node) [][]string {
	visited := make(map[*Node]bool)
	paths := [][]string{}
	path := []string{}

	g.DFS(s, d, visited, path, &paths)

	return paths
}

func (g *Graph) Length(path []string) float64 {
	length := 0.0
	for i := 0; i < len(path)-1; i++ {
		for _, e := range g.edges {
			if e.source.name == path[i] && e.dest.name == path[i+1] {
				length += e.weight
			}
		}
	}
	return length
}

func (g *Graph) LongestPath() []string {
	longestPathLength := 0.0
	longestPathList := []string{}
	for n1 := range g.nodes {
		for n2 := range g.nodes {
			if n1 == n2 {
				continue
			}
			paths := g.AllPaths(g.nodes[n1], g.nodes[n2])
			for _, p := range paths {
				length := g.Length(p)
				if length > longestPathLength {
					longestPathLength = length
					longestPathList = p
				}
			}
		}
	}

	return longestPathList
}
