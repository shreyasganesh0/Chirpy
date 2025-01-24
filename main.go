package main

import(
    "os"
    "database/sql"
    "log"
    "sync/atomic"
    "net/http"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
    "github.com/shreyasganesh0/chirpy/database"
)

type apiConfig struct {
    file_server_hits atomic.Int32
    queries *database.Queries
    platform string
};

func (cfg *apiConfig) metrics_middleware(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    cfg.file_server_hits.Add(1); //increment view counter
    next.ServeHTTP(w, r);
    });
}

func register_api_endpoints(serv_mux *http.ServeMux, conf *apiConfig){

    serv_mux.HandleFunc("GET /api/healthz",  readiness_handler); 
    serv_mux.HandleFunc("GET /admin/metrics",  conf.metrics_handler); 
    serv_mux.HandleFunc("POST /admin/reset",  conf.reset_metrics_handler);
    serv_mux.HandleFunc("POST /api/validate_chirp", validate_chirp_handler);
    serv_mux.HandleFunc("POST /api/users", conf.users_handler);
}

func main(){

    var conf apiConfig;
    godotenv.Load();

    conf.platform = os.Getenv("PLATFORM");

    db_url := os.Getenv("DB_URL");
    log.Printf("db url %v\n",db_url);
    db, err := sql.Open("postgres", db_url);
    if err != nil{
        log.Printf("%v\n",err);
        return;
    }
    queries := database.New(db);
    conf.queries = queries;
    
    serv_mux := http.NewServeMux();

    server := &http.Server{
        Handler: serv_mux,
        Addr: ":8080",
    };
    
    file_server_handler := http.FileServer(http.Dir(".")); // the . is the local file dir from which files will be served

    cache_control_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Cache-Control", "no-store");
        file_server_handler.ServeHTTP(w, r);
    }); 

    serv_mux.Handle("/app/", http.StripPrefix("/app", conf.metrics_middleware(cache_control_handler)));

    register_api_endpoints(serv_mux, &conf);

    log.Printf("Server Starting!\n");
    log.Fatal(server.ListenAndServe()); //tcp listener and create a new service for each conncection
    log.Printf("Server Stopped!\n");

    return; 
}
