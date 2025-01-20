package main

import(
    "net/http"
    "log"
)

type apiHandler struct{};

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request){};

func main(){
    
    serv_mux := http.NewServeMux();

    server := http.Server{
        Handler: serv_mux,
        Addr: ":8080",
    };
    
    file_server_handler := http.FileServer(http.Dir("."));
    serv_mux.Handle("/", file_server_handler); // register the handle with the pattern /

    log.Printf("Server Starting!\n");
    log.Fatal(server.ListenAndServe()); //tcp listener and create a new service for each conncection
    log.Printf("Server Stopped!\n");

    return; 
}
