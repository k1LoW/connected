package conn

import (
	"context"
	"errors"
	"os/exec"
	"strings"
)

const wifiConnCheckCmd = "networksetup -getairportnetwork en0"
const wifiConnDefaultState = "Current Wi-Fi Network:"

type Wifi struct {
	state string
}

func NewWifi(ctx context.Context) (*Wifi, error) {
	c := &Wifi{}
	err := c.Check(ctx)
	if err != nil {
		return nil, err
	}
	return c, err
}

func (c *Wifi) Name() string {
	return "Wi-Fi"
}

func (c *Wifi) State() string {
	if c.state == "" {
		return "Wi-Fi is disconnected"
	}
	return c.state
}

func (c *Wifi) Check(ctx context.Context) error {
	contains := wifiConnDefaultState
	if c.state != "" {
		contains = c.state
	}
	o, _ := exec.CommandContext(ctx, "sh", "-c", wifiConnCheckCmd).Output()
	if !strings.Contains(string(o), contains) {
		return errors.New("Wi-Fi is disconnected or changed")
	}
	c.state = strings.Trim(string(o), "\n")
	return nil
}
