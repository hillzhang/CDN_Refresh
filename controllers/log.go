package controllers
import (
    "github.com/getsentry/raven-go"
    l"log"
)

var sentryDSN string = "http://e492319790d84e278d1bed556c518696:085ba1b15a7945b7bead386e067f7abd@10.100.6.75:9000/9"

func CaptureLog(log string){
    client, err := raven.NewClient(sentryDSN, nil)
    if err != nil {
        l.Println(err)
    }
    packet := raven.NewPacket(log, nil,nil)
    packet.Level = raven.INFO
    _, ch := client.Capture(packet, nil)
    if err = <-ch; err != nil {
        l.Println(err)
    }
}

func CaptureError(error string){
    client, err := raven.NewClient(sentryDSN, nil)
    if err != nil {
        l.Println(err)
    }
    packet := raven.NewPacket(error, nil,nil)
    packet.Level = raven.ERROR
    _, ch := client.Capture(packet, nil)
    if err = <-ch; err != nil {
        l.Println(err)
    }
}

