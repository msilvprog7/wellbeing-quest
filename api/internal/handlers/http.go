package handlers

import (
	"time"

	"api.wellbeingquest.app/internal/dtos"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostActivity(c *gin.Context) {
  // bind and validate activity
	var activity dtos.Activity

	if err := c.BindJSON(&activity); err != nil {
		c.IndentedJSON(http.StatusBadRequest, dtos.ErrorMessage{ Message: fmt.Sprintf("Post activity expected an activity, error: %s", err) })
		return
	}

	if activity.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, dtos.ErrorMessage{ Message: "Post activity expected an activity with non-empty 'name'" })
		return
	}

	if len(activity.Feelings) == 0 {
		c.IndentedJSON(http.StatusBadRequest, dtos.ErrorMessage{ Message: "Post activity expected an activity with non-empty 'feelings'" })
		return
	}

	// add entry
	activity.Created = time.Now()
	entry := AddActivity(activity)

	// set week
	activity.Week = entry.Week

	// return json
	c.IndentedJSON(http.StatusCreated, activity)
}
