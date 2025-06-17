package models

type Courses struct {
	ID    int64  `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Title string `json:"title" db:"title" gorm:"type:varchar(255);not null"`
	Users []User `json:"users_cadastred" gorm:"many2many:user_courses;"`
}
