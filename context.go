package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	// "Context.Req.URL.Query()" will execute parse action everytime
	// So cache it here for repeat usage
	parsedQuery url.Values
}

func (c *Context) BindJSON(val any) error {
	// "c.Req.Body" is an interface "io.ReadCloser", so it can only be read once
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) StringValue {
	// "c.Req.ParseFrom" already have builtin logic to ensure that the parse action
	// will only be executed once. So no need to control it here.
	err := c.Req.ParseForm()
	if err != nil {
		return StringValue{err: err}
	}
	return StringValue{str: c.Req.FormValue(key)}
}

func (c *Context) QueryValue(key string) StringValue {
	// "Context.Req.URL.Query()" will execute parse action everytime
	// So only execute it when its cache "Context.parsedQuery" is nil
	if c.parsedQuery == nil {
		c.parsedQuery = c.Req.URL.Query()
	}

	if vs := c.parsedQuery[key]; len(vs) > 1 {
		return StringValue{str: vs[0]}
	}
	return StringValue{str: "", err: errors.New(key + ": key not found")}
}

func (c *Context) PathParamValue(key string) StringValue {
	v, ok := c.PathParams[key]
	if !ok {
		return StringValue{str: "", err: errors.New(key + ": key not found")}
	}
	return StringValue{str: v}
}

func (c *Context) RespJSONOK(v any) error {
	return c.RespJSON(http.StatusOK, v)
}

func (c *Context) RespJSON(code int, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.Resp.WriteHeader(code)

	// normally, no need to check write result
	_, err = c.Resp.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// StringValue
// For convenient convert of return value
// for example:
//
//	int64Val, err := context.BindJSON({"64", nil}).AsInt64()
type StringValue struct {
	str string
	err error
}

func (sv StringValue) AsInt32() (int32, error) {
	if sv.err != nil {
		return 0, sv.err
	}
	i64, err := strconv.ParseInt(sv.str, 10, 64)
	if err != nil {
		return 0, sv.err
	}
	return int32(i64), nil
}

func (sv StringValue) AsInt64() (int64, error) {
	if sv.err != nil {
		return 0, sv.err
	}
	return strconv.ParseInt(sv.str, 10, 64)
}
