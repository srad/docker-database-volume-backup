package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"net/http"
	"time"
)

type CronEntry struct {
	Id   cron.EntryID `json:"id"`
	Next time.Time    `json:"next"`
}

func GetCron(c echo.Context, cron *cron.Cron) error {
	return c.JSON(http.StatusOK, getCronJobs(cron))
}

func getCronJobs(c *cron.Cron) []CronEntry {
	entries := c.Entries()
	var result = make([]CronEntry, len(entries))
	for _, entry := range entries {
		result = append(result, CronEntry{
			Id:   entry.ID,
			Next: entry.Next,
		})
	}

	return result
}
