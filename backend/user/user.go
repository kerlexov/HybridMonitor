package user

import (
	"github.com/kerlexov/HybridMonitor/models"
	"gorm.io/gorm"
)

type LoginUser struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type Session struct {
	Token string
}

func IsUserValid(db *gorm.DB, u *LoginUser) bool {
	var user models.User
	result := db.Where(&models.User{
		Username: u.Username,
		Password: u.Password,
	}).First(&user)

	if result.Error != nil {
		return false
	}

	return true
}
