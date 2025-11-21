package env

type EnvironSettings struct {
	POSTGRES_DSN string
	APP_PORT int
}

var env *EnvironSettings = nil;

func Env() *EnvironSettings {
	if env == nil {
		env = &EnvironSettings{
			POSTGRES_DSN: LoadEnvironStringOrPanic("POSTGRES_DSN"),
			APP_PORT: ToIntOrPanic(LoadEnvironStringWithDefault("APP_PORT", "8000")),
		}
	}

	return env
}
