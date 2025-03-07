package reago

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type XMLNode struct {
	XMLName xml.Name   `xml:""`
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",chardata"`
	Nodes   []XMLNode  `xml:",any"`
}

func (node *XMLNode) GetTag() string {
	return strings.ToLower(node.XMLName.Local)
}

func (node *XMLNode) HasAttr(name string) bool {
	for _, attr := range node.Attrs {
		if attr.Name.Local == name {
			return true
		}
	}
	return false
}

func (node *XMLNode) GetAttr(name string) string {
	for _, attr := range node.Attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

func (node *XMLNode) GetAttrInt(name string) int {
	value := node.GetAttr(name)
	if value != "" {
		if result, err := strconv.Atoi(value); err == nil {
			return result
		}
	}
	return 0
}

func (node *XMLNode) GetAttrInt32(name string) int32 {
	return int32(node.GetAttrInt(name))
}

func (node *XMLNode) GetAttrFloat(name string) float64 {
	value := node.GetAttr(name)
	if value != "" {
		if result, err := strconv.ParseFloat(value, 64); err == nil {
			return result
		}
	}
	return 0
}

func (node *XMLNode) GetAttrFloat32(name string) float32 {
	return float32(node.GetAttrFloat(name))
}

func (node *XMLNode) GetAttrBool(name string) bool {
	value := node.GetAttr(name)
	return value == "true" || value == "1"
}

func (node *XMLNode) GetText() string {
	return strings.TrimSpace(node.Content)
}

func (node *XMLNode) HasBind(name string) bool {
	for _, attr := range node.Attrs {
		if attr.Name.Space == "bind" && attr.Name.Local == name {
			return true
		}
	}
	return false
}

func (node *XMLNode) GetBind(name string) string {
	for _, attr := range node.Attrs {
		if attr.Name.Space == "bind" && attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}
