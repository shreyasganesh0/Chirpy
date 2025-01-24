package main


import(
    "log"
    "fmt"
	"time"
    "strings"
    "net/http"
    "encoding/json"
	"github.com/google/uuid"

)

func (cfg *apiConfig) users_handler(w http.ResponseWriter, req *http.Request) {
    type user_t struct {
        Email string `json:"email"`
    };
    type resp_t struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Email     string `json:"email"`
    };
    var user_req user_t;

    decoder := json.NewDecoder(req.Body);
    err := decoder.Decode(&user_req);
    if err != nil {
        log.Printf("%v\n", err);
        return;
    }

    user, err1 := cfg.queries.CreateUser(req.Context(), user_req.Email); 
    if err1 != nil{
        log.Printf("%v\n", err);
        return;
    }

    resp := resp_t{
        ID: user.ID,       
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt, 
        Email: user.Email,
    };
    resp_byte, err2 := json.Marshal(&resp);
    if err2 != nil{
        log.Printf("%v\n", err);
        return; 
    }
    w.Header().Set("Content-Type", "application/json");
    w.WriteHeader(http.StatusCreated);
    _, err3 := w.Write(resp_byte);
    if err3 != nil{
        log.Printf("%v\n", err);
        return;
    }
    return;
}


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

    w.Header().Set("Content-Type", "text/html");
    _, err := w.Write([]byte(html));
    if err != nil {
        log.Printf("Error writing to body %v\n", err);
    }
    return;
}

func (cfg *apiConfig) reset_metrics_handler(w http.ResponseWriter, req *http.Request){

    if cfg.platform != "dev"{
        log.Printf("This endpoint cannot be accessed by PLATFORM %v\n", cfg.platform);
        w.WriteHeader(http.StatusForbidden);
        body := "forbidden";
        w.Write([]byte(body));
        return;
    }
    err := cfg.queries.ResetTables(req.Context());
    if err != nil{
        log.Printf("Failed to reset tables\n");
        return;
    }
    log.Printf("Tables reset, all rows deleted\n");

    cfg.file_server_hits.Store(0);
    log.Printf("Server hits count reset\n");
    return;
}

func validate_chirp_handler(w http.ResponseWriter, req *http.Request){

    type validation_req struct {
        Body string `json:"body"`
    };

    type len_err_response struct {
        Error string `json:"error"`
    };
    
    type valid_response struct {
        CleanedBody string `json:"cleaned_body"` 
    };

    var val_req validation_req
    var err_resp len_err_response
    var valid_resp valid_response
    var status_code int
    profane_word_list := []string{"kerfuffle", "sharbert", "fornax"};

   // decoding req part 
    decoder := json.NewDecoder(req.Body);
    err := decoder.Decode(&val_req);
    if err != nil {
        log.Printf("error parsing json\n%verr\n");
        err_resp.Error = "Something went wrong";
        status_code = http.StatusBadRequest;

        resp, err := json.Marshal(err_resp);
        if err != nil{
            log.Printf("Failed response marshalling\n");
            return;
        }
        w.Header().Set("Content-Type", "application/json");
        w.WriteHeader(status_code);
        _, err1 := w.Write(resp); 
        if err1 != nil {
            log.Printf("Error writing to body %v\n", err1);
        }
        return;
    }

    // validation part
    body_len := len(val_req.Body)
    if body_len > 140 {
        log.Printf("Characters length too long: %v\n", body_len);
        err_resp.Error = "Chirp is too long";
        status_code = http.StatusBadRequest;
        resp, err := json.Marshal(err_resp);
        if err != nil{
            log.Printf("Failed response marshalling\n");
            return;
        }
        w.Header().Set("Content-Type", "application/json");
        w.WriteHeader(status_code);
        _, err1 := w.Write(resp); 
        if err1 != nil {
            log.Printf("Error writing to body %v\n", err1);
        }
        return;
    }

    // valid response part
    status_code = http.StatusOK;
    
    response_string_slice := make([]string, 1);
    curr_word := "";
    for _, char := range val_req.Body{ // convert it to a slice of words from a single string
        if (char == ' '){
            response_string_slice = append(response_string_slice, curr_word);
            curr_word = "";
            continue;
        }
        curr_word =curr_word + string(char);
    }
    response_string_slice = append(response_string_slice, curr_word);
    
    var output_str string;
    p_flag := false;
    for _, word := range response_string_slice{ // remove profane words
        for _, profane_word := range profane_word_list{
            if profane_word == strings.ToLower(word){
                output_str = output_str + " ****";
                p_flag = true;
                break;
            }
        }
        if p_flag == true{
            p_flag = false;
            continue;
        }
        output_str = output_str + " " + word;
    }

    output_str = output_str[2:];
    valid_resp.CleanedBody = output_str; 
    resp, err := json.Marshal(valid_resp);
    if err != nil{
        log.Printf("Failed response marshalling\n");
        return;
    }
    w.Header().Set("Content-Type", "application/json");
    w.WriteHeader(status_code);
    _, err1 := w.Write(resp); 
    if err1 != nil {
        log.Printf("Error writing to body %v\n", err1);
    }
    return;
}
