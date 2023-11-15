package cmd

import (
	"log/slog"
	"time"

	"github.com/jonkerj/chargeflux/internal/submitter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	influxerCmd = &cobra.Command{
		Use:   "influxer",
		Short: "Submit SmartEVSE values to InfluxDB2",
		Run:   influxer,
	}
)

func init() {
	flags := influxerCmd.PersistentFlags()
	flags.Duration("interval", time.Minute, "Interval between polls, expressed as golang duration")
	flags.String("influxdb-url", "influxdb.influxdb", "URL to influxDB2")
	flags.String("influxdb-token", "", "token to authenticate to influxDB2")
	flags.String("influxdb-org", "influxdata", "org in influxDB2")
	flags.String("influxdb-bucket", "smartevse", "bucket in influxDB2")
	flags.String("tag", "location=home", "tag(s) to set on measurements")

	rootCmd.AddCommand(influxerCmd)

	err := viper.BindPFlags(flags)
	if err != nil {
		panic(err)
	}
}

func influxer(cmd *cobra.Command, args []string) {
	sub, err := submitter.NewSubmitter(
		viper.GetString("url"),
		viper.GetString("influxdb-url"),
		viper.GetString("influxdb-token"),
		viper.GetString("influxdb-org"),
		viper.GetString("influxdb-bucket"),
		viper.GetString("tag"),
	)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(viper.GetDuration("interval"))
	for ; true; <-ticker.C {
		slog.Info("polling")
		if err := sub.Work(); err != nil {
			slog.Error("error polling", "error", err)
			return
		}
	}
}
