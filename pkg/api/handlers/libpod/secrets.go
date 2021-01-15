package libpod

import (
	"fmt"
	"net/http"

	"github.com/containers/podman/v2/libpod"
	"github.com/containers/podman/v2/pkg/api/handlers/utils"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/domain/infra/abi"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

func CreateSecret(w http.ResponseWriter, r *http.Request) {
	fmt.Println("creating stuff!! in handler")
	var (
		runtime = r.Context().Value("runtime").(*libpod.Runtime)
		decoder = r.Context().Value("decoder").(*schema.Decoder)
	)
	query := struct {
		Name   string `schema:"name"`
		Driver string `schema:"driver"`
	}{
		// override any golang type defaults
	}
	opts := entities.SecretCreateOptions{}
	if err := decoder.Decode(&query, r.URL.Query()); err != nil {
		utils.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest,
			errors.Wrapf(err, "failed to parse parameters for %s", r.URL.String()))
		return
	}
	opts.Driver = query.Driver
	ic := abi.ContainerEngine{Libpod: runtime}
	// // report, err := ic.SecretRm(r.Context(), name)
	// // if err != nil {
	// // 	utils.InternalServerError(w, err)
	// // 	return
	// // }

	report, _ := ic.SecretCreate(r.Context(), query.Name, r.Body, opts)
	utils.WriteResponse(w, http.StatusOK, report)
}
