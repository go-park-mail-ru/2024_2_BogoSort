-- Заполнение таблицы "static"
INSERT INTO static (id, name, path, created_at)
VALUES
    (uuid_generate_v4(), 'default.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image1.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image2.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image3.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image4.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image5.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image6.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image7.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image8.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image9.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image10.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image11.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image12.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image13.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image14.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image15.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image16.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image17.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image18.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image19.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image20.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image21.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image22.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image23.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image24.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image25.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image26.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image27.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image28.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image29.jpg', 'static/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image30.jpg', 'static/', CURRENT_TIMESTAMP);

-- Заполнение таблицы "user"
INSERT INTO "user" (id, username, email, password_hash, password_salt, phone_number, image_id, status, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'ivan_petrov', 'ivan.petrov@example.com', 'hash1', 'salt1', '+79261234567', (SELECT id FROM static WHERE name = 'default.jpg'), 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'anna_smirnova', 'anna.smirnova@example.com', 'hash2', 'salt2', '+79269876543', (SELECT id FROM static WHERE name = 'default.jpg'), 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'pavel_ivanov', 'pavel.ivanov@example.com', 'hash3', 'salt3', '+79261239876', (SELECT id FROM static WHERE name = 'default.jpg'), 'inactive', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'elena_kuznetsova', 'elena.kuznetsova@example.com', 'hash4', 'salt4', '+79261231234', (SELECT id FROM static WHERE name = 'default.jpg'), 'banned', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'dmitry_sokolov', 'dmitry.sokolov@example.com', 'hash5', 'salt5', '+79261234568', (SELECT id FROM static WHERE name = 'default.jpg'), 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

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

-- Заполнение таблицы "advert"
INSERT INTO advert (id, title, description, price, seller_id, image_id, category_id, created_at, updated_at, location, has_delivery, status)
VALUES
    (uuid_generate_v4(), 'Элегантное вечернее платье', 'Стильное платье из натуральных тканей, идеально подходит для вечерних мероприятий.', 7500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image1.jpg'), 
        (SELECT id FROM category WHERE title = 'Женский гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),
        
    (uuid_generate_v4(), 'Мужская кожаная куртка', 'Высококачественная кожаная куртка для стильных мужчин.', 15000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image2.jpg'), 
        (SELECT id FROM category WHERE title = 'Мужской гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', FALSE, 'active'),

    (uuid_generate_v4(), 'Детский комбинезон', 'Удобный и теплый комбинезон для малышей.', 3000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image3.jpg'), 
        (SELECT id FROM category WHERE title = 'Детский гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Игрушечный набор конструктор', 'Развивающий конструктор для детей от 3 лет.', 2500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image4.jpg'), 
        (SELECT id FROM category WHERE title = 'Детские товары'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', FALSE, 'active'),

    (uuid_generate_v4(), 'Лазерный уровень', 'Профессиональный лазерный уровень для строительных работ.', 5000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image5.jpg'), 
        (SELECT id FROM category WHERE title = 'Стройматериалы и инструменты'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Игровой ноутбук ASUS', 'Мощный игровой ноутбук с высокой производительностью.', 60000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image6.jpg'), 
        (SELECT id FROM category WHERE title = 'Компьютерная техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Многофункциональный гриль', 'Гриль для дачи с несколькими функциями приготовления.', 8000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image7.jpg'), 
        (SELECT id FROM category WHERE title = 'Для дома и дачи'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', FALSE, 'active'),

    (uuid_generate_v4(), 'Стиральная машина Bosch', 'Энергоэффективная стиральная машина с большим объемом загрузки.', 20000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image8.jpg'), 
        (SELECT id FROM category WHERE title = 'Бытовая техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Фитнес-браслет Xiaomi', 'Смарт-браслет с множеством функций для здоровья и спорта.', 3500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image9.jpg'), 
        (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', FALSE, 'active'),

    (uuid_generate_v4(), 'Акварельный набор', 'Полный набор для акварельной живописи.', 4000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image10.jpg'), 
        (SELECT id FROM category WHERE title = 'Хобби и развлечения'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Увлажнитель воздуха', 'Компактный увлажнитель для дома или офиса.', 2500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image11.jpg'), 
        (SELECT id FROM category WHERE title = 'Красота и здоровье'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Платье повседневное', 'Комфортное платье для ежедневной носки.', 4500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image12.jpg'), 
        (SELECT id FROM category WHERE title = 'Женский гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', FALSE, 'active'),

    (uuid_generate_v4(), 'Мужские джинсы Levis', 'Классические джинсы для мужчин.', 5000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image13.jpg'), 
        (SELECT id FROM category WHERE title = 'Мужской гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Детская футболка с принтом', 'Яркая футболка для детей всех возрастов.', 1500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image14.jpg'), 
        (SELECT id FROM category WHERE title = 'Детский гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', FALSE, 'active'),

    (uuid_generate_v4(), 'Набор для рисования малышам', 'Безопасные краски и кисти для маленьких художников.', 2000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image15.jpg'), 
        (SELECT id FROM category WHERE title = 'Детские товары'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Шуруповерт Makita', 'Профессиональный шуруповерт для строительных работ.', 7000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image16.jpg'), 
        (SELECT id FROM category WHERE title = 'Стройматериалы и инструменты'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Рабочий стол компьютерный', 'Прочный и удобный рабочий стол для дома или офиса.', 12000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image17.jpg'), 
        (SELECT id FROM category WHERE title = 'Компьютерная техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', FALSE, 'active'),

    (uuid_generate_v4(), 'Комплект мебели для дачи', 'Стол и стулья из устойчивых к погоде материалов.', 9000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image18.jpg'), 
        (SELECT id FROM category WHERE title = 'Для дома и дачи'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Холодильник Samsung', 'Энергосберегающий холодильник с большой вместимостью.', 25000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image19.jpg'), 
        (SELECT id FROM category WHERE title = 'Бытовая техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', FALSE, 'active'),

    (uuid_generate_v4(), 'Велотренажер домашний', 'Удобный велотренажер для поддержания формы дома.', 8000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image20.jpg'), 
        (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Набор акварели для начинающих', 'Все необходимое для первых шагов в акварели.', 3000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image21.jpg'), 
        (SELECT id FROM category WHERE title = 'Хобби и развлечения'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Массажер для шеи', 'Эффективный массажер для снятия напряжения.', 3500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image22.jpg'), 
        (SELECT id FROM category WHERE title = 'Красота и здоровье'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', FALSE, 'active'),

    (uuid_generate_v4(), 'Вечернее платье с шлейфом', 'Элегантное платье с длинным шлейфом для особых случаев.', 10000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image23.jpg'), 
        (SELECT id FROM category WHERE title = 'Женский гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Мужские спортивные брюки', 'Комфортные брюки для занятий спортом и активного отдыха.', 4000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image24.jpg'), 
        (SELECT id FROM category WHERE title = 'Мужской гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', FALSE, 'active'),

    (uuid_generate_v4(), 'Детские сапоги зимние', 'Теплые и водонепроницаемые сапоги для детей.', 2500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image25.jpg'), 
        (SELECT id FROM category WHERE title = 'Детский гардероб'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Развивающая игра "Пазлы"', 'Набор пазлов для развития моторики и логики у детей.', 1800, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image26.jpg'), 
        (SELECT id FROM category WHERE title = 'Детские товары'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Электродрель Bosch', 'Мощная электродрель для профессиональных и домашних работ.', 12000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image27.jpg'), 
        (SELECT id FROM category WHERE title = 'Стройматериалы и инструменты'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', FALSE, 'active'),

    (uuid_generate_v4(), 'Игровой монитор LG 27"', 'Высококачественный монитор с разрешением 4K для геймеров и профессионалов.', 22000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image28.jpg'), 
        (SELECT id FROM category WHERE title = 'Компьютерная техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Кресло офисное Ergonomic', 'Удобное офисное кресло с поддержкой спины.', 8000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image29.jpg'), 
        (SELECT id FROM category WHERE title = 'Для дома и дачи'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', FALSE, 'active'),

    (uuid_generate_v4(), 'Микроволновая печь Panasonic', 'Компактная микроволновая печь с множеством функций.', 7000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image30.jpg'), 
        (SELECT id FROM category WHERE title = 'Бытовая техника'), 
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
    (uuid_generate_v4(), 
        (SELECT id FROM "user" WHERE username = 'ivan_petrov'), 
        (SELECT id FROM advert WHERE title = 'Элегантное вечернее платье'), 
        CURRENT_TIMESTAMP),

    (uuid_generate_v4(), 
        (SELECT id FROM "user" WHERE username = 'anna_smirnova'), 
        (SELECT id FROM advert WHERE title = 'Многофункциональный гриль'), 
        CURRENT_TIMESTAMP),

    (uuid_generate_v4(), 
        (SELECT id FROM "user" WHERE username = 'pavel_ivanov'), 
        (SELECT id FROM advert WHERE title = 'Мужские джинсы Levis'), 
        CURRENT_TIMESTAMP),

    (uuid_generate_v4(), 
        (SELECT id FROM "user" WHERE username = 'elena_kuznetsova'), 
        (SELECT id FROM advert WHERE title = 'Холодильник Samsung'), 
        CURRENT_TIMESTAMP),

    (uuid_generate_v4(), 
        (SELECT id FROM "user" WHERE username = 'dmitry_sokolov'), 
        (SELECT id FROM advert WHERE title = 'Детские сапоги зимние'), 
        CURRENT_TIMESTAMP);

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
        (SELECT id FROM advert WHERE title = 'Элегантное вечернее платье')),
    (uuid_generate_v4(),
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'anna_smirnova')),
        (SELECT id FROM advert WHERE title = 'Многофункциональный гриль')),
    (uuid_generate_v4(),
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'pavel_ivanov')),
        (SELECT id FROM advert WHERE title = 'Мужские джинсы Levis')),
    (uuid_generate_v4(),
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'elena_kuznetsova')),
        (SELECT id FROM advert WHERE title = 'Холодильник Samsung')),
    (uuid_generate_v4(),
        (SELECT id FROM cart WHERE user_id = (SELECT id FROM "user" WHERE username = 'dmitry_sokolov')),
        (SELECT id FROM advert WHERE title = 'Детские сапоги зимние'));

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
