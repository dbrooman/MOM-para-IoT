package main

import (
        "fmt"
        "bufio"
        "os"
        "strings"
        "DeviceProxy"
        "DeviceLookupProxy"
)


var Proxy DeviceProxy.DeviceProxy
var Lookup DeviceLookupProxy.DeviceLookupProxy


func main() {
  Proxy.New("127.0.0.1:8081", false)
  Lookup.New("127.0.0.1:8082", false)

  Lookup.Bind("Atuador 1", "Google", "amp1", "Lampada")

  Proxy.QueueDeclare("Read")

  reader := bufio.NewReader(os.Stdin)            //lê msg escrita pelo cliente

  for { 
    text, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

   if strings.TrimRight(text, "\r\n") == "Recebe"{
      msg := Proxy.Consume("Read", true)
      //Proxy.SendAck("Lista 1", false)
      if msg != "" {
        fmt.Println(msg)
      }
    }
   
  }

}