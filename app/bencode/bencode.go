package bencode

import (
	"bufio"
	"fmt"
	"strconv"
)

// will read from the buffer based on following conditions:
// until delimiter is seen or the capacity is reached.
// it'll not go beyond the size of the buffer.
func readBytesTillDelim(buffer *bufio.Reader, delim *byte, cap int) ([]byte, error) {
	remainingBytes := buffer.Buffered()
	n := min(cap, remainingBytes)
	var bytesBuffer []byte

	for range n {
		b, err := buffer.ReadByte()
		if err != nil {
			return nil, err
		}
		bytesBuffer = append(bytesBuffer, b)
		if delim != nil && b == *delim {
			break
		}
	}

	return bytesBuffer, nil
}

func DecodeBencode(buffer *bufio.Reader) (interface{}, error) {
	char, err := buffer.ReadByte()
	if err != nil {
		return "", err
	}

	var delim byte

	switch char {

	case 'i':
		// handle parsing of integers
		delim = byte('e')
		bytes, err := readBytesTillDelim(buffer, &delim, buffer.Buffered())
		if err != nil {
			return nil, err
		}

		bytes = bytes[:len(bytes)-1] // remove e from end.

		val64, err := strconv.ParseInt(string(bytes), 10, 64)
		if err != nil {
			return "", err
		}
		return val64, nil

	case 'l':
		// handle parsing of lists
		list := []interface{}{}
		for {
			c, err := buffer.ReadByte() // read a char to check if list ends, an empty list: example: "le"
			if err == nil {
				if c == 'e' {
					return list, nil
				} else {
					buffer.UnreadByte()
				}
			}

			value, err := DecodeBencode(buffer) // decode the remainder of the string.
			if err != nil {
				return nil, err
			}

			list = append(list, value)
		}
	case 'd':
		// handle parsing of dictionary
		dictionary := map[string]interface{}{}
		for {
			c, rbErr := buffer.ReadByte()
			if rbErr == nil {
				if c == 'e' {
					return dictionary, nil
				} else {
					buffer.UnreadByte()
				}
			}

			// Decode dictionary key
			value, err := DecodeBencode(buffer)
			if err != nil {
				return nil, err
			}

			key, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("bencode decode: non string key")
			}

			// Decode value
			value, err = DecodeBencode(buffer)
			if err != nil {
				return nil, err
			}

			dictionary[key] = value
		}

	default:
		// handle parsing of strings
		delim = byte(':')
		buffer.UnreadByte()
		// Read the length of string first.
		bytes, err := readBytesTillDelim(buffer, &delim, buffer.Buffered())
		if err != nil {
			return nil, err
		}

		bytes = bytes[:len(bytes)-1] // remove : from end.
		strLength, err := strconv.ParseInt(string(bytes), 10, 64)
		if err != nil {
			return "", err
		}
		// fmt.Println("Length of string:", strLength)

		stringContent, err := readBytesTillDelim(buffer, nil, int(strLength))
		if err != nil {
			return "", err
		}

		return string(stringContent), nil
	}

	return "", fmt.Errorf("Only strings are supported at the moment")
}
