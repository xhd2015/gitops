package git_isolated

// Config customizes isolated git subprocess defaults.
type Config struct {
	UserEmail string
	UserName  string
	ExtraEnv  []string
}

func (c Config) userEmail() string {
	if c.UserEmail != "" {
		return c.UserEmail
	}
	return DefaultUserEmail
}

func (c Config) userName() string {
	if c.UserName != "" {
		return c.UserName
	}
	return DefaultUserName
}

// DefaultConfig is the standard isolated runner used by tests and fixtures.
var DefaultConfig = Config{}