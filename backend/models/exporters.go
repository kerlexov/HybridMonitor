package models

import "gorm.io/gorm"

type RedfishHost struct {
	gorm.Model
	Host     string
	Username string
	Password string
}

type VsphereHost struct {
	gorm.Model
	Host     string
	Username string
	Password string
}
