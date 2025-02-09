## Like Shopee: Order management system project to lean gRPC

- API Gateway: GraphQL with Apollo/Federation
    - Multi client
- Database: PostgreSQL, MongoDB (for product catalog), Redis (for caching)
    - use 2 databases for learning
- Login -> Authentication service
    - Jwt based by phone
- Go to list Products -> Product Service
    - Get lists: name, description, price, categories, amount(inventory tracking)
    - Get reviews
    - Add to carts product + quantity
- Add to carts
    - Checkout -> Payment service
        - invoice generate
        - refund
    - Choose address
    - Status shipping -> Order Service
        - pending -> shopper confirm shipping for shipper
        - delivering -> customer confirm at address 
- Notification service

## run protoc
```
protoc --go_out=authentication/proto/account --go_opt=paths=source_relative \
       --go-grpc_out=authentication/proto/account --go-grpc_opt=paths=source_relative \
       proto/account.proto
```