-- Обновление существующих пользователей
UPDATE "user"
SET username = 'Пользователь'
WHERE username IS NULL;

-- Установка значения по умолчанию для новых пользователей
ALTER TABLE "user"
ALTER COLUMN username SET DEFAULT 'Пользователь';