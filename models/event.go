package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const EventTimeFormat = "2023-07-21T11:11"

type EventTimestamp time.Time

type Event struct {
	gorm.Model
	//Date             datatypes.Date
	Date             EventDate `form:"date"`
	Name             string    `form:"name"`
	Description      string    `form:"description"`
	MinimumAttendees int8      `form:"minimumAttendees"`
}

//type EventBinder struct{}

type EventDate struct {
	datatypes.Date
}

func (ed *EventDate) UnmarshalParam(src string) error {
	//ts, err := time.Parse(time.RFC3339, src)
	t, err := time.Parse("2006-01-02T15:04", src)
	if err != nil {
		log.WithError(err).WithField("source", src).Error("could not unmarshal eventdate param")
		return err
	}

	*ed = EventDate{datatypes.Date(t)}
	return nil
}

/*func (eb *EventBinder) BindBody(i interface{}, c echo.Context) error {

	// Intercept form date fields and convert them to golang time.Time values
	if date := c.FormValue("date"); date != "" {
		format := "2006-01-02T15:04"
		ts, err := time.Parse(format, date)
		if err != nil {
			return err
		}
		c.Request().Form.Set("date", ts.String())
	}

	// Pass to default binder after converting the timestamp
	db := new(echo.DefaultBinder)
	if err := db.BindBody(c, i); err != nil {
		return err
	}

	return nil
}
func (eb *EventBinder) Bind(i interface{}, c echo.Context) (err error) {

	if c.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationForm {
		err = eb.BindBody(i, c)
		return
	}

	db := new(echo.DefaultBinder)
	if err := db.Bind(i, c); err != nil {
		return err
	}
	//contentType := c.Request().Header.Get(echo.HeaderContentType)
	//err = echo.ErrUnsupportedMediaType
	//if strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
	//	if err = json.NewDecoder(c.Request().Body).Decode(i); err != nil {
	//		return http.ErrBodyNotAllowed
	//	}
	//	if tagName := c.Get("ValidateTagName"); tagName == nil {
	//		c.Set("ValidateTagName", "validate")
	//	}
	//	config := &validator.Config{TagName: c.Get("ValidateTagName").(string)}
	//	validate := validator.New(config)
	//	if err = validate.Struct(i); err != nil {
	//		return err
	//	}
	//}
	return err
}*/
