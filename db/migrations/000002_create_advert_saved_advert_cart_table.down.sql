-- Удаление таблиц которые по внешнему ключу зависят от advert
ALTER TABLE saved_advert DROP CONSTRAINT IF EXISTS saved_advert_advert_id_fkey;
ALTER TABLE cart DROP CONSTRAINT IF EXISTS cart_advert_id_fkey;
ALTER TABLE cart_advert DROP CONSTRAINT IF EXISTS cart_advert_advert_id_fkey;

-- Удаление таблиц которые по внешнему ключу зависят от cart
ALTER TABLE cart_advert DROP CONSTRAINT IF EXISTS cart_advert_cart_id_fkey;
ALTER TABLE purchase DROP CONSTRAINT IF EXISTS purchase_cart_id_fkey;

-- Удаление таблиц
DROP TABLE IF EXISTS advert;
DROP TABLE IF EXISTS saved_advert;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS static;
DROP TABLE IF EXISTS cart_advert;
DROP TABLE IF EXISTS purchase;