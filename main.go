package main

import (
	"github.com/gin-gonic/gin"
	"flag"
	"fmt"
	"net/http"
	"mime/multipart"
	"log"
)

// TODO: do we want to support more than just JSON?
// TODO: want to add some field for responding in some thread vs starting a new thread
// OR potentially add two different endpoints depending on this
type PostBody struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

func handleFile(c *gin.Context, file *multipart.FileHeader) {
	log.Println(file.Filename)

	c.SaveUploadedFile(file, "./tmp_files/" + file.Filename)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!\r\n", file.Filename))
}

func handleBody(c *gin.Context) {
	var postbody PostBody
	if c.BindJSON(&postbody) == nil {
		log.Println(postbody.Name)
		log.Println(postbody.Body)
		c.String(http.StatusOK, "bound by JSON!")
	}
}

func fullPost(c *gin.Context) {
	file, err := c.FormFile("file")

	handleBody(c)

	if err == nil {
		log.Println("handling file...")
		handleFile(c, file)
	} else {
		log.Println("not handling file")
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

