# Post user

Para cargar a Fulanito como usuario se utiliza el método POST: 
```
curl -v -X POST -H "content-type: application/json" http://localhost:8080/users/create -d '{"name": "Fulanito", "email": "fulanito@mail.com", "age": 21}'
```

# Get user

```
curl -v http://localhost:8080/users/Fulanito
```

# Patch user

```
curl -v -X PATCH -H "content-type: application/json" http://localhost:8080/users/Fulanito -d '{"email": "newEmail@mail.com", "age": 40}'
```

# Patch user

```
curl -v -X DELETE  http://localhost:8080/users/Fulanito
```
