package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"net/http"
	"github.com/influxdb/influxdb/client/v2"
	"crypto/aes"
	"encoding/hex"
	"encoding/json"
	"crypto/cipher"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const(
	influxDbTag = "/influxdbUrl"
	portNum = "18080"
	defaultStartTime = "1970-01-01 00:00:00.000"
	timeLengthIndex int = 19
	timeStampSuffix = ".000"
)

//Note: parameter name in this struct needs to start with upper case letter
//Pod_id and Metric are required fields, the rest are optional
type json_struct struct {
	Pod_id string
	TimeStart string
	TimeEnd string
	Limit int
	Metric string
}

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}


func main() {
	http.HandleFunc(influxDbTag, influxDBHandler)
	log.Fatal(http.ListenAndServe(":"+portNum, nil))
}

func influxDBHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t json_struct   
	err := decoder.Decode(&t)
	if err != nil {
		return
	}

	log.Println(t.Pod_id)
	log.Println(t.TimeStart)
	log.Println(t.TimeEnd)
	log.Println(t.Limit)
	log.Println(t.Metric)

	podId := t.Pod_id
	startTime := t.TimeStart
	endTime := t.TimeEnd
	limitNum := strconv.Itoa(t.Limit)
	metrics := t.Metric

	if len(metrics) <= 0 {
		log.Println("no metrics")
		return
	}
	if len(podId) <= 0 {
		log.Println("no pod_id")
		return
	}

	if len(startTime) <= 0 {
		startTime = defaultStartTime
	}

	if len(endTime) <= 0 {
		x := time.Now().String()
		endTime = x[:timeLengthIndex] + timeStampSuffix
	}
	sql := "SELECT * FROM '" + metrics + "' WHERE time >= '" + startTime + "' AND time <='" + endTime + "' AND pod_id='" + podId + "'"
	if len(limitNum) <= 0 {
		sql += " LIMIT " +  limitNum
	}

	log.Println(sql)
	res,err := readInfluxDb(sql, metrics)
	if err != nil {
		return;
	}
	a, _ := json.Marshal(res)
	rw.Write(a)
	return;
	//fmt.Fprintf(rw, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func readInfluxDb(command string, metrics string) (res []client.Result, err error) {
	credential,err := getCredentials()

	if(err != nil){
		return res,err
	}

	//test username and pw

	if len(credential[0]) <= 0 {
		fmt.Println("no username")
		return res, errors.New("no username")
	}

	if len(credential[1]) <= 0 {
		fmt.Println("no password")
		return res, errors.New("no password")
	}

	if len(credential[2]) <= 0 {
		fmt.Println("no dbURL")
		return res, errors.New("no dbURL")
	}
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: credential[2],
		Username: credential[0],
		Password: credential[1],
	})
	if err != nil {
		fmt.Println("Error connecting InfluxDB Client: ", err.Error())
		return res, nil
	}
	defer c.Close()
	
	q := client.NewQuery(command, metrics, "ns")
	if response, err := c.Query(q);  err == nil{
		if err != nil {
			fmt.Println("query error")
			return res, err
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
	credential =  make([]string, 3)
	absPath, _ := filepath.Abs("influxdbUrl/credential.config")
	file, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	username := ""
	password := ""
	dbUrl := ""
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
		if strings.Contains(line, "l=") {
			dbUrl = line[2:]
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

	if len(dbUrl) <=0 {
		fmt.Println("no dbUrl")
		return credential,errors.New("no dbUrl")
	}
	realUrl,_ := decypher(dbUrl)
	fmt.Println(realUrl)

	credential[0] = realUsername
	credential[1] = realPassword
	credential[2] = realUrl
	return credential,nil
}

