package usecase

const (
	StatusUndefined uint = 0
	StatusCreated   uint = 1
	StatusInProcess uint = 2
	StatusCompleted uint = 3
	StatusFallBack  uint = 4
	StatusError     uint = 255
)

const (
	SagaRegistration uint = 0
)
