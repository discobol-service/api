package entity

import (
	"github.com/discobol-service/api/io"
	"context"
	"log"
)

// Rec - ...
type Discount struct {
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
	Weight      int64    `json:"weight"`
	Scores      float64  `json:"scores"`
	//Meta []interface{} `json:"meta"`
}


func GetDiscountsMap(limit int, regionName string) (map[string]*Discount, error) {
	discountsMap := make(map[string]*Discount)

	sql := `

		select
			d.udid, d.title, d.description,
			d.thumb, d.type, d.url,
			d.amount_from, d.amount_to,
			d.currency,

			d.brands, d.tags, d.weight,

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

	pg := io.GetPg()
	rows, err := pg.Connect.Query(context.Background(), sql, regionName, limit)
	if err != nil {
		log.Println("Database error: " + err.Error())
		return discountsMap, err
	}

	for rows.Next() {
		d := &Discount{}
		err = rows.Scan(
			&d.UDID, &d.Title, &d.Description, &d.Thumb, &d.Type, &d.URL,
			&d.AmountFrom, &d.AmountTo, &d.Currency, &d.Brands, &d.Tags, &d.Weight,
			&d.SellerName, &d.SellerSite, &d.SellerLogo,
		)

		if err != nil {
			log.Println(err.Error())
			return discountsMap, err
		}

		discountsMap[d.UDID] = d
	}

	return discountsMap, nil
}











