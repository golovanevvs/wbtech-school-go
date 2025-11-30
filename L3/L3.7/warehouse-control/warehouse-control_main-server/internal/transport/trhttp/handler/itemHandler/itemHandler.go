package itemHandler

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/model"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// IItemService interface for item service
type IItemService interface {
	Create(ctx context.Context, item *model.Item, userID int, userRole, userName string) (*model.Item, error)
	GetAll(ctx context.Context, userID int) ([]model.Item, error)
	GetByID(ctx context.Context, id int, userID int) (*model.Item, error)
	Update(ctx context.Context, item *model.Item, userID int, userRole, userName string) error
	Delete(ctx context.Context, id int, userID int, userRole, userName string) error
	GetHistory(ctx context.Context, itemID int, userID int, userRole string) ([]model.ItemAction, error)
	GetAllHistory(ctx context.Context, userID int, userRole string) ([]model.ItemAction, error)
	ExportHistoryToCSV(ctx context.Context, itemID int, userID int, userRole string) ([]map[string]interface{}, error)
}

// ItemHandler handles item-related HTTP requests
type ItemHandler struct {
	lg *zlog.Zerolog
	sv IItemService
}

// New creates a new ItemHandler
func New(lg *zlog.Zerolog, sv IItemService) *ItemHandler {
	return &ItemHandler{
		lg: lg,
		sv: sv,
	}
}

// RegisterProtectedRoutes registers protected routes for items
func (hd *ItemHandler) RegisterProtectedRoutes(rt *ginext.RouterGroup) {
	rt.POST("/items", hd.createItem)
	rt.GET("/items", hd.getAllItems)
	rt.GET("/items/:id", hd.getItemByID)
	rt.PUT("/items/:id", hd.updateItem)
	rt.DELETE("/items/:id", hd.deleteItem)

	rt.GET("/items/:id/history", hd.getItemHistory)
	rt.GET("/items/:id/history/export", hd.exportItemHistory)
	rt.GET("/history", hd.getAllHistory)
}

// createItem creates a new item
func (hd *ItemHandler) createItem(c *gin.Context) {
	var item model.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	userName, exists := c.Get("user_name")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user name not found",
		})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "user_id", userID)

	createdItem, err := hd.sv.Create(ctx, &item, userID.(int), userRole.(string), userName.(string))
	if err != nil {
		hd.lg.Err(err).Msg("failed to create item")
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"item": createdItem,
	})
}

// getAllItems returns all items
func (hd *ItemHandler) getAllItems(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user ID type",
		})
		return
	}

	ctx := c.Request.Context()

	items, err := hd.sv.GetAll(ctx, userIDInt)
	if err != nil {
		hd.lg.Err(err).Msg("failed to get all items")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve items",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}

// getItemByID returns an item by ID
func (hd *ItemHandler) getItemByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user ID type",
		})
		return
	}

	ctx := c.Request.Context()

	item, err := hd.sv.GetByID(ctx, id, userIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "item not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item": item,
	})
}

// updateItem updates an item
func (hd *ItemHandler) updateItem(c *gin.Context) {
	var item model.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	idParam := c.Param("id")
	itemID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})
		return
	}
	item.ID = itemID

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	userName, exists := c.Get("user_name")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user name not found",
		})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "user_id", userID)

	err = hd.sv.Update(ctx, &item, userID.(int), userRole.(string), userName.(string))
	if err != nil {
		hd.lg.Err(err).Msg("failed to update item")
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "item updated successfully",
	})
}

// deleteItem deletes an item
func (hd *ItemHandler) deleteItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	userName, exists := c.Get("user_name")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user name not found",
		})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "user_id", userID)

	err = hd.sv.Delete(ctx, id, userID.(int), userRole.(string), userName.(string))
	if err != nil {
		hd.lg.Err(err).Msg("failed to delete item")
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "item deleted successfully",
	})
}

// getItemHistory returns the history of changes for a specific item
func (hd *ItemHandler) getItemHistory(c *gin.Context) {
	idParam := c.Param("id")
	itemID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	ctx := c.Request.Context()

	history, err := hd.sv.GetHistory(ctx, itemID, userID.(int), userRole.(string))
	if err != nil {
		hd.lg.Err(err).Msg("failed to get item history")
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}

// getAllHistory returns all item history
func (hd *ItemHandler) getAllHistory(c *gin.Context) {
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	ctx := c.Request.Context()

	history, err := hd.sv.GetAllHistory(ctx, userID.(int), userRole.(string))
	if err != nil {
		hd.lg.Err(err).Msg("failed to get all history")
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}

// exportItemHistory exports item history to CSV
func (hd *ItemHandler) exportItemHistory(c *gin.Context) {
	idParam := c.Param("id")
	itemID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	ctx := c.Request.Context()

	csvData, err := hd.sv.ExportHistoryToCSV(ctx, itemID, userID.(int), userRole.(string))
	if err != nil {
		hd.lg.Err(err).Msg("failed to export item history")
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=item_history.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	headers := []string{"ID", "Товар", "Действие", "Пользователь", "Дата", "Изменения"}
	writer.Write(headers)

	for _, row := range csvData {
		record := []string{
			fmt.Sprintf("%v", row["ID"]),
			fmt.Sprintf("%v", row["Товар"]),
			fmt.Sprintf("%v", row["Действие"]),
			fmt.Sprintf("%v", row["Пользователь"]),
			fmt.Sprintf("%v", row["Дата"]),
			fmt.Sprintf("%v", row["Изменения"]),
		}
		writer.Write(record)
	}
}
