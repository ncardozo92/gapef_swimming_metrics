package user

import (
	"context"

	"github.com/ncardozo92/gapef_swimming_metrics/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	USER_COLLECTION string = "users"
)

type Repository interface {
	FindByUsername(id string) (Entity, error, bool)
	GetUsers(page, size int64) ([]Entity, error)
	Create(entity Entity) error
}

type UserRepository struct {
	Database *mongo.Database
}

// creates a UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{Database: persistence.GetDatabase()}
}

func (repository UserRepository) FindByUsername(username string) (Entity, error, bool) {
	user := Entity{}
	mongoErr := repository.Database.Collection(USER_COLLECTION).FindOne(context.TODO(), bson.D{{Key: "username", Value: username}}).Decode(&user)

	if mongoErr != nil {
		var notFound bool
		if mongoErr == mongo.ErrNoDocuments {
			notFound = true
		} else {
			notFound = false
		}
		return user, mongoErr, notFound
	}

	return user, nil, false
}

// gets all the users in the collections in pages
func (repository UserRepository) GetUsers(page, size int64) ([]Entity, error) {
	actualContext := context.TODO()
	usersList := []Entity{}
	findingOptions := options.Find().SetLimit(size).SetSkip((page * size) + 1)

	usersCursor, findUsersErr := repository.Database.Collection(USER_COLLECTION).Find(actualContext, bson.D{}, findingOptions)

	if findUsersErr != nil {
		return nil, findUsersErr
	}

	for usersCursor.Next(actualContext) {
		user := Entity{}

		if cursorErr := usersCursor.Decode(&user); cursorErr != nil {
			return nil, cursorErr
		}

		usersList = append(usersList, user)
	}

	return usersList, nil
}

func (repository UserRepository) Create(entity Entity) error {

	_, insertErr := repository.Database.Collection(USER_COLLECTION).InsertOne(context.TODO(), entity)

	if insertErr != nil {
		return insertErr
	} else {
		return nil
	}
}
