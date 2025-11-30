package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/model"
)

// IItemRp interface for item repository
type IItemRp interface {
	Create(item *model.Item, userID int, userName string) (*model.Item, error)
	GetAll() ([]model.Item, error)
	GetByID(id int) (*model.Item, error)
	Update(item *model.Item, userID int, userName string) error
	Delete(id int, userID int, userName string) error
}

// IItemHistoryRp interface for item history repository
type IItemHistoryRp interface {
	GetByItemID(itemID int) ([]model.ItemAction, error)
	GetAll() ([]model.ItemAction, error)
	ExportToCSV(itemID int) ([]map[string]interface{}, error)
	CreateAction(itemID int, actionType string, userID int, userName string, changes map[string]interface{}) error
}

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

// ItemService service for working with items
type ItemService struct {
	itemRp    IItemRp
	historyRp IItemHistoryRp
}

// NewItemService creates a new ItemService
func NewItemService(itemRp IItemRp, historyRp IItemHistoryRp) *ItemService {
	return &ItemService{
		itemRp:    itemRp,
		historyRp: historyRp,
	}
}

// Create creates a new item (only for Кладовщик)
func (sv *ItemService) Create(ctx context.Context, item *model.Item, userID int, userRole, userName string) (*model.Item, error) {
	// Проверяем права доступа
	if userRole != "Кладовщик" {
		return nil, fmt.Errorf("access denied: only Кладовщик can create items")
	}

	// Валидация данных
	if item.Name == "" {
		return nil, fmt.Errorf("item name cannot be empty")
	}
	if item.Price < 0 {
		return nil, fmt.Errorf("item price cannot be negative")
	}
	if item.Quantity < 0 {
		return nil, fmt.Errorf("item quantity cannot be negative")
	}

	// Устанавливаем временные метки
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// Создаем товар (триггер автоматически создаст запись в истории)
	return sv.itemRp.Create(item, userID, userName)
}

// GetAll returns all items (for all authenticated users)
func (sv *ItemService) GetAll(ctx context.Context, userID int) ([]model.Item, error) {
	return sv.itemRp.GetAll()
}

// GetByID returns an item by ID (for all authenticated users)
func (sv *ItemService) GetByID(ctx context.Context, id int, userID int) (*model.Item, error) {
	return sv.itemRp.GetByID(id)
}

// Update updates an item (only for Кладовщик)
func (sv *ItemService) Update(ctx context.Context, item *model.Item, userID int, userRole, userName string) error {
	// Проверяем права доступа
	if userRole != "Кладовщик" {
		return fmt.Errorf("access denied: only Кладовщик can update items")
	}

	// Валидация данных
	if item.Name == "" {
		return fmt.Errorf("item name cannot be empty")
	}
	if item.Price < 0 {
		return fmt.Errorf("item price cannot be negative")
	}
	if item.Quantity < 0 {
		return fmt.Errorf("item quantity cannot be negative")
	}

	// Обновляем время изменения
	item.UpdatedAt = time.Now()

	// Обновляем товар (триггер автоматически создаст запись в истории)
	return sv.itemRp.Update(item, userID, userName)
}

// Delete deletes an item (only for Кладовщик)
func (sv *ItemService) Delete(ctx context.Context, id int, userID int, userRole, userName string) error {
	// Проверяем права доступа
	if userRole != "Кладовщик" {
		return fmt.Errorf("access denied: only Кладовщик can delete items")
	}

	// Удаляем товар (триггер автоматически создаст запись в истории)
	return sv.itemRp.Delete(id, userID, userName)
}

// GetHistory returns item history (only for Аудитор)
func (sv *ItemService) GetHistory(ctx context.Context, itemID int, userID int, userRole string) ([]model.ItemAction, error) {
	// Проверяем права доступа
	if userRole != "Аудитор" {
		return nil, fmt.Errorf("access denied: only Аудитор can view item history")
	}

	return sv.historyRp.GetByItemID(itemID)
}

// GetAllHistory returns all item history (only for Аудитор)
func (sv *ItemService) GetAllHistory(ctx context.Context, userID int, userRole string) ([]model.ItemAction, error) {
	// Проверяем права доступа
	if userRole != "Аудитор" {
		return nil, fmt.Errorf("access denied: only Аудитор can view item history")
	}

	return sv.historyRp.GetAll()
}

// ExportHistoryToCSV exports item history to CSV format (only for Аудитор)
func (sv *ItemService) ExportHistoryToCSV(ctx context.Context, itemID int, userID int, userRole string) ([]map[string]interface{}, error) {
	// Проверяем права доступа
	if userRole != "Аудитор" {
		return nil, fmt.Errorf("access denied: only Аудитор can export item history")
	}

	return sv.historyRp.ExportToCSV(itemID)
}
