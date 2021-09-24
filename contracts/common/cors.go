package common

func GetCorsHeaders() map[string]string {
	return map[string]string{
		ContentTypeHeader:              ContentTypeApplicationJSON,
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "*",
		"Access-Control-Allow-Headers": "*",
	}
}
