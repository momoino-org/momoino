package usermgt

import "wano-island/common/core"

type UserModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasUpdatedAtColumn

	Username  string `gorm:"type:string;size:64;not null;unique"`
	Email     string `gorm:"type:string;size:256;not null;unique"`
	Password  string `gorm:"type:string;size:256;not null"`
	FirstName string `gorm:"type:string;size:64"`
	LastName  string `gorm:"type:string;size:64"`
	Locale    string `gorm:"type:string;not null;default:en"`
}

func (UserModel) TableName() string {
	return "public.users"
}
