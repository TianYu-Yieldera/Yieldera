package bridge

import (
	"container/heap"
	"database/sql"
	"log"
	"sync"
	"time"

	"loyalty-points-system/internal/models"
)

// Priority levels for bridge messages
const (
	PriorityLow       = 0
	PriorityNormal    = 1
	PriorityHigh      = 2
	PriorityEmergency = 3
)

// MessageQueue manages pending bridge operations with priority support
type MessageQueue struct {
	database *sql.DB
	pq       *PriorityQueue
	mu       sync.RWMutex
	items    map[string]*QueueItem // messageHash -> item
}

// QueueItem represents a bridge message in the queue
type QueueItem struct {
	Message  *models.BridgeEvent
	Priority int
	AddedAt  time.Time
	index    int // Index in heap
}

// PriorityQueue implements heap.Interface
type PriorityQueue []*QueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Higher priority comes first
	if pq[i].Priority != pq[j].Priority {
		return pq[i].Priority > pq[j].Priority
	}
	// If same priority, earlier timestamp comes first (FIFO)
	return pq[i].AddedAt.Before(pq[j].AddedAt)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*QueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// NewMessageQueue creates a new message queue
func NewMessageQueue(database *sql.DB) *MessageQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return &MessageQueue{
		database: database,
		pq:       &pq,
		items:    make(map[string]*QueueItem),
	}
}

// LoadPendingMessages loads pending messages from database into queue
func (mq *MessageQueue) LoadPendingMessages() error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	query := `
		SELECT message_hash, direction, user_address, amount, status,
		       l1_tx_hash, l2_tx_hash, l1_block_number, l2_block_number,
		       EXTRACT(EPOCH FROM initiated_at)::bigint,
		       COALESCE(EXTRACT(EPOCH FROM confirmed_at)::bigint, 0),
		       retry_count, COALESCE(error_msg, '')
		FROM bridge_messages
		WHERE status IN ('initiated', 'pending')
		  AND retry_count < 10
		ORDER BY initiated_at ASC
	`

	rows, err := mq.database.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		evt := &models.BridgeEvent{}
		err := rows.Scan(
			&evt.MessageHash,
			&evt.Direction,
			&evt.UserAddress,
			&evt.Amount,
			&evt.Status,
			&evt.L1TxHash,
			&evt.L2TxHash,
			&evt.L1BlockNumber,
			&evt.L2BlockNumber,
			&evt.InitiatedAt,
			&evt.ConfirmedAt,
			&evt.RetryCount,
			&evt.ErrorMsg,
		)
		if err != nil {
			return err
		}

		// Determine priority based on direction and retry count
		priority := mq.determinePriority(evt)

		// Add to queue
		item := &QueueItem{
			Message:  evt,
			Priority: priority,
			AddedAt:  time.Now(),
		}

		heap.Push(mq.pq, item)
		mq.items[evt.MessageHash] = item
		count++
	}

	if count > 0 {
		log.Printf("📥 [Message Queue] Loaded %d pending messages", count)
	}

	return rows.Err()
}

// determinePriority assigns priority based on message characteristics
func (mq *MessageQueue) determinePriority(msg *models.BridgeEvent) int {
	// Emergency withdrawals (L2→L1 with high retry count) get highest priority
	if msg.Direction == "L2_TO_L1" && msg.RetryCount >= 5 {
		return PriorityEmergency
	}

	// L2→L1 withdrawals get high priority
	if msg.Direction == "L2_TO_L1" {
		return PriorityHigh
	}

	// Messages with retries get higher priority
	if msg.RetryCount > 0 {
		return PriorityHigh
	}

	// Normal L1→L2 deposits
	return PriorityNormal
}

// Enqueue adds a new message to the queue
func (mq *MessageQueue) Enqueue(msg *models.BridgeEvent, priority int) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	// Check if already in queue
	if _, exists := mq.items[msg.MessageHash]; exists {
		return
	}

	item := &QueueItem{
		Message:  msg,
		Priority: priority,
		AddedAt:  time.Now(),
	}

	heap.Push(mq.pq, item)
	mq.items[msg.MessageHash] = item

	log.Printf("📩 [Message Queue] Enqueued message %s (priority: %d)", msg.MessageHash, priority)
}

// Dequeue removes and returns the highest priority message
func (mq *MessageQueue) Dequeue() *models.BridgeEvent {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if mq.pq.Len() == 0 {
		return nil
	}

	item := heap.Pop(mq.pq).(*QueueItem)
	delete(mq.items, item.Message.MessageHash)

	return item.Message
}

// Peek returns the highest priority message without removing it
func (mq *MessageQueue) Peek() *models.BridgeEvent {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	if mq.pq.Len() == 0 {
		return nil
	}

	return (*mq.pq)[0].Message
}

// Remove removes a specific message from the queue
func (mq *MessageQueue) Remove(messageHash string) bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	item, exists := mq.items[messageHash]
	if !exists {
		return false
	}

	heap.Remove(mq.pq, item.index)
	delete(mq.items, messageHash)

	log.Printf("🗑️  [Message Queue] Removed message %s", messageHash)
	return true
}

// UpdatePriority changes the priority of a message in the queue
func (mq *MessageQueue) UpdatePriority(messageHash string, newPriority int) bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	item, exists := mq.items[messageHash]
	if !exists {
		return false
	}

	item.Priority = newPriority
	heap.Fix(mq.pq, item.index)

	log.Printf("🔄 [Message Queue] Updated priority for %s to %d", messageHash, newPriority)
	return true
}

// Size returns the number of messages in the queue
func (mq *MessageQueue) Size() int {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	return mq.pq.Len()
}

// IsEmpty returns true if queue is empty
func (mq *MessageQueue) IsEmpty() bool {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	return mq.pq.Len() == 0
}

// Contains checks if a message is in the queue
func (mq *MessageQueue) Contains(messageHash string) bool {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	_, exists := mq.items[messageHash]
	return exists
}

// GetQueueSnapshot returns a snapshot of all messages in the queue
func (mq *MessageQueue) GetQueueSnapshot() []*QueueItem {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	snapshot := make([]*QueueItem, len(*mq.pq))
	for i, item := range *mq.pq {
		snapshot[i] = &QueueItem{
			Message:  item.Message,
			Priority: item.Priority,
			AddedAt:  item.AddedAt,
		}
	}

	return snapshot
}

// GetByPriority returns all messages with a specific priority
func (mq *MessageQueue) GetByPriority(priority int) []*models.BridgeEvent {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	var messages []*models.BridgeEvent
	for _, item := range *mq.pq {
		if item.Priority == priority {
			messages = append(messages, item.Message)
		}
	}

	return messages
}

// Clear removes all messages from the queue
func (mq *MessageQueue) Clear() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.pq = &PriorityQueue{}
	heap.Init(mq.pq)
	mq.items = make(map[string]*QueueItem)

	log.Println("🧹 [Message Queue] Cleared all messages")
}

// GetStats returns queue statistics
func (mq *MessageQueue) GetStats() *QueueStats {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	stats := &QueueStats{
		TotalMessages: mq.pq.Len(),
		ByPriority:    make(map[int]int),
	}

	for _, item := range *mq.pq {
		stats.ByPriority[item.Priority]++

		// Track oldest message
		if stats.OldestMessage.IsZero() || item.AddedAt.Before(stats.OldestMessage) {
			stats.OldestMessage = item.AddedAt
		}
	}

	return stats
}

// QueueStats contains queue statistics
type QueueStats struct {
	TotalMessages  int
	ByPriority     map[int]int
	OldestMessage  time.Time
}

// PriorityName returns human-readable priority name
func PriorityName(priority int) string {
	switch priority {
	case PriorityEmergency:
		return "Emergency"
	case PriorityHigh:
		return "High"
	case PriorityNormal:
		return "Normal"
	case PriorityLow:
		return "Low"
	default:
		return "Unknown"
	}
}
