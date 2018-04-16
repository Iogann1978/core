package logger

type Logger interface {
	Debug(msg string, args map[string]interface{})
	Info(msg string, args map[string]interface{})
	Fatal(msg string, args map[string]interface{})
	Panic(msg string, args map[string]interface{})
	Warn(msg string, args map[string]interface{})
}

// KV helper to map
func KV(key string, val interface{}) map[string]interface{} {
	return map[string]interface{}{key: val}
}

func Err(err error) map[string]interface{} {
	return map[string]interface{}{`err`: err.Error()}
}
