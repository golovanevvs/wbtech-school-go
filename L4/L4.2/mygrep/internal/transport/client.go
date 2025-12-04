package transport

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/model"
)

// Client represents a network client for peer communication
type Client struct {
	config *model.Config
}

// NewClient creates a new network client
func NewClient(config *model.Config) *Client {
	return &Client{
		config: config,
	}
}

// SendJobToPeer sends a job to a specific peer and returns the result
func (c *Client) SendJobToPeer(peerAddr string, job map[string]interface{}) (*model.JobResult, error) {
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Create job request message
	msg := model.Message{
		Type:      "job_request",
		From:      c.config.ServerID,
		To:        peerAddr,
		Data:      job,
		Timestamp: time.Now(),
	}

	// Send message
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(msg); err != nil {
		return nil, err
	}

	// Read response
	var response model.Message
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	if response.Type != "job_response" {
		return nil, fmt.Errorf("unexpected response from %s: type=%s", peerAddr, response.Type)
	}

	// Convert response to JobResult
	return c.convertToJobResult(response.Data, peerAddr)
}

// SendJobsToPeers sends jobs to multiple peers with quorum tracking
func (c *Client) SendJobsToPeers(jobs []map[string]interface{}) (*QuorumResult, error) {
	numPeers := len(c.config.Peers)
	if numPeers == 0 {
		return nil, fmt.Errorf("no peers configured")
	}

	// Create quorum (all servers: current + peers)
	totalServers := numPeers + 1
	quorum := model.NewQuorumStatus(totalServers)

	fmt.Printf("Starting quorum: %d servers, need %d votes\n",
		totalServers, quorum.RequiredVotes)

	var wg sync.WaitGroup
	errors := make(chan error, totalServers)

	// Send jobs to peers in parallel
	for i, peer := range c.config.Peers {
		wg.Add(1)
		go func(peerAddr string, jobIndex int) {
			defer wg.Done()

			result, err := c.SendJobToPeer(peerAddr, jobs[jobIndex])
			if err != nil {
				errors <- fmt.Errorf("error sending job to peer %s: %v", peerAddr, err)
				return
			}

			// Add result to quorum
			if quorum.AddResult(result) {
				fmt.Println("Quorum achieved!")
			}
		}(peer, i)
	}

	// Wait for all jobs to complete or quorum to be achieved
	go func() {
		wg.Wait()
		close(errors)
	}()

	// Wait for quorum or completion
	return c.waitForQuorum(quorum, errors)
}

// ConnectToPeers connects to all configured peers
func (c *Client) ConnectToPeers() error {
	for _, peer := range c.config.Peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", peer, err)
			continue
		}

		// Send status request
		msg := model.Message{
			Type:      "status_request",
			From:      c.config.ServerID,
			To:        peer,
			Data:      nil,
			Timestamp: time.Now(),
		}

		json.NewEncoder(conn).Encode(msg)
		conn.Close()
	}
	return nil
}

// convertToJobResult converts response data to JobResult
func (c *Client) convertToJobResult(data interface{}, serverID string) (*model.JobResult, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid result data format")
	}

	jobID, _ := dataMap["job_id"].(string)
	processed, _ := dataMap["processed"].(float64)
	success, _ := dataMap["success"].(bool)
	completedAt, _ := dataMap["completed_at"].(time.Time)

	// Convert matches (placeholder - currently empty array)
	matches := make([]model.GrepResult, 0)

	result := &model.JobResult{
		JobID:       jobID,
		ServerID:    serverID,
		Matches:     matches,
		Processed:   int(processed),
		Success:     success,
		CompletedAt: completedAt,
	}

	if !success {
		result.Error, _ = dataMap["error"].(string)
	}

	return result, nil
}

// waitForQuorum waits for quorum achievement or completion
func (c *Client) waitForQuorum(quorum *model.QuorumStatus, errors <-chan error) (*QuorumResult, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if quorum.IsCompleted() {
				fmt.Println("Quorum achieved, combining results...")
				return c.combineResults(quorum.GetResults())
			}
		case err := <-errors:
			fmt.Printf("Execution error: %v\n", err)
		case <-time.After(30 * time.Second): // 30 second timeout
			fmt.Println("Quorum timeout")
			return c.combineResults(quorum.GetResults())
		}
	}
}

// combineResults combines results from all servers
func (c *Client) combineResults(results map[string]*model.JobResult) (*QuorumResult, error) {
	fmt.Printf("Combining results from %d servers\n", len(results))

	totalMatches := 0
	totalProcessed := 0
	allSuccessful := true

	for serverID, result := range results {
		fmt.Printf("Server %s: processed %d lines, found %d matches\n",
			serverID, result.Processed, len(result.Matches))

		totalProcessed += result.Processed
		totalMatches += len(result.Matches)

		if !result.Success {
			allSuccessful = false
			fmt.Printf("Error on server %s: %s\n", serverID, result.Error)
		}
	}

	quorumResult := &QuorumResult{
		TotalMatches:   totalMatches,
		TotalProcessed: totalProcessed,
		AllSuccessful:  allSuccessful,
		Results:        results,
	}

	if allSuccessful {
		fmt.Printf("Success: processed %d lines, found %d matches\n", totalProcessed, totalMatches)
	} else {
		fmt.Printf("Partial success: processed %d lines, found %d matches\n", totalProcessed, totalMatches)
	}

	return quorumResult, nil
}

// QuorumResult represents the combined result from quorum processing
type QuorumResult struct {
	TotalMatches   int
	TotalProcessed int
	AllSuccessful  bool
	Results        map[string]*model.JobResult
}
