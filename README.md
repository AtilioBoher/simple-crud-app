# Create user

```
curl -v -X POST -H "content-type: application/json" http://localhost:8080/users/create -d '{"name": "Fulanito", "email": "fulanito@mail.com", "age": 21}'
```

# Read user

```
curl -v http://localhost:8080/users/Fulanito
```

# Update user

```
curl -v -X PATCH -H "content-type: application/json" http://localhost:8080/users/Fulanito -d '{"email": "newEmail@mail.com", "age": 40}'
```

# Delete user

```
curl -v -X DELETE  http://localhost:8080/users/Fulanito
```
