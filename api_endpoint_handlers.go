package main


import(
    "log"
    "net/http"
    "strconv"
)

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

func (cfg *apiConfig) metrics_handler(w http.ResponseWriter, req *http.Request){

    body := "Hits: " + strconv.FormatInt(int64(cfg.file_server_hits.Load()), 10); 
    _, err := w.Write([]byte(body));
    if err != nil {
        log.Printf("Error writing to body %v\n", err);
    }
    return;
}

func (cfg *apiConfig) reset_metrics_handler(w http.ResponseWriter, req *http.Request){

    cfg.file_server_hits.Store(0);
    log.Printf("Server hits count reset\n");
    return;
}
