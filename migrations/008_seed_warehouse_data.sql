-- 1. Создаем основной склад
INSERT INTO warehouse (name, address, contact_person, contact_phone, total_capacity, current_usage) 
SELECT 'Основной склад', 'ул. Складская, 15, Москва', 'Смирнов Александр', '+7-999-111-22-33', 5000, 0
WHERE NOT EXISTS (SELECT 1 FROM warehouse WHERE name = 'Основной склад');

-- 2. Добавляем категории товаров
INSERT INTO warehouse_categories (name, description) VALUES
('Вендинговые автоматы', 'Различные модели вендинговых автоматов для игрушек'),
('Игрушки', 'Игрушки для наполнения автоматов'),
('Капсулы', 'Капсулы для упаковки игрушек')
ON CONFLICT (name) DO NOTHING;

-- 3. Добавляем инвентарь на склад
DO $$
DECLARE
    warehouse_id BIGINT;
    cat_machines BIGINT;
    cat_toys BIGINT;
    cat_capsules BIGINT;
BEGIN
    -- Получаем ID склада и категорий
    SELECT id INTO warehouse_id FROM warehouse WHERE name = 'Основной склад' LIMIT 1;
    SELECT id INTO cat_machines FROM warehouse_categories WHERE name = 'Вендинговые автоматы' LIMIT 1;
    SELECT id INTO cat_toys FROM warehouse_categories WHERE name = 'Игрушки' LIMIT 1;
    SELECT id INTO cat_capsules FROM warehouse_categories WHERE name = 'Капсулы' LIMIT 1;

    -- Вендинговые автоматы
    INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
    (warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 3000', 'Вендинговый автомат премиум-класса, вместимость 100 игрушек', 5, 2, 10, 50000.00, 'VM-TM3000'),
    (warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 2000', 'Вендинговый автомат стандарт-класса, вместимость 80 игрушек', 3, 1, 5, 35000.00, 'VM-TM2000'),
    (warehouse_id, cat_machines, 'vending_machine', 'ToyMaster 1000', 'Компактный вендинговый автомат, вместимость 50 игрушек', 2, 1, 3, 25000.00, 'VM-TM1000')
    ON CONFLICT (sku) DO NOTHING;

    -- Игрушки
    INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
    (warehouse_id, cat_toys, 'toy', 'Мягкие игрушки (набор)', 'Набор из 10 мягких игрушек разных животных', 150, 50, 500, 150.00, 'TOY-SOFT-10'),
    (warehouse_id, cat_toys, 'toy', 'Фигурки супергероев', 'Коллекционные фигурки популярных супергероев', 200, 100, 1000, 200.00, 'TOY-HEROES-1'),
    (warehouse_id, cat_toys, 'toy', 'Машинки миниатюрные', 'Набор миниатюрных машинок разных моделей', 180, 80, 800, 120.00, 'TOY-CARS-5'),
    (warehouse_id, cat_toys, 'toy', 'Конструктор мини', 'Мини-конструктор для сборки различных моделей', 120, 60, 600, 180.00, 'TOY-CONSTRUCT-1')
    ON CONFLICT (sku) DO NOTHING;

    -- Капсулы
    INSERT INTO warehouse_inventory (warehouse_id, category_id, item_type, item_name, description, quantity, min_stock_level, max_stock_level, unit_price, sku) VALUES
    (warehouse_id, cat_capsules, 'capsule', 'Капсулы стандартные (прозрачные)', 'Стандартные прозрачные капсулы для игрушек, 100 шт.', 80, 20, 200, 300.00, 'CAP-STD-100'),
    (warehouse_id, cat_capsules, 'capsule', 'Капсулы цветные (набор)', 'Набор цветных капсул, 5 цветов по 20 шт.', 60, 15, 150, 350.00, 'CAP-COLOR-100'),
    (warehouse_id, cat_capsules, 'capsule', 'Капсулы премиум (золотые)', 'Премиум капсулы золотого цвета, 50 шт.', 30, 10, 100, 500.00, 'CAP-GOLD-50')
    ON CONFLICT (sku) DO NOTHING;
END $$;

-- 4. Добавляем ожидаемые поставки
DO $$
DECLARE
    warehouse_id BIGINT;
    toy_item_id BIGINT;
    capsule_item_id BIGINT;
    machine_item_id BIGINT;
    supply_id BIGINT;
BEGIN
    SELECT id INTO warehouse_id FROM warehouse WHERE name = 'Основной склад' LIMIT 1;
    
    -- Получаем ID товаров для поставок
    SELECT id INTO toy_item_id FROM warehouse_inventory WHERE sku = 'TOY-SOFT-10' LIMIT 1;
    SELECT id INTO capsule_item_id FROM warehouse_inventory WHERE sku = 'CAP-STD-100' LIMIT 1;
    SELECT id INTO machine_item_id FROM warehouse_inventory WHERE sku = 'VM-TM3000' LIMIT 1;

    -- Поставка игрушек и капсул
    INSERT INTO warehouse_supplies (warehouse_id, supplier_name, supply_date, expected_date, status, total_amount, notes) 
    SELECT warehouse_id, 'ООО "ИгрушкиОпт"', '2024-02-01', '2024-02-10', 'ordered', 75000.00, 'Регулярная поставка игрушек и капсул'
    WHERE NOT EXISTS (SELECT 1 FROM warehouse_supplies WHERE supplier_name = 'ООО "ИгрушкиОпт"' AND supply_date = '2024-02-01')
    RETURNING id INTO supply_id;

    IF FOUND THEN
        INSERT INTO supply_items (supply_id, inventory_item_id, quantity_ordered, unit_price) VALUES
        (supply_id, toy_item_id, 300, 150.00),
        (supply_id, capsule_item_id, 100, 300.00);
    END IF;

    -- Поставка новых автоматов
    INSERT INTO warehouse_supplies (warehouse_id, supplier_name, supply_date, expected_date, status, total_amount, notes) 
    SELECT warehouse_id, 'Завод "ВендингМаш"', '2024-02-05', '2024-02-15', 'ordered', 200000.00, 'Поставка новых автоматов ToyMaster 3000'
    WHERE NOT EXISTS (SELECT 1 FROM warehouse_supplies WHERE supplier_name = 'Завод "ВендингМаш"' AND supply_date = '2024-02-05')
    RETURNING id INTO supply_id;

    IF FOUND THEN
        INSERT INTO supply_items (supply_id, inventory_item_id, quantity_ordered, unit_price) VALUES
        (supply_id, machine_item_id, 4, 50000.00);
    END IF;
END $$;

-- 5. Добавляем отгрузки
DO $$
DECLARE
    warehouse_id BIGINT;
    location1_id BIGINT;
    location2_id BIGINT;
    toy_item_id BIGINT;
    capsule_item_id BIGINT;
    machine_item_id BIGINT;
    machine_item2_id BIGINT;
    vending_machine_id BIGINT;
    shipment_id BIGINT;
BEGIN
    SELECT id INTO warehouse_id FROM warehouse WHERE name = 'Основной склад' LIMIT 1;
    SELECT id INTO location1_id FROM locations WHERE name = 'ТЦ "Москва"' LIMIT 1;
    SELECT id INTO location2_id FROM locations WHERE name = 'ТРК "Европа"' LIMIT 1;
    
    -- Получаем ID товаров
    SELECT id INTO toy_item_id FROM warehouse_inventory WHERE sku = 'TOY-SOFT-10' LIMIT 1;
    SELECT id INTO capsule_item_id FROM warehouse_inventory WHERE sku = 'CAP-STD-100' LIMIT 1;
    SELECT id INTO machine_item_id FROM warehouse_inventory WHERE sku = 'VM-TM3000' LIMIT 1;
    SELECT id INTO machine_item2_id FROM warehouse_inventory WHERE sku = 'VM-TM2000' LIMIT 1;
    SELECT id INTO vending_machine_id FROM vending_machines WHERE serial_number = 'VM004' LIMIT 1;

    -- Отгрузка автомата в локацию
    INSERT INTO warehouse_shipments (warehouse_id, shipment_type, target_location_id, shipment_date, status, notes) 
    SELECT warehouse_id, 'to_location', location1_id, '2024-01-20', 'delivered', 'Установка нового автомата в ТЦ "Москва"'
    WHERE NOT EXISTS (SELECT 1 FROM warehouse_shipments WHERE target_location_id = location1_id AND shipment_date = '2024-01-20')
    RETURNING id INTO shipment_id;

    IF FOUND THEN
        INSERT INTO shipment_items (shipment_id, inventory_item_id, vending_machine_id, quantity) VALUES
        (shipment_id, machine_item_id, vending_machine_id, 1);
    END IF;

    -- Отгрузка игрушек и капсул курьеру
    INSERT INTO warehouse_shipments (warehouse_id, shipment_type, courier_info, shipment_date, status, notes) 
    SELECT warehouse_id, 'to_courier', 'Курьер: Иванов П.С., тел. +7-999-444-55-66', '2024-01-25', 'shipped', 'Отгрузка для пополнения автоматов в ТРК "Европа"'
    WHERE NOT EXISTS (SELECT 1 FROM warehouse_shipments WHERE courier_info LIKE '%Иванов П.С.%' AND shipment_date = '2024-01-25')
    RETURNING id INTO shipment_id;

    IF FOUND THEN
        INSERT INTO shipment_items (shipment_id, inventory_item_id, quantity) VALUES
        (shipment_id, toy_item_id, 50),
        (shipment_id, capsule_item_id, 20);
    END IF;

    -- Отгрузка для технического обслуживания
    INSERT INTO warehouse_shipments (warehouse_id, shipment_type, target_location_id, shipment_date, status, notes) 
    SELECT warehouse_id, 'to_location', location2_id, '2024-01-28', 'preparing', 'Замена неисправного автомата'
    WHERE NOT EXISTS (SELECT 1 FROM warehouse_shipments WHERE target_location_id = location2_id AND shipment_date = '2024-01-28')
    RETURNING id INTO shipment_id;

    IF FOUND THEN
        INSERT INTO shipment_items (shipment_id, inventory_item_id, quantity) VALUES
        (shipment_id, machine_item2_id, 1);
    END IF;
END $$;

-- Обновляем текущее использование склада
UPDATE warehouse 
SET current_usage = (
    SELECT COALESCE(SUM(quantity), 0) 
    FROM warehouse_inventory 
    WHERE warehouse_id = warehouse.id
)
WHERE name = 'Основной склад';
