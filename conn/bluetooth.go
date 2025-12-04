package conn

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
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

			dm, ok := d.(map[interface{}]interface{})
			if !ok {
				continue
			}
			for k, v := range dm {
				vv, ok := v.(map[interface{}]interface{})
				if !ok {
					continue
				}
				ks, ok := k.(string)
				if !ok {
					continue
				}
				connected, ok := vv["Connected"].(bool)
				if !ok {
					continue
				}
				s[ks] = connected
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
