package utils

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"s7ab-platform-hyperledger/platform/core/entities"
	"strings"
)

// Org 3
const ORG3 = "Org3MSP"

// Org 2
const ORG2 = "Org2MSP"

// Org 1
const ORG1 = "Org1MSP"

func Error(w http.ResponseWriter, r *http.Request, err string, code int) {

	logrus.WithFields(logrus.Fields{
		"path": r.RequestURI,
	}).Error(err)

	res := &entities.Response{
		Success: false,
		Error:   err,
	}

	bytes, _err := json.Marshal(res)

	if _err != nil {
		http.Error(w, _err.Error(), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(bytes), code)
}

func Success(w http.ResponseWriter, r *http.Request, payload interface{}) {

	logrus.WithFields(logrus.Fields{
		"path": r.RequestURI,
	}).Info(payload)

	res := &entities.Response{
		Success: true,
		Result:  payload,
	}

	bytes, err := json.Marshal(res)

	if err != nil {
		Error(w, r, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func ArrayToChaincodeArgs(args []string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func ToChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}
