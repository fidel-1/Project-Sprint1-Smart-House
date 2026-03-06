# Проект "Система умного домма Спринт1"


## Задание 1. Анализ и планирование

### 1. Описание функциональности монолитного приложения

**Управление отоплением:**

- Пользователи могут удалённо включать/выключать отопление в своём доме через веб-интерфейс.

- Система отправляет команды на установленный в доме модуль управления отоплением (датчик/реле) по сети.

- Подключение нового дома/оборудования выполняется специалистом (самостоятельное подключение пользователем отсутствует).

**Мониторинг температуры:**

- Пользователи могут просматривать температуру в доме через веб-интерфейс.
- Система получает данные о температуре синхронно: сервер опрашивает датчик/модуль (pull-модель), не получает события от устройства, не работает в ассинхроном режиме (push-модель).
- Анализ историчности значений, отчеты для пользователей не доступны.

### 2. Анализ архитектуры монолитного приложения

-Язык: Go.

-База данных: PostgreSQL.

-Архитектура: монолит (веб-слой, бизнес-логика и доступ к данным находятся в одном приложении).

-Взаимодействие: синхронное, обработка запросов последовательно, без асинхронных очередей/шины событий.

-Интеграция с устройствами: управление и чтение температуры инициируется сервером, который обращается к устройствам по сети.

-Масштабируемость: Ограничена, так как монолит сложно масштабировать по частям.

### 3. Определение доменов и границы контекстов

To-Be (целевые домены для экосистемы умного дома):

-Каталог и регистрация устройств — самостоятельное подключение устройств, возможность поиска и/или выбора типов устройств из правочника, привязка к дому.
-Возможность интеграции с последующей загрузкой устройств с другими приложениями умный дом сторонных производителей.

-Управление устройствами (Device Control) — единый домен команд/состояний для разных модулей (отопление, свет, ворота, наблюдение).

-Телеметрия и мониторинг (Telemetry & Monitoring) — приём потоков данных от устройств, хранение, алерты.

-Сценарии/автоматизация (Automation / Rules) — “если температура < X, включить отопление”, расписания, триггеры.

-Интеграция партнёрских устройств (Partner Integration) — стандартные протоколы, адаптеры, расширяемость под новые типы устройств.

-(опционально для SaaS) Биллинг/подписки (Billing / Subscriptions) — модули, тарифы, доступ к функциям.

-Модуль отчетности. Просмотр исторических данных подключенных устройств.

### **4. Проблемы монолитного решения**

Слабая масштабируемость по частям: при росте количества устройств/пользователей нельзя отдельно масштабировать контуры телеметрии или управления.

Сложность развития функциональности: добавление новых устройств (свет, ворота, камеры) приводит к росту монолита и усложняет тестирование/релизы.

Релизы рискованные: развертывание/обновление затрагивает всю систему; ошибка в одном модуле влияет на всех пользователей.

Нет самообслуживания: подключение устройств требует выезда специалиста, что не соответствует SaaS-модели.

Синхронная pull-модель телеметрии плохо подходит для больших потоков данных и near-real-time сценариев.



### 5. Визуализация контекста системы — диаграмма С4


```markdown
- Исходник PlantUML: [C4 Context (To-Be)](diagrams/context-to-be.puml)

![C4 Context (To-Be)](diagrams/context-to-be.png)
```

## Задание 2. Проектирование микросервисной архитектуры


**Диаграмма контейнеров (Containers)**

```markdown
- Исходник PlantUML: [C4 Containers](diagrams/containers.puml)

![C4 Containers](diagrams/containers.png)
```

Ниже представлены диаграммы компонентов для ключевых микросервисов системы.

**1. Веб-приложение (Web Application)**
- Исходник PlantUML: [C4 Component - Web Application](diagrams/components/web_app.puml)

![C4 Component - Web Application](diagrams/components/web_app.png)

**2. Реестр устройств (Device Registry)**
- Исходник PlantUML: [C4 Component - Device Registry](diagrams/components/device_registry.puml)

![C4 Component - Device Registry](diagrams/components/device_registry.png)

**3. Управление устройствами (Device Control)**
- Исходник PlantUML: [C4 Component - Device Control](diagrams/components/device_control.puml)

![C4 Component - Device Control](diagrams/components/device_control.png)

**4. Телеметрия и Мониторинг (Telemetry & Monitoring)**
- Исходник PlantUML: [C4 Component - Telemetry & Monitoring](diagrams/components/telemetry_monitoring.puml)

![C4 Component - Telemetry & Monitoring](diagrams/components/telemetry_monitoring.png)

**5. Автоматизация и Правила (Automation / Rules)**
- Исходник PlantUML: [C4 Component - Automation / Rules](diagrams/components/automation_rules.puml)

![C4 Component - Automation / Rules](diagrams/components/automation_rules.png)

**6. Интеграция с партнерами (Partner Integration)**
- Исходник PlantUML: [C4 Component - Partner Integration](diagrams/components/partner_integration.puml)

![C4 Component - Partner Integration](diagrams/components/partner_integration.png)

**7. Модуль отчетности (Reporting Module)**
- Исходник PlantUML: [C4 Component - Reporting Module](diagrams/components/reporting_module.puml)

![C4 Component - Reporting Module](diagrams/components/reporting_module.png)


**Диаграмма кода (Code)**

- Исходник PlantUML: [C4 Code - Web App Device Controller](diagrams/code/web_app_code.puml)

![C4 Code - Web App Device Controller](diagrams/code/web_app_code.png)

## Задание 3. Разработка ER-диаграммы

- Исходник PlantUML: [ER Diagram](diagrams/er/er_diagram.puml)

![ER Diagram](diagrams/er/er_diagram.png)

## Задание 4. Создание и документирование API

### 1. Тип API
В системе используется гибридная модель API.
REST API применяется для синхронных операций управления устройствами и получения данных.
AsyncAPI используется для обработки потоковых событий телеметрии и автоматизации

### 2. Документация API


- **REST API (OpenAPI)**: [openapi.yaml](api/openapi.yaml)
- **AsyncAPI**: [asyncapi.yaml](api/asyncapi.yaml)

# Задание 5. Работа с docker и docker-compose

### Проверка работы приложения

Для проверки использована коллекция Postman `smarthome-api.postman_collection.json`.

Выполнены запросы:

- `Create Sensor`
- `Get All Sensors`

Результат:
- после запуска контейнеров `temperature-api`, `smarthome-app` и `postgres` приложение корректно работает;
- `temperature-api` возвращает случайное значение температуры;
- при вызове `Get All Sensors` отображаются разные значения температуры для сенсоров, что подтверждает корректную интеграцию `smart_home` с `temperature-api`.

Используемые технологии:

-Go
-PostgreSQL
-Docker
-Docker Compose
-PlantUML
-C4 Model
-OpenAPI
-AsyncAPI