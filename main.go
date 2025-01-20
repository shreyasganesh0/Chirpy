package main

import(
    "net/http"
    "fmt"
)

func main(){
    
    serv_mux := http.NewServeMux();

    server := http.Server{
        Handler: serv_mux,
        Addr: ":8080",
    };
    
    err := server.ListenAndServe();
    fmt.Printf("Server Started!\n");

    if err != nil{
        fmt.Printf("Error with the server %v", err);
        return;
    }

   return; 
}
