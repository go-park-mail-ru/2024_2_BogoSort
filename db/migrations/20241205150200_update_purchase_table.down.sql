ALTER TABLE purchase
ADD COLUMN cart_id UUID;

UPDATE purchase p
SET cart_id = (
    SELECT ca.cart_id
    FROM cart_advert ca
    JOIN purchase_advert pa ON pa.advert_id = ca.advert_id
    WHERE pa.purchase_id = p.id
    LIMIT 1
);

ALTER TABLE purchase
ADD CONSTRAINT purchase_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE;

ALTER TABLE purchase
DROP CONSTRAINT purchase_seller_fk,
DROP COLUMN seller_id;

DROP TABLE IF EXISTS purchase_advert;
