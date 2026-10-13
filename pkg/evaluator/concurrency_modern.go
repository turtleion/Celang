package evaluator

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"celang/pkg/uast"
)

// Ticket represents a concurrent task handle created by 'titip'.
type Ticket struct {
	mu     sync.Mutex
	done   chan struct{}
	result Object
	err    *Error
	isDone bool
}

func NewTicket() *Ticket {
	return &Ticket{
		done: make(chan struct{}),
	}
}

func (t *Ticket) Type() ObjectType { return TICKET_OBJ }

func (t *Ticket) Inspect() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.isDone {
		return "<Tiket: sedang berjalan>"
	}
	if t.err != nil {
		return fmt.Sprintf("<Tiket: gagal (%s)>", t.err.Message)
	}
	resStr := "null"
	if t.result != nil {
		resStr = t.result.Inspect()
	}
	return fmt.Sprintf("<Tiket: selesai (%s)>", resStr)
}

func (t *Ticket) Resolve(val Object) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.isDone {
		return
	}
	if ret, ok := val.(*ReturnValue); ok {
		val = ret.Value
	}
	t.result = val
	t.isDone = true
	close(t.done)
}

func (t *Ticket) Fail(err *Error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.isDone {
		return
	}
	t.err = err
	t.isDone = true
	close(t.done)
}

func (t *Ticket) Await(timeout time.Duration, meta *uast.NodeMetadata) Object {
	if timeout <= 0 {
		<-t.done
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.err != nil {
			return t.err
		}
		if t.result == nil {
			return NULL
		}
		return t.result
	}

	select {
	case <-t.done:
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.err != nil {
			return t.err
		}
		if t.result == nil {
			return NULL
		}
		return t.result
	case <-time.After(timeout):
		return newError(meta,
			fmt.Sprintf("Batas waktu jemput tiket habis (%v)", timeout),
			"Gunakan waktu tunggu yang lebih lama atau optimalkan proses yang dititipkan.")
	}
}

// Channel represents Celang's CSP channel (Saluran) for inter-goroutine communication.
type Channel struct {
	Ch       chan Object
	Capacity int
	Closed   bool
	ElemType string
	mu       sync.Mutex
}

func NewChannel(capacity int, elemType string) *Channel {
	if capacity < 0 {
		capacity = 0
	}
	return &Channel{
		Ch:       make(chan Object, capacity),
		Capacity: capacity,
		ElemType: elemType,
	}
}

func (c *Channel) Type() ObjectType { return CHANNEL_OBJ }

func (c *Channel) Inspect() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	status := "terbuka"
	if c.Closed {
		status = "tertutup"
	}
	typeStr := ""
	if c.ElemType != "" && c.ElemType != "Sembarang" && c.ElemType != "sembarang" {
		typeStr = fmt.Sprintf("<%s>", c.ElemType)
	}
	return fmt.Sprintf("<Saluran%s: kapasitas=%d, isi=%d, status=%s>", typeStr, c.Capacity, len(c.Ch), status)
}

func (c *Channel) Send(val Object, meta *uast.NodeMetadata) Object {
	c.mu.Lock()
	if c.Closed {
		c.mu.Unlock()
		return newError(meta,
			"Tidak dapat mengirim data ke saluran yang sudah ditutup",
			"Pastikan saluran masih terbuka sebelum memanggil kirim().")
	}
	elemType := c.ElemType
	c.mu.Unlock()

	// Validate element type constraint if specified
	if elemType != "" && elemType != "Sembarang" && elemType != "sembarang" && elemType != "Any" {
		if !checkValueMatchesType(val, elemType) {
			return newError(meta,
				fmt.Sprintf("Tipe data tidak cocok: Saluran mengharuskan '%s', tetapi menerima '%s'", elemType, objectTypeName(val)),
				fmt.Sprintf("Kirim nilai yang sesuai dengan tipe %s.", elemType))
		}
	}

	c.Ch <- val
	return val
}

func (c *Channel) Receive(timeout time.Duration) Object {
	if timeout <= 0 {
		val, ok := <-c.Ch
		if !ok {
			return NULL
		}
		return val
	}

	select {
	case val, ok := <-c.Ch:
		if !ok {
			return NULL
		}
		return val
	case <-time.After(timeout):
		return NULL
	}
}

func (c *Channel) Close() Object {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.Closed {
		c.Closed = true
		close(c.Ch)
	}
	return NULL
}

// Type checking and constraint helpers
func objectTypeName(obj Object) string {
	if obj == nil {
		return "null"
	}
	switch o := obj.(type) {
	case *Integer:
		return "Angka"
	case *Float:
		return "Desimal"
	case *String:
		return "Teks"
	case *Boolean:
		return "Kebenaran"
	case *Array:
		return "Daftar"
	case *Map:
		return "Kamus"
	case *Channel:
		return "Saluran"
	case *Ticket:
		return "Tiket"
	case *Function:
		return "Fungsi"
	case *StructInstance:
		return o.Def.Name
	default:
		return string(obj.Type())
	}
}

func checkValueMatchesType(val Object, expectedType string) bool {
	if expectedType == "" || expectedType == "Sembarang" || expectedType == "sembarang" || expectedType == "Any" || expectedType == "any" {
		return true
	}
	if val == nil || val == NULL {
		return false
	}
	switch strings.ToLower(expectedType) {
	case "angka", "int", "integer":
		_, _, _, ok := CoerceToNumeric(val)
		return ok
	case "desimal", "float":
		_, ok := val.(*Float)
		return ok
	case "teks", "string":
		_, ok := val.(*String)
		return ok
	case "kebenaran", "bool", "boolean":
		_, ok := val.(*Boolean)
		return ok
	case "daftar", "larik", "array", "list":
		_, ok := val.(*Array)
		return ok
	case "kamus", "peta", "map":
		_, ok := val.(*Map)
		return ok
	case "saluran", "channel":
		_, ok := val.(*Channel)
		return ok
	case "tiket", "ticket":
		_, ok := val.(*Ticket)
		return ok
	case "fungsi", "function":
		_, ok := val.(*Function)
		return ok
	default:
		if inst, ok := val.(*StructInstance); ok {
			return strings.EqualFold(inst.Def.Name, expectedType)
		}
		return false
	}
}

func isTypeCompatible(argType, constraint string) bool {
	if constraint == "" || constraint == "Sembarang" || constraint == "sembarang" || constraint == "Any" || constraint == "any" {
		return true
	}
	if strings.EqualFold(argType, constraint) {
		return true
	}
	lowerCons := strings.ToLower(constraint)
	lowerArg := strings.ToLower(argType)
	if lowerCons == "angka" && (lowerArg == "int" || lowerArg == "float" || lowerArg == "desimal" || lowerArg == "integer") {
		return true
	}
	return false
}

func parseTimeoutValue(tVal Object) time.Duration {
	if tVal == nil || tVal == NULL {
		return 0
	}
	if iVal, fVal, isInt, ok := CoerceToNumeric(tVal); ok {
		if isInt {
			if iVal <= 60 {
				return time.Duration(iVal) * time.Second
			}
			return time.Duration(iVal) * time.Millisecond
		}
		return time.Duration(fVal * float64(time.Second))
	}
	return 0
}
