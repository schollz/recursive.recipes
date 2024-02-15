package recurse

import (
	"fmt"
	"os"
	"sort"

	log "github.com/schollz/logger"
	"github.com/schollz/recursive.recipes/src/dfs"
	"github.com/schollz/recursive.recipes/src/recipe"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type RecipeGraph struct {
	g            *simple.WeightedDirectedGraph
	graphviz     string
	Nodes        map[string]recipe.Ingredient
	recipeLookup map[string]recipe.Recipe
}

func (r *RecipeGraph) RecipeOrdering() (recipes []recipe.Recipe) {
	// copy the graph, leaving out anything not in recipeLookup
	dfs := dfs.NewGraph()

	// Copy nodes
	nodes := r.g.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		for singular, object := range r.Nodes {
			if object.Node.ID() == node.ID() {
				if _, ok := r.recipeLookup[singular]; ok {
					dfs.AddNode(singular, r.recipeLookup[singular].TotalHours)
				}
			}
		}
	}

	// Copy edges
	edges := r.g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		from := edge.From().ID()
		to := edge.To().ID()
		for singular, object := range r.Nodes {
			if object.Node.ID() == from {
				if _, ok := r.recipeLookup[singular]; ok {
					for singular2, object2 := range r.Nodes {
						if object2.Node.ID() == to {
							if _, ok := r.recipeLookup[singular2]; ok {
								dfs.AddEdge(singular2, singular)
							}
						}
					}
				}
			}
		}
	}

	for _, recipeName := range dfs.LongestOrdering() {
		recipes = append(recipes, r.recipeLookup[recipeName])
	}

	return
}

func (r *RecipeGraph) NodeNumber(name string) int64 {
	i := int64(0)
	for _, node := range r.Nodes {
		if node.Singular == name {
			return int64(i)
		}
		i++
	}
	return -1
}

func (r *RecipeGraph) NodeName(i int64) string {
	for _, node := range r.Nodes {
		if node.Node.ID() == int64(i) {
			return fmt.Sprintf("node %d (%s)", i, node.Singular)
		}
	}
	return "not found"
}

func NewRecipeGraph() *RecipeGraph {
	return &RecipeGraph{
		recipeLookup: make(map[string]recipe.Recipe),
		g:            simple.NewWeightedDirectedGraph(0, 0),
		graphviz: `digraph G {
			color="#FFFFFF"
			bgcolor="#357EDD00" # RGBA (with alpha)
	`,
		Nodes: make(map[string]recipe.Ingredient),
	}
}

func (r *RecipeGraph) Graphviz() string {
	return r.graphviz + "}\n"
}

func (r *RecipeGraph) Recipes() (recipes []recipe.Recipe) {
	// Perform a topological sort to get the ingredients in order of dependencies
	order, err := topo.Sort(r.g)
	if err != nil {
		fmt.Println("There was an error with the topological sort:", err)
		return
	}

	// reverse sort
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	recipes = make([]recipe.Recipe, len(order))
	i := 0
	for _, node := range order {
		if r.g.From(node.ID()).Len() > 0 {
			for singular, object := range r.Nodes {
				if object.Node.ID() == node.ID() {
					if _, ok := r.recipeLookup[singular]; ok {
						recipes[i] = r.recipeLookup[singular]
						i++
					}
				}
			}
		}
	}
	recipes = recipes[:i]

	return recipes
}

func (r *RecipeGraph) AddRecipe(name string) (err error) {
	var r2 recipe.Recipe
	r2, err = recipe.Parse(name)
	if err != nil {
		log.Errorf("could not parse recipe: %s", err)
		return
	}
	r.recipeLookup[r2.RecipeSingular] = r2

	for i, _ := range r2.Inputs {
		r.AddNode(r2.Inputs[i])
	}
	for i, _ := range r2.Outputs {
		r.AddNode(r2.Outputs[i])
	}

	// Define edges based on dependencies
	for _, input := range r2.Inputs {
		for _, output := range r2.Outputs {
			r.AddEdge(input.Singular, output.Singular)
		}
	}
	return
}

func (r *RecipeGraph) AddNode(ingredient recipe.Ingredient) {
	if _, ok := r.Nodes[ingredient.Singular]; ok {
		// TODO combine the ingredient amounts!
		return
	}
	ingredient.Node = r.g.NewNode()

	r.graphviz += fmt.Sprintf("\"%s\"  [color=\"black\", fontcolor=\"black\",fillcolor=\"#%s\",style=\"filled\"];\n", ingredient.Singular, ingredient.Color)
	r.Nodes[ingredient.Singular] = ingredient
	r.g.AddNode(ingredient.Node)
}

func (r *RecipeGraph) AddEdge(from, to string) {
	if _, ok := r.Nodes[from]; !ok {
		return
	}
	if _, ok := r.Nodes[to]; !ok {
		return
	}
	// make sure edge doesn't already exist
	edges := r.g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		if edge.From().ID() == r.Nodes[to].Node.ID() && edge.To().ID() == r.Nodes[from].Node.ID() {
			return
		}
	}

	weight := r.recipeLookup[r.Nodes[from].Singular].TotalHours + r.recipeLookup[r.Nodes[to].Singular].TotalHours
	r.g.SetWeightedEdge(r.g.NewWeightedEdge(r.Nodes[to].Node, r.Nodes[from].Node, weight))
	r.graphviz += fmt.Sprintf("\"%s\" -> \"%s\"  [color=\"black\", style=\"filled\"];\n", r.Nodes[from].Singular, r.Nodes[to].Singular)
}

func (r *RecipeGraph) Ingredients() (ingredients []recipe.Ingredient, err error) {
	// Perform a topological sort to get the ingredients in order of dependencies
	order, err := topo.Sort(r.g)
	if err != nil {
		fmt.Println("There was an error with the topological sort:", err)
		return
	}

	// reverse sort
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	ingredients = make([]recipe.Ingredient, len(order))
	ingredientI := 0
	for _, node := range order {
		if r.g.From(node.ID()).Len() == 0 {
			for ing_singular, ing_object := range r.Nodes {
				if ing_object.Node.ID() == node.ID() {
					ingredients[ingredientI] = r.Nodes[ing_singular]
					ingredientI++
				}
			}
		}
	}
	ingredients = ingredients[:ingredientI]

	// sort ingredients by type
	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].Type < ingredients[j].Type
	})

	return
}

func Recipe(name string, from_scratch []string, folder string) (recipes []recipe.Recipe, ingredients []recipe.Ingredient, graphviz string, err error) {
	filename := folder + name + ".json"
	// check if file exists
	_, err = os.Stat(filename)
	if err != nil {
		return
	}
	// Create a new directed graph
	rg := NewRecipeGraph()

	// load recipe
	err = rg.AddRecipe(filename)
	if err != nil {
		log.Error(err)
		return
	}

	for _, ingredient_string := range from_scratch {
		rg.AddRecipe(folder + ingredient_string + ".json")
	}

	// get the graphviz
	graphviz = rg.Graphviz()

	// list of ingredients
	ingredients, err = rg.Ingredients()
	if err != nil {
		log.Error(err)
	}

	recipes = rg.RecipeOrdering()

	return
}
