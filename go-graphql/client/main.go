package main

import (
	"context"
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/shurcooL/graphql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var Conf conf

type conf struct {
	SkyUrl string `mapstructure:"sky-url" yaml:"sky-url"`
}

func main() {
	Viper()
	err := run()
	if err != nil {
		log.Println(err)
	}
}

func run() error {

	client := graphql.NewClient(Conf.SkyUrl, nil)
	log.Info(Conf.SkyUrl)
	/*
		query {
		  queryBasicTraces(condition: {
		    serviceId: "YmZm.1",
		    queryDuration: {
		      start: "2023-01-10 0643",
		      end: "2023-01-10 0713",
		      step: MINUTE
		    }
		    traceState: ALL
		    queryOrder: BY_START_TIME
		    paging: {pageSize: 20, pageNum: 1}
		 }){
		  traces{
		    segmentId
		    endpointNames
		    duration
		    start
		    isError
		    traceIds
		  }
		  }
		}
	*/
	type BasicTrace struct {
		segmentId     graphql.String
		endpointNames graphql.String
		duration      graphql.Int
		start         graphql.String
		traceIds      graphql.String
	}
	var q struct {
		queryBasicTraces struct {
			//traces graphql.String `graphql:"traces:{key: segmentId endpointNames  duration start isError  traceIds}"`
			traces []struct {
				segmentId     graphql.String
				endpointNames graphql.String
				duration      graphql.Int
				start         graphql.String
				traceIds      graphql.String
			}
			//duration      graphql.Int
			//start         graphql.String
			//traceIds      graphql.String
			//traces struct {
			//	segmentId     graphql.String
			//	endpointNames graphql.String
			//	duration      graphql.Int
			//	start         graphql.String
			//	traceIds      graphql.String
			//} `graphql:"traces"`
		} `graphql:"queryBasicTraces(condition: {serviceId:\"YmZm.1\",queryDuration:{start:\"2023-01-10 0643\",end:\"2023-01-10 0713\",step:MINUTE}traceState:ALL,queryOrder:BY_START_TIME,paging:{pageSize:200,pageNum:1}})"`
	}

	//variables := map[string]interface{}{
	//	"condition": struct {
	//		queryDuration struct {
	//			start string
	//			end   string
	//			step  string
	//		}
	//		traceState string
	//		queryOrder string
	//		paging     struct {
	//			pageNum  int
	//			pageSize int
	//		}
	//	}{},
	//	"id": graphql.String("11"),
	//}

	err := client.Query(context.Background(), &q, nil)
	if err != nil {
		return err
	}
	print(q)
	return nil
}

// Viper is a method to read the contents of a configuration file into a Struct
// Author [Silent Mo]
func Viper() {

	v := viper.New()
	v.SetConfigFile("config.yaml")
	//v.AddConfigPath("")
	//v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Errorf("fatal error config file: %s", err)
		panic("fatal error config file")
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config file changed:%v", e.Name)
	})
	if err := v.Unmarshal(&Conf); err != nil {
		log.Error(err)
		panic("viper unmarshal error")
	}
	//_ = v.WriteConfig()
	//log.Info(Conf)
	//return v

}

// print pretty prints v to stdout. It panics on any error.
func print(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
