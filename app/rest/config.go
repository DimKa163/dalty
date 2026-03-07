package rest

type Config struct {
	Addr         string `env:"ADDR" envDefault:":8080"`
	Database     string `env:"DATABASE,required"`
	UseProfiling bool   `env:"USE_PROFILING" envDefault:"false"`
	PProfAddr    string `env:"PPROF_PORT" envDefault:":6060"`
}
