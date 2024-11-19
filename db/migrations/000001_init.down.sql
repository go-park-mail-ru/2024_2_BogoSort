-- Откат миграции 000001_init.up.sql

-- Удаление типов
DROP TYPE IF EXISTS user_status CASCADE;
DROP TYPE IF EXISTS payment_method CASCADE;
DROP TYPE IF EXISTS delivery_method CASCADE;
DROP TYPE IF EXISTS purchase_status CASCADE;
DROP TYPE IF EXISTS cart_status CASCADE;
-- Удаление триггеров для автоматического обновления поля updated_at
DROP TRIGGER IF EXISTS update_seller_updated_at ON seller;
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";
DROP TRIGGER IF EXISTS update_advert_updated_at ON advert;
DROP TRIGGER IF EXISTS update_cart_updated_at ON cart;
DROP TRIGGER IF EXISTS update_purchase_updated_at ON purchase;  

-- Удаление GIN индекса для полнотекстового поиска
DROP INDEX IF EXISTS idx_advert_fulltext;

-- Удаление триггера для установки image_id по умолчанию
DROP TRIGGER IF EXISTS trg_set_default_image_id ON "user";

-- Удаление функции обновления updated_at
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS set_default_image_id();

-- Удаление таблиц в обратном порядке
DROP TABLE IF EXISTS purchase CASCADE;
DROP TABLE IF EXISTS cart_advert CASCADE;
DROP TABLE IF EXISTS cart CASCADE;
DROP TABLE IF EXISTS saved_advert CASCADE;
DROP TABLE IF EXISTS advert CASCADE;
DROP TABLE IF EXISTS static CASCADE;
DROP TABLE IF EXISTS category CASCADE;
DROP TABLE IF EXISTS subscription CASCADE;
DROP TABLE IF EXISTS seller CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;

DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;