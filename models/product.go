package models

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID              string         `json:"_id" gorm:"primaryKey;type:varchar(255)"`
	Name            string         `json:"name" gorm:"type:text;index"`
	Brand           string         `json:"brand" gorm:"type:varchar(255);index"`
	Price           float64        `json:"price" gorm:"type:decimal(10,2)"`
	Currency        string         `json:"currency" gorm:"type:varchar(10)"`
	PriceInRubles   *float64       `json:"priceInRubles" gorm:"type:decimal(10,2)"`
	DiscountedPrice *float64       `json:"discountedPrice" gorm:"type:decimal(10,2)"`
	Description     string         `json:"description" gorm:"type:text"`
	Images          datatypes.JSON `json:"images" gorm:"type:jsonb"`
	Sizes           datatypes.JSON `json:"sizes" gorm:"type:jsonb"`
	Colors          datatypes.JSON `json:"colors" gorm:"type:jsonb"`
	ProductURL      string         `json:"productUrl" gorm:"type:text"`
	Store           string         `json:"store" gorm:"type:varchar(255);index"`
	Category        string         `json:"category" gorm:"type:varchar(255);index"`
	ProcessedAt     string         `json:"processedAt" gorm:"type:varchar(50)"`
	IsActive        bool           `json:"isActive" gorm:"default:true;index"`
	StockStatus     string         `json:"stockStatus" gorm:"type:varchar(50);index"`
	Stock           datatypes.JSON `json:"stock" gorm:"type:jsonb"`
	CreatedAt       time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
}

type Size struct {
	SizeName string `json:"sizeName"`
	OnStock  bool   `json:"onStock"`
}

type Color struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

type Stock struct {
	Quantity  int  `json:"quantity"`
	IsInStock bool `json:"isInStock"`
}

// TableName specifies the table name for GORM
func (Product) TableName() string {
	return "products"
}