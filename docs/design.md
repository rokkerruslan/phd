# Design (Draft)

Этот документ содержит описание сущностей системы и варианты использования
системы пользователями.

## Stages

Шаги необходимые пройти до первого релиза.

1. Проработка дизайн документа.
2. Разработка прототипа.
3. Deployment Infrastructure.
4. Integration with payment systems.
5. Стратегия привлечения пользователей.

X. Альфа релиз. Доступ по ссылкам. Создаём публичные мероприятия, ходим
фотографируем отрабатываем все механизмы.
Y. Бета релиз. Доступ для всех. Берём символическую плату за использования сервиса или
предоставляем бесплатно.
Z. Релиз. 

## System Entities

`Public Event` - событие на которое можно отзываться и публиковать фотографии без
подтверждения владельца события. Эти события могут создавать только админы.

`Private Event` - событие на которое можно откликаться всем, но публиковать
фотографии можно только после одобрения создателя события.

`Event Owner` - создатель события. В публичных событиях играет малую роль.

`Collection` - набор фотографий.

`Offer` - предложение об участии в событии, на публичные события `Offer` принимается
автоматически. На приватные события `Offer` должен быть подтверждён владельцем события.

## Use Cases

Предварительные сценарии взаимодействия с сервисом.

- Пользователь создаёт `Public Event`
    1. Create `Public Event`.
    2. Ждать появления фотографий.
    3. Добавлять фото в коллекцию (платно или бесплатно).

- Пользователь получает фотографии с `Public Event`:
    1. Найти событие.
    2. Добавлять фото в коллекцию (платно или бесплатно).

- Пользователь публикует фотографии к `Public Event`:
    1. Найти событие.
    2. Заявить об участии (make `Offer`)
    3. Загрузить фотографии к событию.
    4. Опубликовать фотографии (есть возможность указания стоимости).

- Пользователь создаёт `Private Event`
    1. Создать `Private Event`.
    2. Дождаться появления `Offer`.
    3. Одобрить один или несколько `Offer`.
    4. Оплатить выбранные `Offer`, деньги переходят на счёт `ph`, не уходят фотографу сразу.
    4. Дождаться загрузки/получения фотографий.
    5. Подтвердить работу.
    6. После подтверждения деньги со счёта переходят на счёта фотографов.

- Пользователь участвует в `Private Event`
    1. Пользователь находит событие.
    2. Создаёт `Offer`.
    3. Дожидается одобрения от `Event Owner`.
    4. Загружает фотографии к событию.
    5. Ждёт подтверждения работы.
    6. Получает деньги на счёт.

## Money flows

Можно получать деньги несколькими путями.

Доходы от приватных событий:

1. Пользователь после одобрения одного или нескольких фотографов
оплачивает стоимость каждого + маржа. В дальнейшем мы переведём
фотографам каждую из сумм. Маржа остаётся у нас.

2. Фотографы оплачивают отклики (не очень схема).

3. Фотографы оплачивают подписку и им разрешено откликаться на определённое кол-во
событий и загружать определённое число фотографий.

Доходы от публичных событий:

4. Фотографы выставляют фотографии по определённым ценам, пользователи могу купить фотографии.
Разницу между покупкой-продажей себе.
