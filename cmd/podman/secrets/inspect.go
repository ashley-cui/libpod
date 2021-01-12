package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"text/tabwriter"

	"github.com/containers/common/pkg/report"
	"github.com/containers/podman/v2/cmd/podman/common"
	"github.com/containers/podman/v2/cmd/podman/parse"
	"github.com/containers/podman/v2/cmd/podman/registry"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	// podman secret _inspect_
	inspectCmd = &cobra.Command{
		Use:     "inspect SECRET",
		Short:   "Inspect a secret",
		Long:    "asjkdghlaiuehaudfhgilkadlh",
		RunE:    inspect,
		Example: "podman secret inspect SECRET",
		Args:    cobra.MinimumNArgs(1),
	}
)

var format string

func init() {
	// Subscribe inspect sub command to manifest command
	registry.Commands = append(registry.Commands, registry.CliCommand{
		// _podman manifest inspect_ will support both ABIMode and TunnelMode
		Mode: []entities.EngineMode{entities.ABIMode},
		// The definition for this command
		Command: inspectCmd,
		// The parent command to proceed this command on the CLI
		Parent: secretCmd,
	})
	flags := inspectCmd.Flags()
	formatFlagName := "format"
	flags.StringVar(&format, formatFlagName, "", "Format volume output using Go template")
	_ = inspectCmd.RegisterFlagCompletionFunc(formatFlagName, common.AutocompleteJSONFormat)
	// This is where you would configure the cobra flags using inspectCmd.Flags()
}

// Business logic: cmd is inspectCmd, args is the positional arguments from os.Args
func inspect(cmd *cobra.Command, args []string) error {
	fmt.Println("what a great inspect!")
	inspected, _, _ := registry.ContainerEngine().SecretInspect(context.Background(), args)
	fmt.Println(inspected)

	if cmd.Flags().Changed("format") {
		row := report.NormalizeFormat(format)
		formatted := parse.EnforceRange(row)

		tmpl, err := template.New("inspect secret").Parse(formatted)
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 12, 2, 2, ' ', 0)
		defer w.Flush()
		tmpl.Execute(w, inspected)
	} else {
		buf, err := json.MarshalIndent(inspected, "", "    ")
		if err != nil {
			return err
		}
		_, err = fmt.Println(string(buf))
	}
	return nil
}
