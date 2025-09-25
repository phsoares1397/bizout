package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"

	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var db *sql.DB
var rdb *redis.Client
var ctx = context.Background()

func main() {

	var err error
	// Conexão com MySQL
	db, err = sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_users?allowNativePasswords=true")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tools/Server iniciando...")

	initRedis()

	http.HandleFunc("/editor", Editor)
	http.HandleFunc("/mkrs", markers)
	http.HandleFunc("/lists", lists)
	http.HandleFunc("/cmtsr", commentsRead)
	http.HandleFunc("/cmtsc", commentsCreate)
	http.HandleFunc("/cmtse", commentsEdit)
	http.HandleFunc("/cmtsd", commentsDelete)
	http.HandleFunc("/report", createReport)
	http.HandleFunc("/new_assts", newAssts)
	http.HandleFunc("/filters_read", filtersRead)

	http.ListenAndServe(":4040", nil)

	fmt.Println("Tudo ok...")
}

// ----------- REDIS INIT -----------
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // se configurou senha no redis.conf, coloque aqui
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Erro ao conectar no Redis:", err)
	}
	fmt.Println("Conectado ao Redis!")
}

func response(m *map[int]interface{}, w http.ResponseWriter) {
	data, _ := json.Marshal(m)
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Connection", "close")
	w.Write([]byte(data))
}

func captcha(token string) bool {

	resp, err := http.Get("https://www.google.com/recaptcha/api/siteverify?secret=6LdL49IkAAAAAHiUr0cPZn3kAlvErQ4BsGJEs4xz&response=" + token)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	capRes := string(body)
	var result map[string]interface{}
	json.Unmarshal([]byte(capRes), &result)

	if result["success"] != nil {
		if result["success"] == true {
			return true
		}
	}

	return true
}

func remove_data[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func newAssts(w http.ResponseWriter, r *http.Request) {

	m := make(map[int]interface{})

	if r.Method != "POST" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	op := r.FormValue("op")

	var dataIds []int
	json.Unmarshal([]byte(op), &dataIds)

	db, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_qst?allowNativePasswords=true")
	defer db.Close()

	sql := "SELECT id, assts_ids, assts_names, likes FROM new_assts_ids WHERE id = " + strconv.Itoa(dataIds[0])

	for _, i := range dataIds {
		sql += " OR id = " + strconv.Itoa(i)
	}

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	i := 0
	x := make(map[int]interface{})
	var nId int
	var nAssstIds string
	var nAsstsNames string
	var nLikes int

	for rows.Next() {
		x = make(map[int]interface{})
		rows.Scan(&nId, &nAssstIds, &nAsstsNames, &nLikes)
		x[0] = nId
		x[1] = nAssstIds
		x[2] = nAsstsNames
		x[3] = nLikes
		m[i] = x
		i++
	}

	response(&m, w)
}

func createReport(w http.ResponseWriter, r *http.Request) {
	m := make(map[int]interface{})

	if r.Method != "POST" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	token := r.FormValue("token")
	data := r.FormValue("data")
	if token == "" || data == "" {
		w.Write([]byte("param.err()"))
	} else {
		if !captcha(token) {
			m[0] = 111
			m[1] = "Captcha invalido"
			response(&m, w)
		} else {
			db, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_qst?allowNativePasswords=true")
			defer db.Close()

			_, err := db.Exec("INSERT INTO `pending` (`id`, `data`) VALUES (NULL, ?)", &data)
			if err != nil {
				m[0] = 320
				m[1] = "Erro ao inserir report :insert:tools:320"
				response(&m, w)
			} else {
				m[0] = 200
				response(&m, w)
			}
		}
	}

}

func commentsCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	m := make(map[int]interface{})

	if r.Method != "POST" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	token := r.FormValue("token")
	data := r.FormValue("data")

	if token == "" || data == "" {
		m[0] = 101
		m[1] = "Parâmetros ausentes"
		response(&m, w)
		return
	}

	// 🔐 valida captcha
	if !captcha(token) {
		m[0] = 111
		m[1] = "Captcha inválido"
		response(&m, w)
		return
	}

	// 🔑 valida sessão no Redis
	cookie, err := r.Cookie("session_id")
	if err != nil {
		m[0] = 401
		m[1] = "Sessão não encontrada"
		response(&m, w)
		return
	}

	sessionData, err := rdb.Get(ctx, "session:"+cookie.Value).Result()
	if err == redis.Nil {
		m[0] = 401
		m[1] = "Sessão inválida ou expirada"
		response(&m, w)
		return
	} else if err != nil {
		m[0] = 500
		m[1] = "Erro ao acessar Redis"
		response(&m, w)
		return
	}

	var session struct {
		UserID int `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		m[0] = 500
		m[1] = "Erro ao decodificar sessão"
		response(&m, w)
		return
	}

	// 🔎 pega dados básicos do usuário
	var user, fname string
	err = db.QueryRow("SELECT user, fname FROM udbMain WHERE id=?", session.UserID).Scan(&user, &fname)
	if err == sql.ErrNoRows {
		m[0] = 401
		m[1] = "Usuário não encontrado"
		response(&m, w)
		return
	} else if err != nil {
		m[0] = 500
		m[1] = "Erro ao buscar usuário"
		response(&m, w)
		return
	}

	// 🎯 processa criação do comentário
	qdb, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_qst?allowNativePasswords=true")
	defer qdb.Close()

	var dataJson []string
	json.Unmarshal([]byte(data), &dataJson)

	if len(dataJson) < 2 {
		m[0] = 107
		m[1] = "Erro nos dados do comentário"
		response(&m, w)
		return
	} else if dataJson[1] == "" {
		m[0] = 109
		m[1] = "Comentário vazio"
		response(&m, w)
		return
	}

	// estrutura "own" com informações do autor
	own := "[\"" + dataJson[0] + "\",\"" + user + "\",\"" + fname + "\",\"" + strconv.Itoa(session.UserID) + "\"]"

	table := "cmts_" + dataJson[2] // dataJson[2] = tabela/assunto da questão
	_, err = qdb.Exec("INSERT INTO "+table+" (`id`, `data`, `own`, `like`, `dislike`, `cmt`) VALUES (NULL, DEFAULT, ?, 0, 0, ?)", own, dataJson[1])
	if err != nil {
		m[0] = 101
		m[1] = "Erro ao tentar registrar comentário"
		response(&m, w)
		return
	}

	// pega ID do comentário recém-criado
	var cmtId int
	qdb.QueryRow("SELECT LAST_INSERT_ID()").Scan(&cmtId)

	// associa comentário à questão
	tableIds := table + "_ids"
	var checkIfExist int
	qdb.QueryRow("SELECT id FROM "+tableIds+" WHERE id = ?", dataJson[0]).Scan(&checkIfExist)

	if checkIfExist == 0 {
		_, err := qdb.Exec("INSERT INTO "+tableIds+" (`id`, `cmts`) VALUES (?, ?)", dataJson[0], "["+strconv.Itoa(cmtId)+"]")
		if err != nil {
			m[0] = 103
			m[1] = "Erro ao tentar vincular comentário à questão"
			response(&m, w)
			return
		}
	} else {
		var cmts string
		qdb.QueryRow("SELECT `cmts` FROM "+tableIds+" WHERE id = ?", dataJson[0]).Scan(&cmts)
		var cmtsJson []int
		json.Unmarshal([]byte(cmts), &cmtsJson)
		cmtsJson = append(cmtsJson, cmtId)
		cmtsString, _ := json.Marshal(cmtsJson)

		_, err := qdb.Exec("UPDATE "+tableIds+" SET `cmts` = ? WHERE id = ?", cmtsString, dataJson[0])
		if err != nil {
			m[0] = 103
			m[1] = "Erro ao atualizar lista de comentários da questão"
			response(&m, w)
			return
		}
	}

	// atualiza contador de comentários
	tableCont := "cont_" + dataJson[2]
	_, err = qdb.Exec("UPDATE "+tableCont+" SET `qcmt` = `qcmt` + 1 WHERE id = ?", dataJson[0])
	if err != nil {
		m[0] = 104
		m[1] = "Erro ao atualizar contador de comentários"
		response(&m, w)
		return
	}

	// insere no pending
	var dataPending []string
	dataPending = append(dataPending, "cmt_insert")
	dataPending = append(dataPending, dataJson[2]) // op / tipo de questão
	dataPending = append(dataPending, strconv.Itoa(cmtId))
	dataPending = append(dataPending, dataJson[0]) // id da questão
	dataPendingString, _ := json.Marshal(dataPending)

	_, err = qdb.Exec("INSERT INTO `pending` (`id`, `data`) VALUES (NULL, ?)", dataPendingString)
	if err != nil {
		m[0] = 103
		m[1] = "Erro ao tentar registrar no pending"
		response(&m, w)
		return
	}

	// sucesso 🎉
	var qmtQnt int
	qdb.QueryRow("SELECT `qcmt` FROM "+tableCont+" WHERE id = ?", dataJson[0]).Scan(&qmtQnt)

	m[0] = 200
	m[1] = qmtQnt
	response(&m, w)
}

func commentsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	m := make(map[int]interface{})

	if r.Method != "POST" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	token := r.FormValue("token")
	data := r.FormValue("data")

	if token == "" || data == "" {
		m[0] = 101
		m[1] = "Parâmetros ausentes"
		response(&m, w)
		return
	}

	// 🔐 valida captcha
	if !captcha(token) {
		m[0] = 111
		m[1] = "Captcha inválido"
		response(&m, w)
		return
	}

	// 🔑 valida sessão no Redis
	cookie, err := r.Cookie("session_id")
	if err != nil {
		m[0] = 401
		m[1] = "Sessão não encontrada"
		response(&m, w)
		return
	}

	sessionData, err := rdb.Get(ctx, "session:"+cookie.Value).Result()
	if err == redis.Nil {
		m[0] = 401
		m[1] = "Sessão inválida ou expirada"
		response(&m, w)
		return
	} else if err != nil {
		m[0] = 500
		m[1] = "Erro ao acessar Redis"
		response(&m, w)
		return
	}

	var session struct {
		UserID int `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		m[0] = 500
		m[1] = "Erro ao decodificar sessão"
		response(&m, w)
		return
	}

	// 🔎 pega dados básicos do usuário
	var user string
	err = db.QueryRow("SELECT user FROM udbMain WHERE id=?", session.UserID).Scan(&user)
	if err == sql.ErrNoRows {
		m[0] = 401
		m[1] = "Usuário não encontrado"
		response(&m, w)
		return
	} else if err != nil {
		m[0] = 500
		m[1] = "Erro ao buscar usuário"
		response(&m, w)
		return
	}

	// 🎯 processa deleção do comentário
	qdb, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_qst?allowNativePasswords=true")
	defer qdb.Close()

	var dataJson []string
	json.Unmarshal([]byte(data), &dataJson)

	if len(dataJson) < 3 {
		m[0] = 107
		m[1] = "Dados insuficientes para deletar comentário"
		response(&m, w)
		return
	}

	// dataJson[0] = id da questão
	// dataJson[1] = id do comentário
	// dataJson[2] = tabela/assunto

	table := "cmts_" + dataJson[2]
	tableIds := table + "_ids"
	tableCont := "cont_" + dataJson[2]

	// 1️⃣ verifica se o comentário existe
	var checkUser string
	err = qdb.QueryRow("SELECT own FROM "+table+" WHERE id=?", dataJson[1]).Scan(&checkUser)
	if err == sql.ErrNoRows {
		m[0] = 108
		m[1] = "Comentário não encontrado"
		response(&m, w)
		return
	} else if err != nil {
		m[0] = 500
		m[1] = "Erro ao buscar comentário"
		response(&m, w)
		return
	}

	// opcional: só o autor pode deletar
	if !strings.Contains(checkUser, strconv.Itoa(session.UserID)) {
		m[0] = 403
		m[1] = "Permissão negada"
		response(&m, w)
		return
	}

	// 2️⃣ remove comentário da tabela cmts
	_, err = qdb.Exec("DELETE FROM "+table+" WHERE id=?", dataJson[1])
	if err != nil {
		m[0] = 101
		m[1] = "Erro ao deletar comentário"
		response(&m, w)
		return
	}

	// 3️⃣ remove comentário da lista de ids da questão
	var cmts string
	qdb.QueryRow("SELECT `cmts` FROM "+tableIds+" WHERE id=?", dataJson[0]).Scan(&cmts)
	var cmtsJson []int
	json.Unmarshal([]byte(cmts), &cmtsJson)

	// remove o comentário do array
	idCmt, _ := strconv.Atoi(dataJson[1])
	newCmts := []int{}
	for _, c := range cmtsJson {
		if c != idCmt {
			newCmts = append(newCmts, c)
		}
	}
	cmtsString, _ := json.Marshal(newCmts)
	_, err = qdb.Exec("UPDATE "+tableIds+" SET `cmts`=? WHERE id=?", cmtsString, dataJson[0])
	if err != nil {
		m[0] = 103
		m[1] = "Erro ao atualizar lista de comentários"
		response(&m, w)
		return
	}

	// 4️⃣ decrementa contador de comentários
	_, err = qdb.Exec("UPDATE "+tableCont+" SET `qcmt` = GREATEST(`qcmt`-1,0) WHERE id=?", dataJson[0])
	if err != nil {
		m[0] = 104
		m[1] = "Erro ao atualizar contador de comentários"
		response(&m, w)
		return
	}

	// 5️⃣ registra no pending
	var dataPending []string
	dataPending = append(dataPending, "cmt_delete")
	dataPending = append(dataPending, dataJson[2])
	dataPending = append(dataPending, dataJson[1])
	dataPending = append(dataPending, dataJson[0])
	dataPendingString, _ := json.Marshal(dataPending)

	_, err = qdb.Exec("INSERT INTO `pending` (`id`, `data`) VALUES (NULL, ?)", dataPendingString)
	if err != nil {
		m[0] = 103
		m[1] = "Erro ao registrar no pending"
		response(&m, w)
		return
	}

	// sucesso
	m[0] = 200
	m[1] = "Comentário deletado com sucesso"
	response(&m, w)
}

func commentsEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	m := make(map[int]interface{})

	if r.Method != "POST" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	token := r.FormValue("token")
	data := r.FormValue("data")

	if token == "" || data == "" {
		m[0] = 101
		m[1] = "Parâmetros ausentes"
		response(&m, w)
		return
	}

	// 🔐 valida captcha
	if !captcha(token) {
		m[0] = 111
		m[1] = "Captcha inválido"
		response(&m, w)
		return
	}

	// 🔑 valida sessão (igual ao meHandler)
	cookie, err := r.Cookie("session_id")
	if err != nil {
		m[0] = 401
		m[1] = "Sessão não encontrada"
		response(&m, w)
		return
	}

	sessionData, err := rdb.Get(ctx, "session:"+cookie.Value).Result()
	if err == redis.Nil {
		m[0] = 401
		m[1] = "Sessão inválida ou expirada"
		response(&m, w)
		return
	} else if err != nil {
		m[0] = 500
		m[1] = "Erro ao acessar Redis"
		response(&m, w)
		return
	}

	var session struct {
		UserID int `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		m[0] = 500
		m[1] = "Erro ao decodificar sessão"
		response(&m, w)
		return
	}

	// 🔎 valida se usuário ainda existe no banco
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM udbMain WHERE id=?", session.UserID).Scan(&exists)
	if err != nil || exists == 0 {
		m[0] = 401
		m[1] = "Usuário não encontrado"
		response(&m, w)
		return
	}

	// 🎯 processa edição do comentário
	qdb, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_qst?allowNativePasswords=true")
	defer qdb.Close()

	var dataJson []string
	json.Unmarshal([]byte(data), &dataJson)

	var dataPending []string
	dataPending = append(dataPending, "cmt_edit")
	dataPending = append(dataPending, dataJson[0])                  // tabela
	dataPending = append(dataPending, dataJson[1])                  // comentário id
	dataPending = append(dataPending, strconv.Itoa(session.UserID)) // quem editou
	dataPendingString, _ := json.Marshal(dataPending)

	table := "cmts_" + dataJson[0]
	_, _ = qdb.Exec("INSERT INTO `pending` (`id`, `data`) VALUES (NULL, ?)", dataPendingString)
	_, err = qdb.Exec("UPDATE "+table+" SET cmt=? WHERE id=?", dataJson[2], dataJson[1])

	if err != nil {
		m[0] = 120
		m[1] = "Erro ao tentar editar comentário"
		fmt.Println(err)
		response(&m, w)
		return
	}

	m[0] = 200
	m[1] = "Comentário editado com sucesso"
	response(&m, w)
}

func commentsRead(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := make(map[int]interface{})

	op := ""
	id := ""

	id_, ok := r.URL.Query()["id"]
	op_, ok2 := r.URL.Query()["op"]
	if !ok || !ok2 {
		w.Write([]byte("param.err(op)"))
		return
	} else {
		id = id_[0]
		op = op_[0]
	}

	db, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_qst?allowNativePasswords=true")
	defer db.Close()

	var cmts string
	var cmtsIds []int
	table_ids := "cmts_" + op + "_ids"
	sql_ids := "SELECT `cmts` FROM `" + table_ids + "` WHERE id =" + id
	db.QueryRow(sql_ids).Scan(&cmts)

	if cmts == "" || cmts == "null" || cmts == "[]" {
		m[0] = 200
		response(&m, w)
	} else {
		json.Unmarshal([]byte(cmts), &cmtsIds)
		table := "cmts_" + op
		sql := "SELECT `id`, `data`, `own`, `like`, `dislike`, `cmt` FROM `" + table + "` WHERE `id`=" + strconv.Itoa(cmtsIds[0])
		cmtsIds = cmtsIds[1:]

		for _, i := range cmtsIds {
			sql = sql + " OR `id`=" + strconv.Itoa(i)
		}

		var id int
		var data string
		var own string
		var like int
		var dislike int
		var cmt string
		rows, err := db.Query(sql)
		defer rows.Close()
		if err != nil {
			panic(err)
		} else {
			m[0] = 200
			var i = 1
			for rows.Next() {
				rows.Scan(&id, &data, &own, &like, &dislike, &cmt)
				var x []interface{}
				x = append(x, id)
				x = append(x, data)
				x = append(x, own)
				x = append(x, like)
				x = append(x, dislike)
				x = append(x, cmt)
				m[i] = x
				i++
			}
			response(&m, w)
		}
	}

}

func lists(w http.ResponseWriter, r *http.Request) {
	m := make(map[int]interface{})

	if r.Method != "POST" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	op := r.FormValue("op")
	data := r.FormValue("data")
	var dataArr []interface{}
	json.Unmarshal([]byte(data), &dataArr)
	hash := r.FormValue("hash")
	uId := r.FormValue("u")
	if op == "" || data == "" || hash == "" || uId == "" {
		m[0] = 100
		m[1] = "Parâmetros invalidos"
		response(&m, w)
		return
	}

	db, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_users?allowNativePasswords=true")
	defer db.Close()

	if op == "1" { //Create list
		res, err := db.Exec("INSERT INTO `udata01` (`id`, `uid`,`hash`, `data`) VALUES (NULL, ?, ? , '[]')", uId, hash)
		_ = res
		if err != nil {
			panic(err)
		} else {
			var id int
			db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&id)
			var lists string
			if dataArr[1] == "conc" {
				db.QueryRow("SELECT `lists_conc` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
			} else if dataArr[1] == "vest" {
				db.QueryRow("SELECT `lists_vest` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
			} else if dataArr[1] == "enem" {
				db.QueryRow("SELECT `lists_enem` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
			}
			var y map[string]interface{}
			json.Unmarshal([]byte(lists), &y)
			y[strconv.Itoa(id)] = []interface{}{dataArr[0]}
			yJson, _ := json.Marshal(y)
			if dataArr[1] == "conc" {
				res, err = db.Exec("UPDATE `udataMain` SET `lists_conc` = ? WHERE id = ? AND hash = ?", string(yJson), uId, hash)
			} else if dataArr[1] == "vest" {
				res, err = db.Exec("UPDATE `udataMain` SET `lists_vest` = ? WHERE id = ? AND hash = ?", string(yJson), uId, hash)
			} else if dataArr[1] == "enem" {
				res, err = db.Exec("UPDATE `udataMain` SET `lists_enem` = ? WHERE id = ? AND hash = ?", string(yJson), uId, hash)
			}
			if err != nil {
				panic(err)
			} else {
				m[0] = "200"
				m[1] = y
				response(&m, w)
			}
		}
	} else if op == "2" { //Delete list

		var lists string
		if dataArr[1] == "conc" {
			db.QueryRow("SELECT `lists_conc` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
		} else if dataArr[1] == "vest" {
			db.QueryRow("SELECT `lists_vest` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
		} else if dataArr[1] == "enem" {
			db.QueryRow("SELECT `lists_enem` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
		}

		var listsObject = make(map[string]interface{})
		json.Unmarshal([]byte(lists), &listsObject)

		if len(dataArr) > 0 {
			daraArrString, ok := dataArr[0].(string)
			if ok {
				delete(listsObject, daraArrString)

				listsString, _ := json.Marshal(listsObject)
				listsString_ := string(listsString)

				var err any

				if dataArr[1] == "conc" {
					_, err = db.Exec("UPDATE `udataMain` SET `lists_conc` = ? WHERE id = ? AND hash = ?", listsString_, uId, hash)
				} else if dataArr[1] == "vest" {
					_, err = db.Exec("UPDATE `udataMain` SET `lists_vest` = ? WHERE id = ? AND hash = ?", listsString_, uId, hash)
				} else if dataArr[1] == "enem" {
					_, err = db.Exec("UPDATE `udataMain` SET `lists_enem` = ? WHERE id = ? AND hash = ?", listsString_, uId, hash)
				}

				if err != nil {
					panic(err)
					m[0] = 100
					response(&m, w)
				} else {
					lId, _ := strconv.Atoi(daraArrString)
					_, err = db.Exec("DELETE FROM `udata01` WHERE id = ? AND hash = ?", lId, hash)
					if err != nil {
						panic(err)
						m[0] = 100
						response(&m, w)
					} else {
						m[0] = 200
						m[1] = listsObject
						response(&m, w)
					}
				}
			}
		}
	} else if op == "3" { // Insert into list

		if len(dataArr) > 0 {
			var qsts string
			dataInt_list := int(dataArr[1].(float64))
			db.QueryRow("SELECT `data` FROM `udata01` WHERE id = ? AND hash = ?", dataInt_list, hash).Scan(&qsts)

			if qsts == "" {
				m[0] = 110
				response(&m, w)
				return
			}

			var qstsObject []int
			json.Unmarshal([]byte(qsts), &qstsObject)

			dataInt_qst := int(dataArr[0].(float64))
			if slices.Contains(qstsObject, dataInt_qst) {
				m[0] = 210
				response(&m, w)
			} else {
				qstsObject = append(qstsObject, dataInt_qst)
				listsString, _ := json.Marshal(qstsObject)
				_, err := db.Exec("UPDATE `udata01` SET `data` = ? WHERE id = ?", listsString, dataInt_list)
				if err != nil {
					panic(err)
					m[0] = 100
					response(&m, w)
				} else {
					m[0] = 200
					response(&m, w)
				}
			}
		} else {
			m[0] = 100
			m[1] = "Parâmetro(s) invalido(s)"
			response(&m, w)
		}
	} else if op == "4" { //Read lists names
		var lists string

		if len(dataArr) < 2 {
			m[0] = 110
			m[1] = "Algum parametro incorreto"
			response(&m, w)
			return
		}

		if dataArr[1] == "conc" {
			db.QueryRow("SELECT `lists_conc` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
		} else if dataArr[1] == "vest" {
			db.QueryRow("SELECT `lists_vest` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
		} else if dataArr[1] == "enem" {
			db.QueryRow("SELECT `lists_enem` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&lists)
		}
		var y map[string]interface{}
		json.Unmarshal([]byte(lists), &y)
		m[0] = "200"
		m[1] = y
		response(&m, w)
	} else if op == "5" { //remove from list
		if len(dataArr) > 0 {
			var qsts string
			dataInt_list := int(dataArr[1].(float64))
			db.QueryRow("SELECT `data` FROM `udata01` WHERE id = ? AND hash = ?", dataInt_list, hash).Scan(&qsts)

			if qsts == "" {
				m[0] = 110
				response(&m, w)
				return
			}

			var qstsObject []int
			json.Unmarshal([]byte(qsts), &qstsObject)

			dataInt_qst := int(dataArr[0].(float64))

			for i, element := range qstsObject {
				if element == dataInt_qst {
					qstsObject = remove_data(qstsObject, i)
					break
				}
			}

			listsString, _ := json.Marshal(qstsObject)
			_, err := db.Exec("UPDATE `udata01` SET `data` = ? WHERE id = ?", listsString, dataInt_list)
			if err != nil {
				panic(err)
				m[0] = 100
				response(&m, w)
			} else {
				m[0] = 200
				response(&m, w)
			}
		} else {
			m[0] = 100
			m[1] = "Parâmetro(s) invalido(s)"
			response(&m, w)
		}
	}
}

func markers(w http.ResponseWriter, r *http.Request) {
	m := make(map[int]interface{})

	if r.Method != "POST" {
		m[0] = 100
		m[1] = "Método inválido"
		response(&m, w)
		return
	}

	op := r.FormValue("op")
	code := r.FormValue("code")
	var codeArr []interface{}
	json.Unmarshal([]byte(code), &codeArr)
	hash := r.FormValue("hash")
	uId := r.FormValue("u")
	_ = uId
	if op == "" || code == "" {
		m[0] = 100
		m[1] = "Parâmetros invalidos"
		response(&m, w)
		return
	}

	db, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_users?allowNativePasswords=true")
	defer db.Close()

	if op == "1" { //Create marker

		res, err := db.Exec("INSERT INTO `udata00` (`id`, `hash`, `nc`, `ne`, `q`, `last`) VALUES (NULL, ? , 0, 0, '', CURRENT_TIMESTAMP)", hash)
		_ = res
		if err != nil {
			panic(err)
		} else {
			var id int
			db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&id)
			var mkrs string
			if codeArr[2] == "conc" {
				db.QueryRow("SELECT `mkrs_conc` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&mkrs)
			} else if codeArr[2] == "vest" {
				db.QueryRow("SELECT `mkrs_vest` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&mkrs)
			} else if codeArr[2] == "enem" {
				db.QueryRow("SELECT `mkrs_enem` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&mkrs)
			}
			var y map[string]interface{}
			json.Unmarshal([]byte(mkrs), &y)
			y[strconv.Itoa(id)] = []interface{}{codeArr[0], codeArr[1]}
			yJson, _ := json.Marshal(y)
			if codeArr[2] == "conc" {
				res, err = db.Exec("UPDATE `udataMain` SET `mkrs_conc` = ? WHERE id = ? AND hash = ?", string(yJson), uId, hash)
			} else if codeArr[2] == "vest" {
				res, err = db.Exec("UPDATE `udataMain` SET `mkrs_vest` = ? WHERE id = ? AND hash = ?", string(yJson), uId, hash)
			} else if codeArr[2] == "enem" {
				res, err = db.Exec("UPDATE `udataMain` SET `mkrs_enem` = ? WHERE id = ? AND hash = ?", string(yJson), uId, hash)
			}
			if err != nil {
				panic(err)
			} else {
				m[0] = "200"
				m[1] = y
				response(&m, w)
			}
		}
	} else if op == "2" { //Delete marker

		var mkrs string
		if codeArr[1] == "conc" {
			db.QueryRow("SELECT `mkrs_conc` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&mkrs)
		} else if codeArr[1] == "vest" {
			db.QueryRow("SELECT `mkrs_vest` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&mkrs)
		} else if codeArr[1] == "enem" {
			db.QueryRow("SELECT `mkrs_enem` FROM `udataMain` WHERE id = ? AND hash = ?", uId, hash).Scan(&mkrs)
		}

		var mkrsObject = make(map[string]interface{})
		json.Unmarshal([]byte(mkrs), &mkrsObject)

		if len(codeArr) > 0 {
			_, ok := codeArr[0].(string)
			if ok {
				delete(mkrsObject, codeArr[0].(string))

				mkrsString, _ := json.Marshal(mkrsObject)
				mkrsString_ := string(mkrsString)

				var err any

				if codeArr[1] == "conc" {
					_, err = db.Exec("UPDATE `udataMain` SET `mkrs_conc` = ? WHERE id = ? AND hash = ?", mkrsString_, uId, hash)
				} else if codeArr[1] == "vest" {
					_, err = db.Exec("UPDATE `udataMain` SET `mkrs_vest` = ? WHERE id = ? AND hash = ?", mkrsString_, uId, hash)
				} else if codeArr[1] == "enem" {
					_, err = db.Exec("UPDATE `udataMain` SET `mkrs_enem` = ? WHERE id = ? AND hash = ?", mkrsString_, uId, hash)
				}

				if err != nil {
					panic(err)
					m[0] = 100
					response(&m, w)
				} else {
					mId, _ := strconv.Atoi(codeArr[0].(string))
					_, err = db.Exec("DELETE FROM `udata00` WHERE id = ? AND hash = ?", mId, hash)
					if err != nil {
						panic(err)
						m[0] = 100
						response(&m, w)
					} else {
						m[0] = 200
						m[1] = mkrsString_
						response(&m, w)
					}
				}
			}
		}
	} else if op == "3" { //Clear marker

		mId, _ := strconv.Atoi(codeArr[0].(string))
		_, err := db.Exec("UPDATE `udata00` SET `q` = ?, `nc` = 0, `ne` = 0  WHERE id = ? AND hash = ?", "", mId, hash)
		if err != nil {
			panic(err)
			m[0] = 100
			response(&m, w)
		} else {
			m[0] = 200
			response(&m, w)
		}

	} else if op == "4" { //Read marker, full
		var quest string
		db.QueryRow("SELECT `q` FROM `udata00` WHERE id = ? AND hash = ?;", uId, hash).Scan(&quest)
		m[0] = 200
		m[1] = quest
		response(&m, w)
	} else if op == "5" { //Inster data into marker CORRECT
		var quest string
		currentTime := time.Now()
		date := currentTime.Format("2006/01/02")
		var string_insert = ""
		db.QueryRow("SELECT `q` FROM `udata00` WHERE id = ? AND hash = ?;", uId, hash).Scan(&quest)
		if !strings.Contains(quest, "#d@"+date+"#d#") {
			string_insert = "#d@" + date + "#d#"
		}
		if !strings.Contains(quest, "#"+code+"#") {
			string_insert += "c#" + code + "#"
			_, err := db.Exec("UPDATE `udata00` SET `q` = CONCAT(q, ?), `nc` = `nc` + 1, `last`=now() WHERE id = ? AND hash = ?", string_insert, uId, hash)
			if err != nil {
				panic(err)
				m[0] = 100
				response(&m, w)
			} else {
				m[0] = 200
				response(&m, w)
			}
		}
	} else if op == "6" { //Inster data into marker WRONG
		var quest string
		currentTime := time.Now()
		date := currentTime.Format("2006/01/02")
		var string_insert = ""
		db.QueryRow("SELECT `q` FROM `udata00` WHERE id = ? AND hash = ?;", uId, hash).Scan(&quest)
		if !strings.Contains(quest, "#d@"+date+"#d#") {
			string_insert = "#d@" + date + "#d#"
		}
		if !strings.Contains(quest, "#"+code+"#") {
			string_insert += "e#" + code + "#"
			_, err := db.Exec("UPDATE `udata00` SET `q` = CONCAT(q, ?), `ne` = `ne` + 1, `last`=now() WHERE id = ? AND hash = ?", string_insert, uId, hash)
			if err != nil {
				panic(err)
				m[0] = 100
				response(&m, w)
			} else {
				m[0] = 200
				response(&m, w)
			}
		}
	} else if op == "7" { //Reade created markes
		var mkrs string
		if codeArr[0] == "conc" {
			db.QueryRow("SELECT `mkrs_conc` FROM `udataMain` WHERE id = ? AND hash = ?;", uId, hash).Scan(&mkrs)
		} else if codeArr[0] == "vest" {
			db.QueryRow("SELECT `mkrs_vest` FROM `udataMain` WHERE id = ? AND hash = ?;", uId, hash).Scan(&mkrs)
		} else if codeArr[0] == "enem" {
			db.QueryRow("SELECT `mkrs_enem` FROM `udataMain` WHERE id = ? AND hash = ?;", uId, hash).Scan(&mkrs)
		}
		m[0] = 200
		m[1] = mkrs
		response(&m, w)
	} else if op == "8" { //Read marker, null
		var quest string
		db.QueryRow("SELECT `q` FROM `udata02` WHERE id = ? AND hash = ?;", uId, hash).Scan(&quest)
		m[0] = 200
		m[1] = quest
		response(&m, w)
	} else if op == "9" { //Inster data into marker null CORRECT
		var quest string
		currentTime := time.Now()
		date := currentTime.Format("2006/01/02")
		var string_insert = ""
		db.QueryRow("SELECT `q` FROM `udata02` WHERE id = ? AND hash = ?;", uId, hash).Scan(&quest)
		if !strings.Contains(quest, "#d@"+date+"#d#") {
			string_insert = "#d@" + date + "#d#"
		}
		if !strings.Contains(quest, "#"+code+"#") {
			string_insert += "c#" + code + "#"
			_, err := db.Exec("UPDATE `udata02` SET `q` = CONCAT(q, ?), `nc` = `nc` + 1, `last`=now() WHERE id = ? AND hash = ?", string_insert, uId, hash)
			if err != nil {
				panic(err)
				m[0] = 100
				response(&m, w)
			} else {
				m[0] = 200
				response(&m, w)
			}
		}
	} else if op == "10" { //Inster data into marker null WRONG
		var quest string
		currentTime := time.Now()
		date := currentTime.Format("2006/01/02")
		var string_insert = ""
		db.QueryRow("SELECT `q` FROM `udata02` WHERE id = ? AND hash = ?;", uId, hash).Scan(&quest)
		if !strings.Contains(quest, "#d@"+date+"#d#") {
			string_insert = "#d@" + date + "#d#"
		}
		if !strings.Contains(quest, "#"+code+"#") {
			string_insert += "e#" + code + "#"
			_, err := db.Exec("UPDATE `udata02` SET `q` = CONCAT(q, ?), `ne` = `ne` + 1, `last`=now() WHERE id = ? AND hash = ?", string_insert, uId, hash)
			if err != nil {
				panic(err)
				m[0] = 100
				response(&m, w)
			} else {
				m[0] = 200
				response(&m, w)
			}
		}
	}
}

func Editor(w http.ResponseWriter, r *http.Request) {
	m := make(map[int]interface{})

	if r.Method != "POST" {
		m[0] = "Método inválido"
		response(&m, w)
		return
	}

	tool := r.FormValue("tool")
	user := r.FormValue("user")
	hash := r.FormValue("hash")

	if user == "" || hash == "" {
		m[0] = "Sem credenciais"
		response(&m, w)
		return
	}

	if tool == "1" {
		avatar := r.FormValue("avatar")

		params := strings.Split(avatar, "/")
		file := params[len(params)-1]
		file_name := strings.Split(file, ".")

		db, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_users?allowNativePasswords=true")
		defer db.Close()

		res, err := db.Exec("UPDATE udbMain SET img = ? WHERE (user = ? AND hash = ?)", file_name[0], user, hash)
		_ = res
		if err != nil {
			panic(err.Error())
			m[0] = 330 //Não foi possível alterar a imagem
			response(&m, w)
			return
		} else {
			m[0] = 230
			m[1] = file_name[0]
			response(&m, w)
		}
	} else {
		m[0] = "Ferramenta invalida"
		response(&m, w)
	}
}

func filtersRead(w http.ResponseWriter, r *http.Request) {
	m := make(map[int]interface{})

	asst := ""
	bnc := ""
	org := ""
	crg := ""
	inst := ""
	db_name := ""

	db_name_, ok := r.URL.Query()["db"]
	if !ok {
		w.Write([]byte("param.err(db)"))
		return
	} else {
		db_name = db_name_[0]
	}
	asst_, ok := r.URL.Query()["asst"]
	if ok {
		asst = asst_[0]
	}
	bnc_, ok := r.URL.Query()["bnc"]
	if ok {
		bnc = bnc_[0]
	}
	org_, ok := r.URL.Query()["org"]
	if ok {
		org = org_[0]
	}
	crg_, ok := r.URL.Query()["crg"]
	if ok {
		crg = crg_[0]
	}
	inst_, ok := r.URL.Query()["crg"]
	if ok {
		inst = inst_[0]
	}

	db, _ := sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bzt_data?allowNativePasswords=true")
	defer db.Close()

	result_name := ""
	var array []int
	index := 0

	if asst != "" {
		err := json.Unmarshal([]byte(asst), &array)
		if err == nil {
			for _, i := range array {
				result_name = ""
				sql := "SELECT asst_name FROM assuntos_" + db_name + " WHERE id=" + strconv.Itoa(i) + " LIMIT 1"
				_ = db.QueryRow(sql).Scan(&result_name)
				m[index] = []string{"asst", result_name, strconv.Itoa(i)}
				index++
			}
		}
	}
	if bnc != "" {
		err_bnc := json.Unmarshal([]byte(bnc), &array)
		if err_bnc == nil {
			for _, i := range array {
				result_name = ""
				sql := "SELECT info FROM bancas_conc WHERE id=" + strconv.Itoa(i) + " LIMIT 1"
				_ = db.QueryRow(sql).Scan(&result_name)
				m[index] = []string{"bnc", result_name, strconv.Itoa(i)}
				index++
			}
		}
	}
	if org != "" {
		err_org := json.Unmarshal([]byte(org), &array)
		if err_org == nil {
			for _, i := range array {
				result_name = ""
				sql := "SELECT info FROM orgaos_conc WHERE id=" + strconv.Itoa(i) + " LIMIT 1"
				_ = db.QueryRow(sql).Scan(&result_name)
				m[index] = []string{"org", result_name, strconv.Itoa(i)}
				index++
			}
		}
	}
	if crg != "" {
		err_crg := json.Unmarshal([]byte(crg), &array)
		if err_crg == nil {
			for _, i := range array {
				result_name = ""
				sql := "SELECT crg FROM cargos_conc WHERE id=" + strconv.Itoa(i) + " LIMIT 1"
				_ = db.QueryRow(sql).Scan(&result_name)
				m[index] = []string{"crg", result_name, strconv.Itoa(i)}
				index++
			}
		}
	}
	if inst != "" {
		inst_err := json.Unmarshal([]byte(inst), &array)
		if inst_err == nil {
			for _, i := range array {
				result_name = ""
				sql := "SELECT info FROM inst_vest WHERE id=" + strconv.Itoa(i) + " LIMIT 1"
				_ = db.QueryRow(sql).Scan(&result_name)
				m[index] = []string{"inst", result_name, strconv.Itoa(i)}
				index++
			}
		}
	}

	data, _ := json.Marshal(m)
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(data))
}
