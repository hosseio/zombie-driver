package bootstrap

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"server"`
	Nsq struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"nsq"`
	urlConfig URLConfig
}

func (c Config) UrlConfig() URLConfig {
	return c.urlConfig
}

func (c Config) ReadURLConfiguration(filename string) error {
	var urlConfig URLConfig
	endpointsFromFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(endpointsFromFile, &urlConfig)
	if err != nil {
		return err
	}

	c.urlConfig = urlConfig

	return urlConfig.Validate()
}

type URLConfig struct {
	Urls []Url `yaml:"urls"`
}

func (c URLConfig) Validate() error {
	for _, url := range c.Urls {
		if err := url.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Url struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Http   Http   `yaml:"http"`
	Nsq    Nsq    `yaml:"nsq"`
}

type Http struct {
	Host string `yaml:"host"`
}

type Nsq struct {
	Topic string `yaml:"topic"`
}

var (
	ErrURLPathEmpty      = errors.New("url path is empty")
	ErrURLMethodEmpty    = errors.New("url methop is empty")
	ErrURLIsHTTPAndNSQ   = errors.New("url has config for http and nsq")
	ErrURLHasNoHTTPOrNSQ = errors.New("missing http or nsq for url")
)

func (u Url) Validate() error {
	if u.Path == "" {
		return ErrURLPathEmpty
	}
	if u.Method == "" {
		return ErrURLMethodEmpty
	}
	if u.Http.Host == "" && u.Nsq.Topic == "" {
		return ErrURLIsHTTPAndNSQ
	}
	if u.Http.Host != "" && u.Nsq.Topic != "" {
		return ErrURLHasNoHTTPOrNSQ
	}

	return nil
}
