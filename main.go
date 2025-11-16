package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/phuslu/log"
	"github.com/sanbei101/cau-ai-backend/handle"
)

func main() {
	handle.InitDish()
	log.DefaultLogger = log.Logger{
		Level: log.InfoLevel,
		Writer: &log.IOWriter{
			Writer: os.Stderr,
		},
	}
	router := handle.InitRouter()
	router.Use(cors.Default())
	log.Info().Uint("port", 8000).Msg("server started")
	router.Run(":8000")
}
