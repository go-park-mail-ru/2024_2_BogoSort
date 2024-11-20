DROP TABLE IF EXISTS viewed_advert CASCADE;
DROP TRIGGER IF EXISTS trg_set_advert_default_image_id ON advert;
DROP FUNCTION IF EXISTS set_advert_default_image_id;
