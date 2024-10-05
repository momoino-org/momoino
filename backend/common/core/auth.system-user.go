package core

import (
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type systemUser struct{}

var _ PrincipalUser = (*systemUser)(nil)

func (u systemUser) GetID() uuid.UUID {
	return uuid.Nil
}

func (u systemUser) GetUsername() string {
	return "system"
}

func (u systemUser) GetEmail() string {
	return "system@internal.com"
}

func (u systemUser) GetGivenName() string {
	return "System"
}

func (u systemUser) GetFamilyName() string {
	return ""
}

func (u systemUser) GetLocale() string {
	return language.English.String()
}

func (u systemUser) GetRoles() []string {
	return []string{}
}

func (u systemUser) GetPermissions() []string {
	return []string{}
}

func NewSystemUser() *systemUser {
	return &systemUser{}
}
