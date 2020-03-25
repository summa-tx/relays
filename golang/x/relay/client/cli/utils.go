package cli

import (
    "io/ioutil"
    "github.com/spf13/cobra"
)

func readJSONFromFile(filename string) ([]byte, error) {
    return ioutil.ReadFile("scripts/json_data/" + filename)
}

func attachFlagFileinput(cmd *cobra.Command) {
    cmd.Flags().Bool("inputfile", false, "Accepts a file as input for each json parameter")
}
