-- Создание сервисной учетной записи
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'emporium_service') THEN
        CREATE USER emporium_service WITH PASSWORD :'SERVICE_PASSWORD';
    END IF;
END
$$;

-- Создание роли с необходимыми правами
CREATE ROLE emporium_service_role;

-- Базовые права на схему
GRANT USAGE ON SCHEMA public TO emporium_service_role;

-- Права на таблицы (только необходимые операции)
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE 
    static,
    "user",
    seller,
    category,
    advert,
    saved_advert,
    cart,
    cart_advert,
    purchase,
    purchase_advert,
    viewed_advert,
    price_history
TO emporium_service_role;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO emporium_service_role;

-- Права на функции
GRANT EXECUTE ON FUNCTION 
    update_updated_at_column,
    set_default_image_id,
    set_advert_default_image_id
TO emporium_service_role;

-- Права на ENUM типы
GRANT USAGE ON TYPE 
    user_status,
    payment_method,
    delivery_method,
    purchase_status,
    cart_status,
    advert_status
TO emporium_service_role;

-- Права на uuid-ossp
GRANT USAGE ON EXTENSION "uuid-ossp" TO emporium_service_role;

-- Права на полнотекстовый поиск
GRANT SELECT ON TABLE pg_ts_dict, pg_ts_parser, pg_ts_config TO emporium_service_role;

GRANT emporium_service_role TO emporium_service;

-- Отзыв опасных прав
REVOKE CREATE, TEMPORARY ON DATABASE emporiumdb FROM emporium_service;
REVOKE CREATE ON SCHEMA public FROM emporium_service;
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM PUBLIC;
