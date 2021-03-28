package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Array []interface{}

func (a Array) getString(pos int) (string, error) {
	if pos >= len(a) || pos < 0 {
		return "", errors.New("index out of bounds")
	}

	switch v := a[pos].(type) {
	case string: return v, nil
	case []byte: return string(v), nil
	default: return "", errors.New("not a string")
	}
}

func parse(input io.Reader) (interface{}, error) {
	reader := bufio.NewReader(input)

	c, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	switch c {
	case '+':
		v, err := reader.ReadString('\n')

		if err != nil {
			return nil, fmt.Errorf("not a valid simple string")
		}

		return strings.Trim(v, "\r\n"), nil
	case '-':
		v, err := reader.ReadString('\n')

		if err != nil {
			return nil, fmt.Errorf("not a valid error")
		}

		return errors.New(strings.Trim(v, "\r\n")), nil
	case ':':
		v, err := reader.ReadString('\n')

		if err != nil {
			return nil, fmt.Errorf("not a valid integer")
		}

		i, err := strconv.Atoi(strings.Trim(v, "\r\n"))

		if err != nil {
			return nil, fmt.Errorf("not a valid integer %f", err)
		}

		return i, nil
	case '*':
		v, err := reader.ReadString('\n')

		if err != nil {
			return nil, fmt.Errorf("array length is not a valid integer")
		}

		size, err := strconv.Atoi(strings.Trim(v, "\r\n"))

		if err != nil {
			return nil, errors.New("array length is not a valid integer")
		}

		arr := make(Array, size)

		for i := 0; i < size; i++ {
			arr[i], err = parse(reader)

			if err != nil {
				return nil, fmt.Errorf("could not parse element in array at position %d: %f", i, err)
			}
		}

		return arr, nil
	case '$':
		v, err := reader.ReadString('\n')

		if err != nil {
			return nil, fmt.Errorf("string length is not a valid integer")
		}

		size, err := strconv.Atoi(strings.Trim(v, "\r\n"))

		if err != nil {
			return nil, errors.New("string length is not a valid integer")
		}

		if size == -1 {
			return nil, nil
		}

		buf := make([]byte, size)

		reader.Read(buf)

		// consume \r\n
		reader.ReadByte()
		reader.ReadByte()

		return buf, nil
	default:
		l, _, err := reader.ReadLine()

		if err != nil {
			return nil, err
		}

		return strings.Split(string(l), " "), nil
	}
}

func toString(input interface{}, writer io.Writer) error {
	switch v := input.(type) {
	case nil:
		writer.Write([]byte("$-1\r\n"))
	case string:
		writer.Write([]byte(fmt.Sprintf("+%s\r\n", v)))
	case []byte:
		writer.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)))
	case int:
		writer.Write([]byte(fmt.Sprintf(":%d\r\n", v)))
	case Array:
		writer.Write([]byte(fmt.Sprintf("*%d\r\n", len(v))))

		for _, elem := range v {
			if err := toString(elem, writer); err != nil {
				return err
			}
		}
	case error:
		writer.Write([]byte(fmt.Sprintf("-%s\r\n", v)))
	}

	return nil
}
