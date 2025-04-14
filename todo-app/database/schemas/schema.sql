CREATE TABLE IF NOT EXISTS todo_items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	item_number INTEGER NOT NULL,
	content TEXT NOT NULL,
	date_created timestamp NOT NULL,
	is_done BOOL DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_info (
	password_hash TEXT NOT NULL,
	date_created timestamp NOT NULL
);
