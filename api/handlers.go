package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var ControllerNamespace string
var ControllerName string

func MapController() error {
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
	return nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/", 302)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /create create Create
	//
	// Creates a new encrypted secret object from raw data
	// ---
	// consumes:
	// - application/json
	// produces:
	// - text/plain
	// parameters:
	// - name: payload
	//   in: body
	//   description: Details of secret to be created
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/secretRequest"
	// responses:
	//   '200':
	//     description: Creation successful
	//     type: string
	//   '400':
	//     description: Invalid request body
	//     type: string
	//   '500':
	//     description: Creation error
	//     type: string

	request := &SecretRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if !(ValidateJSON(err, w)) {
		return
	}

	// Convert Map
	m := ConvertMap(request.Data)
	// Turn request into valid Secret template
	secret := &KubeSecret{
		Metadata: SecretMetadata{
			Name:      request.Name,
			Namespace: request.Namespace,
		},
		Data: m,
	}

	// Convert secret to SealedSecret
	res, err := ConvertSecret(ControllerName, ControllerNamespace, *secret)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert secret")
		return
	}

	log.Print("conversion succeeded")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
	//
}

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /convert convert Convert
	//
	// Converts a supplied K8S Secret manifest (JSON) to an Encrypted Secret
	// ---
	// consumes:
	// - application/json
	// produces:
	// - text/plain
	// parameters:
	// - name: payload
	//   in: body
	//   description: Existing K8S Secret in JSON format.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/kubeSecret"
	// responses:
	//   '200':
	//     description: Conversion successful
	//     type: string
	//   '400':
	//     description: Invalid request body
	//     type: string
	//   '500':
	//     description: Conversion Error
	//     type: string

	secret := &KubeSecret{}
	err := json.NewDecoder(r.Body).Decode(&secret)
	if !(ValidateJSON(err, w)) {
		return
	}

	// Convert secret to SealedSecret
	res, err := ConvertSecret(ControllerName, ControllerNamespace, *secret)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert secret")
		return
	}

	//log.Print(res)
	log.Print("conversion succeeded")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
}

// ConvertSecret converts a native secret into an encrypted one
func ConvertSecret(cn string, ns string, s KubeSecret) (string, error) {

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
		return "", err
	}
	//log.Print(out.String())
	log.Print(serr.String())

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

// ValidateJSON gracefully handles JSON Decoder-related errors
func ValidateJSON(err error, w http.ResponseWriter) bool {
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return false
	}
	return true
}

func ConvertMap(m map[string]string) map[string][]byte {
	log.Print("ConvertMap -")
	byteMap := map[string][]byte{}
	for k, v := range m {
		log.Printf("Converting key: %s", k)
		log.Printf("Value: %s", v)
		log.Printf("Bytes: %v", []byte(m[k]))
		byteMap[k] = []byte(m[k])
	}
	log.Printf("Map: %v", byteMap)
	return byteMap
}
