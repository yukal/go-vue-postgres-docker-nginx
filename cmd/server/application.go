package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"yu/webserver/internal/app/Repository/postgres"

	"github.com/gin-gonic/gin"
)

type Payload struct {
	Filter postgres.Filter `json:"filter" query:"filter"`
}

func getPing(c *gin.Context) {
	// c.JSON(http.StatusNoContent, gin.H{
	// 	"message": "pong",
	// })

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func InitRoutes(router *gin.Engine, db *postgres.Instance) {
	router.GET("/api/ping", getPing)

	// ...
}

func unparseJsonData[T any](body io.ReadCloser, data T) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(&data)
}

func MustReadFile(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return string(content)
}
