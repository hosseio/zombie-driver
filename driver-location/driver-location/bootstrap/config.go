package bootstrap

type Config struct {
	Server struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"server"`
	Nsq struct {
		Addr    string `mapstructure:"addr"`
		Topic   string `mapstructure:"topic"`
		Channel string `mapstructure:"channel"`
	} `mapstructure:"nsq"`
	Redis string `mapstructure:"redis"`
}
