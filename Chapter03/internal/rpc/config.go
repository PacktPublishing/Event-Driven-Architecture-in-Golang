package rpc

import (
	"fmt"
)

type RpcConfig struct {
	Host string `default:"0.0.0.0"`
	Port string `default:":8085"`
}

func (c RpcConfig) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}
