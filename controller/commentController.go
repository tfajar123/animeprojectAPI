package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "kuhakuanime.com/db/sqlc"
)

type CommentsInput struct {
	Content string `form:"content" binding:"required"`
}

func GetCommentsByEpisodeId(c *gin.Context) {

	idStr := c.Param("episodeId")
	commentId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	queries := db.New(db.DBPool)
	comments, err := queries.GetCommentsByEpisodeId(c, int32(commentId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func CreateComment(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("episodeId")
	commentId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var comment CommentsInput
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found"})
		return
	}

	_, err = queries.CreateComment(c, db.CreateCommentParams{
		EpisodeID: int32(commentId),
		UserID:    userID.(int32),
		Content:   comment.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully"})
}

func UpdateComment(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("commentId")
	commentId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var comment CommentsInput
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found"})
		return
	}
	oldComment, err := queries.GetComment(c, int32(commentId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if oldComment.UserID != userID.(int32) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to update this comment"})
		return
	}

	_, err = queries.UpdateComment(c, db.UpdateCommentParams{
		ID:        int32(commentId),
		Content:   comment.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

func DeleteComment(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("commentId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found"})
		return
	}
	oldComment, err := queries.GetComment(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if oldComment.UserID != userID.(int32) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to delete this comment"})
		return
	}

	err = queries.DeleteComment(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}