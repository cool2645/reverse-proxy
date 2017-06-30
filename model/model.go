package model

import (
	"github.com/jinzhu/gorm"
	"time"
	"crypto/sha256"
	"fmt"
	"strings"
)

type User struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Nickname  string
	Email     string `gorm:"type:varchar(100);not null;unique"`
	Password  string `gorm:"not null"`
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Website struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Name      string `gorm:type:varchar(100);not null;unique"`
	Host      string `gorm:"not null"`
	Port      string `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	Domains   []Domain
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Domain struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Name      string `gorm:type:varchar(100);not null"`
	WebsiteID uint    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindWebsiteByDomain(db *gorm.DB, host string) (result bool, website Website) {
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}
	domain_name := host
	var domain Domain
	var count uint
	db.Model(&Domain{}).Where("name = ?", domain_name).Count(&count)
	if count > 0 {
		db.Where("name = ?", domain_name).First(&domain)
		db.Model(&Website{}).Where("id = ?", domain.WebsiteID).Count(&count)
		if count > 0 {
			db.Where("id = ?", domain.WebsiteID).First(&website)
			result = true
			return
		}
	}
	result = false
	return
}

func AddDomain(db *gorm.DB, website_id uint, user_id uint, domain_name string) (result bool, domain Domain) {
	var website Website
	db.Where("id = ?", website_id).First(&website)
	if website.UserID == user_id {
		domain = Domain{
			Name:      domain_name,
			WebsiteID: website.ID,
		}
		db.Create(&domain)
		result = true
		return
	}
	result = false
	return
}

func DelDomain(db *gorm.DB, domain_id uint, user_id uint) (result bool, domain Domain) {
	db.Where("id = ?", domain_id).First(&domain)
	var website Website
	db.Where("id = ?", domain.WebsiteID).First(&website)
	if website.UserID == user_id {
		db.Delete(&domain)
		result = true
		return
	}
	result = false
	return
}

func ListWebsites(db *gorm.DB) (websites []Website) {
	db.Preload("Domains").Find(&websites)
	return
}

func ListUserWebsites(db *gorm.DB, user_id uint) (websites []Website) {
	var user User
	db.Where("id = ?", user_id).First(&user)
	if user.Role == "admin" {
		db.Preload("Domains").Find(&websites)
	} else {
		db.Preload("Domains").Where("user_id = ?", user_id).Find(&websites)
	}
	return
}

func AddWebsite(db *gorm.DB, name string, host string, port string, user_id uint) (result bool, website Website) {
	var count uint
	db.Model(&Website{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		result = false
	} else {
		website = Website{
			Name:   name,
			Host:   host,
			Port:   port,
			UserID: user_id,
		}
		db.Create(&website)
		result = true
	}
	return
}

func DelWebsite(db *gorm.DB, id uint, user_id uint) (result bool, website Website) {
	var user User
	db.Where("id = ?", user_id).First(&user)
	db.Where("id = ?", id).First(&website)
	if website.UserID == user_id || user.Role == "admin" {
		db.Delete(&website)
		db.Where("website_id = ?", website.ID).Delete(Domain{})
		result = true
	} else {
		result = false
	}
	return
}

func AddUser(db *gorm.DB, nickname string, email string, password string) (result bool, user User) {
	var count uint
	db.Model(&User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		result = false
	} else {
		enc := sha256.New()
		enc.Write([]byte(password))
		dest := fmt.Sprintf("%x", enc.Sum(nil))
		user = User{
			Nickname: nickname,
			Email:    email,
			Password: dest,
		}
		db.Create(&user)
		result = true
	}
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
