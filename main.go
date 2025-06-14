package main

import (
	"flag"
	"fmt"

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
		api.GET("/:board/:threadID/get-thread", getThread) // NOTE: couldn't get this to work without a little constant URI field at the end
		api.GET("/:board/get-threads", getActiveThreads)

		// TODO: give DELETE methods some configuration options
		api.DELETE("/:board/:threadID/delete-thread", deleteThread)
		api.DELETE("/:board/:threadID/:replyID/delete-reply", deleteReply)

		// TODO: same treatment as deletion.
		api.PUT("/:board/:threadID/edit-thread", editThread)
		api.PUT("/:board/:threadID/:replyID/edit-reply", editReply)
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
