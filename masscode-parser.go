package main

import (
	"masscode-parser/cli/cmd"
	"runtime"

	"github.com/ondrovic/common/utils/cli"
)

func main() {
	cli.ClearTerminalScreen(runtime.GOOS)
	cmd.Execute()
}
