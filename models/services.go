package models

import "github.com/jinzhu/gorm"

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true) // TODO - remove when env == production
	return &Services{}
}

type Services struct {
	Gallery GalleryService
	User    UserService
}
