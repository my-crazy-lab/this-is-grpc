### ADD TO CART
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "mutation { AddToCart(product_id: 22, user_id: 4, quantity: 2) { id cart_id product_id quantity created_at } }" }' \
http://localhost:9090/graphql
```

### UPDATE INVENTORY
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "mutation { UpdateInventory(product_id: 1, action: \"increase\", quantity: 2) { product_id new_quantity } }" }' \
http://localhost:9090/graphql
```

### GET PRODUCTS
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetProducts(pageSize:10, pageIndex: 1) { total products {id  name description price quantity categories {id name description created_at updated_at} created_at updated_at} } }" }' \
http://localhost:9090/graphql
```

### GET PRODUCT BY ID
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetProduct(id: 1) {id  name description price quantity categories {id name description created_at updated_at} created_at updated_at} }" }' \
http://localhost:9090/graphql
```

### CREATE PRODUCT
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "mutation { CreateProduct(name:\"product 1\", description:\"description hihihi\", price: 200.5, quantity: 2, categories: [1,2]) { id } }" }' \
http://localhost:9090/graphql
```

### GET CATEGORIES
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetCategories { id  name description created_at updated_at } }" }' \
http://localhost:9090/graphql
```

### CREATE CATEGORIES
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "mutation { CreateCategories(name:\"must unique 2\", description:\"description hihihi\") { id } }" }' \
http://localhost:9090/graphql
```

### REGISTER
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer _token_" \
--data '{ "query": "mutation { SignUp(phoneNumber:\"123456\", password:\"hihihi\") { token } }" }' \
http://localhost:9090/graphql
```

### SIGN IN
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer _token_" \
--data '{ "query": "mutation { SignIn(phoneNumber:\"123456\",password:\"hihihi\") { token } }" }' \
http://localhost:9090/graphql
```

### GET USER BY ID
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetUser { id phoneNumber } }" }' \
http://localhost:9090/graphql
```

### GET USERS
```
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetUsers { id phoneNumber } }" }' \
http://localhost:9090/graphql
```

