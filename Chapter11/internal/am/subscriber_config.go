package am

import (
	"time"
)

type AckType int

const (
	AckTypeAuto AckType = iota
	AckTypeManual
)

var defaultAckWait = 30 * time.Second
var defaultMaxRedeliver = 5

type SubscriberConfig struct {
	msgFilter    []string
	groupName    string
	ackType      AckType
	ackWait      time.Duration
	maxRedeliver int
}

func NewSubscriberConfig(options []SubscriberOption) SubscriberConfig {
	cfg := SubscriberConfig{
		msgFilter:    []string{},
		groupName:    "",
		ackType:      AckTypeManual,
		ackWait:      defaultAckWait,
		maxRedeliver: defaultMaxRedeliver,
	}

	for _, option := range options {
		option.configureSubscriberConfig(&cfg)
	}

	return cfg
}

type SubscriberOption interface {
	configureSubscriberConfig(*SubscriberConfig)
}

func (c SubscriberConfig) MessageFilters() []string {
	return c.msgFilter
}

func (c SubscriberConfig) GroupName() string {
	return c.groupName
}

func (c SubscriberConfig) AckType() AckType {
	return c.ackType
}

func (c SubscriberConfig) AckWait() time.Duration {
	return c.ackWait
}

func (c SubscriberConfig) MaxRedeliver() int {
	return c.maxRedeliver
}

type MessageFilter []string

func (s MessageFilter) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.msgFilter = s
}

type GroupName string

func (n GroupName) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.groupName = string(n)
}

func (t AckType) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.ackType = t
}

type AckWait time.Duration

func (w AckWait) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.ackWait = time.Duration(w)
}

type MaxRedeliver int

func (i MaxRedeliver) configureSubscriberConfig(cfg *SubscriberConfig) {
	cfg.maxRedeliver = int(i)
}
