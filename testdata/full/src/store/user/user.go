package demo

type User struct {
	UID int
}

func (u User) ID() int {
	return u.UID
}
