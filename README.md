# influxdbUrl

This is a helper project to let frontend code connect to backend influxdb databse.
By default, this service is bind to port 18080


```
go install influxdbUrl/
go run influxdbUrl/influxDBUrl.go
```

http syntax:
```
http://localhost:8080/influxdb/curl%20-G%20%27${Your URL}?u=${DATABASE USER NAME}&p={DATABASE PASSWORD}&pretty=true%27%20--data-urlencode%20"q=${sql command}"
```

sample link:
```
http://localhost:8080/influxdb/curl%20-G%20%27abc.com/x/y?u=uname&p=pw&pretty=true%27%20--data-urlencode%20"q=select%20*%20from%20pod_id%20limit%201"
```
