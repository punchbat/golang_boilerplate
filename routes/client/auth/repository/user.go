package repository

import (
	"context"
	"health/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	*mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Collection: db.Collection(models.UserCollection),
	}
}

func (r Repository) CreateUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	model := mapToMongoSchema(user)
	res, err := r.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}

	return nil
}

func (r Repository) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	model := mapToMongoSchema(user)

	filter := bson.M{
		"email": model.Email,
	}
	update := bson.M{
		"$set": model,
	}
	_, err := r.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.UserDBSchema)

	filter := bson.M{
		"email": email,
	}

	err := r.FindOne(ctx, filter).Decode(user)

	if err != nil {
		return nil, err
	}

	return mapToDomainModel(user), nil
}

func (r *Repository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user := new(models.UserDBSchema)
	err = r.FindOne(ctx, bson.M{"_id": oid}).Decode(user)

	if err != nil {
		return nil, err
	}

	return mapToDomainModel(user), nil
}

func (r *Repository) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}

func mapToMongoSchema(u *models.User) *models.UserDBSchema {
	rolesLikeID := make([]primitive.ObjectID, len(u.UserRoleIDs))
	for i, roleID := range u.UserRoleIDs {
		if id, err := primitive.ObjectIDFromHex(roleID); err == nil {
			rolesLikeID[i] = id
		}
	}

	return &models.UserDBSchema{
		Email:           u.Email,
		Password:        u.Password,
		PasswordConfirm: u.PasswordConfirm,
		Verified:        u.Verified,
		VerifyCode:      u.VerifyCode,

		FinishedRegistration: u.FinishedRegistration,

		IIN:      u.IIN,
		Name:     u.Name,
		Surname:  u.Surname,
		Birthday: u.Birthday,
		Gender:   u.Gender,
		Address:  models.AddressDBSchema(u.Address),

		UserRoleIDs: rolesLikeID,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func mapToDomainModel(u *models.UserDBSchema) *models.User {
	rolesLikeString := make([]string, len(u.UserRoleIDs))
	for i, roleID := range u.UserRoleIDs {
		rolesLikeString[i] = roleID.Hex()
	}

	return &models.User{
		ID: u.ID.Hex(),

		Email:           u.Email,
		Password:        u.Password,
		PasswordConfirm: u.PasswordConfirm,
		Verified:        u.Verified,
		VerifyCode:      u.VerifyCode,

		FinishedRegistration: u.FinishedRegistration,

		IIN:      u.IIN,
		Name:     u.Name,
		Surname:  u.Surname,
		Birthday: u.Birthday,
		Gender:   u.Gender,
		Address:  models.Address(u.Address),

		UserRoleIDs: rolesLikeString,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
