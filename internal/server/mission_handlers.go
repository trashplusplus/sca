package server

import (
	"net/http"
	"sca/internal/entity/mission"
	"sca/internal/entity/target"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateMissionRequest struct {
	CatId    int             `json:"cat_id"`
	Complete bool            `json:"complete"`
	Targets  []target.Target `json:"targets"`
}

type AssignCatRequest struct {
	CatId int `json:"cat_id"`
}

type DeleteTargetFromMissionRequest struct {
	TargetId int `json:"target_id"`
}

func createMissionHandler(repo *mission.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateMissionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}
		//формуємо місію без кота, назначимо його пізніше
		mission := &mission.Mission{
			CatId:    request.CatId,
			Complete: request.Complete,
			Targets:  request.Targets,
		}

		if err := repo.Create(mission); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, mission)

	}
}

func listMissionHandler(repo *mission.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		missions, err := repo.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, missions)
	}
}

func assignСatMissionHandler(repo *mission.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Беремо ID місії з URL параметра
		missionIDStr := c.Param("id")
		missionID, err := strconv.Atoi(missionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission id"})
			return
		}

		// Читаємо body JSON
		var req AssignCatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body reading error"})
			return
		}

		// Викликаємо репозиторій, щоб оновити cat_id місії
		if err := repo.AssignCat(missionID, req.CatId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "assign error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cat assigned successfully"})
	}
}

func deleteTargetFromMissionHandler(repo *mission.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		missionIdStr := c.Param("id")
		missionId, err := strconv.Atoi(missionIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid mission id"})
			return
		}

		var request DeleteTargetFromMissionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		_, err = repo.GetMissionById(missionId)
		if err != nil {
			c.JSON(404, gin.H{"error": "target not found"})
			return
		}

		if err := repo.DeleteTarget(request.TargetId); err != nil {
			c.JSON(400, gin.H{"error": "failed to delete target from mission"})
			return
		}

		c.JSON(200, gin.H{"message": "Target deleted from mission successfully"})
	}
}

func getMissionByIdHandler(repo *mission.Repository) gin.HandlerFunc {

	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid mission id"})
			return
		}

		mission, err := repo.GetMissionById(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "mission not found"})
			return
		}

		c.JSON(200, mission)

	}

}

func completeMissionHandler(repo *mission.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid mission id"})
			return
		}

		_, err = repo.GetMissionById(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "mission not found"})
			return
		}

		if err := repo.UpdateMissionComplete(id, true); err != nil {
			c.JSON(500, gin.H{"error": "failed to complete mission"})
			return
		}

		c.JSON(200, gin.H{"message": "Mission completed successfully"})
	}
}

func deleteMissionHandler(repo *mission.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid mission id"})
			return
		}

		_, err = repo.GetMissionById(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "mission not found"})
			return
		}

		if err := repo.DeleteMission(id); err != nil {
			c.JSON(400, gin.H{"error": "failed to delete mission (cat is assigned)"})
			return
		}

		c.JSON(200, gin.H{"message": "Mission deleted succesfully"})

	}
}

func addTargetToMission(repo *mission.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		missionIdStr := c.Param("id")
		missionId, err := strconv.Atoi(missionIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid mission id"})
			return
		}

		var request target.Target
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		targetMission, err := repo.GetMissionById(missionId)
		if err != nil {
			c.JSON(404, gin.H{"error": "mission not found"})
			return
		}

		if err := repo.CreateTarget(targetMission.Id, request); err != nil {
			c.JSON(400, gin.H{"error": "failed to add target to mission"})
			return
		}

		c.JSON(200, gin.H{"message": "Target added to mission successfully"})
	}
}
