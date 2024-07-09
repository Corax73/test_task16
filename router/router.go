package router

import (
	"timeTracker/repository"
	"timeTracker/utils"

	"github.com/gin-gonic/gin"
)

func RunRouter() {
	utils.GCRunAndPrintMemory()
	router := gin.Default()
	router.GET("/users", getList)
	router.GET("/tasks", getList)
	router.POST("/tasks/start", startTask)
	router.POST("/tasks/stop", stopTask)
	router.POST("/tasks/complete", completeTask)
	router.Run(":4343")
}

// getList processes the route for obtaining lists of entities.
func getList(c *gin.Context) {
	rep := repository.NewTaskRepository()
	rep.GetList(c)
}

// startTask processes the task start route.
func startTask(c *gin.Context) {
	rep := repository.NewTaskRepository()
	rep.StartTask(c)
}

// stopTask handles the task pause route.
func stopTask(c *gin.Context) {
	rep := repository.NewTaskRepository()
	rep.StopTask(c)
}

// completeTask processes the task completion route.
func completeTask(c *gin.Context) {
	rep := repository.NewTaskRepository()
	rep.CompleteTask(c)
}
