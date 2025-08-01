-- Добавить автоинкремент в таблицу
CREATE SEQUENCE IF NOT EXISTS notes_id_seq OWNED BY notes.id;
ALTER TABLE notes ALTER COLUMN id SET DEFAULT nextval('notes_id_seq');