package user

type Entity struct {
	Id       string `bson:"_id"` // tag used by mongoDB driver for mapping data
	Email    string `bson:"email"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Role     string `bson:"role"`
}
