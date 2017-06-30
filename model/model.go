package model

import (
	"github.com/jinzhu/gorm"
	"time"
	"crypto/sha256"
	"fmt"
)

type User struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Nickname  string
	Email     string `gorm:"type:varchar(100);not null;unique"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Website struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Name      string `gorm:type:varchar(100);not null;unique"`
	Host      string `gorm:"not null"`
	Port      string `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ListWebsites(db *gorm.DB) (websites []Website) {
	db.Find(&websites)
	return
}

func ListUserWebsites(db *gorm.DB, user_id uint) (websites []Website) {
	db.Where("user_id = ?", user_id).Find(&websites)
	return
}

func AddWebsite(db *gorm.DB, name string, host string, port string, user_id uint) (website Website)  {
	db.LogMode(true)
	website = Website{
		Name: name,
		Host: host,
		Port: port,
		UserID: user_id,
	}
	db.Debug().Create(&website)
	return
}

func DelWebsite(db *gorm.DB, id uint, user_id uint) (result bool, website Website) {
	db.Where("id = ?", id).First(&website)
	if website.UserID == user_id {
		db.Delete(&website)
		result = true
	} else {
		result = false
	}
	return
}

func AddUser(db *gorm.DB, nickname string, email string, password string) (user User) {
	db.LogMode(true)
	enc := sha256.New()
	enc.Write([]byte(password))
	dest := fmt.Sprintf("%x", enc.Sum(nil))
	user = User{
		Nickname: nickname,
		Email: email,
		Password: dest,
	}
	db.Debug().Create(&user)
	return
}

func CheckUser(db *gorm.DB, email string, password string) (result bool, user User) {
	db.Where("email = ?", email).First(&user)
	enc := sha256.New()
	enc.Write([]byte(password))
	dest := fmt.Sprintf("%x", enc.Sum(nil))
	if &user != nil && user.Password == dest {
		result = true
	} else {
		result = false
	}
	return
}
