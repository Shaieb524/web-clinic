# Web Clinic API

Web Clinic API is an authenticated backend application that is used to book doctor-patient appointments with proper administration tools for admin users.
The api is built over golang [Fiber](https://github.com/gofiber/fiber) framework.
You can check [diagrams](./diagrams) for better understanding of the application logic.

## Installation

You need to have [GO](https://go.dev/) set up and working on your device and clone the source code:

```bash
git clone https://github.com/Shaieb524/web-clinic.git
```

## Run

Prepare your environment variables (please refer to [.env-example](./.env-example) file)

to run: 
```bash
go run ./main.go
```

also, you can run it in development auto-restart mode using [nodemon](https://www.npmjs.com/package/nodemon)

```bash
nodemon --exec "go run" ./main.go --signal SIGTERM
```
## Usage
### Postmant collection and local environment:
[clinic-web-api.postman_collection.json](./postman/clinic-web-api.postman_collection.json)

[local.postman_environment.json](./postman/local.postman_environment.json)

## Test
Run scripts inside [tests directory](./tests) with :

```bash
go test -v ./tests/test_script.go
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Deployment
Master branch is connected with [Heroku](https://dashboard.heroku.com/) cloud hosting.

```bash
git remote add heroku https://git.heroku.com/web-clinic-api.git
```

```bash
git push heroku master
```

