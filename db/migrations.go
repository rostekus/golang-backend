package db

func (d *Database) MigrateDB() error {

	createTableSQL := `CREATE TABLE IF NOT EXISTS users(
    "id" varchar(60) PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar NOT NULL
);`

	_, err := d.db.Exec(createTableSQL)
	return err
}
