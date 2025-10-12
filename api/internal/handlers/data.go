package handlers

import (
	"time"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/models"
)

type DataHandler interface {
	AddActivity(activity *dtos.Activity, user *dtos.User) (models.Entry, error)
	GetWeekAndActivities(week *dtos.Week, user *dtos.User) (models.Week, []models.Entry, error)
  GetActivitiesAndFeelings(user *dtos.User) ([]models.Activity, []models.Feeling, error)
}

// Helper functions related to representing data
func GetWeek(time *time.Time) string {
	// Calculate the start of the week (Sunday)
	weekday := int(time.Weekday())
	startOfWeek := time.AddDate(0, 0, -weekday)
	return startOfWeek.Format("2006-01-02")
}

func GetTime(week string) (time.Time, error) {
	return time.Parse("2006-01-02", week)
}
