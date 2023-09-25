package main

import (
    "Broker"
)


var QueueManager Broker.Broker


func main() {
    QueueManager.New()
}