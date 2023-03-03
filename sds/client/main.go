/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
        "fmt"
	"log"
	"time"
        "crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
        "b00m.in/crypto/util"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:18000", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
	version = flag.String("version", "3", "snapshot verison")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := secretservice.NewSecretDiscoveryServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//r, err := c.FetchSecrets(ctx, &discoverygrpc.DiscoveryRequest{Node: &v3.Node{Id: "test-id", Cluster: "example_proxy_cluster"}, VersionInfo: *version})
	r, err := c.FetchSecrets(ctx, &discoverygrpc.DiscoveryRequest{Node: &v3.Node{Id: "primero", Cluster: "boomin"}, VersionInfo: *version})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
        fmt.Printf("%s\n", r.GetVersionInfo())
        //if x := r.GetResources(); len(x) > 0 {
        for _, x := range r.GetResources() {
                switch r.TypeUrl {
                case resource.SecretType:
                        if cache.GetResponseType(r.TypeUrl) == types.Secret {
                                //sec := &tlsv3.TlsCertificate{}
                                sec := &tlsv3.Secret{}
                                //if err := anypb.UnmarshalTo(x[0], sec, proto.UnmarshalOptions{} ) ; err != nil {
                                if err := anypb.UnmarshalTo(x, sec, proto.UnmarshalOptions{} ) ; err != nil {
                                        fmt.Printf("unmarshal %v\n", err)
                                }
                                fmt.Printf("%s\n", sec.GetName())
                                cc := sec.GetTlsCertificate().GetCertificateChain().GetInlineBytes()
                                cert, err := util.PEMBytes2DERCert(cc)
                                if err != nil {
                                        fmt.Printf("pem2der %v \n", err)
                                }
                                vo := x509.VerifyOptions{CurrentTime: time.Now()}
                                chains, err := cert.Verify(vo)
                                if err != nil {
                                        fmt.Printf("verify %v %s \n" , err, cert.DNSNames[0])
                                        fmt.Printf("%v %v \n", cert.Issuer, cert.Subject )
                                } else {
                                        fmt.Printf("%v \n", chains)
                                }
                                //pk := sec.GetPrivateKey().GetInlineBytes()
                                //fmt.Printf("%v %v\n", sec.GetType(), sec.GetName())
                                //fmt.Printf("Greeting: %v", x)
                        }
                case resource.RouteType:
                        fmt.Printf("%v\n", x)
                }
        }

        GetRoutes(conn)
}

func GetRoutes(conn *grpc.ClientConn) {

        c := routeservice.NewRouteDiscoveryServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.FetchRoutes(ctx, &discoverygrpc.DiscoveryRequest{Node: &v3.Node{Id: "primero", Cluster: "boomin"}, VersionInfo: *version})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
        fmt.Printf("%s\n", r.GetVersionInfo())
        //if x := r.GetResources(); len(x) > 0 {
        for _, x := range r.GetResources() {
                switch r.TypeUrl {
                case resource.SecretType:
                        if cache.GetResponseType(r.TypeUrl) == types.Secret {
                                //sec := &tlsv3.TlsCertificate{}
                                sec := &tlsv3.Secret{}
                                //if err := anypb.UnmarshalTo(x[0], sec, proto.UnmarshalOptions{} ) ; err != nil {
                                if err := anypb.UnmarshalTo(x, sec, proto.UnmarshalOptions{} ) ; err != nil {
                                        fmt.Printf("unmarshal %v\n", err)
                                }
                                fmt.Printf("%s\n", sec.GetName())
                                cc := sec.GetTlsCertificate().GetCertificateChain().GetInlineBytes()
                                cert, err := util.PEMBytes2DERCert(cc)
                                if err != nil {
                                        fmt.Printf("pem2der %v \n", err)
                                }
                                vo := x509.VerifyOptions{CurrentTime: time.Now()}
                                chains, err := cert.Verify(vo)
                                if err != nil {
                                        fmt.Printf("verify %v %s \n" , err, cert.DNSNames[0])
                                        fmt.Printf("%v %v \n", cert.Issuer, cert.Subject )
                                } else {
                                        fmt.Printf("%v \n", chains)
                                }
                                //pk := sec.GetPrivateKey().GetInlineBytes()
                                //fmt.Printf("%v %v\n", sec.GetType(), sec.GetName())
                                //fmt.Printf("Greeting: %v", x)
                        }
                case resource.RouteType:
                        fmt.Printf("route: %v\n", x)
                }
        }
}
