-- Создание таблицы orders, если она не существует
CREATE TABLE IF NOT EXISTS orders
(
    id         SERIAL PRIMARY KEY,
    order_id   VARCHAR(50) NOT NULL,
    amount     VARCHAR(50) NOT NULL,
    payment_id VARCHAR(50) NOT NULL,
    status     VARCHAR(20) NOT NULL CHECK (status IN ('in_process', 'canceled', 'completed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
