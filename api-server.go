package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

// Rec - ...
type Rec struct {
	UDID        string   `json:"udid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Thumb       string   `json:"thumb"`
	Type        int64    `json:"type"`
	URL         string   `json:"url"`
	AmountFrom  float64  `json:"amount_from"`
	AmountTo    float64  `json:"amount_to"`
	Currency    string   `json:"currency"`
	SellerName  string   `json:"seller_name"`
	SellerSite  string   `json:"seller_site"`
	SellerLogo  string   `json:"seller_logo"`
	Brands      []string `json:"brands"`
	Tags        []string `json:"tags"`
	//Meta []interface{} `json:"meta"`
}

type recsRequest struct {
	RegionID int64
	UserID   int64
}

type banditStat struct {
	Arm    string  `json:"arm"`
	Scores float64 `json:"scores"`
}

var db *pgx.Conn

func init() {
	conn, err := pgx.Connect(context.Background(), "postgresql://shootnix:12345@localhost/discobol")
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err.Error())
		os.Exit(1)
	}
	db = conn
}

func main() {
	router := gin.Default()

	defer db.Close(context.Background())

	router.GET("/v1/recs", recsV1)

	router.Run(":8080")
}

// Взять у бандита рекомендации,
// соорудить из этого ответ в жсоне:
func recsV1(c *gin.Context) {
	// Получить данные сначала из базы
	recsMap, err := getLatestRecs(100, c.DefaultQuery("domain", "default"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	if len(recsMap) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "no recs",
		})
		return
	}

	arms := make([]string, 0, len(recsMap))
	for k := range recsMap {
		arms = append(arms, k)
	}

	jsonStr, err := json.Marshal(arms)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "error marshaling json",
		})
		return
	}

	// Взять веса урлов у бандита
	req, err := http.NewRequest("POST", "http://localhost:4444/stat/list/default", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "bandit request error",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.JSON(404, gin.H{
			"message": "recs not found",
		})
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	stat := []banditStat{}
	if err = json.Unmarshal(body, &stat); err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
		return
	}

	var recs = make([]Rec, 0, len(stat))
	for _, s := range stat {
		recs = append(recs, *recsMap[s.Arm])
	}

	log.Printf("recs = %+v", recs)

	c.JSON(http.StatusOK, recs)
}

func getLatestRecs(limit int, domain string) (map[string]*Rec, error) {
	dbRecsMap := make(map[string]*Rec)

	sql := `

		select
			d.udid, d.title, d.description,
			d.thumb, d.type, d.url,
			d.amount_from, d.amount_to,
			d.currency, 

			d.brands, d.tags,
			
			s.name as seller_name,
			s.site as seller_site,
			s.logo as seller_logo
		from discounts d
		join regions r on r.id = d.region_id 
		join sellers s on s.id = d.seller_id
		where 
			r.name = $1 and 
			d.is_deleted = false and 
			d.is_active = true and
			s.is_deleted = false and
			s.is_active = true
		order by ctime desc 
		limit $2
	
	`

	rows, err := db.Query(context.Background(), sql, domain, limit)
	if err != nil {
		log.Println("Database error: " + err.Error())
		return dbRecsMap, err
	}

	for rows.Next() {
		r := new(Rec)
		err = rows.Scan(
			&r.UDID, &r.Title, &r.Description, &r.Thumb, &r.Type, &r.URL,
			&r.AmountFrom, &r.AmountTo, &r.Currency, &r.Brands, &r.Tags,
			&r.SellerName, &r.SellerSite, &r.SellerLogo,
		)

		if err != nil {
			log.Println(err.Error())
			return dbRecsMap, err
		}

		dbRecsMap[r.UDID] = r
	}

	return dbRecsMap, nil
}
