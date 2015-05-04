package main

import (
	"bytes"
	"net"
	"encoding/binary"
)

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
