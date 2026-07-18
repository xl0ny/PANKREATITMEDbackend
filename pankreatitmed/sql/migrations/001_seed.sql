-- ==== USERS ====
INSERT INTO medusers (login, password, is_moderator) VALUES
                                                         ('demo','demo',true),
                                                         ('moder','moder',true),
                                                         ('alice','alice',false),
                                                         ('bob','bob',false)
    ON CONFLICT (login) DO NOTHING;

-- ==== CRITERIA ====
INSERT INTO criteria
(code, name, description, duration, home_visit, image_url, status, unit, ref_low, ref_high)
VALUES
    ('№1','Оценка возраста пациента',
     'Возраст > 55 лет — критерий при поступлении.',
     '1 календарный день', true,
     'http://127.0.0.1:9000/services-images/n1_age.png',
     'active','лет',NULL,55),

    ('№2','Анализ лейкоцитов крови',
     'Повышение лейкоцитов может указывать на выраженный воспалительный процесс (критерий: > 16 000/мм³).',
     '1 календарный день', true,
     'http://127.0.0.1:9000/services-images/n2_wbc.png',
     'active','10^9/л',4.0,11.0),

    ('№3','Измерение уровня глюкозы',
     'Гипергликемия — один из ранних критериев (критерий: > 200 мг/дл ≈ 11,1 ммоль/л).',
     '1 календарный день', true,
     'http://127.0.0.1:9000/services-images/n3_glucose.png',
     'active','мг/дл',70,99),

    ('№4','Определение уровня ЛДГ',
     'ЛДГ > 350 МЕ/л — критерий тяжести.',
     '1 календарный день', true,
     'http://127.0.0.1:9000/services-images/n4_ldh.png',
     'active','Ед/л',135,225),

    ('№5','Анализ активности АСТ',
     'АСТ > 250 МЕ/л — критерий тяжести.',
     '1 календарный день', true,
     'http://127.0.0.1:9000/services-images/n5_ast.png',
     'active','Ед/л',0,40),

    ('№6','Контроль изменения гематокрита',
     'Снижение гематокрита > 10% за 48 часов — неблагоприятный признак.',
     '1 календарный день через 48 часов', true,
     'http://127.0.0.1:9000/services-images/n6_hct.png',
     'active','%',NULL,NULL),

    ('№7','Измерение уровня мочевины (BUN)',
     'Рост мочевины > 5 мг/дл за 48 часов — критерий ухудшения.',
     '1 календарный день через 48 часов', true,
     'http://127.0.0.1:9000/services-images/n7_bun.png',
     'active','мг/дл',7,20),

    ('№8','Измерение уровня кальция сыворотки',
     'Гипокальциемия (< 8,0 мг/дл ≈ 2,0 ммоль/л) — прогностический критерий.',
     '1 календарный день', true,
     'http://127.0.0.1:9000/services-images/n8_ca.png',
     'active','мг/дл',8.6,10.2),

    ('№9','Измерение PaO₂',
     'PaO₂ < 60 мм рт.ст. — критерий в шкале Рэнсона.',
     '1 календарный день', true,
     'http://127.0.0.1:9000/services-images/n9_pao2.png',
     'active','мм рт.ст.',80,100),

    ('№10','Оценка кислотно-щелочного состояния',
     'Дефицит оснований > 4 мЭкв/л — неблагоприятен.',
     '1 календарный день через 48 часов', true,
     'http://127.0.0.1:9000/services-images/n10_acidbase.png',
     'active','мЭкв/л',-2,2),

    ('№11','Оценка объёма секвестрированной жидкости',
     'Секвестрированная жидкость > 6 л за 48 часов — высокий риск осложнений.',
     '1 календарный день через 48 часов', true,
     'http://127.0.0.1:9000/services-images/n11_sequestration.png',
     'active','л',NULL,NULL)
    ON CONFLICT (code) DO NOTHING;

-- ==== ORDER #1: draft (alice) + items ====
WITH o AS (
INSERT INTO medorders
(status, created_at, creator_id, formed_at, finished_at, moderator_id, ranson_score, mortality_risk)
VALUES
    ('draft', NOW(), (SELECT id FROM medusers WHERE login='alice'), NULL, NULL, NULL, NULL, NULL)
    RETURNING id
    )
INSERT INTO medorderitems (med_order_id, criterion_id, position, value_num, value_indicator)
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№1'), 1, NULL, TRUE
UNION ALL
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№3'), 2, 212.0, TRUE
    ON CONFLICT (med_order_id, criterion_id) DO NOTHING;

-- ==== ORDER #2: formed (alice, moder) + items ====
WITH o AS (
INSERT INTO medorders
(status, created_at, creator_id, formed_at, finished_at, moderator_id, ranson_score, mortality_risk)
VALUES
    ('formed', NOW(), (SELECT id FROM medusers WHERE login='alice'), NOW(), NULL,
    (SELECT id FROM medusers WHERE login='moder'), 3, 'около 16%')
    RETURNING id
    )
INSERT INTO medorderitems (med_order_id, criterion_id, position, value_num, value_indicator)
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№2'), 1, 17.5, TRUE
UNION ALL
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№4'), 2, 420.0, TRUE
UNION ALL
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№8'), 3, 1.95, TRUE
    ON CONFLICT (med_order_id, criterion_id) DO NOTHING;

-- ==== ORDER #3: completed (bob, moder) + items ====
WITH o AS (
INSERT INTO medorders
(status, created_at, creator_id, formed_at, finished_at, moderator_id, ranson_score, mortality_risk)
VALUES
    ('completed', NOW(), (SELECT id FROM medusers WHERE login='bob'), NOW(), NOW(),
    (SELECT id FROM medusers WHERE login='moder'), 4, 'около 40%')
    RETURNING id
    )
INSERT INTO medorderitems (med_order_id, criterion_id, position, value_num, value_indicator)
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№1'), 1, NULL, TRUE
UNION ALL
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№5'), 2, 310.0, TRUE
UNION ALL
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№9'), 3, 58.0, TRUE
UNION ALL
SELECT (SELECT id FROM o), (SELECT id FROM criteria WHERE code='№10'), 4, 5.1, TRUE
    ON CONFLICT (med_order_id, criterion_id) DO NOTHING;