package main

import (
    "os"
    "fmt"
    "log"
    "flag"
    "time"
    "strings"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

func PrintHelp(){
    fmt.Println(`
NAME:
    cli_net - client displays network interfaces and info.
USAGE:
    cli_net [global options] command [command options] [arguments...]
VERSION:
    0.0.0
COMMANDS:
    help, h Shows a list of commands or help for one command ...
GLOBAL OPTIONS
    --version Shows version information
    `)
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getResponse(url string, target interface{}) error {
    r, err := myClient.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Body.Close()
    if (r.StatusCode != http.StatusOK){
        fmt.Println("HTTP Response Status:", r.StatusCode, http.StatusText(r.StatusCode))
        bodyBytes, err2 := ioutil.ReadAll(r.Body)
        if err2 != nil {
            log.Fatal(err2)
        }
        fmt.Println(string(bodyBytes))
        os.Exit(1)
    }
    return json.NewDecoder(r.Body).Decode(target)
}

func list(server string, port string){
    type ResponseType struct {
        Interfaces []string `json:"interfaces"`
    }
    response := new(ResponseType) 
    var url = "http://" + server + ":" + port + "/v1/interfaces"
    getResponse(url, response)
    for _,i := range response.Interfaces{
        fmt.Print("\t" + i)
    }
    fmt.Println("")
}

func show(server string, port string, interface_name string){
    if len(os.Args) < 3 {
        fmt.Println(`
show command requires additional argument - interface name.
example :
    cli_net show eth0 --server 127.0.0.1 --port 8080
        `)
        os.Exit(1)
    }
    type ResponseType struct {
        Name string `json:"name"`
        HardwareAddr string `json:"hw_addr"`
        InetAddr []string `json:"inet_addr"`
        MTU int `json:"MTU"`
    }
    response := new(ResponseType) 
    var url = "http://" + server + ":" + port + "/v1/interface/" + interface_name
    getResponse(url, response)
    fmt.Printf("%s:\t hw addr: %v\n", response.Name, response.HardwareAddr)
    fmt.Printf("\t inet addr: %s\n", strings.Join(response.InetAddr, ", "))
    fmt.Printf("\t MTU: %d\n", response.MTU)
}

func required(arg string, arg_name string) (string) {
    if (arg == ""){
        fmt.Println(`
additional parameter required:
    ` + arg_name + `
        `)
        os.Exit(1)
    }
    return arg
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println(`
One of the following commands required: 
    help, h, list, show
        `)
        os.Exit(1)
    }
    cmd := flag.NewFlagSet("cmd", 0)
    portPtr := cmd.String("port", "", "")
    serverPtr := cmd.String("server", "", "")
    
    switch os.Args[1] {
    case "list":
        cmd.Parse(os.Args[2:])
        list(required(*serverPtr, "server"), required(*portPtr,"port"))
    case "show":
        cmd.Parse(os.Args[3:])
        show(required(*serverPtr, "server"), required(*portPtr,"port"), os.Args[2])
    case "help":
        PrintHelp()
    case "h":
        PrintHelp()
    default:
        fmt.Printf("Unknown command - %s", os.Args[1])
        os.Exit(1)
    }
}