package env

import "github.com/joho/godotenv"

type EnvironSettings struct {
	POSTGRES_DSN string
	APP_PORT int
}

var env *EnvironSettings = nil;

func Env() *EnvironSettings {
	if env == nil {
		godotenv.Load()

		env = &EnvironSettings{
			POSTGRES_DSN: LoadEnvironStringOrPanic("POSTGRES_DSN"),
			APP_PORT: ToIntOrPanic(LoadEnvironStringWithDefault("APP_PORT", "8000")),
		}
	}

	return env
}
