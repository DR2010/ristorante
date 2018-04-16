// Package main is the main package
// -------------------------------------
// .../restauranteapi/btccotacaohandler.go
// -------------------------------------
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"restauranteapi/btcmarkets"
	"time"
)

// Hbtcpreorderadd add orders
func Hbtcpreorderadd(httpwriter http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	bodybyte, _ := ioutil.ReadAll(req.Body)

	var objtoaction btcmarkets.DCPreOrder

	err = json.Unmarshal(bodybyte, &objtoaction)

	ret := btcmarkets.PreOrderAdd(redisclient, objtoaction)

	if ret.IsSuccessful == "Y" {
		// do something

		fmt.Println("Order added successfully:")

		type RespAddOrder struct {
			ID string
		}

	}

	return
}

// Hbtcpreorderlist is a function to return a list of dishes
func Hbtcpreorderlist(httpwriter http.ResponseWriter, req *http.Request) {
	var cotacaolist = btcmarkets.PreorderGetAll(redisclient)
	json.NewEncoder(httpwriter).Encode(&cotacaolist)
}

// Hbtccotacaoadd is
func Hbtccotacaoadd(httpwriter http.ResponseWriter, req *http.Request) {

	cotacaotoadd := btcmarkets.BalanceCrypto{}

	cotacaotoadd.Currency = req.FormValue("cryptoCurrency")
	cotacaotoadd.Balance = req.FormValue("cryptoBalance")
	cotacaotoadd.CotacaoAtual = req.FormValue("cryptoCotacaoAtual")
	cotacaotoadd.ValueInCashAUD = req.FormValue("cryptoValueInCashAUD")
	cotacaotoadd.BestAsk = req.FormValue("cryptoBestAsk")
	cotacaotoadd.BestBid = req.FormValue("cryptoBestBid")
	cotacaotoadd.Volume24 = req.FormValue("cryptoVolume24")
	cotacaotoadd.DateTime = time.Now().String()
	cotacaotoadd.Rotina = req.FormValue("rotina")

	// params := req.URL.Query()
	// cotacaotoadd.Currency = params.Get("Currency")
	// cotacaotoadd.Balance = params.Get("Balance")
	// cotacaotoadd.CotacaoAtual = params.Get("CotacaoAtual")
	// cotacaotoadd.DateTime = params.Get("DateTime")
	// cotacaotoadd.ValueInCashAUD = params.Get("ValueInCashAUD")

	ret := btcmarkets.CryptoCotacaoAdd(redisclient, cotacaotoadd)

	if ret.IsSuccessful == "Y" {
		// do something
	}
}

// Hbtccotacaolist is a function to return a list of dishes
func Hbtccotacaolist(httpwriter http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	var currency = params.Get("currency")
	var rows = params.Get("rows")

	var cotacaolist = btcmarkets.GetAll(redisclient, currency, rows)

	json.NewEncoder(httpwriter).Encode(&cotacaolist)
}

// Hbtccotacaolistbyday is a list to group it all
// func Hbtccotacaolistbyday(httpwriter http.ResponseWriter, req *http.Request) {

// 	params := req.URL.Query()
// 	var currency = params.Get("currency")

// 	// get all lines
// 	var cotacaolist = btcmarkets.GetAllNoLimit(redisclient, currency)

// 	// loop through lines and get all the

// 	var coindaymax = 0
// 	var coindaymin = 0
// 	var coindayclosing = 0

// 	for x := 1; x < len(cotacaolist); x++ {

// 		balance, _ := strconv.ParseFloat(cotacaolist[x].Balance, 32)
// 		balance1, _ := strconv.ParseFloat(cotacaolist[x].Balance, 32)

// 		var coindaymax = (cotacaolist[x].Balance)
// 		var coindaymin = 0
// 		var coindayclosing = 0

// 	}

// 	json.NewEncoder(httpwriter).Encode(&cotacaolist)
// }

// Hbtcimport is a function to return a list of dishes
func Hbtcimport(httpwriter http.ResponseWriter, req *http.Request) {

	var cotacaolist = btcmarkets.Import(redisclient)

	json.NewEncoder(httpwriter).Encode(&cotacaolist)
}

// Hbtccotacaolistdate is a function to return a list of dishes
func Hbtccotacaolistdate(httpwriter http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	var currency = params.Get("currency")
	var yeardaymonth = params.Get("yeardaymonth")
	var yeardaymonthend = params.Get("yeardaymonthend")

	var cotacaolist = btcmarkets.GetDayStats(redisclient, currency, yeardaymonth, yeardaymonthend)

	json.NewEncoder(httpwriter).Encode(&cotacaolist)
}
