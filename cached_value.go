package templates

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

type CachedValueTemplateDetails struct {
	// Has no reset function
	NoReset bool
	// Type of queue element
	ElementType string
	// Type name (defaults to ElementTypeCachedValue)
	TypeName string
}

func CreateCachedValue(fileDetails FileDetails, details CachedValueTemplateDetails) (err error) {
	details.setDefaults()
	// Defaults
	return CreateTemplate(queueTemplate, fileDetails, details)
}

func (d *CachedValueTemplateDetails) setDefaults() {
	if d.TypeName == "" {
		d.TypeName = fmt.Sprintf("Cached%sValue", strcase.ToCamel(d.ElementType))
	}
}

var queueTemplate = `package {{.Package}}

import (
)

type {{.TypeName}} struct {
	Creator  func() {{.ElementType}}
	theValue {{.ElementType}}
	hasValue bool
}

func (v *{{.TypeName}}) GetValue() {{.ElementType}} {
	if !v.hasValue {
		v.theValue = v.Creator()
		v.hasValue = true
	}
	return v.theValue
}
{{- if not .NoReset}}

func (v *{{.TypeName}}) Reset() {
	v.hasValue = false
}
{{- end}}
`
