package model

import (
	_ "database/sql"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Sensor struct {
	Oid      uuid.UUID `gorm:"primary_key"`
	Eid      string    `gorm:"type:varchar(100);unique_index"`
	Unit     string    `gorm:"type:varchar(20)"`
	Readings []Reading `gorm:"foreignkey:SensorOid"`
}

type Reading struct {
	gorm.Model
	Value     int
	Timestamp string
	SensorOid uuid.UUID `gorm:"index"`
}
