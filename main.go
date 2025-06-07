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

//	@title			nodeguy
//	@version		0.1
//	@description	Backend API for a forum website.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

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

//	@Summary	full post
//	@Schemes
//	@Description	makes a post
//	@Success		200
//	@Router			/post [post]
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

// --------------
// creating posts
// --------------
// probably double check multipart/form-data from https://stackoverflow.com/questions/1443158/binary-data-in-json-string-something-better-than-base64

type NewThreadUri struct {
	Board	 string `uri:"board" binding:"required"`
}
//	@Summary	post a new thread
//	@Schemes
//	@Description	Create a new thread in a board for others to reply to.
//	@Success		200
//	@Router			/:board/new-thread [post]
func newThread(c *gin.Context) {
	var newThreadUri NewThreadUri
	if err := c.ShouldBindUri(&newThreadUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
		gin.H{"board": newThreadUri.Board})
	}
	return

	file, err := c.FormFile("file")
	// TODO: check filetype
	// (https://pkg.go.dev/github.com/h2non/filetype#section-readme maybe?)
	if err != nil {
		c.String(http.StatusBadRequest, "No file upload.")
	}
	// -------------------
	// DB INTERACTION HERE
	// -------------------
	log.Println(fmt.Sprintf("storing file %s", file.Filename))

	c.String(http.StatusOK, "0")	// TODO: return thread ID
}

type ReplyUri struct {
	Board	 string `uri:"board" binding:"required"`
	ThreadID string `uri:"threadID" binding:"required"`
}
//	@Summary	reply to a thread
//	@Schemes
// Decription Reply to an existing thread in a board.
//	@Success	200
//	@Router		/:board/:threadID/reply [post]
func replyToThread(c *gin.Context) {
	var replyUri ReplyUri
	if err := c.ShouldBindUri(&replyUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": replyUri.Board, "threadID": replyUri.ThreadID})
	}
	return

	file, err := c.FormFile("file")
	if err != nil {
		// -------------------
		// DB INTERACTION HERE
		// -------------------
		log.Println(fmt.Sprintf("storing file %s", file.Filename))
	} else {
		log.Println("not storing a file")
	}

	// -------------------
	// DB INTERACTION HERE
	// -------------------
	c.String(http.StatusOK, "1")	// TODO: return reply ID
}

func setupRouter(url string) *gin.Engine {
	r := gin.Default()
	// Set a lower memory limit for multipart forms
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	// TODO: check if we can just set the accepted filetypes here

	docs.SwaggerInfo.BasePath = "/"
	api := r.Group("/")
	// TODO: probably want to add an easy way to separate based on boards in some config file
	// very likely to be irrelevant to this block, though. Maybe just some script tool.
	{
		//eg.POST("/post", fullPost)
		// TODO: problem here with URI binding. Need to determine some other way of resolving this tomorrow.
		api.POST("/:board/new-thread", newThread)
		api.POST("/:board/:threadID/reply", replyToThread)
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
