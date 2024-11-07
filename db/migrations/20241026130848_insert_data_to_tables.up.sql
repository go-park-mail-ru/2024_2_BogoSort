-- Заполнение таблицы "static"
INSERT INTO static (id, name, path, created_at)
VALUES
    (uuid_generate_v4(), 'image1.jpg', 'static/images', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image2.jpg', 'static/images', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image3.jpg', 'static/images', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image4.jpg', 'static/images', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image5.jpg', 'static/images', CURRENT_TIMESTAMP);

-- Заполнение таблицы "user"
INSERT INTO "user" (id, username, email, password_hash, password_salt, phone_number, image_id, status, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'ivan_petrov', 'ivan.petrov@example.com', 'hash1', 'salt1', '+79261234567', NULL, 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'anna_smirnova', 'anna.smirnova@example.com', 'hash2', 'salt2', '+79269876543', NULL, 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'pavel_ivanov', 'pavel.ivanov@example.com', 'hash3', 'salt3', '+79261239876', NULL, 'inactive', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'elena_kuznetsova', 'elena.kuznetsova@example.com', 'hash4', 'salt4', '+79261231234', NULL, 'banned', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'dmitry_sokolov', 'dmitry.sokolov@example.com', 'hash5', 'salt5', '+79261234568', NULL, 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Заполнение таблицы "seller"
INSERT INTO seller (id, user_id, description, created_at, updated_at)
VALUES
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'ivan_petrov'), 'Продавец электроники и бытовой техники', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'anna_smirnova'), 'Продавец модной одежды и аксессуаров', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'pavel_ivanov'), 'Продавец спортивного инвентаря', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'elena_kuznetsova'), 'Продавец книг и канцелярии', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'dmitry_sokolov'), 'Продавец мебели и декора', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Заполнение таблицы "category"
INSERT INTO category (id, title, created_at)
VALUES
    (uuid_generate_v4(), 'Электроника', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Одежда', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Техника для дома', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Спорт и отдых', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'Книги и канцелярия', CURRENT_TIMESTAMP);

-- Заполнение таблицы "advert"
INSERT INTO advert (id, title, description, price, seller_id, image_id, category_id, created_at, updated_at, location, has_delivery, status)
VALUES
    (uuid_generate_v4(), 'Смартфон Samsung Galaxy', 'Новый смартфон Samsung с AMOLED экраном', 30000, 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники'), 
        (SELECT id FROM static WHERE name = 'image1.jpg'), 
        (SELECT id FROM category WHERE title = 'Электроника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),
    (uuid_generate_v4(), 'Кроссовки Nike', 'Стильные и удобные кроссовки для повседневного использования', 7000, 
        (SELECT id FROM seller WHERE description = 'Продавец модной одежды и аксессуаров'), 
        (SELECT id FROM static WHERE name = 'image2.jpg'), 
        (SELECT id FROM category WHERE title = 'Одежда'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', FALSE, 'active'),
    (uuid_generate_v4(), 'Гантели 10 кг', 'Комплект гантелей для домашнего использования', 2000, 
        (SELECT id FROM seller WHERE description = 'Продавец спортивного инвентаря'), 
        (SELECT id FROM static WHERE name = 'image3.jpg'), 
        (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),
    (uuid_generate_v4(), 'Книга "Война и мир"', 'Классическое произведение Льва Толстого', 500, 
        (SELECT id FROM seller WHERE description = 'Продавец книг и канцелярии'), 
        (SELECT id FROM static WHERE name = 'image4.jpg'), 
        (SELECT id FROM category WHERE title = 'Книги и канцелярия'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', FALSE, 'active'),
    (uuid_generate_v4(), 'Диван угловой', 'Удобный угловой диван для гостиной', 45000, 
        (SELECT id FROM seller WHERE description = 'Продавец мебели и декора'), 
        (SELECT id FROM static WHERE name = 'image5.jpg'), 
        (SELECT id FROM category WHERE title = 'Техника для дома'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active');

-- Заполнение таблицы "subscription"
INSERT INTO subscription (id, user_id, seller_id)
VALUES
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'ivan_petrov'), 
        (SELECT id FROM seller WHERE description = 'Продавец модной одежды и аксессуаров')),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'anna_smirnova'), 
        (SELECT id FROM seller WHERE description = 'Продавец электроники и бытовой техники')),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'pavel_ivanov'), 
        (SELECT id FROM seller WHERE description = 'Продавец книг и канцелярии')),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'elena_kuznetsova'), 
        (SELECT id FROM seller WHERE description = 'Продавец спортивного инвентаря')),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'dmitry_sokolov'), 
        (SELECT id FROM seller WHERE description = 'Продавец мебели и декора'));

-- Заполнение таблицы "saved_advert"
INSERT INTO saved_advert (id, user_id, advert_id, created_at)
VALUES
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'ivan_petrov'), 
        (SELECT id FROM advert WHERE title = 'Смартфон Samsung Galaxy'), CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'anna_smirnova'), 
        (SELECT id FROM advert WHERE title = 'Кроссовки Nike'), CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'pavel_ivanov'), 
        (SELECT id FROM advert WHERE title = 'Гантели 10 кг'), CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'elena_kuznetsova'), 
        (SELECT id FROM advert WHERE title = 'Книга "Война и мир"'), CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'dmitry_sokolov'), 
        (SELECT id FROM advert WHERE title = 'Диван угловой'), CURRENT_TIMESTAMP);

-- Заполнение таблицы "cart"
INSERT INTO cart (id, user_id, created_at, updated_at)
VALUES
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'ivan_petrov'), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'anna_smirnova'), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'pavel_ivanov'), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'elena_kuznetsova'), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM "user" WHERE username = 'dmitry_sokolov'), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Заполнение таблицы "cart_advert"
INSERT INTO cart_advert (id, cart_id, advert_id)
VALUES
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'ivan_petrov')), 
        (SELECT id FROM advert WHERE title = 'Смартфон Samsung Galaxy')),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'anna_smirnova')), 
        (SELECT id FROM advert WHERE title = 'Кроссовки Nike')),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'pavel_ivanov')), 
        (SELECT id FROM advert WHERE title = 'Гантели 10 кг')),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'elena_kuznetsova')), 
        (SELECT id FROM advert WHERE title = 'Книга "Война и мир"')),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'dmitry_sokolov')), 
        (SELECT id FROM advert WHERE title = 'Диван угловой'));

-- Заполнение таблицы "purchase"
INSERT INTO purchase (id, cart_id, status, adress, payment_method, delivery_method, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'ivan_petrov')), 
        'completed', 'ул. Ленина, д. 1', 'card', 'delivery', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'anna_smirnova')), 
        'pending', 'ул. Пушкина, д. 2', 'cash', 'pickup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'pavel_ivanov')), 
        'in_progress', 'ул. Чехова, д. 3', 'card', 'delivery', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'elena_kuznetsova')), 
        'cancelled', 'ул. Гоголя, д. 4', 'cash', 'pickup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'dmitry_sokolov')), 
        'completed', 'ул. Тургенева, д. 5', 'card', 'delivery', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO category (id, title, created_at)
VALUES
    ('d4d10f10-4f9a-4bd5-ab1e-d2fc3ed35748', 'Женский гардероб', CURRENT_TIMESTAMP),
    ('d49a98a6-f041-4432-b255-f23d4a97edde', 'Мужской гардероб', CURRENT_TIMESTAMP),
    ('f21963b7-fd2b-4770-97f0-8dfac77c6155', 'Детский гардероб', CURRENT_TIMESTAMP),
    ('97f4f702-5412-4588-8e53-b682499df8c7', 'Детские товары', CURRENT_TIMESTAMP),
    ('1a4f92f6-c6f5-4930-91e8-163ec679ed0d', 'Стройматериалы и инструменты', CURRENT_TIMESTAMP),
    ('aeeb6c57-b428-450c-8049-fbb942aa0d1c', 'Компьютерная техника', CURRENT_TIMESTAMP),
    ('e310f974-9ea8-4e78-ad86-4fb49c92842a', 'Для дома и дачи', CURRENT_TIMESTAMP),
    ('c513af49-3189-49cf-aed5-6a23465b5056', 'Бытовая техника', CURRENT_TIMESTAMP),
    ('cb905cad-0bd2-42fd-a3da-712ea07d8a8b', 'Спорт и отдых', CURRENT_TIMESTAMP),
    ('fe767d7e-5754-45b3-9a35-7015ff103aee', 'Хобби и развлечения', CURRENT_TIMESTAMP),
    ('4368d269-9710-448c-8d99-0f76b8e4eb30', 'Красота и здоровье', CURRENT_TIMESTAMP);