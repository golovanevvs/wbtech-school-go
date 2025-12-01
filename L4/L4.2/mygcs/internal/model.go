package model

import (
	"net"
	"time"
)

// Config содержит конфигурацию сервера
type Config struct {
	ServerID     string          `json:"server_id"`
	Port         string          `json:"port"`
	Peers        []string        `json:"peers"` // адреса других серверов
	Pattern      string          `json:"pattern"`
	InputFile    string          `json:"input_file"`
	Flags        map[string]bool `json:"flags"` // флаги grep: color, invert-match, etc.
	Timeout      time.Duration   `json:"timeout"`
	LocalAddress *net.TCPAddr    `json:"-"`
}

// Job представляет задание для обработки части данных
type Job struct {
	ID        string          `json:"id"`
	ServerID  string          `json:"server_id"`
	Pattern   string          `json:"pattern"`
	Data      string          `json:"data"`       // данные для обработки
	StartLine int             `json:"start_line"` // начальная строка
	EndLine   int             `json:"end_line"`   // конечная строка
	Flags     map[string]bool `json:"flags"`
	CreatedAt time.Time       `json:"created_at"`
}

// Result представляет результат обработки задания
type Result struct {
	JobID       string    `json:"job_id"`
	ServerID    string    `json:"server_id"`
	Matches     []Match   `json:"matches"`   // найденные совпадения
	Processed   int       `json:"processed"` // количество обработанных строк
	Error       string    `json:"error,omitempty"`
	Success     bool      `json:"success"`
	CompletedAt time.Time `json:"completed_at"`
}

// Match представляет найденное совпадение
type Match struct {
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
	Column     int    `json:"column,omitempty"` // позиция в строке
}

// ServerInfo содержит информацию о состоянии сервера
type ServerInfo struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	Status    string    `json:"status"` // "online", "offline", "processing"
	JobsCount int       `json:"jobs_count"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Message представляет сообщение для сетевого обмена
type Message struct {
	Type      string      `json:"type"` // "job_request", "job_response", "status_update", "result"
	From      string      `json:"from"`
	To        string      `json:"to,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NetworkMessage используется для сериализации сообщений по сети
type NetworkMessage struct {
	Type      string      `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// GrepFlags содержит флаги для grep
type GrepFlags struct {
	Color        bool `json:"color"`         // --color
	InvertMatch  bool `json:"invert_match"`  // -v
	IgnoreCase   bool `json:"ignore_case"`   // -i
	WholeLine    bool `json:"whole_line"`    // -x
	LineNumber   bool `json:"line_number"`   // -n
	Count        bool `json:"count"`         // -c
	OnlyMatching bool `json:"only_matching"` // -o
}

// JobRequest структура запроса на выполнение задания
type JobRequest struct {
	Job Job `json:"job"`
}

// JobResponse структура ответа с результатом
type JobResponse struct {
	Result Result `json:"result"`
}

// StatusUpdate структура обновления статуса
type StatusUpdate struct {
	ServerInfo ServerInfo `json:"server_info"`
}

// QuorumStatus отслеживает статус кворума
type QuorumStatus struct {
	TotalServers  int               `json:"total_servers"`
	RequiredVotes int               `json:"required_votes"` // N/2+1
	ReceivedVotes int               `json:"received_votes"`
	Results       map[string]Result `json:"results"` // serverID -> result
	Completed     bool              `json:"completed"`
}
