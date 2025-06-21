package user

import "time"

type User struct {
	ID            int            `db:"id" json:"-"`
	UserUUID      string         `db:"user_uuid" json:"user_uuid"`
	FirstName     *string        `db:"first_name" json:"first_name"`
	LastName      *string        `db:"last_name" json:"last_name"`
	UserName      string         `db:"user_name" json:"user_name"`
	Password      string         `db:"password" json:"-"`
	Email         string         `db:"email" json:"email"`
	LoginSession  *string        `db:"login_session" json:"-"`
	ProfilePhoto  *string        `db:"profile_photo" json:"profile_photo"`
	StatusID      int            `db:"status_id" json:"-"`
	Order         *int           `db:"order" json:"-"`
	CreatedBy     int            `db:"created_by" json:"-"`
	CreatedAt     time.Time      `db:"created_at" json:"-"`
	UpdatedBy     *int           `db:"updated_by" json:"-"`
	UpdatedAt     *time.Time     `db:"updated_at" json:"-"`
	DeletedBy     *int           `db:"deleted_by" json:"-"`
	DeletedAt     *time.Time     `db:"deleted_at" json:"-"`
	UserDatabases []UserDatabase `json:"user_databases"`
}

type UserDatabase struct {
	DBUUID string `db:"db_uuid" json:"db_uuid"`
	DBName string `db:"db_name" json:"db_name"`
}

type UserInfoResponse struct {
	UserInfo User `json:"user_info"`
}
