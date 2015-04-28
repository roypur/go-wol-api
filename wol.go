package main

import (
	"bytes"
	"fmt"
	"net"
	"bufio"
	"strings"
	"encoding/json"
	"encoding/binary"
	"io"
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

var path string = "";

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

    if(args[1] == "dir"){
        path = args[2];
        for i := 0;  i < 5; i++{
            args[i] = args[i+2];
        }
    }

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
    fmt.Println("\ndir <conf-dir> (optional)");
    fmt.Println("\nadd|edit|delete|server");
    fmt.Println("\nuser|host <username>|<computer-name>");
    fmt.Println("\n<password>|<mac-address>\n");
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



// This function accepts a MAC Address string, and returns a pointer to
// a MagicPacket object. A Magic Packet is a broadcast frame which
// contains 6 bytes of 0xFF followed by 16 repetitions of a given mac address.

func httpServer(){
    ln, _ := net.Listen("tcp", ":8781")
    
    for{
        conn,_ := ln.Accept()
        go parseHttp(conn)
    }

}



func makePacket(mac string) (*MagicPacket, error) {
    var packet MagicPacket
    var macAddr MACAddress

    hwAddr, err := net.ParseMAC(mac)
    if err != nil {
        return nil, err
    }

    // Copy bytes from the returned HardwareAddr -> a fixed size
    // MACAddress
    for idx := range macAddr {
        macAddr[idx] = hwAddr[idx]
    }

    // Setup the header which is 6 repetitions of 0xFF
    for idx := range packet.header {
        packet.header[idx] = 0xFF
    }

    // Setup the payload which is 16 repetitions of the MAC addr
    for idx := range packet.payload {
        packet.payload[idx] = macAddr
    }

    return &packet, nil
}

func getHost(host string)(string,error){

    file, err := ioutil.ReadFile(path + "hosts.json");
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

    file, err := ioutil.ReadFile(path + "users.json");
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

func modUser(operation string){
    var pass string;
    
    
    if(operation == "del"){
        operation = "delete";
    }
    
    
    if(operation!="delete"){
        pass = encode(string(args[4]));
    };
    
    var user string = encode(string(args[3]));
    
    
    file, err := ioutil.ReadFile(path + "users.json");
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
                hash, err = bcrypt.GenerateFromPassword([]byte(pass), 10);
                users[user] = string(hash);
            }else{
                err = errors.New("User exists")
            }
                
            if(err==nil){
                    
                b, err := json.MarshalIndent(users,"","    ");
                
                if(err==nil){
                
                    f, err := os.Create(path + "users.json")

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






func modHost(operation string){
    var mac string;
    
    if(operation == "del"){
        operation = "delete";
    }
    
    if(operation!="delete"){
        mac = string(args[4]);
    };
    
    var host string = encode(string(args[3]));
    
    
    file, err := ioutil.ReadFile(path + "hosts.json");
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
                
                    f, err := os.Create(path + "hosts.json")

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

func sendPacket(apiData map[string]string)(error){

    host,err := getHost(apiData["host"])
    if(err == nil){

        packet,_ := makePacket(host);
    
        var buf bytes.Buffer
    
        binary.Write(&buf, binary.BigEndian, packet)
    
        remoteAddr,_ := net.ResolveUDPAddr("udp", "255.255.255.255:9")
    
        connection, err := net.DialUDP("udp", nil, remoteAddr)
    
        if(err != nil){
            return err
        }else{
            _,err := connection.Write(buf.Bytes());
            if(err == nil){
                return nil
            }else{
                return err
            }
        }
    }
    return err
}

func parseHttp(conn net.Conn){

    message := bufio.NewReader(conn);
    
    header := parseHeader(*message);
    
    reply := []byte("HTTP/1.1 200 OK\r\nAccess-Control-Allow-Origin: *\r\nAccess-Control-Allow-Headers: x-api\r\nAccess-Control-Allow-Methods: POST\r\nContent-Type: text/html; charset=utf-8\r\n\r\n");
    
    var status string;  
    
    conn.Write(reply);
    
    if(len(header["x-api:"]) > 0){
    
        jsonString := []byte(header["x-api:"][0])
    
        in := make(map[string]string);

        err := json.Unmarshal((jsonString), &in);
        
        fmt.Println(in);
      
        if(isUser(in["user"], in["pass"]) == nil){
        
            if(err == nil){
                if(sendPacket(in) == nil){
                    status = "ok";
                }else{
                    status = "fail";
                }
            }    
        }else{
            status = "denied";
        }   
    }
    
    conn.Write([]byte(status))
    
    conn.Close();
}

func parseHeader(text bufio.Reader)(map[string][]string){

    str := "init"

    header := make(map[string][]string);

    for(len(str) > 3){
    
        str,_ = text.ReadString('\n');
        
        val := strings.Split(strings.TrimSpace(str), " ");
        
        header[val[0]] = val[1:]
        }
    return header;
}
