package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"words_backend/internal/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {

	resetFlag := flag.Bool("reset", false, "reset all worlds data in database")
	flag.Parse()

	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	viper.Set("reset", *resetFlag)

	fmt.Println(viper.Get("db_param"))
	fmt.Println(viper.Get("reset"))

	h := handler.New()
	r := gin.Default()
	r.GET("/ping", h.Ping)
	r.Run() // listen and serve on 0.0.0.0:8080
}
