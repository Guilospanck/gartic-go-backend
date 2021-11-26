# Gartic Go Backend
This is the backend application of the Gartic-like application. The frontend can be found at [Gartic-Like-App](https://github.com/Guilospanck/gartic-like-app/)

## Installation
### Install Go
You can find a tutorial at [Digital Ocean - Go Ubuntu 20.04](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-ubuntu-20-04)

### Database
We're using Postgres as database for this project. Therefore, in order to run the application, you must have it up and running before.
The quicker and most practical way is by using Docker.

After having Docker installed, run
```bash
sudo docker run --rm --name pg-docker -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=default -d -p 5432:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data postgres
```

### Air
Air is a live reload for Go. Install it by following their [GitHub page](https://github.com/cosmtrek/air)

### Repository
Git clone this repository
```bash
git clone https://github.com/Guilospanck/gartic-go-backend.git
```
Change directory into it
```bash
cd gartic-go-backend
```
Run the following command to install dependencies
```bash
go mod tidy
```
And finally
```bash
air -d
```

## Tip
Checkout `template/clean_arch` for the Clean Architecture Template for Golang.
