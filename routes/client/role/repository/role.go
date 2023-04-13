package repository

import (
	"context"
	"health/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	*mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Collection: db.Collection(models.RoleCollection),
	}
}

func (r *Repository) GetRoleByID(ctx context.Context, id string) (*models.Role, error) {
	role := new(models.RoleDBSchema)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": oid,
	}
	err = r.FindOne(ctx, filter).Decode(role)
	if err != nil {
		return nil, err
	}

	return mapToDomainModel(role), nil
}

func (r *Repository) GetRoleByIDs(ctx context.Context, ids []string) ([]*models.Role, error) {
	if len(ids) == 0 {
		return []*models.Role{}, nil
	}

	var roles []*models.Role
	var objectIds []primitive.ObjectID

	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIds = append(objectIds, oid)
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": objectIds,
		},
	}

	cur, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		role := new(models.RoleDBSchema)
		err := cur.Decode(role)
		if err != nil {
			return nil, err
		}
		roles = append(roles, mapToDomainModel(role))
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *Repository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	role := new(models.RoleDBSchema)

	filter := bson.M{
		"name": name,
	}

	if err := r.FindOne(ctx, filter).Decode(role); err != nil {
		return nil, err
	}

	return mapToDomainModel(role), nil
}

func (r Repository) GetRoles(ctx context.Context) ([]*models.Role, error) {
	var roles []*models.Role
	filter := bson.M{}
	cursor, err := r.Find(ctx, filter)
	if err != nil {
		return roles, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var roleDBSchema models.RoleDBSchema
		err = cursor.Decode(&roleDBSchema)
		if err != nil {
			return roles, err
		}

		role := mapToDomainModel(&roleDBSchema)
		roles = append(roles, role)
	}

	if err = cursor.Err(); err != nil {
		return roles, err
	}

	return roles, nil
}

func mapToMongoSchema(i *models.Role) *models.RoleDBSchema {
	return &models.RoleDBSchema{
		IsDefault: i.IsDefault,
		Name:      i.Name,

		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func mapToDomainModel(i *models.RoleDBSchema) *models.Role {
	return &models.Role{
		ID: i.ID.Hex(),

		IsDefault: i.IsDefault,
		Name:      i.Name,

		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}
