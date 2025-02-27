package domain

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey:autoIncrement"`
	Email    string `json:"email" gorm:"unique; not null"`
	Password string `json:"password" gorm:"not null"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
