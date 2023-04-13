package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var RoleCollection string = "roles"

type RoleName string

const (
	RoleNameUser       RoleName = "user"
	RoleNameSpecialist RoleName = "specialist"
	RoleNameMinion     RoleName = "minion"
)

type Role struct {
	ID string

	IsDefault bool
	Name      RoleName

	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoleDBSchema struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	IsDefault bool     `bson:"isDefault"`
	Name      RoleName `bson:"name"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}