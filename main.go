package main

import (
	"fmt"

	"github.com/avrebarra/basabaso/cmd"
	"github.com/fatih/color"
)

func main() {
	c, err := cmd.New(cmd.Config{})
	if err != nil {
		msg := fmt.Sprintf("setup failure: %s", err.Error())
		fmt.Println(color.RedString(msg))
		return
	}

	if err := c.Execute(); err != nil {
		msg := fmt.Sprintf("command failed: %s", err.Error())
		fmt.Println(color.RedString(msg))
		return
	}
}
