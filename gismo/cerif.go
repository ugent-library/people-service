package gismo

import (
	"io"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
)

type CerifValue struct {
	Value     string
	StartDate time.Time
	EndDate   time.Time
}

// Parse parses xml byte stream, removes namespace, and return xmlquery node
func Parse(reader io.Reader) (*xmlquery.Node, error) {
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		return nil, err
	}
	RemoveNamespace(doc)
	return doc, nil
}

func RemoveNamespace(node *xmlquery.Node) {
	node.Prefix = ""
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		RemoveNamespace(child)
	}
}

func ClassID(root *xmlquery.Node, name string) string {
	classIDNode := xmlquery.FindOne(root, "//cfClass[contains(cfURI, '"+name+"')]/cfClassId")
	if classIDNode == nil {
		return ""
	}
	return strings.TrimSpace(classIDNode.InnerText())
}

// NodesByClassName returns nodes matching tag name and sub attribute matching class name (uri)
func NodesByClassName(root *xmlquery.Node, tagName, className string) []*xmlquery.Node {
	classID := ClassID(root, className)
	if classID == "" {
		return nil
	}
	return xmlquery.Find(root, "//"+tagName+"[contains(cfClassId, '"+classID+"')]")
}

// ValuesByClassName transforms nodes from NodesByClassName into array of CerifValue. If set to a non empty value, valueTag is used to populate the cerif value.
func ValuesByClassName(root *xmlquery.Node, tagName, className, valueTag string) []CerifValue {
	nodes := NodesByClassName(root, tagName, className)
	vals := make([]CerifValue, 0, len(nodes))
	for _, node := range nodes {
		val := CerifValue{}
		if valueTag != "" {
			if n := xmlquery.FindOne(node, valueTag); n != nil {
				val.Value = strings.TrimSpace(n.InnerText())
			} else {
				continue
			}
		}
		if n := xmlquery.FindOne(node, "cfStartDate"); n != nil {
			t, err := time.Parse(time.RFC3339, strings.TrimSpace(n.InnerText()))
			if err != nil {
				continue
			}
			val.StartDate = t
		} else {
			continue
		}
		if n := xmlquery.FindOne(node, "cfEndDate"); n != nil {
			t, err := time.Parse(time.RFC3339, strings.TrimSpace(n.InnerText()))
			if err != nil {
				continue
			}
			val.EndDate = t
		} else {
			continue
		}
		vals = append(vals, val)
	}
	return vals
}
