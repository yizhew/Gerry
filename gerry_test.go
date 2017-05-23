package gerry

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type MyType struct {
	Name string
}

func (m *MyType) CallMe() string {
	return fmt.Sprintf("me.%s", m.Name)
}

type OtherType struct {
	MyType
}

func (o *OtherType) CallU() string {
	return fmt.Sprintf("u.%s", o.Name)
}

type MyString string

func (m *MyString) Parts() []string {
	s := string(*m)
	return strings.Split(s, ".")
}

func TestTypeDefine(t *testing.T) {
	s1 := MyString("nodot")
	s2 := MyString("no.dot")
	s3 := MyString("no.dot.")

	assert.Equal(t, 1, len(s1.Parts()), "should follow parts split")
	assert.Equal(t, 2, len(s2.Parts()), "should follow parts split")
	assert.Equal(t, 3, len(s3.Parts()), "should follow parts split")

	// assert.True(t, , ...)
}

func TestMyType(t *testing.T) {
	mt := &MyType{"mt"}

	assert.Equal(t, "me.mt", mt.CallMe(), "call me")

	ot := OtherType{*mt}
	ot.Name = "mt"
	assert.Equal(t, "u.mt", (&ot).CallU(), "call u")
	assert.Equal(t, "me.mt", (&ot).CallMe(), "call me")
	assert.Equal(t, "mt", (&ot).Name, "name")
}
