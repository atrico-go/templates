package collections

import (
	"fmt"

	"github.com/atrico-go/templates"
	"github.com/iancoleman/strcase"
)

type QueueTemplateDetails struct {
	// Type of queue
	QueueType
	// Type of queue element
	ElementType string
	// Type name (defaults to ElementTypeQueue)
	TypeName string
}

type QueueType int

const (
	QueueTypeNormal      QueueType = iota
	QueueTypeThreadSafe  QueueType = iota // Mutex controlled to allow multiple threads to use queue safely
	QueueTypeMultiThread QueueType = iota // Designed for multi-thread queue processing
)

func CreateQueue(fileDetails templates.FileDetails, details QueueTemplateDetails) (err error) {
	d := queueTemplateDetails{
		MultiThread: details.QueueType == QueueTypeMultiThread,
		AccessMutex: details.QueueType == QueueTypeMultiThread || details.QueueType == QueueTypeThreadSafe,
		ElementType: details.ElementType,
		TypeName:    details.TypeName,
	}
	d.setDefaults()
	// Defaults
	return templates.CreateTemplate(queueTemplate, fileDetails, d)
}

type queueTemplateDetails struct {
	MultiThread bool
	AccessMutex bool
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
{{- if .MultiThread}}
	"context"
{{- end}}
	"sort"
{{- if .AccessMutex}}
	"sync"
{{- end}}
{{- if .MultiThread}}

	"github.com/atrico-go/core/syncEx" // >= v1.6.0
{{- end}}
)

type {{.TypeName}} struct {
	queue {{if .AccessMutex}}      {{end}}{{if .MultiThread}}   {{end}}[]{{.ElementType}}
{{- if .AccessMutex}}
	accessMutex {{if .MultiThread}}   {{end}}sync.Mutex
{{- end}}
{{- if .MultiThread}}
	availableEvent syncEx.Event
	emptyEvent     syncEx.Event
{{- end}}
}

func Make{{.TypeName}}(initial ...{{.ElementType}}) {{.TypeName}} {
	q := {{.TypeName}}{}
	for _, el := range initial {
		q.Push(el)
	}
	return q
}

func (q *{{.TypeName}}) IsEmpty() bool {
	return q.Count() == 0
}

func (q *{{.TypeName}}) Count() int {
{{- if .AccessMutex}}
	q.accessMutex.Lock()
	defer q.accessMutex.Unlock()
{{- end}}
	return q.count()
}

func (q *{{.TypeName}}) Push(el {{.ElementType}}) {
{{- if .AccessMutex}}
	q.accessMutex.Lock()
	defer q.accessMutex.Unlock()
{{- end}}
	q.queue = append(q.queue, el)
{{- if .MultiThread}}
	q.availableEvent.Set()
	q.emptyEvent.Reset()
{{- end}}
}

func (q *{{.TypeName}}) Pop({{if .MultiThread}}ctx context.Context{{end}}) (element {{.ElementType}}{{if .MultiThread}}, err error{{end}}) {
	return q.getFirstElement(true{{if .MultiThread}}, ctx{{end}})
}

func (q *{{.TypeName}}) Peek({{if .MultiThread}}ctx context.Context{{end}}) (element {{.ElementType}}{{if .MultiThread}}, err error{{end}}) {
	return q.getFirstElement(false{{if .MultiThread}}, ctx{{end}})
}

// Sort Queue, lowest value will be popped next
func (q *{{.TypeName}}) Sort(before func(i, j {{.ElementType}}) bool) {
{{- if .AccessMutex}}
	q.accessMutex.Lock()
	defer q.accessMutex.Unlock()
{{- end}}
	sort.Slice(q.queue, func(i, j int) bool {return before(q.queue[i], q.queue[j])})
}
{{- if .MultiThread}}

func (q *{{.TypeName}}) WaitUntilEmpty(ctx context.Context) error {
	return q.emptyEvent.Wait(ctx)
}
{{- end}}

func (q *{{.TypeName}}) getFirstElement(remove bool{{if .MultiThread}}, ctx context.Context{{end}}) (element {{.ElementType}}{{if .MultiThread}}, err error{{end}}) {
{{- if .AccessMutex}}
	q.accessMutex.Lock()
	defer q.accessMutex.Unlock()
{{- end}}
{{- if .MultiThread}}
	err = q.availableEvent.Wait(ctx)
	if err == nil {
{{- else}}
	if q.count() == 0 {
		panic("queue is empty")
	}
{{- end}}
	{{if .MultiThread}}	{{end}}element = q.queue[0]
	{{if .MultiThread}}	{{end}}if remove {
		{{if .MultiThread}}	{{end}}q.queue = q.queue[1:]
{{- if .MultiThread}}
			q.availableEvent.SetValue(q.count() > 0)
			q.emptyEvent.SetValue(q.count() == 0)
		}
{{- end}}
	}
	return element{{if .MultiThread}}, err{{end}}
}

func (q *{{.TypeName}}) count() int {
	return len(q.queue)
}
`
