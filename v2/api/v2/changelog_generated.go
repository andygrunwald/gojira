// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/ghostsquad/genny

package v2

import "github.com/valyala/fastjson"

func (x *Changelog) UnmarshalFromObj(o *fastjson.Object) error {
	if o == nil {
		return nil
	}

	return x.Unmarshal(o.String())
}
