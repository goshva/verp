-- Migration: 011_add_vending_machines_and_operations.sql

-- 1. Добавляем 15 вендинговых автоматов на различные локации
DO $$
DECLARE
    location_ids BIGINT[];
    i INT;
    location_index INT;
    serial_num TEXT;
    machine_model TEXT;
    machine_status TEXT;
    toys_count INT;
    cash_amount DECIMAL(10,2);
BEGIN
    -- Получаем ID всех активных локаций
    SELECT array_agg(id) INTO location_ids FROM locations WHERE is_active = true;
    
    -- Добавляем 15 автоматов
    FOR i IN 1..15 LOOP
        -- Выбираем локацию по кругу
        location_index := ((i - 1) % array_length(location_ids, 1)) + 1;
        
        -- Генерируем уникальный серийный номер
        serial_num := 'VM-' || LPAD((100 + i)::text, 3, '0');
        
        -- Выбираем модель в зависимости от номера
        IF i % 3 = 0 THEN
            machine_model := 'ToyMaster 1000';
            toys_count := 30 + (i * 3) % 40;
            cash_amount := 800.00 + (i * 50)::decimal;
        ELSIF i % 3 = 1 THEN
            machine_model := 'ToyMaster 2000';
            toys_count := 45 + (i * 2) % 35;
            cash_amount := 1200.00 + (i * 75)::decimal;
        ELSE
            machine_model := 'ToyMaster 3000';
            toys_count := 60 + i % 30;
            cash_amount := 1500.00 + (i * 100)::decimal;
        END IF;
        
        -- Выбираем статус (большинство активны, некоторые на обслуживании)
        IF i % 7 = 0 THEN
            machine_status := 'maintenance';
            toys_count := 0;
            cash_amount := 0;
        ELSE
            machine_status := 'active';
        END IF;
        
        -- Вставляем автомат, если его еще нет
        IF NOT EXISTS (SELECT 1 FROM vending_machines WHERE serial_number = serial_num) THEN
            INSERT INTO vending_machines (
                serial_number, model, status, location_id, 
                current_toys_count, capacity_toys, cash_amount,
                installation_date, last_maintenance_date, next_maintenance_date
            ) VALUES (
                serial_num,
                machine_model,
                machine_status,
                location_ids[location_index],
                toys_count,
                CASE 
                    WHEN machine_model = 'ToyMaster 1000' THEN 50
                    WHEN machine_model = 'ToyMaster 2000' THEN 80
                    ELSE 100
                END,
                cash_amount,
                CURRENT_DATE - (i * 10 || ' days')::interval,
                CASE WHEN machine_status = 'maintenance' THEN NULL ELSE CURRENT_DATE - (i * 5 || ' days')::interval END,
                CASE WHEN machine_status = 'maintenance' THEN NULL ELSE CURRENT_DATE + ((30 + (i * 2)) || ' days')::interval END
            );
        END IF;
    END LOOP;
END $$;

-- 2. Добавляем операции для автоматов
DO $$
DECLARE
    machine_record RECORD;
    user_ids BIGINT[];
    i INT;
    user_id BIGINT;
    op_date TIMESTAMP;
    operation_type TEXT;
    toys_before INT;
    toys_after INT;
    toys_added INT;
    cash_before DECIMAL(10,2);
    cash_after DECIMAL(10,2);
    cash_collected DECIMAL(10,2);
BEGIN
    -- Получаем пользователей
    SELECT array_agg(id) INTO user_ids FROM users WHERE userrole IN ('admin', 'operator1', 'tech1', 'agent');
    
    -- Для каждого активного автомата добавляем операции
    FOR machine_record IN (
        SELECT id, current_toys_count, cash_amount 
        FROM vending_machines 
        WHERE status = 'active'
    ) LOOP
        -- Начальные значения
        toys_before := machine_record.current_toys_count;
        cash_before := machine_record.cash_amount;
        
        -- Добавляем 2-3 операции для этого автомата
        FOR i IN 1..(2 + (machine_record.id % 2)) LOOP
            -- Выбираем случайного пользователя
            user_id := user_ids[1 + ((machine_record.id + i - 1) % array_length(user_ids, 1))];
            
            -- Определяем дату операции (разные дни) - исправлено приведение типов
            op_date := CURRENT_TIMESTAMP - ((i * 3 + machine_record.id * 2) || ' days')::interval;
            
            -- Определяем тип операции
            IF i % 3 = 0 THEN
                operation_type := 'maintenance';
                toys_after := toys_before;
                toys_added := 0;
                cash_after := cash_before;
                cash_collected := 0;
            ELSIF i % 3 = 1 THEN
                operation_type := 'restock';
                toys_added := 20 + (machine_record.id * i) % 30;
                toys_after := LEAST(toys_before + toys_added, 100);
                cash_after := cash_before;
                cash_collected := 0;
            ELSE
                operation_type := 'collection';
                toys_after := toys_before;
                toys_added := 0;
                cash_collected := cash_before * 0.8; -- собираем 80% денег
                cash_after := cash_before - cash_collected;
            END IF;
            
            -- Вставляем операцию
            INSERT INTO vending_operations (
                vending_machine_id, operation_type, performed_by, operation_date,
                toys_before, toys_after, toys_added,
                cash_before, cash_after, cash_collected,
                notes
            ) VALUES (
                machine_record.id,
                operation_type,
                user_id,
                op_date,
                toys_before,
                toys_after,
                toys_added,
                cash_before,
                cash_after,
                cash_collected,
                CASE 
                    WHEN operation_type = 'restock' THEN 'Пополнение игрушек. ' || toys_added || ' шт. добавлено'
                    WHEN operation_type = 'collection' THEN 'Инкассация денежных средств. Собрано: ' || cash_collected || ' руб.'
                    WHEN operation_type = 'maintenance' THEN 'Плановое техническое обслуживание'
                END
            );
            
            -- Обновляем значения для следующей операции
            toys_before := toys_after;
            cash_before := cash_after;
        END LOOP;
        
        -- Обновляем данные автомата после всех операций
        UPDATE vending_machines 
        SET 
            current_toys_count = toys_after,
            cash_amount = cash_after,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = machine_record.id;
    END LOOP;
END $$;

-- 3. Добавляем дополнительные операции для некоторых автоматов (исторические данные)
DO $$
DECLARE
    machine_record RECORD;
    user_ids BIGINT[];
    user_id BIGINT;
    op_date TIMESTAMP;
BEGIN
    -- Получаем пользователей
    SELECT array_agg(id) INTO user_ids FROM users WHERE userrole IN ('admin', 'operator1');
    
    -- Для первых 5 автоматов добавляем исторические операции
    FOR machine_record IN (
        SELECT id 
        FROM vending_machines 
        ORDER BY id 
        LIMIT 5
    ) LOOP
        user_id := user_ids[1 + (machine_record.id % array_length(user_ids, 1))];
        
        -- Операция 2 месяца назад
        op_date := CURRENT_TIMESTAMP - '60 days'::interval;
        IF NOT EXISTS (
            SELECT 1 FROM vending_operations 
            WHERE vending_machine_id = machine_record.id 
            AND operation_date = op_date
        ) THEN
            INSERT INTO vending_operations (
                vending_machine_id, operation_type, performed_by, operation_date,
                toys_before, toys_after, toys_added,
                cash_before, cash_after, cash_collected,
                notes
            ) VALUES (
                machine_record.id,
                'restock',
                user_id,
                op_date,
                15,
                65,
                50,
                500.00,
                500.00,
                0,
                'Регулярное пополнение игрушек'
            );
        END IF;
        
        -- Операция 1 месяц назад
        op_date := CURRENT_TIMESTAMP - '30 days'::interval;
        IF NOT EXISTS (
            SELECT 1 FROM vending_operations 
            WHERE vending_machine_id = machine_record.id 
            AND operation_date = op_date
        ) THEN
            INSERT INTO vending_operations (
                vending_machine_id, operation_type, performed_by, operation_date,
                toys_before, toys_after, toys_added,
                cash_before, cash_after, cash_collected,
                notes
            ) VALUES (
                machine_record.id,
                'collection',
                user_id,
                op_date,
                40,
                40,
                0,
                1800.00,
                200.00,
                1600.00,
                'Ежемесячная инкассация'
            );
        END IF;
        
        -- Операция 2 недели назад
        op_date := CURRENT_TIMESTAMP - '14 days'::interval;
        IF NOT EXISTS (
            SELECT 1 FROM vending_operations 
            WHERE vending_machine_id = machine_record.id 
            AND operation_date = op_date
        ) THEN
            INSERT INTO vending_operations (
                vending_machine_id, operation_type, performed_by, operation_date,
                toys_before, toys_after, toys_added,
                cash_before, cash_after, cash_collected,
                notes
            ) VALUES (
                machine_record.id,
                'restock',
                user_id,
                op_date,
                20,
                70,
                50,
                600.00,
                600.00,
                0,
                'Пополнение игрушек после выходных'
            );
        END IF;
    END LOOP;
END $$;

-- 4. Обновляем статистику по самым успешным автоматам
DO $$
DECLARE
    machine_record RECORD;
BEGIN
    -- Увеличиваем количество игрушек и денег у самых популярных автоматов
    FOR machine_record IN (
        SELECT id, current_toys_count, cash_amount, capacity_toys
        FROM vending_machines 
        WHERE status = 'active'
        ORDER BY cash_amount DESC 
        LIMIT 3
    ) LOOP
        UPDATE vending_machines 
        SET 
            current_toys_count = LEAST(machine_record.current_toys_count + 15, machine_record.capacity_toys),
            cash_amount = machine_record.cash_amount + 500.00,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = machine_record.id;
    END LOOP;
END $$;

-- 5. Добавляем операции обслуживания для автоматов на техобслуживании
DO $$
DECLARE
    machine_record RECORD;
    user_ids BIGINT[];
    user_id BIGINT;
    op_date TIMESTAMP;
BEGIN
    -- Получаем технических пользователей
    SELECT array_agg(id) INTO user_ids FROM users WHERE userrole = 'technician';
    
    IF array_length(user_ids, 1) > 0 THEN
        -- Для каждого автомата на обслуживании добавляем операцию
        FOR machine_record IN (
            SELECT id 
            FROM vending_machines 
            WHERE status = 'maintenance'
        ) LOOP
            user_id := user_ids[1 + ((machine_record.id - 1) % array_length(user_ids, 1))];
            
            -- Добавляем операцию обслуживания
            op_date := CURRENT_TIMESTAMP - '2 days'::interval;
            IF NOT EXISTS (
                SELECT 1 FROM vending_operations 
                WHERE vending_machine_id = machine_record.id 
                AND operation_type = 'maintenance'
                AND operation_date >= CURRENT_TIMESTAMP - '7 days'::interval
            ) THEN
                INSERT INTO vending_operations (
                    vending_machine_id, operation_type, performed_by, operation_date,
                    toys_before, toys_after, toys_added,
                    cash_before, cash_after, cash_collected,
                    notes
                ) VALUES (
                    machine_record.id,
                    'maintenance',
                    user_id,
                    op_date,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,
                    'Диагностика и ремонт оборудования. Замена механических компонентов.'
                );
            END IF;
        END LOOP;
    END IF;
END $$;

-- 6. Создаем индекс для оптимизации запросов по операциям
CREATE INDEX IF NOT EXISTS idx_vending_operations_machine_date 
ON vending_operations(vending_machine_id, operation_date DESC);

CREATE INDEX IF NOT EXISTS idx_vending_machines_location_status 
ON vending_machines(location_id, status);

-- 7. Обновляем временные метки
UPDATE vending_machines SET updated_at = CURRENT_TIMESTAMP WHERE updated_at < created_at;
UPDATE vending_operations SET updated_at = CURRENT_TIMESTAMP WHERE updated_at < created_at;