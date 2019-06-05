package cmd

import (
	"github.com/spf13/viper"
	pb "tethys-console/services"
	"context"

	"github.com/spf13/cobra"
)

type ResStart struct {
	res *pb.ResStart
}

func (resp ResStart) Info() string {
	return resp.res.Info
}

func (resp ResStart) Error() string {
	return resp.res.Info
}

func (resp ResStart) Success() bool {
	return resp.res.Success
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a LDS",
	Long: `Start to running LDS. Login should be preceded to produce, query the block.`,
	Run: start,
}

func start(cmd *cobra.Command, args []string)  {
	monitorMode := viper.GetBool("monitor")

	var mode pb.ReqStart_Mode

	if monitorMode == true {
		infoLogger.Println("monitor mode")

		mode = pb.ReqStart_MONITOR
	} else {
		infoLogger.Println("default mode")
		mode = pb.ReqStart_DEFAULT
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, client, err := beforeExecute(&ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := client.Start(ctx, &pb.ReqStart{Mode: mode})

	responseStart := ResStart{resp}
	afterExecute(responseStart, err)
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("monitor", "m", false, "")
	viper.BindPFlag("monitor", startCmd.Flags().Lookup("monitor"))
}
