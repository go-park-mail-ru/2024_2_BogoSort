-- Таблица просмотренных объявлений
CREATE TABLE IF NOT EXISTS viewed_advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    user_id UUID NOT NULL,
    advert_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE,
    CONSTRAINT viewed_advert_unique UNIQUE (user_id, advert_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);