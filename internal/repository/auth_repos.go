package repository

import (
	"carWash/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type UserAuthRepos struct {
	db *sqlx.DB
}

func NewUserAuthRepos(db *sqlx.DB) *UserAuthRepos {
	return &UserAuthRepos{db: db}
}

func (u *UserAuthRepos) VerifyExistenceUser(phone string, activated bool) (*domain.User, error) {

	var input domain.User

	query := fmt.Sprintf("SELECT id, user_name FROM %s WHERE phone_number = $1 AND is_activated = $2", userTable)

	err := u.db.Get(&input, query, phone, activated)

	if err != nil {
		return nil, fmt.Errorf("repository.VerifyExistenceUser: %w", err)
	}

	return &input, nil
}

func (u *UserAuthRepos) UpdateUser(inp domain.User, id int) error {

	setValues := make([]string, 0, reflect.TypeOf(domain.User{}).NumField())

	if inp.Name != "" {
		setValues = append(setValues, fmt.Sprintf("user_name=:user_name"))
	}

	if inp.Password != "" {
		setValues = append(setValues, fmt.Sprintf("password=:password"))
	}

	if inp.PhoneNumber != "" {
		setValues = append(setValues, fmt.Sprintf("phone_number=:phone_number"))
	}
	if inp.UserType != "" {
		setValues = append(setValues, fmt.Sprintf("user_type=:user_type"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return fmt.Errorf("repository.Update: %w", errors.New("empty body"))
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", userTable, setQuery, id)

	result, err := u.db.NamedExec(query, inp)

	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.UpdateUser: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.UpdateUser: %w", domain.ErrNotFound)
	}

	return nil
}

func (u *UserAuthRepos) CreateUser(user domain.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s(user_name,phone_number,password,user_type) VALUES($1,$2,$3,$4) RETURNING id", userTable)

	err := u.db.QueryRowx(query, user.Name, user.PhoneNumber, user.Password, user.UserType).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.CreateUser: %w", err)
	}

	return id, nil
}

func (u *UserAuthRepos) Verify(phone string) error {

	tx := u.db.MustBegin()
	var id int

	fmt.Println(phone)

	query := fmt.Sprintf("UPDATE %s SET is_activated = $1 WHERE phone_number = $2 AND is_activated = $3 RETURNING id", userTable)

	err := tx.QueryRowx(query, true, phone, false).Scan(&id)

	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("repository.Verify: %w", err)
		}
		return fmt.Errorf("repository.Verify: %w", domain.ErrUserAlreadyExist)
	}

	queryInsert := fmt.Sprintf("INSERT INTO %s(user_id) VALUES($1)", sessionTable)
	_, err = tx.Exec(queryInsert, id)

	if err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("repository.Verify: %w", err)
		}
		return fmt.Errorf("repository.Verify: %w", err)
	}
	return tx.Commit()
}

func (u *UserAuthRepos) GetUser(id int) (*domain.User, error) {
	var inp domain.User
	query := fmt.Sprintf("SELECT id,user_name,email,phone_number FROM %s WHERE id = $1", userTable)

	err := u.db.Get(&inp, query, id)

	if err != nil {
		return nil, fmt.Errorf("repository.GetUser: %w", err)
	}
	return &inp, nil
}

func (u *UserAuthRepos) SignIn(phone, password string) (*domain.User, error) {

	var input domain.User
	query := fmt.Sprintf("SELECT id,user_type FROM %s WHERE phone_number = $1 AND password = $2 AND is_activated = $3", userTable)

	err := u.db.Get(&input, query, phone, password, true)

	if err != nil {
		return nil, fmt.Errorf("repository.SignIn: %w", domain.ErrUserDoesNotExist)
	}

	return &input, nil
}

func (u *UserAuthRepos) SetSession(userId int, session domain.Session) error {

	setValues := make([]string, 0, reflect.TypeOf(domain.Session{}).NumField())

	if session.RefreshToken != "" {
		setValues = append(setValues, fmt.Sprintf("refresh_token=:refresh_token"))
	}

	if !session.ExpiresAt.IsZero() {
		setValues = append(setValues, fmt.Sprintf("expires_at=:expires_at"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return fmt.Errorf("repository.SetSession: %v", "empty body")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE user_id=%d", sessionTable, setQuery, userId)

	result, err := u.db.NamedExec(query, session)

	if err != nil {
		return fmt.Errorf("repository.SetSession: %w", err)
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("repository.SetSession: %w", err)
	}

	if affected == 0 {
		// fmt.Errorf
		return sql.ErrNoRows
	}

	return nil
}

func (u *UserAuthRepos) VerifyViaPassword(id int, password string) error {
	var input domain.User

	query := fmt.Sprintf("SELECT user_name FROM %s WHERE password = $1 AND id = $2", userTable)
	err := u.db.Get(&input, query, password, id)
	if err != nil {

		return fmt.Errorf("repository.VerifyViaPassword: %w", domain.ErrInvalidPassword)
	}
	return nil
}

func (u *UserAuthRepos) VerifyViaPhoneNumber(phone string) (*domain.User, error) {

	var input domain.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE phone_number = $1 AND is_activated = $2", userTable)

	err := u.db.Get(&input, query, phone, true)

	if err != nil {
		return nil, fmt.Errorf("repository.VerifyViaPhoneNumber: %w", domain.ErrUserNotRegistered)
	}

	return &input, nil
}

func (u *UserAuthRepos) ResetPassword(phone, password string) error {

	query := fmt.Sprintf("UPDATE %s SET password = $1 WHERE phone_number = $2 ", userTable)

	rows, err := u.db.Exec(query, password, phone)

	if err != nil {
		return fmt.Errorf("repository.ResetPassword: %w", err)
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.ResetPassword: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.ResetPassword: %w", sql.ErrNoRows)
	}
	return nil

}

func (u *UserAuthRepos) SetPassword(id int, hashedOldPassword, hashedNewPassword string) error {
	query := fmt.Sprintf("UPDATE %s SET password = $1 WHERE password = $2 AND id = $3", userTable)

	_, err := u.db.Exec(query, hashedNewPassword, hashedOldPassword, id)

	if err != nil {
		return fmt.Errorf("repository.SetPassword: %w", err)
	}
	return nil

}
