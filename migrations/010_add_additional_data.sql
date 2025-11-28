-- Migration: 010_add_additional_data.sql

-- 1. Добавляем 3 дополнительных склада
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM warehouse WHERE name = 'Склад "Северный"') THEN
        INSERT INTO warehouse (name, address, contact_person, contact_phone, total_capacity, current_usage) VALUES
        ('Склад "Северный"', 'ул. Северная, 25, Москва', 'Петров Сергей', '+7-999-222-33-44', 3000, 0);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM warehouse WHERE name = 'Склад "Западный"') THEN
        INSERT INTO warehouse (name, address, contact_person, contact_phone, total_capacity, current_usage) VALUES
        ('Склад "Западный"', 'ул. Западная, 10, Москва', 'Козлова Ольга', '+7-999-333-44-55', 4000, 0);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM warehouse WHERE name = 'Склад "Центральный"') THEN
        INSERT INTO warehouse (name, address, contact_person, contact_phone, total_capacity, current_usage) VALUES
        ('Склад "Центральный"', 'ул. Центральная, 5, Москва', 'Николаев Дмитрий', '+7-999-444-55-66', 6000, 0);
    END IF;
END $$;

-- 2. Добавляем 5 курьеров как пользователей с ролью courier
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'courier1') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('courier1', 'courier1@testsystem.ru', 'courier', 1, '$2y$12$courier1passwordhash', 'Иванов Алексей', 'Курьерская служба', 'Курьер', '+7-999-111-11-11');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'courier2') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('courier2', 'courier2@testsystem.ru', 'courier', 1, '$2y$12$courier2passwordhash', 'Петров Михаил', 'Курьерская служба', 'Курьер', '+7-999-222-22-22');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'courier3') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('courier3', 'courier3@testsystem.ru', 'courier', 1, '$2y$12$courier3passwordhash', 'Сидорова Анна', 'Курьерская служба', 'Курьер', '+7-999-333-33-33');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'courier4') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('courier4', 'courier4@testsystem.ru', 'courier', 1, '$2y$12$courier4passwordhash', 'Кузнецов Денис', 'Курьерская служба', 'Курьер', '+7-999-444-44-44');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'courier5') THEN
        INSERT INTO users (username, email, userrole, status, password, fullusername, companyname, companyrole, phone) VALUES
        ('courier5', 'courier5@testsystem.ru', 'courier', 1, '$2y$12$courier5passwordhash', 'Морозова Екатерина', 'Курьерская служба', 'Курьер', '+7-999-555-55-55');
    END IF;
END $$;

-- 3. Добавляем 40 локаций
DO $$
DECLARE
    location_record RECORD;
BEGIN
    FOR location_record IN (
        SELECT * FROM (VALUES
            ('ТЦ "Авиапарк"', 'Ходынский бульвар, 4, Москва', 'Семенов А.В.', '+7-999-100-01-01', 25000.00, 5),
            ('ТРЦ "РИО"', 'Дмитровское ш., 163А, Москва', 'Калинина М.С.', '+7-999-100-01-02', 18000.00, 10),
            ('ТЦ "МЕГА"', 'Коровинское ш., 10, Москва', 'Орлов Д.Н.', '+7-999-100-01-03', 22000.00, 15),
            ('ТЦ "Весна"', 'ул. Пушкина, 35, Москва', 'Волкова Т.П.', '+7-999-100-01-04', 12000.00, 20),
            ('ТРК "Континент"', 'ул. Ленина, 87, Москва', 'Жуков Р.А.', '+7-999-100-01-05', 19000.00, 25),
            ('ТЦ "Октябрь"', 'пр. Мира, 124, Москва', 'Лебедева О.И.', '+7-999-100-01-06', 15000.00, 12),
            ('ТРЦ "Галерея"', 'ул. Советская, 56, Москва', 'Новиков С.М.', '+7-999-100-01-07', 28000.00, 8),
            ('ТЦ "Парус"', 'ул. Гагарина, 23, Москва', 'Федорова Е.В.', '+7-999-100-01-08', 16000.00, 18),
            ('ТРК "Москва"', 'ул. Тверская, 15, Москва', 'Дмитриев П.К.', '+7-999-100-01-09', 32000.00, 7),
            ('ТЦ "Заря"', 'ул. Кирова, 42, Москва', 'Соколова А.М.', '+7-999-100-01-10', 14000.00, 22),
            ('ТРЦ "Небо"', 'пр. Победы, 67, Москва', 'Комаров В.С.', '+7-999-100-01-11', 26000.00, 9),
            ('ТЦ "Восток"', 'ул. Садовая, 18, Москва', 'Егорова Л.Д.', '+7-999-100-01-12', 13000.00, 17),
            ('ТРК "Юг"', 'ул. Центральная, 91, Москва', 'Григорьев И.А.', '+7-999-100-01-13', 17000.00, 14),
            ('ТЦ "Север"', 'ул. Лесная, 29, Москва', 'Тихонова М.В.', '+7-999-100-01-14', 14500.00, 19),
            ('ТРЦ "Запад"', 'ул. Школьная, 54, Москва', 'Фролов А.Н.', '+7-999-100-01-15', 21000.00, 11),
            ('ТЦ "Лукоморье"', 'ул. Парковая, 33, Москва', 'Мартынова С.П.', '+7-999-100-01-16', 12500.00, 23),
            ('ТРК "Атриум"', 'ул. Новая, 76, Москва', 'Белов К.Д.', '+7-999-100-01-17', 24000.00, 6),
            ('ТЦ "Радуга"', 'ул. Строителей, 48, Москва', 'Крылова Т.С.', '+7-999-100-01-18', 15500.00, 16),
            ('ТРЦ "Планета"', 'ул. Мира, 112, Москва', 'Сорокин М.А.', '+7-999-100-01-19', 19500.00, 13),
            ('ТЦ "Орион"', 'ул. Звездная, 25, Москва', 'Воронова Е.Н.', '+7-999-100-01-20', 13500.00, 21),
            ('ТРК "Глобус"', 'ул. Интернациональная, 39, Москва', 'Лазарев Д.В.', '+7-999-100-02-01', 22500.00, 8),
            ('ТЦ "Феникс"', 'ул. Возрождения, 17, Москва', 'Медведева О.С.', '+7-999-100-02-02', 16500.00, 15),
            ('ТРЦ "Высота"', 'ул. Горная, 63, Москва', 'Савельев Р.П.', '+7-999-100-02-03', 27500.00, 5),
            ('ТЦ "Волна"', 'ул. Речная, 28, Москва', 'Гусева А.К.', '+7-999-100-02-04', 14200.00, 20),
            ('ТРК "Энергия"', 'ул. Энергетиков, 45, Москва', 'Тарасов В.М.', '+7-999-100-02-05', 18800.00, 12),
            ('ТЦ "Спутник"', 'ул. Космонавтов, 31, Москва', 'Комарова Л.В.', '+7-999-100-02-06', 15200.00, 18),
            ('ТРЦ "Меридиан"', 'ул. Параллельная, 52, Москва', 'Ефимов С.Н.', '+7-999-100-02-07', 23200.00, 9),
            ('ТЦ "Альфа"', 'ул. Бета, 27, Москва', 'Одинцова М.Д.', '+7-999-100-02-08', 13800.00, 22),
            ('ТРК "Омега"', 'ул. Гамма, 34, Москва', 'Власов П.С.', '+7-999-100-02-09', 20200.00, 11),
            ('ТЦ "Кварц"', 'ул. Гранитная, 41, Москва', 'Маслова Т.А.', '+7-999-100-02-10', 14800.00, 19),
            ('ТРЦ "Кристалл"', 'ул. Алмазная, 58, Москва', 'Исаев А.В.', '+7-999-100-02-11', 24200.00, 7),
            ('ТЦ "Рубин"', 'ул. Сапфировая, 22, Москва', 'Суханова Е.П.', '+7-999-100-02-12', 12800.00, 24),
            ('ТРК "Изумруд"', 'ул. Изумрудная, 47, Москва', 'Горбунов Д.М.', '+7-999-100-02-13', 19200.00, 14),
            ('ТЦ "Бриллиант"', 'ул. Бриллиантовая, 36, Москва', 'Зайцева С.В.', '+7-999-100-02-14', 26200.00, 6),
            ('ТРЦ "Платина"', 'ул. Платиновая, 29, Москва', 'Семенов К.А.', '+7-999-100-02-15', 17200.00, 16),
            ('ТЦ "Золото"', 'ул. Золотая, 44, Москва', 'Кузнецова Н.С.', '+7-999-100-02-16', 18200.00, 13),
            ('ТРК "Серебро"', 'ул. Серебряная, 51, Москва', 'Виноградов М.П.', '+7-999-100-02-17', 21200.00, 10),
            ('ТЦ "Бронза"', 'ул. Бронзовая, 38, Москва', 'Давыдова О.Н.', '+7-999-100-02-18', 13200.00, 21),
            ('ТРК "Металл"', 'ул. Металлистов, 55, Москва', 'Журавлев С.Д.', '+7-999-100-02-19', 19800.00, 17),
            ('ТЦ "Сталь"', 'ул. Сталеваров, 42, Москва', 'Носова Т.К.', '+7-999-100-02-20', 14200.00, 23)
        ) AS loc(name, address, contact_person, contact_phone, monthly_rent, rent_due_day)
    ) LOOP
        IF NOT EXISTS (SELECT 1 FROM locations WHERE name = location_record.name) THEN
            INSERT INTO locations (name, address, contact_person, contact_phone, monthly_rent, rent_due_day) 
            VALUES (location_record.name, location_record.address, location_record.contact_person, location_record.contact_phone, location_record.monthly_rent, location_record.rent_due_day);
        END IF;
    END LOOP;
END $$;

-- 4. Распределяем товары по складам в реальной пропорции (по 1 единице каждого товара)
DO $$
DECLARE
    main_warehouse_id BIGINT;
    north_warehouse_id BIGINT;
    west_warehouse_id BIGINT;
    central_warehouse_id BIGINT;
    cat_machines BIGINT;
    cat_toys BIGINT;
    cat_capsules BIGINT;
BEGIN
    -- Получаем ID складов
    SELECT id INTO main_warehouse_id FROM warehouse WHERE name = 'Основной склад' LIMIT 1;
    SELECT id INTO north_warehouse_id FROM warehouse WHERE name = 'Склад "Северный"' LIMIT 1;
    SELECT id INTO west_warehouse_id FROM warehouse WHERE name = 'Склад "Западный"' LIMIT 1;
    SELECT id INTO central_warehouse_id FROM warehouse WHERE name = 'Склад "Центральный"' LIMIT 1;
    
    -- Получаем ID категорий
    SELECT id INTO cat_machines FROM warehouse_categories WHERE name = 'Вендинговые автоматы' LIMIT 1;
    SELECT id INTO cat_toys FROM warehouse_categories WHERE name = 'Игрушки' LIMIT 1;
    SELECT id INTO cat_capsules FROM warehouse_categories WHERE name = 'Капсулы' LIMIT 1;

    -- Основной склад: по 1 единице каждого товара
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM3000-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 3000', 'Вендинговый автомат премиум-класса', 1, 1, 5, 50000.00, 'VM-TM3000-MAIN');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM2000-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 2000', 'Вендинговый автомат стандарт-класса', 1, 1, 5, 35000.00, 'VM-TM2000-MAIN');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM1000-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 1000', 'Компактный вендинговый автомат', 1, 1, 5, 25000.00, 'VM-TM1000-MAIN');
    END IF;

    -- Игрушки
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-SOFT-10-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_toys, 'toy', 'Мягкие игрушки (набор)', 'Набор из 10 мягких игрушек', 1, 5, 50, 150.00, 'TOY-SOFT-10-MAIN');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-HEROES-1-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_toys, 'toy', 'Фигурки супергероев', 'Коллекционные фигурки', 1, 5, 50, 200.00, 'TOY-HEROES-1-MAIN');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-CARS-5-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_toys, 'toy', 'Машинки миниатюрные', 'Набор миниатюрных машинок', 1, 5, 50, 120.00, 'TOY-CARS-5-MAIN');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-CONSTRUCT-1-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_toys, 'toy', 'Конструктор мини', 'Мини-конструктор', 1, 5, 50, 180.00, 'TOY-CONSTRUCT-1-MAIN');
    END IF;

    -- Капсулы
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-STD-100-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_capsules, 'capsule', 'Капсулы стандартные', 'Стандартные прозрачные капсулы', 1, 5, 50, 300.00, 'CAP-STD-100-MAIN');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-COLOR-100-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_capsules, 'capsule', 'Капсулы цветные', 'Набор цветных капсул', 1, 5, 50, 350.00, 'CAP-COLOR-100-MAIN');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-GOLD-50-MAIN') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (main_warehouse_id, cat_capsules, 'capsule', 'Капсулы премиум', 'Премиум капсулы золотого цвета', 1, 5, 50, 500.00, 'CAP-GOLD-50-MAIN');
    END IF;

    -- Склад "Центральный": по 1 единице каждого товара
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM3000-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 3000', 'Вендинговый автомат премиум-класса', 1, 1, 5, 50000.00, 'VM-TM3000-CENTRAL');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM2000-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 2000', 'Вендинговый автомат стандарт-класса', 1, 1, 5, 35000.00, 'VM-TM2000-CENTRAL');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM1000-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 1000', 'Компактный вендинговый автомат', 1, 1, 5, 25000.00, 'VM-TM1000-CENTRAL');
    END IF;

    -- Игрушки
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-SOFT-10-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_toys, 'toy', 'Мягкие игрушки (набор)', 'Набор из 10 мягких игрушек', 1, 5, 50, 150.00, 'TOY-SOFT-10-CENTRAL');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-HEROES-1-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_toys, 'toy', 'Фигурки супергероев', 'Коллекционные фигурки', 1, 5, 50, 200.00, 'TOY-HEROES-1-CENTRAL');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-CARS-5-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_toys, 'toy', 'Машинки миниатюрные', 'Набор миниатюрных машинок', 1, 5, 50, 120.00, 'TOY-CARS-5-CENTRAL');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-CONSTRUCT-1-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_toys, 'toy', 'Конструктор мини', 'Мини-конструктор', 1, 5, 50, 180.00, 'TOY-CONSTRUCT-1-CENTRAL');
    END IF;

    -- Капсулы
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-STD-100-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_capsules, 'capsule', 'Капсулы стандартные', 'Стандартные прозрачные капсулы', 1, 5, 50, 300.00, 'CAP-STD-100-CENTRAL');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-COLOR-100-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_capsules, 'capsule', 'Капсулы цветные', 'Набор цветных капсул', 1, 5, 50, 350.00, 'CAP-COLOR-100-CENTRAL');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-GOLD-50-CENTRAL') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (central_warehouse_id, cat_capsules, 'capsule', 'Капсулы премиум', 'Премиум капсулы золотого цвета', 1, 5, 50, 500.00, 'CAP-GOLD-50-CENTRAL');
    END IF;

    -- Склад "Западный": по 1 единице каждого товара
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM3000-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 3000', 'Вендинговый автомат премиум-класса', 1, 1, 5, 50000.00, 'VM-TM3000-WEST');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM2000-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 2000', 'Вендинговый автомат стандарт-класса', 1, 1, 5, 35000.00, 'VM-TM2000-WEST');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM1000-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 1000', 'Компактный вендинговый автомат', 1, 1, 5, 25000.00, 'VM-TM1000-WEST');
    END IF;

    -- Игрушки
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-SOFT-10-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_toys, 'toy', 'Мягкие игрушки (набор)', 'Набор из 10 мягких игрушек', 1, 5, 50, 150.00, 'TOY-SOFT-10-WEST');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-HEROES-1-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_toys, 'toy', 'Фигурки супергероев', 'Коллекционные фигурки', 1, 5, 50, 200.00, 'TOY-HEROES-1-WEST');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-CARS-5-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_toys, 'toy', 'Машинки миниатюрные', 'Набор миниатюрных машинок', 1, 5, 50, 120.00, 'TOY-CARS-5-WEST');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-CONSTRUCT-1-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_toys, 'toy', 'Конструктор мини', 'Мини-конструктор', 1, 5, 50, 180.00, 'TOY-CONSTRUCT-1-WEST');
    END IF;

    -- Капсулы
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-STD-100-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_capsules, 'capsule', 'Капсулы стандартные', 'Стандартные прозрачные капсулы', 1, 5, 50, 300.00, 'CAP-STD-100-WEST');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-COLOR-100-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_capsules, 'capsule', 'Капсулы цветные', 'Набор цветных капсул', 1, 5, 50, 350.00, 'CAP-COLOR-100-WEST');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-GOLD-50-WEST') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (west_warehouse_id, cat_capsules, 'capsule', 'Капсулы премиум', 'Премиум капсулы золотого цвета', 1, 5, 50, 500.00, 'CAP-GOLD-50-WEST');
    END IF;

    -- Склад "Северный": по 1 единице каждого товара
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM3000-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 3000', 'Вендинговый автомат премиум-класса', 1, 1, 5, 50000.00, 'VM-TM3000-NORTH');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM2000-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 2000', 'Вендинговый автомат стандарт-класса', 1, 1, 5, 35000.00, 'VM-TM2000-NORTH');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'VM-TM1000-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 1000', 'Компактный вендинговый автомат', 1, 1, 5, 25000.00, 'VM-TM1000-NORTH');
    END IF;

    -- Игрушки
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-SOFT-10-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_toys, 'toy', 'Мягкие игрушки (набор)', 'Набор из 10 мягких игрушек', 1, 5, 50, 150.00, 'TOY-SOFT-10-NORTH');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-HEROES-1-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_toys, 'toy', 'Фигурки супергероев', 'Коллекционные фигурки', 1, 5, 50, 200.00, 'TOY-HEROES-1-NORTH');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-CARS-5-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_toys, 'toy', 'Машинки миниатюрные', 'Набор миниатюрных машинок', 1, 5, 50, 120.00, 'TOY-CARS-5-NORTH');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'TOY-CONSTRUCT-1-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_toys, 'toy', 'Конструктор мини', 'Мини-конструктор', 1, 5, 50, 180.00, 'TOY-CONSTRUCT-1-NORTH');
    END IF;

    -- Капсулы
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-STD-100-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_capsules, 'capsule', 'Капсулы стандартные', 'Стандартные прозрачные капсулы', 1, 5, 50, 300.00, 'CAP-STD-100-NORTH');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-COLOR-100-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_capsules, 'capsule', 'Капсулы цветные', 'Набор цветных капсул', 1, 5, 50, 350.00, 'CAP-COLOR-100-NORTH');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM warehouse_inventory WHERE sku = 'CAP-GOLD-50-NORTH') THEN
        INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
        (north_warehouse_id, cat_capsules, 'capsule', 'Капсулы премиум', 'Премиум капсулы золотого цвета', 1, 5, 50, 500.00, 'CAP-GOLD-50-NORTH');
    END IF;
END $$;

-- 5. Убираем секцию с денежными средствами, так как тип 'cash' недопустим
-- Вместо этого добавим дополнительную категорию для аксессуаров

-- 6. Обновляем текущее использование всех складов
UPDATE warehouse 
SET current_usage = (
    SELECT COALESCE(SUM(quantity), 0) 
    FROM warehouse_inventory 
    WHERE warehouse_id = warehouse.id
);

-- 7. Добавляем несколько тестовых отгрузок с курьерами
DO $$
DECLARE
    main_warehouse_id BIGINT;
    location_ids BIGINT[];
    courier_ids BIGINT[];
    toy_item_id BIGINT;
    capsule_item_id BIGINT;
    i INT;
    shipment_id BIGINT;
BEGIN
    SELECT id INTO main_warehouse_id FROM warehouse WHERE name = 'Основной склад' LIMIT 1;
    
    -- Получаем ID локаций
    SELECT array_agg(id) INTO location_ids FROM locations LIMIT 10;
    
    -- Получаем ID курьеров
    SELECT array_agg(id) INTO courier_ids FROM users WHERE userrole = 'courier';
    
    -- Получаем ID товаров
    SELECT id INTO toy_item_id FROM warehouse_inventory WHERE sku = 'TOY-SOFT-10-MAIN' LIMIT 1;
    SELECT id INTO capsule_item_id FROM warehouse_inventory WHERE sku = 'CAP-STD-100-MAIN' LIMIT 1;

    -- Создаем 5 тестовых отгрузок с разными курьерами
    FOR i IN 1..5 LOOP
        IF i <= array_length(courier_ids, 1) THEN
            INSERT INTO warehouse_shipments (warehouse_id, shipment_type, courier_info, shipment_date, status, notes) 
            VALUES (
                main_warehouse_id, 
                'to_courier', 
                'Курьер ID: ' || courier_ids[i] || ' - Регулярная доставка', 
                CURRENT_DATE - (i * 2), 
                CASE 
                    WHEN i = 1 THEN 'delivered'
                    WHEN i = 2 THEN 'shipped' 
                    ELSE 'preparing' 
                END,
                'Отгрузка №' || i || ' для пополнения автоматов'
            )
            RETURNING id INTO shipment_id;

            -- Добавляем товары в отгрузку
            INSERT INTO shipment_items (shipment_id, inventory_item_id, quantity) VALUES
            (shipment_id, toy_item_id, 1),
            (shipment_id, capsule_item_id, 1);
        END IF;
    END LOOP;
END $$;