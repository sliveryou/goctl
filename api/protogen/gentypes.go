package protogen

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/logrusorgru/aurora"
	"github.com/zeromicro/go-zero/core/stringx"

	"github.com/sliveryou/goctl/api/spec"
	"github.com/sliveryou/goctl/util"
)

const (
	indent = "  "
)

// BuildTypes gen types to string
func BuildTypes(api *spec.ApiSpec) (string, error) {
	var builder strings.Builder

	for i, tp := range api.Types {
		if i > 0 {
			builder.WriteString("\n\n")
		}
		if err := writeMessage(&builder, tp); err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}

func writeMessage(b *strings.Builder, t spec.Type) error {
	st, ok := t.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unspport struct type: %s", t.Name())
	}

	mfs := parseMessageFields(st)
	b.WriteString(getDoc(st) + "message " + util.Title(st.RawName) + " {\n")

	for i, mf := range mfs {
		repeated := ""
		if mf.IsRepeated {
			repeated = "repeated "
		}
		comment := ""
		if len(mf.Comment) > 0 {
			comment = " // " + mf.Comment
		}
		b.WriteString(fmt.Sprintf("%s%s%s %s = %d;%s\n",
			indent, repeated, mf.FieldType, mf.FieldName, i+1, comment))
	}

	b.WriteByte('}')

	return nil
}

type messageField struct {
	FieldName  string
	FieldType  string
	Comment    string
	IsRepeated bool
}

func parseMessageFields(ds spec.DefineStruct) []messageField {
	var mfs []messageField

	for _, m := range ds.Members {
		var tagName string
		if tag := getUsefulTag(m.Tag); tag != nil {
			tagName = tag.Name
		}
		if tagName == "" {
			if m.Name == "" {
				tagName = strcase.ToSnake(trimPrefix(m.Type.Name()))
			} else {
				tagName = strcase.ToSnake(trimPrefix(m.Name))
			}
		}

		switch mt := convertSpec(m.Type).(type) {
		case spec.PrimitiveType:
			mfs = append(mfs, messageField{
				FieldName:  tagName,
				FieldType:  getFieldType(mt.RawName),
				Comment:    getComment(m.Comment),
				IsRepeated: false,
			})
		case spec.ArrayType:
			mf := messageField{
				FieldName:  tagName,
				FieldType:  trimPrefix(mt.RawName),
				Comment:    getComment(m.Comment),
				IsRepeated: true,
			}
			if mf.FieldType == "byte" {
				mf.FieldType = "bytes"
				mf.IsRepeated = false
			}
			mfs = append(mfs, mf)
		case spec.DefineStruct:
			if m.IsInline && len(ds.Members) > 1 {
				mfs = append(mfs, parseMessageFields(mt)...)
			} else {
				mfs = append(mfs, messageField{
					FieldName:  tagName,
					FieldType:  getFieldType(mt.RawName),
					Comment:    getComment(m.Comment),
					IsRepeated: false,
				})
			}
		default:
			fmt.Println(aurora.Red(fmt.Sprintf("struct type: %s, member name: %s, member type: %s, "+
				"convert message failed.", ds.RawName, m.Name, mt.Name())))
		}
	}

	return mfs
}

func getUsefulTag(tag string) *spec.Tag {
	tags, err := spec.Parse(tag)
	if err != nil {
		return nil
	}

	keys := []string{"json", "form", "path", "header"}
	for _, t := range tags.Tags() {
		if stringx.Contains(keys, t.Key) {
			return t
		}
	}

	return nil
}

func getFieldType(dataType string) string {
	switch dataType {
	case "int", "int64":
		return "int64"
	case "int8", "int16", "int32":
		return "int32"
	case "byte", "uint8", "uint16", "uint32":
		return "uint32"
	case "uint", "uint64":
		return "uint64"
	case "float32", "complex64":
		return "float"
	case "float64", "complex128":
		return "double"
	case "bool":
		return "bool"
	case "string":
		return "string"
	default:
		return dataType
	}
}

func convertSpec(t spec.Type) spec.Type {
	var tt spec.PointerType
	var ok = true

	for ok {
		tt, ok = t.(spec.PointerType)
		if ok {
			t = tt.Type
		}
	}

	return t
}

func getDoc(st spec.DefineStruct) string {
	if docs := st.Docs; len(docs) > 0 {
		return docs[len(docs)-1] + "\n"
	} else {
		return "// " + st.RawName + " 详情信息\n"
	}
}

func trimPrefix(t string) string {
	ft := strings.TrimPrefix(t, "[]")
	return strings.TrimSpace(strings.TrimPrefix(ft, "*"))
}

func getComment(comment string) string {
	return strings.TrimSpace(strings.TrimPrefix(comment, "//"))
}
