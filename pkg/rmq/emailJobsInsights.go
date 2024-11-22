package rmq

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type QueueStatus struct {
	ReadyCount     int64  `json:"ready_count"`
	UnackedCount   int64  `json:"unacked_count"`
	ConsumerCount  int    `json:"consumer_count"`
	RejectedCount  int64  `json:"rejected_count"`
	LastDeliveries string `json:"last_delivery_time,omitempty"`
}

var (
	queueStatus   *QueueStatus
	statusMutex   sync.RWMutex
	refreshTicker *time.Ticker
)

// PollQueueStats refreshes stats every interval
func PollQueueStats(interval time.Duration) {
	refreshTicker = time.NewTicker(interval)
	go func() {
		for range refreshTicker.C {
			status, err := collectQueueStats()
			if err != nil {
				log.Printf("Error collecting queue stats: %v", err)
				continue
			}
			statusMutex.Lock()
			queueStatus = status
			statusMutex.Unlock()
		}
	}()
}

// collectQueueStats fetches live stats from the rmqConnection
func collectQueueStats() (*QueueStatus, error) {
	if rmqConnection == nil {
		return nil, fmt.Errorf("email queue is not initialized")
	}

	stats, err := rmqConnection.CollectStats([]string{"email-queue"})
	if err != nil {
		return nil, fmt.Errorf("failed to collect stats: %v", err)
	}

	queueStats, exists := stats.QueueStats["email-queue"]
	if !exists {
		return nil, fmt.Errorf("no stats found for email-queue")
	}

	return &QueueStatus{
		ReadyCount:    queueStats.ReadyCount,
		UnackedCount:  queueStats.UnackedCount(),
		ConsumerCount: int(queueStats.ConsumerCount()),
		RejectedCount: queueStats.RejectedCount,
	}, nil
}

// GetEmailQueueStatus reads the latest stats from memory
func GetEmailQueueStatus() (*QueueStatus, error) {
	statusMutex.RLock()
	defer statusMutex.RUnlock()
	if queueStatus == nil {
		return nil, fmt.Errorf("queue stats not available yet")
	}
	return queueStatus, nil
}
