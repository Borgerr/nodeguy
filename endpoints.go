package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: probably double check multipart/form-data from https://stackoverflow.com/questions/1443158/binary-data-in-json-string-something-better-than-base64
// once we get to actually submitting and uploading images

type BoardIdentifier struct {
	Board string `uri:"board" binding:"required"`
}
type ThreadIdentifier struct {
	Board    string `uri:"board" binding:"required"`
	ThreadID string `uri:"threadID" binding:"required"`
}
type ReplyIdentifier struct {
	Board    string `uri:"board" binding:"required"`
	ThreadID string `uri:"threadID" binding:"required"`
	ReplyID  string `uri: "replyID" binding:"required"`
}

// @Summary	post a new thread
// @Schemes
// @Description	Create a new thread in a board for others to reply to.
// @Success		200
// @Router			/:board/new-thread [post]
func newThread(c *gin.Context) {
	var boardUri BoardIdentifier
	if err := c.ShouldBindUri(&boardUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": boardUri.Board})
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

	c.String(http.StatusOK, "0") // TODO: return thread ID
}

// @Summary	reply to a thread
// @Schemes
// @Decription	Reply to an existing thread in a board.
// @Success	200
// @Router		/:board/:threadID/reply [post]
func replyToThread(c *gin.Context) {
	var ThreadUri ThreadIdentifier
	if err := c.ShouldBindUri(&ThreadUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": ThreadUri.Board, "threadID": ThreadUri.ThreadID})
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
	c.String(http.StatusOK, "1") // TODO: return reply ID
}

// @Summary	get thread contents
// @Schemes
// @Description	Get thread and replies in specified board with corresponding ID.
// @Success		200
// @Router			/:board/:threadID [get]
func getThread(c *gin.Context) {
	var ThreadUri ThreadIdentifier
	if err := c.ShouldBindUri(&ThreadUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": ThreadUri.Board, "threadID": ThreadUri.ThreadID})
	}
	return

	// -------------------
	// DB INTERACTION HERE
	// -------------------
}

// @Summary	get IDs of active threads
// @Schemes
// @Description	Gets the IDs of active threads in a board, depending on configuration of what constitutes an "active thread".
// @Success		200
// @Router			/:board/get-threads [get]
func getActiveThreads(c *gin.Context) {
	var boardUri BoardIdentifier
	if err := c.ShouldBindUri(&boardUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": boardUri.Board})
	}
	return

	// -------------------
	// DB INTERACTION HERE
	// -------------------
}

// @Summary	Delete thread with ID in board
// @Schemes
// @Description	Given the ID and board of a thread, attempt to delete a thread. May be accepted or denied depending on admin's configuration.
// @Success		200
// @Router			/:board/:threadID/delete-thread [delete]
func deleteThread(c *gin.Context) {
	// TODO:
	// right now we want people to have their digital footprint be "permanent"
	// for comical, embarrassing reasons.
	// People will need to live with what they say on our forum.
	// At some point, I'd like to get some configuration options
	// for a "deletion scheme", or some function implemented
	// that some goroutine drives and wakes back up to execute.
	//
	// Essentially, Say we had some function with pseudo-code:
	// ```
	// func deleteOld() {
	//   for thread in threadDB:
	//     if thread older than 1 day:
	//     delete thread
	// }
	// ```
	// We then would use some sort of branching depending on a flag or field
	// in a local config.json file, pointing towards that,
	// or towards something else.
	// We give the backend administrator different options depending on how
	// they would like to deploy this project.
	//
	// For now, though, as I said, forum users will need to live with their
	// embarrassing degeneracy.
	var ThreadUri ThreadIdentifier
	if err := c.ShouldBindUri(&ThreadUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": ThreadUri.Board,
				"threadID":        ThreadUri.ThreadID,
				"deletion-status": "Denied"}) // TODO: see note above
	}
	return
}

// @Summary	Delete reply with ID in thread with ID in board
// @Schemes
// @Description	Given the ID and board of a thread, as well as the ID of a reply, attempt to delete a reply. May be accepted or denied depending on admin's configuration.
// @Success		200
// @Router			/:board/:threadID/replyID/delete-reply [delete]
func deleteReply(c *gin.Context) {
	var replyUri ReplyIdentifier
	if err := c.ShouldBindUri(&replyUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": replyUri.Board,
				"threadID":        replyUri.ThreadID,
				"replyID":         replyUri.ReplyID,
				"deletion-status": "Denied"}) // TODO: see note above in deleteThread
	}
	return
}

// @Summary	Edit thread with ID in board
// @Schemes
// @Description	Given the ID and board of a thread, attempt to edit a thread. May be accepted or denied depending on admin's configuration.
// @Success		200
// @Router			/:board/:threadID/edit-thread [put]
func editThread(c *gin.Context) {
	var ThreadUri ThreadIdentifier
	if err := c.ShouldBindUri(&ThreadUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": ThreadUri.Board,
				"threadID":        ThreadUri.ThreadID,
				"edit-status": "Denied"}) // TODO: see note above
	}
	return
}

// @Summary	Edit reply with ID in thread with ID in board
// @Schemes
// @Description	Given the ID and board of a thread, as well as the ID of a reply, attempt to edit a reply. May be accepted or denied depending on admin's configuration.
// @Success		200
// @Router			/:board/:threadID/replyID/edit-reply [put]
func editReply(c *gin.Context) {
	var replyUri ReplyIdentifier
	if err := c.ShouldBindUri(&replyUri); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
	} else {
		c.JSON(http.StatusOK,
			gin.H{"board": replyUri.Board,
				"threadID":        replyUri.ThreadID,
				"replyID":         replyUri.ReplyID,
				"edit-status": "Denied"}) // TODO: see note above
	}
	return
}
