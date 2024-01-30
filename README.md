eCRUD
===
A CRUD API for managing employee in-memory records

## Endpoints
### `GET /employees`
`200 OK`
```
{
    "message": "list success",
    "data": [
        {
            "id": 1,
            "firstName": "John",
            "lastName": "Doe",
            "dateOfBirth": "1985-05-15",
            "email": "john.doe@example.com",
            "isActive": true,
            "department": "Engineering",
            "role": "Software Developer"
        },
        {
            "id": 2,
            "firstName": "Jane",
            "lastName": "Smith",
            "dateOfBirth": "1990-09-22",
            "email": "jane.smith@example.com",
            "isActive": true,
            "department": "Marketing",
            "role": "Marketing Specialist"
        }
    ]
}
```
### `POST /employees`
Request sample
```
{
    "firstName": "Jane",
    "lastName": "Smith",
    "dateOfBirth": "1990-09-22",
    "email": "jane.smith@example.com",
    "isActive": true,
    "department": "Marketing",
    "role": "Marketing Specialist"
}
```
`200 OK`
```
{
    "message": "create success",
    "data": {
        "id": 1
    }
}
```
`400 Bad request`
```
{
    "message": "missing/invalid params",
    "data": {
        "fields": ["email", "dateOfBirth"]
    }
}
```
### `GET /employees/{id}`
`200 OK`
```
{
    "message": "get success",
    "data": {
        "id": 1,
        "firstName": "John",
        "lastName": "Doe",
        "dateOfBirth": "1985-05-15",
        "email": "john.doe@example.com",
        "isActive": true,
        "department": "Engineering",
        "role": "Software Developer"
    }
}
```
`404 Not found`
```
{
    "message": "record not found"
}
```
### `PUT /employees/{id}`
`200 OK`
```
{
    "message": "update success",
    "data": {
        "firstName": "Bruce",
        "lastName": "Wayne",
        "role": "CEO"
    }
}
```
### `DELETE /employees/{id}`
`200 OK`
```
{
    "message": "delete success",
    "data": {
        "id": 1
    }
}
```

## Development

:warning: This project requires at least Go 1.13. If you're running anything older, what are we doing here? ;)

### Build and run the server manually
1. `cd path/to/ecrud`
2. `go build -o server cmd/server.go`
3. `./server`