package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var asst_main [][][]uint16
var asst_main_len int32
var atualRequestPG1 int
var atualRequestPGX int

var ad_full_qnt int
var ad_small_qnt int

// caches globais
var totalMatchesCache = make(map[string]int32)
var filteredIDsCache = make(map[string][]int32)
var filteredTagsCache = make(map[string][][]uint16)
var filteredDetsCache = make(map[string][][]uint16)

func main() {

	fmt.Println("Carregando tabela...")

	db, _ := sql.Open("mysql", "****:****:@@tcp(localhost:3306)/****:?allowNativePasswords=true")

	var bnc uint16
	var org uint16
	var nv uint16
	var ano uint16
	var tipo uint16
	var dfcd uint16
	var temp_asst []uint16
	var cargos_i []uint16
	var data [][]uint16

	var asst_string string
	var crg_string string
	var bnc_string string

	count := 0

	sql := "SELECT COUNT(*) FROM srch_conc WHERE 1"
	_ = db.QueryRow(sql).Scan(&count)

	count_string := strconv.Itoa(count)

	for i := 1; i <= 1000; i++ {

		fmt.Printf("\r%d/"+count_string, i)

		sql = "SELECT asst, crg, bnc, org, nv, ano, tipo, dfcd FROM srch_conc WHERE id=" + strconv.Itoa(i) + " LIMIT 1"
		_ = db.QueryRow(sql).Scan(&asst_string, &crg_string, &bnc_string, &org, &nv, &ano, &tipo, &dfcd)

		var assuntos_i []uint16
		var assuntos_root []uint16

		json.Unmarshal([]byte(asst_string), &temp_asst)

		for j, _ := range temp_asst {
			if temp_asst[j] < 209 {
				assuntos_root = append(assuntos_root, temp_asst[j])
				assuntos_i = append(assuntos_i, temp_asst[j])
			} else {
				assuntos_i = append(assuntos_i, temp_asst[j])
			}
		}

		json.Unmarshal([]byte(crg_string), &cargos_i)

		tmp, err := strconv.Atoi(bnc_string)
		_ = err
		bnc = uint16(tmp)

		data = append(data, []uint16{bnc, org, nv, ano, tipo, dfcd})
		data = append(data, assuntos_root)
		data = append(data, assuntos_i)
		data = append(data, cargos_i)

		asst_main = append(asst_main, data)

		data = nil
		assuntos_i = nil
		cargos_i = nil
	}

	asst_main_len = int32(len(asst_main) - 1)

	fmt.Println("Atualizando valores...")
	updateValues(nil, nil)
	fmt.Println("Tudo pronto!")

	http.HandleFunc("/v", questFetch)
	http.HandleFunc("/update1397", updateValues)
	http.HandleFunc("/console", consoleHandler)
	http.HandleFunc("/id", questByID)
	http.ListenAndServe(":3030", nil)
}

// gera chave do cache baseada nos filtros
func generateCacheKey(set int32, assuntos, cargos, banca, orgao, ano, nv, tipo, dfcd []uint16) string {
	keyParts := []string{"set=" + strconv.Itoa(int(set))}
	appendSorted := func(prefix string, arr []uint16) {
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		for _, v := range arr {
			keyParts = append(keyParts, prefix+strconv.Itoa(int(v)))
		}
	}
	appendSorted("asst_", assuntos)
	appendSorted("crg_", cargos)
	appendSorted("bnc_", banca)
	appendSorted("org_", orgao)
	appendSorted("ano_", ano)
	appendSorted("nv_", nv)
	appendSorted("tp_", tipo)
	appendSorted("dfcd_", dfcd)

	rawKey := ""
	for _, k := range keyParts {
		rawKey += k + "|"
	}
	h := sha1.New()
	h.Write([]byte(rawKey))
	return hex.EncodeToString(h.Sum(nil))
}

// função principal
func questFetch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if atualRequestPG1 > 150 || atualRequestPGX > 300 {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("0x503"))
		r.Body.Close()
		return
	}

	// ----------- parâmetros básicos -----------
	var pg, set int32
	id_, ok := r.URL.Query()["pg"]
	if !ok || len(id_[0]) < 1 {
		w.Write([]byte("param.err(pg)"))
		r.Body.Close()
		return
	}
	n, err := strconv.Atoi(id_[0])
	if err != nil {
		w.Write([]byte("param.err(pg)"))
		r.Body.Close()
		return
	}
	pg = int32(n)

	set_, hasSet := r.URL.Query()["set"]
	if hasSet {
		if _, err := strconv.Atoi(set_[0]); err == nil {
			set = 1
		}
	}

	// ----------- filtros -----------
	getUint16Slice := func(key string) ([]uint16, bool) {
		if len(r.URL.Query()[key]) > 0 {
			var arr []uint16
			_ = json.Unmarshal([]byte(r.URL.Query()[key][0]), &arr)
			return arr, true
		}
		return nil, false
	}

	fill_assuntos, _ := getUint16Slice("asst")
	fill_assuntos_roots, _ := getUint16Slice("asst_rt")
	fill_cargos, _ := getUint16Slice("crg")
	fill_banca, _ := getUint16Slice("bnc")
	fill_org, _ := getUint16Slice("org")
	fill_ano, _ := getUint16Slice("ano")
	fill_nv, _ := getUint16Slice("nv")
	fill_tipo, _ := getUint16Slice("tp")
	fill_dfcd, _ := getUint16Slice("dfcd")

	key := generateCacheKey(set, fill_assuntos, fill_cargos, fill_banca, fill_org, fill_ano, fill_nv, fill_tipo, fill_dfcd)

	var ids []int32
	var tags [][]uint16
	var dets [][]uint16
	var totalMatches int32

	cachedIDs, ok := filteredIDsCache[key]
	if ok {
		totalMatches = totalMatchesCache[key]
		// pegar IDs da página diretamente
		start := (pg - 1) * 10
		end := start + 10
		if int(start) >= len(cachedIDs) {
			ids = []int32{}
		} else {
			if int(end) > len(cachedIDs) {
				end = int32(len(cachedIDs))
			}
			ids = cachedIDs[start:end]
			tags = filteredTagsCache[key][start:end]
			dets = filteredDetsCache[key][start:end]
		}
	} else {
		// percorrer asst_main e aplicar filtros
		for sizeOut := asst_main_len; sizeOut >= 0; sizeOut-- {
			checkAssuntos := len(fill_assuntos) > 0 && len(fill_assuntos_roots) > 0
			checkCargos := len(fill_cargos) > 0
			checkBanca := len(fill_banca) > 0
			checkOrgao := len(fill_org) > 0
			checkAno := len(fill_ano) > 0
			checkNv := len(fill_nv) > 0
			checkTipo := len(fill_tipo) > 0
			checkDfcd := len(fill_dfcd) > 0

			passed := true

			// -------- filtros --------
			if checkAssuntos {
				// 1️⃣ Checa roots
				hasRoot := false
				for _, root := range fill_assuntos_roots {
					for _, j := range asst_main[sizeOut][1] {
						if j == root {
							hasRoot = true
							break
						}
					}
					if hasRoot {
						break
					}
				}
				if !hasRoot {
					passed = false
				} else {
					// 2️⃣ Checa assuntos não-root (>=250)
					nonRootAssuntos := []uint16{}
					for _, fa := range fill_assuntos {
						if fa >= 250 {
							nonRootAssuntos = append(nonRootAssuntos, fa)
						}
					}

					if len(nonRootAssuntos) > 0 {
						foundAny := false
						for _, fa := range nonRootAssuntos {
							for _, v := range asst_main[sizeOut][2] {
								if fa == v {
									foundAny = true
									break
								}
							}
							if foundAny {
								break
							}
						}
						if !foundAny {
							passed = false
						}
					}
					// Se não houver assuntos não-root, a presença dos roots já garante passed = true
				}
			}

			if passed && checkCargos {
				found := false
				for _, fc := range fill_cargos {
					for _, v := range asst_main[sizeOut][3] {
						if fc == v {
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if !found {
					passed = false
				}
			}

			if passed && checkBanca {
				ok := false
				for _, fb := range fill_banca {
					if fb == asst_main[sizeOut][0][0] {
						ok = true
						break
					}
				}
				if !ok {
					passed = false
				}
			}

			if passed && checkOrgao {
				ok := false
				for _, fo := range fill_org {
					if fo == asst_main[sizeOut][0][1] {
						ok = true
						break
					}
				}
				if !ok {
					passed = false
				}
			}

			if passed && checkNv {
				ok := false
				for _, fn := range fill_nv {
					if fn == asst_main[sizeOut][0][2] {
						ok = true
						break
					}
				}
				if !ok {
					passed = false
				}
			}

			if passed && checkAno {
				anoQuestao := asst_main[sizeOut][0][3]
				found := false
				for _, ano := range fill_ano {
					if anoQuestao == ano {
						found = true
						break
					}
				}
				if !found {
					passed = false
				}
			}

			if passed && checkTipo {
				ok := false
				for _, ft := range fill_tipo {
					if ft == asst_main[sizeOut][0][4] {
						ok = true
						break
					}
				}
				if !ok {
					passed = false
				}
			}

			if passed && checkDfcd {
				ok := false
				for _, fd := range fill_dfcd {
					if fd == asst_main[sizeOut][0][5] {
						ok = true
						break
					}
				}
				if !ok {
					passed = false
				}
			}

			// -------- resultado --------
			if passed {
				totalMatches++
				cachedIDs = append(cachedIDs, int32(sizeOut+1))
				filteredTagsCache[key] = append(filteredTagsCache[key], asst_main[sizeOut][2])
				filteredDetsCache[key] = append(filteredDetsCache[key], asst_main[sizeOut][0])
			}
		}

		// salvar no cache
		totalMatchesCache[key] = totalMatches
		filteredIDsCache[key] = cachedIDs

		// pegar IDs da página
		start := (pg - 1) * 10
		end := start + 10
		if int(start) >= len(cachedIDs) {
			ids = []int32{}
		} else {
			if int(end) > len(cachedIDs) {
				end = int32(len(cachedIDs))
			}
			ids = cachedIDs[start:end]
			tags = filteredTagsCache[key][start:end]
			dets = filteredDetsCache[key][start:end]
		}
	}

	// ----------- acessar banco -----------
	var m = make(map[int]interface{})
	if len(ids) > 0 {
		db, _ := sql.Open("mysql", "****:****:@@tcp(localhost:3306)/****:?allowNativePasswords=true")
		defer db.Close()

		sqlQ := "SELECT id, cont, c_itens, asst, crg, bnc, org, qcmt FROM cont_conc WHERE id = " + strconv.Itoa(int(ids[0]))
		for j := 1; j < len(ids); j++ {
			sqlQ += " OR id = " + strconv.Itoa(int(ids[j]))
		}
		sqlQ += " ORDER BY id DESC LIMIT 10"

		rows, err := db.Query(sqlQ)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for i := 0; rows.Next(); i++ {
			var id int32
			var cont, c_itens, ass, crg, bnc, org string
			var qcmt int
			rows.Scan(&id, &cont, &c_itens, &ass, &crg, &bnc, &org, &qcmt)

			x := make(map[int]interface{})
			x[0] = id
			x[1] = tags[i]
			x[2] = dets[i]
			x[3] = ass
			x[4] = crg
			x[5] = bnc
			x[6] = org
			x[7] = cont
			x[8] = c_itens
			x[9] = qcmt
			m[i] = x
		}
	}

	// último item do JSON: totalMatches
	m[len(m)] = totalMatches

	data, _ := json.Marshal(m)
	w.Header().Set("Content-type", "application/json")
	w.Write(data)

	if pg == 1 {
		atualRequestPG1--
	} else {
		atualRequestPGX--
	}
	r.Body.Close()
}

func questByID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// ---------- validação do parâmetro ----------
	idParam, ok := r.URL.Query()["idQst"]
	if !ok || len(idParam[0]) < 1 {
		w.Write([]byte("param.err(id)"))
		r.Body.Close()
		return
	}

	id, err := strconv.Atoi(idParam[0])
	if err != nil {
		w.Write([]byte("param.err(id)"))
		r.Body.Close()
		return
	}

	// ---------- acessar banco ----------
	db, err := sql.Open("mysql", "****:****@@tcp(localhost:3306)/****?allowNativePasswords=true")
	if err != nil {
		http.Error(w, "db.conn.error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	sqlQ := "SELECT id, cont, c_itens, asst, crg, bnc, org, qcmt FROM cont_conc WHERE id = ? LIMIT 1"
	row := db.QueryRow(sqlQ, id)

	var (
		qID     int32
		cont    string
		c_itens string
		ass     string
		crg     string
		bnc     string
		org     string
		qcmt    int
	)

	m := make(map[int]interface{})
	err = row.Scan(&qID, &cont, &c_itens, &ass, &crg, &bnc, &org, &qcmt)
	if err == sql.ErrNoRows {
		m[len(m)] = 0
		data, _ := json.Marshal(m)
		w.Header().Set("Content-type", "application/json")
		w.Write(data)
		return
	} else if err != nil {
		http.Error(w, "db.query.error", http.StatusInternalServerError)
		return
	}

	// ---------- montar resposta no mesmo formato ----------
	x := make(map[int]interface{})
	x[0] = qID
	x[1] = []uint16{} // sem tags (não vem de cache, mas mantém a posição)
	x[2] = []uint16{} // sem dets (mesma ideia)
	x[3] = ass
	x[4] = crg
	x[5] = bnc
	x[6] = org
	x[7] = cont
	x[8] = c_itens
	x[9] = qcmt
	m[0] = x

	// último item = totalMatches (sempre 1 nesse caso)
	m[1] = 1

	data, _ := json.Marshal(m)
	w.Header().Set("Content-type", "application/json")
	w.Write(data)

	r.Body.Close()
}

func updateValues(w http.ResponseWriter, r *http.Request) {
	db_ads, _ := sql.Open("mysql", "****:****@@tcp(localhost:3306)/****?allowNativePasswords=true")

	db_ads.QueryRow("SELECT COUNT(*) FROM small WHERE 1").Scan(&ad_small_qnt)
	db_ads.QueryRow("SELECT COUNT(*) FROM full WHERE 1").Scan(&ad_full_qnt)
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica parâmetro de segurança
	pass := r.URL.Query().Get("pass")
	if pass != "phSS1397" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("access denied"))
		return
	}

	// Opcional: pode filtrar quais caches mostrar com parâmetros extras
	show := r.URL.Query().Get("show") // ex: show=totalMatches,ids

	response := make(map[string]interface{})

	// Mostra totalMatchesCache
	if show == "" || show == "totalMatches" {
		response["totalMatchesCache_size"] = len(totalMatchesCache)
		response["totalMatchesCache_keys"] = func() []string {
			keys := []string{}
			count := 0
			for k := range totalMatchesCache {
				keys = append(keys, k)
				count++
				if count >= 10 { // limita a 10 chaves para não poluir
					break
				}
			}
			return keys
		}()
	}

	// Mostra filteredIDsCache
	if show == "" || show == "ids" {
		response["filteredIDsCache_size"] = len(filteredIDsCache)
		response["filteredIDsCache_keys"] = func() []string {
			keys := []string{}
			count := 0
			for k := range filteredIDsCache {
				keys = append(keys, k)
				count++
				if count >= 10 {
					break
				}
			}
			return keys
		}()
	}

	// Mostra filteredTagsCache e filteredDetsCache
	if show == "" || show == "tags_dets" {
		response["filteredTagsCache_size"] = len(filteredTagsCache)
		response["filteredDetsCache_size"] = len(filteredDetsCache)
	}

	data, _ := json.MarshalIndent(response, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
