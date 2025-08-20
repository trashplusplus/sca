package server

import (
	"sca/internal/database"
	"sca/internal/entity/cat"
	"sca/internal/entity/mission"
	"sca/internal/entity/target"
	"sca/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init(ip, port string) {

	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	catRepository := cat.NewRepository(db)
	missionRepository := mission.NewRepository(db)
	targetRepository := target.NewRepository(db)

	router := gin.Default()
	router.Use(middleware.Logger()) //мідлевейр для логнування

	{
		//хендлери в cat_handlers.go
		//todo: validate breed with catAPI | зробив
		cat := router.Group("/cat")
		cat.POST("/create", createCatHandler(catRepository))
		cat.GET("/list", listCatHandler(catRepository))
		cat.PUT("/update", updateCatHandler(catRepository))
		cat.DELETE("/delete/:id", deleteCatHandler(catRepository))
		cat.GET("/get/:id", getCatHandler(catRepository))

	}

	{
		target := router.Group("/target")
		target.POST("/complete/:id", completeTargetHandler(targetRepository))
		target.PUT("/update-notes/:id", updateNotesTargetHandler(targetRepository)) //оновлюється лише коли target або mission - false
	}

	{
		mission := router.Group("/mission")
		mission.POST("/create", createMissionHandler(missionRepository))
		mission.GET("/list", listMissionHandler(missionRepository))
		mission.PUT("/:id/assign-cat", assignСatMissionHandler(missionRepository)) //можливо переписати
		mission.GET("/get/:id", getMissionByIdHandler(missionRepository))
		mission.PUT("/complete/:id", completeMissionHandler(missionRepository)) //updateMission
		mission.DELETE("/delete/:id", deleteMissionHandler(missionRepository))
		mission.DELETE("/delete-target/:id", deleteTargetFromMissionHandler(missionRepository))
		mission.POST("/add-target/:id", addTargetToMission(missionRepository))
	}

	router.Run(ip + ":" + port)

}
