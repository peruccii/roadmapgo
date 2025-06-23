package models

type PaymentStatus int

const (
	StatusPENDING PaymentStatus = iota
	StatusAPPROVED
	StatusCANCELLED
)

func (p PaymentStatus) String() string {
	switch p {
	case StatusAPPROVED:
		return "APPROVED"
	case StatusPENDING:
		return "PENDING"
	case StatusCANCELLED:
		return "CANCELLED"
	}

	return "UNKNOWN"
}

type Payment struct {
	ID       int64         `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Amount   float64       `json:"amount" db:"amount"`
	Status   PaymentStatus `json:"payment_status" db:"payment_status"`
	Currency string
}
