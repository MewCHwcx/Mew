package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

type Warehouse struct {
	Id      string `json:"id"`
	Details string `json:"details"`
}

type Rice_transfer_record struct {
	Id             int    `json:"id"`
	Upd_date       string `json:"upd_date"` //DATE
	RiceName       string `json:"rice_name"`
	Amount         int    `json:"amount"` //DOUBLE
	Status         string `json:"status"`
	From_warehouse string `json:"from_warehouse"`
	To_warehouse   string `json:"to_warehouse"`
	Rice_id_from   string `json:"rice_id_from"`
	Rice_id_to     string `json:"rice_id_to"`
}

type Rice_info struct {
	Rice_info_id int    `json:"rice_info_id"`
	Rice_id      string `json:"rice_id"`
	Rice_name    string `json:"rice_name"`
	Rice_unit    string `json:"rice_unit"`
	Rice_on_hand string `json:"rice_on_hand"`
	Upd_date     string `json:"upd_date"`
	Shortname    string `json:"shortname"`
}

type TransferRice struct {
	Id             int
	Upd_date       string `json:"upd_date"` //DATE
	RiceName       string `json:"rice_name"`
	Amount         int    `json:"amount"` //DOUBLE
	Status         string `json:"status"`
	From_warehouse string `json:"from_warehouse"`
	To_warehouse   string `json:"to_warehouse"`
}

type RiceResult struct {
	Update_date string `json:"update_date"`
	RiceName    string `json:"rice_name"`
	Weight      int    `json:"weight"`
	Well        int    `json:"well"`
	Lose        int    `json:"lose"`
}

type Rice_packing_daily struct {
	Id           int    `json:"id"`
	Warehouse_id string `json:"warehouse_id"`
	Rice_id      string `json:"rice_id"`
	Amount       int    `json:"amount"`
	Upd_date     string `json:"upd_date"`
	Activity     string `json:"activity"`
}

type Rice_sale struct {
	Id           int    `json:"id"`
	Warehouse_id string `json:"warehouse"`
	Rice_id      string `json:"rice_id"`
	Amount       int    `json:"amount"`
	Upd_date     string `json:"upd_date"`
}

func NewWarehouse() *Warehouse {
	return new(Warehouse)
}

func (data *Warehouse) GetAllWarehouses() ([]Warehouse, error) {
	warehouses := make([]Warehouse, 0)
	db, err := GetDB()
	if err != nil {
		return warehouses, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id, details FROM warehouse`)
	if err != nil {
		return warehouses, err
	}
	for rows.Next() {
		warehouse := Warehouse{}
		rows.Scan(&warehouse.Id, &warehouse.Details)
		warehouses = append(warehouses, warehouse)
	}
	return warehouses, nil
}

func (data *Warehouse) FindByID(ID string) (Warehouse, error) {
	if ID == "" {
		return Warehouse{}, errors.New("ID can not be empty")
	}
	db, err := GetDB()
	warehouse := Warehouse{}
	if err != nil {
		return warehouse, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id, details FROM warehouse WHERE id = ` + ID)
	if err != nil {
		return warehouse, err
	}
	if rows.Next() {
		rows.Scan(&warehouse.Id, &warehouse.Details)
		return warehouse, nil
	}
	return Warehouse{}, err
}

func (data *Warehouse) Update(newWarehouse Warehouse) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE warehouse SET id = ?, details = ? WHERE id = ?`, newWarehouse.Id, newWarehouse.Details, newWarehouse.Id)
	if err != nil {
		return err
	}
	return nil
}

func Index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	_, err := w.Write([]byte("Welcome!"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// FindAllWarehousesreturn all warehouses in database
func GetAllWarehouses(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	users, err := NewWarehouse().GetAllWarehouses()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (data *Warehouse) SaveWarehouse() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO warehouse (id, details) VALUES (?, ?)`, data.Id, data.Details)
	if err != nil {
		return err
	}
	return nil
}

//create a warehouses
func CreateWarehouse(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	body := req.Body
	warehouse := NewWarehouse()
	err := json.NewDecoder(body).Decode(warehouse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer body.Close()
	err = warehouse.SaveWarehouse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// FindWarehouseByID return a warehouse
func FindWarehouseByID(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	warehouse, err := NewWarehouse().FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(warehouse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateWarehouse update an existing warehouse
func UpdateWarehouse(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ID := params.ByName("id")
	var warehouse Warehouse
	err := json.NewDecoder(req.Body).Decode(&warehouse)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if warehouse.Id == "" {
		http.Error(w, "Please set ID in warehouse information", http.StatusBadRequest)
		return
	}
	_, err = NewWarehouse().FindByID(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = NewWarehouse().Update(warehouse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteWarehouse delete a warehouse from database
func (data *Warehouse) Delete(ID string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM warehouse WHERE id = ?`, data.Id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteWarehouse(w http.ResponseWriter, req *http.Request, prm httprouter.Params) {
	ID := prm.ByName("id")
	warehouse, err := NewWarehouse().FindByID(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = warehouse.Delete(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewTransfer() *TransferRice {
	return new(TransferRice)
}

// FindTransferRiceByID return a transferRice
func FindTransferRiceByID(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	transferRice, err := NewTransfer().FindTransferRiceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(transferRice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (data *TransferRice) FindTransferRiceByID(ID string) ([]TransferRice, error) {
	transferRices := make([]TransferRice, 0)
	db, err := GetDB()
	if err != nil {
		return transferRices, err
	}
	if ID == "" {
		return transferRices, errors.New("ID can not be empty")
	}

	defer db.Close()
	rows, err := db.Query(`SELECT ID, Update_date, RiceName, Amount, Status, From_warehouse_name, To_warehouse_name 
	FROM (SELECT rf.id AS ID, rf.upd_date AS Update_date, ri.rice_name AS RiceName, rf.amount AS Amount, "โยกย้าย" AS Status, w2.details AS From_warehouse_name, w1.details AS To_warehouse_name 
	FROM rice_transfer_record rf 
	JOIN warehouse w1 ON w1.id = rf.to_warehouse 
	JOIN warehouse w2 ON w2.id = rf.from_warehouse 
	JOIN rice_info ri ON rf.rice_id_from = ri.rice_id 
	JOIN rice_info ri2 ON rf.rice_id_to = ri2.rice_id WHERE w1.id = ` + ID + ` OR w2.id = ` + ID +
		`UNION
	SELECT rp.id AS ID ,rp.upd_date AS Update_date, ri.rice_name AS RiceName, rp.amount AS Amount, "การผลิต" AS Status, "-" AS From_warehouse_name, "-" AS To_warehouse_name 
	FROM rice_packing_daily rp 
	JOIN rice_info ri ON rp.rice_id = ri.rice_id WHERE rp.warehouse_id = ` + ID + ` AND activity = "source"
	UNION
	SELECT rs.id AS ID ,rs.upd_date AS Update_date, ri.rice_name AS RiceName, rs.amount AS Amount, "การขายออก" AS Status, "-" AS From_warehouse_name, "-" AS To_warehouse_name 
	FROM rice_sale rs 
	JOIN rice_info ri ON rs.rice_id = ri.rice_id WHERE rs.warehouse_id = ` + ID + `) AS result ORDER BY Update_date DESC LIMIT 10`)
	if err != nil {
		return transferRices, err
	}
	for rows.Next() {
		transferRice := TransferRice{}
		rows.Scan(&transferRice.Id, &transferRice.Upd_date, &transferRice.RiceName, &transferRice.Amount, &transferRice.Status, &transferRice.From_warehouse, &transferRice.To_warehouse)
		transferRices = append(transferRices, transferRice)
	}
	return transferRices, err
}

func NewTransferRice() *Rice_transfer_record {
	return new(Rice_transfer_record)
}

func (data *Rice_transfer_record) SaveRiceTransfer() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO rice_transfer_record VALUES (NULL,SYSDATE(),?,?,?,?,?)`, data.From_warehouse, data.To_warehouse, data.Amount, data.Rice_id_from, data.Rice_id_to)
	if err != nil {
		return err
	}
	return nil
}

func CreateTransfer(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	body := req.Body
	transferRice := NewTransferRice()
	err := json.NewDecoder(body).Decode(transferRice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer body.Close()
	err = transferRice.SaveRiceTransfer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewPackingDaily() *Rice_packing_daily {
	return new(Rice_packing_daily)
}

func (data *Rice_packing_daily) SaveRicePackingDaily() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO rice_packing_daily VALUES (NULL,?,?,?,SYSDATE(),?)`, data.Warehouse_id, data.Rice_id, data.Amount, data.Activity)
	if err != nil {
		return err
	}
	return nil
}

func CreatePacking(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	body := req.Body
	packing := NewPackingDaily()
	err := json.NewDecoder(body).Decode(packing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer body.Close()
	err = packing.SaveRicePackingDaily()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewRiceSale() *Rice_sale {
	return new(Rice_sale)
}

func (data *Rice_sale) SaveRiceSale() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO rice_sale VALUES (NULL,?,?,?,SYSDATE())`, data.Warehouse_id, data.Rice_id, data.Amount)
	if err != nil {
		return err
	}
	return nil
}

func CreateRiceSale(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	body := req.Body
	sale := NewRiceSale()
	err := json.NewDecoder(body).Decode(sale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer body.Close()
	err = sale.SaveRiceSale()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// START TRANSACTION;
// 	INSERT INTO rice_transfer_record VALUES (NULL,"2021-05-27","ก1","ก5",20,"g56005-1-02","g56005-1-02");
// 	INSERT INTO rice_packing_daily VALUES (NULL,?,?,?,SYSDATE()"source"), data.Warehouse_id,data.Rice_id, data.Amount, data.;
// 	INSERT INTO rice_sale VALUES (NULL,"ก1","g56005-1-02",20,"2021-05-27");
// COMMIT;

// (data *TransferRice) FindTransferRiceByID(ID string) ([]TransferRice, error)
// INSERT INTO warehouse (id, details) VALUES (?, ?)`, data.Id, data.Details

// NewRouter return all router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index)
	//Warehouse
	router.GET("/warehouses", GetAllWarehouses)
	router.POST("/warehouses", CreateWarehouse)
	router.GET("/warehouses/:id", FindWarehouseByID)
	router.PUT("/warehouses/:id", UpdateWarehouse)
	router.DELETE("/warehouses/:id", DeleteWarehouse)
	//TransferRice
	router.GET("/transfers/:id", FindTransferRiceByID)
	router.POST("/transfers", CreateTransfer)
	router.POST("/packing", CreatePacking)
	router.POST("/sale", CreateRiceSale)

	// router.GET("/transfers/:id", FindTransferByID)
	return router
}

const (
	DBDriver   = "mysql"
	DBName     = "rice_mill_process"
	DBUser     = "root"
	DBPassword = "123456"
	DBURL      = DBUser + ":" + DBPassword + "@tcp(127.0.0.1:3306)/" + DBName
)

func GetDB() (*sql.DB, error) {
	db, err := sql.Open(DBDriver, DBURL)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func main() {
	log.Println("Server is up on 9000 port")
	router := NewRouter()
	log.Fatalln(http.ListenAndServe(":9000", router))
}

// func GetAllTransfer(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
// 	w.Header().Set("Content-Type", "application/json")
// 	users, err := NewTranfer().GetAllTransfers()
// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	err = json.NewEncoder(w).Encode(users)
// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// func (data *TransferRice) GetAllTransfers() ([]TransferRice, error) {
// 	transferRices := make([]TransferRice, 0)
// 	db, err := GetDB()
// 	if err != nil {
// 		return transferRices, err
// 	}
// 	defer db.Close()
// 	rows, err := db.Query(`SELECT rf.id as Id, rf.upd_date as Upd_date, rf.amount as Amount ,w1.details as To_warehouse_name, w2.details as From_warehouse_name , ri.rice_name as RiceNameFrom, ri2.rice_name as RiceNameTo,"โยกย้าย" as Status FROM rice_transfer_record rf join warehouse w1 on w1.id = rf.to_warehouse join warehouse w2 on w2.id = rf.from_warehouse join rice_info ri on rf.rice_id_from = ri.rice_id join rice_info ri2 on rf.rice_id_to = ri2.rice_id where w1.id = 'ก1' or w2.id = 'ก1' limit 10`)
// 	if err != nil {
// 		return transferRices, err
// 	}
// 	for rows.Next() {
// 		transferRice := TransferRice{}
// 		rows.Scan(&transferRice.Id, &transferRice.Upd_date, &transferRice.Amount, &transferRice.To_warehouse_name, &transferRice.From_warehouse_name, &transferRice.RiceNameFrom, &transferRice.RiceNameTo, &transferRice.Status)
// 		transferRices = append(transferRices, transferRice)
// 	}
// 	return transferRices, nil
// }

// func getResult(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var tests []RiceResult
// 	sql := "select rice_result.upd_date as Update_date , rice_info.rice_name as RiceName, rice_result.weight as Weight, rice_result.well as Well, rice_result.lose as Lose from rice_result, rice_info where rice_result.rice_id = rice_info.rice_id"
// 	result, err := db.Query(sql)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer result.Close()
// 	var test RiceResult
// 	for result.Next() {
// 		err = result.Scan(&test.Update_date, &test.RiceName, &test.Weight, &test.Well, &test.Lose)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		tests = append(tests, test)

// 	}
// 	json.NewEncoder(w).Encode(tests)
// }
