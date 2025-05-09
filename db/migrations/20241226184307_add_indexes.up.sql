-- Добавление индекса на поле title
CREATE INDEX IF NOT EXISTS idx_advert_title ON advert(title);

-- Добавление индекса на поле status
CREATE INDEX IF NOT EXISTS idx_advert_status ON advert(status);