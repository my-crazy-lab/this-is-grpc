/*
 *
 * Copyright 2018 gRPC authors.
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

// The server demonstrates how to consume and validate OAuth2 tokens provided by
// clients for each RPC.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"

	"github.com/my-crazy-lab/this-is-grpc/authentication/pg"
	authPb "github.com/my-crazy-lab/this-is-grpc/authentication/proto/account"
	pb "github.com/my-crazy-lab/this-is-grpc/authentication/proto/echo"
	"github.com/my-crazy-lab/this-is-grpc/authentication/secure"

	"github.com/my-crazy-lab/this-is-grpc/authentication/server"
)

var port = flag.Int("port", 50051, "the port to serve on")

func main() {
	flag.Parse()

	pg.InitDB()
	defer pg.CloseDB()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(data.Path("x509/server_cert.pem"), data.Path("x509/server_key.pem"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	opts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(secure.EnsureValidToken),
		// Enable TLS for all incoming connections.
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	s := grpc.NewServer(opts...)

	pb.RegisterEchoServer(s, server.NewEchoServer())
	authPb.RegisterAccountServer(s, server.NewAuthServer())

	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
