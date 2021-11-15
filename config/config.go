package config

// VolumeMountConfig provides a configuration structure that can be used to
// mount a volume for a job container.
type VolumeMountConfig struct {
	Volume   string   `mapstructure:"volume"`
	Tags     []string `mapstructure:"tags"`
	Readonly bool     `mapstructure:"readonly"`
}

// Config defines the configuration map structure.
type Config struct {
	Debug        bool                `mapstructure:"debug"`
	DefaultImage string              `mapstructure:"default_image"`
	WorkingDir   string              `mapstructure:"working_dir"`
	CacheDir     string              `mapstructure:"cache_dir"`
	AuthFile     string              `mapstructure:"auth_file"`
	Arguments    map[string]string   `mapstructure:"extra_arguments"`
	Mounts       []VolumeMountConfig `mapstructure:"mounts"`
}
