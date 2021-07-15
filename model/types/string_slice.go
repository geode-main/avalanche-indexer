package types

type StringSlice []string

func (slice *StringSlice) Scan(src interface{}) error {
	return nil
}
