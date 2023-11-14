package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func checkAuth(c *gin.Context) {
	session := sessions.Default(c)
	authenticated := session.Get("authenticated")
	if authenticated == true {
		c.JSON(http.StatusOK, gin.H{"authenticated": true})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
	}
}

func login(c *gin.Context) {
	var loginForm struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	// This will bind the form data to the struct
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check the credentials: if they are valid, set the session
	if loginForm.Username == "test" && loginForm.Password == "dummy" {
		session := sessions.Default(c)
		session.Set("authenticated", true)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
	}
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) ([]byte, error) {
	// Generate a salt with a cost of 12 (adjust the cost as needed)
	saltCost := 12
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), saltCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// checkPassword checks if a password matches the hashed password
func checkPassword(password string, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}

/* Example password hashing and checking
// Create and hash a password
    password := "mySecretPassword"
    hashedPassword, err := hashPassword(password)
    if err != nil {
        fmt.Println("Error hashing password:", err)
        return
    }

    // Create a new user
    newUser := User{
        Email:        "user@example.com",
        PasswordHash: hashedPassword,
    }

    // Insert the user into the database
    db.Create(&newUser)

    // Query the user from the database
    var queriedUser User
    db.Where("email = ?", "user@example.com").First(&queriedUser)

    // Verify the password
    isValid := checkPassword("mySecretPassword", queriedUser.PasswordHash)
    if isValid {
        fmt.Println("Password is valid.")
    } else {
        fmt.Println("Password is invalid.")
    }
*/
