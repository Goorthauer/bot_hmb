package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type SchoolTraining struct {
	SchoolID uuid.UUID `json:"school_id" gorm:"column:school_id"`

	Price       int      `json:"price"  gorm:"column:price"`
	Description string   `json:"description"  gorm:"column:description"`
	Schedule    Schedule `json:"schedule" gorm:"column:schedule;type:jsonb"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

type TrainingDay struct {
	Day         TrainingDayDay  `json:"day"`
	Description string          `json:"description"`
	Time        TrainingDayTime `json:"time"`
}

type TrainingDayTime struct {
	Open   string `json:"open"`
	Closed string `json:"closed"`
}

func (t TrainingDayTime) GetTime() string {
	return fmt.Sprintf(`c %s по %s`, t.Open, t.Closed)
}

type TrainingDayDay string

const (
	TrainingDayMonday    TrainingDayDay = `пн`
	TrainingDayTuesday   TrainingDayDay = `вт`
	TrainingDayWednesday TrainingDayDay = `ср`
	TrainingDayThursday  TrainingDayDay = `чт`
	TrainingDayFriday    TrainingDayDay = `пт`
	TrainingDaySaturday  TrainingDayDay = `сб`
	TrainingDaySunday    TrainingDayDay = `вс`
)

func (t TrainingDayDay) getTimeDay() time.Weekday {
	switch t {
	case TrainingDayTuesday:
		return time.Tuesday
	case TrainingDayWednesday:
		return time.Wednesday
	case TrainingDayThursday:
		return time.Thursday
	case TrainingDayFriday:
		return time.Friday
	case TrainingDaySaturday:
		return time.Saturday
	case TrainingDaySunday:
		return time.Sunday
	default:
		return time.Monday
	}
}

func (s *SchoolTraining) FindNextTrainingDay() (time.Time, string, error) {
	const daysOnWeek = 7
	if len(s.Schedule) == 0 {
		return time.Time{}, "", errors.New("заполните тренировочное расписание корректно")
	}
	trainingDays := make(map[time.Weekday]string)
	for _, v := range s.Schedule {
		trainingDays[v.Day.getTimeDay()] = v.Time.GetTime()
	}
	currentTime := time.Now()
	currentDay := currentTime.Weekday()

	for day, t := range trainingDays {
		// Если тренировка в тот же день или позже по неделе
		if day >= currentDay {
			// Сколько дней до этой тренировки
			daysUntilTraining := day - currentDay
			return currentTime.AddDate(0, 0, int(daysUntilTraining)), t, nil
		}
	}

	// Если все тренировки уже прошли на этой неделе, выбираем первую тренировку на следующей неделе
	firstDay := s.Schedule[0]
	daysUntilNextWeekTraining := (daysOnWeek - int(currentDay)) + int(firstDay.Day.getTimeDay())
	return currentTime.AddDate(0, 0, daysUntilNextWeekTraining), firstDay.Time.GetTime(), nil
}

func (t TrainingDayDay) GetNextWeekThisDay() time.Time {
	currentTime := time.Now()
	// Вычисляем, сколько дней до следующего такого же дня
	daysUntil := (t.getTimeDay() - currentTime.Weekday() + 7) % 7
	if daysUntil == 0 {
		daysUntil = 7
	}
	return currentTime.AddDate(0, 0, int(daysUntil))
}

func (t TrainingDayDay) IsValid() bool {
	return t == TrainingDayMonday ||
		t == TrainingDayTuesday ||
		t == TrainingDayWednesday ||
		t == TrainingDayThursday ||
		t == TrainingDayFriday ||
		t == TrainingDaySaturday ||
		t == TrainingDaySunday

}

func (t TrainingDayDay) String() string {
	return string(t)
}

type Schedule []TrainingDay

// Scan gorm Scanner.
func (d *Schedule) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := []TrainingDay{}
	err := json.Unmarshal(bytes, &result)
	*d = result
	return err
}

// Value gorm Valuer.
func (d Schedule) Value() (driver.Value, error) {
	j, err := json.Marshal(&d)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
