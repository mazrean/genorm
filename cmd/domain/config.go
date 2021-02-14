package domain

// DBMS dbms(mysql,postgres,etc)
type DBMS string

const (
	// MySQL mysql
	MySQL DBMS = "mysql"
)

// Config 全体に関わる設定(yaml用)
type Config struct {
	DBMS
	Version  string
	Database string
}
