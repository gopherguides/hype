package demo

type Validatable interface {
	Model
	Validate() error
}
