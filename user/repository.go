package user

import (
	"context"

	"github.com/ncardozo92/gapef_swimming_metrics/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const USER_COLLECTION string = "users"

type Repository interface {
	doSomething() string
	FindByUsername(id string) (User, error, bool)
}

type UserRepository struct {
	Database *mongo.Database
}

// creates a UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{Database: persistence.GetDatabase()}
}

func (repository UserRepository) FindByUsername(username string) (User, error, bool) {
	user := User{}
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

// Method for testing purposes
func (ur *UserRepository) doSomething() string {
	return "Doing something..."
}
