package secrets

import (
	"context"
	"errors"
	"fmt"

	"github.com/containers/podman/v2/cmd/podman/registry"
	"github.com/containers/podman/v2/cmd/podman/utils"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:     "rm [options] SECRET [SECRET...]",
		Short:   "Remove one or more secrets",
		RunE:    rm,
		Example: "podman secret rm mysecret1 mysecret2",
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode},
		Command: rmCmd,
		Parent:  secretCmd,
	})
	flags := rmCmd.Flags()
	flags.BoolVarP(&rmOptions.All, "all", "a", false, "Remove all secrets")
}

var (
	rmOptions = entities.SecretRmOptions{}
)

func rm(cmd *cobra.Command, args []string) error {
	var (
		errs utils.OutputErrors
	)
	if (len(args) > 0 && rmOptions.All) || (len(args) < 1 && !rmOptions.All) {
		return errors.New("`podman secret rm` requires one argument, or the --all flag")
	}
	responses, err := registry.ContainerEngine().SecretRm(context.Background(), args, rmOptions)
	if err != nil {
		return err
	}
	for _, r := range responses {
		if r.Err == nil {
			fmt.Println(r.ID)
		} else {
			errs = append(errs, r.Err)
		}
	}
	return errs.PrintErrors()
}
