package env

import (
	"os"
	"strings"
)

// LoadDynamicEnv loads environment variables from the environment
// based on a prefix.
// Example:
// LoadDynamicEnv("GITHUB_")
// Will load all environment variables starting with GITHUB_
// into the environment.
func LoadDynamicEnv(prefix string) map[string]string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			// skip if not the right prefix
			continue
		}
		split := strings.SplitN(env, "=", 2)
		envs[split[0]] = split[1]
	}
	return envs
}
