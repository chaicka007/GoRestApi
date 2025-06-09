package models

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)

var validStatuses = map[string]bool{
	StatusPending:    true,
	StatusInProgress: true,
	StatusCompleted:  true,
}

func IsValidStatus(s string) bool {
	return validStatuses[s]
}
