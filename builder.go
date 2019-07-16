package resolver

import (
	"fmt"
	"strings"
)

const (
	String typesBuilder = 1
	Int    typesBuilder = 2
)

type (
	typesBuilder uint8

	Request interface {
		Extract() string
	}

	Builder struct {
		rawString string
		props     []prop
	}

	prop struct {
		key   string
		types typesBuilder
		value interface{}
	}
)

func (b *Builder) ComposeStr(key, value string) *Builder {
	b.props = append(b.props, prop{
		key:   key,
		types: String,
		value: value,
	})
	return b
}

func (b *Builder) ComposeInt(key string, value int) *Builder {
	b.props = append(b.props, prop{
		key:   key,
		types: Int,
		value: value,
	})

	return b
}

func (b *Builder) Extract() string {
	if b.rawString != "" {
		return b.rawString
	}

	var str []string
	for _, prop := range b.props {
		switch prop.types {
		case String:
			str = append(str, fmt.Sprintf(`"%s":"%v"`, prop.key, prop.value))
		case Int:
			str = append(str, fmt.Sprintf(`"%s":%v`, prop.key, prop.value))
		default:
			str = append(str, fmt.Sprintf(`"%s":%v`, prop.key, prop.value))
		}
	}

	return fmt.Sprintf("{%s}", strings.Join(str, ","))
}

func Adapter(req string) Request {
	b := Builder{
		rawString: req,
	}
	return &b
}
