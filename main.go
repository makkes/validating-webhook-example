package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	admissionv1 "k8s.io/api/admission/v1"
)

type AdmissionRequest struct {
	UID string `json:"uid"`
}

type AdmissionResponse struct {
	UID     string `json:"uid"`
	Allowed bool   `json:"allowed"`
}

type AdmissionReview struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Request    AdmissionRequest  `json:"request,omitempty"`
	Response   AdmissionResponse `json:"response,omitempty"`
}

func applyForAdmission(req admissionv1.AdmissionReview) admissionv1.AdmissionReview {
	res := admissionv1.AdmissionReview{
		Response: &admissionv1.AdmissionResponse{
			UID:     req.Request.UID,
			Allowed: true,
		},
	}
	res.APIVersion = req.APIVersion
	res.Kind = req.Kind
	return res
}

func main() {
	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		listenPort = "443"
	}
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("%s\n", body)

		var admissionReview admissionv1.AdmissionReview
		if err := json.Unmarshal(body, &admissionReview); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling request: %v\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := applyForAdmission(admissionReview)
		respBody, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling response: %v\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("Responding: %s\n", respBody)
		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("content-type", "application/json")
		_, err = rw.Write(respBody)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not write response: %v\n", err)
			return
		}
	})
	panic(http.ListenAndServeTLS("0.0.0.0:"+listenPort, "/tls/cert/tls.crt", "/tls/cert/tls.key", nil))
}
