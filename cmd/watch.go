/*
Copyright © 2020 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/k1LoW/connected/conn"
	"github.com/spf13/cobra"
)

var (
	defaultCmd = `osascript -e "set Volume 5"; say -v Alex "Disconnected."`
	interval   int
	command    string
	wifi       bool
	bluetooth  bool
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch connection",
	Long:  `watch connection.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var (
			cn  conn.Conn
			err error
		)
		switch {
		case wifi:
			cn, err = conn.NewWifi(ctx)
		case bluetooth:
			cn, err = conn.NewBluetooth(ctx)
		default:
			cn, err = conn.NewPower(ctx)
		}

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("%s. Connect and execute command again.", err))
			os.Exit(1)
		}

		fmt.Println(cn.State())

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		var c []string

		if len(args) > 0 {
			c = args
		} else {
			c = []string{"sh", "-c", command}
		}
		sigChan := make(chan os.Signal, 1)
		signal.Ignore()
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("Start watching connection (%s).", cn.Name()))

		lock := false
	L:
		for {
			select {
			case <-sigChan:
				break L
			case <-ticker.C:
				if lock {
					continue
				}
				if err := cn.Check(ctx); err != nil {
					go func() {
						lock = true
						_, _ = fmt.Fprintln(os.Stderr, err)
						_ = exec.CommandContext(ctx, c[0], c[1:]...).Run()
						time.Sleep(500 * time.Millisecond)
						lock = false
					}()
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().IntVarP(&interval, "interval", "i", 3, "watch interval (second)")
	if os.Getenv("DEBUG") != "" {
		defaultCmd = `osascript -e "set Volume 2"; say -v Alex "Disconnected."`
	}
	watchCmd.Flags().StringVarP(&command, "command", "c", defaultCmd, "command to execute when disconnected")
	watchCmd.Flags().BoolVarP(&wifi, "wifi", "", false, "watch Wi-Fi connection")
	watchCmd.Flags().BoolVarP(&bluetooth, "bluetooth", "", false, "watch Bluetooth connection")
}
