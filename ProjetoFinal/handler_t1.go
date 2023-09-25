package main

import (
        "fmt"
/*        "bufio"
        "os"
        "strings"*/
        "DeviceProxy"
//        "DeviceLookupProxy"
        "time"
)


var Proxy DeviceProxy.DeviceProxy
//var Lookup DeviceLookupProxy.DeviceLookupProxy


func main() {
  Proxy.New("127.0.0.1:8081", true)
//  Lookup.New("127.0.0.1:8082", false)

//  Lookup.Bind("Sensor 1", "Samsung", "s2", "Temperatura")

//  device := Lookup.LookupAll()
//  fmt.Println(device)

/*  reader := bufio.NewReader(os.Stdin)            //lê msg escrita pelo cliente

  for { 
    text, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    if strings.TrimRight(text, "\r\n") == "Lista"{
      Proxy.QueueDeclare("Lista 1")
    } else if strings.TrimRight(text, "\r\n") == "Recebe"{
      msg := Proxy.Consume("Lista 1", true)
      //Proxy.SendAck("Lista 1", false)
      if msg != "" {
        fmt.Println(msg)
      }
    } else if strings.TrimRight(text, "\r\n") == "Fecha" {
      Proxy.Fecha()
    } else if strings.TrimRight(text, "\r\n") == "Abre" {
      Proxy.Abre()
    } else {
      Proxy.Publish("Lista 1", "text", false)
    }
   
  }*/

  Proxy.QueueDeclare("Lista 1")

  fmt.Println("Começa Publish")
  for i := 0; i < 250; i++ {
  //  start := time.Now()
    Proxy.Publish("Lista 1", "text", false)
    duration := time.Second
    time.Sleep(duration)
  //  fmt.Println(start.Format("15;04;05.000000000"))
  }

  fmt.Println("Acaba Publish")

  for i := 0; i < 250; i++ {
    start := time.Now()
    Proxy.Consume("Lista 1", true)
    end := time.Now()
    duration := time.Second
    time.Sleep(duration)
    fmt.Println(start.Format("15;04;05.000000000"),";",end.Format("15;04;05.000000000"))
  }

}