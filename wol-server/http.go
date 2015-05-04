package main

import (
	"net"
	"bufio"
	"strings"
	"encoding/json"
)

func httpServer(){
    ln, _ := net.Listen("tcp", config["listen"])
    
    for{
        conn,_ := ln.Accept()
        go parseHttp(conn)
    }

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
