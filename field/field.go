package field

type FieldOption func(func(key string, value any))

// Value
func Value(key string, value any) FieldOption {
	return func(fn func(string, any)) { fn(key, value) }
}

// Error
func Error(err error) FieldOption {
	return Value("error", err.Error())
}
