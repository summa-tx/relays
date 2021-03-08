package cli

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

func readJSONFromFile(filename string) ([]byte, error) {
	// get path to root directory
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// if running this function from cli_test directory do not include it in golang path
	path = strings.TrimSuffix(path, "/cli_test")
	return ioutil.ReadFile("/" + path + "/scripts/json_data/" + filename)
}

func attachFlagFileinput(cmd *cobra.Command) {
	cmd.Flags().Bool("inputfile", false, "Accepts a file as input for each json parameter")
}
