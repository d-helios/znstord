package main

import (
	"encoding/json"
	"flag"
	"github.com/d-helios/znstor"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)


const (
	defaultConfigFile = "/etc/znstor/config.json"
	defaultLogFile = "/var/adm/znstord.log"
)

func main() {

	// configuration file in json format
	configFile := flag.String("config", defaultConfigFile, "znstord configuration file in json format")
	logFile := flag.String("log", defaultLogFile, "znstord log file")
	flag.Parse()

	// read configuration
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// try to open logfile
	lf, err := os.Create(*logFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	var config znstor.ServerData

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// load routes
	router := znstor.NewRouter(io.Writer(lf), config.Auth.UserName, config.Auth.UserPassword)

	// start http server
	log.Fatal(http.ListenAndServe(config.Listen + ":10987", router))
}
