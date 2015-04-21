package main

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
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


func main(){
    
    if(os.Args[1] == "add"){
        modUser("add");
    }else if(os.Args[1] == "edit"){
        modUser("edit");
    }else if(os.Args[1] == "server"){
        httpServer();
    };
}

// This function accepts a MAC Address string, and returns a pointer to
// a MagicPacket object. A Magic Packet is a broadcast frame which
// contains 6 bytes of 0xFF followed by 16 repetitions of a given mac address.

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

    file, err := ioutil.ReadFile("hosts.json");
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

    file, err := ioutil.ReadFile("users.json");
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

func modUser(operation string)(error){
    var pass string = url.QueryEscape(string(os.Args[3]));
    var user string = url.QueryEscape(string(os.Args[2]));
    
    
    file, err := ioutil.ReadFile("users.json");
    users := make(map[string]string);
    
    if(err == nil){
        err = json.Unmarshal(file, &users);
        if(err == nil){
            
            if(len(users[user]) == 0 || operation == "edit"){
            
            
                hash, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
            
                if(err==nil){
                    users[user] = string(hash);
            
                    b, err := json.MarshalIndent(users,"","    ");
                
                    if(err==nil){
                
                        f, err := os.Create("users.json")
    
                        if(err != nil){
                            fmt.Println(err);
                        }else{
        
                            _,err := io.WriteString(f, string(b));
                            
                            fmt.Println(err);
                        }   
                
            
                        if(err == nil){
                            fmt.Println(string(hash));
                            return nil;
                        }
                    }
                }
            }
        }
    }
    return errors.New("something went wrong");
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

func httpServer(){
    ln, _ := net.Listen("tcp", ":8081")
    
    for{
        conn,_ := ln.Accept()
        go parseHttp(conn)
    }

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

func parseHttp(conn net.Conn){

    message := bufio.NewReader(conn);
    
    header := parseHeader(*message);
    
    reply := []byte("HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nAccess-Control-Allow-Headers: x-api\nAccess-Control-Allow-Methods: POST\nContent-Type: text/html; charset=utf-8\n\n");
    
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

