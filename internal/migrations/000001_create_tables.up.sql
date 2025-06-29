CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(10),
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(255),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(255),
    address TEXT,
    region VARCHAR(255),
    email VARCHAR(255)
);
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(10),
    provider VARCHAR(255),
    amount INTEGER,
    payment_dt BIGINT,
    bank VARCHAR(255),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INTEGER,
    track_number VARCHAR(255),
    price INTEGER,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INTEGER,
    size VARCHAR(50),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(255),
    status INTEGER
);
-- Вставка тестовых данных
INSERT INTO orders (
        order_uid,
        track_number,
        entry,
        locale,
        internal_signature,
        customer_id,
        delivery_service,
        shardkey,
        sm_id,
        date_created,
        oof_shard
    )
VALUES (
        'b563feb7b2b84b6test',
        'WBILMTESTTRACK',
        'WBIL',
        'en',
        '',
        'test',
        'meest',
        '9',
        99,
        '2021-11-26T06:22:19Z',
        '1'
    );
INSERT INTO deliveries (
        order_uid,
        name,
        phone,
        zip,
        city,
        address,
        region,
        email
    )
VALUES (
        'b563feb7b2b84b6test',
        'Test Testov',
        '+9720000000',
        '2639809',
        'Kiryat Mozkin',
        'Ploshad Mira 15',
        'Kraiot',
        'test@gmail.com'
    );
INSERT INTO payments (
        order_uid,
        transaction,
        request_id,
        currency,
        provider,
        amount,
        payment_dt,
        bank,
        delivery_cost,
        goods_total,
        custom_fee
    )
VALUES (
        'b563feb7b2b84b6test',
        'b563feb7b2b84b6test',
        '',
        'USD',
        'wbpay',
        1817,
        1637907727,
        'alpha',
        1500,
        317,
        0
    );
INSERT INTO items (
        order_uid,
        chrt_id,
        track_number,
        price,
        rid,
        name,
        sale,
        size,
        total_price,
        nm_id,
        brand,
        status
    )
VALUES (
        'b563feb7b2b84b6test',
        9934930,
        'WBILMTESTTRACK',
        453,
        'ab4219087a764ae0btest',
        'Mascaras',
        30,
        '0',
        317,
        2389212,
        'Vivienne Sabo',
        202
    );