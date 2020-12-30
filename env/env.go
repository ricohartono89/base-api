package env

import "os"

// AppName application name
func AppName() string {
	return os.Getenv("AppName")
}

// AppPort application port
func AppPort() string {
	return os.Getenv("AppPort")
}

// DBUser ps username
func DBUser() string {
	return os.Getenv("DBUser")
}

// DBPassword ps password
func DBPassword() string {
	return os.Getenv("DBPassword")
}

// DBAddress address host:port
func DBAddress() string {
	return os.Getenv("DBAddress")
}

// DBDatabase db name
func DBDatabase() string {
	return os.Getenv("DBDatabase")
}

// RedisUrl redis url
func RedisUrl() string {
	return os.Getenv("RedisUrl")
}

// RedisPassword ...
func RedisPassword() string {
	return os.Getenv("RedisPassword")
}

// RedisPort ...
func RedisPort() string {
	return os.Getenv("RedisPort")
}

// CorsOrigin url
func CorsOrigin() string {
	return os.Getenv("CorsOrigin")
}

// APIToken apitoken
func APIToken() string {
	return os.Getenv("APIToken")
}

// AuthJwtSigningKey ...
func AuthJwtSigningKey() string {
	return os.Getenv("AuthJwtSigningKey")
}

// AuthJwtTokenKey ...
func AuthJwtTokenKey() interface{} {
	return os.Getenv("AuthJwtTokenKey")
}

// Env apikey
func Env() string {
	return os.Getenv("Env")
}
