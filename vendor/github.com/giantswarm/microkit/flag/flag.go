package flag

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	microerror "github.com/giantswarm/microkit/error"
)

func Init(f interface{}) {
	b, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	for k, v := range m {
		m[k] = toValue([]string{toCase(k)}, k, v)
	}
	b, err = json.Marshal(m)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, f)
	if err != nil {
		panic(err)
	}
}

// Merge merges the given flag set with an internal viper configuration. That
// way command line flags, environment variables and config files will be
// merged.
func Merge(v *viper.Viper, fs *pflag.FlagSet, dirs, files []string) error {
	// Use the given viper for internal configuration management. We merge the
	// defined flags with their upper case counterparts from the environment.
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.BindPFlags(fs)

	for _, configDir := range dirs {
		v.AddConfigPath(configDir)
	}

	for _, configFile := range files {
		// Check the defined config file.
		v.SetConfigName(configFile)
		err := v.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// In case there is no config file given we simply go ahead to check
				// the other ones. If we do not find any configuration using config
				// files, we go ahead to check the process environment.
			} else {
				return microerror.MaskAny(err)
			}
		}

		fs.VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				// The current flag was set via the command line. We definitly want to use
				// the set value. Therefore we do not merge anything into it.
				return
			}
			if !v.IsSet(f.Name) {
				// There is neither configuration in the provided config file nor in the
				// process environment. That means we cannot use it to merge it into any
				// defined flag.
				return
			}

			f.Value.Set(v.GetString(f.Name))
		})
	}

	return nil
}

func toValue(path []string, key string, val interface{}) interface{} {
	m, ok := val.(map[string]interface{})
	if ok {
		for k, v := range m {
			m[k] = toValue(append([]string{toCase(k)}, path...), k, v)
		}

		return m
	}

	res := strings.Join(reverse(path), ".")
	return res
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func toCase(k string) string {
	a := []rune(k)
	d := false

	for i, c := range a {
		if d {
			return string(a)
		}

		if unicode.IsUpper(c) {
			a[i] = unicode.ToLower(a[i])
		} else {
			d = true
		}
	}

	return string(a)
}
