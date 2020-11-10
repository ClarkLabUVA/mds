package identifier

import (
	"strings"
	"net/http"
    "github.com/dgrijalva/jwt-go"
    "os"
    "context"
    "io/ioutil"
    "encoding/json"

    "fmt"
    "github.com/rs/zerolog"
)

var (
    authLogger = zerolog.New(os.Stderr).With().Timestamp().Str("backend", "auth").Logger()
)

var jwtSecret  []byte
var authURI string


func init() {

    jwtENV, ok := os.LookupEnv("JWT_SECRET")

    if !ok {
        jwtENV = "test secret"
    }

    jwtSecret = []byte(jwtENV) 

    authENV, ok := os.LookupEnv("AUTH_URI")

    if !ok {
        authENV = "http://auth"
    }

    authURI = authENV
}

/*
//AuthConfig is the struct
type AuthConfig struct {
    Enabled         bool
    AuthServiceURI  string
    JWTSecret       []byte
}
*/


//User is the struct used by handlers for determining privleges
type User struct {
    ID      string   `json:"@id" bson:"@id"`
	Type    string   `json:"@type" bson:"@type"`
	Name    string   `json:"name" bson:"name"`
	Email   string   `json:"email" bson:"email"`
    Role    string   `json:"role" bson:"-"`
    Groups  []string   `json:"groups" bson:"-"`
}


//AllowedAccess determines if a user is allowed access to a 
func (u *User) AllowedAccess(r Resource) (allowed bool) {

    if u.Role == "admin" {
        return true
    }

    if u.Role == "user" {

        // check if user is the owner
        if r.Owner == u.ID {
            return true
        }

        // check if user id is allowed in 
        allowedUsers := strings.Join(r.Groups, ";")

        if strings.Contains(allowedUsers, u.ID) {
            return true
        }

        // check if any group on the resource is one our user belongs to
        userGroups := strings.Join(u.Groups, ";")

        for _, resourceGroup := range r.Groups {
            if strings.Contains(userGroups, resourceGroup) {
                return true
            }
        }



    }


    return false
}



//UserTokenClaims handles the custom claims for fairscape json web tokens
type UserTokenClaims struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Role string `json:"role"`
	Groups []string `json:"groups"`
	jwt.StandardClaims
}

// AuthMiddleware is a general authentication middleware for fairscape services
// this middleware checks that token is present and valid for a user
// before unpacking the token with its custom fields and adding it to the request context.
// This is implemented as a negroni middleware.
func AuthMiddleware (next http.Handler) http.Handler {
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        // read bearer token from request
        var authToken string

        authToken = r.Header.Get("Authorization")

        // if bearer token doesn't exist
        if authToken == "" {

            // check cookies of reqest
            authCookie, err := r.Cookie("fairscapeAuth") 

            if err != nil {
                w.Write([]byte(`{"error": "request missing authorization token"}`))
                w.WriteHeader(403)
                return
            }

            authToken = authCookie.Value
            
        }

        // decrypt and inspect jwt
        userToken, err := jwt.ParseWithClaims(
            authToken, 
            &UserTokenClaims{}, 
            func(tk *jwt.Token) (interface{}, error) {
            // we are just using the signing secret
            return jwtSecret, nil
        })

        // error handling for token error
        if err != nil {

            authLogger.Error().
                Err(err).
                Msg("Error Parsing Token")

            
            w.Write([]byte(`{"message": "invalid token", "error": "`+ err.Error() + `"}`))        
            w.WriteHeader(401)
            return
        }

        claims := userToken.Claims.(*UserTokenClaims)

        authLogger.Info().
            Str("claims", fmt.Sprintf("%v+", claims)).
            Msg("Claims from Token")

        u := User{
            ID: claims.Subject,
            Type: "Person",
            Name: claims.Name,
            Email: claims.Email,
            Role: claims.Role,
            Groups: claims.Groups,
        }

        authLogger.Info().
            Dict("user", 
                zerolog.Dict().
                    Str("ID", u.ID).
                    Str("Name", u.Name).
                    Str("Email", u.Email).
                    Str("Role", u.Role).
                    Interface("Groups", u.Groups),
                ).
            Msg("User Data from Claims")


        // take token and put into context
        ctx := context.WithValue(r.Context(), "user", u)

        // Call the next handler 
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}


//Resource is the struct for the 
type Resource struct {
	ID    	string `json:"@id" bson:"@id"`
	Type  	string `json:"@type" bson:"@type"` 
	Owner 	string `json:"owner" bson:"owner"`
	Users 	[]string `json:"users" bson:"users"`
	Groups 	[]string `json:"groups" bson: "groups"`
}

//AuthGetACL queries the ACL for the specified resource at the auth service
func AuthGetACL(id string) (r Resource, err error) {

    resourceURL := authURI + "/" + id

    response, err := http.Get(resourceURL)

    if err != nil {
        err = fmt.Errorf("AuthGetACL: Failed to Get Resource (%w)", err)
        return
    }

    responseBody, err := ioutil.ReadAll(response.Body)

    if err != nil {
        err = fmt.Errorf("AuthGetACL: Error Reading Auth Service Response (%w)", err)
        return
    }

    // decode response body into response
    err = json.Unmarshal(responseBody, &r)

    if err != nil {
        err = fmt.Errorf("AuthGetACL: Error Unmarshaling Auth Service Response (%w)", err)
    }

    return
    
}