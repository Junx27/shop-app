package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/helper"
	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	repository entity.MenuRepository
}

func NewEventHandler(repository entity.MenuRepository) *MenuHandler {
	return &MenuHandler{repository: repository}
}

func (h *MenuHandler) GetMany(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	nameFilter := ctx.DefaultQuery("name", "")
	categoryFilter := ctx.DefaultQuery("category", "")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	menus, totalItems, err := h.repository.GetMany(ctx, pageInt, limitInt, nameFilter, categoryFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limitInt)))

	if pageInt > totalPages {
		ctx.JSON(http.StatusNotFound, helper.FailedResponse("Data not found"))
		return
	}
	ctx.JSON(http.StatusOK, helper.PaginationResponse("Fetch data successfully", pageInt, limitInt, totalPages, totalItems, menus))
}

func (h *MenuHandler) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	menu, err := h.repository.GetOne(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}

	ctx.JSON(http.StatusOK, helper.SuccessResponse(("Fetch data successfully"), menu))
}

func hashString(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func (h *MenuHandler) CreateOne(ctx *gin.Context) {
	menu := &entity.Menu{}
	if err := ctx.ShouldBind(&menu); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Invalid input"))
		return
	}

	userID, err := helper.GetUserIDFromCookie(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	menu.UserID = userID
	menu.Name = ctx.PostForm("name")
	menu.Price, _ = strconv.Atoi(ctx.PostForm("price"))
	menu.Category = ctx.PostForm("category")
	menu.Quantity, _ = strconv.Atoi(ctx.PostForm("quantity"))
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "no image file provided"})
		return
	}

	hashedName := hashString(strconv.Itoa(int(userID)) + file.Filename)
	fileExtension := strings.ToLower(strings.Split(file.Filename, ".")[1])
	hashedFileName := hashedName + "." + fileExtension
	uploadPath := "./uploads/"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create upload directory"})
		return
	}
	filePath := uploadPath + hashedFileName
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to upload image"})
		return
	}
	menu.Image = filePath
	createMenu, err := h.repository.CreateOne(ctx, menu)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to create data"))
		return
	}

	ctx.JSON(http.StatusOK, helper.SuccessResponse("Create data successfully", createMenu))
}

func (h *MenuHandler) UpdateOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	menu, err := h.repository.GetOne(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	file, _ := ctx.FormFile("image")
	if file != nil {
		if err := os.MkdirAll("./uploads/", os.ModePerm); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create upload directory"})
			return
		}
		if menu.Image != "" {
			if err := os.Remove(menu.Image); err != nil {
				log.Println("Failed to delete old image:", err)
			}
		}

		hashedName := hashString(strconv.Itoa(int(menu.UserID)) + file.Filename)
		fileExtension := strings.ToLower(strings.Split(file.Filename, ".")[1])
		hashedFileName := hashedName + "." + fileExtension
		uploadPath := "./uploads/"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create upload directory"})
			return
		}

		filePath := uploadPath + hashedFileName
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to upload image"})
			return
		}
		menu.Image = filePath
	}

	updateData := entity.Menu{
		Name: ctx.PostForm("name"),
		Price: func() int {
			price, _ := strconv.Atoi(ctx.PostForm("price"))
			return price
		}(),
		Category: ctx.PostForm("category"),
		Quantity: func() int {
			quantity, _ := strconv.Atoi(ctx.PostForm("quantity"))
			return quantity
		}(),
		Image: menu.Image,
	}
	if err := ctx.ShouldBind(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Invalid input"))
		return
	}

	updateFields := map[string]interface{}{
		"id":       menu.ID,
		"user_id":  menu.UserID,
		"name":     updateData.Name,
		"price":    updateData.Price,
		"category": updateData.Category,
		"quantity": updateData.Quantity,
		"image":    updateData.Image,
	}

	updatedMenu, err := h.repository.UpdateOne(ctx, uint(id), updateFields)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to update data"))
		return
	}

	ctx.JSON(http.StatusOK, helper.SuccessResponse("Update data successfully", updatedMenu))
}

func (h *MenuHandler) DeleteOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	menu, err := h.repository.GetOne(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	if menu.Image != "" {
		if err := os.Remove(menu.Image); err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to delete image"))
			return
		}
	}
	err = h.repository.DeleteOne(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to delete data"))
		return
	}

	ctx.JSON(http.StatusOK, helper.SuccessResponse("Delete data successfully", nil))
}

func (h *MenuHandler) DownloadImage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	menu, err := h.repository.GetOne(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	filePath := menu.Image
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
		return
	}

	_, fileName := filepath.Split(filePath)
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.File(filePath)
}
