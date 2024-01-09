package patch

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
)

type UpdateFunc func(target *interface{}) error

type Accessor interface {
	get(object interface{}, key string) error
	set(object interface{}, key string, value interface{}) error
}

type defaultAccessor struct {
}

func (d *defaultAccessor) set(object interface{}, key string, update UpdateFunc) error {
	if v, ok := object.(*[]interface{}); ok {
		object = *v
	} else if v, ok := object.(*map[string]interface{}); ok {
		object = *v
	}

	switch object.(type) {
	case nil:
		return fmt.Errorf("empty object")
	case []interface{}:
		sli := object.([]interface{})
		i, err := parseIndex(key, len(sli))
		if err != nil {
			return err
		}
		cpy := sli[i]
		if err := update(&cpy); err != nil {
			return err
		}
		sli[i] = cpy
		return nil
	case map[string]interface{}:
		cpy, ok := object.(map[string]interface{})[key]
		if !ok {
			return fmt.Errorf("failed to find %s in object", key)
		}
		if err := update(&cpy); err != nil {
			return err
		}
		object.(map[string]interface{})[key] = cpy
		return nil
	default:
		return fmt.Errorf("unknown object type %s", reflect.TypeOf(object).String())
	}
}

func (d *defaultAccessor) get(object interface{}, key string) (interface{}, error) {

	if v, ok := object.(*[]interface{}); ok {
		object = *v
	} else if v, ok := object.(*map[string]interface{}); ok {
		object = *v
	}

	switch object.(type) {
	case nil:
		return nil, fmt.Errorf("empty object")
	case []interface{}:
		sli := object.([]interface{})
		i, err := parseIndex(key, len(sli))
		if err != nil {
			return nil, err
		}
		return sli[i], nil
	case map[string]interface{}:
		v, ok := object.(map[string]interface{})[key]
		if !ok {
			return nil, fmt.Errorf("failed to find %s in object", key)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown object type %s", reflect.TypeOf(object).String())
	}
}

func PartialPatch(object interface{}, path string, patch UpdateFunc) error {
	keys := strings.Split(path, ".")

	if reflect.ValueOf(object).Kind() == reflect.Struct {
		b, err := json.Marshal(&object)
		if err != nil {
			return err
		}
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}

		object = m
	}

	a := &defaultAccessor{}
	field := object
	for i, key := range keys {
		if i < len(keys)-1 {
			v, err := a.get(field, keys[i])
			if err != nil {
				return err
			}
			field = v
		} else {
			if err := a.set(field, key, patch); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseIndex(index string, length int) (int, error) {
	i, err := strconv.Atoi(index)
	if err != nil {
		return 0, fmt.Errorf("failed to parse index: %s", err)
	}
	if i < 0 || i > length {
		return 0, fmt.Errorf("invalid index value: %d", i)
	}
	return i, nil
}

func StructToMap(obj interface{}, m map[string]interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	return nil
}

func MapToStruct(m map[string]interface{}, obj interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	return nil
}

func MergePatchGenerator(mutate func(patchDoc json.RawMessage) (patch interface{}, err error)) UpdateFunc {
	return func(target *interface{}) error {
		origin, err := json.Marshal(*target)
		if err != nil {
			return err
		}

		cpy := json.RawMessage(origin)

		patchedObject, err := mutate(cpy)
		if err != nil {
			return err
		}

		patch, err := json.Marshal(&patchedObject)
		if err != nil {
			return err
		}

		patched, err := jsonpatch.MergePatch(origin, patch)
		if err != nil {
			return err
		}

		var patchedMap map[string]interface{}
		if err := json.Unmarshal(patched, &patchedMap); err != nil {
			return err
		}
		*target = patchedMap
		return nil
	}
}
