package bencode

import (
	"fmt"
	"io"
	"reflect"
	"sort"
)

func Encode(encodedWriter io.Writer, val interface{}) error {
	switch v := val.(type) {
	case string:
		_, err := fmt.Fprintf(encodedWriter, "%d:%s", len(v), v)
		return err

	case []byte:
		_, err := fmt.Fprintf(encodedWriter, "%d:%s", len(v), string(v))
		return err

	case int, int8, int64:
		_, err := fmt.Fprintf(encodedWriter, "i%de", v)
		return err

	case map[string]interface{}:
		return encodeMap(encodedWriter, v)

	default:
		// Use reflection for custom struct and list
		reflectValue := reflect.ValueOf(val)
		reflectType := reflect.TypeOf(val)

		if reflectType.Kind() == reflect.Struct {
			m := make(map[string]interface{})

			for i := 0; i < reflectType.NumField(); i++ {
				field := reflectType.Field(i)
				key := field.Tag.Get("bencode")
				if key == "" {
					key = field.Name
				}
				m[key] = reflectValue.Field(i).Interface()
			}
			return encodeMap(encodedWriter, m)
		} else if reflectType.Kind() == reflect.Slice {
			_, err := encodedWriter.Write([]byte("l"))
			if err != nil {
				return err
			}
			for i := range reflectValue.Len() {
				listEntry := reflectValue.Index(i).Interface()
				if err := Encode(encodedWriter, listEntry); err != nil {
					return err
				}
			}
			_, err = encodedWriter.Write([]byte("e"))
			return err
		}

		return fmt.Errorf("unsupported type: %T", val)
	}
}

func encodeMap(w io.Writer, m map[string]interface{}) error {
	_, err := w.Write([]byte("d"))
	if err != nil {
		return err
	}

	// Bencode requires dictionary keys to be sorted
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		// Encode key (must be string)
		if err := Encode(w, k); err != nil {
			return err
		}
		// Encode value
		if err := Encode(w, m[k]); err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("e"))
	return err
}
