package v1

import (
	"banner/internal/entity"
	"banner/internal/service"
	"banner/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BannerController struct {
	bannerService service.Service
	l             logger.Logger
}

func NewBannerController(bannerService service.Service, logger logger.Logger) *BannerController {
	return &BannerController{
		bannerService: bannerService,
		l:             logger,
	}
}
func (h *BannerController) createBanner(c *gin.Context) {
	token := c.GetBool("isAdmin")
	if !token {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	var banner entity.Banner
	err := c.ShouldBindJSON(&banner)
	if err != nil {
		h.l.Error("Failed to parse request data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	bannerID, err := h.bannerService.Save(c.Request.Context(), &banner)
	if err != nil {
		h.l.Error("Failed to create banner: %v", err)
		if err.Error() == "record with same featureId and tagId already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Баннер с таким featureId и tagId уже существует"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}
	banner.ID = bannerID
	h.l.Info("Banner created successfully")
	c.JSON(http.StatusCreated, gin.H{"banner_id": banner.ID})
}
func (h *BannerController) getBanner(c *gin.Context) {
	isAdmin, exists := c.Get("isAdmin")
	if !exists {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	lastRevision, err := strconv.ParseBool(c.DefaultQuery("use_last_revision", "false"))
	if err != nil {
		lastRevision = false
	}
	tagID, err := strconv.ParseInt(c.Query("tag_id"), 10, 32)
	if err != nil {
		h.l.Error("Failed to parse tag ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректные данные tagId"})
		return
	}
	featureID, err := strconv.ParseInt(c.Query("feature_id"), 10, 32)
	if err != nil {
		h.l.Error("Failed to parse feature ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректные данные featureId"})
		return
	}
	content, err := h.bannerService.GetForUser(c.Request.Context(), int32(tagID), int32(featureID), isAdmin.(bool), lastRevision)
	if err != nil {
		if err.Error() == "no banner found" {
			h.l.Info("No banner found for tag ID: %d, feature ID: %d", tagID, featureID)
			c.JSON(http.StatusNotFound, nil)
			return
		}
		h.l.Error("Failed to get content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}
	h.l.Info("Content retrieved successfully")
	c.JSON(http.StatusOK, content)

}
func (h *BannerController) getBanners(c *gin.Context) {
	token := c.GetBool("isAdmin")
	if !token {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	featureID, _ := strconv.ParseInt(c.DefaultQuery("feature_id", "0"), 10, 32)
	tagID, _ := strconv.ParseInt(c.DefaultQuery("tag_id", "0"), 10, 32)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	var featureIDPtr, tagIDPtr, limitPtr, offsetPtr *int32
	if featureID != 0 {
		featureIDConverted := int32(featureID)
		featureIDPtr = &featureIDConverted
	}
	if tagID != 0 {
		tagIDConverted := int32(tagID)
		tagIDPtr = &tagIDConverted
	}
	if limit != 0 {
		limitConverted := int32(limit)
		limitPtr = &limitConverted
	}
	if offset != 0 {
		offsetConverted := int32(offset)
		offsetPtr = &offsetConverted
	}
	banners, err := h.bannerService.GetBanners(c.Request.Context(), featureIDPtr, tagIDPtr, limitPtr, offsetPtr)
	if err != nil {
		h.l.Error("Failed to get banners: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}
	h.l.Info("Banners retrieved successfully")
	c.JSON(http.StatusOK, banners)
}
func (h *BannerController) deleteBanner(c *gin.Context) {
	token := c.GetBool("isAdmin")
	if !token {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	bannerID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		h.l.Error("Failed to parse banner ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные bannerID"})
		return
	}
	err = h.bannerService.Delete(c.Request.Context(), int32(bannerID))
	if err != nil {
		if err.Error() == "no banner found" {
			h.l.Info("No banner found with ID: %d", bannerID)
			c.JSON(http.StatusNotFound, nil)
			return
		}
		h.l.Error("Failed to delete banner: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}
	h.l.Info("Banner deleted successfully")
	c.JSON(http.StatusNoContent, nil)
}
func (h *BannerController) updateBanner(c *gin.Context) {
	token := c.GetBool("isAdmin")
	if !token {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	bannerID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		h.l.Error("Failed to parse banner ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные bannerID"})
		return
	}
	var bannerUpdate entity.BannerUpdate
	if err := c.ShouldBindJSON(&bannerUpdate); err != nil {
		h.l.Error("Failed to bind banner JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при разборе данных баннера"})
		return
	}
	bannerIDConverted := int32(bannerID)
	bannerUpdate.ID = &bannerIDConverted
	err = h.bannerService.Update(c.Request.Context(), &bannerUpdate)
	if err != nil {
		if err.Error() == "no banner found" {
			h.l.Info("No banner found with ID: %d", bannerID)
			c.JSON(http.StatusNotFound, nil)
			return
		}
		h.l.Error("Failed to update banner: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}
	h.l.Info("Banner updated successfully")
	c.JSON(http.StatusNoContent, nil)
}
func (h *BannerController) getBannersHistoryByID(c *gin.Context) {
	token := c.GetBool("isAdmin")
	if !token {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	bannerID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		h.l.Error("Failed to parse banner ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные bannerID"})
		return
	}
	banners, err := h.bannerService.GetBannersHistoryByID(c.Request.Context(), int32(bannerID))
	if err != nil {
		h.l.Error("Failed to get banner history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}
	h.l.Info("Banner history retrieved successfully")
	c.JSON(http.StatusOK, banners)
}
