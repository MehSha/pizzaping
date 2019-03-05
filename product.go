package main

import (
	"log"

	"github.com/MehSha/basicdam"
	"github.com/jmoiron/sqlx"
)

type Product struct {
	VendorID       int    `db:"id"`
	VendorCode     string `db:"code"`
	VendorTitle    string `db:"title"`
	CategoryID     int    `db:"categoryid"`
	CategoryTitle  string `db:"categorytitle"`
	ProductID      int    `db:"pid"`
	ProductTitle   string `db:"Producttitle"`
	VariationID    int    `db:"pvId"`
	VariationTitle int    `db:"pvtitle"`
	Price          int    `db:"pvprice"`
}

var productDAM *basicdam.BasicDAM

func initProduct(db *sqlx.DB) {
	productDAM = basicdam.NewDAM(&Product{}, db)
	err := productDAM.SyncDB()
	if err != nil {
		log.Fatalln(err)
	}
}
