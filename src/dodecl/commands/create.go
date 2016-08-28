package commands

import (
	"dodecl"
	"dodecl/util"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createFile, "filename", "f", "", "filename")
}

var createFile string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources by filename or stdin",
	Run:   runCmd(createCmdFn),
}

func createCmdFn(cmd *cobra.Command, args []string) error {
	var fileReader io.Reader

	switch createFile {
	case "":
		return errors.New("create requires a file name or -")

	case "-":
		fileReader = os.Stdin

		fmt.Println("reading from stdin")

	default:
		file, err := os.Open(createFile)
		if err != nil {
			return errors.New("file does not exist or is not readable")
		}

		fileReader = file

		fmt.Println("reading", createFile)
	}

	guts, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return err
	}

	d, err := dodecl.ReadFromYAML(guts)
	if err != nil {
		return errors.Wrap(err, "read yaml from string failure")
	}

	d.ID = util.RandID(5)

	fmt.Printf("config: %#v\n", d)

	p := dodecl.NewPlanner()
	plan, err := p.Plan(d)
	if err != nil {
		return errors.Wrap(err, "plan failure")
	}

	e := dodecl.NewExecer()
	if err := e.Exec(plan); err != nil {
		return errors.Wrap(err, "exec failure")
	}

	return nil
}
