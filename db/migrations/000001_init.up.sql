-- Создание расширения для генерации UUID, если оно еще не существует
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
        CREATE TYPE user_status AS ENUM ('active', 'inactive', 'banned');
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_method') THEN
        CREATE TYPE payment_method AS ENUM ('cash', 'card');
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'delivery_method') THEN
        CREATE TYPE delivery_method AS ENUM ('pickup', 'delivery');
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'purchase_status') THEN
        CREATE TYPE purchase_status AS ENUM ('pending', 'in_progress', 'completed', 'cancelled');
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cart_status') THEN
        CREATE TYPE cart_status AS ENUM ('active', 'inactive', 'deleted');
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'advert_status') THEN
        CREATE TYPE advert_status AS ENUM ('active', 'inactive', 'reserved');
    END IF;
END $$;

-- Таблица для хранения статических файлов
CREATE TABLE IF NOT EXISTS static (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT
        CONSTRAINT upload_name_length CHECK (LENGTH(name) <= 255) NOT NULL,
    path TEXT
        CONSTRAINT upload_path_length CHECK (LENGTH(path) <= 255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT
        CONSTRAINT username_length CHECK (LENGTH(username) <= 50),
    email TEXT NOT NULL
        CONSTRAINT email_unique UNIQUE,
    password_hash bytea NOT NULL,
    password_salt bytea NOT NULL,
    phone_number TEXT
        CONSTRAINT phone_number_length CHECK (LENGTH(phone_number) <= 20),
    image_id UUID
        CONSTRAINT image_id_fk REFERENCES static(id) ON DELETE SET NULL,
    status user_status DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица продавцов
CREATE TABLE IF NOT EXISTS seller (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL
        CONSTRAINT seller_user_fk REFERENCES "user"(id) ON DELETE CASCADE,
    description TEXT
        CONSTRAINT description_length CHECK (LENGTH(description) <= 1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица подписок
CREATE TABLE IF NOT EXISTS subscription (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL
        CONSTRAINT subscription_user_fk REFERENCES "user"(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL
        CONSTRAINT subscription_seller_fk REFERENCES seller(id) ON DELETE CASCADE,
    CONSTRAINT subscription_unique UNIQUE (user_id, seller_id)
);

-- Таблица категорий
CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    title TEXT
        CONSTRAINT category_title_length CHECK (LENGTH(title) <= 100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица объявлений
CREATE TABLE IF NOT EXISTS advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    title TEXT
        CONSTRAINT advert_title_length CHECK (LENGTH(title) <= 255) NOT NULL,
    description TEXT
        CONSTRAINT advert_description_length CHECK (LENGTH(description) <= 3000) NOT NULL,
    price INTEGER NOT NULL,
    seller_id UUID NOT NULL,
    image_id UUID,
    category_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    location TEXT
        CONSTRAINT advert_location_length CHECK (LENGTH(location) <= 150) NOT NULL,
    has_delivery BOOLEAN NOT NULL,
    status advert_status DEFAULT 'active',
    FOREIGN KEY (seller_id) REFERENCES seller(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES static(id) ON DELETE SET NULL,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE SET NULL
);

-- Таблица сохраненных объявлений
CREATE TABLE IF NOT EXISTS saved_advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    user_id UUID NOT NULL,
    advert_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE,
    CONSTRAINT saved_advert_unique UNIQUE (user_id, advert_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица корзины
CREATE TABLE IF NOT EXISTS cart (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    user_id UUID NOT NULL,
    status cart_status DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);

-- Таблица для связи между корзиной и объявлениями
CREATE TABLE IF NOT EXISTS cart_advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    cart_id UUID NOT NULL,
    advert_id UUID NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE,
    CONSTRAINT cart_advert_unique UNIQUE (cart_id, advert_id)
);

-- Таблица для покупки
CREATE TABLE IF NOT EXISTS purchase (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    cart_id UUID NOT NULL,
    status purchase_status DEFAULT 'pending',
    adress TEXT
        CONSTRAINT adress_length CHECK (LENGTH(adress) <= 150),
    payment_method payment_method DEFAULT 'cash',
    delivery_method delivery_method DEFAULT 'pickup',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE
);

-- Функция для автоматического обновления поля updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггеры для автоматического обновления поля updated_at
CREATE TRIGGER update_advert_updated_at
BEFORE UPDATE ON advert
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_seller_updated_at
BEFORE UPDATE ON seller
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cart_updated_at
BEFORE UPDATE ON cart
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_purchase_updated_at
BEFORE UPDATE ON purchase
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Создание триггерной функции для установки image_id по умолчанию
CREATE OR REPLACE FUNCTION set_default_image_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.image_id IS NULL THEN
        SELECT id INTO NEW.image_id FROM static WHERE name = 'default.jpg' LIMIT 1;
        
        IF NEW.image_id IS NULL THEN
            RAISE EXCEPTION 'Изображение "default.jpg" не найдено в таблице "static". Пожалуйста, добавьте его перед созданием пользователя.';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создание триггера для таблицы "user"
CREATE TRIGGER trg_set_default_image_id
BEFORE INSERT ON "user"
FOR EACH ROW
EXECUTE FUNCTION set_default_image_id();