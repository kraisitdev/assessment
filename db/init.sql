DROP TABLE IF EXISTS expenses;

CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);

INSERT INTO expenses (id, title, amount, note, tags) values ('100', 'strawberry smoothie', 79, 'night market promotion discount 10 bath', ARRAY ['food','beverage']);

INSERT INTO expenses (id, title, amount, note, tags) values ('200', 'strawberry smoothie', 79, 'night market promotion discount 10 bath', ARRAY ['food','beverage']);

