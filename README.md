## Getting started

Start docker compose:
- Api server
- Postgres
- Grafana
- Prometheus

```docker compose up -d --build```

Start the server in develop mode:

```
cd backend
go run main.go
```

Start the client:

```
cd frontend
npm start
```


Login:

``` 
curl --location 'http://localhost:9393/api/login' \
--header 'Content-Type: application/json' \
--data '{
    "username":"user",
    "password":"test"
}'
```

Use token to authenticate requests

Add redfish host:
```
curl --location 'http://localhost:9393/api/v1/redfish/add' \
--header 'Content-Type: application/json' \
--header 'Authorization: $TOKEN' \
--data '{
    "Host":"192.168.1.125",
    "User":"worker",
    "Password":"123"
}'
```

Add vCenter host:
```
curl --location 'http://localhost:9393/api/v1/vsphere/add' \
--header 'Content-Type: application/json' \
--header 'Authorization: $TOKEN' \
--data '{
    "Host":"https://vcenter.host/sdk",
    "User":"user",
    "Password":"123"
}'
```

Check Grafana dashboard default creds admin:admin on localhost:3000