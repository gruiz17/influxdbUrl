package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"net/url"
	"github.com/influxdb/influxdb/client/v2"
	"github.com/gorilla/mux"
)

const(
	influxDbTag = "/hello"
	portNum = "8080"
	sqlTag = "q="
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc(influxDbTag + "/{rest:.*}", influxDBHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":" + portNum, router))
}

func influxDBHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /influxDbTag request")
	log.Println(r.UserAgent())

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "your influxDB result:", r)

	t := r.URL.String()
	t = t [7:len(t)]
	fmt.Fprintln(w, "t=" + t)

	x,_ := url.QueryUnescape(t)
	//add a '/' after http:/
	if(strings.Contains(x, "http:/") &&  !strings.Contains(x, "http://")){
		slashIndex := strings.Index(x, "http:/") + 5
		x = x[:slashIndex] + "/" + x[slashIndex:]
	}
	readInfluxDb(x);
}


func readInfluxDb(command string) (res []client.Result, err error) {
	input := command
	if(strings.Contains(input, "http:/") &&  !strings.Contains(input, "http://")){
		slashIndex := strings.Index(input, "http:/") + 5
		input = input[:slashIndex] + "/" + input[slashIndex:]
	}
	fmt.Println("input=" + input)
	elem := strings.Split(input, " ")
	urlInput := ""
	sqlInput := ""
	for k, _ := range elem{
		if strings.HasPrefix(elem[k], "'http"){
			urlInput = elem[k][1:len(elem[k]) - 1]
			fmt.Println(urlInput)
			break;
		}
	}

	if(strings.Contains(input, sqlTag)){
		sqlInput = input [strings.Index(input, sqlTag)+2:len(input) - 1]
		fmt.Println(sqlInput)
	}

	if len(urlInput) <= 0 {
		fmt.Println("invalid url")
		return res, nil
	}

	if len(sqlInput) <= 0 {
		fmt.Println("invalid sql")
		return res, nil
	}

	sqlElem := strings.Split(sqlInput, " ")

	tableName := ""
	for k, _ := range sqlElem{
		if strings.HasPrefix(sqlElem[k], "from"){
			tableName = sqlElem[k+1]
			fmt.Println(tableName)
			break;
		}
	}

	if len(tableName) <= 0 {
		fmt.Println("cannot find table name")
		return res, nil
	}

	urlAddr := urlInput [:strings.Index(urlInput, "?")]
	fmt.Println("urlAddr:" + urlAddr)

	u, err := url.Parse(urlInput)
	if err != nil {
		log.Fatal(err)
	}

	m, _ := url.ParseQuery(u.RawQuery)

	//debug
	/************/
	fmt.Println(m)

	for k, _ := range m {
		fmt.Println(k,m[k][0])
	}
	/************/


	//test username and pw
	if _, ok := m["u"]; !ok {
		fmt.Println("no username")
		return res, nil
	}

	if len(m["u"][0]) <= 0 {
		fmt.Println("no username")
		return res, nil
	}

	if _, ok := m["p"]; !ok {
		fmt.Println("no password")
		return res, nil
	}

	if len(m["p"][0]) <= 0 {
		fmt.Println("no password")
		return res, nil
	}

	/*
	 * 1. if it's delete or drop or remove
	 * 2. add/create is not allowed
	 */

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: urlAddr,
		Username: m["u"][0],
		Password: m["p"][0],
	})
	if err != nil {
		fmt.Println("Error connecting InfluxDB Client: ", err.Error())
		return res, nil
	}
	defer c.Close()
	
	q := client.NewQuery(sqlInput, tableName, "ns")
	if response, err := c.Query(q);  err == nil{
		if response.Error() != nil {
			fmt.Println("query error")
			return res, response.Error()
		}
		fmt.Println(response.Results)
		res = response.Results
	}
	return res, nil
}
