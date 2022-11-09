package main

import (
	"context"
	"os"
	"time"
	"wikibot/pkg/config"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func execute(ctx context.Context, event lambdaTestEvent) error {
	config := &config.Config{}

	cmd := cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.RegisterConfigFile()
			if err != nil {
				return err
			}

			if !event.shouldBeIgnored() {
				period, err := time.ParseDuration(event.Period)
				if err != nil {
					return err
				}
				config.WikiClientConfig.RefreshDelay = period
				config.TelegramClientConfig.ChanID = event.ChanID
			}

			err = run(ctx, config)
			return err
		},
	}

	config.RegisterFlags(cmd.PersistentFlags())

	err := cmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Getenv("ISLAMBDA")) == 0 {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := execute(ctx, lambdaTestEvent{})
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		lambda.Start(handler)
	}
}

func handler(ctx context.Context, event lambdaTestEvent) error {
	return execute(ctx, event)
}

type lambdaTestEvent struct {
	ChanID int    `json:"chan_id"`
	Period string `json:"period"`
}

func (e *lambdaTestEvent) shouldBeIgnored() bool {
	if e.ChanID == 0 || e.Period == "" {
		return true
	}
	return false
}
