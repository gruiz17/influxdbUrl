package main

import (
    "bufio"
	"fmt"
	"log"
	"strings"
	"net/http"
	"net/url"
	"github.com/influxdb/influxdb/client/v2"
	"github.com/gorilla/mux"
	"crypto/aes"
	"encoding/hex"
	"crypto/cipher"
    "errors"
    "os"
    "path/filepath"
)

const(
	influxDbTag = "/influxdb"
	portNum = "18080"
	sqlTag = "q="
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}


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
	readInfluxDb(x)
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

    credential,err := getCredentials()

    if(err == nil){
    	return res,err
    }

	//test username and pw

	if len(credential[0]) <= 0 {
		fmt.Println("no username")
		return res, nil
	}

	if len(credential[1]) <= 0 {
		fmt.Println("no password")
		return res, nil
	}

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: urlAddr,
		Username: credential[0],
		Password: credential[1],
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

func decypher(command string) (s string,err error){
	key_text := "astaxie12798akljzmknm.ahkjkljl;k"

	// Create the aes encryption algorithm 
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
		os.Exit(-1)
	}

	//fmt.Print("here1")
	ciphertext2,_ := hex.DecodeString(command)
	//fmt.Print("here2")
	// Decrypt strings
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	//fmt.Print("here3")
	plaintextCopy := make([]byte, len(ciphertext2))
	//fmt.Print("here4")
	cfbdec.XORKeyStream(plaintextCopy, ciphertext2)
	//fmt.Print("here5")
	fmt.Printf("%x=>%s\n", ciphertext2, plaintextCopy)
	s = string(plaintextCopy[:])
	return s, nil
}

func getCredentials() (credential []string,err error) {
    credential =  make([]string, 2)
    absPath, _ := filepath.Abs("influxdbUrl/credential.config")
    file, err := os.Open(absPath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    username := ""
    password := ""
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "u=") {
            username = line[2:]
            continue;
        }
        if strings.Contains(line, "p=") {
            password = line[2:]
            continue;
        }
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    //test username and pw
    if len(username) <= 0 {
        fmt.Println("no username")
        return credential,errors.New("no username")
    }

    realUsername,_ := decypher(username)

    fmt.Println(realUsername)

    if len(password) <= 0 {
        fmt.Println("no password")
        return credential,errors.New("no password")
    }
    realPassword,_ := decypher(password)
    fmt.Println(realPassword)

    credential[0] = realUsername
    credential[1] = realPassword
    return credential,nil
}

