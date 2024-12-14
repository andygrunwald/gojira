package v2

import (
	"fmt"
	"github.com/valyala/fastjson"
)

var issueJsonParserPool fastjson.ParserPool

// Issue represents a Jira issue.
type Issue struct {
	Expand         string               `json:"expand"`
	ID             string               `json:"id"`
	Self           string               `json:"self"`
	Key            string               `json:"key"`
	Fields         *IssueFields         `json:"field"`
	RenderedFields *IssueRenderedFields `json:"renderedFields"`
	Changelog      *Changelog           `json:"changelog"`
	Transitions    []Transition         `json:"transitions"`
	Names          map[string]string    `json:"names"`
}

// Unmarshal unmarshals JSON s to issue.
func (x *Issue) Unmarshal(j string) error {

	p := issueJsonParserPool.Get()
	defer issueJsonParserPool.Put(p)
	v, err := p.Parse(j)

	if err != nil {
		return fmt.Errorf("failed to parse issue json: %w", err)
	}

	x.Expand = string(v.GetStringBytes("expand"))
	x.ID = string(v.GetStringBytes("id"))
	x.Self = string(v.GetStringBytes("self"))
	x.Key = string(v.GetStringBytes("key"))

	{
		y := &IssueFields{}
		err := y.UnmarshalFromObj(v.GetObject("fields"))
		if err != nil {
			return err
		}
	}

	{
		y := &IssueRenderedFields{}
		err = y.UnmarshalFromObj(v.GetObject("renderedFields"))
		if err != nil {
			return err
		}
		x.RenderedFields = y
	}

	{
		y := &Changelog{}
		err = y.UnmarshalFromObj(v.GetObject("changelog"))
		if err != nil {
			return err
		}
		x.Changelog = y
	}

	{
		for _, jV := range v.GetArray("transitions") {
			y, err := UnmarshalTransitionFromValue(jV)
			if err != nil {
				return err
			}
			x.Transitions = append(x.Transitions, y)
		}
	}

	{
		y := make(map[string]string)
		n := v.GetObject("names")
		if n != nil {
			n.Visit(func(key []byte, v *fastjson.Value) {
				y[string(key)] = string(v.GetStringBytes())
			})
		}
	}

	return nil
}
