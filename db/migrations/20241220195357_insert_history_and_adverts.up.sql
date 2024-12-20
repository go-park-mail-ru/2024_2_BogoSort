-- Заполнение таблицы "static" для новых объявлений
INSERT INTO static (id, name, path, created_at)
VALUES
    (uuid_generate_v4(), 'image31.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image32.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image33.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image34.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image35.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image36.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image37.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image38.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image39.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image40.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image41.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image42.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image43.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image44.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image45.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image46.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image47.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image48.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image49.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image50.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image51.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image52.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image53.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image54.jpg', 'static_files/images/', CURRENT_TIMESTAMP),
    (uuid_generate_v4(), 'image55.jpg', 'static_files/images/', CURRENT_TIMESTAMP);
   

-- Заполнение таблицы "advert" для новых объявлений
INSERT INTO advert (id, title, description, price, seller_id, image_id, category_id, created_at, updated_at, location, has_delivery, status)
VALUES
   (uuid_generate_v4(), 'Горный велосипед', 'Надежный горный велосипед для активного отдыха.', 25000, 
       (SELECT id FROM seller LIMIT 1 OFFSET 0), 
       (SELECT id FROM static WHERE name = 'image31.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Шоссейный велосипед', 'Легкий шоссейный велосипед для длительных поездок.', 30000, 
       (SELECT id FROM seller LIMIT 1 OFFSET 1), 
       (SELECT id FROM static WHERE name = 'image32.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', TRUE, 'active'),

    (uuid_generate_v4(), 'Детский велосипед', 'Яркий детский велосипед для маленьких гонщиков.', 15000, 
       (SELECT id FROM seller LIMIT 1 OFFSET 2), 
       (SELECT id FROM static WHERE name = 'image33.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Электровелосипед', 'Современный электровелосипед для комфортных поездок.', 50000, 
       (SELECT id FROM seller LIMIT 1 OFFSET 3), 
       (SELECT id FROM static WHERE name = 'image34.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', TRUE, 'active'),

    (uuid_generate_v4(), 'Велосипед для трюков', 'Велосипед для выполнения трюков и фристайла.', 20000, 
       (SELECT id FROM seller LIMIT 1 OFFSET 4), 
       (SELECT id FROM static WHERE name = 'image35.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Велосипедный шлем', 'Безопасный шлем для велосипедистов.', 3000, 
       (SELECT id FROM seller LIMIT 1 OFFSET 0), 
       (SELECT id FROM static WHERE name = 'image36.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Спортивный шлем', 'Шлем для защиты во время занятий спортом.', 3500, 
       (SELECT id FROM seller LIMIT 1 OFFSET 1), 
       (SELECT id FROM static WHERE name = 'image37.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', TRUE, 'active'),

    (uuid_generate_v4(), 'Футбольный мяч', 'Качественный футбольный мяч для тренировок.', 2000, 
       (SELECT id FROM seller LIMIT 1 OFFSET 2), 
       (SELECT id FROM static WHERE name = 'image38.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Теннисная ракетка', 'Легкая ракетка для игры в теннис.', 4000, 
       (SELECT id FROM seller LIMIT 1 OFFSET 3), 
       (SELECT id FROM static WHERE name = 'image39.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', TRUE, 'active'),

    (uuid_generate_v4(), 'Йога-мат', 'Удобный мат для занятий йогой.', 1500, 
       (SELECT id FROM seller LIMIT 1 OFFSET 4), 
       (SELECT id FROM static WHERE name = 'image40.jpg'), 
       (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
       CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Волейбольный мяч', 'Качественный волейбольный мяч для тренировок.', 1200, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image41.jpg'), 
        (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Хоккейная клюшка', 'Новая хоккейная клюшка.', 12000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image42.jpg'), 
        (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', TRUE, 'active'),

    (uuid_generate_v4(), 'Коврик для занятий', 'Удобный коврик для занятий йогой.', 1500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image43.jpg'), 
        (SELECT id FROM category WHERE title = 'Спорт и отдых'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    -- Новые товары для категории "Бытовая техника"
    (uuid_generate_v4(), 'Микроволновая печь', 'Компактная микроволновая печь с множеством функций.', 7000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image44.jpg'), 
        (SELECT id FROM category WHERE title = 'Бытовая техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', TRUE, 'active'),

    (uuid_generate_v4(), 'Холодильник', 'Энергосберегающий холодильник с большой вместимостью.', 25000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image45.jpg'), 
        (SELECT id FROM category WHERE title = 'Бытовая техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Стиральная машина', 'Энергоэффективная стиральная машина с большим объемом загрузки.', 20000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image46.jpg'), 
        (SELECT id FROM category WHERE title = 'Бытовая техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    -- Новые товары для категории "Компьютерная техника"
    (uuid_generate_v4(), 'Игровой ноутбук', 'Мощный игровой ноутбук с высокой производительностью.', 60000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image47.jpg'), 
        (SELECT id FROM category WHERE title = 'Компьютерная техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', TRUE, 'active'),

    (uuid_generate_v4(), 'Монитор 27', 'Высококачественный монитор с разрешением 4K.', 22000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image48.jpg'), 
        (SELECT id FROM category WHERE title = 'Компьютерная техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Клавиатура механическая', 'Удобная механическая клавиатура для геймеров.', 5000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image49.jpg'), 
        (SELECT id FROM category WHERE title = 'Компьютерная техника'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', TRUE, 'active'),

    -- Новые товары для категории "Красота и здоровье"
    (uuid_generate_v4(), 'Фен для волос', 'Мощный фен с несколькими режимами.', 3000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image50.jpg'), 
        (SELECT id FROM category WHERE title = 'Красота и здоровье'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active'),

    (uuid_generate_v4(), 'Массажер для спины', 'Эффективный массажер для снятия напряжения.', 3500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 0), 
        (SELECT id FROM static WHERE name = 'image51.jpg'), 
        (SELECT id FROM category WHERE title = 'Красота и здоровье'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Москва', TRUE, 'active'),

    (uuid_generate_v4(), 'Крем для лица', 'Увлажняющий крем для ежедневного использования.', 1500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 1), 
        (SELECT id FROM static WHERE name = 'image52.jpg'), 
        (SELECT id FROM category WHERE title = 'Красота и здоровье'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Санкт-Петербург', TRUE, 'active'),

    -- Новые товары для категории "Хобби и развлечения"
    (uuid_generate_v4(), 'Набор для рисования', 'Полный набор для акварельной живописи.', 4000, 
        (SELECT id FROM seller LIMIT 1 OFFSET 2), 
        (SELECT id FROM static WHERE name = 'image53.jpg'), 
        (SELECT id FROM category WHERE title = 'Хобби и развлечения'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Новосибирск', TRUE, 'active'),

    (uuid_generate_v4(), 'Настольная игра', 'Интересная настольная игра для всей семьи.', 2500, 
        (SELECT id FROM seller LIMIT 1 OFFSET 3), 
        (SELECT id FROM static WHERE name = 'image54.jpg'), 
        (SELECT id FROM category WHERE title = 'Хобби и развлечения'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Екатеринбург', TRUE, 'active'),

    (uuid_generate_v4(), 'Пазлы', 'Развивающая игра для детей.', 1800, 
        (SELECT id FROM seller LIMIT 1 OFFSET 4), 
        (SELECT id FROM static WHERE name = 'image55.jpg'), 
        (SELECT id FROM category WHERE title = 'Хобби и развлечения'), 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Казань', TRUE, 'active');

-- Заполнение таблицы "history"
INSERT INTO price_history (id, advert_id, old_price, new_price, changed_at)
VALUES
    -- Для "Футбольный мяч"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Волейбольный мяч'), 2000, 1800, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Волейбольный мяч'), 1800, 1600, CURRENT_TIMESTAMP),
    -- Для "Теннисная ракетка"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Теннисная ракетка'), 4000, 3500, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Теннисная ракетка'), 3500, 3200, CURRENT_TIMESTAMP),
    -- Для "Йога-мат"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Йога-мат'), 1500, 1400, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Йога-мат'), 1400, 1300, CURRENT_TIMESTAMP),
    -- Для "Микроволновая печь"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Микроволновая печь'), 7000, 6500, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Микроволновая печь'), 6500, 6000, CURRENT_TIMESTAMP),
    -- Для "Холодильник"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Холодильник'), 25000, 24000, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Холодильник'), 24000, 23000, CURRENT_TIMESTAMP),
    -- Для "Стиральная машина"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Стиральная машина'), 20000, 19500, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Стиральная машина'), 19500, 19000, CURRENT_TIMESTAMP),
    -- Для "Игровой ноутбук"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Игровой ноутбук'), 60000, 58000, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Игровой ноутбук'), 58000, 55000, CURRENT_TIMESTAMP),
    -- Для "Монитор 27"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Монитор 27'), 22000, 21000, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Монитор 27'), 21000, 20000, CURRENT_TIMESTAMP),
    -- Для "Клавиатура механическая"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Клавиатура механическая'), 5000, 4800, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Клавиатура механическая'), 4800, 4500, CURRENT_TIMESTAMP),
    -- Для "Фен для волос"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Фен для волос'), 3000, 2800, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Фен для волос'), 2800, 2500, CURRENT_TIMESTAMP),
    -- Для "Массажер для шеи"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Массажер для спины'), 3500, 3300, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Массажер для спины'), 3300, 3000, CURRENT_TIMESTAMP),
    -- Для "Крем для лица"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Крем для лица'), 1500, 1400, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Крем для лица'), 1400, 1300, CURRENT_TIMESTAMP),
    -- Для "Набор для рисования"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Набор для рисования'), 4000, 3800, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Набор для рисования'), 3800, 3600, CURRENT_TIMESTAMP),
    -- Для "Настольная игра"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Настольная игра'), 2500, 2300, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Настольная игра'), 2300, 2100, CURRENT_TIMESTAMP),
    -- Для "Пазлы"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Пазлы'), 1800, 1700, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Пазлы'), 1700, 1600, CURRENT_TIMESTAMP),

    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Горный велосипед'), 25000, 24000, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Горный велосипед'), 24000, 23000, CURRENT_TIMESTAMP),
        -- Для "Шоссейный велосипед"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Шоссейный велосипед'), 30000, 29000, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Шоссейный велосипед'), 29000, 28000, CURRENT_TIMESTAMP),
        -- Для "Детский велосипед"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Детский велосипед'), 15000, 14000, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Детский велосипед'), 14000, 13000, CURRENT_TIMESTAMP),
        -- Для "Электровелосипед"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Электровелосипед'), 50000, 48000, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Электровелосипед'), 48000, 46000, CURRENT_TIMESTAMP),
        -- Для "Велосипед для трюков"
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Велосипед для трюков'), 20000, 19000, CURRENT_TIMESTAMP),
    (uuid_generate_v4(), (SELECT id FROM advert WHERE title = 'Велосипед для трюков'), 19000, 18000, CURRENT_TIMESTAMP);