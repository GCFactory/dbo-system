package models

import "github.com/google/uuid"

type Account struct {
	Acc_uuid         uuid.UUID `json:"acc_uuid" db:"acc_uuid" validate:"len=36 required unique"`
	Acc_status       uint8     `json:"acc_status" db:"acc_status" validate:"min=0 max=255 required oneof=0 10 22 30 40 50 255"`
	Acc_culc_number  string    `json:"acc_culc_number" db:"acc_culc_number" validate:"len=20 required"`
	Acc_corr_number  string    `json:"acc_corr_number" db:"acc_corr_number" validate:"len=20 required"`
	Acc_bic          string    `json:"acc_bic" db:"acc_bic" validate:"len=9 required"`
	Acc_cio          string    `json:"acc_cio" db:"acc_cio" validate:"len=9 required"`
	Acc_money_value  uint      `json:"acc_money_value" db:"acc_money_value" validate:"min=0 max=1 oneof=0 1 2 3"`
	Acc_money_amount float64   `json:"acc_money_amount" db:"acc_money_amount" validate:"required"`
}

type ReserverReason struct {
	Acc_uuid uuid.UUID `json:"acc_uuid" db:"acc_uuid" validate:"len=36 required unique"`
	Reason   string    `json:"reserve_reason" db:"reserve_reason" validate:"required"'`
}
