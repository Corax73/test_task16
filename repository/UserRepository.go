package repository

import (
	"net/http"
	"strconv"
	"timeTracker/customDb"
	"timeTracker/customLog"
	"timeTracker/models"
	"timeTracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRepository struct {
	OriginalRep Repository
}

// NewTaskRepository returns a pointer to the initiated repository instance.
func NewUserRepository() *UserRepository {
	rep := UserRepository{
		OriginalRep: Repository{
			SomethingWrong: "try later",
			NoRecords:      "not found",
			LimitDefault:   5,
		},
	}
	return &rep
}

// Create based on the data from the request, if they are not empty, creates a User record and returns it and an empty error,
// otherwise returns an empty User instance and a non-empty error.
func (rep *UserRepository) Create(c *gin.Context) (*models.User, error) {
	var err error
	user := models.User{}
	resp := &user
	newId := uuid.New()
	name := c.DefaultPostForm("name", "")
	passportNumber, errNumber := strconv.Atoi(c.DefaultPostForm("passportNumber", "0"))
	passportSeries, errSeries := strconv.Atoi(c.DefaultPostForm("passportSeries", "0"))
	if errNumber == nil && errSeries == nil {
		if name != "" && passportNumber != 0 && passportSeries != 0 {
			database := customDb.GetConnect()
			user := models.User{ID: newId, Name: name, PassportNumber: passportNumber, PassportSeries: passportSeries}
			result := database.Create(&user) // pass pointer of data to Create
			if result.Error == nil {
				utils.GCRunAndPrintMemory()
				c.JSON(200, "created")
			} else {
				customLog.Logging(result.Error)
				utils.GCRunAndPrintMemory()
				c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			}
			resp = &user
		}
	} else {
		if errNumber != nil {
			err = errNumber
		} else {
			err = errSeries
		}
	}
	return resp, err
}
