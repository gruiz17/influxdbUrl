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

Copy and paste the result generated from the above program to encryptionGenerator.config


Use the following command to start the backend service
```
go run influxdbUrl/influxDBUrl2.go
```

Note: because this program use current working path for config file, it can only work correctly using the above relative path. Do not run it from a different path


http syntax:
```
http://localhost:18080/influxdb/curl%20-G%20%27${Your URL}?u=${DATABASE USER NAME}&p={DATABASE PASSWORD}&pretty=true%27%20--data-urlencode%20"q=${sql command}"
```

sample link:
```
http://localhost:18080/influxdb/curl%20-G%20%27abc.com/x/y?u=uname&p=pw&pretty=true%27%20--data-urlencode%20"q=select%20*%20from%20pod_id%20limit%201"
```
