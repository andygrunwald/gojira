// +build generated

package api

import (
	"github.com/ghostsquad/genny/generic"
	"github.com/valyala/fastjson"
)

type GenericType generic.Type

func (x *GenericType) UnmarshalFromObj(o *fastjson.Object) error {
	if o == nil {
		return nil
	}

	return x.Unmarshal(o.String())
}
