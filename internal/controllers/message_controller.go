package controllers

import (
	"fmt"
	"net/http"

	//  "sit-iot-message-mng-api/internal/middleware"
	"sit-iot-message-mng-api/internal/services"
	"sit-iot-message-mng-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	MessageService services.MessageService
}

func NewMessageController(messageService services.MessageService) *MessageController {
	return &MessageController{
		MessageService: messageService,
	}
}

func (mc *MessageController) GetMessage(c *gin.Context) {
	id := c.Param("id")

	message, err := mc.MessageService.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (mc *MessageController) ListMessagesByDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")

	// Parse query params with defaults
	rangeParam := c.DefaultQuery("range", "[0,9]")
	sortParam := c.DefaultQuery("sort", `["timestamp","DESC"]`)

	// Parse range
	var rangeArr [2]int
	if err := utils.ParseJSON(rangeParam, &rangeArr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range parameter"})
		return
	}
	skip := rangeArr[0]
	limit := rangeArr[1] - rangeArr[0] + 1

	// Parse sort
	var sortArr [2]string
	if err := utils.ParseJSON(sortParam, &sortArr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort parameter"})
		return
	}
	sortField := sortArr[0]
	sortOrder := sortArr[1]

	messages, total, err := mc.MessageService.ListMessagesByDeviceID(c.Request.Context(), deviceID, nil, sortField, sortOrder, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set Content-Range header
	end := skip + len(messages) - 1
	if len(messages) == 0 {
		end = skip - 1
	}
	contentRange := fmt.Sprintf("items %d-%d/%d", skip, end, total)
	c.Header("Content-Range", contentRange)
	c.JSON(http.StatusOK, messages)
}

// GetAggregatedDataByDevice returns aggregated data for a device for graphing max, min, avg
func (mc *MessageController) GetAggregatedDataByDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")

	aggregations, err := mc.MessageService.GetAggregatedDataByDeviceID(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"device_id":    deviceID,
		"aggregations": aggregations,
	}
	c.JSON(http.StatusOK, response)
}
