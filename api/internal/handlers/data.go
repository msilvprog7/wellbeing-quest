package handlers

import (
	"time"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"
)

type DataHandler interface {
	AddActivity(activity *dtos.Activity) (models.Entry, error)
	GetWeekAndActivities(week *dtos.Week) (models.Week, []models.Entry, error)
  GetActivitiesAndFeelings() ([]models.Activity, []models.Feeling, error)
}

// Helper functions related to representing data
func getWeek(time *time.Time) string {
	// Calculate the start of the week (Sunday)
	weekday := int(time.Weekday())
	startOfWeek := time.AddDate(0, 0, -weekday)
	return startOfWeek.Format("2006-01-02")
}

func getTime(week string) (time.Time, error) {
	return time.Parse("2006-01-02", week)
}
