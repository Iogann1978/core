// Package router provides base router for using in chaincode Invoke function
package router

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"sort"
)

var (
	ErrMethodNotFound = errors.New(`method not found`)
	ErrNoRoutes       = errors.New(`no routes presented`)
)

type Group struct {
	prefix  string
	methods map[string]func(stub shim.ChaincodeStubInterface) peer.Response
}

// Handle used for using in CC Invoke function
// Must be called after adding new routes using Add function
func (g *Group) Handle(stub shim.ChaincodeStubInterface) peer.Response {
	fnString, _ := stub.GetFunctionAndParameters()
	if fn, ok := g.methods[fnString]; ok {
		return fn(stub)
	}
	return shim.Error(ErrMethodNotFound.Error())
}

// Group gets new group using presented path
// New group can be used as independent
func (g *Group) Group(path string) *Group {
	return &Group{prefix: g.prefix + path, methods: g.methods}
}

// Add adds new handler using presented path
// Sets methods handler container
func (g *Group) Add(path string, fn func(stub shim.ChaincodeStubInterface) peer.Response) *Group {
	g.methods[g.prefix+path] = fn
	return g
}

// Routes
// Get ordered string view of routes
func (g *Group) Routes() ([]string, error) {
	rLen := len(g.methods)
	if rLen == 0 {
		return nil, ErrNoRoutes
	}
	r := make([]string, len(g.methods))
	i := 0
	for k, _ := range g.methods {
		r[i] = k
		i++
	}
	sort.Strings(r)
	return r, nil
}

func New() *Group {
	g := new(Group)
	g.methods = make(map[string]func(stub shim.ChaincodeStubInterface) peer.Response)
	return g
}
