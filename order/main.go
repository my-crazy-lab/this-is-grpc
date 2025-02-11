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
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/my-crazy-lab/this-is-grpc/product/pg"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/utils"

	"github.com/my-crazy-lab/this-is-grpc/order/server"
)

var port = flag.Int("port", 50053, "the port to serve on")

func main() {
	flag.Parse()

	pg.InitDB()
	defer pg.CloseDB()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := utils.GetOptsServer()

	s := grpc.NewServer(opts...)

	order.RegisterOrderServer(s, server.NewOrderServer())

	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
