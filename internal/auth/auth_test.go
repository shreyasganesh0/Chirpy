package myauth

import(
    "testing"
	"github.com/google/uuid"
    "time"
)

func TestMakeJWT(t *testing.T){
    uid := uuid.New();
    secret := "secret";
    duration := 2 * time.Hour;

    jwt, err := MakeJWT(uid, secret, duration);
    if err != nil{
        t.Errorf("Failed to create JWT\n %v", err);
    }

    test_uid, err := ValidateJWT(jwt, secret);
    if (test_uid != uid){
        t.Errorf("Mismatch when expected match for uid before and after %v != %v", test_uid, uid);
    }
    return ;
}
