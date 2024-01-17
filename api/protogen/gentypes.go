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
	indent = "    "
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
		optional := ""
		if mf.IsRepeated {
			repeated = "repeated "
		} else if mf.IsPointer {
			optional = "optional "
		}
		comment := mf.Comment
		if mf.IsOptional {
			if comment == "" {
				comment = "非必填"
			} else {
				comment += "，非必填"
			}
		}
		if comment != "" {
			comment = " // " + comment
		}
		b.WriteString(fmt.Sprintf("%s%s%s%s %s = %d;%s\n",
			indent, repeated, optional, mf.FieldType, mf.FieldName, i+1, comment))
	}

	b.WriteByte('}')

	return nil
}

type messageField struct {
	FieldName      string
	FieldType      string
	Comment        string
	IsRepeated     bool
	IsPointer      bool
	IsOptional     bool
	MessageName    string
	MessageComment string
}

func parseMessageFields(ds spec.DefineStruct) []messageField {
	var mfs []messageField
	replacer := strings.NewReplacer("__", "_", ".", "_", " ", "_")

	for _, m := range ds.Members {
		var tagName string
		var isOptional bool
		if tag := getUsefulTag(m.Tag); tag != nil {
			tagName = replacer.Replace(tag.Name)
			isOptional = stringx.Contains(tag.Options, "optional")
			if tagSnakeName := strcase.ToSnake(tagName); tagSnakeName != tagName {
				tagName = tagSnakeName
			}
		}
		if tagName == "" {
			if m.Name == "" {
				tagName = strcase.ToSnake(trimPrefix(m.Type.Name()))
			} else {
				tagName = strcase.ToSnake(trimPrefix(m.Name))
			}
		}

		mt, isPointer := convertSpec(m.Type)
		switch mt := (mt).(type) {
		case spec.PrimitiveType:
			mfs = append(mfs, messageField{
				FieldName:  tagName,
				FieldType:  getFieldType(mt.RawName),
				Comment:    getComment(m.Comment),
				IsRepeated: false,
				IsPointer:  isPointer,
				IsOptional: isOptional,
			})
		case spec.ArrayType:
			mf := messageField{
				FieldName:  tagName,
				FieldType:  trimPrefix(mt.RawName),
				Comment:    getComment(m.Comment),
				IsRepeated: true,
				IsPointer:  isPointer,
				IsOptional: isOptional,
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
					IsPointer:  isPointer,
					IsOptional: isOptional,
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

func convertSpec(t spec.Type) (spec.Type, bool) {
	var (
		tt        spec.PointerType
		ok        = true
		isPointer = false
	)

	for ok {
		tt, ok = t.(spec.PointerType)
		if ok {
			t = tt.Type
			isPointer = true
		}
	}

	return t, isPointer
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
