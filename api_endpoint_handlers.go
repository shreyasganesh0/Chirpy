package main


import(
    "log"
    "fmt"
	"time"
    "net/http"
    "encoding/json"
	"github.com/google/uuid"
    "github.com/shreyasganesh0/Chirpy/database"
    "github.com/shreyasganesh0/Chirpy/auth"

)

type user_resp_t struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string `json:"email"`
    Token string `json:"token"`
};

type validation_req struct {
    Body string `json:"body"`
};

type chirp_resp_t struct {
    ID        uuid.UUID `json:"id"` 
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body string `json:"body"`
    UserID uuid.UUID `json:"user_id"` 
};

type user_t struct {
    Password string `json:"password"`
    Email string `json:"email"`
    ExpiresInSeconds int `json:"expires_in_seconds"`
};

func (cfg *apiConfig) login_handler(w http.ResponseWriter, req *http.Request) {
    // POST /api/login

    var login_inp user_t;

    decoder := json.NewDecoder(req.Body);
    err1 := decoder.Decode(&login_inp);
    if err1 != nil{
        log.Printf("%v\n", err1);
        return;
    }

    user, err := cfg.queries.GetUserByEmail(req.Context(), login_inp.Email);    
    if err != nil{
        log.Printf("%v\n", err);
        return;
    }

    err_hash := myauth.ComparePassHash(login_inp.Password, user.HashedPassword);
    if err_hash != nil{
        w.WriteHeader(http.StatusUnauthorized);
        w.Write([]byte("incorrect username or password\n"));
        return;
    }

    if login_inp.ExpiresInSeconds < 1 || login_inp.ExpiresInSeconds > 3600 {
        login_inp.ExpiresInSeconds = 3600;
    }

    jwt_token, err_jwt := myauth.MakeJWT(user.ID, cfg.jwt_secret, time.Duration(login_inp.ExpiresInSeconds) * time.Second);
    if err_jwt != nil{
        log.Printf("Error in jwt creation: %v\n", err_jwt);
        return;
    }

    user_resp := user_resp_t{
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Email: user.Email,
        Token: jwt_token,
    };

    resp_byte, err_marsh := json.Marshal(&user_resp);
    if err_marsh != nil{
        log.Printf("%v\n", err);
        return;
    }

    w.WriteHeader(http.StatusOK);
    _, err_write := w.Write(resp_byte);
    if err_write != nil{
        log.Printf("%v\n", err_write);
        return;
    }

    return;
} 

func (cfg *apiConfig) get_chirp_by_id_handler(w http.ResponseWriter, req *http.Request) {
    uuid_val,err := uuid.Parse(req.PathValue("chirpID")); if err != nil{
        log.Printf("%v\n", err);
        return;
    }
    chirp, err := cfg.queries.GetChirpByID(req.Context(), uuid_val); 
    if err != nil {
        log.Printf("%v\n", err);
        return;
    }

    chirp_resp := chirp_resp_t{
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body: chirp.Body,
        UserID: chirp.UserID,
    };
    chirp_byte, err1 := json.Marshal(chirp_resp);
    if err1 != nil {
        log.Printf("%v\n", err);
        return;
    }
    w.Header().Set("Content-Type", "application/json");
    w.WriteHeader(http.StatusOK);
    _, err2 := w.Write(chirp_byte);
    if err2 != nil{
        log.Printf("%v\n", err2);
        return;
    }
    return;
}
func (cfg *apiConfig) get_chirps_handler(w http.ResponseWriter, req *http.Request) {

    chirps, err := cfg.queries.GetAllChirps(req.Context());
    if err != nil {
        log.Printf("%v\n", err);
        return;
    }

    var chirps_resp []chirp_resp_t;
    for _, chirp := range chirps{
        chirp_resp := chirp_resp_t{
            ID: chirp.ID,
            CreatedAt: chirp.CreatedAt,
            UpdatedAt: chirp.UpdatedAt,
            Body: chirp.Body,
            UserID: chirp.UserID,
        };
        chirps_resp = append(chirps_resp, chirp_resp);
    }
    chirps_byte, err1 := json.Marshal(chirps_resp);
    if err1 != nil {
        log.Printf("%v\n", err);
        return;
    }
    w.Header().Set("Content-Type", "application/json");
    w.WriteHeader(http.StatusOK);
    _, err2 := w.Write(chirps_byte);
    if err2 != nil{
        log.Printf("%v\n", err2);
        return;
    }
    return;
}


func (cfg *apiConfig) chirps_handler(w http.ResponseWriter, req *http.Request) {
    //POST /api/chirps

    var val_req validation_req;
    err := validate_chirp(w, req, &val_req);
    if err != nil{
        log.Printf("%v\n", err);
        return;
    }

    bearer_token, err_bearer := myauth.GetBearerToken(req.Header);
    if err_bearer != nil{
        log.Printf("Error while retrieving bearer token: %v\n", err_bearer);
        return;
    }
    
    log.Printf("Bearer token: %v\n", bearer_token);
    user_id, err_val := myauth.ValidateJWT(bearer_token, cfg.jwt_secret);
    if err_val != nil{
        log.Printf("Error while validating jwt %v\n", err_val);
        return;
    }

    var err_id uuid.UUID
    if user_id == err_id{
        w.WriteHeader(http.StatusUnauthorized);
        _, err := w.Write([]byte("Unauthorized user for post"));
        if err != nil {
           log.Print("Error sending message to unauth user while creating chrip\n"); 
        }
        return;
    }
    
    query_args := database.CreateChirpParams{
        Body: val_req.Body,
        UserID: user_id,
    };

    chirp, err := cfg.queries.CreateChirp(req.Context(), query_args);
    chirp_resp := chirp_resp_t{
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body: chirp.Body,
        UserID:  chirp.UserID,
    };

    resp, err_json := json.Marshal(&chirp_resp);
    if err_json != nil {
        fmt.Printf("Error writing to body %v\n", err_json);
        return;
    }
    w.Header().Set("Content-Type", "application/json");
    w.WriteHeader(http.StatusCreated);
    _, err1 := w.Write(resp); 
    if err1 != nil {
        fmt.Printf("Error writing to body %v\n", err1);
    }
    return;
}


func (cfg *apiConfig) users_handler(w http.ResponseWriter, req *http.Request) {
    var user_req user_t;

    decoder := json.NewDecoder(req.Body);
    err := decoder.Decode(&user_req);
    if err != nil {
        log.Printf("%v\n", err);
        return;
    }

    hash_pass, err_hash := myauth.EncryptPassword(user_req.Password);
    if err_hash != nil{
        log.Printf("%v\n", err);
        return;
    }

    query_args := database.CreateUserParams{
        Email: user_req.Email,
        HashedPassword: hash_pass,
    };

    user, err1 := cfg.queries.CreateUser(req.Context(), query_args); 
    if err1 != nil{
        log.Printf("%v\n", err);
        return;
    }

    resp := user_resp_t{
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

func validate_chirp(w http.ResponseWriter, req *http.Request, val_req *validation_req) error{

    type len_err_response struct {
        Error string `json:"error"`
    };
    
    var err_resp len_err_response

   // decoding req part 
    decoder := json.NewDecoder(req.Body);
    err := decoder.Decode(&val_req);
    if err != nil {
        log.Printf("error parsing json\n%v",err);
        err_resp.Error = "Something went wrong";
    }

    // validation part
    body_len := len(val_req.Body)
    if body_len > 140 {
        err_resp.Error = "Chirp is too long";
    }

    if err_resp.Error != ""{
        status_code := http.StatusBadRequest;
        resp, err := json.Marshal(err_resp);
        if err != nil{
            return fmt.Errorf("Failed response marshalling %v\n", err);
        }
        w.Header().Set("Content-Type", "application/json");
        w.WriteHeader(status_code);
        w.Write(resp);
        return fmt.Errorf("Chirp invalid, request sent");
    }
    return nil

}

