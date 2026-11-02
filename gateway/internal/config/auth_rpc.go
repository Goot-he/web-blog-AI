package config

type AuthConf struct {
	JwtSecret     string `yaml:"jwt_secret"`
	AccessExpire  int    `yaml:"access_expire"`
	RefreshExpire int    `yaml:"refresh_expire"`
	CookieSecure  bool   `yaml:"cookie_secure"`
}
