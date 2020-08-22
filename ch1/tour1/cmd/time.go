package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
	"time"
	"tour1/internal/timer"
)

var calculateTime string
var duration string
var timeCmd = &cobra.Command{
	Use: "time",
	Short: "time format convert",
	Long: "time format convert",
	Run: func (cmd *cobra.Command, args []string){},
}

var nowCommand = &cobra.Command{
	Use: "now",
	Short: "get current time",
	Long: "get current time",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		log.Printf("out: %s, %d\n", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}


var calculateTimeCmd = &cobra.Command{
	Use: "calc",
	Short: "calc time",
	Long: "calc time",
	Run: func(cmd *cobra.Command, args []string) {
		var currentTimer time.Time
		var layout = "2006-01-02 15:04:05"
		if calculateTime == "" {
			currentTimer = timer.GetNowTime()
		} else {
			var err error
			if !strings.Contains(calculateTime, " ") {
				layout = "2006-01-02"
			}
			currentTimer, err = time.Parse(layout, calculateTime)
			if err != nil {
				t, _ := strconv.Atoi(calculateTime)
				currentTimer = time.Unix(int64(t), 0)
			}
		}
		calculateTime, err := timer.GetCalculateTime(currentTimer, duration)
		if err != nil {
			log.Fatalf("timer.GetCalculateTime err: %v\n", err)
		}
		log.Printf("out: %s, %d\n", calculateTime.Format(layout), calculateTime.Unix())
	},
}

func init() {
	timeCmd.AddCommand(nowCommand)
	timeCmd.AddCommand(calculateTimeCmd)

	calculateTimeCmd.Flags().StringVarP(&calculateTime, "calculate", "c", "", "time")
	calculateTimeCmd.Flags().StringVarP(&duration, "duration", "d", "", "duration")
}

