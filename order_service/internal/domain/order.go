package domain

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index"           json:"user_id"`
	ItemName  string    `gorm:"not null"                 json:"item_name"`
	Amount    float64   `gorm:"not null"                 json:"amount"`
	Status    string    `gorm:"default:'PENDING'"        json:"status"`
	CreatedAt time.Time `                                json:"created_at"`
}
