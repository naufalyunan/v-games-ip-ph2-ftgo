package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID    uint        `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CartID;references:ID" json:"cart_id"`
	GameID    uint        `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:GameID;references:ID" json:"game_id"`
	StartDate *CustomDate `gorm:"type:datetime; not null" json:"start_date"`
	EndDate   *CustomDate `gorm:"type:datetime; not null" json:"end_date"`
	Quantity  int         `gorm:"type:int; not null" json:"quantity"`

	//Assc
	Cart Cart  `gorm:"foreignKey:CartID;references:ID" json:"-"`
	Game *Game `gorm:"foreignKey:GameID;references:ID" json:"game,omitempty"`
}

// ValidateDates checks if end_date is before start_date
func (ci *CartItem) ValidateDates() error {
	if time.Time(*ci.EndDate).Before(time.Time(*ci.StartDate)) {
		return errors.New("end_date cannot be before start_date")
	}
	return nil
}

func (ci *CartItem) CalculateDaysDifference() float64 {
	start := time.Time(*ci.StartDate)
	end := time.Time(*ci.EndDate)

	// Calculate the difference in days
	duration := end.Sub(start)
	days := int(duration.Hours() / 24) // Convert duration to days

	return math.Abs(float64(days))
}

type CustomDate time.Time

const CustomDateFormat = "2006-01-02"

// Value implements the driver.Valuer interface
func (cd CustomDate) Value() (driver.Value, error) {
	return time.Time(cd).Format(CustomDateFormat), nil
}

// Scan implements the sql.Scanner interface
func (cd *CustomDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*cd = CustomDate(v)
		return nil
	case string:
		t, err := time.Parse(CustomDateFormat, v)
		if err != nil {
			return err
		}
		*cd = CustomDate(t)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

// MarshalJSON for JSON serialization
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(cd).Format(CustomDateFormat))
}

// UnmarshalJSON for JSON deserialization
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	str := string(b)
	t, err := time.Parse(`"`+CustomDateFormat+`"`, str)
	if err != nil {
		return err
	}
	*cd = CustomDate(t)
	return nil
}
