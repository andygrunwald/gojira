// +build generated

package api

import (
	"github.com/ghostsquad/genny/generic"
	"github.com/valyala/fastjson"
)

type GenericType generic.Type

func UnmarshalGenericTypeFromValue(v *fastjson.Value) (GenericType, error) {
	if v == nil {
		return GenericType{}, nil
	}

	return UnmarshalGenericType(v.String())
}
