DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'update_advert_updated_at'
    ) THEN
        CREATE TRIGGER update_advert_updated_at
        BEFORE UPDATE ON advert
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

DROP TABLE IF EXISTS viewed_advert CASCADE;
DROP TRIGGER IF EXISTS trg_set_advert_default_image_id ON advert;
DROP FUNCTION IF EXISTS set_advert_default_image_id;
