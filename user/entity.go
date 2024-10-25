package user

type User struct {
	Id       string `bson:"_id"` // tag used by mongoDB driver for mapping data
	Email    string
	Username string
	Password string
	Role     string
}
