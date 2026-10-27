package domain

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null"     json:"username"`
	Password  string    `gorm:"not null"                 json:"-"`
	FullName  string    `gorm:"not null"                 json:"full_name"`
	AvatarURL string    `                                json:"avatar_url"`
	CreatedAt time.Time `                                json:"created_at"`
}
