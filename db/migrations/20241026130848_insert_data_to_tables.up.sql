CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Заполнение таблицы "user"
INSERT INTO "user" (id, username, email, password_hash, password_salt, phone_number, image_id, status)
VALUES
    (uuid_generate_v4(), 'ivan_petrov', 'ivan.petrov@example.com', 'hash1', 'salt1', '+79261234567', NULL, 'active'),
    (uuid_generate_v4(), 'anna_smirnova', 'anna.smirnova@example.com', 'hash2', 'salt2', '+79269876543', NULL, 'active');

-- Заполнение таблицы "seller" с ссылкой на таблицу "user"
INSERT INTO "seller" (id, user_id, description)
VALUES
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'ivan_petrov'), 'Продавец электроники и бытовой техники'),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'anna_smirnova'), 'Продавец модной одежды и аксессуаров');

-- Заполнение таблицы "category"
INSERT INTO category (id, title)
VALUES
    (uuid_generate_v4(), 'Электроника'),
    (uuid_generate_v4(), 'Одежда'),
    (uuid_generate_v4(), 'Техника для дома'),
    (uuid_generate_v4(), 'Спорт и отдых');

-- Заполнение таблицы "static" для хранения путей к изображениям
INSERT INTO static (id, name, path)
VALUES
    (uuid_generate_v4(), 'image1.jpg', '/images/image1.jpg'),
    (uuid_generate_v4(), 'image2.jpg', '/images/image2.jpg');

-- Заполнение таблицы "advert" с ссылками на "seller", "static" и "category"
INSERT INTO advert (id, title, description, price, seller_id, image_id, category_id, location, has_delivery, status)
VALUES
    (uuid_generate_v4(), 'Смартфон Samsung Galaxy', 'Новый смартфон Samsung с AMOLED экраном', 30000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        (SELECT id FROM static WHERE name = 'image1.jpg'), 
        (SELECT id FROM category WHERE title = 'Электроника'), 
        'Москва', TRUE, 'доступен'),
    (uuid_generate_v4(), 'Кроссовки Nike', 'Стильные и удобные кроссовки для повседневного использования', 7000, 
        (SELECT id FROM seller WHERE description = 'Продавец модной одежды и аксессуаров'), 
        (SELECT id FROM static WHERE name = 'image2.jpg'), 
        (SELECT id FROM category WHERE title = 'Одежда'), 
        'Санкт-Петербург', FALSE, 'доступен'),
    (uuid_generate_v4(), 'Стиральная машина LG', 'Энергосберегающая стиральная машина на 7 кг', 20000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        NULL, 
        (SELECT id FROM category WHERE title = 'Техника для дома'), 
        'Екатеринбург', TRUE, 'доступен'),
    (uuid_generate_v4(), 'Электросамокат Xiaomi', 'Новый электросамокат для поездок по городу', 15000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        NULL, 
        (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
        'Москва', FALSE, 'доступен'),
    -- 16 дополнительных записей для "advert"
    (uuid_generate_v4(), 'Холодильник Samsung', 'Двухкамерный холодильник с морозильной камерой', 35000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        NULL, 
        (SELECT id FROM category WHERE title = 'Техника для дома'), 
        'Москва', TRUE, 'доступен'),
    (uuid_generate_v4(), 'Ноутбук Lenovo', 'Ноутбук для работы и игр', 45000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        NULL, 
        (SELECT id FROM category WHERE title = 'Электроника'), 
        'Москва', TRUE, 'доступен'),
    (uuid_generate_v4(), 'Пылесос Dyson', 'Мощный беспроводной пылесос', 25000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        NULL, 
        (SELECT id FROM category WHERE title = 'Техника для дома'), 
        'Санкт-Петербург', TRUE, 'доступен'),
    (uuid_generate_v4(), 'Часы Apple Watch', 'Умные часы с множеством функций', 30000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        NULL, 
        (SELECT id FROM category WHERE title = 'Электроника'), 
        'Москва', TRUE, 'доступен'),
    (uuid_generate_v4(), 'Телевизор LG OLED', 'Сочный OLED экран с высоким разрешением', 70000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        NULL, 
        (SELECT id FROM category WHERE title = 'Электроника'), 
        'Екатеринбург', TRUE, 'доступен');

-- Заполнение таблицы "subscription" с ссылками на "user" и "seller"
INSERT INTO subscription (id, user_id, seller_id)
VALUES
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'ivan_petrov'), 
        (SELECT id FROM seller WHERE description = 'Продавец модной одежды и аксессуаров')),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'anna_smirnova'), 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'));

-- Заполнение таблицы "saved_advert" с ссылками на "user" и "advert"
INSERT INTO saved_advert (id, user_id, advert_id)
VALUES
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'ivan_petrov'), 
        (SELECT id FROM advert WHERE title = 'Смартфон Samsung Galaxy')),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'anna_smirnova'), 
        (SELECT id FROM advert WHERE title = 'Кроссовки Nike'));

-- Заполнение таблицы "cart" с ссылками на "user" и "advert"
INSERT INTO cart (id, user_id, advert_id)
VALUES
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'ivan_petrov'), 
        (SELECT id FROM advert WHERE title = 'Стиральная машина LG')),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'anna_smirnova'), 
        (SELECT id FROM advert WHERE title = 'Электросамокат Xiaomi'));

-- Заполнение таблицы "cart_advert" с ссылками на "cart" и "advert"
INSERT INTO cart_advert (id, cart_id, advert_id)
VALUES
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'ivan_petrov')), 
        (SELECT id FROM advert WHERE title = 'Смартфон Samsung Galaxy')),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'anna_smirnova')), 
        (SELECT id FROM advert WHERE title = 'Кроссовки Nike'));

-- Заполнение таблицы "purchase" с ссылкой на "cart"
INSERT INTO purchase (id, cart_id, status)
VALUES
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'ivan_petrov')), 
        'завершена'),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'anna_smirnova')), 
        'ожидание оплаты');
