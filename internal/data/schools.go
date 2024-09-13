package data

import (
	"time"

	"Apis_go.sahil.net/internal/validator"
)

type School struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	Contact   string    `json:"contract"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email,omitempty"`
	Website   string    `json:"website,omitempty"`
	Address   string    `json:"address"`
	Mode      []string  `json:"mode"`
	Version   int32     `json:"version"`
}

func ValidateSchool(v *validator.Validator, input *School) {
	//use check method to execute validation check
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(len(input.Name) <= 200, "name", "must not be more than 200 bytes")

	v.Check(input.Level != "", "level", "must be provided")
	v.Check(len(input.Level) <= 200, "level", "must not be more than 200 bytes")

	v.Check(input.Contact != "", "contact", "must be provided")
	v.Check(len(input.Contact) <= 200, "contact", "must not be more than 200 bytes")

	v.Check(input.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(input.Phone, validator.PhoneRX), "phone", "must be valid phone number")

	v.Check(input.Email != "", "email", "must be provided")
	v.Check(validator.Matches(input.Email, validator.EmailRX), "email", "must be valid email address")

	v.Check(input.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(input.Website), "website", "must be valid phone number")

	v.Check(input.Address != "", "address", "must be provided")
	v.Check(len(input.Address) <= 500, "address", "must not be more than 500 bytes")

	v.Check(input.Mode != nil, "mode", "must be provided")
	v.Check(len(input.Mode) >= 1, "mode", "must contain atleast one mode")
	v.Check(len(input.Mode) <= 5, "mode", "must contain atmost five mode")
	v.Check(validator.Unique(input.Mode), "mode", "must not contain duplicate entries")

}
