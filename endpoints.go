package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
