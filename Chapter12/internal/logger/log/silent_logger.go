package log

import (
	std "log"
)

type SilentLogger struct{}

func (*SilentLogger) Fatal(v ...interface{})                 { std.Fatal(v...) }
func (*SilentLogger) Fatalf(format string, v ...interface{}) { std.Fatalf(format, v...) }
func (*SilentLogger) Print(...interface{})                   {}
func (*SilentLogger) Println(...interface{})                 {}
func (*SilentLogger) Printf(string, ...interface{})          {}
