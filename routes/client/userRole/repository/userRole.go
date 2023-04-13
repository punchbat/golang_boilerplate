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
		Collection: db.Collection(models.UserRoleCollection),
	}
}

func (r *Repository) CreateUserRole(ctx context.Context, userRole *models.UserRole) error {
	userRole.CreatedAt = time.Now()
	userRole.UpdatedAt = time.Now()

	model := mapToMongoSchema(userRole)
	res, err := r.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		userRole.ID = oid.Hex()
	}

	return nil
}

func (r *Repository) GetUserRoleByID(ctx context.Context, id string) (*models.UserRole, error) {
	userRole := new(models.UserRoleDBSchema)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": oid,
	}
	err = r.FindOne(ctx, filter).Decode(userRole)
	if err != nil {
		return nil, err
	}

	return mapToDomainModel(userRole), nil
}

func (r *Repository) GetUserRoleByIDs(ctx context.Context, ids []string) ([]*models.UserRole, error) {
	var userRoles []*models.UserRole
	var objectIDs []primitive.ObjectID

	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, oid)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	cur, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var userRoleDBSchema models.UserRoleDBSchema
		if err := cur.Decode(&userRoleDBSchema); err != nil {
			return userRoles, err
		}

		userRole := mapToDomainModel(&userRoleDBSchema)
		userRoles = append(userRoles, userRole)
	}

	if err := cur.Err(); err != nil {
		return userRoles, err
	}

	return userRoles, nil
}

func (r *Repository) DeleteUserRoleByID(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": oid,
	}

	_, err = r.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func mapToMongoSchema(i *models.UserRole) *models.UserRoleDBSchema {
	UserOid, err := primitive.ObjectIDFromHex(i.UserID)
	if err != nil {
		return nil
	}

	RoleOid, err := primitive.ObjectIDFromHex(i.RoleID)
	if err != nil {
		return nil
	}

	return &models.UserRoleDBSchema{
		UserID: UserOid,
		RoleID: RoleOid,
		Status: i.Status,

		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func mapToDomainModel(i *models.UserRoleDBSchema) *models.UserRole {
	return &models.UserRole{
		ID: i.ID.Hex(),

		UserID: i.UserID.Hex(),
		RoleID: i.RoleID.Hex(),
		Status: i.Status,

		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}
