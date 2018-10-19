package types

type Config struct {
	WebUrl         string
	Port           int
	KeyFile        string
	CertFile       string
	Database       string
	User           string
	Password       string
	SendgridKey    string
	SendgridEmail  string
	SendgridSender string
}
