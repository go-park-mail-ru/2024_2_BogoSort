-- Обновление существующих пользователей
UPDATE "user"
SET username = NULL
WHERE username LIKE 'Пользователь%';  -- Удаление пользователей, начинающихся с 'Пользователь'

-- Установка значения по умолчанию для новых пользователей обратно на NULL
ALTER TABLE "user"
ALTER COLUMN username SET DEFAULT NULL;
