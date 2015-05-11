package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type MACAddress [6]byte

// A MagicPacket is constituted of 6 bytes of 0xFF followed by
// 16 groups of the destination MAC address.

type MagicPacket struct {
    header  [6]byte
    payload [16]MACAddress
}

var args map[int]string;


func main(){
    
    args = make(map[int]string);
    
    for i := 0;  i < 7; i++{
        if(len(os.Args) > i){
            args[i] = os.Args[i];
        }else{
            args[i] = "";
        }
    }
    
    setPath();
    
    getConf();

    if((args[1] == "delete") || (args[1] == "del")){
        if(args[2] == "user"){
            modUser(args[1]);
        }else if(args[2] == "host"){
            modHost(args[1]);
        }else{
            showHelp();
        }     
    }else if(((args[1] == "add") || (args[1] == "edit")) && ((len(args[3]) > 0) && (len(args[4]) > 0))){
        if(args[2] == "user"){
            modUser(args[1]);
        }else if(args[2] == "host"){
            modHost(args[1]);
        }else{
            showHelp();
        }  
    }else if(args[1] == "server"){
        httpServer();
    }else{
        showHelp();
    }
}



func showHelp(){
    fmt.Println("\n");
    fmt.Println("# # # # # # # # # # # # # # # # # # # # # #");
    fmt.Println("#                                         #");
    fmt.Println("#  Usage:                                 #");
    fmt.Println("#                                         #");
    fmt.Println("#  dir <conf-dir> (optional)              #");
    fmt.Println("#                                         #");
    fmt.Println("#  add|edit|delete|server                 #");
    fmt.Println("#                                         #");
    fmt.Println("#  user|host <username>|<computer-name>   #");
    fmt.Println("#                                         #");
    fmt.Println("#  <password>|<mac-address>               #");
    fmt.Println("#                                         #");
    fmt.Println("# # # # # # # # # # # # # # # # # # # # # #");
    fmt.Println("\n");
}


func encode(raw string)(string){

    translate := make(map[string]string);
    
    var out string = strings.Replace(raw, "-", "-a", -1);

    translate["["] = "-b"
    translate["]"] = "-c"
    translate["{"] = "-d"
    translate["}"] = "-e"
    translate["'"] = "-f"
   translate["\""] = "-g"
    translate[" "] = "-h"
    translate[","] = "-i"
    translate[":"] = "-j"
    
    for k, v := range translate {
        out = strings.Replace(out, k, v, -1);
        //fmt.Println(out);
    }
    return out;
}

func getHost(host string)(string,error){

    file, err := ioutil.ReadFile(configPath + "/hosts.json");
    hosts := make(map[string]string);
    
    if(err == nil){
        err = json.Unmarshal(file, &hosts);
        if((err == nil)  && (len(hosts[host]) > 0)){
            
            return hosts[host], nil;
        }
    }
    return "", errors.New("host not found" + host);
    
    
}

func isUser(user string, pass string)(error){

    file, err := ioutil.ReadFile(configPath + "/users.json");
    users := make(map[string]string);
    
    if(err == nil){
        err = json.Unmarshal(file, &users);
        if((err == nil)  && (len(users[user]) > 0)){
            
            err = bcrypt.CompareHashAndPassword([]byte(users[user]), []byte(pass));
            
            if(err == nil){            
                return nil;
            }
        }
    }
    return errors.New("user not found" + user);
 
}
