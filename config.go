package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const (
	introLong = `
Start K8S Watcher`
)

var RootCmd = &cobra.Command{
	Use:   "USe flags to change default configuration",
	Short: "OCS Watcher",
	Long:  `OSE Container setvice watcher`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

type Watcher struct {
	RestConfig restclient.Config
	DebugLevel uint8
}

type WatcherOptions struct {
	Username    string
	Password    string
	Insecure    bool
	Host        string
	BearerToken string
	DebugLevel  log.Level
}

func (options *WatcherOptions) init() {
	RootCmd.PersistentFlags().StringVar(&options.Host, "host", "localhost:8443", "hostname for connection")
	RootCmd.PersistentFlags().StringVar(&options.Username, "username", "admin", "Username")
	RootCmd.PersistentFlags().StringVar(&options.Password, "password", "admin", "Password")
	RootCmd.PersistentFlags().StringVar(&options.BearerToken, "token", "", "SA Token")
	RootCmd.PersistentFlags().BoolVar(&options.Insecure, "insecure", true, "Insecure connection?")
	RootCmd.Execute()
}

func (options *WatcherOptions) Validate() (watcher Watcher) {

	config := restclient.Config{}

	watcher.RestConfig.Host = options.Host
	watcher.RestConfig.Insecure = options.Insecure
	if len(options.BearerToken) > 0 {
		log.Info("Token set, will not use Username and Pass")
		config = restclient.Config{
			Host:        options.Host,
			BearerToken: options.BearerToken,
			Insecure:    options.Insecure,
		}
	} else {
		log.Info("Token Not set, will  use Username and Pass")
		config = restclient.Config{
			Host:     options.Host,
			Username: options.Username,
			Password: options.Password,
			Insecure: options.Insecure,
		}
	}

	watcher.RestConfig = config
	return watcher
}
