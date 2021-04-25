package service

type db struct {
	User   string
	Pass   string
	Addr   string
	DBName string
}

type server struct {
	Port string
	Sign string
}

type path struct {
	Result  string
	History string
	Log     string
}

type conf struct {
	DB     db
	DBDev  db
	Server server
	Path   path
}
