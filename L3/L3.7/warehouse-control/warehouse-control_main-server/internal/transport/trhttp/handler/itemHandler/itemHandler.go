package itemHandler

import (
	"context"
	"encoding/csv"
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
	// CRUD операции с товарами
	rt.POST("/items", hd.createItem)
	rt.GET("/items", hd.getAllItems)
	rt.GET("/items/:id", hd.getItemByID)
	rt.PUT("/items/:id", hd.updateItem)
	rt.DELETE("/items/:id", hd.deleteItem)

	// История изменений
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

	// Получаем информацию о пользователе из контекста
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

	// Создаем контекст с информацией о пользователе
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
	// Получаем информацию о пользователе из Gin контекста
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

	// Получаем информацию о пользователе из Gin контекста
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

	// Получаем ID из URL
	idParam := c.Param("id")
	itemID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})
		return
	}
	item.ID = itemID

	// Получаем информацию о пользователе
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

	// Создаем контекст с информацией о пользователе
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

	// Получаем информацию о пользователе
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

	// Создаем контекст с информацией о пользователе
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

	// Получаем роль пользователя
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	// Получаем ID пользователя
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
	// Получаем роль пользователя
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	// Получаем ID пользователя
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

	// Получаем роль пользователя
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	// Получаем ID пользователя
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

	// Создаем CSV файл
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=item_history.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Записываем заголовки
	headers := []string{"ID", "Товар", "Действие", "Пользователь", "Дата", "Изменения"}
	writer.Write(headers)

	// Записываем данные
	for _, row := range csvData {
		record := []string{
			row["ID"].(string),
			row["Товар"].(string),
			row["Действие"].(string),
			row["Пользователь"].(string),
			row["Дата"].(string),
			row["Изменения"].(string),
		}
		writer.Write(record)
	}
}
