package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var rdb *redis.Client
var ctx = context.Background()

// ----------- MAIN -----------
func main() {
	var err error
	// Conexão com MySQL
	db, err = sql.Open("mysql", "pedrosilva:phSS1397@@tcp(localhost:3306)/bizout_users?allowNativePasswords=true")
	if err != nil {
		log.Fatal(err)
	}

	// Conexão com Redis
	initRedis()

	// Rotas
	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/me", meHandler)
	mux.HandleFunc("/userimage", profileAvatarHandler)
	mux.HandleFunc("/changepassword", changePasswordHandler)
	mux.HandleFunc("/devices", devicesHandler)

	mux.HandleFunc("/devices/logout/", logoutDeviceHandler) // <-- Observe a barra no final

	// Middleware de CORS
	handler := corsMiddleware(mux)

	fmt.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
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

// ----------- STRUCTS -----------

type UserRequest struct {
	User  string `json:"user"`
	Mail  string `json:"mail"`
	Pass  string `json:"pass"`
	Fname string `json:"fname"`
	Cies  string `json:"cies,omitempty"`
	Dnas  string `json:"dnas"` // <-- aceita nulo
}

type LoginRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MeResponse struct {
	ID    int    `json:"id"`
	User  string `json:"user"`
	Fname string `json:"fname"`
	Image string `json:"img"`
}

type MeResponseFull struct {
	ID    int    `json:"id"`
	User  string `json:"user"`
	Fname string `json:"fname"`
	Image string `json:"img"`
	Mail  string `json:"mail"`
	Cies  string `json:"cies"`
	Dnas  string `json:"dnas"`
}

type SessionData struct {
	SessionID string `json:"session_id"` // chave no Redis
	UserID    int    `json:"user_id"`
	Device    string `json:"device"`
	IP        string `json:"ip"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	IsCurrent bool   `json:"is_current"`
}

// ----------- UTILS -----------
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": message})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allowedOrigins := map[string]bool{
		// 	"http://localhost:5173": true,
		// 	"http://localhost:5174": true,
		// 	"https://bizout.com.br": true,
		// }

		// origin := r.Header.Get("Origin")
		// if allowedOrigins[origin] {
		// 	w.Header().Set("Access-Control-Allow-Origin", origin)
		// }

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ----------- REGISTER -----------
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var u UserRequest
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if u.User == "" || u.Mail == "" || u.Pass == "" || u.Fname == "" {
		writeJSONError(w, http.StatusBadRequest, "Campos obrigatórios faltando")
		return
	}

	// Verificar se usuário ou email já existem
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM udbMain WHERE user=? OR mail=?", u.User, u.Mail).Scan(&exists)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao verificar usuário")
		return
	}
	if exists > 0 {
		writeJSONError(w, http.StatusConflict, "Usuário ou email já existe")
		return
	}

	// Hash da senha
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao gerar hash da senha")
		return
	}

	// Inserir no banco
	stmt, err := db.Prepare(`INSERT INTO udbMain (user, mail, pass, fname, cies, dnas) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao preparar query")
		return
	}
	defer stmt.Close()

	var birth sql.NullTime
	if u.Dnas != "" {
		t, err := time.Parse("2006-01-02", u.Dnas)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Erro ao inserir usuário: "+err.Error())
			return
		}
		birth = sql.NullTime{Time: t, Valid: true}
	} else {
		birth = sql.NullTime{Valid: false}
	}

	_, err = stmt.Exec(u.User, u.Mail, string(hashedPass), u.Fname, u.Cies, birth)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao inserir usuário: "+err.Error())
		return
	}

	// Sucesso
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Usuário criado com sucesso"})
}

// ----------- LOGIN -----------
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.User == "" || req.Pass == "" {
		writeJSONError(w, http.StatusBadRequest, "Dados inválidos")
		return
	}

	// Buscar usuário
	var userID int
	var hashedPass string
	err = db.QueryRow("SELECT id, pass FROM udbMain WHERE user=? OR mail=?", req.User, req.User).Scan(&userID, &hashedPass)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusUnauthorized, "Usuário ou senha incorretos")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao buscar usuário")
		return
	}

	// Comparar senha
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Pass))
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Usuário ou senha incorretos")
		return
	}

	// Criar sessão
	sessionID := uuid.NewString()

	// Descobre IP do cliente
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Detectar user-agent (dispositivo/navegador)
	userAgent := r.Header.Get("User-Agent")
	device := "Desconhecido"
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "mobile") {
		device = "Celular"
	} else if strings.Contains(ua, "windows") {
		device = "Windows"
	} else if strings.Contains(ua, "mac") {
		device = "MacOS"
	} else if strings.Contains(ua, "linux") {
		device = "Linux"
	}

	// Struct da sessão
	sessionData := struct {
		UserID    int    `json:"user_id"`
		Device    string `json:"device"`
		IP        string `json:"ip"`
		Location  string `json:"location"`
		CreatedAt string `json:"created_at"`
	}{
		UserID:    userID,
		Device:    device,
		IP:        ip,
		Location:  "Desconhecido", // dá pra melhorar com API de geolocalização
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Serializa para JSON
	jsonData, _ := json.Marshal(sessionData)

	// Salva no Redis (24h)
	err = rdb.Set(ctx, "session:"+sessionID, jsonData, 24*time.Hour).Err()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao criar sessão")
		return
	}

	// Retorna cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400,
	})

	// Resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Status:  "ok",
		Message: "Login realizado com sucesso",
	})
}

// ----------- ME -----------
func meHandler(w http.ResponseWriter, r *http.Request) {
	// Pegar cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão não encontrada")
		return
	}

	op := r.URL.Query().Get("op")

	// Buscar sessão no Redis
	sessionData, err := rdb.Get(ctx, "session:"+cookie.Value).Result()
	if err == redis.Nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão inválida ou expirada")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao acessar Redis")
		return
	}

	// Decodificar JSON da sessão
	var session struct {
		UserID int `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao decodificar sessão")
		return
	}

	// Buscar dados do usuário no MySQL
	var user, fname, img, mail, cies, dnas string
	if op != "full" {
		err = db.QueryRow("SELECT user, fname, img FROM udbMain WHERE id=?", session.UserID).Scan(&user, &fname, &img)
	} else {
		err = db.QueryRow("SELECT user, fname, img, mail, cies, dnas FROM udbMain WHERE id=?", session.UserID).Scan(&user, &fname, &img, &mail, &cies, &dnas)
	}

	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusUnauthorized, "Usuário não encontrado")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao buscar usuário")
		return
	}

	// Resposta JSON
	w.Header().Set("Content-Type", "application/json")
	if op != "full" {
		json.NewEncoder(w).Encode(MeResponse{
			ID:    session.UserID,
			User:  user,
			Fname: fname,
			Image: img,
		})
	} else {
		json.NewEncoder(w).Encode(MeResponseFull{
			ID:    session.UserID,
			User:  user,
			Fname: fname,
			Image: img,
			Mail:  mail,
			Cies:  cies,
			Dnas:  dnas,
		})
	}
}

func profileAvatarHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar sessão
	cookie, err := r.Cookie("session_id")
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão não encontrada")
		return
	}

	val, err := rdb.Get(ctx, "session:"+cookie.Value).Result()
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão inválida")
		return
	}

	var sessionData struct {
		UserID int `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(val), &sessionData); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao decodificar sessão")
		return
	}

	userID := sessionData.UserID

	// Parse multipart
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Erro ao processar arquivo")
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		writeJSONError(w, http.StatusBadRequest, "Erro ao receber arquivo")
		return
	}

	var avatarData string
	if file != nil {
		// Caso arquivo enviado
		defer file.Close()
		filename := fmt.Sprintf("/usr/share/nginx/guia/users/uploads/avatars/%d_%s", userID, handler.Filename)
		dst, err := os.Create(filename)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Erro ao salvar arquivo")
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Erro ao salvar arquivo")
			return
		}

		// Remove parte do caminho desnecessária
		filename = strings.TrimPrefix(filename, "/usr/share/nginx/guia")

		// Prepara JSON
		avatarDataJSON := map[string]interface{}{
			"image": true,
			"data":  filename,
		}
		jsonBytes, _ := json.Marshal(avatarDataJSON)
		avatarData = string(jsonBytes)

	} else {
		// Caso sem arquivo, pega cor enviada do front
		avatarField := r.FormValue("avatar") // será algo como {"image":false,"data":"#f87171"}
		if avatarField == "" {
			writeJSONError(w, http.StatusBadRequest, "Nenhum arquivo ou cor enviada")
			return
		}

		// Valida JSON
		var tmp map[string]interface{}
		err := json.Unmarshal([]byte(avatarField), &tmp)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "JSON inválido para avatar")
			return
		}

		avatarData = avatarField
	}

	// Salva no banco
	_, err = db.Exec("UPDATE udbMain SET img=? WHERE id=?", avatarData, userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao atualizar avatar no banco")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Avatar salvo"})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No session found", http.StatusUnauthorized)
		return
	}

	// Remove a sessão do Redis
	err = rdb.Del(ctx, "session:"+cookie.Value).Err()
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Remove o cookie do cliente
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expira imediatamente
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","message":"Logout realizado com sucesso"}`))
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	// Verifica sessão
	cookie, err := r.Cookie("session_id")
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão não encontrada")
		return
	}

	val, err := rdb.Get(ctx, "session:"+cookie.Value).Result()
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão inválida")
		return
	}

	var sessionData struct {
		UserID int `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(val), &sessionData); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao decodificar sessão")
		return
	}

	userID := sessionData.UserID

	// Estrutura para JSON recebido
	var req struct {
		Current string `json:"current"`
		New     string `json:"new"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if req.Current == "" || req.New == "" {
		writeJSONError(w, http.StatusBadRequest, "Campos obrigatórios faltando")
		return
	}

	// Busca senha atual
	var hashedPass string
	err = db.QueryRow("SELECT pass FROM udbMain WHERE id=?", userID).Scan(&hashedPass)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusUnauthorized, "Usuário não encontrado")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao buscar senha")
		return
	}

	// Confere senha antiga
	if bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Current)) != nil {
		writeJSONError(w, http.StatusUnauthorized, "Senha atual incorreta")
		return
	}

	// Gera hash da nova senha
	newHashed, err := bcrypt.GenerateFromPassword([]byte(req.New), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao gerar hash da nova senha")
		return
	}

	// Atualiza no banco
	_, err = db.Exec("UPDATE udbMain SET pass=? WHERE id=?", string(newHashed), userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao atualizar senha")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Senha alterada com sucesso"})
}

// Função utilitária para buscar todas as sessões ativas de um usuário
func getUserSessions(userID int) ([]SessionData, error) {
	var sessions []SessionData

	iter := rdb.Scan(ctx, 0, "session:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		// Pega os dados da sessão
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		// Desserializa
		var data SessionData
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			continue
		}

		// Filtra pelo userID
		if data.UserID == userID {
			data.SessionID = key // a chave "session:<id>" do Redis
			sessions = append(sessions, data)
		}
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// Endpoint para listar os dispositivos conectados do usuário logado
func devicesHandler(w http.ResponseWriter, r *http.Request) {
	// Pegar cookie da sessão atual
	cookie, err := r.Cookie("session_id")
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão não encontrada")
		return
	}

	// Remove o prefixo "session:" caso exista
	currentSessionID := strings.TrimPrefix(cookie.Value, "session:")

	// Buscar dados da sessão atual
	val, err := rdb.Get(ctx, "session:"+currentSessionID).Result()
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão inválida ou expirada")
		return
	}

	// Desserializa JSON da sessão atual
	var currentSession SessionData
	if err := json.Unmarshal([]byte(val), &currentSession); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao ler sessão atual")
		return
	}

	currentUserID := currentSession.UserID

	// Buscar todas as sessões do usuário
	allSessions, err := getUserSessions(currentUserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao buscar sessões")
		return
	}

	// Marca qual é a sessão atual
	for i := range allSessions {
		allSessions[i].IsCurrent = strings.TrimPrefix(allSessions[i].SessionID, "session:") == currentSessionID
	}

	// Retorna JSON com todas as sessões
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"sessions": allSessions,
	})
}

func logoutDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	// Pegar cookie da sessão atual
	cookie, err := r.Cookie("session_id")
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão não encontrada")
		return
	}

	// Descobre userID da sessão atual
	val, err := rdb.Get(ctx, "session:"+cookie.Value).Result()
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão inválida ou expirada")
		return
	}

	// Desserializa o JSON
	var sessionData SessionData
	if err := json.Unmarshal([]byte(val), &sessionData); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao ler sessão")
		return
	}

	currentUserID := sessionData.UserID

	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Sessão inválida ou expirada")
		return
	}

	// Pega sessionId da URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeJSONError(w, http.StatusBadRequest, "Session ID não informado")
		return
	}
	targetSessionID := parts[3]

	// Verifica se a sessão alvo pertence ao usuário
	val, err = rdb.Get(ctx, targetSessionID).Result()
	if err == redis.Nil {
		writeJSONError(w, http.StatusNotFound, "Sessão não encontrada")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao acessar Redis")
		return
	}

	// Parse do valor JSON da sessão
	var sessData struct {
		UserID int `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(val), &sessData); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao decodificar sessão")
		return
	}

	if sessData.UserID != currentUserID {
		writeJSONError(w, http.StatusForbidden, "Sessão não pertence ao usuário")
		return
	}

	// Remove sessão do Redis
	if err := rdb.Del(ctx, targetSessionID).Err(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Erro ao remover sessão")
		return
	}

	// Sucesso
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Sessão encerrada com sucesso",
	})
}
