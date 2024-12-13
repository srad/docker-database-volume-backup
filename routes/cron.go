package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/srad/docker-database-volume-backup/services"
	"net/http"
	"time"
)

type CronJobDto struct {
	Id   int       `json:"id"`
	Next time.Time `json:"next"`
}

// GetCron godoc
// @Summary Get list of cron jobs
// @Description Retrieves a list of all active cron jobs.
// @Tags Cron Jobs
// @Accept json
// @Produce json
// @Success 200 {array} []CronJobDto "List of cron jobs"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /cron [get]
func GetCron(c echo.Context) error {
	jobs, err := services.GetCronJobs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var jobList []CronJobDto
	for _, job := range jobs {
		jobList = append(jobList, CronJobDto{
			Id:   int(job.Id),
			Next: job.Next,
		})
	}

	return c.JSON(http.StatusOK, jobList)
}
