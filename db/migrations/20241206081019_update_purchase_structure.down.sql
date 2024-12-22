-- Удаление индексов
DROP INDEX IF EXISTS idx_purchase_advert_purchase_id;
DROP INDEX IF EXISTS idx_purchase_advert_seller_id;

-- Удаление ограничений для проверки методов и статуса
ALTER TABLE purchase 
    DROP CONSTRAINT IF EXISTS purchase_status_check,
    DROP CONSTRAINT IF EXISTS payment_method_check,
    DROP CONSTRAINT IF EXISTS delivery_method_check;

-- Возврат типов колонок к исходным
ALTER TABLE purchase 
    ALTER COLUMN status TYPE VARCHAR,
    ALTER COLUMN payment_method TYPE VARCHAR,
    ALTER COLUMN delivery_method TYPE VARCHAR;

-- Удаление таблицы purchase_advert
DROP TABLE IF EXISTS purchase_advert CASCADE;

-- Удаление внешних ключей
ALTER TABLE purchase
    DROP CONSTRAINT IF EXISTS purchase_seller_id_fkey,
    DROP CONSTRAINT IF EXISTS purchase_customer_id_fkey;

-- Удаление колонок
ALTER TABLE purchase
    DROP COLUMN IF EXISTS seller_id,
    DROP COLUMN IF EXISTS customer_id;