package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type GameState struct {
	repo []string
	word 		 string
	guess 	 string
}

type Request struct {
	Letter string `json:"letter"` // Correct: exported and tagged
}

var game GameState

func getGuess(c echo.Context) error {
	guess := fmt.Sprintf("World: %s", game.guess)

	return c.String(http.StatusOK, guess)
} 

func postGuess(c echo.Context) error {
	var req Request

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	processGuess(req.Letter)

	guess := fmt.Sprintf("World: %s", game.guess)
	return c.String(http.StatusOK, guess)
} 

func getRestart(c echo.Context) error {
	game.guess = ""

	start()

	guess := fmt.Sprintf("World: %s", game.guess)
	return c.String(http.StatusOK, guess)
} 

func start() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	game.repo = []string{"skate"}
	game.word = game.repo[r.Intn(len(game.repo))]

	for range game.word {
		game.guess += "_"
	}
}

func processGuess(letter string) {
	letter = strings.ToLower(letter)

	if strings.Contains(game.word, letter) {
		for i, c := range game.word {
			if string(c) == letter {
				game.guess = game.guess[:i] + letter + game.guess[i+1:]
			}
		}
	}

	if !strings.Contains(game.guess, "_") {
		 game.guess = "Parabens, voce acertou!!!\n A palavra era " + game.guess
	}
}

func main() {
	start()

	e := echo.New()

	e.GET("/guess", getGuess)
	e.GET("/restart", getRestart)

	e.POST("/guess", postGuess)
	
	e.Logger.Fatal(e.Start(":5000"))
}