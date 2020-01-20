package conn

import (
	"context"
	"errors"
	"os/exec"
	"strings"
)

const powerConnCheckCmd = "ioreg -rn AppleSmartBattery | grep ExternalConnected"
const powerConnDefaultState = "Yes"

type Power struct {
	state string
}

func NewPower(ctx context.Context) (*Power, error) {
	c := &Power{}
	err := c.Check(ctx)
	if err != nil {
		return nil, err
	}
	return c, err
}

func (c *Power) Name() string {
	return "Power cable"
}

func (c *Power) State() string {
	if c.state == "" {
		return "Power cable is disconnected"
	}
	return "Power cable is connected"
}

func (c *Power) Check(ctx context.Context) error {
	contains := powerConnDefaultState
	if c.state != "" {
		contains = c.state
	}
	o, _ := exec.CommandContext(ctx, "sh", "-c", powerConnCheckCmd).Output()
	if !strings.Contains(string(o), contains) {
		return errors.New("Power cable is disconnected")
	}
	c.state = strings.Trim(string(o), "\n")
	return nil
}
