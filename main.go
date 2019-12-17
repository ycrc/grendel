package main

import (
	"fmt"
	"io/ioutil"
	golog "log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/ubccr/grendel/cmd"
	"github.com/ubccr/grendel/logger"
	"github.com/ubccr/grendel/util"
	"github.com/urfave/cli"
)

var release = "(version not set)"
var log = logger.GetLogger("main")

func init() {
	viper.SetConfigName("grendel")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/grendel/")
}

func main() {
	app := cli.NewApp()
	app.Name = "grendel"
	app.Version = release
	app.Usage = "provisioning system for high-performance Linux clusters"
	app.Author = "Andrew E. Bruno"
	app.Email = "aebruno2@buffalo.edu"
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "conf,c", Usage: "Path to conf file"},
		&cli.BoolFlag{Name: "verbose", Usage: "Print verbose messages"},
		&cli.BoolFlag{Name: "debug", Usage: "Print debug messages"},
	}
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			log.Logger.SetLevel(logrus.DebugLevel)
		} else if c.GlobalBool("verbose") {
			log.Logger.SetLevel(logrus.InfoLevel)
			golog.SetOutput(ioutil.Discard)
		} else {
			log.Logger.SetLevel(logrus.WarnLevel)
			golog.SetOutput(ioutil.Discard)
		}

		conf := c.GlobalString("conf")
		if len(conf) > 0 {
			viper.SetConfigFile(conf)

			err := viper.ReadInConfig()
			if err != nil {
				return fmt.Errorf("Failed reading config file: %s", err)
			}
		}

		if !viper.IsSet("secret") {
			secret, err := util.GenerateSecret(32)
			if err != nil {
				return err
			}

			viper.Set("secret", secret)
		}

		return nil
	}
	app.Commands = []cli.Command{
		cmd.NewCertsCommand(),
		cmd.NewServeCommand(),
		cmd.NewHostCommand(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
