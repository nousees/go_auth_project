package entities

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey:autoIncrement"`
	Email    string `json:"email" gorm:"not null;uniqueIndex"`
	Password string `json:"password" gorm:"not null"`
}
