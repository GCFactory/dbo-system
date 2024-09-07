package models

import (
	"github.com/google/uuid"
)

type Passport struct {
	Passport_uuid       uuid.UUID `json:"passport_uuid" db:"passport_uuid" validate:"required,uuid4"`
	Passport_series     string    `json:"passport_series" db:"passport_series" validate:"required,len=4,number"`
	Passport_number     string    `json:"passport_number" db:"passport_number" validate:"required,len=6,number"`
	Name                string    `json:"name" db:"name" validate:"required"`
	Surname             string    `json:"surname" db:"surname" validate:"required"`
	Patronimic          string    `json:"patronimic" db:"patronimic" validate:"required"`
	Birth_date          string    `json:"birth_date" db:"birth_date" validate:"required,datetime"`
	Birth_location      string    `json:"birth_location" db:"birth_location" validate:"required"`
	Pick_up_point       string    `json:"pick_up_point" db:"pick_up_point" validate:"required"`
	Authority           string    `json:"authority" db:"authority" validate:"required,len=7"`
	Authority_date      string    `json:"authority_date" db:"authority_date" validate:"required,datetime"`
	Registration_adress string    `json:"registration_adress" db:"registration_adress" validate:"required"`
}

type User struct {
	User_uuid     uuid.UUID `json:"user_uuid" db:"user_uuid" validate:"required,uuid4"`
	Passport_uuid uuid.UUID `json:"passport_uuid" db:"passport_uuid" validate:"required,uuid4"`
	User_inn      string    `json:"user_inn" db:"user_inn" validate:"required,number,len=20"`
	User_accounts string    `json:"user_accounts" db:"user_accounts" validate:"required,json"`
}

type User_full_data struct {
	User     *User
	Passport *Passport
}
