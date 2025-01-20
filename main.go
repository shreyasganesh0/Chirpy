package main

import(
    "net/http"
    "log"
)


var file_server_handler http.Handler;

func readiness_handler(w http.ResponseWriter,req *http.Request){
    content_type := make([]string,1);
    content_type[0] = "text/plain; charset=utf-8";

    header_map := w.Header(); 
    header_map["Content-Type"] = content_type;

    w.WriteHeader(http.StatusOK); // techinically not needed if we call Write it automatically sets content-type, statusOk and, if small enough message, content-length
    body := "OK";
    _, err := w.Write([]byte(body));
    if err != nil{
        log.Printf("Error Writing %v\n", err);
    }
    return;
}

func main(){
    
    serv_mux := http.NewServeMux();

    server := http.Server{
        Handler: serv_mux,
        Addr: ":8080",
    };
    
    file_server_handler = http.FileServer(http.Dir(".")); // the . is the local file dir from which files will be served
    serv_mux.Handle("/app/", http.StripPrefix("/app",http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Cache-Control", "no-store");
        file_server_handler.ServeHTTP(w, r);
    }))); // This is wrapper on the file_server_handler passed to prevent etag caching (might remove later)
         // Remove /app prefix from the url req path /app/home.png -> /home.png which will then be served by ./home.png

    serv_mux.HandleFunc("/healthz",  readiness_handler); // pass a handler func for the readiness check endpoint

    log.Printf("Server Starting!\n");
    log.Fatal(server.ListenAndServe()); //tcp listener and create a new service for each conncection
    log.Printf("Server Stopped!\n");

    return; 
}
