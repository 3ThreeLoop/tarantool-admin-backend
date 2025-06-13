package tarantool

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
)

type ColumnMeta struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type QueryResult struct {
	Columns []ColumnMeta             `json:"columns"`
	Data    []map[string]interface{} `json:"data"`
}

func TestTarantoolConnection(host string, port int, username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := pool.Connect(ctx, []pool.Instance{{
		Name: fmt.Sprintf("%s:%d", host, port),
		Dialer: tarantool.NetDialer{
			Address:  fmt.Sprintf("%s:%d", host, port),
			User:     username,
			Password: password,
		},
		Opts: tarantool.Opts{MaxReconnects: 1},
	}})
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func ConnectTarantool(host string, port int, username, password string) (*pool.ConnectionPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := pool.Connect(ctx, []pool.Instance{{
		Name: fmt.Sprintf("%s:%d", host, port),
		Dialer: tarantool.NetDialer{
			Address:  fmt.Sprintf("%s:%d", host, port),
			User:     username,
			Password: password,
		},
		Opts: tarantool.Opts{MaxReconnects: 1},
	}})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func MapToDetailedQueryResult(result interface{}) (*QueryResult, error) {
	// fmt.Printf("Tarantool raw result: %#v\n", result)

	var data_map map[string]interface{}

	switch val := result.(type) {
	case map[string]interface{}:
		data_map = val
	case []interface{}:
		if len(val) == 0 {
			return nil, fmt.Errorf("empty result array")
		}
		converted := convertMapInterfaceToString(val[0])
		m, ok := converted.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected result format in array (converted type: %T)", converted)
		}
		data_map = m
	default:
		return nil, fmt.Errorf("invalid Tarantool result type: %T", val)
	}

	if errMsg, exists := data_map["error"].(string); exists {
		return nil, fmt.Errorf("tarantool error: %s", errMsg)
	}

	// parse metadata
	meta_raw, ok := data_map["metadata"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing metadata in result")
	}

	var columns []ColumnMeta
	var column_names []string
	for _, meta := range meta_raw {
		meta_map, ok := meta.(map[string]interface{})
		if !ok {
			continue
		}
		name := getString(meta_map, "name")
		typ := getString(meta_map, "type")
		columns = append(columns, ColumnMeta{Name: name, Type: typ})
		column_names = append(column_names, name)
	}

	// parse rows
	rows_raw, ok := data_map["rows"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing rows in result")
	}

	var data []map[string]interface{}
	for _, row := range rows_raw {
		row_slice, ok := row.([]interface{})
		if !ok {
			continue
		}

		row_map := make(map[string]interface{})
		for i, val := range row_slice {
			if i < len(column_names) {
				row_map[column_names[i]] = val
			}
		}
		data = append(data, row_map)
	}

	return &QueryResult{
		Columns: columns,
		Data:    data,
	}, nil
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func convertMapInterfaceToString(input interface{}) interface{} {
	switch val := input.(type) {
	case map[interface{}]interface{}:
		mapped := make(map[string]interface{})
		for k, v := range val {
			key_str := fmt.Sprintf("%v", k)
			mapped[key_str] = convertMapInterfaceToString(v)
		}
		return mapped
	case []interface{}:
		for i, elem := range val {
			val[i] = convertMapInterfaceToString(elem)
		}
	}
	return input
}
