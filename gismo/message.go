package gismo

import (
	"time"

	"github.com/ugent-library/people-service/models"
)

/*
	See also https://github.com/ugent-library/soap-bridge/blob/main/main.go
*/

type Attribute struct {
	Name      string     `json:"name"`
	Value     string     `json:"value"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type Message struct {
	ID         string      `json:"id,omitempty"`
	Date       string      `json:"date,omitempty"`
	Language   string      `json:"language"`
	Attributes []Attribute `json:"attributes"`
	Source     string      `json:"source"`
}

func (m *Message) getAllAttributesAt(t time.Time) []Attribute {
	attrs := make([]Attribute, 0, len(m.Attributes))
	for _, attr := range m.Attributes {
		if !attr.ValidAt(t) {
			continue
		}
		attrs = append(attrs, attr)
	}
	return attrs
}

func (m *Message) GetAttributeAt(name string, t time.Time) (string, error) {
	for _, attr := range m.getAllAttributesAt(t) {
		if attr.Name == name {
			return attr.Value, nil
		}
	}
	return "", models.ErrNotFound
}

func (m *Message) GetAttributesAt(name string, t time.Time) []string {
	values := make([]string, 0)
	for _, attr := range m.getAllAttributesAt(t) {
		if attr.Name == name {
			values = append(values, attr.Value)
		}
	}
	return values
}

func (attr *Attribute) ValidAt(t time.Time) bool {
	if !attr.StartDate.Before(t) {
		return false
	}
	return attr.EndDate == nil || attr.EndDate.After(t)
}
