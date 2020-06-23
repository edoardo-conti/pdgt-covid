package models

// User struttura per la gestione di utenti
// nota: il campo 'avatar_url' non viene salvato in database perchè non necessario ai fini della rappresentazione
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	Avatar   string `json:"avatar_url"`
}
