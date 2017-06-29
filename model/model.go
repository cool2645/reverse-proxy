package model

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Nickname  string
	Email     string `gorm:"type:varchar(100);not null;unique"`
	CreateAt  time.Time
	UpdatedAt time.Time
}

type Website struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Name      string `gorm:type:varchar(100);not null;unique"`
	Host      string `gorm:"not null"`
	Port      string `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	CreateAt  time.Time
	UpdatedAt time.Time
}

func ListWebsites(db *gorm.DB) (websites []Website, err error) {
	err = db.Table("websites").Scan(&websites).Error
	if err != nil {
		err = errors.Wrap(err, "ListUser")
		return
	}
	return
}
