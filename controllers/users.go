package controllers

import (
	"context"
	"log"
	"net/http"
	"psql-typesense/database"
	"psql-typesense/models"
	"psql-typesense/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var user models.Users

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "Invalid Request Body"})
		return
	}

	validate := validator.New()

	err = validate.Struct(user)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&user).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email already exists already exists"})
		return
	}

	hashPass, err := utils.HashPassword(user.Password)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error encrypting password"})
		return
	}

	user.Password = hashPass
	user.UserCreatedAt = time.Now().Format(time.RFC3339)

	err = database.DB.Create(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in creating account"})
		return
	}

	tysenseDoc := models.TypesenseUser{
		ID:             strconv.Itoa(user.ID),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Username:       user.Username,
		UserCreatedAt:  user.UserCreatedAt,
		InactiveStatus: user.InactiveStatus,
		Role:           user.Role,
	}

	_, err = database.TypesenseClient.Collection("users").Documents().Create(ctx, tysenseDoc, &api.DocumentIndexParameters{})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user index is not synced"})
		return
	}
	err = database.DB.Model(&models.Users{}).Where("id = ?", user.ID).Update("typesense_synced", true).Error
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registration completed"})
}

func SearchUsers(c *gin.Context) {

	query := c.DefaultQuery("q", "*")

	searchParameters := &api.SearchCollectionParams{
		Q:       pointer.String(query),
		Infix:   pointer.String("fallback"),
		QueryBy: pointer.String("first_name,last_name,user_created_at"),
		//FilterBy: pointer.String("num_employees:>100"),
		SortBy: pointer.String("user_created_at:desc"),
	}

	searchResult, err := database.TypesenseClient.Collection("users").Documents().Search(context.Background(), searchParameters)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not able to get users"})
		return
	}

	c.JSON(http.StatusOK, searchResult)
}

func UpdatePassword(c *gin.Context) {
	var updatepass models.UpdatePassword
	user := &models.Users{}
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Id is missing"})
		return
	}
	err := c.ShouldBindJSON(&updatepass)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "Invalid request body"})
		return
	}

	if updatepass.CurrentPassword == "" || updatepass.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both are required"})
		return
	}

	validate := validator.New()
	err = validate.Struct(updatepass)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	err = database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	err = utils.ComparePass(updatepass.CurrentPassword, user.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Incorrect username or password"})
		return
	}
	hashPass, err := utils.HashPassword(updatepass.NewPassword)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "Password not encrypted"})
		return
	}

	updatepass.CurrentPassword = hashPass
	err = database.DB.Model(&models.Users{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"password":            hashPass,
			"password_changed_at": time.Now().Format(time.RFC3339),
		}).Error

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func Login(c *gin.Context) {

	var userLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	err := c.ShouldBindJSON(&userLogin)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	validate := validator.New()
	err = validate.Struct(userLogin)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// if userLogin.Email == "" || userLogin.Password == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"Error": "Both are required"})
	// 	return
	// }

	var user models.Users

	err = database.DB.Where("email = ?", userLogin.Email).First(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "User not found"})
		return
	}
	err = utils.ComparePass(userLogin.Password, user.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Incorrect username or password"})
		return
	}
	token, err := utils.SignedToken(user.FirstName, user.LastName, user.Email, user.Username, user.Role)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Token not created"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Bearer",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   false,
		HttpOnly: true,
		//SameSite: http.SameSiteStrictMode,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"success": "login successfull",
		"token":   token,
	})
}

func Logout(c *gin.Context) {

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Second),
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		//SameSite: http.SameSiteStrictMode,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, gin.H{"success": "Logout successfully"})
}
