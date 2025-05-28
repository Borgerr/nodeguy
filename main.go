package main

import (
	"github.com/gin-gonic/gin"
	"flag"
	"fmt"
	//"io"
	"net/http"
)

func postBody(c *gin.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.PostForm("name")
	message := c.PostForm("message")

	fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
}

func postImage(c *gin.Context) {
	file, _ := c.FormFile("file")
	fmt.Println(file.Filename)

	// Upload the file to a specific dst
	// TODO: change to upload a blob to our DB
	// may want to change depending on testing vs prod
	c.SaveUploadedFile(file, "./tmp_files/" + file.Filename)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!\r\n", file.Filename))
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	// Set a lower memory limit for multipart forms
	router.MaxMultipartMemory = 8 << 20		// 8 MiB

	router.POST("/post", postBody)

	router.POST("/upload", postImage)

	return router
}

func main() {
	wordPtr := flag.String("port", "8080", "open port for API")
	flag.Parse()

	router := setupRouter()

	router.Run(fmt.Sprintf(":%s", *wordPtr))
}

