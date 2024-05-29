
# Run it

Start the server: 

```bash
go run . 
```


```bash
curl http://localhost:8090/add\?key=niklas
```

pass payload as well

```bash
curl -d '{"login":"my_login","password":"my_password"}' -H "Content-Type: application/json" http://localhost:8080/add\?key=niklas
```

# 
