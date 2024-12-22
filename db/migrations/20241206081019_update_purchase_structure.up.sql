-- Обновление таблицы purchase с добавлением колонок (изначально допускающих NULL)
ALTER TABLE purchase 
    ADD COLUMN IF NOT EXISTS seller_id UUID,
    ADD COLUMN IF NOT EXISTS customer_id UUID;

-- Заполнение существующих записей значениями по умолчанию
UPDATE purchase SET 
    seller_id = (SELECT id FROM seller LIMIT 1),
    customer_id = (SELECT id FROM "user" LIMIT 1)
WHERE seller_id IS NULL OR customer_id IS NULL;

-- Добавление ограничения NOT NULL после заполнения данных
ALTER TABLE purchase 
    ALTER COLUMN seller_id SET NOT NULL,
    ALTER COLUMN customer_id SET NOT NULL;

-- Добавление внешних ключей с правильными ссылками
ALTER TABLE purchase
    ADD CONSTRAINT purchase_seller_id_fkey 
        FOREIGN KEY (seller_id) REFERENCES seller(id) ON DELETE CASCADE,
    ADD CONSTRAINT purchase_customer_id_fkey 
        FOREIGN KEY (customer_id) REFERENCES "user"(id) ON DELETE CASCADE;

-- Создание таблицы для связи между покупкой и объявлениями
CREATE TABLE IF NOT EXISTS purchase_advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    purchase_id UUID NOT NULL,
    advert_id UUID NOT NULL,
    FOREIGN KEY (purchase_id) REFERENCES purchase(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Добавление индекса для оптимизации поиска по purchase_id, если он не существует
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_purchase_advert_purchase_id') THEN
        CREATE INDEX idx_purchase_advert_purchase_id ON purchase_advert(purchase_id);
    END IF;
END $$;

-- Обновление таблицы purchase
ALTER TABLE purchase 
    ALTER COLUMN status TYPE purchase_status 
        USING status::purchase_status,
    ALTER COLUMN payment_method TYPE payment_method 
        USING payment_method::payment_method,
    ALTER COLUMN delivery_method TYPE delivery_method 
        USING delivery_method::delivery_method;

-- Добавление ограничения для проверки статуса
ALTER TABLE purchase 
    ADD CONSTRAINT purchase_status_check 
    CHECK (status IN ('pending', 'completed', 'in_progress', 'cancelled'));

-- Добавление ограничения для проверки метода оплаты
ALTER TABLE purchase 
    ADD CONSTRAINT payment_method_check 
    CHECK (payment_method IN ('cash', 'card'));

-- Добавление ограничения для проверки метода доставки
ALTER TABLE purchase 
    ADD CONSTRAINT delivery_method_check 
    CHECK (delivery_method IN ('pickup', 'delivery'));