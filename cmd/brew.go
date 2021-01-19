/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/spf13/cobra"
)

var queues chan termbox.Event
var startDone bool
var startX, startY int

func finished() {
	fmt.Println("Your coffee is ready.")

}
func draw(d time.Duration) {
	w, h := termbox.Size()
	clear()

	str := format(d)
	text := toText(str)

	if !startDone {
		startDone = true
		startX, startY = w/2-text.width()/2, h/2-text.height()/2
	}

	x, y := startX, startY
	for _, s := range text {
		echo(s, x, y)
		x += s.width()
	}

	flush()
}

func format(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h < 1 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// brewCmd represents the brew command
var brewCmd = &cobra.Command{
	Use:   "brew",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var exitCode int
		err := termbox.Init()
		if err != nil {
			panic(err)
		}

		queues = make(chan termbox.Event)
		go func() {
			for {
				queues <- termbox.PollEvent()
			}
		}()
		timeLeft := 4 * time.Minute
		tick := time.Second

		ticker := time.NewTicker(tick)
		timer := time.NewTimer(timeLeft)

		draw(timeLeft)

	loop:
		for {
			select {
			case ev := <-queues:
				if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC) {
					exitCode = 1
					break loop
				}
				if ev.Ch == 'p' || ev.Ch == 'P' {
					timer.Stop()
					ticker.Stop()
				}
				if ev.Ch == 'c' || ev.Ch == 'C' {
					ticker = time.NewTicker(tick)
					timer = time.NewTimer(timeLeft)
				}
			case <-ticker.C:
				timeLeft -= time.Duration(tick)
				draw(timeLeft)
			case <-timer.C:
				break loop
			}
		}
		termbox.Close()

		finished()

		if exitCode != 0 {
			os.Exit(exitCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(brewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
