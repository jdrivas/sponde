package config

import "github.com/spf13/viper"

// Debug returns whether debug mode is set.
func Debug() bool {
	return viper.GetBool("debug")
}

// Verbose returs whether verbose mode is set.
func Verbose() bool {
	return viper.GetBool("verbose")
}
