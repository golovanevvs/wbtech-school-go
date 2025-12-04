package transport

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/model"
)

// Server represents a TCP server for distributed communication
type Server struct {
	config      *model.Config
	listener    net.Listener
	servers     map[string]*model.ServerInfo
	serverMutex sync.RWMutex
	quit        chan bool
	jobHandler  JobHandler
}

// JobHandler interface for handling job requests
type JobHandler interface {
	HandleJobRequest(job map[string]interface{}) (*model.JobResult, error)
}

// NewServer creates a new TCP server
func NewServer(config *model.Config, jobHandler JobHandler) *Server {
	return &Server{
		config:     config,
		servers:    make(map[string]*model.ServerInfo),
		quit:       make(chan bool),
		jobHandler: jobHandler,
	}
}

// Start starts the TCP server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.config.Port)
	if err != nil {
		return err
	}
	s.listener = listener

	go s.acceptConnections()
	fmt.Printf("TCP server started on port %s\n", s.config.Port)
	return nil
}

// Stop stops the TCP server
func (s *Server) Stop() error {
	close(s.quit)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// acceptConnections accepts incoming TCP connections
func (s *Server) acceptConnections() {
	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				continue
			}
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a single TCP connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	var msg model.Message
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&msg); err != nil {
		return
	}

	switch msg.Type {
	case "status_request":
		s.sendStatusResponse(conn, msg.From)
	case "job_request":
		s.handleJobRequest(conn, msg)
	default:
		fmt.Printf("Received unknown message type: %s\n", msg.Type)
	}
}

// sendStatusResponse sends status response to client
func (s *Server) sendStatusResponse(conn net.Conn, from string) {
	status := map[string]interface{}{
		"server_id": s.config.ServerID,
		"status":    "online",
		"timestamp": time.Now(),
	}

	response := model.Message{
		Type:      "status_response",
		From:      s.config.ServerID,
		To:        from,
		Data:      status,
		Timestamp: time.Now(),
	}

	json.NewEncoder(conn).Encode(response)
}

// handleJobRequest handles job request from client
func (s *Server) handleJobRequest(conn net.Conn, msg model.Message) {
	fmt.Printf("Received job request from %s\n", msg.From)

	jobData, ok := msg.Data.(map[string]interface{})
	if !ok {
		fmt.Println("Invalid job data format")
		s.sendErrorResponse(conn, msg.From, "Invalid job data format")
		return
	}

	// Execute job using the job handler
	result, err := s.jobHandler.HandleJobRequest(jobData)
	if err != nil {
		fmt.Printf("Job execution error: %v\n", err)
		s.sendErrorResponse(conn, msg.From, err.Error())
		return
	}

	// Send result back
	s.sendJobResponse(conn, msg.From, result)
}

// sendJobResponse sends job result response
func (s *Server) sendJobResponse(conn net.Conn, to string, result *model.JobResult) {
	response := model.Message{
		Type:      "job_response",
		From:      s.config.ServerID,
		To:        to,
		Data:      result,
		Timestamp: time.Now(),
	}

	if err := json.NewEncoder(conn).Encode(response); err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	}
}

// sendErrorResponse sends error response
func (s *Server) sendErrorResponse(conn net.Conn, to string, errorMsg string) {
	result := &model.JobResult{
		JobID:       "",
		ServerID:    s.config.ServerID,
		Processed:   0,
		Success:     false,
		Error:       errorMsg,
		CompletedAt: time.Now(),
	}

	s.sendJobResponse(conn, to, result)
}
