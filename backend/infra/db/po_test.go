package db

import "gorm.io/gorm"

type somePO struct {
	ID          int `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true;comment:ID" json:"id"`
	Name        string
	Description string
	DeletedAt   gorm.DeletedAt
}
