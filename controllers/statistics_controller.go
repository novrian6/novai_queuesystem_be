package controllers

import (
	"net/http"
	"queue-system-backend/models"

	"github.com/gin-gonic/gin"
)

type StatisticsController struct{}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{}
}

type StatisticsRequest struct {
	CounterID uint `json:"counter_id"`
	ServiceID uint `json:"service_id"`
	VenueID   uint `json:"venue_id"`
}

// getStatsFilter now as method
func (sc *StatisticsController) getStatsFilter(c *gin.Context) models.StatisticsFilter {
	var req StatisticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return models.StatisticsFilter{}
	}

	//claims := c.MustGet("claims").(*utils.Claims)

	filter := models.StatisticsFilter{
		CounterID: req.CounterID,
		ServiceID: req.ServiceID,
		VenueID:   req.VenueID,
	}

	// Non-admin users can only see their venue's data
	//if !strings.EqualFold(claims.Role, "admin") {
	//	filter.VenueID = claims.VenueID
	//}

	return filter
}

// GetActiveQueues now as method
func (sc *StatisticsController) GetActiveQueues(c *gin.Context) {
	stats := &models.QueueStatistics{}
	param := sc.getStatsFilter(c)
	results, err := stats.GetActiveQueues(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch active queues statistics",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, results)
}

// GetAverageWaitTime now as method
func (sc *StatisticsController) GetAverageWaitTime(c *gin.Context) {
	stats := &models.QueueStatistics{}
	results, err := stats.GetAverageWaitTime(sc.getStatsFilter(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch wait time statistics",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, results)
}

// GetTotalServed now as method
func (sc *StatisticsController) GetTotalServed(c *gin.Context) {
	stats := &models.QueueStatistics{}
	results, err := stats.GetTotalServed(sc.getStatsFilter(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch total served statistics",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, results)
}
