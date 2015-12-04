# influxdbUrl

This is a helper project to let frontend code connect to backend influxdb databse.
By default, this service is bind to port 18080

follow the instructions here to install go on your environment:
https://golang.org/doc/install



build all go file in influxdbUrl project
```
go install influxdbUrl/
```

If you have a username and password to encrypt.

```
go run influxdbUrl/encryptionGenerator.go ${YOUR_USERNAME}
go run influxdbUrl/encryptionGenerator.go ${YOUR_PASSWORD}
```
This will generate a encryted version of your username and password

Copy and paste the result generated from the above program to credential.config


Use the following command to start the backend service
```
go run influxdbUrl/influxDBUrl2.go
```

Note: because this program use current working path for config file, it can only work correctly using the above relative path. Do not run it from a different path


http syntax:
```
curl -H "Content-Type: application/json" -X POST -d "{\"pod_id\": \"sdfd\", \"timeStart\": \"2015-11-03\",\"timeEnd\": \"2015-11-04\",\"limit\": 1000,\"metric\": \"uptime_ms_cumulative\" }"    http://localhost:18080/influxdbUrl
```

Note: after the above config, the last step is to change the file permmision of credential.config to read only
```
chmod 444 credential.config
```