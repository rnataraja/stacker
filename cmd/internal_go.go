package main

import (
	"os"

	"github.com/anuvu/stacker/lib"
	"github.com/anuvu/stacker/log"
	"github.com/anuvu/stacker/overlay"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var internalGoCmd = cli.Command{
	Name:   "internal-go",
	Hidden: true,
	Subcommands: []cli.Command{
		/*
		 * these are not actually used by stacker, but are entrypoints
		 * to the code for use in the test suite.
		 */
		cli.Command{
			Name:   "testsuite-check-overlay",
			Action: doTestsuiteCheckOverlay,
		},
		cli.Command{
			Name:   "copy",
			Action: doImageCopy,
		},
	},
	Before: doBeforeUmociSubcommand,
}

func doBeforeUmociSubcommand(ctx *cli.Context) error {
	log.Debugf("stacker subcommand: %v", os.Args)
	return nil
}

// doTestsuiteCheckOverlay is only called from the stacker test suite to
// determine if the kernel is new enough to run the full overlay test suite as
// the user it is run as.
//
// If it can do the overlay operations it exit(0)s. It prints overlay error
// returned if it cannot, and exit(50)s in that case. This way we can test for
// that error code in the test suite, vs. a standard exit(1) or exit(2) from
// urfave/cli when bad arguments are passed in the eventuality that we refactor
// this command.
func doTestsuiteCheckOverlay(ctx *cli.Context) error {
	err := os.MkdirAll(config.RootFSDir, 0755)
	if err != nil {
		return errors.Wrapf(err, "couldn't make rootfs dir for testsuite check")
	}

	err = overlay.CanDoOverlay(config)
	if err != nil {
		log.Infof("%s", err)
		os.Exit(50)
	}

	return nil
}

func doImageCopy(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return errors.Errorf("wrong number of args")
	}

	return lib.ImageCopy(lib.ImageCopyOpts{
		Src:      ctx.Args()[0],
		Dest:     ctx.Args()[1],
		Progress: os.Stdout,
	})
}
