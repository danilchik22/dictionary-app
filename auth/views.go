package auth

import (
	"bytes"
	"dictionary_app/config"
	sl "dictionary_app/utils/logger"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserLogin(ctx *gin.Context) {
	logger := sl.GetLogger()
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}
	username := loginData.Username
	password := loginData.Password
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "username was missed",
		})
		return
	}
	if password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "password was missed",
		})
		return
	}

	urlAuthServer := config.GetConfig().AuthServerAddress
	requestBody := RequestBody{
		Username: username,
		Password: password,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlAuthServer+"/login", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Info(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "password is not corrected",
		})
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	var data JSONResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	if data.AccessToken == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error in auth server: access token is not correct",
		})
		return
	}
	if data.RefreshToken == "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "refresh token is not correct",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  data.AccessToken,
		"refresh_token": data.RefreshToken,
	})
	return

}

func CreateNewUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	if user.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Username is required",
		})
		return
	}

	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Password is required",
		})
		return
	}
	requestBody := User{
		Username: user.Username,
		Password: user.Password,
		Sex:      user.Sex,
		Age:      user.Age,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	urlAuthServer := config.GetConfig().AuthServerAddress

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlAuthServer+"/new_user", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if resp.StatusCode == http.StatusBadRequest {
		var data ErrorResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": data.Error,
		})
		return
	}
	var data NewUserResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": data.Message,
		"user_id": data.UserId,
	})
	return

}

// func Logout(ctx *gin.Context) {
// 	tokenString := ctx.GetHeader("Authorization")
// 	successful := c.service.Logout(tokenString)
// 	if successful == false {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"message": "logout is unseccessful",
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "logout is succecful",
// 	})
// 	return

// }

func Refresh(ctx *gin.Context) {
	var refreshData struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ShouldBindJSON(&refreshData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request format",
		})
		return
	}

	requestBody := RefreshBody{
		RefreshToken: refreshData.RefreshToken,
	}
	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	urlAuthServer := config.GetConfig().AuthServerAddress

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlAuthServer+"/refresh", bytes.NewBuffer(jsonData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if resp.StatusCode == http.StatusBadRequest {
		var data ErrorResponse
		err = json.Unmarshal(body, &data)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": data.Error,
		})
		return
	}
	var data JSONResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  data.AccessToken,
		"refresh_token": data.RefreshToken,
	})
	return
}
