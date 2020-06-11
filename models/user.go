package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
// HashPassword ...
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost for hashing the password.
	return string(bytes), err
}

// CheckPasswordHash ...
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}

// BeforeSave hashes user password
func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedpassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = string(hashedpassword)
	return nil
}

// Prepare strips user input of any white spaces
func (u *User) Prepare() {
	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password)
}
*/
