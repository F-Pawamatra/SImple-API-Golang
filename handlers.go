package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kylelemons/godebug/pretty"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var jwtKey = []byte("secret_key")

var users = map[string]string{
	"xxx": "3327xxx001",
	"xxx002": "3327xxx002",
}

type Credentials struct {
	NOPPBB 	string `json:"noppbb"`
	NIK 	string `json:"nik"`
}

type Claims struct {
	NOPPBB string `json:"noppbb"`
	jwt.StandardClaims
}

type Data struct {
	NOPPBB    string    `json:"noppbb"`
	Nama      string    `json:"nama"`
	Alamat    string    `json:"alamat"`
	Kabupaten string    `json:"kabupaten"`
	Kecamatan string    `json:"kecamatan"`
	Desa      string    `json:"desa"`
	Rt        string    `json:"rt"`
	Rw        string    `json:"rw"`
	DataPbbs  []DataPbb `gorm:"foreignKey:IdNOPBB" json:"datapbbs"`
}

type DataPbb struct {
	Tahun    string `json:"tahun"`
	Pajak    int64  `json:"pajak"`
	Denda    int64  `json:"denda"`
	IdNOPPBB string `json:"id_noppbb"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedNIK, ok := users[credentials.NOPPBB]

	if !ok || expectedNIK != credentials.NIK {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		NOPPBB: credentials.NOPPBB,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

	w.Write([]byte(fmt.Sprintf(tokenString)))
}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.NOPPBB)))
	GetData(claims.NOPPBB, w)
}

func GetData(NoppbbLogin string, w http.ResponseWriter,) {
	dsn := "host=localhost user=postgres password=postgre dbname=db_example port=25432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	rows, err := db.Raw(`select * from "data" d inner join data_pbb dp on dp.id_noppbb = d.noppbb where d.noppbb = ?`, NoppbbLogin).Rows()
	
	if err != nil {
		log.Panic(err)
	}

	defer rows.Close()
	// Values to load into
	newData := &Data{}
	newData.DataPbbs = make([]DataPbb, 0)

	for rows.Next() {
		dataPbb := DataPbb{}
		err = rows.Scan(&newData.NOPPBB, &newData.Nama, &newData.Alamat, &newData.Kabupaten, &newData.Kecamatan, &newData.Desa, &newData.Rt, &newData.Rw, &dataPbb.Tahun, &dataPbb.Pajak, &dataPbb.Denda, &dataPbb.IdNOPPBB)
		if err != nil {
			log.Panic(err)
		}
		newData.DataPbbs = append(newData.DataPbbs, dataPbb)
	}
	log.Print(pretty.Sprint(newData))
}

