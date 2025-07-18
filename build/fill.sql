
-- Заполняем таблицу Categories
-- INSERT INTO categories (ID, Name, Description) VALUES
-- (1, 'Еда', 'Питание, кафе и рестораны'),
-- (2, 'Транспорт', 'Расходы на транспорт и бензин'),
-- (3, 'Развлечения', 'Кино, театр, концерты'),
-- (4, 'Одежда', 'Покупка одежды и аксессуаров'),
-- (5, 'Коммунальные услуги', 'Оплата счетов за жилье и ЖКУ');

INSERT INTO recommendations (id, title, description) VALUES
(1, 'Повысились расходы в категории «%s»',
 'Я вижу, что в этом месяце твои траты по категории «%s» выросли аж на %s по сравнению с прошлым месяцем. Давай попробуем разобраться: это было осознанное решение или импульсивные траты? Может, тебе пришлось взять что-то срочное или просто расслабился и позволил себе больше обычного? Если это что-то важное — отлично! Но если эти расходы могли подождать или были вызваны эмоциями, возможно, стоит пересмотреть подход. В долгосрочной перспективе контроль над повторяющимися расходами поможет сохранить баланс и достичь твоих целей.'),

(2, 'Баланс снизился по сравнению с прошлым месяцем',
'Обрати внимание, прошлый месяц закрылся с большим остатком на счету аж на %s, чем этот. Это знак того, что либо доходы уменьшились, либо расходы увеличились. Если это временные трудности — ничего страшного, главное, чтобы это не стало системной проблемой. Возможно, стоит на какое-то время сфокусироваться на накоплении? Например, можно пересмотреть лишние покупки или попробовать откладывать хотя бы 10%% от каждого поступления себе «на подушку безопасности». Даже небольшие регулярные накопления дадут тебе финансовую уверенность!'),

(3, 'Доходы не изменились, но расходы увеличились',
 'Ты зарабатываешь столько же, но почему-то тратить начал больше аж на %s. Это может быть случайность, но если это повторяется — стоит понять, куда утекают деньги. Иногда мы даже не замечаем, как мелкие ежедневные покупки складываются в ощутимую сумму. Попробуй хотя бы неделю фиксировать абсолютно все расходы — возможно, откроешь для себя неожиданное!'),

(4, 'Постоянное превышение расходов над доходами',
 'Кажется, что твои расходы стабильно превышают доходы. В этом месяце ты потратил %s, а накопил %s. Это тревожный звоночек! Если так продолжать, можно быстро оказаться в ситуации, где придется брать в долг. Давай вместе посмотрим, как можно сбалансировать бюджет. Может, есть ненужные подписки, которые можно отменить? Или привычки, которые забирают больше денег, чем ты осознавал? Контролировать траты — это не про ограничения, а про осознанный выбор!'),

(5, 'Основная часть бюджета уходит на категорию «%s»',
 'Я заметил, что в этом месяце большая часть твоих расходов сосредоточена в категории «%s». Это осознанный приоритет или просто так сложилось? Если ты действительно получаешь от этого удовольствие и пользу — окей. Но если это нечто, что можно оптимизировать, почему бы не попробовать? Например, если это еда вне дома, можно часть недели готовить самому, а если это развлечения — искать более бюджетные альтернативы.'),

(6, 'Нет накоплений в этом месяце',
 'Дружище, кажется, в этом месяце у тебя не было отложенных денег на будущее. Конечно, жизнь должна приносить радость, но даже небольшие накопления помогут тебе избежать неожиданных финансовых стрессов. Представь, если завтра сломается телефон или придется срочно потратиться на здоровье — будет неприятно, если свободных денег не окажется. Почему бы не ввести правило: откладывать хотя бы 5-10% от поступлений?'),

(7, 'Частые импульсивные покупки',
 'Похоже, что в этом месяце было много несистемных трат — то есть покупок, которые не были запланированы заранее. Это окей, если они приносят тебе пользу или радость. Но импульсивные покупки могут незаметно съедать гораздо больше бюджета, чем кажется. Лайфхак: перед покупкой чего-то незапланированного попробуй подождать 24 часа. Если на следующий день все еще хочется это купить — значит, действительно нужно!')
;

