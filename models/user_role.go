package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var UserRoleCollection string = "user.roles"

type UserRoleStatus string

const (
	UserRoleStatusCanceled UserRoleStatus = "canceled"
	UserRoleStatusPending  UserRoleStatus = "pending"
	UserRoleStatusApproved UserRoleStatus = "approved"
)

type UserRole struct {
	ID string

	UserID string
	RoleID string
	Status UserRoleStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRoleDBSchema struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	UserID primitive.ObjectID `bson:"userId"`
	RoleID primitive.ObjectID `bson:"roleId"`
	Status UserRoleStatus     `bson:"status"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type UserRoleWithRole struct {
	IsDefault bool
	Name      RoleName
	Status    UserRoleStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}