package main

import (
	"net/http"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"
	"fmt"
	_ "strconv"

	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var orgIds []uint16
var orgInfos []string
var orgInfos_low []string

var crgIds []uint16
var crgInfos []string
var crgInfos_low []string

var bncIds []uint16
var bncInfos []string
var bncInfos_low []string

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func main() {

	fmt.Println("Elastic/Server iniciando...")
	fmt.Println("Carregando tabela(orgaos)...")

	db, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_data?allowNativePasswords=true")

	var transformed_string string
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	// Carregar orgaos

	sql := "SELECT id, info FROM orgaos_conc WHERE 1 ORDER BY id ASC"
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	i := 0
	var _orgId uint16
	var _orgInfo string

	for rows.Next() {
		rows.Scan(&_orgId, &_orgInfo)
		orgIds = append(orgIds, _orgId)
		orgInfos = append(orgInfos, _orgInfo)
		transformed_string = strings.ToLower(_orgInfo)
		transformed_string, _, _ := transform.String(t, transformed_string)
		orgInfos_low = append(orgInfos_low, transformed_string)
		i++
	}

	// Carregar cargos

	sql = "SELECT id, crg FROM cargos_conc WHERE 1 ORDER BY id ASC"
	rows, err = db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	i = 0
	var _crgId uint16
	var _crgInfo string

	for rows.Next() {
		rows.Scan(&_crgId, &_crgInfo)
		crgIds = append(crgIds, _crgId)
		crgInfos = append(crgInfos, _crgInfo)
		transformed_string = strings.ToLower(_crgInfo)
		transformed_string, _, _ := transform.String(t, transformed_string)
		crgInfos_low = append(crgInfos_low, transformed_string)
		i++
	}

	// Carregar bancas

	sql = "SELECT id, info FROM bancas_conc WHERE 1 ORDER BY id ASC"
	rows, err = db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	i = 0
	var _bncId uint16
	var _bncInfo string

	for rows.Next() {
		rows.Scan(&_bncId, &_bncInfo)
		bncIds = append(bncIds, _bncId)
		bncInfos = append(bncInfos, _bncInfo)
		transformed_string = strings.ToLower(_bncInfo)
		transformed_string, _, _ := transform.String(t, transformed_string)
		bncInfos_low = append(bncInfos_low, transformed_string)
		i++
	}

	http.HandleFunc("/org", orgaos)
	http.HandleFunc("/crg", cargos)
	http.HandleFunc("/bnc", bancas)

	http.ListenAndServe(":4050", nil)
}

func response(m *map[int]interface{}, w http.ResponseWriter) {
	data, _ := json.Marshal(m)
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Connection", "close")
	w.Write([]byte(data))
}

func removeAcentos(s string) string {
	t := transformString(norm.NFD, s)
	return strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1 // remove acento
		}
		return r
	}, t)
}

func transformString(form norm.Form, s string) string {
	return form.String(s)
}

func orgaos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := make(map[int]interface{})

	if r.Method != "GET" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	srch := ""
	_ = srch

	srch_, ok := r.URL.Query()["srch"]
	if !ok {
		w.Write([]byte("param.err(srch)"))
		return
	} else {
		srch = srch_[0]
	}

	if len(srch) < 3 {
		m[0] = 101
		m[1] = "Poucos caracteres"
		response(&m, w)
		return
	}

	srch = strings.ToLower(removeAcentos(srch_[0])) // <-- AQUI
	srch_split := strings.Split(srch, " ")
	split_length := len(srch_split)

	for i := 0; i < len(orgIds); i++ {
		org_info_normalizado := removeAcentos(strings.ToLower(orgInfos[i]))

		matches := 0
		for _, match := range srch_split {
			if strings.Contains(org_info_normalizado, match) {
				matches++
			}
		}
		if matches == split_length {
			m[int(orgIds[i])] = orgInfos[i]
		}
	}

	response(&m, w)
}

func cargos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := make(map[int]interface{})

	if r.Method != "GET" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	srch := ""
	_ = srch

	srch_, ok := r.URL.Query()["srch"]
	if !ok {
		w.Write([]byte("param.err(srch)"))
		return
	} else {
		srch = srch_[0]
	}

	if len(srch) < 3 {
		m[0] = 101
		m[1] = "Poucos caracteres"
		response(&m, w)
		return
	}

	srch = strings.ToLower(removeAcentos(srch_[0])) // <-- AQUI
	srch_split := strings.Split(srch, " ")
	split_length := len(srch_split)

	for i := 0; i < len(crgIds); i++ {
		crg_info_normalizado := removeAcentos(strings.ToLower(crgInfos[i]))

		matches := 0
		for _, match := range srch_split {
			if strings.Contains(crg_info_normalizado, match) {
				matches++
			}
		}
		if matches == split_length {
			m[int(crgIds[i])] = crgInfos[i]
		}
	}

	response(&m, w)
}

func bancas(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := make(map[int]interface{})

	if r.Method != "GET" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	srch := ""
	_ = srch

	srch_, ok := r.URL.Query()["srch"]
	if !ok {
		w.Write([]byte("param.err(srch)"))
		return
	} else {
		srch = srch_[0]
	}

	if len(srch) < 3 {
		m[0] = 101
		m[1] = "Poucos caracteres"
		response(&m, w)
		return
	}

	srch = strings.ToLower(removeAcentos(srch_[0])) // <-- AQUI
	srch_split := strings.Split(srch, " ")
	split_length := len(srch_split)

	for i := 0; i < len(bncIds); i++ {
		cbnc_info_normalizado := removeAcentos(strings.ToLower(bncInfos[i]))

		matches := 0
		for _, match := range srch_split {
			if strings.Contains(cbnc_info_normalizado, match) {
				matches++
			}
		}
		if matches == split_length {
			m[int(bncIds[i])] = bncInfos[i]
		}
	}

	response(&m, w)
}
