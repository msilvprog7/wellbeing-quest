package e2etests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api.wellbeingquest.app/internal/dtos"
	"api.wellbeingquest.app/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	services.LoadDotEnv(".env")
	dataHandler := services.GetDataHandler()
	return services.SetupRouter(dataHandler)
}

func TestAddActivityWhenNoEntries(t *testing.T) {
	// Arrange
	router := setupRouter()
	activity := dtos.Activity{
		Name: "Read",
		Feelings: []string{"Relaxed", "Accomplished"},
	}

	// Act
	w := httptest.NewRecorder()
	json, _ := json.Marshal(activity)
	req, _ := http.NewRequest("POST", "/activities/v1", strings.NewReader(string(json)))
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, 201, w.Code)
}
