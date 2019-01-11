package types

type Config struct {
	WebUrl          string
	Address         string
	Port            int
	ApiPrefix       string
	KeyFile         string
	CertFile        string
	DatabaseAddress string
	Database        string
	User            string
	Password        string
	SendgridKey     string
	SendgridEmail   string
	SendgridSender  string
}
