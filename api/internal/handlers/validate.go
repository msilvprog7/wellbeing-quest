package handlers

import (
	"api.wellbeingquest.app/internal/dtos"

	"errors"
)

func ValidateActivity(activity *dtos.Activity, previousErr error) error {
	if previousErr != nil {
		return previousErr
	}

	if activity.Name == "" {
		return errors.New("activity has an empty 'name'")
	}

	if len(activity.Feelings) == 0 {
		return errors.New("activity has empty 'feelings'")
	}

	return nil
}

func ValidateWeek(week *dtos.Week, previousErr error) error {
	if previousErr != nil {
		return previousErr
	}

	if week.Name == "" {
		return errors.New("week has an empty 'name'")
	}

	return nil
}