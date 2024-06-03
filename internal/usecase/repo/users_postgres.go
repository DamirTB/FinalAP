package repo

import (
	"context"
	"crypto/sha256"
	"damir/internal/filters"
	"damir/internal/entity"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *entity.User) error {
	query := `
	INSERT INTO user_info (fname, sname, email, password_hash, activated)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, version`
	args := []any{user.Name, user.Surname, user.Email, user.Password.Hash, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return entity.ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*entity.User, error) {
	query := `
	SELECT id, created_at, fname, email, password_hash, activated, version
	FROM user_info
	WHERE email = $1`
	var user entity.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, entity.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *entity.User) error {
	query := `
	UPDATE user_info
	SET fname = $1, sname=$2, email = $3, activated = $4, version = version + 1
	WHERE id = $5 
	RETURNING version`
	args := []any{
		user.Name,
		user.Surname,
		user.Email,
		user.Activated,
		user.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return entity.ErrDuplicateEmail

		case errors.Is(err, sql.ErrNoRows):
			return entity.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*entity.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
	SELECT user_info.id, user_info.created_at, user_info.fname, user_info.sname, user_info.email, user_info.password_hash, user_info.user_role, user_info.balance, user_info.activated, user_info.version
	FROM user_info
	INNER JOIN tokens
	ON user_info.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`
	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user entity.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.Password.Hash,
		&user.Role,
		&user.Balance,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, entity.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.
	return &user, nil
}

func (m UserModel) Delete(id int64) error {
	query := `
		DELETE FROM user_info
		where id = $1
	`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return entity.ErrRecordNotFound
	}
	return nil
}

func (m UserModel) Get(id int64) (*entity.User, error) {
	query := `
		SELECT *
		FROM user_info
		WHERE id = $1`

	var user entity.User
	var passwordHash []byte
	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Updatedat,
		&user.Name,
		&user.Surname,
		&user.Email,
		&passwordHash,
		&user.Role,
		&user.Balance,
		&user.Activated, 
		&user.Version,
	)
	err = user.Password.SetFromHash(passwordHash)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, entity.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}


func (m UserModel) GetAll(name string, filters filters.Filters) ([]*entity.User, error) {
	query := fmt.Sprintf(`
	SELECT *
	FROM user_info
	WHERE (to_tsvector('simple', fname) @@ plainto_tsquery('simple', $1) OR $1 = '')
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3`, filters.SortColumn(), filters.SortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{name, filters.Limit(), filters.Offset()}
	// query := `
	// 	SELECT *
	// 	FROM user_info`
	rows, err := m.DB.QueryContext(ctx, query, args...)	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User

	for rows.Next() {
		var user entity.User
		var passwordHash []byte
		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.Updatedat,
			&user.Name,
			&user.Surname,
			&user.Email,
			&passwordHash,
			&user.Role,
			&user.Balance,
			&user.Activated, 
			&user.Version,
		)
		err = user.Password.SetFromHash(passwordHash)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

return users, nil
}

