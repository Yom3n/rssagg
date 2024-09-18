To run migrations cd to schema folder and run

goose postgres postgres://postgres:admin@localhost:5432/rssagg up

To generate sqlc generation in root run:
sqlc generate