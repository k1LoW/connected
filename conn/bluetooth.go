package conn

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/goccy/go-yaml"
)

const bluetoothConnCheckCmd = "system_profiler SPBluetoothDataType"

type bstate map[string]bool

func (s bstate) String() string {
	o := []string{}
	for k, v := range s {
		if v {
			o = append(o, fmt.Sprintf("%s: Connected", k))
		} else {
			o = append(o, fmt.Sprintf("%s: Disconnected", k))
		}
	}
	return strings.Join(o, "\n")
}

type Bluetooth struct {
	state bstate
}

func NewBluetooth(ctx context.Context) (*Bluetooth, error) {
	c := &Bluetooth{}
	err := c.Check(ctx)
	if err != nil {
		return nil, err
	}
	return c, err
}

func (c *Bluetooth) Name() string {
	return "Bluetooth"
}

func (c *Bluetooth) State() string {
	if len(c.state) == 0 {
		return "Bluetooth devices are disconnected"
	}
	return c.state.String()
}

func (c *Bluetooth) Check(ctx context.Context) error {
	o, _ := exec.CommandContext(ctx, "sh", "-c", bluetoothConnCheckCmd).Output()
	v := make(map[string]map[string]interface{})
	err := yaml.Unmarshal(o, &v)
	if err != nil {
		return err
	}
	s := make(bstate)
	b := v["Bluetooth"]
	for k, v := range b {
		if strings.Contains(k, "Devices") {
			d := v
			for k, v := range d.(map[string]interface{}) {
				vv := v.(map[string]interface{})
				if vv["Connected"].(string) == "Yes" {
					s[k] = true
				} else {
					s[k] = false
				}
			}
		}
	}
	if len(s) == 0 {
		return errors.New("Bluetooth devices are disconnected")
	}

	disconnectedAll := true
	for k, connected := range s {
		if connected {
			disconnectedAll = false
		} else {
			if prev, ok := c.state[k]; ok {
				if prev {
					return fmt.Errorf("Bluetooth device (%s) is disconnected", k)
				}
			}
		}
	}

	if disconnectedAll {
		return errors.New("Bluetooth devices are disconnected")
	}

	if len(c.state) == 0 {
		c.state = s
	}
	return nil
}
