CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица категорий
CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    title TEXT 
        CONSTRAINT category_title_length CHECK (LENGTH(title) <= 100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для хранения статических файлов
CREATE TABLE IF NOT EXISTS static (
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT 
        CONSTRAINT upload_name_length CHECK (LENGTH(name) <= 255) NOT NULL,
    path TEXT 
        CONSTRAINT upload_path_length CHECK (LENGTH(path) <= 255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
    status TEXT NOT NULL 
        CONSTRAINT advert_status_length CHECK (LENGTH(status) <= 100) NOT NULL,
    FOREIGN KEY (seller_id) REFERENCES seller(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES static(id) ON DELETE SET NULL,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE SET NULL
);

-- Таблица сохраненных объявлений
CREATE TABLE IF NOT EXISTS saved_advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    user_id UUID NOT NULL, 
    advert_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE
);

-- Таблица корзины
CREATE TABLE IF NOT EXISTS cart (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    user_id UUID NOT NULL, 
    advert_id UUID NOT NULL, 
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE
);

-- Таблица для связи между корзиной и объявлениями
CREATE TABLE IF NOT EXISTS cart_advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    cart_id UUID NOT NULL, 
    advert_id UUID NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE
);

-- Таблица для покупки
CREATE TABLE IF NOT EXISTS purchase (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    cart_id UUID NOT NULL,
    status TEXT NOT NULL 
        CONSTRAINT status_length CHECK (LENGTH(status) <= 255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE
);