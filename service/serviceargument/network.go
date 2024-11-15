package serviceargument

type EmailContent struct {
	Subject    string
	SenderName string
	ToEmails   []string
	Content    string
}
