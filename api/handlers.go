package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /convert convert Convert
	//
	// Converts a supplied key to a sealed key
	// ---
	// consumes:
	// - application/json
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
	if err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			fmt.Fprintf(w, "Unable to decode JSON body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		case *json.InvalidUTF8Error:
			log.Print("Invalid UTF8 Error")
		case *json.InvalidUnmarshalError:
			log.Print("Invalid Marshaller Error")
		case *json.MarshalerError:
			log.Print("Marshaller Error")
		case *json.UnsupportedTypeError:
			log.Print("Unsupported Type Error")
		case *json.UnsupportedValueError:
			log.Print("Unsupported Value Error")
		case *json.UnmarshalFieldError:
			log.Print("Unmarshal field error")
		case *json.UnmarshalTypeError:
			log.Print("Unmarshal type error")
		default:
			log.Printf("Undefined error: %+v", err)
		}
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
