package models

import (
	"queue-system-backend/database"
	"time"
)

type QueueStatistics struct {
	CounterID    uint      `json:"counter_id" gorm:"column:counter_id"`
	ServiceID    uint      `json:"service_id" gorm:"column:service_id"`
	VenueID      uint      `json:"venue_id" gorm:"column:venue_id"`
	ActiveQueues int       `json:"active_queues"`
	WaitTime     float64   `json:"wait_time"`
	TotalServed  int       `json:"total_served"`
	CreatedAt    time.Time `json:"created_at"`
}

type StatisticsFilter struct {
	CounterID uint
	ServiceID uint
	VenueID   uint
}

func (QueueStatistics) TableName() string {
	return "QueueTickets"
}

func (s *QueueStatistics) GetActiveQueues(filter StatisticsFilter) ([]QueueStatistics, error) {
	var stats []QueueStatistics
	query := database.DB.Table("QueueTickets").
		Select("counter_id, service_id, venue_id, COUNT(*) as active_queues").
		Where("status = ?", "waiting")

	if filter.VenueID != 0 {
		query = query.Where("venue_id = ?", filter.VenueID)
	}
	if filter.ServiceID != 0 {
		query = query.Where("service_id = ?", filter.ServiceID)
	}
	if filter.CounterID != 0 {
		query = query.Where("counter_id = ?", filter.CounterID)
	}

	query = query.Group("venue_id, service_id, counter_id")
	err := query.Find(&stats).Error
	return stats, err
}

func (s *QueueStatistics) GetAverageWaitTime(filter StatisticsFilter) ([]QueueStatistics, error) {
	var stats []QueueStatistics
	query := database.DB.Table("QueueTickets").
		Select("counter_id, service_id, venue_id, AVG(TIMESTAMPDIFF(MINUTE, created_at, called_at)) as wait_time").
		Where("status IN (?, ?) AND called_at IS NOT NULL", "called", "completed")

	if filter.VenueID != 0 {
		query = query.Where("venue_id = ?", filter.VenueID)
	}
	if filter.ServiceID != 0 {
		query = query.Where("service_id = ?", filter.ServiceID)
	}
	if filter.CounterID != 0 {
		query = query.Where("counter_id = ?", filter.CounterID)
	}

	query = query.Group("venue_id, service_id, counter_id")
	err := query.Find(&stats).Error
	return stats, err
}

func (s *QueueStatistics) GetTotalServed(filter StatisticsFilter) ([]QueueStatistics, error) {
	var stats []QueueStatistics
	query := database.DB.Table("QueueTickets").
		Select("counter_id, service_id, venue_id, COUNT(*) as total_served").
		Where("status = ? AND completed_at IS NOT NULL", "completed")

	if filter.VenueID != 0 {
		query = query.Where("venue_id = ?", filter.VenueID)
	}
	if filter.ServiceID != 0 {
		query = query.Where("service_id = ?", filter.ServiceID)
	}
	if filter.CounterID != 0 {
		query = query.Where("counter_id = ?", filter.CounterID)
	}

	query = query.Group("venue_id, service_id, counter_id")
	err := query.Find(&stats).Error
	return stats, err
}
