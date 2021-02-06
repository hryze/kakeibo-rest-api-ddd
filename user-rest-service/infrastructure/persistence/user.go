package persistence

import (
	"database/sql"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/domain/userdomain"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/errors"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/config"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-rest-api-ddd/user-rest-service/infrastructure/persistence/datasource"
)

type userRepository struct {
	*config.RedisHandler
	*config.MySQLHandler
}

func NewUserRepository(redisHandler *config.RedisHandler, mysqlHandler *config.MySQLHandler) *userRepository {
	return &userRepository{redisHandler, mysqlHandler}
}

func (r *userRepository) FindSignUpUserByUserID(userID string) (*userdomain.SignUpUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email,
            password
        FROM
            users
        WHERE
            user_id = ?`

	var signUpUserDto datasource.SignUpUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, userID).StructScan(&signUpUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	signUpUser := userdomain.NewSignUpUserFromDataSource(signUpUserDto.UserID, signUpUserDto.Name, signUpUserDto.Email, signUpUserDto.Password)

	return signUpUser, nil
}

func (r *userRepository) FindSignUpUserByEmail(email string) (*userdomain.SignUpUser, error) {
	query := `
        SELECT
            user_id,
            name,
            email,
            password
        FROM
            users
        WHERE
            email = ?`

	var signUpUserDto datasource.SignUpUser
	if err := r.MySQLHandler.Conn.QueryRowx(query, email).StructScan(&signUpUserDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	signUpUser := userdomain.NewSignUpUserFromDataSource(signUpUserDto.UserID, signUpUserDto.Name, signUpUserDto.Email, signUpUserDto.Password)

	return signUpUser, nil
}

func (r *userRepository) CreateSignUpUser(signUpUser *userdomain.SignUpUser) error {
	query := `
        INSERT INTO users
            (user_id, name, email, password)
        VALUES
            (?,?,?,?)`

	if _, err := r.MySQLHandler.Conn.Exec(query, signUpUser.UserID(), signUpUser.Name(), signUpUser.Email(), signUpUser.Password()); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteSignUpUser(signUpUser *userdomain.SignUpUser) error {
	query := `
        DELETE
        FROM
            users
        WHERE
            user_id = ?`

	_, err := r.MySQLHandler.Conn.Exec(query, signUpUser.UserID())

	return err
}
