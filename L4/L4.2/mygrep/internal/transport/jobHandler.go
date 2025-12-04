package transport

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/service"
)

// JobHandlerAdapter adapts GrepService to JobHandler interface
type JobHandlerAdapter struct {
	grepService *service.GrepService
}

// NewJobHandlerAdapter creates a new job handler adapter
func NewJobHandlerAdapter(grepService *service.GrepService) JobHandler {
	return &JobHandlerAdapter{
		grepService: grepService,
	}
}

// HandleJobRequest implements JobHandler interface
func (a *JobHandlerAdapter) HandleJobRequest(job map[string]interface{}) (*model.JobResult, error) {
	return a.grepService.ExecuteJobLocally(job)
}
