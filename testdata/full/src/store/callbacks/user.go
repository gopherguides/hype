package demo

type User struct {
	UID int
}

func (u User) ID() int {
	return u.UID
}

func (u User) Validate() error {
	return nil
}
