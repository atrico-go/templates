package collections

import (
	"fmt"

	"github.com/atrico-go/templates"
	"github.com/iancoleman/strcase"
)

type QueueTemplateDetails struct {
	// Type of queue element
	ElementType string
	// Type name (defaults to ElementTypeQueue)
	TypeName string
}

func CreateQueue(fileDetails templates.FileDetails, details QueueTemplateDetails) (err error) {
	d := queueTemplateDetails{
		ElementType: details.ElementType,
		TypeName:    details.TypeName,
	}
	d.setDefaults()
	// Defaults
	return templates.CreateTemplate(queueTemplate, fileDetails, d)
}

func CreateMultiThreadQueue(fileDetails templates.FileDetails, details QueueTemplateDetails) (err error) {
	d := queueTemplateDetails{
		MultiThread: true,
		ElementType: details.ElementType,
		TypeName:    details.TypeName,
	}
	d.setDefaults()
	// Defaults
	return templates.CreateTemplate(queueTemplate, fileDetails, d)
}

type queueTemplateDetails struct {
	MultiThread bool
	ElementType string
	TypeName    string
}

func (d *queueTemplateDetails) setDefaults() {
	if d.TypeName == "" {
		d.TypeName = fmt.Sprintf("%sQueue", strcase.ToCamel(d.ElementType))
	}
}

var queueTemplate = `package {{.Package}}

import (
	"sort"
{{- if .MultiThread}}
	"sync"
	"time"

	"github.com/atrico-go/core/syncEx" // >= v1.6.0
{{- end}}
)

{{- $lowTypeName := ToLowerCamel .TypeName}}

type {{.TypeName}} interface {
	IsEmpty() bool
	Count() int
	Push(el {{.ElementType}})
	Pop({{if .MultiThread}}timeout time.Duration{{end}}) (element {{.ElementType}}{{if .MultiThread}}, ok bool{{end}})
	Peek({{if .MultiThread}}timeout time.Duration{{end}}) (element {{.ElementType}}{{if .MultiThread}}, ok bool{{end}})
	Sort(less func(i, j {{.ElementType}}) bool)
{{- if .MultiThread}}
	WaitUntilEmpty(timeout time.Duration) bool
{{- end}}
}

func New{{.TypeName}}(initial ...{{.ElementType}}) {{.TypeName}} {
	count := len(initial)
	q := {{$lowTypeName}}{
		queue: make([]{{.ElementType}}, count),
		count: count,
	}
	for i, el := range initial {
		q.queue[i] = el
	}
{{- if .MultiThread}}
	q.availableEvent.SetValue(count > 0)
	q.emptyEvent.SetValue(count == 0)
{{- end}}
	return &q
}

// ----------------------------------------------------------------------------------------------------------------------------
// Implementation
// ----------------------------------------------------------------------------------------------------------------------------

type {{$lowTypeName}} struct {
	queue {{if .MultiThread}}         {{end}}[]{{.ElementType}}
	count {{if .MultiThread}}         {{end}}int
{{- if .MultiThread}}
	accessMutex    sync.Mutex
	availableEvent syncEx.Event
	emptyEvent     syncEx.Event
{{- end}}
}

func (q *{{$lowTypeName}}) IsEmpty() bool {
	return q.Count() == 0
}

func (q *{{$lowTypeName}}) Count() int {
	return q.count
}

func (q *{{$lowTypeName}}) Push(el {{.ElementType}}) {
{{- if .MultiThread}}
	q.accessMutex.Lock()
	defer q.accessMutex.Unlock()
{{- end}}
	q.queue = append(q.queue, el)
	q.count++
{{- if .MultiThread}}
	q.availableEvent.Set()
	q.emptyEvent.Reset()
{{- end}}
}

func (q *{{$lowTypeName}}) Pop({{if .MultiThread}}timeout time.Duration{{end}}) (element {{.ElementType}}{{if .MultiThread}}, ok bool{{end}}) {
{{- if .MultiThread}}
	q.accessMutex.Lock()
	defer q.accessMutex.Unlock()
	element, ok = q.peek(timeout)
	if ok {
		q.queue = q.queue[1:]
		q.count--
		q.availableEvent.SetValue(!q.IsEmpty())
		q.emptyEvent.SetValue(q.IsEmpty())
	}
{{- else}}
	element = q.Peek()
	q.queue = q.queue[1:]
	q.count--
{{- end}}
	return element{{if .MultiThread}}, ok{{end}}
}

func (q *{{$lowTypeName}}) Peek({{if .MultiThread}}timeout time.Duration{{end}}) (element {{.ElementType}}{{if .MultiThread}}, ok bool{{end}}) {
{{- if .MultiThread}}
	q.accessMutex.Lock()
	defer q.accessMutex.Unlock()
	return q.peek(timeout)
{{- else}}
	if q.IsEmpty() {
		panic("queue is empty")
	}
	return q.queue[0]
{{- end}}
}

// Sort Queue, lowest value will be popped next
func (q *{{$lowTypeName}}) Sort(before func(i, j {{.ElementType}}) bool) {
{{- if .MultiThread}}
	q.accessMutex.Lock()
	defer q.accessMutex.Unlock()
{{- end}}
	sort.Slice(q.queue, func(i, j int) bool {return before(q.queue[i], q.queue[j])})
}
{{- if .MultiThread}}

func (q *{{$lowTypeName}}) WaitUntilEmpty(timeout time.Duration) bool {
	return q.emptyEvent.Wait(timeout)
}

func (q *{{$lowTypeName}}) peek(timeout time.Duration) (element {{.ElementType}}, ok bool) {
	ok = q.availableEvent.Wait(timeout)
	if ok {
		element = q.queue[0]
	}
	return element, ok
}
{{- end}}
`
