package mongo

import (
	"context"
	"time"

	"github.com/labasubagia/realworld-backend/internal/adapter/repository/mongo/model"
	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
	"go.mongodb.org/mongo-driver/bson"
)

type userRepo struct {
	db DB
}

func NewUserRepository(db DB) port.UserRepository {
	return &userRepo{
		db: db,
	}
}
func (r *userRepo) CreateUser(ctx context.Context, arg domain.User) (domain.User, error) {
	user := model.AsUser(arg)
	_, err := r.db.Collection(CollectionUser).InsertOne(ctx, user)
	if err != nil {
		return domain.User{}, intoException(err)
	}
	return user.ToDomain(), nil
}

func (r *userRepo) FilterFollow(ctx context.Context, arg port.FilterUserFollowPayload) ([]domain.UserFollow, error) {
	query := []bson.M{}
	if len(arg.FollowerIDs) > 0 {
		query = append(query, bson.M{"follower_id": bson.M{"$in": arg.FollowerIDs}})
	}
	if len(arg.FolloweeIDs) > 0 {
		query = append(query, bson.M{"followee_id": bson.M{"$in": arg.FolloweeIDs}})
	}
	filter := bson.M{}
	if len(query) > 0 {
		filter = bson.M{"$and": query}
	}

	cursor, err := r.db.Collection(CollectionUserFollow).Find(ctx, filter)
	if err != nil {
		return []domain.UserFollow{}, intoException(err)
	}

	result := []domain.UserFollow{}
	for cursor.Next(ctx) {
		data := model.UserFollow{}
		if err := cursor.Decode(&data); err != nil {
			return []domain.UserFollow{}, intoException(err)
		}
		result = append(result, data.ToDomain())
	}

	return result, nil
}

func (r *userRepo) FilterUser(ctx context.Context, arg port.FilterUserPayload) ([]domain.User, error) {

	query := []bson.M{}
	if len(arg.IDs) > 0 {
		query = append(query, bson.M{"id": bson.M{"$in": arg.IDs}})
	}
	if len(arg.Emails) > 0 {
		query = append(query, bson.M{"email": bson.M{"$in": arg.Emails}})
	}
	if len(arg.Usernames) > 0 {
		query = append(query, bson.M{"username": bson.M{"$in": arg.Usernames}})
	}

	filter := bson.M{}
	if len(query) > 0 {
		filter = bson.M{"$and": query}
	}

	cursor, err := r.db.Collection(CollectionUser).Find(ctx, filter)
	if err != nil {
		return []domain.User{}, intoException(err)
	}

	result := []domain.User{}
	for cursor.Next(ctx) {
		data := model.User{}
		if err := cursor.Decode(&data); err != nil {
			return []domain.User{}, intoException(err)
		}
		result = append(result, data.ToDomain())
	}

	return result, nil
}

func (r *userRepo) FindOne(ctx context.Context, arg port.FilterUserPayload) (domain.User, error) {
	users, err := r.FilterUser(ctx, arg)
	if err != nil {
		return domain.User{}, intoException(err)
	}
	if len(users) == 0 {
		return domain.User{}, exception.New(exception.TypeNotFound, "user not found", nil)
	}
	return users[0], nil
}

func (r *userRepo) Follow(ctx context.Context, arg domain.UserFollow) (domain.UserFollow, error) {
	follow := model.AsUserFollow(arg)
	_, err := r.db.Collection(CollectionUserFollow).InsertOne(ctx, follow)
	if err != nil {
		return domain.UserFollow{}, intoException(err)
	}
	return follow.ToDomain(), nil
}

func (r *userRepo) UnFollow(ctx context.Context, arg domain.UserFollow) (domain.UserFollow, error) {
	_, err := r.db.Collection(CollectionUserFollow).DeleteOne(ctx, bson.M{
		"follower_id": arg.FollowerID,
		"followee_id": arg.FolloweeID,
	})
	if err != nil {
		return domain.UserFollow{}, intoException(err)
	}
	return arg, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, arg domain.User) (domain.User, error) {

	filter := bson.M{"id": arg.ID}
	fields := bson.M{}
	if arg.Email != "" {
		fields["email"] = arg.Email
	}
	if arg.Username != "" {
		fields["username"] = arg.Username
	}
	if arg.Bio != "" {
		fields["bio"] = arg.Bio
	}
	if arg.Image != "" {
		fields["image"] = arg.Image
	}
	if arg.Password != "" {
		fields["password"] = arg.Password
	}
	if len(fields) > 0 {
		fields["updated_at"] = time.Now()
	}

	_, err := r.db.Collection(CollectionUser).UpdateOne(ctx, filter, bson.M{"$set": fields})
	if err != nil {
		return domain.User{}, intoException(err)
	}

	updated, err := r.FindOne(ctx, port.FilterUserPayload{IDs: []domain.ID{arg.ID}})
	if err != nil {
		return domain.User{}, intoException(err)
	}

	return updated, nil
}
