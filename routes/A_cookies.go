package  routes


import (
	"net/http"
	"strings"
	"time"
	"fmt"
)

func SanitizeCookieValue(value string) string {
    return strings.NewReplacer(
        "\r", "",
        "\n", "",
        ";", "",
        " ", "_",

    ).Replace(value)
}
func SanitizeStringTwo(value string) string{
	return strings.NewReplacer(
	"_", " ",	
	).Replace(value)
}


func CreateCookie(user_name, user_id string, w http.ResponseWriter, r *http.Request){

	
	if cookie, err := r.Cookie("user_info"); err == nil {
		// Cookie exists - verify its value matches current user
		parts := strings.Split(cookie.Value, ":")
		if len(parts) == 2 && parts[0] == user_id && parts[1] == user_name {
			
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Login successful"))
			return
		}
	
	}
	value := SanitizeCookieValue(fmt.Sprintf("%s:%s", user_id, user_name))
	
	http.SetCookie(w, &http.Cookie{
		Name:     "user_info",
		Value:    value,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Println("Created Cookie")
}


func GetUserName(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user_info")
	if err != nil {
		return "", fmt.Errorf("cookie not found")
	}


	parts := strings.Split(cookie.Value, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid cookie format")
	}

	admin_name := SanitizeStringTwo(parts[1])

	return admin_name, nil
}
