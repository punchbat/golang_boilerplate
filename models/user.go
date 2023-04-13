package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var UserCollection string = "users"

type Gender int

const (
	Male Gender = iota + 1
	Female
	NonBinary
)

func (gender Gender) String() string {
	terms := []string{"Male", "Female", "NonBinary"}
	if gender < Male || gender > NonBinary {
		return "unknown"
	}
	return terms[gender]
}

type Address struct {
	Country     string
	City        string
	Street      string
	HouseNumber string
}

type AddressDBSchema struct {
	Country     string `bson:"country"`
	City        string `bson:"city"`
	Street      string `bson:"street"`
	HouseNumber string `bson:"houseNumber"`
}

type User struct {
	ID string

	Email           string
	Password        string
	PasswordConfirm string
	Verified        bool
	VerifyCode      string

	UserRoleIDs []string
	UserRoles   []*UserRoleWithRole

	FinishedRegistration bool

	IIN      int
	Name     string
	Surname  string
	Birthday time.Time
	Gender   Gender
	Address  Address

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserDBSchema struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Email           string `bson:"email"`
	Password        string `bson:"password"`
	PasswordConfirm string `bson:"passwordConfirm"`
	Verified        bool   `bson:"verified"`
	VerifyCode      string `bson:"verifyCode"`

	UserRoleIDs          []primitive.ObjectID `bson:"userRoleIds"`
	FinishedRegistration bool                 `bson:"finishedRegistration"`

	IIN      int             `bson:"IIN"`
	Name     string          `bson:"name"`
	Surname  string          `bson:"surname"`
	Birthday time.Time       `bson:"birthday"`
	Gender   Gender          `bson:"gender"`
	Address  AddressDBSchema `bson:"address"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}