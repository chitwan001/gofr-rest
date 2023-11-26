
# GoFr REST-API

A REST API built using [GoFr](https://gofr.dev/). This project illustrates all the basics of GoLang required for building a **REST API** like **Middlewares, Authentication using JWT and Password hashing using ``bcrypt``**.


## Prerequisites

This project requires [GoLang](https://go.dev/dl/), [GoLand](https://www.jetbrains.com/go/download/) and [MySQL](https://dev.mysql.com/downloads/mysql/).
## Installation

Firstly, **we will setup the database**:

1.Create a new database and remember the name.

2.Create a table named USERS using this SQL query:
```mysql
CREATE TABLE users(
  id int AUTO_INCREMENT PRIMARY KEY,
  name varchar(255),
  email varchar(255),
  hash_pass varchar(255));
```

3.Locate ``.env`` file in the ``config`` directory and update all the environment variables with valid values for you and we are good to go.


Now for **running the project follow**:


For **first time running** after cloning the repository, run:

```bash
  npm run tidy-dev

```

For **successively running** the project **untill no new dependency is installed**, run:
```bash
  npm run dev
```



## API Reference

#### Base

```http
  GET /
```
Responds with message that server is up!

#### Register a new user

```http
  POST /user/create
```

| JSON Key | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `name` | `string` | **Required**. User's name |
| `email` | `string` | **Required**. Users's email (UNIQUE) |
| `pass` | `string` | **Required**. User's pasword |

In response, a JWT will be sent. **For all the protected routes, add Authorization header like ``Bearer <TOKEN>``**

#### Login a user

```http
  POST /user/login
```

| JSON Key | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `email`      | `string` | **Required**. User's email |
| `password`      | `string` | **Required**. User's password |


In response, a JWT will be sent. **For all the protected routes, add Authorization header like ``Bearer <TOKEN>``**

#### Get user info (protected)

```http
  GET /user/me
```

| Headers | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Authorization`      | `string` | **Required**. Bearer <JWT Token> |


#### Delete a registered user (protected)


```http
  DELETE /user/delete
```

| Headers | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Authorization`      | `string` | **Required**. Bearer <JWT Token> |


## Contributing

**LEARNING NEVER STOPS :)**

Contributions are always welcome!

Raise a new Pull Request with the suggested changes that can improve the code quality and to add other features.


## Support

For support, email [me](mailto:chitwan001@gmail.com).

Please give this repository a star ⭐ so that it reaches maximum audience. View all the star gazers [here](https://github.com/chitwan001/golang-rest/stargazers).


## References

1.[GoFr Documentation](https://gofr.dev/docs)