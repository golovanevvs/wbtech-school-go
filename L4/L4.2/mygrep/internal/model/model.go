package model

import (
	"io"
	"net"
	"sync"
	"time"
)

// Config contains application configuration
type Config struct {
	ServerID      string        `json:"server_id"` // server identifier
	Port          string        `json:"port"`
	Peers         []string      `json:"peers"` // addresses of other servers
	Pattern       string        `json:"pattern"`
	Files         []string      `json:"files"`          // list of files to process
	Input         io.Reader     `json:"-"`              // input reader (stdin by default)
	Output        io.Writer     `json:"-"`              // output writer (stdout by default)
	IsDistributed bool          `json:"is_distributed"` // distributed mode flag
	Flags         GrepFlags     `json:"flags"`          // grep flags
	Timeout       time.Duration `json:"timeout"`
	LocalAddress  *net.TCPAddr  `json:"-"`
}

// GrepResult represents a grep search result
type GrepResult struct {
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
	Match      string `json:"match"`
}

// JobResult represents job execution result
type JobResult struct {
	JobID       string       `json:"job_id"`
	ServerID    string       `json:"server_id"`
	Matches     []GrepResult `json:"matches"`   // found matches
	Processed   int          `json:"processed"` // number of processed lines
	Error       string       `json:"error,omitempty"`
	Success     bool         `json:"success"`
	CompletedAt time.Time    `json:"completed_at"`
}

// Job represents a task for processing data chunks
type Job struct {
	ID        string    `json:"id"`
	ServerID  string    `json:"server_id"`
	Pattern   string    `json:"pattern"`
	Data      string    `json:"data"`       // data for processing
	StartLine int       `json:"start_line"` // starting line
	EndLine   int       `json:"end_line"`   // ending line
	Flags     GrepFlags `json:"flags"`
	CreatedAt time.Time `json:"created_at"`
}

// Result represents task processing result
type Result struct {
	JobID       string    `json:"job_id"`
	ServerID    string    `json:"server_id"`
	Matches     []Match   `json:"matches"`   // found matches
	Processed   int       `json:"processed"` // number of processed lines
	Error       string    `json:"error,omitempty"`
	Success     bool      `json:"success"`
	CompletedAt time.Time `json:"completed_at"`
}

// Match represents a found match
type Match struct {
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
	Column     int    `json:"column,omitempty"` // position in line
}

// ServerInfo contains server state information
type ServerInfo struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	Status    string    `json:"status"` // "online", "offline", "processing"
	JobsCount int       `json:"jobs_count"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Message represents a network communication message
type Message struct {
	Type      string      `json:"type"` // "job_request", "job_response", "status_update", "result"
	From      string      `json:"from"`
	To        string      `json:"to,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NetworkMessage used for message serialization over network
type NetworkMessage struct {
	Type      string      `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// GrepFlags contains grep flags
type GrepFlags struct {
	Color        bool `json:"color"`         // --color
	InvertMatch  bool `json:"invert_match"`  // -v
	IgnoreCase   bool `json:"ignore_case"`   // -i
	WholeLine    bool `json:"whole_line"`    // -x
	LineNumber   bool `json:"line_number"`   // -n
	Count        bool `json:"count"`         // -c
	OnlyMatching bool `json:"only_matching"` // -o
}

// JobRequest structure for job execution request
type JobRequest struct {
	Job Job `json:"job"`
}

// JobResponse structure for result response
type JobResponse struct {
	Result Result `json:"result"`
}

// StatusUpdate structure for status update
type StatusUpdate struct {
	ServerInfo ServerInfo `json:"server_info"`
}

// QuorumStatus tracks quorum status
type QuorumStatus struct {
	TotalServers  int                   `json:"total_servers"`
	RequiredVotes int                   `json:"required_votes"` // N/2+1
	ReceivedVotes int                   `json:"received_votes"`
	Results       map[string]*JobResult `json:"results"` // serverID -> result
	Completed     bool                  `json:"completed"`
	mu            sync.RWMutex
}

// NewQuorumStatus creates a new QuorumStatus
func NewQuorumStatus(totalServers int) *QuorumStatus {
	return &QuorumStatus{
		TotalServers:  totalServers,
		RequiredVotes: totalServers/2 + 1,
		Results:       make(map[string]*JobResult),
		Completed:     false,
	}
}

// AddResult adds a result and checks for quorum achievement
func (q *QuorumStatus) AddResult(result *JobResult) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Results[result.ServerID] = result
	q.ReceivedVotes++

	// Check quorum achievement
	if q.ReceivedVotes >= q.RequiredVotes {
		q.Completed = true
		return true
	}

	return false
}

// IsCompleted checks if quorum is achieved
func (q *QuorumStatus) IsCompleted() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.Completed
}

// GetResults returns all results
func (q *QuorumStatus) GetResults() map[string]*JobResult {
	q.mu.RLock()
	defer q.mu.RUnlock()

	// Create copy for safety
	results := make(map[string]*JobResult)
	for k, v := range q.Results {
		results[k] = v
	}
	return results
}
