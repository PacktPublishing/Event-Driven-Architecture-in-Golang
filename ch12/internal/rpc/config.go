package rpc

import (
	"fmt"
	"strings"
)

type Services map[string]string

type RpcConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     string `default:":9000"`
	Services Services
}

func (c RpcConfig) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

func (c RpcConfig) Service(service string) string {
	if address, exists := c.Services[service]; exists {
		return address
	}
	return c.Address()
}

func (s *Services) Decode(v string) error {
	services := map[string]string{}

	pairs := strings.Split(v, ",")
	for _, pair := range pairs {
		p := strings.TrimSpace(pair)
		if len(p) == 0 {
			continue
		}
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			return fmt.Errorf("invalid service pair: %q", p)
		}
		services[strings.ToUpper(kv[0])] = kv[1]
	}

	*s = services
	return nil
}
