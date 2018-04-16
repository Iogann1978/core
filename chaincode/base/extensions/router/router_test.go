package router

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var handlerFunc = func(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte(`123`))
}

// example using Convey
func TestGroup(t *testing.T) {
	Convey(`Testing new group instance`, t, func() {
		g := New()
		Convey(`new group instance should be not nil`, func() {
			So(g, ShouldNotBeNil)
			So(g.methods, ShouldBeEmpty)
			So(g.prefix, ShouldEqual, ``)
			Convey(`adding new handler...`, func() {
				g.Add(`/handler`, handlerFunc)
				So(g.methods, ShouldContainKey, `/handler`)
				Convey(`and testing calling handler with mockstub...`, func() {
					//ms := shim.NewMockStub()
				})
			})
		})
	})
}
