package graph

type Config struct {
	Addr     string `env:"ADDR" envDefault:":8080"`
	Database string `env:"DATABASE,required"`
}
