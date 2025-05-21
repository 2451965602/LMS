package config

type server struct {
	Addr string
	Port int64
}

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type config struct {
	Server server
	MySQL  mySQL
}
