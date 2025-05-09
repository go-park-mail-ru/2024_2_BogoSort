-- Сохранение текущего DDL базы данных
\copy (SELECT table_name FROM information_schema.tables WHERE table_schema='public') TO 'tables_list.txt'

-- Экспорт всех таблиц
pg_dump -U postgres -s -f init_ddl.sql your_database_name
