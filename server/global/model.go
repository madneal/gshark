package global

import (
	"time"
)

type GVA_MODEL struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
