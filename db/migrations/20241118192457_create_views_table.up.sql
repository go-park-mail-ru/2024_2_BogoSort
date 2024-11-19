-- Таблица просмотренных объявлений
CREATE TABLE IF NOT EXISTS viewed_advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    user_id UUID NULL,
    advert_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание триггерной функции для установки image_id по умолчанию
CREATE OR REPLACE FUNCTION set_default_image_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.image_id IS NULL THEN
        SELECT id INTO NEW.image_id FROM static WHERE name = 'default_advert.jpg' LIMIT 1;
        
        IF NEW.image_id IS NULL THEN
            RAISE EXCEPTION 'Изображение "default_advert.jpg" не найдено в таблице "static". Пожалуйста, добавьте его перед созданием объявления.';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создание триггера для таблицы "user"
CREATE TRIGGER trg_set_default_image_id
BEFORE INSERT ON advert
FOR EACH ROW
EXECUTE FUNCTION set_default_image_id();