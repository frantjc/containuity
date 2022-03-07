package conf

import "github.com/spf13/pflag"

func BindPortFlag(flag *pflag.Flag) {
	sviper.BindPFlag(portKey, flag)
}

func BindVerboseFlag(flag *pflag.Flag) {
	sviper.BindPFlag(verboseKey, flag)
}

func BindSocketFlag(flag *pflag.Flag) {
	sviper.BindPFlag(socketKey, flag)
}

func BindGitHubTokenFlag(flag *pflag.Flag) {
	sviper.BindPFlag(githubTokenKey, flag)
}

func BindRuntimeImageFlag(flag *pflag.Flag) {
	sviper.BindPFlag(runtimeImageKey, flag)
}

func BindRuntimeNameFlag(flag *pflag.Flag) {
	sviper.BindPFlag(runtimeNameKey, flag)
}

func BindSecretsFlag(flag *pflag.Flag) {
	sviper.BindPFlag(secretsKey, flag)
}
