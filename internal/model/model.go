package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Author struct {
	CreatedBy uint `json:"-"`
	UpdatedBy uint `json:"-"`
	DeletedBy uint `json:"-"`
	Creator   User `json:"creator" gorm:"foreignkey:CreatedBy"`
	Updator   User `json:"updator" gorm:"foreignkey:UpdatedBy"`
	Deletor   User `json:"deletor" gorm:"foreignkey:DeletedBy"`
}
