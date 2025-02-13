package command

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func NewCommandReader() *CommandReader {
	ctx := &CommandReader{
		cmd: &CommandFrags{},
	}

	return ctx
}

func (ctx *CommandReader) Read(p []byte) (int, error) {
	return ctx.b.Read(p)
}

func (ctx *CommandReader) Write(p []byte) (int, error) {
	return ctx.b.Write(p)
}

func (ctx *CommandReader) WriteVar(name string, val string) (int, error) {
	var (
		err error
		b   []byte
		v   reflect.Value = reflect.ValueOf(ctx.cmd)
	)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return 0, fmt.Errorf("expected a pointer to a struct")
	}

	v = v.Elem()

	field := v.FieldByName(name)
	newValueReflect := reflect.ValueOf(val)
	if newValueReflect.Type() != field.Type() {
		return 0, fmt.Errorf("type mismatch: expected %s but got %s", field.Type(), newValueReflect.Type())
	}
	field.Set(newValueReflect)

	if b, err = json.Marshal(ctx.cmd); err != nil {
		return 0, err
	}
	ctx.b.Reset()
	return ctx.Write(b)
}

func (ctx *CommandReader) Close() error {
	return nil
}
