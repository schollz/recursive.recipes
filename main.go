package main

import (
	"flag"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/schollz/logger"
	"github.com/schollz/recursive.recipes/src/recipe"
	"github.com/schollz/recursive.recipes/src/recurse"
)

var flagPort int

func init() {
	flag.IntVar(&flagPort, "port", 8334, "port to run the server on")
}

func main() {
	var err error
	log.Debugf("listening on :%d", flagPort)
	http.HandleFunc("/", handler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", flagPort), nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().UTC()
	// Redirect URLs with trailing slashes (except for the root "/")
	if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, strings.TrimRight(r.URL.Path, "/"), http.StatusPermanentRedirect)
		return
	}
	err := handle(w, r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Debugf("%v %v %v %s\n", r.RemoteAddr, r.Method, r.URL.Path, time.Since(t))
}

func handle(w http.ResponseWriter, r *http.Request) (err error) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// very special paths
	if r.URL.Path == "/ws" {
		return handleWebsocket(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/static") {
		filename := r.URL.Path[1:]
		// serve filename
		mimeType := mime.TypeByExtension(filepath.Ext(filename))
		log.Tracef("serving %s as %s", filename, mimeType)
		w.Header().Set("Content-Type", mimeType)
		var b []byte
		b, err = os.ReadFile(filename)
		if err != nil {
			log.Errorf("could not read %s: %s", filename, err.Error())
			return
		}
		w.Write(b)
	} else if r.URL.Path == "/" {
		var b []byte
		b, err = os.ReadFile("templates/base.html")
		if err != nil {
			log.Errorf("could not read template %s: %s", "index.html", err.Error())
			return
		}
		tmpl, errTemplate := template.New("index").Parse(string(b))
		// tmpl, errTemplate := template.New("index").Delims("[[", "]]").Parse(string(b))
		if errTemplate != nil {
			log.Errorf("could not parse template %s: %s", "index.html", errTemplate.Error())
			return
		}

		recipes := []string{
			"apple pie",
			"chocolate chip cookie",
			"chocolate cake",
			"bread",
			"egg sandwich",
		}
		sort.Strings(recipes)
		data := struct {
			Recipes []string
		}{
			Recipes: recipes,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Errorf("could not execute template %s: %s", "index.html", err.Error())
		}
		return
	} else {
		var b []byte
		b, err = os.ReadFile("templates/index.html")
		if err != nil {
			log.Errorf("could not read template %s: %s", "index.html", err.Error())
			return
		}
		tmpl, errTemplate := template.New("index").Parse(string(b))
		// tmpl, errTemplate := template.New("index").Delims("[[", "]]").Parse(string(b))
		if errTemplate != nil {
			log.Errorf("could not parse template %s: %s", "index.html", errTemplate.Error())
			return
		}

		routeSplit := strings.Split(r.URL.Path[1:], "/")
		recipeName := routeSplit[0]
		fromScratch := routeSplit[1:]

		var recipes []recipe.Recipe
		var ingredients []recipe.Ingredient
		var graphviz string
		recipes, ingredients, graphviz, err = recurse.Recipe(recipeName, fromScratch, "static/data/")
		if err != nil {
			log.Errorf("could not parse recipe: %s", err.Error())
			return
		}

		// check if ingredients have data
		for i, _ := range ingredients {
			filename := "static/data/" + ingredients[i].Singular + ".json"
			_, err = os.Stat(filename)
			if err == nil {
				ingredients[i].HasData = true
			}
		}

		data := struct {
			URL            string
			VersionCurrent string
			RecipeName     string
			Graphviz       string
			FromScratch    []string
			FromScratchLen int
			Recipes        []recipe.Recipe
			Ingredients    []recipe.Ingredient
		}{
			URL:            r.URL.Path,
			VersionCurrent: "v1.1.0",
			Graphviz:       graphviz,
			RecipeName:     strings.Title(recipeName),
			Recipes:        recipes,
			FromScratch:    fromScratch,
			FromScratchLen: len(fromScratch) - 1,
			Ingredients:    ingredients,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Errorf("could not execute template %s: %s", "index.html", err.Error())
		}
		return
	}

	return
}

var upgrader = websocket.Upgrader{} // use default options

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

var hasReloaded = false

func handleWebsocket(w http.ResponseWriter, r *http.Request) (err error) {
	// use gorilla to open websocket
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	if !hasReloaded {
		c.WriteJSON(Message{Action: "do_reload"})
		hasReloaded = true
	}

	for {
		var message Message
		err := c.ReadJSON(&message)
		if err != nil {
			break
		}
		log.Tracef("message: %+v", message)
		if message.Action == "getinfo" {
		}
	}
	return
}
