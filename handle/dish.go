package handle

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/phuslu/log"
	"github.com/sanbei101/cau-ai-backend/database"
	"github.com/sanbei101/cau-ai-backend/utils/response"
	"github.com/sanbei101/cau-ai-backend/utils/validate"
	"gorm.io/datatypes"
)

var PG = database.PG

type Dish struct {
	ID      uuid.UUID      `json:"id" gorm:"primaryKey;default:uuidv7()"`
	Name    string         `json:"name"`
	Tag     string         `json:"tag"`
	Canteen datatypes.JSON `json:"canteen" gorm:"type:jsonb"`
}

func InitDish() {
	PG.Exec("DROP TABLE IF EXISTS dishes")
	if err := PG.AutoMigrate(&Dish{}); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate dish")
	}

	file, err := os.Open("resource/dishs.csv")
	if err != nil {
		log.Error().Err(err).Msg("failed to open dishs.csv")
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to read dishs.csv")
	}
	var dishs []Dish
	for _, record := range records[1:] {
		dishName := strings.TrimSpace(record[0])
		canteenStr := strings.TrimSpace(record[1])
		canteen := strings.Split(canteenStr, ",")
		tag := strings.TrimSpace(record[2])
		canteenJSON, _ := json.Marshal(canteen)
		dish := Dish{
			Name:    dishName,
			Tag:     tag,
			Canteen: canteenJSON,
		}
		dishs = append(dishs, dish)
	}
	if err := PG.CreateInBatches(dishs, 100).Error; err != nil {
		log.Fatal().Err(err).Msg("failed to create dishs")
	}
}

type ListDishReq struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Name     string `query:"name"`
	Canteen  string `query:"canteen"`
	Tag      string `query:"tag"`
}

func ListDish(c *gin.Context) {
	var req ListDishReq
	if err := validate.ParseQuery(c, &req); err != nil {
		log.Error().Err(err).Msg("failed to parse query")
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	pageInfo := response.PageInfo{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	baseQuery := PG.WithContext(c).Model(&Dish{})
	if req.Name != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Canteen != "" {
		baseQuery = baseQuery.Where("canteen @> jsonb_build_array(?::text)", req.Canteen)
	}
	if req.Tag != "" {
		baseQuery = baseQuery.Where("tag LIKE ?", "%"+req.Tag+"%")
	}
	var dishs []Dish
	var total int64
	err := baseQuery.Find(&dishs).Count(&total).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to find dishs")
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = baseQuery.Scopes(pageInfo.Paginate()).Find(&dishs).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to find dishs")
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(
		c,
		dishs,
		total,
		pageInfo.Page,
		pageInfo.PageSize,
	)
}
