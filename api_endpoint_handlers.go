package main


import(
    "log"
    "net/http"
    "fmt"
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

    var header_elements []tag_level;
    body_elements := make([]tag_level, 2);
    body_elements[0] = tag_level{
        tag: H1,
        level: 1,
        text: "Welcome, Chirpy Admin",
    };

    body_elements[1] = tag_level{
        tag: P,
        level: 1,
        text: fmt.Sprintf("Chirpy has been visited %v times!", int64(cfg.file_server_hits.Load())),
    };

    html, err_html :=  generate_html(header_elements, body_elements);
    if err_html != nil{
        log.Printf("%v", err_html);
    }

    _, err := w.Write([]byte(html));
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
