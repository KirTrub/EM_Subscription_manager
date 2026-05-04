package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SubscriptionDate time.Time

const DateFormat = "01-2006"

func (sd SubscriptionDate) IsZero() bool {
	return time.Time(sd).IsZero()
}

func (sd SubscriptionDate) Before(u SubscriptionDate) bool {
	return time.Time(sd).Before(time.Time(u))
}

func (sd SubscriptionDate) After(u SubscriptionDate) bool {
	return time.Time(sd).After(time.Time(u))
}

func (sd SubscriptionDate) Time() time.Time {
	return time.Time(sd)
}

func (sd *SubscriptionDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == "" || s == `""` {
		return nil
	}
	t, err := time.Parse(`"`+DateFormat+`"`, s)
	if err != nil {
		return err
	}
	*sd = SubscriptionDate(t)
	return nil
}

func (sd SubscriptionDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(sd).Format(DateFormat))), nil
}

func (sd SubscriptionDate) Value() (driver.Value, error) {
	return time.Time(sd), nil
}

func (sd *SubscriptionDate) Scan(value interface{}) error {
	if t, ok := value.(time.Time); ok {
		*sd = SubscriptionDate(t)
		return nil
	}
	return fmt.Errorf("cannot scan %T into SubscriptionDate", value)
}

type Subscription struct {
	Id          int               `json:"id"`
	ServiceName string            `json:"service_name"`
	Price       int               `json:"price"`
	UserId      uuid.UUID         `json:"user_id"`
	StartDate   SubscriptionDate  `json:"start_date"`
	EndDate     *SubscriptionDate `json:"end_date,omitempty"`
}

func ParseMonthYear(s string) (time.Time, error) {
	return time.Parse(DateFormat, s)
}
