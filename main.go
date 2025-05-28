package main

import (
	"github.com/gin-gonic/gin"
	//"encoding/json"
	"flag"
	"fmt"
	//"io"
	"net/http"
	"mime/multipart"
)

func handleFile(c *gin.Context, file *multipart.FileHeader) {
	fmt.Println(file.Filename)

	c.SaveUploadedFile(file, "./tmp_files/" + file.Filename)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!\r\n", file.Filename))
}

func handleBody(c *gin.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.PostForm("name")
	message := c.PostForm("message")

	fmt.Printf("id: %s; page: %s; name: %s; message: %s\r\n", id, page, name, message)
}

func fullPost(c *gin.Context) {
	file, err := c.FormFile("file")

	handleBody(c)

	if err == nil {
		handleFile(c, file)
	} else {
		c.String(http.StatusOK, "No file upload? Totally cool.\r\n")
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	// Set a lower memory limit for multipart forms
	router.MaxMultipartMemory = 8 << 20		// 8 MiB

	router.POST("/post", fullPost)

	return router
}

func main() {
	wordPtr := flag.String("port", "8080", "open port for API")
	flag.Parse()

	router := setupRouter()

	router.Run(fmt.Sprintf(":%s", *wordPtr))
}

