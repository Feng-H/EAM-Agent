package service

import (
	"errors"
	"fmt"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/pkg/database"
	"gorm.io/gorm"
)

// SparePartService
type SparePartService struct {
	partRepo        *repository.SparePartRepository
	inventoryRepo   *repository.SparePartInventoryRepository
	consumptionRepo *repository.SparePartConsumptionRepository
	txRepo          *repository.SparePartTransactionRepository
	userRepo        *repository.UserRepository
}

func NewSparePartService() *SparePartService {
	return &SparePartService{
		partRepo:        repository.NewSparePartRepository(),
		inventoryRepo:   repository.NewSparePartInventoryRepository(),
		consumptionRepo: repository.NewSparePartConsumptionRepository(),
		txRepo:          repository.NewSparePartTransactionRepository(),
		userRepo:        repository.NewUserRepository(),
	}
}

func (s *SparePartService) CreatePart(code, name, specification, unit string, factoryID *uint, safetyStock int) (*model.SparePart, error) {
	// Check if code exists
	if _, err := s.partRepo.GetByCode(code); err == nil {
		return nil, errors.New("code already exists")
	}

	part := &model.SparePart{
		Code:          code,
		Name:          name,
		Specification: specification,
		Unit:          unit,
		FactoryID:     factoryID,
		SafetyStock:   safetyStock,
	}

	if err := s.partRepo.Create(part); err != nil {
		return nil, err
	}

	return s.partRepo.GetByID(part.ID)
}

func (s *SparePartService) GetPartByID(id uint) (*model.SparePart, error) {
	return s.partRepo.GetByID(id)
}

func (s *SparePartService) ListParts(filter repository.SparePartFilter) (*SparePartListResult, error) {
	parts, total, err := s.partRepo.List(filter)
	if err != nil {
		return nil, err
	}

	return &SparePartListResult{
		Items: parts,
		Total: total,
	}, nil
}

func (s *SparePartService) UpdatePart(id uint, code, name, specification, unit string, factoryID *uint, safetyStock int) error {
	part, err := s.partRepo.GetByID(id)
	if err != nil {
		return ErrNotFound
	}

	// Check if new code conflicts
	if code != part.Code {
		if existing, _ := s.partRepo.GetByCode(code); existing != nil {
			return errors.New("code already exists")
		}
	}

	part.Code = code
	part.Name = name
	part.Specification = specification
	part.Unit = unit
	part.FactoryID = factoryID
	part.SafetyStock = safetyStock

	return s.partRepo.Update(part)
}

func (s *SparePartService) DeletePart(id uint) error {
	return s.partRepo.Delete(id)
}

// Inventory operations
func (s *SparePartService) StockIn(partID, factoryID uint, quantity int, remark string, userID uint) error {
	// Verify part exists
	_, err := s.partRepo.GetByID(partID)
	if err != nil {
		return ErrNotFound
	}

	// Update inventory
	if err := s.inventoryRepo.UpdateQuantity(partID, factoryID, quantity); err != nil {
		return err
	}

	// Create transaction record
	tx := &model.SparePartTransaction{
		SparePartID: partID,
		FactoryID:   factoryID,
		Type:        "in",
		Quantity:    quantity,
		OperatorID:  userID,
		Remark:      remark,
	}

	return s.txRepo.Create(tx)
}

func (s *SparePartService) StockOut(partID, factoryID uint, quantity int, orderID, taskID *uint, remark string, userID uint) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Verify part exists
	_, err := s.partRepo.GetByID(partID)
	if err != nil {
		return ErrNotFound
	}

	db := database.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Atomic decrement with stock floor condition (no overselling)
		res := tx.Model(&model.SparePartInventory{}).
			Where("spare_part_id = ? AND factory_id = ? AND quantity >= ?", partID, factoryID, quantity).
			Update("quantity", gorm.Expr("quantity - ?", quantity))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("insufficient stock: requested %d", quantity)
		}

		// 2. Create consumption record within the same transaction
		consumption := &model.SparePartConsumption{
			SparePartID: partID,
			OrderID:     orderID,
			TaskID:      taskID,
			Quantity:    quantity,
			UserID:      userID,
		}
		if err := tx.Create(consumption).Error; err != nil {
			return err
		}

		// 3. Create transaction record within the same transaction
		txRecord := &model.SparePartTransaction{
			SparePartID: partID,
			FactoryID:   factoryID,
			Type:        "out",
			Quantity:    -quantity,
			OperatorID:  userID,
			RelatedID:   orderID,
			Remark:      remark,
		}
		if txRecord.RelatedID == nil {
			txRecord.RelatedID = taskID
		}

		return tx.Create(txRecord).Error
	})
}

func (s *SparePartService) GetInventory(filter repository.InventoryFilter) (*InventoryListResult, error) {
	inventories, total, err := s.inventoryRepo.List(filter)
	if err != nil {
		return nil, err
	}

	return &InventoryListResult{
		Items: inventories,
		Total: total,
	}, nil
}

func (s *SparePartService) GetLowStockAlerts() ([]repository.LowStockAlert, error) {
	return s.inventoryRepo.GetLowStockAlerts()
}

// Consumption operations
func (s *SparePartService) RecordConsumption(partID uint, quantity int, orderID, taskID *uint, userID uint, remark string) error {
	// Verify part exists
	_, err := s.partRepo.GetByID(partID)
	if err != nil {
		return ErrNotFound
	}

	consumption := &model.SparePartConsumption{
		SparePartID: partID,
		OrderID:     orderID,
		TaskID:      taskID,
		Quantity:    quantity,
		UserID:      userID,
	}

	if err := s.consumptionRepo.Create(consumption); err != nil {
		return err
	}

	_ = remark
	return nil
}

func (s *SparePartService) GetConsumption(filter repository.ConsumptionFilter) (*ConsumptionListResult, error) {
	consumptions, total, err := s.consumptionRepo.List(filter)
	if err != nil {
		return nil, err
	}

	return &ConsumptionListResult{
		Items: consumptions,
		Total: total,
	}, nil
}

func (s *SparePartService) GetStatistics() (*SparePartStatistics, error) {
	stats, err := s.consumptionRepo.GetStatistics()
	if err != nil {
		return nil, err
	}

	return &SparePartStatistics{
		TotalParts:        stats["total_parts"],
		LowStockCount:     stats["low_stock_count"],
		TotalStockValue:   0, // Calculate if needed
		MonthlyConsumption: stats["monthly_consumption"],
	}, nil
}

// Types
type SparePartListResult struct {
	Items []model.SparePart
	Total int64
}

type InventoryListResult struct {
	Items []model.SparePartInventory
	Total int64
}

type ConsumptionListResult struct {
	Items []model.SparePartConsumption
	Total int64
}

type SparePartStatistics struct {
	TotalParts        int64
	LowStockCount     int64
	TotalStockValue   float64
	MonthlyConsumption int64
}
