package entity

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"damir/internal/validator"
)

var AnonymousUser = &User{}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID        	int64     	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	Updatedat 	time.Time 	`json:"updated_at"`
	Name      	string    	`json:"name"`
	Surname   	string    	`json:"surname"`
	Email     	string    	`json:"email"`
	Password  	Password  	`json:"-"`
	Activated 	bool      	`json:"activated"`
	Version   	int       	`json:"-"`
	Role 		string		`json:"role"`
	Balance 	int32			`json:"balance"`
}

type Password struct {
	Plaintext *string
	Hash      []byte
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *Password) SetFromHash(hash []byte) error {
    p.Hash = hash
    return nil
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)
	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}
	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}