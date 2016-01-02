package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"github.com/roypur/goapi/src"
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
        goapi.Listen(httpHandler, config["listen"], config["cert"], config["key"]);
    }else{
        showHelp();
    }
}

func httpHandler(req goapi.Request){
        
        var status string;
        
        var user string = req.Header["Goapi-User"]
        
        var pass string = req.Header["Goapi-Pass"]
        
        var host string = req.Header["Goapi-Host"]
        
        var list bool = req.Header["Goapi-List"] == "true"
        
        
        if(isUser(user, pass) == nil){
            if(list){                            
                //var jsonList string;
                
                var list []string = []string{}
                
                for k,_ := range getHosts(){
                    list = append(list,k)
                }
                
                b,err := json.Marshal(list)

                if(err==nil){
                    status = string(b)
                }else{
                    status = "no-hosts-found"
                }
            }else{
                if(sendPacket(host) == nil){
                    status = "ok";
                }else{
                    status = "fail";
                }
            }
        }else{
            status = "denied";
        }
        req.Write(goapi.Ok)
        req.Write(req.Resp)
        req.Write("\r\n\r\n")
        req.Write(status)
        req.Close()
        
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

func getHost(host string)(string,error){

    hosts := getHosts()
       
    if(len(hosts[host]) > 0){    
        return hosts[host], nil;
    }
    return "", errors.New("host not found" + host); 
}

func getHosts()(map[string]string){

    file, err := ioutil.ReadFile(configPath + "/hosts.json");
        
    hosts := make(map[string]string);

    if(err == nil){
        err = json.Unmarshal(file, &hosts);
        
        if(err == nil){
            return hosts;
        }
    }
    return nil;
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
