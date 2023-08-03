package models

type ProcessState string

const (
	Created ProcessState = "Created"
	Running ProcessState = "Running"
	Done    ProcessState = "Done"
	Failed  ProcessState = "Failed"
)

type CrawlerProcess struct {
	CPF      string       `json:"cpf"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	State    ProcessState `json:"process_state"`
}

// Separating Id field because it causes error with elasticsearch
type CrawlerProcessWithId struct {
	ID string `json:"_id"`
	CrawlerProcess
}
