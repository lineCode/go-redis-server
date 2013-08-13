package redis

import (
	"strconv"
)

type Request struct {
	name string
	args [][]byte
}

func (r *Request) HasArgument(index int) bool {
	return len(r.args) >= index+1
}

func (r *Request) ExpectArgument(index int) ReplyWriter {
	if !r.HasArgument(index) {
		return NewError("Not enough arguments")
	}
	return nil
}

func (r *Request) GetString(index int) (string, ReplyWriter) {
	if reply := r.ExpectArgument(index); reply != nil {
		return "", reply
	}
	return string(r.args[index]), nil
}

func (r *Request) GetInteger(index int) (int, ReplyWriter) {
	if reply := r.ExpectArgument(index); reply != nil {
		return -1, reply
	}
	i, err := strconv.Atoi(string(r.args[index]))
	if err != nil {
		return -1, NewError("Expected integer")
	}
	return i, nil
}

func (r *Request) GetPositiveInteger(index int) (int, ReplyWriter) {
	i, reply := r.GetInteger(index)
	if reply != nil {
		return -1, reply
	}
	if i < 0 {
		return -1, NewError("Expected positive integer")
	}
	return i, nil
}

func (r *Request) GetMap(index int) (map[string][]byte, ReplyWriter) {
	count := len(r.args) - index
	if count <= 0 {
		return nil, NewError("Expected at least one key val pair")
	}
	if count%2 != 0 {
		return nil, NewError("Got uneven number of key val pairs")
	}
	values := make(map[string][]byte)
	for i := index; i < len(r.args); i += 2 {
		key, reply := r.GetString(i)
		if reply != nil {
			return nil, reply
		}
		values[key] = r.args[i+1]
	}
	return values, nil
}
