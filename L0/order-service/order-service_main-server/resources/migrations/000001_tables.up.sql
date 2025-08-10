BEGIN;

CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    order_uid TEXT,
    track_number TEXT,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INT,
    date_created TIMESTAMP,
    oof_chard TEXT
);

CREATE TABLE IF NOT EXISTS delivery(
    id SERIAL PRIMARY KEY,
    name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT,
    order_id INT, 
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE IF NOT EXISTS payment(
    id SERIAL PRIMARY KEY,
    transaction TEXT,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INT,
    payment_dt INT,
    bank TEXT,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    order_id INT, 
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE IF NOT EXISTS items(
    id SERIAL PRIMARY KEY,
    chrt_id INT,
    track_number TEXT,
    price INT,
    rid TEXT,
    name TEXT,
    sale INT,
    size TEXT,
    total_price INT,
    nm_id INT,
    brand TEXT,
    status INT,
    order_id INT,
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

COMMIT;
