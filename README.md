# User Management System Blueprint Server

![Blueprint](/static/logos/logo_banner.png)

## Description

This project provides a user management system for the Blueprint Server. Stevens Blueprint uses Authelia as
SSO application. Authelia provides a YAML file called users.yaml where all the authorized users are stored.
This API intends to provide a way to manage the users.yaml. The API should provide endpoints to create, add,
delete, update and disable users. The service will follow a RESTful API architecture. 

## Installation

1. Clone the repository:
```
git clone https://github.com/your-username/user-management-system.git
```

2. Build the project:

```
go build
```

3. To run the server. Use the -dev flag to run the service in dev mode.
```
go run main.go
```

## Configuration

Make sure go is installed in your local. You can verify this by running
```
go version
```
If go is not installed, you can refer to this documentation to install it. https://go.dev/doc/install

Install dependencies

```
go mod download
```

Add .env file

Get the absolute PATH of your working directory. If you are using VSCode in the terminal run
```
pwd
```
Paste the output of the command in the .env file with the following format
```
PATH=(Output of pwd)
```

## Running the service
Run
```
go run main.go
```

## Docs
Available Endpoints
``
GET /api/v1/all
``

## Contributing

Contributions are welcome! If you would like to contribute to this project, please follow these guidelines:


1. Create a new branch for your feature or bug fix.

2. Make your changes and commit them with descriptive commit messages.

3. Push your changes to your branch.

4. Submit a pull request, explaining the changes you have made.

## Testing

To run the tests for this project, use the following command:

```
cd tests
go test
```

## License

This project is licensed under the [MIT License](LICENSE).
