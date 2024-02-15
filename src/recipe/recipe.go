package recipe

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"

	"github.com/hansrodtang/randomcolor"
	log "github.com/schollz/logger"
	"github.com/schollz/recursive.recipes/src/humanize"
	"gonum.org/v1/gonum/graph"
)

type Recipe struct {
	Directions     []string     `json:"directions"`
	Hours          Hours        `json:"hours"`
	Inputs         []Ingredient `json:"input"`
	Outputs        []Ingredient `json:"output"`
	Recipe         string       `json:"recipe"`
	RecipePlural   string       `json:"recipe_plural"`
	RecipeSingular string       `json:"recipe_singular"`
	Duration       string
	TotalHours     float64
	InputCost      float64
}
type Hours struct {
	Parallel float64 `json:"parallel"`
	Serial   float64 `json:"serial"`
}
type Amount struct {
	Count float64 `json:"count"`
	Unit  string  `json:"unit"`
}
type Ingredient struct {
	Amount     Amount  `json:"amount"`
	Ingredient string  `json:"ingredient"`
	Singular   string  `json:"singular"`
	StoreCost  float64 `json:"store_cost"`
	Type       string  `json:"type"`
	Color      string  `json:"color"`
	Node       graph.Node
	HasData    bool
}

func Parse(fname string) (r Recipe, err error) {
	// read the file
	b, err := os.ReadFile(fname)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		log.Error(err)
		return
	}
	// create the color for every recipe using randomcolor
	for i, _ := range r.Inputs {
		// seed with ingredient name
		name := r.Inputs[i].Ingredient
		// hash the name to int64
		hash := int64(0)
		for _, c := range name {
			hash = int64(c) + (hash << 6) + (hash << 16) - hash
		}
		rand.Seed(hash)
		var hue randomcolor.Color
		switch r.Inputs[i].Type {
		case "animal":
			hue = randomcolor.Red
		case "natural":
			hue = randomcolor.Green
		case "spice":
			hue = randomcolor.Pink
		case "plants":
			hue = randomcolor.Green
		case "grain":
			hue = randomcolor.Yellow
		case "fungi":
			hue = randomcolor.Purple
		default:
			hue = randomcolor.Blue
		}
		rgba := randomcolor.New(hue, randomcolor.LIGHT)
		// convert rgba to hex
		hexColor := ""
		cc1, cc2, cc3, _ := rgba.RGBA()
		ccs := []uint8{uint8(cc1), uint8(cc2), uint8(cc3), uint8(100)}
		for _, c := range ccs {
			hexColor += string("0123456789ABCDEF"[c>>4]) + string("0123456789ABCDEF"[c&15])
		}
		r.Inputs[i].Color = hexColor
	}
	// create the color for every recipe using randomcolor
	for i, _ := range r.Outputs {
		// seed with ingredient name
		name := r.Outputs[i].Ingredient
		// hash the name to int64
		hash := int64(0)
		for _, c := range name {
			hash = int64(c) + (hash << 6) + (hash << 16) - hash
		}
		rand.Seed(hash)
		var hue randomcolor.Color
		switch r.Outputs[i].Type {
		case "animal":
			hue = randomcolor.Red
		case "natural":
			hue = randomcolor.Green
		case "spice":
			hue = randomcolor.Pink
		case "plants":
			hue = randomcolor.Green
		case "grain":
			hue = randomcolor.Yellow
		case "fungi":
			hue = randomcolor.Purple
		default:
			hue = randomcolor.Blue
		}
		rgba := randomcolor.New(hue, randomcolor.LIGHT)
		// convert rgba to hex
		hexColor := ""
		cc1, cc2, cc3, _ := rgba.RGBA()
		ccs := []uint8{uint8(cc1), uint8(cc2), uint8(cc3), uint8(100)}
		for _, c := range ccs {
			hexColor += string("0123456789ABCDEF"[c>>4]) + string("0123456789ABCDEF"[c&15])
		}
		r.Outputs[i].Color = hexColor
	}

	// calculate ingredient cost
	r.InputCost = 0.0
	for _, ingredient := range r.Inputs {
		r.InputCost += ingredient.StoreCost
	}

	// for every direction, replace recipe word with link
	for i, direction := range r.Directions {
		r.Directions[i] = strings.Replace(direction, r.Recipe, "<span class='recipe'>"+r.Recipe+"</span>", -1)
		for _, ingredient := range r.Inputs {
			r.Directions[i] = strings.Replace(r.Directions[i], ingredient.Ingredient+" ", "<span class='ingredient' style='background-color:#"+ingredient.Color+"'>"+ingredient.Ingredient+"</span> ", -1)
			r.Directions[i] = strings.Replace(r.Directions[i], ingredient.Singular+" ", "<span class='ingredient' style='background-color:#"+ingredient.Color+"'>"+ingredient.Ingredient+"</span> ", -1)
		}
		for _, ingredient := range r.Outputs {
			r.Directions[i] = strings.Replace(r.Directions[i], ingredient.Ingredient+" ", "<span class='ingredient' style='background-color:#"+ingredient.Color+"'>"+ingredient.Ingredient+"</span> ", -1)
			r.Directions[i] = strings.Replace(r.Directions[i], ingredient.Singular+" ", "<span class='ingredient' style='background-color:#"+ingredient.Color+"'>"+ingredient.Ingredient+"</span> ", -1)
		}
	}

	// make sure directions has a period at the end
	for i, direction := range r.Directions {
		if direction[len(direction)-1] != '.' {
			r.Directions[i] = direction + "."
		}
	}

	// convert to nice case
	r.Recipe = strings.Title(r.Recipe)

	r.TotalHours = r.Hours.Serial + r.Hours.Parallel
	r.Duration = humanize.Duration(int64((r.Hours.Serial + r.Hours.Parallel) * 60 * 60 * 1000))
	return
}
