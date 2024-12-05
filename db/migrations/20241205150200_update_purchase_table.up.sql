ALTER TABLE purchase
ADD COLUMN seller_id UUID;

UPDATE purchase p
SET seller_id = (
    SELECT DISTINCT a.seller_id
    FROM cart_advert ca
    JOIN advert a ON ca.advert_id = a.id
    WHERE ca.cart_id = p.cart_id
    LIMIT 1
);

CREATE TABLE IF NOT EXISTS purchase_advert (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    purchase_id UUID NOT NULL,
    advert_id UUID NOT NULL,
    FOREIGN KEY (purchase_id) REFERENCES purchase(id) ON DELETE CASCADE,
    FOREIGN KEY (advert_id) REFERENCES advert(id) ON DELETE CASCADE,
    CONSTRAINT purchase_advert_unique UNIQUE (purchase_id, advert_id)
);

INSERT INTO purchase_advert (id, purchase_id, advert_id)
SELECT 
    uuid_generate_v4(),
    p.id,
    ca.advert_id
FROM purchase p
JOIN cart_advert ca ON ca.cart_id = p.cart_id;

ALTER TABLE purchase
ALTER COLUMN seller_id SET NOT NULL;

ALTER TABLE purchase
ADD CONSTRAINT purchase_seller_fk FOREIGN KEY (seller_id) REFERENCES seller(id) ON DELETE CASCADE;

ALTER TABLE purchase
DROP CONSTRAINT purchase_cart_id_fkey,
DROP COLUMN cart_id;