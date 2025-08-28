package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"hls-on-the-fly/internal/environment"
)

type Server struct {
	env *environment.Env
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	fmt.Println("Running on:", port)
	NewServer := &Server{
		env: environment.GetEnv(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", NewServer.env.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
