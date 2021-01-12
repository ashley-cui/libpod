package compat

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/containers/podman/v2/libpod"
	"github.com/containers/podman/v2/pkg/api/handlers/utils"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/containers/podman/v2/pkg/domain/infra/abi"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

func ListSecrets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("listing stuff!!")
	var (
		runtime = r.Context().Value("runtime").(*libpod.Runtime)
	)

	ic := abi.ContainerEngine{Libpod: runtime}
	reports, err := ic.SecretList(r.Context())
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	utils.WriteResponse(w, http.StatusOK, reports)
}

func InspectSecret(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inspecting stuff!! in handler")
	var (
		runtime = r.Context().Value("runtime").(*libpod.Runtime)
	)

	name := utils.GetName(r)
	names := []string{name}
	ic := abi.ContainerEngine{Libpod: runtime}
	reports, _, err := ic.SecretInspect(r.Context(), names)
	// if len(errs) != 0 {
	// 	fmt.Println("here!")
	// 	msg := fmt.Sprintf("No such secret: %s", name)
	// 	utils.Error(w, msg, http.StatusNotFound, err)
	// 	return
	// }
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	utils.WriteResponse(w, http.StatusOK, reports[0])
}

func RemoveSecret(w http.ResponseWriter, r *http.Request) {
	fmt.Println("removing stuff!! in handler")
	var (
		runtime = r.Context().Value("runtime").(*libpod.Runtime)
	)

	name := utils.GetName(r)
	ic := abi.ContainerEngine{Libpod: runtime}
	_, err := ic.SecretRm(r.Context(), name)
	if err != nil {
		// fmt.Println("HOOOO HEE")
		// fmt.Println(err.Error())
		// if strings.Contains(err.Error(), "no such secret") {
		// 	fmt.Println("adfg")
		// 	utils.Error(w, "name AIM LITERADFJTGH", http.StatusNotFound, err)
		// 	return
		// }
		utils.InternalServerError(w, err)
		return
	}
	utils.WriteResponse(w, http.StatusNoContent, nil)
}

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
	ic := abi.ContainerEngine{Libpod: runtime}

	var createParams entities.SecretCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		utils.Error(w, "Something went wrong.", http.StatusInternalServerError, errors.Wrap(err, "Decode()"))
		return
	}
	fmt.Println(createParams)
	decoded, _ := base64.StdEncoding.DecodeString(createParams.Data)
	reader := bytes.NewReader(decoded)
	opts.Driver = createParams.Driver.Name
	fmt.Println("LOOKERE")
	fmt.Println(createParams.Name)
	report, err := ic.SecretCreate(r.Context(), createParams.Name, reader, opts)
	if err != nil {
		if err.Error() == "secret name in use" {
			utils.Error(w, "name conflicts with an existing object", http.StatusConflict, err)
			return
		}
		utils.InternalServerError(w, err)
	}
	utils.WriteResponse(w, http.StatusOK, report)
}
