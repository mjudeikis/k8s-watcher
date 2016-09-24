package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"k8s.io/kubernetes/pkg/client/restclient"
	"os"
)

type Watcher struct {
	Config     restclient.Config
	DebugLevel uint8
}

type WatcherOptions struct {
	Username    string
	Password    string
	Insecure    bool
	Host        string
	BearerToken string
	Watcher     Watcher
}

func (options WatcherOptions) init(args []string) (config *restclient.Config) {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host",
			Value:       "localhost:8443",
			Usage:       "Master API host",
			Destination: &options.Host,
		},
		cli.StringFlag{
			Name:        "username",
			Value:       "admin",
			Usage:       "K8S/OSE username",
			Destination: &options.Username,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "admin",
			Usage:       "K8S/OSE password",
			Destination: &options.Password,
		},
		cli.StringFlag{
			Name:        "token",
			Value:       "...",
			Usage:       "K8S/OSE token",
			Destination: &options.BearerToken,
		},
	}

	app.Name = "K8S Watcher"
	app.Usage = "Watch some stuff"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}

	app.Run(os.Args)

	//add check here with termination if values empty

	log.Debugf("Token %s", options.BearerToken)
	if len(options.BearerToken) > 0 {
		log.Info("Token set, will not use Username and Pass")
		log.Debugf("Token %s...", options.BearerToken[0:9])
		config2 := &restclient.Config{
			Host:        options.Host,
			BearerToken: options.BearerToken,
			Insecure:    true,
		}
		return config2
	} else {
		log.Info("Token Not set, will  use Username and Pass")
		config2 := &restclient.Config{
			Host:     options.Host,
			Username: options.Username,
			Password: options.Password,
			Insecure: true,
		}
		return config2
	}

}
