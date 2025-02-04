package auth
import(
    "github.com/golang-jwt/jwt/v5"
    "time"
    "github.com/google/uuid"
    "fmt"
    "net/http"
    "strings"
)

func MakeJWT(userID uuid.UUID, tokenSecret string)(string, error){  
    expireTime := time.Duration(3600)*time.Millisecond
    claims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expireTime)),
        IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
        NotBefore: jwt.NewNumericDate(time.Now().UTC()),
        Issuer:    "chirpy",
        Subject:   userID.String(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    ss, err := token.SignedString([]byte(tokenSecret))
    if err != nil{
        return "", err
    }
    return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
    claims := &jwt.RegisteredClaims{}
    _, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("Wrong signing method: %v", token.Header["alg"])
            }
            return []byte(tokenSecret), nil
    })

    if err != nil {
        return uuid.UUID{}, fmt.Errorf("token parsing failed: %w", err)
    }

    userID, err := uuid.Parse(claims.Subject) // Access Subject directly
    if err != nil {
            return uuid.UUID{}, fmt.Errorf("Unable to parse subject for uuid: %w", err)
    }
    return userID, nil
}

func GetBearerToken(headers http.Header) (string, error){
    authHeader := headers.Get("Authorization")
    if authHeader == ""{
        return "", fmt.Errorf("No authorized header")
    }
    splitHeader := strings.Split(authHeader," ")
    if len(splitHeader)!=2 || splitHeader[0] != "Bearer"{
        return "", fmt.Errorf("Malformed header")
    } 
    return splitHeader[1], nil
}