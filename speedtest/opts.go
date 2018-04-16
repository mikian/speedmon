package speedtest

import (
	"flag"
	"time"
)

type Opts struct {
	SpeedInBytes bool
	Quiet        bool
	List         bool
	Server       ServerID
	Interface    string
	Timeout      time.Duration
	Secure       bool
	Help         bool
	Version      bool
}

func ParseOpts() *Opts {
	opts := new(Opts)

	flag.BoolVar(&opts.SpeedInBytes, "bytes", false,
		"Display values in bytes instead of bits. Does not affect the image generated by -share")
	flag.BoolVar(&opts.Quiet, "quiet", true, "Suppress verbose output, only show basic information")
	flag.BoolVar(&opts.List, "list", false, "Display a list of speedtest.net servers sorted by distance")
	flag.Uint64Var((*uint64)(&opts.Server), "server", 0, "Specify a server ID to test against")
	flag.StringVar(&opts.Interface, "interface", "", "IP address of network interface to bind to")
	flag.DurationVar(&opts.Timeout, "timeout", 10 * time.Second, "HTTP timeout duration. Default 10s")
	flag.BoolVar(&opts.Secure, "secure", false,
		"Use HTTPS instead of HTTP when communicating with speedtest.net operated servers")
	flag.BoolVar(&opts.Help, "help", false, "Show usage information and exit")
	flag.BoolVar(&opts.Help, "h", false, "Shorthand for -help option")
	flag.BoolVar(&opts.Version, "version", false, "Show the version number and exit")

	flag.Parse();

	return opts
}
