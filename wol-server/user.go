package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"io"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)


func modUser(operation string){
    var pass string;
    
    
    if(operation == "del"){
        operation = "delete";
    }
    
    
    if(operation!="delete"){
        pass = args[4];
    };
    
    var user string = args[3];
    
    
    file, err := ioutil.ReadFile(configPath + "/users.json");
    users := make(map[string]string);
    
    if(err == nil){
        err = json.Unmarshal(file, &users);
        if(err == nil){
        
            var hash []byte;
            
            if(operation=="delete"){
                if(len(users[user]) == 0){
                    fmt.Println("User<" + args[3] + "> Doesn't exist");
                    os.Exit(-1);
                }
                delete(users,user);
                
            }else if(len(users[user]) == 0){
                
                work,_:= strconv.Atoi(config["work"]);
            
                hash, err = bcrypt.GenerateFromPassword([]byte(pass), work);
                users[user] = string(hash);
            }else{
                err = errors.New("User exists")
            }
                
            if(err==nil){
                    
                b, err := json.MarshalIndent(users,"","    ");
                
                if(err==nil){
                
                    f, err := os.Create(configPath + "/users.json")

                    if(err != nil){
                        fmt.Println(err);
                    }else{
                        _,err := io.WriteString(f, string(b));
                        
                        if(err == nil){
                            
                            var grammar string = strings.Replace(operation + "ed","ee","e",-1);
                            
                            fmt.Println("user<" + args[3] + "> " + grammar);
                        }
                    }
                }
            }else{
                fmt.Println(err);
            }
        }
    }
}


