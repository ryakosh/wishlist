package graph

import "github.com/jinzhu/gorm"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	DB *gorm.DB
}
