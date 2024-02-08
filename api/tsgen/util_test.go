package tsgen

import (
	"testing"

	"github.com/sliveryou/goctl/api/spec"
	"github.com/stretchr/testify/assert"
)

func TestGenTsType(t *testing.T) {
	member := spec.Member{
		Name:     "foo",
		Type:     spec.PrimitiveType{RawName: "string"},
		Tag:      `json:"foo,options=foo|bar|options|123"`,
		Comment:  "",
		Docs:     nil,
		IsInline: false,
	}
	ty, err := genTsType(member)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `'foo' | 'bar' | 'options' | '123'`, ty)

	member.IsInline = true
	ty, err = genTsType(member)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `'foo' | 'bar' | 'options' | '123'`, ty)

	member.Type = spec.PrimitiveType{RawName: "int"}
	member.Tag = `json:"foo,options=1|3|4|123"`
	ty, err = genTsType(member)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `1 | 3 | 4 | 123`, ty)
}
