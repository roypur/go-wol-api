package main

import (
	"fmt"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var config map[string]string;
var configPath string = "";

func setPath(){
    if(args[1] == "dir"){
        configPath,_ = filepath.Abs(args[2]);
        for i := 0;  i < 5; i++{
            args[i] = args[i+2];
        }
    }else{
        configPath,_ = filepath.Abs("");
    }
}

func getConf(){

    _,confErr := os.Stat(configPath + "/config.json");
    _,hostErr := os.Stat(configPath + "/hosts.json");
    _,userErr := os.Stat(configPath + "/users.json");
    
    config = make(map[string]string);
    
    if(os.IsNotExist(userErr)){
        
        f, err := os.Create(configPath + "/users.json")
        if(err==nil){
        
            var content string = "{}";
            
            _,err = io.WriteString(f, content);
            if(err!=nil){
                fmt.Println(err);       
            }
        }else{
            fmt.Println(err);
        }
    }
    if(os.IsNotExist(hostErr)){
        
        f, err := os.Create(configPath + "/hosts.json")
        if(err==nil){

            var content string = "{}";
        
            _,err = io.WriteString(f, content);
            if(err!=nil){
                fmt.Println(err);       
            }
        }else{
            fmt.Println(err);
        }
        
    }    
    if(os.IsNotExist(confErr)){
        
        f, err := os.Create(configPath + "/config.json")
        if(err==nil){    
        
            config["listen"] = ":8781";
            config["work"] = "10";
        
            b, err := json.MarshalIndent(config,"","    ");
        
            _,err = io.WriteString(f, string(b));
            if(err!=nil){
                fmt.Println(err);       
            }
        }else{
            fmt.Println(err);
        }
    }else{
        file, err := ioutil.ReadFile(configPath + "/config.json");
        if(err==nil){
            err = json.Unmarshal(file, &config);
            if(err!=nil){
                fmt.Println(err);       
            }
        }else{
            fmt.Println(err);
        }
    }
}

