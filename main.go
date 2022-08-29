package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Types used in templates
type Game struct {
	Board      Board
	GameNumber int
	Latency    int
}

type PlayerData struct {
	Games    []Game
	Latency  int
	NumGames int
}

// Templates
var files = []string{"./web/template/layout.html", "./web/template/myGames.html", "./web/template/game.html", "./web/template/home.html"}
var templ = template.Must(template.ParseFiles(files...))

// Our "database" and db helper functions
var boardsByPlayer = make(map[string][]Board)

func getBoards(player string, numGames int) []Board {
	if numGames == 0 {
		numGames = 4 // default to 4 games per player
	}

	boards, hasBoards := boardsByPlayer[player]
	if !hasBoards || numGames != len(boards) {
		boards = make([]Board, numGames)
		for i := 0; i < numGames; i++ {
			boards[i] = NewBoard()
		}
		boardsByPlayer[player] = boards
	}
	return boards
}

// Helper functions to extract query params, path vars
func getLatency(req *http.Request) int {
	latency, _ := strconv.Atoi(req.URL.Query().Get("latency"))
	return latency
}

func getNumGames(req *http.Request) int {
	latency, _ := strconv.Atoi(req.URL.Query().Get("num-games"))
	return latency
}

func getPlayer(req *http.Request) string {
	vars := mux.Vars(req)
	return vars["player"]
}

// Helper function to prepare game data for templates
func getGames(boards []Board, latency int) []Game {
	games := make([]Game, len(boards))
	for i, board := range boards {
		games[i] = Game{board, i, latency}
	}
	return games
}

// Route handlers
func handleLayout(w http.ResponseWriter, req *http.Request) {
	numGames := getNumGames(req)
	latency := getLatency(req)
	player := getPlayer(req)

	boards := getBoards(player, numGames)
	games := getGames(boards, latency)

	playerData := PlayerData{games, latency, numGames}
	err := templ.ExecuteTemplate(w, "layout", playerData)
	if err != nil {
		w.WriteHeader(500)
		log.Fatal("game:", err)
	}
}

func handleMove(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	latency := getLatency(req)
	player := getPlayer(req)

	boards, hasBoards := boardsByPlayer[player]
	if !hasBoards {
		log.Fatal("Player has no boards cannot make move, path ", req.URL.Path)
		w.WriteHeader(500)
		return
	}

	gameNumber, gameNumberErr := strconv.Atoi(vars["gameNumber"])
	if gameNumberErr != nil || gameNumber > len(boards)-1 {
		log.Fatal("Invalid game number ", vars["gameNumber"], " path ", req.URL.Path)
		w.WriteHeader(500)
		return
	}

	cell, cellErr := strconv.Atoi(vars["cell"])
	if cellErr != nil || cell > 8 {
		w.WriteHeader(500)
		log.Fatal("Invalid cell number ", vars["cell"], "path ", req.URL.Path)
		return
	}

	board := boards[gameNumber]
	board.Move(cell, X)

	game := Game{board, gameNumber, latency}

	// sleep
	log.Println("latency", latency)
	time.Sleep(time.Duration(latency * int(time.Millisecond)))

	tmplErr := templ.ExecuteTemplate(w, "game", game)
	if tmplErr != nil {
		log.Fatal("template error ", tmplErr)
		w.WriteHeader(500)
	}
}

func handleResetGames(w http.ResponseWriter, req *http.Request) {
	player := getPlayer(req)
	numGames := getNumGames(req)
	latency := getLatency(req)

	boardsByPlayer[player] = []Board{}

	boards := getBoards(player, numGames)
	games := getGames(boards, latency)

	playerData := PlayerData{games, latency, numGames}
	err := templ.ExecuteTemplate(w, "myGames", playerData)
	if err != nil {
		w.WriteHeader(500)
		log.Fatal("game:", err)
	}
}

func handleHome(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		player := req.FormValue("player")

		const DefaultNumGames = 4
		const DefaultLatency = 0

		gamesUrl := fmt.Sprintf("/players/%s/games?num-games=%d&latency=%d", player, DefaultNumGames, DefaultLatency)
		http.Redirect(w, req, gamesUrl, 302)
		return
	}

	err := templ.ExecuteTemplate(w, "home", nil)
	if err != nil {
		w.WriteHeader(500)
		log.Fatal("home:", err)
	}
}

// Main routine to setup routes and start web server
func main() {
	// Serve up static files, for stylesheets, htmx
	fs := http.FileServer(http.Dir("./web/static/"))

	// Setup our routes
	route := mux.NewRouter()
	route.HandleFunc("/", handleHome)
	route.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	route.HandleFunc("/players/{player}/games", handleLayout)
	route.HandleFunc("/players/{player}/games/{gameNumber}/cells/{cell}", handleMove)
	route.HandleFunc("/players/{player}/reset", handleResetGames)
	http.Handle("/", route)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
