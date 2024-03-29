package util

import (
	"database/sql/driver"

	"github.com/go-openapi/strfmt"
	"github.com/shopspring/decimal"
)

type Decimal decimal.Decimal

func (Decimal) Validate(formats strfmt.Registry) error {
	return nil
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	base := decimal.Decimal(d)
	return base.MarshalJSON()
}

func (d *Decimal) UnmarshalJSON(value []byte) error {
	base := decimal.Decimal(*d)
	if err := base.UnmarshalJSON(value); err != nil {
		return err
	}
	*d = Decimal(base)
	return nil
}

func (d Decimal) MarshalBinary() (data []byte, err error) {
	base := decimal.Decimal(d)
	return base.MarshalBinary()
}

func (d *Decimal) UnmarshalBinary(data []byte) error {
	base := decimal.Decimal(*d)
	if err := base.UnmarshalBinary(data); err != nil {
		return err
	}
	*d = Decimal(base)
	return nil
}

func (d *Decimal) Scan(value interface{}) error {
	base := decimal.Decimal(*d)
	if err := base.Scan(value); err != nil {
		return err
	}
	*d = Decimal(base)
	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	base := decimal.Decimal(d)
	return base.String(), nil
}
