package userdomain

type Repository interface {
	FindSignUpUserByUserID(userID string) (*SignUpUser, error)
	FindSignUpUserByEmail(email string) (*SignUpUser, error)
	CreateSignUpUser(user *SignUpUser) error
	DeleteSignUpUser(signUpUser *SignUpUser) error
}
