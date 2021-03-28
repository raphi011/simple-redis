package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)


type RedisServer struct {
	in  io.Reader
	out io.Writer

	keys map[string]Value
}

type Value struct {
	val interface{}
	expiry *time.Time
}

func (p *RedisServer) handleRequest() {
	for {
		t, err := parse(p.in)

		if err != nil {
			if err != io.EOF {
				toString(err, p.out)
			}

			break
		}

		switch r := t.(type) {
		case Array:
			if err := p.handleCommand(r, p.out); err != nil {
				toString(err, p.out)
			}
		default:
			toString(fmt.Errorf("request type %T not implemented yet: %+v", r, r), p.out)
		}
	}
}

func (p *RedisServer) handleCommand(request []interface{}, out io.Writer) error {
	var cmd string

	switch v := request[0].(type) {
	case []byte:
		cmd = string(v)
	default:
		return errors.New("command is not a BulkString")
	}

	switch strings.ToUpper(cmd) {
	case "PING":
		toString("PONG", out)
	case "ECHO":
		toString(request[1], out)
	case "SET":
		p.setCommand(request, out)
	case "GET":
		p.getCommand(request, out)
	default:
		return fmt.Errorf("unknown command %s", cmd)
	}

	return nil
}

func (p *RedisServer) getCommand(request Array, out io.Writer) {
	key, err := request.getString(1)

	if err != nil {
		toString(errors.New("did not pass key"), out)
		return
	}

	val, found := p.keys[key]
	expired := val.expiry != nil && val.expiry.Before(time.Now())

	if found && expired  {
		delete(p.keys, key)
	}

	if found && !expired {
		toString(val.val, out)
	} else {
		toString(nil, out)
	}
}

func (p *RedisServer) setCommand(request Array, out io.Writer) {
	key, err1 := request.getString(1)
	val, err2 := request.getString(2)

	if err1 != nil || err2 != nil {
		toString(errors.New("invalid cmd params"), out)
		return
	}

	if option, _ := request.getString(3); strings.ToUpper(option) == "PX" {
		expiryString, err := request.getString(4)

		if err != nil {
			toString(errors.New("px arg is not a string"), out)
			return
		}

		expiryInt, err := strconv.Atoi(expiryString)

		if err != nil {
			toString(errors.New("cannot convert px arg to int"), out)
			return
		}

		expiry := time.Now().Add(time.Duration(expiryInt) * time.Millisecond)

		p.keys[key] = Value{val: val, expiry: &expiry }
	} else {
		p.keys[key] = Value{val: val }
	}

	toString("OK", out)
}
