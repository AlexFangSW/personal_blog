package entities

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// encrypted password
	Password string `json:"password"`
	JWT      string `json:"jwt"`
}

func NewUser(name, password, jwt string) *User {
	return &User{
		ID:       0,
		Name:     name,
		Password: password,
		JWT:      jwt,
	}
}

type InUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewInUser(name, password string) *InUser {
	return &InUser{
		Name:     name,
		Password: password,
	}
}
