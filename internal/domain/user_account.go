package domain

import (
	"time"

	"github.com/lib/pq"
)

// UserAccount defines the one to many relationship of an user to an account. This
// will enable a single user access to multiple accounts without having duplicate
// users. Each association of a user to an account has a set of roles and a status
// defined for the user. The roles will be applied to enforce ACLs across the
// application. The status will allow users to be managed on by account with users
// being global to the application.
type UserAccount struct {
	ID         int64             `json:"id" validate:"required" example:"72938896-a998-4258-a17b-6418dcdb80e3"`
	UserID     int               `json:"user_id" validate:"required" example:"d69bdef7-173f-4d29-b52c-3edc60baf6a2"`
	AccountID  int               `json:"account_id" validate:"required" example:"c4653bf9-5978-48b7-89c5-95704aebb7e2"`
	Roles      UserAccountRoles  `json:"roles" validate:"required,dive,oneof=manager user" enums:"manager,user" swaggertype:"array,string" example:"manager"`
	Status     UserAccountStatus `json:"status" validate:"omitempty,oneof=active invited disabled" enums:"active,invited,disabled" swaggertype:"string" example:"active"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	ArchivedAt *pq.NullTime      `json:"archived_at,omitempty" swaggertype:"string"`
}

// UserAccountStatus represents the status of a user for an account.
type UserAccountStatus string

// UserAccountStatus values define the status field of a user account.
const (
	// UserAccountStatus_Active defines the state when a user can access an account.
	UserAccountStatus_Active UserAccountStatus = "active"
	// UserAccountStatus_Invited defined the state when a user has been invited to an
	// account.
	UserAccountStatus_Invited UserAccountStatus = "invited"
	// UserAccountStatus_Disabled defines the state when a user has been disabled from
	// accessing an account.
	UserAccountStatus_Disabled UserAccountStatus = "disabled"
)

// UserAccountStatus_Values provides list of valid UserAccountStatus values.
var UserAccountStatus_Values = []UserAccountStatus{
	UserAccountStatus_Active,
	UserAccountStatus_Invited,
	UserAccountStatus_Disabled,
}

// String converts the UserAccountStatus value to a string.
func (s UserAccountStatus) String() string {
	return string(s)
}

// UserAccountRole represents the role of a user for an account.
type UserAccountRole string

// UserAccountRoles represents a set of roles for a user for an account.
type UserAccountRoles []UserAccountRole

// UserAccountRole values define the role field of a user account.
const (
	// UserAccountRole_Manager defines the state of a user when they have manager
	// privileges for accessing an account. This role provides a user with full
	// access to an account.
	UserAccountRole_Manager UserAccountRole = "Manager"
	// UserAccountRole_User defines the state of a user when they have basic
	// privileges for accessing an account. This role provies a user with the most
	// limited access to an account.
	UserAccountRole_User UserAccountRole = "User"
)
