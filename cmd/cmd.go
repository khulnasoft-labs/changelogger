package cmd

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sulaiman-coder/goeventbus"

	"github.com/khulnasoft-labs/changelogger/changelogger"
	"github.com/khulnasoft-labs/changelogger/internal/config"
	"github.com/khulnasoft-labs/changelogger/internal/log"
	"github.com/khulnasoft-labs/go-logger/adapter/logrus"
)

var (
	appConfig         *config.Application
	eventBus          *eventbus.Bus
	eventSubscription *eventbus.Subscription // nolint
)

func init() {
	cobra.OnInitialize(
		initCmdAliasBindings,
		initAppConfig,
		initLogging,
		logAppConfig,
		initEventBus,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, color.Red.Sprint(err.Error()))
		os.Exit(1)
	}
}

// we must setup the config-cli bindings first before the application configuration is parsed. However, this cannot
// be done without determining what the primary command that the config options should be bound to since there are
// shared concerns (the root-create alias).
func initCmdAliasBindings() {
	activeCmd, _, err := rootCmd.Find(os.Args[1:])
	if err != nil {
		panic(err)
	}

	if activeCmd == createCmd || activeCmd == rootCmd {
		// note: we need to lazily bind config options since they are shared between both the root command
		// and the create command. Otherwise there will be global viper state that is in contention.
		// See for more details: https://github.com/spf13/viper/issues/233 . Additionally, the bindings must occur BEFORE
		// reading the application configuration, which implies that it must be an initializer (or rewrite the command
		// initialization structure against typical patterns used with cobra, which is somewhat extreme for a
		// temporary alias)
		if err = bindCreateConfigOptions(activeCmd.Flags()); err != nil {
			panic(err)
		}
	} else {
		// even though the root command or create command is NOT being run, we still need default bindings
		// such that application config parsing passes.
		if err = bindCreateConfigOptions(createCmd.Flags()); err != nil {
			panic(err)
		}
	}
}

func initAppConfig() {
	cfg, err := config.LoadApplicationConfig(viper.GetViper(), persistentOpts)
	if err != nil {
		fmt.Printf("failed to load application config: \n\t%+v\n", err)
		os.Exit(1)
	}

	appConfig = cfg
}

func initLogging() {
	lgr, err := logrus.New(logrus.Config{
		EnableConsole: (appConfig.Log.FileLocation == "" || appConfig.CliOptions.Verbosity > 0) && !appConfig.Quiet,
		FileLocation:  appConfig.Log.FileLocation,
		Level:         appConfig.Log.LevelOpt,
	})

	if err != nil {
		panic(err)
	}

	changelogger.SetLogger(lgr)
}

func logAppConfig() {
	log.Debugf("application config:\n%+v", color.Magenta.Sprint(appConfig.String()))
}

func initEventBus() {
	eventBus = eventbus.NewBus()
	eventSubscription = eventBus.Subscribe()
	changelogger.SetBus(eventBus)
}
