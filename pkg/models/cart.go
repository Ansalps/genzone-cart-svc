package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID string `validate:"required"`
	//User   User `gorm:"foriegnkey:UserID;references:ID"`
	// CartID      string  `validate:"required,numeric"`
	// Cart        Cart    `gorm:"foriegnkey:CartID;references:ID"`
	ProductID string `validate:"required,numeric"`
	//Product     Product `gorm:"foriegnkey:ProductID;references:ID"`
	Qty    uint    `gorm:"default:0"`
	Price  float64 `gorm:"type:decimal(10,2)" `
	Amount float64 `gorm:"type:decimal(10,2);default:0.00"  `

	//Discount    float64 `gorm:"default:0.00"`
	//FinalAmount float64
}
