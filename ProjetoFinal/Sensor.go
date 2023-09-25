package main

import (
        "fmt"
        "bufio"
        "os"
        "DeviceProxy"
        "DeviceLookupProxy"
)


var Proxy DeviceProxy.DeviceProxy
var Lookup DeviceLookupProxy.DeviceLookupProxy


func main() {
  Proxy.New("127.0.0.1:8081", false)
  Lookup.New("127.0.0.1:8082", false)

  Lookup.Bind("Sensor 1", "Lego", "s2", "Temperatura")

  Proxy.QueueDeclare("Write")
  
  reader := bufio.NewReader(os.Stdin)            //lê msg escrita pelo cliente

  for { 
    text, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    Proxy.Publish("Write", text, false)
   
  }

}