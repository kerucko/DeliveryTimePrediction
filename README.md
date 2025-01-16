# DeliveryTimePrediction - сервис предсказания времени доставки в зависимости от условий

# Авторы

| ФИО               | Группа          |
|--------------------|-----------------|
| Катаев Юрий       | М80-410Б-21     |
| Куликова Анастасия | М80-410Б-21     |
| Цыкин Павел       | М80-410Б-21     |


# Презентация проекта

https://docs.google.com/presentation/d/1p5WcG8X9Cig24ZBOzoyLY_eVTPVThjmrjyxVMkeGS_c/edit#slide=id.p


# Отчет 

## Бизнес-цель

Бизнес-цель данного проекта заключается в разработке системы прогнозирования времени доставки на основе ключевых характеристик заказа, маршрута и внешних факторов. Продукт позволит повысить точность прогнозов, автоматизировать процесс расчета, сократить затраты на ручной анализ и улучшить клиентский опыт за счет предоставления более точной информации о времени получения заказов.

## ML-цель

## Обучение модели

Для обучения использовался датасет [Food Delivery time Analysis and Prediction](https://www.kaggle.com/code/a3amat02/food-delivery-time-analysis-and-prediction)

код представлен в файле [training.ipynb](training.ipynb)

метрики на тествой выборке
* MAE: 6.48
* MSE: 87.36
* RMSE: 9.35

## Архитектура

![архитектура](docs/architecture.png)


## Стек

Проект состоит из двух микросервисов, Kafka и Postgres БД:
 - backend-gateway - сервис на Go, отвечающий за взаимодействие с клиентами и доступ к БД. Пользователь может отправить запрос на обработку и проверить готовность запроса по id.
 - predictor - сервис на Python с использованием FastAPI. Прогнозирует время доставки с помощью предобученной модели.
 - Kafka - брокер сообщений, который используется для связи микросервисов с двумя топиками: tasks и completed. Благодаря ей достигается асинхронность в обработке запросов.
 - PostgreSQL - БД для сохранения результатов. Клиент может получить результат запроса по id.



