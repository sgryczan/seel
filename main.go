package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var Items = []string{}

var (
	versionMajor string
	versionMinor string
	versionPatch string
)

var ControllerNamespace string
var ControllerName string

// kubeSecret
// swagger:model kubeSecret
type kubeSecret struct {
	// required: true
	// example: Secret
	Kind string `json:"kind"`
	// required: true
	// example: v1
	APIVersion string         `json:"apiVersion"`
	Metadata   secretMetadata `json:"metadata"`
	// required: true
	// example: {"admin":"YWRtaW4=", "password":"cGFzc3dvcmQ="}
	Data map[string][]byte `json:"data"`
}

// Secret Metadata
// swagger:model secretMetadata
type secretMetadata struct {
	// required: true
	// example: myAwesomeSecret
	Name string `json:"name"`
	// required: true
	// example: destinationNamespace
	Namespace string `json:"namespace"`
}

type secretData struct {
	Admin    []byte `json:"admin"`
	Password []byte `json:"password"`
}

func main() {

	r := mux.NewRouter()
	fmt.Println("Seel v%d.%d.%d", versionMajor, versionMinor, versionPatch)

	ControllerName = os.Getenv("CONTROLLER_NAME")
	if ControllerName == "" {
		log.Print("CONTROLLER_NAME not provided. Defaulting to 'sealed-secrets'")
		ControllerName = "sealed-secrets"
	}

	ControllerNamespace = os.Getenv("CONTROLLER_NAMESPACE")
	if ControllerNamespace == "" {
		log.Print("CONTROLLER_NAMESPACE not provided. Defaulting to 'adm'")
		ControllerNamespace = "adm"
	}

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/convert", ConvertHandler).Methods("POST")

	sh := http.StripPrefix("/api",
		http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/api/").Handler(sh)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/", 302)
}

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /convert convert Convert
	//
	// Converts a supplied key to a sealed key
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	// parameters:
	// - name: payload
	//   in: body
	//   description: Name to be added.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/kubeSecret"
	// responses:
	//   '200':
	//     description: Add an item to the database
	//     type: string

	w.WriteHeader(http.StatusOK)

	secret := &kubeSecret{}
	err := json.NewDecoder(r.Body).Decode(&secret)
	if err != nil {
		log.Print(err)
	}

	// Convert secret to SealedSecret
	res, err := ConvertSecret(ControllerName, ControllerNamespace, *secret)
	if err != nil {
		log.Print(err)
	}
	log.Print(res)

	fmt.Fprintf(w, "%s", res)
}

// ConvertSecret converts a native secret into an encrypted one
func ConvertSecret(cn string, ns string, s kubeSecret) (string, error) {

	json, err := json.Marshal(s)
	if err != nil {
		log.Print(err)
	}

	kubesealCmd := exec.Command("/usr/local/bin/kubeseal", "--format=yaml", fmt.Sprintf("--controller-namespace=%s", ns), fmt.Sprintf("--controller-name=%s", cn))
	fmt.Printf("Kubeseal command: %v\n", kubesealCmd)

	log.Printf("input json: %s", json)
	kubesealCmd.Stdin = strings.NewReader(fmt.Sprintf("%s", json))

	var out bytes.Buffer
	var serr bytes.Buffer
	kubesealCmd.Stdout = &out
	kubesealCmd.Stderr = &serr
	err = kubesealCmd.Run()
	if err != nil {
		log.Printf("Error running conversion command: %v", err)
	}
	//log.Print(out.String())
	//log.Print(serr.String())

	return out.String(), nil
}

// FetchCert grabs cert from kubeseal
func FetchCert(ns string, cn string) (string, error) {
	cmd := exec.Command("/usr/local/bin/kubeseal", "--fetch-cert", fmt.Sprintf("--controller-namespace=%s", ns), fmt.Sprintf("--controller-name=%s", cn))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Error running command: %v", err)
	}

	return out.String(), nil
}
