package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

func applyForAdmission(req AdmissionReview) AdmissionReview {
	return AdmissionReview{
		APIVersion: req.APIVersion,
		Kind:       req.Kind,
		Response: AdmissionResponse{
			UID:     req.Request.UID,
			Allowed: true,
		},
	}
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
		var admissionReview AdmissionReview
		err = json.Unmarshal(body, &admissionReview)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := applyForAdmission(admissionReview)
		respBody, err := json.Marshal(resp)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
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
