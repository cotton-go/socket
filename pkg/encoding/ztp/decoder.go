package ztp

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"sync"
)

type Decoder struct {
	mutex sync.Mutex   // each item must be received atomically
	r     io.Reader    // source of the data
	buf   bytes.Buffer // buffer for more efficient i/o from r
	err   error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (dec *Decoder) Decode(e any) error {
	if e == nil {
		return dec.DecodeValue(reflect.Value{})
	}

	value := reflect.ValueOf(e)
	// If e represents a value as opposed to a pointer, the answer won't
	// get back to the caller. Make sure it's a pointer.
	if value.Type().Kind() != reflect.Pointer {
		dec.err = errors.New("gob: attempt to decode into a non-pointer")
		return dec.err
	}

	return dec.DecodeValue(value)
}

func (dec *Decoder) DecodeValue(v reflect.Value) error {
	if v.IsValid() {
		if v.Kind() == reflect.Pointer && !v.IsNil() {
			// That's okay, we'll store through the pointer.
		} else if !v.CanSet() {
			return errors.New("gob: DecodeValue of unassignable value")
		}
	}

	dec.buf.Reset() // In case data lingers from previous invocation.
	dec.err = nil
	id := dec.decodeTypeSequence(false)
	if dec.err == nil {
		dec.decodeValue(id, v)
	}
	return dec.err
}

func (dec *Decoder) decodeValue(_ any, v reflect.Value) {

}
