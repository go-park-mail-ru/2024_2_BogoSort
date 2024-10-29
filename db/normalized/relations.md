### Таблица пользователя
```
Relation user:

{id} -> username, email, phone_number, password_hash, password_salt, image_id, status, created_at, updated_at

```
### Таблица продавца
```
Relation seller:

{id} -> user_id, description, created_at, updated_at

```
### Таблица подписки
```
Relation subscription:

{id} -> user_id, seller_id, created_at, updated_at

```
### Таблица объявлений
```
Relation advert:

{id} -> title, description, price, image_id, category_id, has_delivery, status, seller_id, created_at, updated_at

{seller_id} -> location

```
### Таблица сохраненных объявлений
```
Relation saved_advert:

{id} -> user_id, advert_id, created_at, updated_at

```
### Таблица корзины
```
Relation cart:

{id} ->  user_id, created_at, updated_at

```
### Таблица связи между корзиной и объявлениями
```
Relation cart_advert:

{id} -> cart_id, advert_id, created_at, updated_at

```
### Таблица покупок
```
Relation purchase:

{id} -> cart_id, status, created_at, updated_at

```
### Таблица категорий
```
Relation category:

{id} -> title, created_at, updated_at

```
### Таблица для хранения статических файлов
```
Relation static:

{id} -> name, path, created_at, updated_at

```