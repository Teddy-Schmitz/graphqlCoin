package graph

import (
	"io"

	"strconv"

	"fmt"

	"github.com/vektah/gqlgen/graphql"
)

func MarshalInt64(i int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.FormatInt(i, 10)))
	})
}

func UnmarshalInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("%T is not a int", v)
	}
}

func MarshalUInt64(i uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.FormatUint(i, 10)))
	})
}

func UnmarshalUInt64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	case int:
		if v < 0 {
			return 0, fmt.Errorf("%v is less then 0", v)
		}
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint64:
		return v, nil
	default:
		return 0, fmt.Errorf("%T is not a int", v)
	}
}
