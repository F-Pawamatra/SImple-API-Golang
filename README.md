# Initialization
Run 
```bash
go mod latihanJwt
```

make sure that you has a databases, and a table. here is mine
```go
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
```

connect your database by editing dsn in file handlers.go line 133 
```go
dsn := "host=localhost user=postgres password=postgre dbname=db_example port=25432 sslmode=disable TimeZone=Asia/Jakarta"
```

then, run 
```bash
go build .
```

after that, an executable file called latihanJwt.exe will be created. run it
