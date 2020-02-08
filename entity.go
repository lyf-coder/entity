// Copyright © 2020 - present. liyongfei <liyongfei@walktotop.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package entity // import "github.com/lyf-coder/entity"
import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cast"
)

// Entity is json access type like github.com/spf13/viper
type Entity struct {
	// Delimiter that separates a list of keys
	// used to access a nested value in one go
	keyDelim string

	data map[string]interface{}
}

// New returns an initialized Entity instance.
func New(data map[string]interface{}) *Entity {
	entity := new(Entity)
	entity.keyDelim = ":"
	entity.data = data
	return entity
}

// NewByJSON returns an initialized Entity instance by json byte[].
func NewByJSON(data []byte) *Entity {
	mapData := make(map[string]interface{})
	err := json.Unmarshal(data, &mapData)
	if err != nil {
		log.Println(err)
	}
	return New(mapData)
}

// deepSearch scans deep maps, following the key indexes listed in the
// sequence "path".
// The last value is expected to be another map, and is returned.
//
// In case intermediate keys do not exist, or map to a non-map value,
// a new map is created and inserted, and the search continues from there:
// the initial map "m" may be modified!
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}

// Set sets the value for the key in the Entity
func (entity *Entity) Set(key string, value interface{}) *Entity {
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, entity.keyDelim)
	lastKey := path[len(path)-1]
	deepestMap := deepSearch(entity.data, path[0:len(path)-1])

	// set innermost value
	deepestMap[lastKey] = value

	return entity
}

// toCaseInsensitiveValue checks if the value is a  map;
// if so, create a copy and recursively.
func toCaseInsensitiveValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		value = copyAndInsensitiveMap(cast.ToStringMap(v))
	case map[string]interface{}:
		value = copyAndInsensitiveMap(v)
	}

	return value
}

// copyAndInsensitiveMap  creates a copy of any map it makes case insensitive.
func copyAndInsensitiveMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})

	for key, val := range m {
		switch v := val.(type) {
		case map[interface{}]interface{}:
			nm[key] = copyAndInsensitiveMap(cast.ToStringMap(v))
		case map[string]interface{}:
			nm[key] = copyAndInsensitiveMap(v)
		default:
			nm[key] = v
		}
	}

	return nm
}

// searchMap recursively searches for a value for path in source map.
// Returns nil if not found.
func (entity *Entity) searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// Fast path
		if len(path) == 1 {
			return next
		}

		// Nested case
		switch next := next.(type) {
		case map[interface{}]interface{}:
			return entity.searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return entity.searchMap(next, path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		}
	}
	return nil
}

// isPathShadowedInDeepMap makes sure the given path is not shadowed somewhere
// on its path in the map.
// e.g., if "foo.bar" has a value in the given map, it “shadows”
//       "foo.bar.baz" in a lower-priority map
func (entity *Entity) isPathShadowedInDeepMap(path []string, m map[string]interface{}) string {
	var parentVal interface{}
	for i := 1; i < len(path); i++ {
		parentVal = entity.searchMap(m, path[0:i])
		if parentVal == nil {
			// not found, no need to add more path elements
			return ""
		}
		switch parentVal.(type) {
		case map[interface{}]interface{}:
			continue
		case map[string]interface{}:
			continue
		default:
			// parentVal is a regular value which shadows "path"
			return strings.Join(path[0:i], entity.keyDelim)
		}
	}
	return ""
}

// find
func (entity *Entity) find(key string) interface{} {
	var (
		val    interface{}
		path   = strings.Split(key, entity.keyDelim)
		nested = len(path) > 1
	)

	val = entity.searchMap(entity.data, path)
	if val != nil {
		return val
	}

	// compute the path through the nested maps to the nested value
	if nested && entity.isPathShadowedInDeepMap(path, entity.data) != "" {
		return nil
	}

	return nil
}

// Get can retrieve any value given the key to use.
// Get returns an interface. For a specific value use one of the Get____ methods.
func (entity *Entity) Get(key string) interface{} {
	val := entity.find(key)
	if val == nil {
		return nil
	}
	return val
}

// GetString returns the value associated with the key as a string.
func (entity *Entity) GetString(key string) string {
	return cast.ToString(entity.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (entity *Entity) GetBool(key string) bool {
	return cast.ToBool(entity.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (entity *Entity) GetInt(key string) int {
	return cast.ToInt(entity.Get(key))
}

// GetInt32 returns the value associated with the key as an integer.
func (entity *Entity) GetInt32(key string) int32 {
	return cast.ToInt32(entity.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (entity *Entity) GetInt64(key string) int64 {
	return cast.ToInt64(entity.Get(key))
}

// GetUint returns the value associated with the key as an unsigned integer.
func (entity *Entity) GetUint(key string) uint {
	return cast.ToUint(entity.Get(key))
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func (entity *Entity) GetUint32(key string) uint32 {
	return cast.ToUint32(entity.Get(key))
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (entity *Entity) GetUint64(key string) uint64 {
	return cast.ToUint64(entity.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func (entity *Entity) GetFloat64(key string) float64 {
	return cast.ToFloat64(entity.Get(key))
}

// GetTime returns the value associated with the key as time.
func (entity *Entity) GetTime(key string) time.Time {
	return cast.ToTime(entity.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (entity *Entity) GetDuration(key string) time.Duration {
	return cast.ToDuration(entity.Get(key))
}

// GetSlice returns the value associated with the key as a slice.
func (entity *Entity) GetSlice(key string) []interface{} {
	return cast.ToSlice(entity.Get(key))
}

// GetStringMapSlice returns the value associated with the key as a []map[string]interface{}  slice.
func (entity *Entity) GetStringMapSlice(key string) []map[string]interface{} {
	v, _ := ToStringMapSlice(entity.Get(key))
	return v
}

// ToStringMapSlice casts an interface to a []map[string]interface{} type.
func ToStringMapSlice(i interface{}) ([]map[string]interface{}, error) {
	var s []map[string]interface{}

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			s = append(s, cast.ToStringMap(u))
		}
		return s, nil
	case []map[string]interface{}:
		s = append(s, v...)
		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []map[string]interface{}", i, i)
	}
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func (entity *Entity) GetIntSlice(key string) []int {
	return cast.ToIntSlice(entity.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (entity *Entity) GetStringSlice(key string) []string {
	return cast.ToStringSlice(entity.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (entity *Entity) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(entity.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (entity *Entity) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(entity.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (entity *Entity) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(entity.Get(key))
}

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func (entity *Entity) GetSizeInBytes(key string) uint {
	sizeStr := cast.ToString(entity.Get(key))
	return parseSizeInBytes(sizeStr)
}

func safeMul(a, b uint) uint {
	c := a * b
	if a > 1 && b > 1 && c/b != a {
		return 0
	}
	return c
}

// parseSizeInBytes converts strings like 1GB or 12 mb into an unsigned integer number of bytes
func parseSizeInBytes(sizeStr string) uint {
	sizeStr = strings.TrimSpace(sizeStr)
	lastChar := len(sizeStr) - 1
	multiplier := uint(1)

	if lastChar > 0 {
		if sizeStr[lastChar] == 'b' || sizeStr[lastChar] == 'B' {
			if lastChar > 1 {
				switch unicode.ToLower(rune(sizeStr[lastChar-1])) {
				case 'k':
					multiplier = 1 << 10
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'm':
					multiplier = 1 << 20
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'g':
					multiplier = 1 << 30
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				default:
					multiplier = 1
					sizeStr = strings.TrimSpace(sizeStr[:lastChar])
				}
			}
		}
	}

	size := cast.ToInt(sizeStr)
	if size < 0 {
		size = 0
	}

	return safeMul(uint(size), multiplier)
}
