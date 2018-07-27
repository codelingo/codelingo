package main

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

const hubFrontend = "http://localhost:8080"

func main() {
	r := gin.Default()
	r.GET("/*path", func(c *gin.Context) {

		path := ".." + c.Param("path")
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		c.Header("Access-Control-Allow-Origin", hubFrontend)
		c.String(200, string(data))
	})

	r.Run(":3000")
}
