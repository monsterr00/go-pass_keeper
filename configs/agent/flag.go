package agent

import (
	"flag"
)

func init() {
	flag.StringVar(&ClientOptions.Command, "c", "", "command")
	flag.StringVar(&ClientOptions.Data, "d", "", "data")
}
