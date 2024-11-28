package user

type Entity struct {
	Id       string `bson:"_id,omitempty"`
	Email    string `bson:"email"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Role     string `bson:"role"`
}
