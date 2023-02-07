package result

type Result int

const (
	OK      Result = 0
	Unknown        = 1
	Warning        = 2
	Error          = 3
)
