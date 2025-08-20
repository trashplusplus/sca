package server

import (
	"net/http"
	"sca/internal/entity/target"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotesUpdate struct {
	Notes string `json:"notes"`
}

func completeTargetHandler(repo *target.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(404, gin.H{"error": "invalid target id"})
			return
		}
		currentTarget, err := repo.GetTargetById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "target not found"})
			return
		}

		currentTarget.Complete = true //завершуємо ціль

		if err := repo.Update(currentTarget); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete target"})
			return
		}
		c.JSON(200, gin.H{"message": "Target completed successfully"})
	}
}

func updateNotesTargetHandler(repo *target.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(404, gin.H{"error": "invalid target id"})
			return
		}

		currentTarget, err := repo.GetTargetById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "target not found"})
			return
		}

		var notesUpdate NotesUpdate
		if err := c.ShouldBindJSON(&notesUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		currentTarget.Notes = notesUpdate.Notes

		if err := repo.Update(currentTarget); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update target (target or mission is completed)"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Target notes updated successfully"})
	}
}
