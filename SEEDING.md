### Authentication service
```
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password TEXT NOT NULL
);
```

### Product Service
```
CREATE TABLE categories (  
    id          SERIAL PRIMARY KEY,  
    name        VARCHAR(255) NOT NULL UNIQUE,  
    description TEXT, 
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  
);  

CREATE TABLE products (  
    id          SERIAL PRIMARY KEY,  
    name        VARCHAR(255) NOT NULL,  
    description TEXT,  
    price       DECIMAL(10, 2) NOT NULL,  
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  
);  

CREATE TABLE product_categories (  
    product_id  INT NOT NULL,  
    category_id INT NOT NULL,  
    PRIMARY KEY (product_id, category_id),  
    FOREIGN KEY (product_id) REFERENCES products(id),  
    FOREIGN KEY (category_id) REFERENCES categories(id)  
);  

CREATE TABLE reviews (  
    id         SERIAL PRIMARY KEY,  
    product_id INT NOT NULL,  
    user_id    INT NOT NULL,  
    rating     INT CHECK (rating BETWEEN 1 AND 5),  
    comment   TEXT,  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (product_id) REFERENCES products(id),  
    FOREIGN KEY (user_id) REFERENCES users(id)  
);  

CREATE TABLE inventory (  
    product_id INT NOT NULL,  
    quantity  INT NOT NULL DEFAULT 0,  
    PRIMARY KEY (product_id),  
    FOREIGN KEY (product_id) REFERENCES products(id)  
);  
```

### Payment Service
```
CREATE TABLE payment_methods (  
    id          SERIAL PRIMARY KEY,  
    user_id     INT NOT NULL,  
    type        VARCHAR(50) NOT NULL CHECK (type IN ('credit_card', 'paypal', 'bank_transfer')),  
    provider    VARCHAR(50),  
    account_number VARCHAR(255),  
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (user_id) REFERENCES users(id)  
);  

CREATE TABLE transactions (  
    id         SERIAL PRIMARY KEY,  
    user_id    INT NOT NULL,  
    order_id    INT,  
    payment_method_id INT,  
    amount     DECIMAL(10, 2) NOT NULL,  
    status     VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (user_id) REFERENCES users(id),  
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)  
);  

CREATE TABLE invoices (  
    id         SERIAL PRIMARY KEY,  
    transaction_id INT NOT NULL,  
    invoice_number VARCHAR(255) NOT NULL UNIQUE,  
    invoice_date DATE NOT NULL,  
    total       DECIMAL(10, 2) NOT NULL,  
    status     VARCHAR(50) NOT NULL CHECK (status IN ('paid', 'unpaid', 'refunded')),  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)  
);  

CREATE TABLE refunds (  
    id         SERIAL PRIMARY KEY,  
    transaction_id INT NOT NULL,  
    amount     DECIMAL(10, 2) NOT NULL,  
    reason     TEXT,  
    status     VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processed', 'failed')),  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)  
);  

```

### Notification service
```
CREATE TABLE notification_types (  
    id          SERIAL PRIMARY KEY,  
    type_name  VARCHAR(255) NOT NULL UNIQUE,  
    description TEXT,  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  
);  

CREATE TABLE notifications (  
    id          SERIAL PRIMARY KEY,  
    user_id     INT,  
    order_id     INT,  
    type_id     INT NOT NULL,  
    message     TEXT NOT NULL,  
    status      VARCHAR(50) NOT NULL CHECK (status IN ('sent', 'failed', 'pending')),  
    sent_at     TIMESTAMP,  
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (user_id) REFERENCES users(id),  
    FOREIGN KEY (order_id) REFERENCES orders(id),  
    FOREIGN KEY (type_id) REFERENCES notification_types(id)  
);  
```

### Order service
```
CREATE TABLE carts (  
    id       SERIAL PRIMARY KEY,  
    user_id  INT NOT NULL,  
    status   VARCHAR(50) NOT NULL CHECK (status IN ('active', 'checkout')),  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (user_id) REFERENCES users(id)  
);  

CREATE TABLE cart_items (  
    id       SERIAL PRIMARY KEY,  
    cart_id  INT NOT NULL,  
    product_id INT NOT NULL,  
    quantity INT NOT NULL DEFAULT 1,  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (cart_id) REFERENCES carts(id),  
    FOREIGN KEY (product_id) REFERENCES products(id)  
);  

CREATE TABLE orders (  
    id           SERIAL PRIMARY KEY,  
    user_id      INT NOT NULL,  
    cart_id      INT NOT NULL,  
    payment_id   INT NOT NULL,  
    shipping_address_id INT,  
    total        DECIMAL(10, 2) NOT NULL,  
    status       VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),  
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (user_id) REFERENCES users(id),  
    FOREIGN KEY (cart_id) REFERENCES carts(id),  
    FOREIGN KEY (payment_id) REFERENCES transactions(id),  
    FOREIGN KEY (shipping_address_id) REFERENCES shipping_addresses(id)  
);  

CREATE TABLE order_items (  
    id       SERIAL PRIMARY KEY,  
    order_id INT NOT NULL,  
    product_id INT NOT NULL,  
    quantity INT NOT NULL DEFAULT 1,  
    price    DECIMAL(10, 2) NOT NULL,  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (order_id) REFERENCES orders(id),  
    FOREIGN KEY (product_id) REFERENCES products(id)  
);  

CREATE TABLE shipping_addresses (  
    id       SERIAL PRIMARY KEY,  
    user_id  INT NOT NULL,  
    address  TEXT NOT NULL,  
    city     VARCHAR(255) NOT NULL,  
    state    VARCHAR(255),  
    country  VARCHAR(255) NOT NULL,  
    zip_code VARCHAR(50) NOT NULL,  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (user_id) REFERENCES users(id)  
);  

CREATE TABLE order_statuses (  
    id       SERIAL PRIMARY KEY,  
    order_id INT NOT NULL,  
    status   VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),  
    note     TEXT,  
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (order_id) REFERENCES orders(id)  
);  
```