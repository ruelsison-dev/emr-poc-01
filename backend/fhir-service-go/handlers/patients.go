package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NOTE: This scaffold uses an in-memory map. Phase 1 will wire DynamoDB persistence.
var patients = map[string]map[string]interface{}{}

func createPatient(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_payload"})
		return
	}
	id := "pat_" + fmt.Sprint(time.Now().UnixNano())
	payload["id"] = id
	patients[id] = payload
	c.JSON(http.StatusCreated, payload)
}

func getPatient(c *gin.Context) {
	id := c.Param("id")
	p, ok := patients[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
