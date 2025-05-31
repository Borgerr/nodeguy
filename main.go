package main

import (
	"flag"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Borgerr/nodeguy/docs"
)

// @title           nodeguy
// @version         0.1
// @description     Backend API for a forum website.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// TODO: do we want to support more than just JSON?
// TODO: want to add some field for responding in some thread vs starting a new thread
// OR potentially add two different endpoints depending on this
type PostBody struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

func handleFile(c *gin.Context, file *multipart.FileHeader) {
	log.Println(file.Filename)

	c.SaveUploadedFile(file, "./tmp_files/"+file.Filename)
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

// @Summary full post
// @Schemes
// @Description makes a post
// @Success 200
// @Router /post [post]
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

func setupRouter(url string) *gin.Engine {
	r := gin.Default()
	// Set a lower memory limit for multipart forms
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/")
		{
			eg.POST("/post", fullPost)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// TODO: determine what this is used for, docs make it seem important
	/*
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL(url),
			ginSwagger.DefaultModelsExpandDepth(-1))
	*/

	return r
}

func main() {
	hostPtr := flag.String("host", "localhost", "host for API")
	portPtr := flag.String("port", "8080", "open port for API")
	flag.Parse()

	router := setupRouter(fmt.Sprintf("http://%s:%s/swagger/doc.json", *hostPtr, *portPtr))

	router.Run(fmt.Sprintf(":%s", *portPtr))
}
