package customDb

// GetDsnString from the passed map returns a string of settings for the database if there are keys, otherwise an empty string.
func GetDsnString(envData map[string]string) string {
	var dsnStr string
	if val, ok := envData["DB_USER"]; ok {
		dsnStr += "user=" + val + " "
	}
	if val, ok := envData["DB_PASSWORD"]; ok {
		dsnStr += "password=" + val + " "
	}
	if val, ok := envData["DB_NAME"]; ok {
		dsnStr += "dbname=" + val + " "
	}
	if val, ok := envData["DB_SSLMODE"]; ok {
		dsnStr += "sslmode=" + val + " "
	}
	return dsnStr
}
