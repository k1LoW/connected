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
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

const powerConnCheckCmd = "ioreg -rn AppleSmartBattery | grep ExternalConnected"
const powerConnContains = "Yes"
const powerConn = "Power cable"
const wifiConnCheckCmd = "networksetup -getairportnetwork en0"
const wifiConnContains = "Current Wi-Fi Network:"
const wifiConn = "Wi-Fi"

var (
	defaultCmd = `osascript -e "set Volume 5"; say -v Alex "Disconnected."`
	interval   int
	command    string
	wifi       bool
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch connection",
	Long:  `watch connection.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		checkCmd := powerConnCheckCmd
		contains := powerConnContains
		target := powerConn
		if wifi {
			checkCmd = wifiConnCheckCmd
			contains = wifiConnContains
			target = wifiConn
		}

		o, _ := exec.CommandContext(ctx, "sh", "-c", checkCmd).Output()
		if !strings.Contains(string(o), contains) {
			_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("%s is disconnected. Connect and execute command again.", target))
			os.Exit(1)
		}

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

		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("Start watching connection (%s).", target))

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
				o, _ := exec.CommandContext(ctx, "sh", "-c", checkCmd).Output()
				if !strings.Contains(string(o), contains) {
					go func() {
						lock = true
						_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("%s disconnected.", target))
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
}
