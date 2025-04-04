package controllers

import (
	"net/http"
	"strconv"

	"queue-system-backend/models"

	"github.com/gin-gonic/gin"
)

func GetQueueTicketByTokenHandler(c *gin.Context) {
	token := c.Param("token")

	ticket, err := models.GetQueueTicketByToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serviceName, err := models.GetServiceNameByID(*ticket.ServiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	venueName, err := models.GetVenueNameByID(*ticket.VenueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	averageQueueTime, err := models.CalculateAverageQueuingTime(*ticket.VenueID, *ticket.ServiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ticket":             ticket,
		"service_name":       serviceName,
		"venue_name":         venueName,
		"average_queue_time": averageQueueTime,
	})
}

func GetWaitingQueueTicketsHandler(c *gin.Context) {
	venueID := c.Query("venue_id")
	serviceID := c.Query("service_id")

	// Convert venueID and serviceID to uint
	venueIDUint, err := strconv.ParseUint(venueID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid venue_id"})
		return
	}

	serviceIDUint, err := strconv.ParseUint(serviceID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service_id"})
		return
	}

	/*waitingTickets, err := models.GetQueueTicketsSorted("waiting", uint(venueIDUint), uint(serviceIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	calledTickets, err := models.GetQueueTicketsSorted("called", uint(venueIDUint), uint(serviceIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	completedTickets, err := models.GetQueueTicketsSorted("all", uint(venueIDUint), uint(serviceIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}



	tickets := append(completedTickets, waitingTickets...)
	tickets = append(tickets, calledTickets...)
	*/

	tickets, err := models.GetQueueTicketsSorted("all", uint(venueIDUint), uint(serviceIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}
