# Нормализация

### Таблица пользователя
Таблица, описывающая параметры сущности пользователя
```
Relation user:

{id} -> username, email, phone_number, password_hash, password_salt, image_id, status, created_at, updated_at
```

### Таблица продавца
Таблица, описывающая продавца
```
Relation seller:

{id, user_id} -> description, created_at, updated_at

```

### Таблица подписки
Таблица, связывающая покупателя и продавца
```
Relation subscription:

{id, user_id, seller_id} -> created_at, updated_at

```

### Таблица объявлений
Таблица, определяющая параметры объявления
```
Relation advert:

{id} -> title, description, price, image_id, category_id, has_delivery, location, status, seller_id, created_at, updated_at

```

### Таблица сохраненных объявлений
Таблица объявления, добавленного в избранное пользователем
```
Relation saved_advert:

{id, user_id, advert_id} -> created_at, updated_at

```

### Таблица корзины
Таблица, создающая корзину для пользователя
```
Relation cart:

{id, user_id} -> created_at, updated_at

```

### Таблица связи между корзиной и объявлениями
Таблица, связывающая корзину с объявлениями
```
Relation cart_advert:

{id, cart_id, advert_id} -> created_at, updated_at

```

### Таблица покупок
Таблица, описывающая статус и содержимое покупок
```
Relation purchase:

{id, cart_id} -> status, created_at, updated_at

```

### Таблица категорий
Таблица, описывающая категории объявлений
```
Relation category:

{id} -> title, created_at, updated_at

```

### Таблица для хранения статических файлов
```
Relation static:

{id} -> name, path, created_at, updated_at

```