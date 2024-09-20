package entity

type Step struct {
	ChatID    int64    `json:"chat_id" gorm:"column:chat_id"`
	Step      StepType `json:"step"  gorm:"column:step"`
	Username  string   `json:"username" gorm:"column:username"`
	Firstname string   `json:"firstname" gorm:"column:firstname"`
	Lastname  string   `json:"lastname" gorm:"column:lastname"`
	Phone     string   `json:"phone" gorm:"column:phone"`
}

type StepType int64

const (
	StepPhone StepType = iota + 1
	StepFirstName
	StepLastname
	StepSchool
)

var validSteps = map[StepType]struct{}{
	StepPhone:     {},
	StepFirstName: {},
	StepLastname:  {},
	StepSchool:    {},
}

func (s StepType) IsValid() bool {
	_, exists := validSteps[s]
	return exists
}

func (s StepType) Next() StepType {
	switch s {
	case StepSchool:
		return s
	default:
		return s + 1
	}
}

func (s StepType) Prev() StepType {
	switch s {
	case StepPhone:
		return s
	default:
		return s - 1
	}
}
