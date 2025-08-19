CREATE TABLE IF NOT EXISTS todo_items (
	item_id INTEGER PRIMARY KEY AUTOINCREMENT,
	list_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	date_created timestamp NOT NULL,
	is_done BOOL DEFAULT FALSE NOT NULL,
	foreign key(list_id) references lists(list_id)
);

CREATE TABLE IF NOT EXISTS user_info (
	password_hash TEXT NOT NULL,
	date_created timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS lists (
	list_id INTEGER PRIMARY KEY AUTOINCREMENT,
	list_name TEXT NOT NULL,
	icon_color TEXT NOT NULL,
	date_created timestamp NOT NULL
);

PRAGMA journal_mode = WAL;
