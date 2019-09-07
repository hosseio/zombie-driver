package bootstrap

type Config struct {
	Server struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"server"`
	DriverLocation struct {
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"driver_location"`
}
