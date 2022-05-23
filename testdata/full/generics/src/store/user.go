package demo

type User struct {
	Email string
}

func (u User) ID() string {
	return u.Email
}
