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

## features
```
1. Product Service

Display product listings with filtering and sorting options.
Show detailed product information, including descriptions, specifications, and images.
Manage product reviews and ratings.
Track inventory levels and product availability.

GET /products: Retrieve a list of products with optional filters and pagination.
GET /products/{id}: Get detailed information for a specific product.
POST /products/{id}/reviews: Submit a review for a product.
GET /products/{id}/reviews: Retrieve reviews for a specific product.

2. Authentication Service

User registration and login.
Token generation for authenticated sessions.
User profile management.
Token validation for other services.

POST /auth/register: Create a new user account.
POST /auth/login: Authenticate a user and generate a token.
GET /auth/profile: Retrieve the authenticated user's profile information.
POST /auth/validate: Validate a token for authentication.

3. Payment Service

Process payment transactions.
Generate invoices and receipts.
Handle refunds and payment method management.
Verify payment status.

POST /payment/checkout: Process a payment transaction.
GET /payment/invoice/{id}: Retrieve an invoice for a specific order.
POST /payment/refund: Initiate a refund for a transaction.
GET /payment/methods: List available payment methods for the user.
POST /payment/verify: Verify the status of a payment transaction.

4. Notification Service

Send various types of notifications via email, SMS, or in-app alerts.
Trigger notifications for order confirmations, shipment updates, payment receipts, and password reset requests.

POST /notification/order: Send an order confirmation or update notification.
POST /notification/payment: Send a payment receipt or confirmation.
POST /notification/reset-password: Send a password reset link.
POST /notification/shipping: Send shipping status updates.

5. Order Service

Manage shopping carts and wishlists.
Handle the checkout process and order creation.
Update order status (e.g., pending, delivering, completed).
Manage order cancellations and returns.
Communicate with other services for payment processing and notifications.

POST /order/cart: Add or update items in the shopping cart.
GET /order/cart: Retrieve the current shopping cart contents.
POST /order/checkout: Initiate the checkout process.
GET /order/{id}: Retrieve details for a specific order.
POST /order/{id}/status: Update the status of an order.
POST /order/{id}/cancel: Cancel an order.
POST /order/{id}/return: Initiate a return for an order.

6. Shopping Cart Management

Add and remove items from the cart.
Update item quantities.
View cart contents.

POST /cart/items: Add an item to the cart.
DELETE /cart/items/{id}: Remove an item from the cart.
PATCH /cart/items/{id}: Update the quantity of an item in the cart.
GET /cart: Retrieve the current state of the cart.

7. Order Status and Shipping

Track order status from pending to delivery.
Update shipping information and track packages.

GET /order/{id}/status: Retrieve the current status of an order.
POST /order/{id}/status: Update the status of an order (e.g., mark as delivered).
GET /order/{id}/shipping: Retrieve shipping details for an order.

```