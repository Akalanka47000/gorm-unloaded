package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm/schema"
)

// KVToJSONMap is a custom serializer that bridges two representations of the same data:
//
//	Go field  →  key=value string     "theme=dark,lang=en,timezone=UTC"
//	DB column →  JSONB map            {"theme":"dark","lang":"en","timezone":"UTC"}
//
// This matters when: the application passes settings as simple key=value pairs
// (config files, env vars, HTTP headers) but the DB needs a structured map for
// querying with JSON operators like settings->>'theme' = 'dark'.
type KVToJSONMap struct{}

// Value is called before every write. Parses "k=v,k=v" → JSON map object for the DB.
func (KVToJSONMap) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue any) (any, error) {
	str, _ := fieldValue.(string)
	if str == "" {
		return "{}", nil
	}
	m := make(map[string]string)
	for _, pair := range strings.Split(str, ",") {
		k, v, found := strings.Cut(pair, "=")
		if found {
			m[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	data, err := json.Marshal(m)
	return string(data), err
}

// Scan is called after every read. Converts the JSONB map → "k=v,k=v" Go string.
func (KVToJSONMap) Scan(_ context.Context, field *schema.Field, dst reflect.Value, dbValue any) error {
	var raw []byte
	switch v := dbValue.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return fmt.Errorf("KVToJSONMap: cannot scan %T", dbValue)
	}

	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil {
		return err
	}
	pairs := make([]string, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, k+"="+v)
	}
	field.ReflectValueOf(context.Background(), dst).SetString(strings.Join(pairs, ","))
	return nil
}

func init() {
	schema.RegisterSerializer("kv_json", KVToJSONMap{})
}
