# Greenlight

Practice Go course of Alex Edwards - Let's Go Further. Consists in a JSON API managing information about movies.

## Project stack
- Go 1.24
- Postgres 

## Requirements

Before you run this project you have to have the Go installed on your machine.

### Create database in Postgres

This project was developed using the Postgres database. You have to provide it to your deployment environment.

After that, make sure to create the database called `greenlight`.

```
$ postgres# CREATE DATABASE greenlight;
```

### Environment Variable
Please make sure to add an environment variable called `GREENLIGHT_DB_DSN`. It should contain a connection string to your database just like this:

```
postgres://<user>:<pass>@<server>/greenlight?sslmode=disable
``` 

## Endpoints

## Authentication

## Contribution