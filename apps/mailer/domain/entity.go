package domain

type Mail struct {
	Sender    string `json:"sender" validate:"required,email,min=3,max=100"`
	Recipient string `json:"recipient" validate:"required,email,min=3,max=100"`
	Subject   string `json:"subject" validate:"required,min=3,max=100"`
	Body      string `json:"body" validate:"required,min=3,max=1000"`
	MIMEType  string `json:"mimeType"`
}
