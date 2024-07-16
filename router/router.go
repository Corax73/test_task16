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
	router.POST("/users", create)
	router.PUT("users", update)
	router.DELETE("/users/:id", delete)
	router.POST("/users/time", getExecTime)
	router.GET("/tasks", getList)
	router.PUT("tasks", update)
	router.POST("/tasks/start", startTask)
	router.POST("/tasks/stop", stopTask)
	router.POST("/tasks/complete", completeTask)
	router.Run(":4343")
}

// getList processes the route for obtaining lists of entities.
func getList(c *gin.Context) {
	rep := repository.NewTaskRepository()
	rep.OriginalRep.GetList(c)
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

// createUser handles the user creation route.
func create(c *gin.Context) {
	rep := repository.NewUserRepository()
	rep.Create(c)
}

// getExecTime returns a slice of data with task IDs and their execution time, in descending order of time.
func getExecTime(c *gin.Context) {
	rep := repository.NewUserRepository()
	rep.GetTaskExecutionTime(c)
}

// delete deletes an entity using the passed ID.
func delete(c *gin.Context) {
	rep := repository.NewRepository()
	rep.Delete(c)
}

// update processes the entity data update route.
func update(c *gin.Context) {
	rep := repository.NewRepository()
	rep.Update(c)
}
