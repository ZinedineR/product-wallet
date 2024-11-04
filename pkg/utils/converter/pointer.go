package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ToString(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	default:
		val := reflect.ValueOf(data)

		// Handle if the input is a pointer to a struct
		if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
			// Dereference the pointer
			val = val.Elem()
			v = val.Interface()
		}

		// Check if it's now a struct after dereferencing
		if val.Kind() == reflect.Struct {
			jsonData, err := json.Marshal(v)
			if err != nil {
				// Return an error message or a fallback string if marshaling fails
				return fmt.Sprintf("Error converting to JSON: %v", err)
			}
			return string(jsonData) // Return JSON string
		}

		// Fallback for any other types (e.g., slices, maps, etc.)
		return fmt.Sprintf("%v", v)
	}
}

// ToInt converts any data to int
func ToInt(data interface{}) (int, error) {
	switch v := data.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return i, nil
	default:
		return 0, fmt.Errorf("cannot convert type %s to int", reflect.TypeOf(data))
	}
}

// ToInt64 converts any data to int64
func ToInt64(data interface{}) (int64, error) {
	switch v := data.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	default:
		return 0, fmt.Errorf("cannot convert type %s to int64", reflect.TypeOf(data))
	}
}

// ToUint64 converts any data to uint64
func ToUint64(data interface{}) (uint64, error) {
	switch v := data.(type) {
	case int:
		if v < 0 {
			return 0, fmt.Errorf("cannot convert negative int to uint64")
		}
		return uint64(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("cannot convert negative int32 to uint64")
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("cannot convert negative int64 to uint64")
		}
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("cannot convert negative float64 to uint64")
		}
		return uint64(v), nil
	case string:
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	default:
		return 0, fmt.Errorf("cannot convert type %s to uint64", reflect.TypeOf(data))
	}
}

// ToFloat64 converts any data to float64
func ToFloat64(data interface{}) (float64, error) {
	switch v := data.(type) {
	case float64:
		return v, nil
	case float32:
		return strconv.ParseFloat(strconv.FormatFloat(float64(v), 'f', -1, 64), 64)
	case int:
		return strconv.ParseFloat(strconv.FormatInt(int64(v), 10), 64)
	case int32:
		return strconv.ParseFloat(strconv.FormatInt(int64(v), 10), 64)
	case int64:
		return strconv.ParseFloat(strconv.FormatInt(v, 10), 64)
	case uint:
		return strconv.ParseFloat(strconv.FormatUint(uint64(v), 10), 64)
	case uint32:
		return strconv.ParseFloat(strconv.FormatUint(uint64(v), 10), 64)
	case uint64:
		return strconv.ParseFloat(strconv.FormatUint(v, 10), 64)
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert type %s to float64", reflect.TypeOf(data))
	}
}
