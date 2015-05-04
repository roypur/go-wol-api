package main

import (
	"fmt"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
        configPath,_ = filepath.Abs("config");
    }
}

func getConf(){

    _,confErr := os.Stat(configPath + "/config.json");
    _,hostErr := os.Stat(configPath + "/hosts.json");
    _,userErr := os.Stat(configPath + "/users.json");
    
    var accept bool = false; 
    
    var folder bool = false;
    
    if(os.IsNotExist(confErr) || os.IsNotExist(hostErr) || os.IsNotExist(userErr)){
        var answer string;
        
        fmt.Println("\nThe following config files are missing:\n");
        
        if(os.IsNotExist(confErr)){
            fmt.Println("config.json")
        }
        if(os.IsNotExist(confErr)){
            fmt.Println("users.json")
        }        
        if(os.IsNotExist(confErr)){
            fmt.Println("hosts.json")
        }
        
        fmt.Println("\nDo you want them to be automatically created in the following directory?\n");
        fmt.Println(configPath + "/");
        fmt.Println("\n(Yes/No)\n");
        fmt.Scanln(&answer);
        
        fmt.Println();
        
        answer = strings.Trim(answer," ");
        
        if((strings.ToLower(answer)=="y") || (strings.ToLower(answer)=="yes")){
            err := os.MkdirAll(configPath, 0755);
            accept = true;
            if(err==nil){
                folder = true;
            }else{
                fmt.Println(err);
            }
        }
    }
    
    
    
    config = make(map[string]string);
    
    
    
    if(os.IsNotExist(userErr) && accept && folder){
        
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
    if(os.IsNotExist(hostErr) && accept && folder){
        
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
    if(os.IsNotExist(confErr) && accept && folder){
        
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
    }else if(accept == false && !os.IsNotExist(confErr)){
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

