package compat

import (
	"net/http"

	"github.com/containers/podman/v3/libpod"
	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/api/handlers/utils"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/domain/infra/abi"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

func StopContainer(w http.ResponseWriter, r *http.Request) {
	runtime := r.Context().Value("runtime").(*libpod.Runtime)
	decoder := r.Context().Value("decoder").(*schema.Decoder)
	// Now use the ABI implementation to prevent us from having duplicate
	// code.
	containerEngine := abi.ContainerEngine{Libpod: runtime}

	// /{version}/containers/(name)/stop
	query := struct {
		Ignore        bool `schema:"ignore"`
		DockerTimeout uint `schema:"t"`
		LibpodTimeout uint `schema:"timeout"`
	}{
		// override any golang type defaults
	}
	if err := decoder.Decode(&query, r.URL.Query()); err != nil {
		utils.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest,
			errors.Wrapf(err, "failed to parse parameters for %s", r.URL.String()))
		return
	}

	name := utils.GetName(r)

	options := entities.StopOptions{
		Ignore: query.Ignore,
	}
	if utils.IsLibpodRequest(r) {
		if query.LibpodTimeout > 0 {
			options.Timeout = &query.LibpodTimeout
		}
	} else {
		if query.DockerTimeout > 0 {
			options.Timeout = &query.DockerTimeout
		}
	}
	con, err := runtime.LookupContainer(name)
	if err != nil {
		utils.ContainerNotFound(w, name, err)
		return
	}
	state, err := con.State()
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	if state == define.ContainerStateStopped || state == define.ContainerStateExited {
		utils.WriteResponse(w, http.StatusNotModified, nil)
		return
	}
	report, err := containerEngine.ContainerStop(r.Context(), []string{name}, options)
	if err != nil {
		if errors.Cause(err) == define.ErrNoSuchCtr {
			utils.ContainerNotFound(w, name, err)
			return
		}

		utils.InternalServerError(w, err)
		return
	}

	if len(report) > 0 && report[0].Err != nil {
		utils.InternalServerError(w, report[0].Err)
		return
	}

	// Success
	utils.WriteResponse(w, http.StatusNoContent, nil)
}
