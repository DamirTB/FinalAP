package entity

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"errors"
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
	Password  	password  	`json:"-"`
	Activated 	bool      	`json:"activated"`
	Version   	int       	`json:"-"`
	Role 		string		`json:"role"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) SetFromHash(hash []byte) error {
    p.hash = hash
    return nil
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
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