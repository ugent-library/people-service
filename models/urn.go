package models

import (
	"strings"
)

type URN struct {
	Namespace string
	Value     string
}

func NewURN(ns string, val string) *URN {
	return &URN{
		Namespace: ns,
		Value:     val,
	}
}

func (urn *URN) Dup() *URN {
	return NewURN(urn.Namespace, urn.Value)
}

func (urn *URN) String() string {
	return "urn:" + urn.Namespace + ":" + urn.Value
}

// ParseURN parses "urn:<namespace>:<value>" into struct models.URN
func ParseURN(v string) (*URN, error) {
	parts := strings.Split(v, ":")
	if len(parts) < 3 {
		return nil, ErrInvalidURN
	}
	return &URN{Namespace: parts[1], Value: strings.Join(parts[2:], ":")}, nil
}

type ByURN []*URN

func (urns ByURN) Len() int {
	return len(urns)
}

func (urns ByURN) Swap(i, j int) {
	urns[i], urns[j] = urns[j], urns[i]
}

func (urns ByURN) Less(i, j int) bool {
	if urns[i].Namespace != urns[j].Namespace {
		return urns[i].Namespace < urns[j].Namespace
	}
	return urns[i].Value < urns[j].Value
}
