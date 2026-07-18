package handler

import (
	"fmt"
	"net/http"
	"pankreatitmed/internal/app/authctx"
	"pankreatitmed/internal/app/dto"
	"pankreatitmed/internal/app/dto/request"
	"pankreatitmed/internal/app/dto/response"
	"pankreatitmed/internal/app/mapper"
	"strconv"

	"github.com/gin-gonic/gin"

	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// CriteriaList godoc
// @Summary      Список услуг
// @Description  Возвращает список критериев (услуг) с фильтрацией по подстроке названия
// @Tags         services
// @Produce      json
// @Param        query   query     string  false  "Поиск по названию (ILIKE)"
// @Success      200 {object} response.SendPankreatitOrder
// @Failure      400 {object} map[string]any "bad request"
// @Failure      500 {object} map[string]any "internal error"
// @Router       /criteria [get]
func (h *Handler) CriteriaList(c *gin.Context) {
	var query request.GetCriteria
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list, err := h.svcs.Criteria.List(query.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := mapper.CriterionsToSendCrtierions(list)
	res := dto.List[response.SendCriterion]{Items: items}

	c.JSON(http.StatusOK, res)

}

// CriteriaGet godoc
// @Summary      Получить одну услугу
// @Tags         services
// @Produce      json
// @Param        id   path      int  true  "ID услуги"
// @Success      200 {object} response.SendCriterion
// @Failure      400 {object} map[string]any "bad request"
// @Failure      404 {object} map[string]any "not found"
// @Router       /criteria/{id} [get]
func (h *Handler) CriteriaGet(c *gin.Context) {
	var id request.GetCriterion
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	println("id", id.ID)
	criterion, err := h.svcs.Criteria.Get(id.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	crit := mapper.CritertionToSendCriterionLink(criterion)
	c.JSON(http.StatusOK, crit)
}

// CriteriaCreate godoc
// @Summary      Добавить услугу
// @Description  Создаёт новую услугу (без изображения)
// @Tags         services
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body  request.CreateCriterion  true  "Данные услуги"
// @Success      201 {string} string "created"
// @Failure      400 {object} map[string]any "validation error"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden"
// @Router       /criteria [post]

// TODO: возвращать критерию
func (h *Handler) CriteriaCreate(c *gin.Context) {
	var criterion request.CreateCriterion
	if err := c.ShouldBindJSON(&criterion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		println(1)
		return
	}
	crit, err := mapper.CreateCriterionToCriterion(criterion)
	if err != nil {
		println(2)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := h.svcs.Criteria.Create(&crit); err != nil {
		println(3)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusCreated)
}

// TODO: разобраться почему не кидает ошибку
// CriteriaUpdate godoc
// @Summary      Изменить услугу
// @Tags         services
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path  int                       true  "ID услуги"
// @Param        input  body  request.UpdateCriterion   true  "Изменяемые поля"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden"
// @Failure      404 {object} map[string]any "not found"
// @Router       /criteria/{id} [put]
func (h *Handler) CriteriaUpdate(c *gin.Context) {
	var id request.GetCriterion
	var criterion request.UpdateCriterion
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&criterion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svcs.Criteria.Update(id.ID, &criterion); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}

// TODO: перенести удаления изображения в services, сделать проверку на дублирующие изображения, чтобы не удалять, если у 1-го 2 изображения
// CriteriaDelete godoc
// @Summary      Удалить услугу (со встроенным удалением изображения)
// @Tags         services
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "ID услуги"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden"
// @Failure      404 {object} map[string]any "not found"
// @Router       /criteria/{id} [delete]
func (h *Handler) CriteriaDelete(c *gin.Context) {
	var id request.GetCriterion
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svcs.Criteria.Delete(id.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	client := connectMinio() // из предыдущего примера
	if err := h.svcs.Criteria.DeleteImage(client, id.ID, c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}

// TODO: настроить нормализацию последовательности БД
// AddCriteriaToDraft godoc
// @Summary      Добавить услугу в заявку-черновик
// @Description  Создаёт черновик автоматически (если нет) и добавляет выбранную услугу
// @Tags         services
// @Security     BearerAuth
// @Produce      json
// @Param        id   path int true "ID услуги"
// @Success      201 {string} string "created"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      404 {object} map[string]any "not found"
// @Router       /criteria/{id}/add-to-draft [post]
func (h *Handler) AddCriteriaToDraft(c *gin.Context) {
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	var id request.GetCriterion
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svcs.Criteria.ToDraft(id.ID, usr.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusCreated)
}

// TODO: сделать проверку на существовние критерия (если его нет, все равно добавляет изображение)
// TODO: если картинка уже есть и меняешь, в минио админе почему то изображение не отображается, хотя по ссылке оно висит
// UploadCriterionImage godoc
// @Summary      Загрузить/заменить изображение услуги
// @Description  Загружает файл в MinIO по ID услуги; старое изображение удаляется
// @Tags         services
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        id     path     int   true  "ID услуги"
// @Param        image  formData file  true  "Изображение (jpg/png/webp)"
// @Success      200 {object} map[string]any "url: ссылка на изображение"
// @Failure      400 {object} map[string]any "bad request / image is required"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden"
// @Failure      500 {object} map[string]any "minio/internal error"
// @Router       /criteria/{id}/image [post]
func (h *Handler) UploadCriterionImage(c *gin.Context) {
	// ID услуги из URL
	var id request.GetCriterion
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}

	// Открываем файл
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Загружаем в MinIO
	client := connectMinio() // из предыдущего примера
	bucket := "services-images"
	objectName := fmt.Sprintf("service_%s/%s", strconv.Itoa(int(id.ID)), fileHeader.Filename)

	_, err = client.PutObject(
		c, bucket, objectName, file, fileHeader.Size,
		minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	imgname := fmt.Sprintf("http://localhost:9000/%s/%s", bucket, objectName)
	crit := request.UpdateCriterion{ImageURL: &imgname}
	fmt.Println(crit.ImageURL)
	if err := h.svcs.Criteria.DeleteImage(client, id.ID, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.svcs.Criteria.Update(id.ID, &crit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"url":    fmt.Sprintf("http://localhost:9000/%s/%s", bucket, objectName),
	})
}

func connectMinio() *minio.Client {
	endpoint := "localhost:9000" // адрес контейнера
	accessKey := "minio"         // MINIO_ROOT_USER
	secretKey := "minio124"      // MINIO_ROOT_PASSWORD
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
