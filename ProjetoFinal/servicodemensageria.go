package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "os"
    "Marshaller"
    "Encryptor"
    "IoTMessage"
)

var deviceList []*Device          //lista de clientes no chat
var queueList []*Queue

var marshar Marshaller.Marshaller

var encrypter Encryptor.Encryptor

//Estrutura para armazenar as informações de cada cliente
type Device struct {
  outgoing       []string         //buffer de saída de mensagens, a serem enviadas para o cliente
  outgoingCounter int 
  contactList    []*Device        //lista de clientes com os quais ele pode se comunicar
  name           string           //nome do cliente
}

type Queue struct {
	name			string
	messages		[]string
	messagesCounter	int
}

/*type Msg struct {
	CommandId	int
	ListName	string
	Message 	string
}*/


//Função de leitura de msg vinda do cliente
func (device *Device) Read(connection net.Conn) {
	texto := make([]byte, 1024)
	for {
		n, err := connection.Read(texto)
		if err != nil {
		    fmt.Println(err. Error())
		    //os.Exit(2)
		    break
		}
		msg := texto[:n]
		decryptedMsg := encrypter.Decrypt(msg, "senha")
		//Mensagem := Msg{}
		Mensagem := IoTMessage.IoTMessage{}
		err = marshar.Unmarshal(decryptedMsg, &Mensagem)
		if err != nil {
			fmt.Println(err. Error())
			break
	    }

		command := Mensagem.Header.CommandId

		switch command {
		case 1:
			QueueDeclare(Mensagem)

		case 2:
			Publish(Mensagem, device.name, false)

		case 3:
			Publish(Mensagem, device.name, true)

		case 4:
			Consume(Mensagem, connection)

		case 5:
			Remove(Mensagem)

		}
	}

	fmt.Println(device.name, "saiu do chat")
}

func QueueDeclare(msg IoTMessage.IoTMessage){
	listExist := CheckQueue(msg)
	if listExist == false {
		queue := &Queue{                        //cria novo cliente e define algumas informações
			name: msg.Header.ListName,
			messagesCounter: 0,
		}

		queueList = append(queueList, queue)   //adiciona novo cliente na lista geral de clientes
	}

	for i := range queueList {
			fmt.Println(queueList[i].name)
	}
}

func CheckQueue(msg IoTMessage.IoTMessage) bool{
	flag := false
	for i := range queueList{
		if msg.Header.ListName == queueList[i].name{
			flag = true
			break
		}
	}

	return flag
}

func Consume(msg IoTMessage.IoTMessage, connection net.Conn) {
	for i := range queueList {
		if msg.Header.ListName == queueList[i].name{
			if queueList[i].messagesCounter >  0 {              //se o contador de mensagens disponíveis estiver maior que 0
			    line := queueList[i].messages[0]                   //é salva a linha de mensagem na primeira posição do buffer

			    msgRequest := IoTMessage.IoTMessage{}
			    msgRequest.Header.CommandId = msg.Header.CommandId
			    msgRequest.Header.ListName = msg.Header.ListName
			    msgRequest.Body.Message = line

    			msgMarshada, _ := marshar.Marshal(msgRequest)
    			cryptedMsg := encrypter.Encrypt(msgMarshada, "senha")
		 		connection.Write(cryptedMsg)

		 	} else {
		 		cryptedMsg := encrypter.Encrypt([]byte("\nil"), "senha")
		 		connection.Write(cryptedMsg)
		 	}
		 	
			break
		}
	}
}

func Remove(msg IoTMessage.IoTMessage) {
	for i := range queueList {
		if msg.Header.ListName == queueList[i].name{
			if queueList[i].messagesCounter >  0 {              //se o contador de mensagens disponíveis estiver maior que 0
			    if msg.Body.Message == "true" {
			    	queueList[i].messages = append(queueList[i].messages[:0], queueList[i].messages[1:]...)   //a mensagem na primeira posição do buffer é retirada anexando a partir da segunda, criando um fila tipo FIFO
					queueList[i].messagesCounter--	                 //o contador de mensagens do cliente é decrementado
			    } else {
			    	line := queueList[i].messages[0]                   //é salva a linha de mensagem na primeira posição do buffer
			    	queueList[i].messages = append(queueList[i].messages[:0], queueList[i].messages[1:]...)   //a mensagem na primeira posição do buffer é retirada
					queueList[i].messages = append(queueList[i].messages, line)								  //a mensagem salva é colocana no final da fila
			    }
			    
		 	}		 	
			break
		}
	}
}


func Publish(msg IoTMessage.IoTMessage, device string, priority bool) {
	line := device + " diz: " + string(msg.Body.Message)
	fmt.Println(line)
	for i := range queueList {
		if msg.Header.ListName == queueList[i].name{
			if priority == false {
				queueList[i].messages = append(queueList[i].messages, line)
			} else {
				queueList[i].messages = append([]string{line}, queueList[i].messages...)
			}
			
			queueList[i].messagesCounter++
			break
		}
	}
}


//Função de adição de novo cliente no chat
func NewDevice(connection net.Conn) {

  device := &Device{                        //cria novo cliente e define algumas informações
      outgoingCounter: 0,
  }

  n, err := bufio.NewReader(connection).ReadString('\n')   //primeia msg é o nome do cliente
  if err != nil {
    fmt.Println(err. Error())
    os.Exit(2)
  } else {
    device.name = strings.TrimRight(n, "\r\n")             //armazena nome do cliente
    fmt.Println(device.name, "entrou no chat")
  }

  for _, oldDevice := range deviceList {                  //roda a lista de cliente
    device.contactList = append(device.contactList, oldDevice)    //adiciona clientes antigos como contato do novo
    oldDevice.contactList = append(oldDevice.contactList, device) //adiciona cliente novo como contato dos antigos
  }

  deviceList = append(deviceList, device)   //adiciona novo cliente na lista geral de clientes

  go device.Read(connection)                          //go rotine de leitura de msgs do novo cliente
}


func main() {
  listener, _ := net.Listen("tcp", ":8081")
  fmt.Println("Servidor pronto...")
  for {
    conn, err := listener.Accept()      //fica aguardando novos clientes entrarem no chat
    if err != nil {
        fmt.Println(err. Error())
        os.Exit(0)
    }
    NewDevice(conn)                     //chama função para adicionar o cliente que entrou
  }
}