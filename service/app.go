package main

import (
    "log"
    "net"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

// <- version ->
type Version struct {
    State string `json:"version,omitempty"`
}

var currentVersion = Version{ State:"v1" }

func GetVersion(w http.ResponseWriter, r *http.Request) {
	if (currentVersion.State != ""){ // ???
		json.NewEncoder(w).Encode(currentVersion)
	} else {
    	http.Error(w, "Could not get current version", http.StatusInternalServerError)
	}
}

// <- interfaces ->
func GetInterfaces(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    if vars["api_version"] != currentVersion.State {
        http.Error(w, "Wrong api version. Request /version to get current one.", http.StatusNotFound)
        return
    }
    ifaces, err := net.Interfaces()
    if err != nil {
        http.Error(w, "Could not get existing interfaces : " + err.Error(), http.StatusInternalServerError)
        return
    }
    type ResponseType struct {
        Interfaces []string `json:"interfaces"`
    }
    var ret = ResponseType{};
    for _, i := range ifaces {
        ret.Interfaces = append(ret.Interfaces, i.Name);
    }
    json.NewEncoder(w).Encode(ret)
}

// <- interface information ->
func GetInterfaceInformation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    if vars["api_version"] != currentVersion.State {
        http.Error(w, "Wrong api version. Request /version to get current one.", http.StatusNotFound)
        return
    }
    iface, ierr := net.InterfaceByName(vars["interface_name"])
    if ierr != nil {
        http.Error(w, "Could not get existing interfaces : " + ierr.Error(), http.StatusInternalServerError)
        return
    }
    type ResponseType struct {
        Name string `json:"name"`
        HardwareAddr string `json:"hw_addr"`
        InetAddr []string `json:"inet_addr"`
        MTU int `json:"MTU"`
    }
    var ret = ResponseType{ Name : iface.Name, MTU : iface.MTU }
    ret.HardwareAddr = iface.HardwareAddr.String()
    // SOMEHOW in case of 00:00:00:00:00:00 address standart golang function HardwareAddr.String() provides WRONG output
    if (ret.HardwareAddr == ""){
        ret.HardwareAddr = "00.00.00.00.00.00"
    }
    var addrs, aerr = iface.Addrs()
    if aerr != nil {
        http.Error(w, "Could not get addresses for requested interface : " + ierr.Error(), http.StatusInternalServerError)
        return
    }
    for _, a := range addrs {
        ret.InetAddr = append(ret.InetAddr, a.String())
    }
    json.NewEncoder(w).Encode(ret)
}

// <- main ->
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/{api_version}/interfaces", GetInterfaces).Methods("GET")
    router.HandleFunc("/{api_version}/interface/{interface_name}", GetInterfaceInformation).Methods("GET")
    router.HandleFunc("/version", GetVersion).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", router))
}