docker run -p 5432:5432 -v ./postgres-data:/var/lib/postgresql/data -it --name postgres --rm -d -e POSTGRES_PASSWORD=postgres postgres
# psql "sslmode=disable host=localhost user=postgres password=postgres dbname=postgres"