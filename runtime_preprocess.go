package parsex

import (
	"log"
	"reflect"

	"github.com/bbfh-dev/berr"
	"github.com/bbfh-dev/parsex/v2/internal"
	"github.com/iancoleman/strcase"
)

func (runtime *runtimeType) preprocess() error {
	if runtime.data == nil {
		return nil
	}

	typePtr := reflect.TypeOf(runtime.data)
	if typePtr.Kind() != reflect.Pointer {
		return berr.WithContext(
			runtime.name,
			ErrMustbePointer,
			"Provided Program.Data", typePtr,
		)
	}

	typeElem := typePtr.Elem()
	if typeElem.Kind() != reflect.Struct {
		return berr.WithContext(
			runtime.name,
			ErrMustPointToStruct,
			"Provided *Program.Data", typePtr,
		)
	}

	valueElem := reflect.ValueOf(runtime.data).Elem()
	numOfFields := typeElem.NumField()

	runtime.genOptions.Clear()
	runtime.genOptions.Add("help", internal.HelpOption)
	if runtime.version != "" {
		runtime.genOptions.Add("version", internal.VersionOption)
	}

	for i := range numOfFields {
		fieldType := typeElem.Field(i)
		fieldValue := valueElem.Field(i)
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			log.Printf("WARNING: (Parsex reflection) not possible to modify field %+v", fieldType)
			continue
		}

		name := strcase.ToKebab(fieldType.Name)
		option := internal.Option{
			Name:    name,
			Alt:     fieldType.Tag.Get("alt"),
			Desc:    fieldType.Tag.Get("desc"),
			Default: fieldType.Tag.Get("default"),
			Type:    fieldType.Type,
			Ref:     &fieldValue,
		}
		option.Set(option.Default)
		runtime.genOptions.Add(name, option)
		if alt := fieldType.Tag.Get("alt"); alt != "" {
			runtime.genOptionAlts[alt] = name
		}
	}

	return nil
}
