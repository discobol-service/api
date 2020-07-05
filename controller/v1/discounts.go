package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/discobol-service/api/entity"
	"github.com/discobol-service/api/io"
	"net/http"
	"encoding/json"
	"sort"
	"log"
)

/*
 *	GET /v1/discounts
 *
 *	Взять у бандита рекомендации,
 *	соорудить из этого ответ в жсоне:
 */
func GetDiscounts(c *gin.Context) {
	domain := c.DefaultQuery("domain", "default")

	// Получить данные сначала из базы:
	discountsMap, err := entity.GetDiscountsMap(100, domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
		})
		return
	}

	if len(discountsMap) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "No Data",
		})
		return
	}

	// Потом стату по этим данным от бандита:
	requestJSON, err := getBanditStatRequestJSON(discountsMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server Error",
		})

		return
	}

	bandit := io.GetBandit("http://localhost:4444")
	stat, err := bandit.GetStat(domain, requestJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "System Error",
		})
		log.Println(err)

		return
	}

	var ds = make([]entity.Discount, 0, len(stat))
	for _, s := range stat {
		d := *discountsMap[s.Arm]
		d.Scores = float64(d.Weight) + s.Scores
		ds = append(ds, d)
	}

	sort.Slice(ds, func(i, j int) bool {
		return ds[i].Scores > ds[j].Scores
	})

	log.Printf("ds = %+v", ds)

	c.JSON(http.StatusOK, ds)
}

func getBanditStatRequestJSON (dsMap map[string]*entity.Discount) ([]byte, error) {
	arms := make([]string, 0, len(dsMap))
	for arm := range dsMap {
		arms = append(arms, arm)
	}

	requestJSON, err := json.Marshal(arms)
	if err != nil {
		return []byte{}, err
	}

	return requestJSON, nil
}