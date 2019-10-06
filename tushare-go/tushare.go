package tushare

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const tsurl string = "http://api.tushare.pro"

type Tushare struct {
	Token string
}

type postbody struct {
	ApiName string                 `json:"api_name"`
	Token   string                 `json:"token"`
	Params  map[string]interface{} `json:"params"`
	Fields  []string               `json:"fields"`
}

type Tdata struct {
	Fields []string   `json:"fields"`
	Items  [][]string `json:"items"`
}

type Tresp struct {
	Data      Tdata  `json:"data"`
	RequestId string `json:"request_id"`
	Code      int    `json:"code`
	Msg       string `json:"msg"`
}

func (ts *Tushare) Query(apiname string, params map[string]interface{}, fields []string) (*Tresp, error) {

	postdata, err := json.Marshal(
		&postbody{
			ApiName: apiname,
			Token:   ts.Token,
			Params:  params,
			Fields:  fields,
		},
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(tsurl, "application/json;charset=utf-8;Content-Encoding=gzip", bytes.NewBuffer(postdata))
	if err != nil {
		return nil, err
	}

	m := &Tresp{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (ts *Tushare) GetStcokBasic(params map[string]interface{}) ([]StockBasic, error) {

	res, err := ts.Query("stock_basic", params, []string{"ts_code", "symbol", "name", "area", "industry", "fullname", "enname", "market", "exchange", "curr_type", "list_status", "list_date", "delist_date", "is_hs"})

	if err != nil {
		return nil, err
	} else if res.Code != 0 {
		return nil, errors.New(res.Msg)
	}

	header := res.Data.Fields
	rows := res.Data.Items

	sbList := []StockBasic{}
	for _, r := range rows {
		row := make(map[string]interface{})
		for index, h := range header {
			row[h] = r[index]
		}
		sbList = append(sbList, StockBasic{
			TsCode:     row["ts_code"].(string),
			Symbol:     row["symbol"].(string),
			Name:       row["name"].(string),
			Area:       row["area"].(string),
			Industry:   row["industry"].(string),
			Fullname:   row["fullname"].(string),
			Enname:     row["enname"].(string),
			Market:     row["market"].(string),
			Exchange:   row["exchange"].(string),
			CurrType:   row["curr_type"].(string),
			ListStatus: row["list_status"].(string),
			ListDate:   row["list_date"].(string),
			DelistDate: row["delist_date"].(string),
			Ishs:       row["is_hs"].(string),
		})
	}
	return sbList, nil
}

type StockBasic struct {
	TsCode     string `json:"ts_code"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	Area       string `json:"area"`
	Industry   string `json:"industry"`
	Fullname   string `json:"fullname"`
	Enname     string `json:"enname"`
	Market     string `json:"market"`
	Exchange   string `json:"exchange"`
	CurrType   string `json:"curr_type"`
	ListStatus string `json:"list_status"`
	ListDate   string `json:"list_date"`
	DelistDate string `json:"delist_date"`
	Ishs       string `json:"is_hs"`
}
