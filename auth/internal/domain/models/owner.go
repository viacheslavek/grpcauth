package models

type Owner struct {
	Id       int64
	Email    string
	Login    string
	Password string
	PassHash []byte
}

// Option is a function that configures an Owner
type Option func(*Owner)

// NewOwner creates a new Owner with the provided options
func NewOwner(opts ...Option) *Owner {
	owner := &Owner{}
	for _, opt := range opts {
		opt(owner)
	}
	return owner
}

// WithId sets the Id of the Owner
func WithId(id int64) Option {
	return func(o *Owner) {
		o.Id = id
	}
}

// WithEmail sets the Email of the Owner
func WithEmail(email string) Option {
	return func(o *Owner) {
		o.Email = email
	}
}

// WithLogin sets the Login of the Owner
func WithLogin(login string) Option {
	return func(o *Owner) {
		o.Login = login
	}
}

// WithPassword sets the Password of the Owner
func WithPassword(password string) Option {
	return func(o *Owner) {
		o.Password = password
	}
}
