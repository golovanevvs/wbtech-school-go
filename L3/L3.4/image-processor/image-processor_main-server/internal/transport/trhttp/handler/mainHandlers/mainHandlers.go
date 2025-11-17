package mainHandlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgErrors"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	AddComment(ctx context.Context, comment *model.Comment) error
	GetCommentsTree(ctx context.Context, parentID *int) ([]*model.Comment, error)
	RemoveComment(ctx context.Context, id int) error
	FindComments(ctx context.Context, query string) ([]*model.Comment, error)
}

type MainHandlers struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService) *MainHandlers {
	lg := parentLg.With().Str("component", "AddShortURL").Logger()
	return &MainHandlers{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *MainHandlers) RegisterRoutes() {
	hd.rt.POST("/comments", hd.AddComment)
	hd.rt.GET("/comments", hd.GetCommentsTree)
	hd.rt.DELETE("/comments/:id", hd.RemoveComment)
	hd.rt.GET("/comments/search", hd.FindComments)
}

func (hd *MainHandlers) AddComment(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "AddComment").Logger()

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Err(pkgErrors.ErrContentTypeAJ).Str("Content-Type", c.ContentType()).Int("status", http.StatusBadRequest).Msgf("%s invalid content-type", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, response{Error: pkgErrors.ErrContentTypeAJ.Error()})
		return
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Int("status", http.StatusBadRequest).Msgf("%s failed to bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, response{Error: err.Error()})
		return
	}

	comment := &model.Comment{
		ParentID:  req.ParentID,
		Text:      req.Text,
		CreatedAt: time.Now(),
	}

	if err := hd.sv.AddComment(c.Request.Context(), comment); err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to add comment", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	resp := convertToResponse(comment)
	c.JSON(http.StatusCreated, resp)

	lg.Debug().Str("text", resp.Text).Msgf("%s comment added successfully", pkgConst.OpSuccess)
}

func (hd *MainHandlers) GetCommentsTree(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetCommentsTree").Logger()

	parentStr := c.Query("parent")
	var parentID *int
	if parentStr != "" {
		id, err := strconv.Atoi(parentStr)
		if err != nil {
			lg.Warn().Err(err).Str("param", parentStr).Int("status", http.StatusBadRequest).Msgf("%s invalid parent ID", pkgConst.Warn)
			c.JSON(http.StatusBadRequest, response{Error: "invalid parent ID"})
			return
		}
		parentID = &id
	}

	comments, err := hd.sv.GetCommentsTree(c.Request.Context(), parentID)
	if err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to get comments tree", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	responses := convertToResponseList(comments)
	c.JSON(http.StatusOK, responses)

	lg.Debug().Msgf("%s comments got successfully", pkgConst.OpSuccess)

}

func (hd *MainHandlers) RemoveComment(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "RemoveComment").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("param", c.Param("id")).Int("status", http.StatusBadRequest).Msgf("%s invalid comment ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, response{Error: "invalid comment ID"})
		return
	}

	if err := hd.sv.RemoveComment(c.Request.Context(), id); err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to remove comment", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})

	lg.Debug().Int("comment_id", id).Msgf("%s comment deleted successfully", pkgConst.OpSuccess)
}

func (hd *MainHandlers) FindComments(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "FindComments").Logger()

	query := c.Query("q")
	if query == "" {
		lg.Warn().Int("status", http.StatusBadRequest).Msgf("%s search query is required", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, response{Error: "search query is required"})
		return
	}

	comments, err := hd.sv.FindComments(c.Request.Context(), query)
	if err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to find comments", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	responses := convertToResponseList(comments)
	c.JSON(http.StatusOK, responses)

	lg.Debug().Msgf("%s comment finded successfully", pkgConst.OpSuccess)
}

func convertToResponse(comment *model.Comment) *response {
	resp := &response{
		ID:        comment.ID,
		ParentID:  comment.ParentID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}

	if comment.UpdatedAt != nil {
		t := comment.UpdatedAt.Format(time.RFC3339)
		resp.UpdatedAt = &t
	}

	if comment.Children != nil {
		resp.Children = convertToResponseList(comment.Children)
	}

	return resp
}

func convertToResponseList(comments []*model.Comment) []*response {
	responses := make([]*response, len(comments))
	for i, comment := range comments {
		responses[i] = convertToResponse(comment)
	}
	return responses
}
