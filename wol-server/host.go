package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"io"
	"io/ioutil"
	"errors"
	"os"
)

func modHost(operation string){
    var mac string;
    
    if(operation == "del"){
        operation = "delete";
    }
    
    if(operation!="delete"){
        mac = string(args[4]);
    };
    
    var host string = encode(string(args[3]));
    
    
    file, err := ioutil.ReadFile(configPath + "/hosts.json");
    hosts := make(map[string]string);
    
    if(err == nil){
        err = json.Unmarshal(file, &hosts);
        if(err == nil){
            
            if(operation=="delete"){
                if(len(hosts[host]) == 0){
                    fmt.Println("Host<" + args[3] + "> Doesn't exist");
                    os.Exit(-1);
                }
                delete(hosts,host);
            }else if(len(hosts[host]) == 0){
                hosts[host] = string(mac);
            }else{
                err = errors.New("Host exists")
            }
                
            if(err==nil){
                    
                b, err := json.MarshalIndent(hosts,"","    ");
                
                if(err==nil){
                
                    f, err := os.Create(configPath + "/hosts.json")

                    if(err != nil){
                        fmt.Println(err);
                    }else{
                        _,err := io.WriteString(f, string(b));
                        
                        if(err == nil){
                            
                            var grammar string = strings.Replace(operation + "ed","ee","e",-1);
                            
                            fmt.Println("host<" + args[3] + "> " + grammar);
                        }
                    }
                }
            }else{
                fmt.Println(err);
            }
        }
    }
}

