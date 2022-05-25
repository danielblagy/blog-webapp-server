# REST API Documentation

## Contents
* [Data Structures](#data-structures)
	* [User](#user)
	* [Article](#article)
* [/users endpoint](#users)
	* [Get all users](#get-all-users)
	* [Get user by id](#get-user-by-id)
	* [Sign up user](#sign-up-user)
	* [Sign in user](#sign-in-user)
	* [Refresh User Tokens](#refresh-user-tokens)
	* [Get my data](#get-my-data)
	* [Update my data](#update-my-data)
	* [Delete my data](#delete-my-data)

## Data structures

### User

| Name | Type | Description |
| --- | --- | --- |
| id | int | Primary key. |
| login | string | Each user has a unique login, max length is 100 characters. |
| fullname | string | First and last name of the user, max length is 300 characters. |
| articles | []Article | An array of articles written by the user. |

JSON Example of User object

```json
{
    "id": 8,
    "login": "jazz_oogie",
    "fullname": "Jamie Dunkin",
    "articles": [
        {
            "id": 5,
            "author_id": 8,
            "title": "my article",
            "content": "hello :)",
            "published": true
        }
    ]
}
```

### Article

| Name | Type | Description |
| --- | --- | --- |
| id | int | Primary key. |
| author_id | int | ID of the user who owns the article. |
| title | string | Title must be unique relative to other articles of the user, max length is 300 characters. |
| content | string | The content of the article. |
| published | boolean | If true, it's public and can be read by other users, it's private otherwise. |

JSON Example of Article object

```json
{
    "id": 1,
    "author_id": 10,
    "title": "New Title",
    "content": "Updated content.",
    "published": true
}
```

## users/

### *Get all users*
### GET users/

#### Response

| Case | Status | Body |
| --- | --- | --- |
| On success | `200 OK` | An array named 'users' of User objects |
| On failure | `404 Not Found` | `{ "message": [error message] }` |

#### Example

Response on successful retrieval (`200 OK`)
```json
{
    "users": [
        {
            "id": 8,
            "login": "jazz_oogie",
            "fullname": "Jamie Dunkin",
            "articles": null
        },
        {
            "id": 10,
            "login": "danielblagy",
            "fullname": "Daniel Blagy",
            "articles": null
        },
        {
            "id": 11,
            "login": "tomalberto",
            "fullname": "Thomas Alberto",
            "articles": null
        },
        {
            "id": 12,
            "login": "sergey",
            "fullname": "Sergey Urtugov",
            "articles": null
        },
        {
            "id": 13,
            "login": "tomsanders",
            "fullname": "Tom Sanders",
            "articles": null
        }
    ]
}
```

### *Get user by id*
### GET users/:id

#### Request

`id` must correspond to user id

#### Response

| Case | Status | Body |
| --- | --- | --- |
| On success | `200 OK` | User object |
| On failure | `404 Not Found` | `{ "message": "record not found" }` |

#### Example

Request GET users/8

Response on success (`200 OK`)
```json
{
    "id": 8,
    "login": "jazz_oogie",
    "fullname": "Jamie Dunkin",
    "articles": [
        {
            "id": 5,
            "author_id": 8,
            "title": "my article",
            "content": "hello :)",
            "published": true
        }
    ]
}
```

### *Sign up user*
### POST users/signup

#### Request

Request body structure (example)

```json
{
    "login": "danielblagy",
    "fullname": "Daniel Blagy",
    "password": "danielblagypassword"
}
```

login, fullname, and password must not be empty strings.

#### Response

| Case | Status | Body |
| --- | --- | --- |
| On success | `201 Created` | User object of newly created user. |
| Not all required fields provided | `400 Bad Request` | `{ "message": "invalid user data" }` |
| Login is taken | `409 Conflict` | `{ "message": "this login is taken" }` |
| On failure | `500 Internal Server Error` | `{ "message": [server error] }` |

#### Example

Request POST users/signup

Request body
```json
{
    "login": "johnpeterson",
    "fullname": "John Peterson",
    "password": "johnp"
}
```

Response on success (`201 Created`)
```json
{
    "id": 14,
    "login": "johnpeterson",
    "fullname": "John Peterson",
    "articles": null
}
```

### *Sign in user*
### POST users/signin

#### Request

Request body structure (example)

```json
{
    "login": "danielblagy",
    "password": "danielblagypassword"
}
```

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `200 OK` | `{ "access_token": [], "refresh_token": [] }` |
| Request body is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| No user with login | `404 Not Found` | `{ "message": "user with this login doesn't exist" }` |
| Incorrect password | `401 Unauthorized` | `{ "message": [error message] }` |
| Server error | `500 Internal Server Error` | `{ "message": [server error] }` |

#### Example

Request POST users/signin

Request body
```json
{
    "login": "johnpeterson",
    "password": "johnp"
}
```

Response on success (`200 OK`)
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTMzODQ5ODEsImp0aSI6IjE0In0.UPMVXrXya9McDhECCE1OrnPayya6UQFvqtU67MdIJBE",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTUxOTg0ODEsImp0aSI6IjE0In0.le7yayzG8U-JHm0Vd2O8uR1doWaCKbCV899_8qQZih8"
}
```

### *Refresh User Tokens*
### POST users/refresh

User must be signed in.

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `200 OK` | `{ "access_token": [], "refresh_token": [] }` |
| Refresh Token cookie doesn't exist | `400 Bad Request` | `{ "message": [error message] }` |
| Not logged it / Refresh Token has expired | `401 Unauthorized` | `{ "message": [error message] }` |
| Server error | `500 Internal Server Error` | `{ "message": [server error] }` |

#### Example

Request POST users/refresh

Response on success (`200 OK`)
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTMzODUyMjksImp0aSI6IjE0In0._RHg0S8TyWzNcXM5FjWzR6gHOA14Bq9YQexnw4uYfck",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTUxOTg3MjksImp0aSI6IjE0In0.p4AaMdA7EA6jsCRqWTkhBzY-_ZqKfFSE8reBLj72ylY"
}
```

### *Get my data*
### GET users/me

User must be signed in.

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `200 OK` | User object |
| Access Token is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Couldn't get user by id / User doesn't exist | `404 Not Found` | `{ "message": [error message] }` |
| Not logged it / Access Token has expired | `401 Unauthorized` | `{ "message": [error message] }` |
| Server error | `500 Internal Server Error` | `{ "message": [server error] }` |

#### Example

Request GET users/me

Response on success (`200 OK`)
```json
{
    "id": 14,
    "login": "johnpeterson",
    "fullname": "John Peterson",
    "articles": []
}
```

### *Update my data*
### PUT users/

User must be signed in.

#### Request

Request body structure (example)

```json
{
    "fullname": "Daniel Updated Blagy",
    "password": "myupdatedpassword"
}
```

Only `fullname` and `password` fields of User object can be updated.

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `200 OK` | User object |
| Request body is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Access Token is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Not logged it / Access Token has expired | `401 Unauthorized` | `{ "message": [error message] }` |
| Couldn't get user with id / User doesn't exist | `404 Not Found` | `{ "message": [error message] }` |
| Server error | `500 Internal Server Error` | `{ "message": [server error] }` |

#### Example

Request PUT users/

Request body
```json
{
    "fullname": "John Derek Peterson",
    "password": "johnp2"
}
```

Response on success (`200 OK`)
```json
{
    "id": 14,
    "login": "johnpeterson",
    "fullname": "John Derek Peterson",
    "articles": null
}
```

### *Delete my data*
### DELETE users/

User must be signed in.

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `200 OK` | User object |
| Request body is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Access Token is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Not logged it / Access Token has expired | `401 Unauthorized` | `{ "message": [error message] }` |
| Couldn't get user with id / User doesn't exist | `404 Not Found` | `{ "message": [error message] }` |
| Server error | `500 Internal Server Error` | `{ "message": [server error] }` |

#### Example

Request DELETE users/

Response on success (`200 OK`)
```json
{
    "id": 14,
    "login": "johnpeterson",
    "fullname": "John Derek Peterson",
    "articles": []
}
```

## /articles endpoint

### *Get articles*
### GET articles/

Get all published articles.

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `200 OK` | An array of Article objects. |
| Failure | `404 Not Found` | `{ "message": [error message] }` |

#### Example

Request GET articles/

Response on success (`200 OK`)
```json
[
    {
        "id": 5,
        "author_id": 8,
        "title": "my article",
        "content": "hello :)",
        "published": true
    },
    {
        "id": 1,
        "author_id": 10,
        "title": "New Title",
        "content": "Updated content.",
        "published": true
    }
]
```

### *Get article by id*
### GET article/:id

If the user is authorized, they can access their private articles.

#### Request

`id` must correspond to article id.

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `200 OK` | Article object |
| Failure / Accessing private article while unauthorized | `404 Not Found` | `{ "message": "record not found" }` |
| No access to a private article when authorized | `401 Unauthorized` | `{ "message": "article is private" }` |

#### Example

Request GET articles/5

Response on success (`200 OK`)
```json
{
    "id": 5,
    "author_id": 8,
    "title": "my article",
    "content": "hello :)",
    "published": true
}
```

### *Create article*
### POST articles/

User must be signed in.

#### Request

Request body structure (example)

```json
{
    "title": "Green Leopards",
    "content": "Have u seen them?",
    "published": false
}
```

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `201 Created` | Article object of newly created article. |
| Request body is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Access Token is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Not logged it / Access Token has expired | `401 Unauthorized` | `{ "message": [error message] }` |
| User already has article with that title | `409 Conflict` | `{ "message": "user already has article with this title" }` |
| Server Error | `500 Internal Server Error` | `{ "message": [server error] }` |

#### Example

Request POST articles/

Request body
```json
{
    "title": "Green Leopards",
    "content": "Have u seen them?",
    "published": false
}
```

Response on success (`201 Created`)
```json
{
    "id": 8,
    "author_id": 12,
    "title": "Green Leopards",
    "content": "Have u seen them?",
    "published": false
}
```

### *Update article*
### PUT articles/:id

User must be signed in.

`id` must correspond to article id.

#### Request

Request body structure (example)

```json
{
    "title": "Title updated",
    "content": "Content updated",
    "published": true
}
```

Only `title`, `content`, and `published` fields of Article object can be updated.

#### Response

| Case | Status | Body |
| --- | --- | --- |
| Success | `200 OK` | Article object |
| Request body is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Access Token is invalid | `400 Bad Request` | `{ "message": [error message] }` |
| Not logged it / Access Token has expired | `401 Unauthorized` | `{ "message": [error message] }` |
| User doesn't own the article | `401 Unauthorized` | `{ "message": "access denied" }` |
| Couldn't get article with id / Article doesn't exist / Failure | `404 Not Found` | `{ "message": [error message] }` |
| Server error | `500 Internal Server Error` | `{ "message": [server error] }` |

#### Example

Request PUT articles/8

Request body
```json
{
    "content": "Have u seen them? I bet you haven't.",
    "published": true
}
```

Response on success (`200 OK`)
```json
{
    "id": 8,
    "author_id": 12,
    "title": "Green Leopards",
    "content": "Have u seen them? I bet you haven't.",
    "published": true
}
```