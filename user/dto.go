package user

type DTO struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginDTO struct {
	Token string `json:"token"`
}

func toDTO(e Entity) DTO {
	return DTO{
		Id:       e.Id,
		Email:    e.Email,
		Username: e.Username,
		Role:     e.Role,
	}
}

func fromDTO(d DTO) Entity {
	return Entity{
		Email:    d.Email,
		Username: d.Username,
		Password: d.Password,
		Role:     d.Role,
	}
}
