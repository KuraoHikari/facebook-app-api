package entity

type Post struct {
	ID          uint64 `gorm:"primary_key:auto_increment" json:"id"`
	ImageLink   string `gorm:"type:text" json:"image_link"`
	Description string `gorm:"not null;type:text" json:"description"`
	UserID      uint64 `gorm:"not null" json:"-"`
	User        User   `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
}