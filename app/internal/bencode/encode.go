package bencode

import (
	"fmt"
	"io"
	"reflect"
	"sort"
)

func Encode(w io.Writer, val interface{}) error {
	switch v := val.(type) {
	case string:
		_, err := fmt.Fprintf(w, "%d:%s", len(v), v)
		return err

	case int:
		_, err := fmt.Fprintf(w, "i%de", v)
		return err

	case int64:
		_, err := fmt.Fprintf(w, "i%de", v)
		return err

	case []interface{}:
		_, err := w.Write([]byte("l"))
		if err != nil {
			return err
		}
		for _, item := range v {
			if err := Encode(w, item); err != nil {
				return err
			}
		}
		_, err = w.Write([]byte("e"))
		return err

	case map[string]interface{}:
		return encodeMap(w, v)

	default:
		// Use reflection for custom struct
		rv := reflect.ValueOf(val)
		rt := reflect.TypeOf(val)

		if rt.Kind() == reflect.Struct {
			m := make(map[string]interface{})

			for i := 0; i < rt.NumField(); i++ {
				field := rt.Field(i)
				key := field.Tag.Get("bencode")
				if key == "" {
					key = field.Name
				}
				m[key] = rv.Field(i).Interface()
			}
			return encodeMap(w, m)
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
