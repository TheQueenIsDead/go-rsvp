package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const EventTimeFormat = "2006-01-02T15:04"

var (
	nilTime = (time.Time{}).UnixNano()
)

type EventDate struct {
	datatypes.Date
}
type EventTime struct {
	datatypes.Time
}

type Event struct {
	gorm.Model
	Date             EventDate `form:"date"`
	Time             EventTime `form:"time"`
	Name             string    `form:"name"`
	Description      string    `form:"description"`
	MinimumAttendees int8      `form:"minimumAttendees"`
}

func (ed *EventDate) UnmarshalParam(param string) error {

	t, err := time.Parse(`2006-01-02`, param)
	if err != nil {
		return err
	}

	date := datatypes.Date(t)

	ed.Date = date
	return nil
}
func (et *EventTime) UnmarshalParam(param string) error {

	t, err := time.Parse(`15:04`, param)
	if err != nil {
		return err
	}

	tt := datatypes.NewTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())

	et.Time = tt
	return nil
}

//type EventTime struct {
//	Time     time.Time
//	cdcdscds datatypes.Time
//}
//
//func (et *EventTime) Scan(src interface{}) (err error) {
//	switch v := src.(type) {
//	case []byte:
//		//et.setFromString(string(v))
//		et.Time, err = time.Parse(EventTimeFormat, string(v))
//	case string:
//		//et.setFromString(v)
//		et.Time, err = time.Parse(EventTimeFormat, v)
//	case time.Time:
//		et.Time = v
//	default:
//		return errors.New(fmt.Sprintf("failed to scan value: %v", v))
//	}
//
//	return err
//}
//
//// Value implements driver.Valuer interface and returns string format of Time.
//func (et EventTime) Value() (driver.Value, error) {
//	return et.String(), nil
//}
//
//// String implements fmt.Stringer interface.
//func (et EventTime) String() string {
//	return et.Time.GoString()
//
//}
//
//func (et *EventTime) UnmarshalJSON(b []byte) (err error) {
//	s := strings.Trim(string(b), "\"")
//	if s == "null" {
//		et.Time = time.Time{}
//		return
//	}
//	et.Time, err = time.Parse(EventTimeFormat, s)
//	return
//}
//
//func (et *EventTime) MarshalJSON() ([]byte, error) {
//	if et.Time.UnixNano() == nilTime {
//		return []byte("null"), nil
//	}
//	return []byte(fmt.Sprintf("\"%s\"", et.Time.Format(EventTimeFormat))), nil
//}
//
////// EventDate is a custom type that extends the sql ORM library datatype for database operation compatibility,
////// while allowing us to implement unmarshalling functionality as part of the echo web framework struct binding process.
////type EventDate struct {
////	datatypes.Time
////}
//
//// UnmarshalParam decodes and assigns a value from a form or query param.
//// It is invoked as part of the web framework entity binding process, and is triggered by the inclusion of the `form`
//// tag on an EventDate object
//func (et *EventTime) UnmarshalParam(src string) error {
//	t, err := time.Parse("2006-01-02T15:04", src)
//	if err != nil {
//		log.WithError(err).WithField("source", src).Error("could not unmarshal eventdate param")
//		return err
//	}
//
//	et.Time = t
//
//	// t := time.Time(dbRecord.Date)
//	return nil
//}
