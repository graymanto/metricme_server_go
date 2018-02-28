package main

import (
	"runtime"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type systemConfig struct {
	flushInterval     int
	flushRepeater     bool
	flushGraphite     bool
	flushConsole      bool
	graphiteAddress   string
	graphiteTimeout   int
	graphiteNamespace string
}

const (
	flushIntKey       = "flushInterval"
	flushRepKey       = "flushRep"
	flushGraphiteKey  = "flushGraphite"
	flushConsoleKey   = "flushConsole"
	graphAddressKey   = "graphiteAddress"
	graphTimeoutKey   = "graphiteTimeout"
	graphNamespaceKey = "graphiteNamespace"
)

func setDefaultConfigVals() {
	viper.SetDefault(flushIntKey, 10000)
	viper.SetDefault(flushRepKey, false)
	viper.SetDefault(flushGraphiteKey, false)
	viper.SetDefault(flushConsoleKey, true)
	viper.SetDefault(graphAddressKey, "")
	viper.SetDefault(graphTimeoutKey, 0)
	viper.SetDefault(graphNamespaceKey, "")
}

func loadConfigurationFile() {
	viper.SetConfigName("mm_config")

	var confFilePath string
	if runtime.GOOS == "windows" {
		confFilePath = "$HOME/_metricme"

	} else {
		confFilePath = "$HOME/.metricme"
	}

	viper.AddConfigPath(confFilePath)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		// Do nothing for now, as we don't expect a config file
	}
}

func loadAndBindFlags() {
	flag.Int(flushIntKey, 10000, "Flush interval")
	flag.Bool(flushRepKey, false, "Flush to repeater enabled")
	flag.Bool(flushGraphiteKey, false, "Flush to graphite enabled")
	flag.Bool(flushConsoleKey, false, "Flush to console enabled")
	flag.String(graphAddressKey, "", "Graphite server address in form localhost:2004")
	flag.Int(graphTimeoutKey, 0, "Graphite client connection timeout in seconds")
	flag.String(graphNamespaceKey, "", "A prefix for all metrics flushed to graphite")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
}

func loadConfig() systemConfig {
	setDefaultConfigVals()
	loadConfigurationFile()
	loadAndBindFlags()

	return systemConfig{
		flushInterval:     viper.GetInt(flushIntKey),
		flushRepeater:     viper.GetBool(flushRepKey),
		flushGraphite:     viper.GetBool(flushGraphiteKey),
		flushConsole:      viper.GetBool(flushConsoleKey),
		graphiteAddress:   viper.GetString(graphAddressKey),
		graphiteTimeout:   viper.GetInt(graphTimeoutKey),
		graphiteNamespace: viper.GetString(graphNamespaceKey),
	}
}
