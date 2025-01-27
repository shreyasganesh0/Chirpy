package myauth

import(
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)


func EncryptPassword(password string) (string, error){

    hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), 10) // the 10 is just a preset value for the cost of the hash
    if err != nil{
        return "", err;
    }

    return string(hashed_pass), nil;
}

func ComparePassHash(password string, hash string) error{
    
    err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password));
    if err != nil{
        return err;
    }

    return nil;

}

func MakeJWT(user_id uuid.UUID, token_secret string, expires_in time.Duration) (string, error){
    
    current_time := time.Now();
    claims := jwt.RegisteredClaims{
        Issuer: "chirpy",
        Subject: user_id.String(),
        IssuedAt: jwt.NewNumericDate(current_time),
        ExpiresAt: jwt.NewNumericDate(current_time.Add(expires_in)),
    };

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims);
    ss, err := token.SignedString([]byte(token_secret));
    if err != nil{
        return "", err;
    }
    
    return ss, nil;
}

func ValidateJWT(token_string string, token_secret string) (uuid.UUID, error){

    var err_uid uuid.UUID;
    var claim jwt.RegisteredClaims;
    token, err := jwt.ParseWithClaims(token_string, &claim, func(token *jwt.Token) (interface{}, error){
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Wrong signing method %v\n", token.Header["alg"]);
        }

        return []byte(token_secret), nil;
    });

    if err != nil {
        return err_uid, err;
    }

    user_id, err := token.Claims.GetSubject();
    if err != nil {
        return err_uid, err;
    }
        
    uuid_user, err1 := uuid.Parse(user_id);
    if err1 != nil {
        return err_uid, err1;
    }
    return uuid_user, nil;
}



