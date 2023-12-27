package main

import (
	m "Shopping/model"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/sessions"

	"bytes"

	"gopkg.in/gomail.v2"
)

var db *sql.DB
var tpl *template.Template
var tplMail *template.Template

type Result struct {
	T string
	F string
}

func main() {
	tpl, _ = tpl.ParseGlob("templates/*.html")
	tplMail, _ = tpl.ParseGlob("mail/*.html")
	var err error
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	db, err = sql.Open("mysql", "root:1234567890@tcp(localhost:3306)/shopcart")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//**Customer**
	http.HandleFunc("/loginCus", loginCustomer)
	http.HandleFunc("/", home)
	http.HandleFunc("/home", home)
	http.HandleFunc("/detailProduct", detailProduct)
	http.HandleFunc("/registerCustomer", registerCustomer)
	http.HandleFunc("/carts", carts)
	http.HandleFunc("/newCart", newCart)
	http.HandleFunc("/deleteCart", deleteCart)
	http.HandleFunc("/info-acc", information)
	http.HandleFunc("/change-password", changePassword)
	http.HandleFunc("/myOrder", myOrder)
	http.HandleFunc("/myOrderDetail", myOrderDetail)
	http.HandleFunc("/signOut", signOut)
	http.HandleFunc("/forgetPassword", forgetPassword)
	http.HandleFunc("/searchCate", searchCate)
	http.HandleFunc("/face", face)
	http.HandleFunc("/goo", goo)
	//**Mail customer**
	http.HandleFunc("/addMail", addMail)
	http.HandleFunc("/verify", verify)
	//**Beta**
	http.HandleFunc("/searchBeta", searchBeta)
	http.HandleFunc("/myOrderBeta", myOrderBeta)
	http.HandleFunc("/myOrderDetailBeta", myOrderDetailBeta)
	http.HandleFunc("/informationBeta", informationBeta)
	//**Guest**
	http.HandleFunc("/cartGuest", cartGuest)
	http.HandleFunc("/createGuest", createGuest)
	http.HandleFunc("/png", png)
	//**Seller**
	http.HandleFunc("/loginSel", loginSel)
	http.HandleFunc("/registerSel", registerSel)
	http.HandleFunc("/homeSel", homeSel)
	http.HandleFunc("/infoSeller", infoSeller)
	http.HandleFunc("/statistical", statistical)
	http.HandleFunc("/stat30days", stat30days)
	http.HandleFunc("/products", products)
	http.HandleFunc("/info", infoProduct)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/new", new)
	http.HandleFunc("/images", subImage)
	http.HandleFunc("/deleteImage", deleteImage)
	http.HandleFunc("/cusOrder", cusOrder)
	http.HandleFunc("/detailCusOrder", detailCusOrder)
	http.HandleFunc("/forgetPasswordForSeller", forgetPasswordForSeller)
	//**Mail seller**
	http.HandleFunc("/verifySeller", verifySeller)
	//**Admins**
	http.HandleFunc("/homeAdmin", homeAdmin)
	http.HandleFunc("/cate", categories)
	http.HandleFunc("/infoCate", infoCate)
	http.HandleFunc("/newCate", newCate)
	http.HandleFunc("/deleteCate", deleteCate)
	http.HandleFunc("/banner", banner)
	http.HandleFunc("/newBanner", newBanner)
	http.HandleFunc("/deleteBanner", deleteBanner)
	http.HandleFunc("/footerBanner", footerBanner)
	http.HandleFunc("/newFooterBanner", newFooterBanner)
	http.HandleFunc("/deleteFooterBanner", deleteFooterBanner)
	http.HandleFunc("/colors", subColor)
	http.HandleFunc("/deleteColor", deleteColor)
	//JSON Response
	http.HandleFunc("/pn", pn)
	http.HandleFunc("/upd", upd)
	http.HandleFunc("/upw", upw)
	http.HandleFunc("/jc", json_cart)
	http.HandleFunc("/buy", buy)
	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe("localhost:9990", nil))
}

func test(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "1test.html", nil)
	}
}

func goo(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("acc")
	if cookie != nil {
		http.Redirect(w, r, "home", http.StatusSeeOther)
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "1test.html", nil)
	}
	var username string
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		// fmt.Println(id, name, email)
		var cusId string
		check := checkGoo(id)
		if check == "" {
			db.Exec("INSERT INTO `shopcart`.`google` (`googleID`) VALUES ( ? );", id)
			row, _ := db.Query("SELECT id FROM shopcart.google where googleID = ?", id)
			for row.Next() {
				row.Scan(&cusId)
			}
			username = getRandomUserName(5)
			db.Exec("INSERT INTO `shopcart`.`information` (`GooID`, `FullName`, `UserName`, `Email`) VALUES (?, ?, ?, ?);", cusId, name, username, email)

			row, _ = db.Query("SELECT id FROM shopcart.information where GooID = ?", cusId)
			for row.Next() {
				row.Scan(&cusId)
			}

			db.Exec("INSERT INTO `shopcart`.`carts` (`CustomerID`) VALUES (?);", cusId)
		} else {
			row, _ := db.Query("SELECT id, UserName FROM shopcart.information where GooID = ?", check)
			for row.Next() {
				row.Scan(&cusId, &username)
			}
		}
		cookie = &http.Cookie{
			Name:    "acc",
			Value:   cusId + "/" + username,
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, cookie)
	}
}

func checkGoo(id string) string {
	row, _ := db.Query("SELECT id FROM shopcart.google where googleID = ?", id)
	id = ""
	for row.Next() {
		row.Scan(&id)
	}
	return id
}

func face(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("acc")
	if cookie != nil {
		http.Redirect(w, r, "home", http.StatusSeeOther)
	}
	var username string
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		var cusId string
		check := checkFace(id)
		if check == "" {
			fmt.Println(id, name)
			db.Exec("INSERT INTO `shopcart`.`facebook` (`facebookID`) VALUES ( ? );", id)
			row, _ := db.Query("SELECT id FROM shopcart.facebook where facebookID = ?", id)
			for row.Next() {
				row.Scan(&cusId)
			}
			username = getRandomUserName(5)
			db.Exec("INSERT INTO `shopcart`.`information` (`FaceID`, `FullName`, `UserName`, `Email`) VALUES (?, ?, ?, ?);", cusId, name, username, email)

			row, _ = db.Query("SELECT id FROM shopcart.information where FaceID = ?", cusId)
			for row.Next() {
				row.Scan(&cusId)
			}

			db.Exec("INSERT INTO `shopcart`.`carts` (`CustomerID`) VALUES (?);", cusId)
		} else {
			row, _ := db.Query("SELECT id, UserName FROM shopcart.information where FaceID = ?", check)
			for row.Next() {
				row.Scan(&cusId, &username)
			}
		}
		cookie = &http.Cookie{
			Name:    "acc",
			Value:   cusId + "/" + username,
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, cookie)
	}
}

func checkFace(id string) string {
	row, _ := db.Query("SELECT id FROM shopcart.facebook where facebookID = ?", id)
	id = ""
	for row.Next() {
		row.Scan(&id)
	}
	return id
}

func signOut(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "acc",
		MaxAge: 0,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func home(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT idPro, NamePro, Price, Discount, Star, Image FROM shopcart.products;")
	defer rows.Close()
	var p1 m.Product
	var ListPro []m.Product
	for rows.Next() {
		rows.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Discount, &p1.Star, &p1.Image)
		p1.Discount = math.Round(p1.Price*(1-p1.Discount)*100) / 100
		ListPro = append(ListPro, p1)
	}
	ListCate := listCategories()
	rows, _ = db.Query("SELECT idBanner, Image FROM shopcart.banner;")
	var b m.Banner
	var ListBanner []m.Banner
	for rows.Next() {
		rows.Scan(&b.IdBannerImg, &b.Image)
		ListBanner = append(ListBanner, b)
	}
	rows, _ = db.Query("SELECT idFB, Title, Content FROM shopcart.footer_banner;")
	var fb m.FooterBanner
	var ListFB []m.FooterBanner
	for rows.Next() {
		rows.Scan(&fb.IdBanner, &fb.Title, &fb.Content)
		ListFB = append(ListFB, fb)
	}
	rows, _ = db.Query("SELECT idPro, NamePro, Price, Discount, Star, Image FROM shopcart.products order by Discount desc limit 5")
	var listDeals []m.Product
	for rows.Next() {
		rows.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Discount, &p1.Star, &p1.Image)
		p1.Discount = math.Round(p1.Price*(1-p1.Discount)*100) / 100
		listDeals = append(listDeals, p1)
	}

	rows, _ = db.Query("SELECT idPro, NamePro, Price, Discount, Star, Image FROM shopcart.products order by Date desc limit 7")
	var listNews []m.Product
	for rows.Next() {
		rows.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Discount, &p1.Star, &p1.Image)
		p1.Discount = math.Round(p1.Price*(1-p1.Discount)*100) / 100
		listNews = append(listNews, p1)
	}

	rows, _ = db.Query("SELECT idPro, NamePro, Price, Discount, Star, Image FROM shopcart.products order by Date asc limit 7")
	var listLatest []m.Product
	for rows.Next() {
		rows.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Discount, &p1.Star, &p1.Image)
		p1.Discount = math.Round(p1.Price*(1-p1.Discount)*100) / 100
		listLatest = append(listLatest, p1)
	}

	rows, _ = db.Query("SELECT idPro, NamePro, Price, Discount, Star, Image FROM shopcart.products order by Sold desc limit 5")
	var listSpecials []m.Product
	for rows.Next() {
		rows.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Discount, &p1.Star, &p1.Image)
		p1.Discount = math.Round(p1.Price*(1-p1.Discount)*100) / 100
		listSpecials = append(listSpecials, p1)
	}

	type list struct {
		ListPro    []m.Product
		ListCate   []m.Category
		ListBanner []m.Banner
		ListFB     []m.FooterBanner
		Deals      []m.Product
		News       []m.Product
		Specials   []m.Product
		Latests    []m.Product
		S          []string
	}
	// var list1 = list{ListPro: ListPro, ListCate: ListCate, ListBanner: ListBanner, ListFB: ListFB, Deals: listDeals, News: listNews, Specials: listSpecials, Latests: listLatest}
	var list2 list
	list2.ListPro = append(list2.ListPro, ListPro...)
	list2.ListCate = append(list2.ListCate, ListCate...)
	list2.ListBanner = append(list2.ListBanner, ListBanner...)
	list2.ListFB = append(list2.ListFB, ListFB...)
	list2.Deals = append(list2.Deals, listDeals...)
	list2.News = append(list2.News, listNews...)
	list2.Specials = append(list2.Specials, listSpecials...)
	list2.Latests = append(list2.Latests, listLatest...)
	data := []string{"Chuỗi 1", "Chuỗi 2", "Chuỗi 3", "Chuỗi 4", "Chuỗi 5"}
	list2.S = append(list2.S, data...)
	tpl.ExecuteTemplate(w, "home.html", list2)
}

func myOrder(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		listOder := getInformationOrderByCustomerId(w, r)
		listCate := listCategories()
		type send struct {
			ListOrder []m.Order
			ListCate  []m.Category
		}
		send1 := send{ListOrder: listOder, ListCate: listCate}
		tpl.ExecuteTemplate(w, "cusOrder.html", send1)
		return
	}
}

func myOrderBeta(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		ty := r.FormValue("ty")
		_ = ty
		var listOB []m.OrderBeta
		cusId := getCustomerId(w, r)
		rows, _ := db.Query("SELECT id FROM shopcart.orders where CustomerID = ? order by id desc", cusId)
		defer rows.Close()
		var orderId string
		for rows.Next() {
			rows.Scan(&orderId)
			row, _ := db.Query("")
			defer rows.Close()
			if ty == "" {
				row, _ = db.Query("SELECT os.id , StoreID , s.Name , os.StatusID, st.Name FROM shopcart.orders_store os, shopcart.stores s, shopcart.status st where OrderID = ? and os.StoreID = s.id and os.StatusID = st.id ", orderId)
			} else {
				row, _ = db.Query("SELECT os.id , StoreID , s.Name , os.StatusID, st.Name FROM shopcart.orders_store os, shopcart.stores s, shopcart.status st where OrderID = ? and os.StoreID = s.id and os.StatusID = st.id and os.StatusID = ?", orderId, ty)
			}
			_ = row
			for row.Next() {
				var ob m.OrderBeta
				ob.Id = orderId
				row.Scan(&ob.OsId, &ob.Store.Id, &ob.Store.Name, &ob.Status.Id, &ob.Status.Name)
				row1, _ := db.Query("SELECT od.id , od.ProductID, od.Price, od.Amount, p.Image , p.NamePro FROM shopcart.orders_detail od, shopcart.products p where OrderSID = ? and od.ProductID = p.idPro", ob.OsId)
				defer row1.Close()
				for row1.Next() {
					var product m.Product
					var id string
					row1.Scan(&id, &product.Id, &product.Price, &product.Sold, &product.Image, &product.Name)
					ob.AllPrice += product.Price * float64(product.Sold)
					row2, _ := db.Query("SELECT c.id, c.NameColor FROM shopcart.orders_detail od ,shopcart.color c where OrderSID = ? and od.id = ? and od.ColorID = c.id;", ob.OsId, id)
					defer row2.Close()
					for row2.Next() {
						row2.Scan(&product.Color.Id, &product.Color.ColorUrl)
					}
					ob.ListProduct = append(ob.ListProduct, product)
				}
				listOB = append(listOB, ob)
			}
		}
		listCate := listCategories()
		type send struct {
			ListCate []m.Category
			ListOB   []m.OrderBeta
		}
		send1 := send{ListCate: listCate, ListOB: listOB}
		tpl.ExecuteTemplate(w, "myOrderBeta.html", send1)
	}
	id := r.FormValue("id")
	db.Exec("UPDATE `shopcart`.`orders_store` SET `StatusID` = '5' WHERE (`id` = ?);", id)
}

func searchBeta(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		key := "\"%" + r.FormValue("key") + "%\""
		sort := r.FormValue("sort")
		or := r.FormValue("or")
		star := r.FormValue("star")
		min := r.FormValue("minPrice")
		max := r.FormValue("maxPrice")

		var query string
		query = "SELECT idPro, NamePro, Price, Discount, Image, Star FROM shopcart.products where (NamePro like " + key + " or Description like " + key + ")"
		if star != "" {
			query += " and Star >= " + star
		}
		if min != "" {
			query += " and Price >= " + min
		}
		if max != "" {
			query += " and Price <= " + max
		}

		if sort != "" && or == "" {
			if sort == "2" {
				query += " order by Date DESC"
			} else if sort == "3" {
				query += " order by Sold DESC"
			}
		} else if sort != "" && or != "" {
			if or == "1" {
				query += " order by Price asc"
			} else if or == "2" {
				query += " order by Price DESC"
			}
		}

		rows, _ := db.Query(query)

		defer rows.Close()
		var p m.Product
		var listProduct []m.Product
		for rows.Next() {
			rows.Scan(&p.Id, &p.Name, &p.Price, &p.Discount, &p.Image, &p.Star)
			if p.Discount != 0 {
				p.Discount = math.Round(p.Price*(1-p.Discount)*100) / 100
			}
			listProduct = append(listProduct, p)
		}
		listCate := listCategories()
		type send struct {
			ListProduct []m.Product
			ListCate    []m.Category
		}
		tpl.ExecuteTemplate(w, "searchBeta.html", send{ListProduct: listProduct, ListCate: listCate})
	}
}

func searchCate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idCate := r.FormValue("id")
		_ = idCate
		rows, _ := db.Query("SELECT idPro, NamePro, Price, Discount, Image, Star FROM shopcart.products where idCate = ?", idCate)
		defer rows.Close()
		var p m.Product
		var listProduct []m.Product
		for rows.Next() {
			rows.Scan(&p.Id, &p.Name, &p.Price, &p.Discount, &p.Image, &p.Star)
			if p.Discount != 0 {
				p.Discount = math.Round(p.Price*(1-p.Discount)*100) / 100
			}
			listProduct = append(listProduct, p)
		}
		listCate := listCategories()
		type send struct {
			ListProduct []m.Product
			ListCate    []m.Category
		}
		tpl.ExecuteTemplate(w, "searchBeta.html", send{ListProduct: listProduct, ListCate: listCate})
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func myOrderDetailBeta(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		ido := r.FormValue("ido")
		idos := r.FormValue("idos")
		rows, _ := db.Query("SELECT os.id, o.NameCustomer, o.Phone, o.Address, o.Saddress, d.Name as Delivery, p.Name as Payment, os.StoreID, s.Name as Store, os.StatusID, st.Name as Status, os.DateOrder, os.DateReceive FROM shopcart.orders o, shopcart.orders_store os, shopcart.delivery d, shopcart.payment p, shopcart.stores s, shopcart.status st where o.id = ? and os.id = ? and o.DeliveryID = d.id and o.PaymentID = p.id and os.StoreID = s.id and os.StatusID = st.id", ido, idos)
		defer rows.Close()
		var o m.OrderDetailBeta
		for rows.Next() {
			rows.Scan(&o.OsId, &o.Info.FullName, &o.Info.PhoneNumber, &o.Info.Address, &o.Info.Saddress, &o.Delivery.Name, &o.Payment.Name, &o.Store.Id, &o.Store.Name, &o.Status.Id, &o.Status.Name, &o.DateOrder, &o.DateReceive)
		}
		row1, _ := db.Query("SELECT od.id , od.ProductID, od.Price, od.Amount, p.Image , p.NamePro FROM shopcart.orders_detail od, shopcart.products p where OrderSID = ? and od.ProductID = p.idPro", o.OsId)
		defer row1.Close()
		for row1.Next() {
			var product m.Product
			var id string
			row1.Scan(&id, &product.Id, &product.Price, &product.Sold, &product.Image, &product.Name)
			o.AllPrice += product.Price * float64(product.Sold)
			row2, _ := db.Query("SELECT c.id, c.NameColor FROM shopcart.orders_detail od ,shopcart.color c where OrderSID = ? and od.id = ? and od.ColorID = c.id;", o.OsId, id)
			defer row2.Close()
			for row2.Next() {
				row2.Scan(&product.Color.Id, &product.Color.ColorUrl)
			}
			o.ListProduct = append(o.ListProduct, product)
		}
		listCate := listCategories()
		type send struct {
			Od       m.OrderDetailBeta
			ListCate []m.Category
		}
		send1 := send{Od: o, ListCate: listCate}
		tpl.ExecuteTemplate(w, "myOrderDetailBeta.html", send1)
	}
}

func myOrderDetail(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		id := r.FormValue("id")
		row, _ := db.Query("SELECT id, CustomerID, NameCustomer, Phone, Address, Saddress, Date, DeliveryID, PaymentID FROM shopcart.orders where id = ? ", id)
		defer row.Close()
		var order m.Order
		for row.Next() {
			row.Scan(&order.Id, &order.CustomerID, &order.CustomerName, &order.Phone, &order.Address, &order.Saddress, &order.Date, &order.Delivery.Id, &order.Payment.Id)
			order.Delivery.Name = getNameDeliveryById(order.Delivery.Id)
			order.Payment.Name = getNamePaymentById(order.Payment.Id)
		}
		var o m.OrderDetail
		var listO []m.OrderDetail
		rows, _ := db.Query("SELECT id, OrderID, StoreID, StatusID FROM shopcart.orders_store where OrderID = ?", id)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&o.Id, &o.OrderID, &o.Store.Id, &o.Status.Id)
			row, _ = db.Query("SELECT Name FROM shopcart.stores where StoreID = ?", o.Store.Id)
			for row.Next() {
				row.Scan(&o.Store.Name)
			}
			row, _ = db.Query("SELECT Name FROM shopcart.payment WHERE StatusID = ?", o.Status.Id)
			for row.Next() {
				row.Scan(&o.Status.Name)
			}

			var po m.ProductOrder
			rows1, _ := db.Query("SELECT ProductID, ColorID, Amount, Price FROM shopcart.orders_detail where OrderSID = ?", o.Id)
			defer rows1.Close()
			for rows1.Next() {
				rows1.Scan(&po.ProductName, &po.ColorName, &po.Amount, &po.Price)
				row, _ = db.Query("SELECT NamePro FROM shopcart.products where idPro = ?", po.ProductName)
				for row.Next() {
					row.Scan(&po.ProductName)
				}
				if po.Price == "0" {
					po.ColorName = "None"
				} else {
					row, _ = db.Query("SELECT NameColor FROM shopcart.color where id = ?", po.ColorName)
					for row.Next() {
						row.Scan(&po.ColorName)
					}
				}
				o.ListPro = append(o.ListPro, po)
			}
			listO = append(listO, o)
		}

		listCate := listCategories()
		type send struct {
			Order     m.Order
			ListOrder []m.OrderDetail
			ListCate  []m.Category
		}
		send1 := &send{Order: order, ListOrder: listO, ListCate: listCate}
		tpl.ExecuteTemplate(w, "", send1)
		return
	}
	id := r.FormValue("id")
	status, _ := strconv.Atoi(r.FormValue("status"))
	if status == 5 {
		db.Exec("UPDATE `shopcart`.`orders_store` SET `StatusID` = ? WHERE (`id` = ?)", status, id)
	} else {
		db.Exec("UPDATE `shopcart`.`orders_store` SET `StatusID` = ? WHERE (`id` = ?)", status+1, id)
	}
}

func buy(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		cusId := getCustomerId(w, r)
		name := r.FormValue("name")
		address := r.FormValue("address")
		s_address := r.FormValue("s-address")
		phone := r.FormValue("phone")
		delivery := r.FormValue("delivery")
		pay := r.FormValue("ra-payment")
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006/01/02 15:04:05")

		db.Exec("BEGIN;")
		db.Exec("SELECT * FROM shopcart.products for update")

		db.Exec("INSERT INTO `shopcart`.`orders` (`CustomerID`, `NameCustomer`, `Phone`, `Address`,`Saddress`, `Date`, `DeliveryID`,`PaymentID`) VALUES (?,?,?,?,?,?,?,?);", cusId, name, phone, address, s_address, formattedTime, delivery, pay)
		var idCart string

		row, _ := db.Query("SELECT id FROM shopcart.carts where CustomerID = ?", cusId)
		defer row.Close()
		for row.Next() {
			row.Scan(&idCart)
		}

		var c m.Cart
		var listCart []m.Cart

		row1, _ := db.Query("SELECT id, CartID, ProductID, ColorID, Amount FROM shopcart.carts_product where CartID = ?", idCart)
		defer row1.Close()
		for row1.Next() {
			row1.Scan(&c.Id, &c.CartID, &c.ProductID, &c.ColorID, &c.Amount)
			listCart = append(listCart, c)
		}

		var listStores []string
		var idStore string
	k:
		for _, store := range listCart {
			idStore = getStoreIdByProductId(store.ProductID)
			for _, store := range listStores {
				if store == idStore {
					continue k
				}
			}
			listStores = append(listStores, idStore)
		}

		var idOrder string
		row, _ = db.Query("SELECT id FROM shopcart.orders where CustomerID = ? and NameCustomer = ? and Phone = ?;", cusId, name, phone)
		for row.Next() {
			row.Scan(&idOrder)
		}

		for _, store := range listStores {
			db.Exec("INSERT INTO `shopcart`.`orders_store` (`OrderID`, `StoreID`, `StatusID`) VALUES (?, ?, '1');", idOrder, store)
		}

		var dis1 string
		var dis2 string
		var price string
		var orderSID string
		done := true
		for _, cart := range listCart {
			var am int
			row, _ := db.Query("SELECT InStock FROM shopcart.products where idPro = ?", cart.ProductID)
			for row.Next() {
				row.Scan(&am)
			}
			if am > cart.Amount {
				done = false
				break
			}
			if cart.ColorID == "0" {
				row, _ = db.Query("SELECT p.StoreID, p.Discount, p.Price FROM shopcart.`products` p where p.idPro = ?", cart.ProductID)
				for row.Next() {
					row.Scan(&idStore, &dis1, &price)
				}
				dis2 = "0"
			} else {
				row, _ = db.Query("SELECT p.StoreID, p.Discount, p.Price, c.Discount_color FROM shopcart.`products` p , shopcart.product_color c where p.idPro = ? and c.id = ?", cart.ProductID, cart.ColorID)
				for row.Next() {
					row.Scan(&idStore, &dis1, &price, &dis2)
				}
			}
			d1, _ := strconv.ParseFloat(dis1, 64)
			d2, _ := strconv.ParseFloat(dis2, 64)
			p, _ := strconv.ParseFloat(price, 64)
			total := math.Round((p-(p*(d1+d2)))*100) / 100
			row, _ = db.Query("SELECT id FROM shopcart.orders_store where OrderID = ? and StoreID = ?", idOrder, idStore)
			for row.Next() {
				row.Scan(&orderSID)
			}
			db.Exec("INSERT INTO `shopcart`.`orders_detail` (`OrderSID`, `ProductID`, `ColorID`, `Amount`, `Price`) VALUES (?, ?, ?, ?, ?);", orderSID, cart.ProductID, cart.ColorID, cart.Amount, total)
			db.Exec("UPDATE `shopcart`.`products` SET `InStock` = ? WHERE (`idPro` = ? );", am-cart.Amount, cart.ProductID)
		}

		if done {
			db.Exec("DELETE FROM `shopcart`.`carts_product` WHERE (`CartID` = ?)", idCart)
			db.Exec("COMMIT;")
			tpl.ExecuteTemplate(w, "1test.html", "Buy successfully")
		} else {
			db.Exec("ROLLBACK;")
			db.Exec("COMMIT;")
			tpl.ExecuteTemplate(w, "1test.html", "Buy failed")
		}

	}
}

func information(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		cusId := getCustomerId(w, r)
		row, _ := db.Query("SELECT id, AccID, FullName, Image, PhoneNumber, Email FROM shopcart.information where id = ? ", cusId)
		defer row.Close()
		var info m.Info
		for row.Next() {
			row.Scan(&info.Id, &info.AccID, &info.FullName, &info.Image, &info.PhoneNumber, &info.Email)
		}
		listCate := listCategories()
		type in struct {
			I        m.Info
			ListCate []m.Category
			Oke      string
			Err      string
		}
		enter := in{I: info, ListCate: listCate}
		// enter.Oke = "Update successfully"
		tpl.ExecuteTemplate(w, "information_account.html", enter)
		return
	}
	cusId := getCustomerId(w, r)
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	checkImg := r.FormValue("checkImg")
	img := ""
	if checkImg != "" {
		file, fileHeader, _ := r.FormFile("img")
		defer file.Close()
		contentType := fileHeader.Header["Content-Type"][0]
		fmt.Println(fileHeader.Header["Content-Type"][0])
		var osFile *os.File
		defer osFile.Close()
		if contentType == "image/jpeg" || contentType == "image/png" {
			osFile, _ = ioutil.TempFile("static/imgUser/", "*.jpg")
		} else {
			tpl.ExecuteTemplate(w, "/home", nil)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		osFile.Write(fileBytes)
		s := osFile.Name()
		fmt.Println(osFile.Name())
		img = "../" + s
	}
	db.Exec("UPDATE `shopcart`.`information` SET `FullName` = ?, `Image` = ?, `PhoneNumber` = ?, `Email` = ? WHERE (`id` = ?);", name, img, phone, email, cusId)

	row, _ := db.Query("SELECT id, AccID, FullName, Image, PhoneNumber, Email FROM shopcart.information where id = ? ", cusId)
	defer row.Close()
	var info m.Info
	for row.Next() {
		row.Scan(&info.Id, &info.AccID, &info.FullName, &info.Image, &info.PhoneNumber, &info.Email)
	}
	listCate := listCategories()
	type in struct {
		I        m.Info
		ListCate []m.Category
		Oke      string
		Err      string
	}
	enter := in{I: info, ListCate: listCate}
	enter.Oke = "Update successfully"
	tpl.ExecuteTemplate(w, "information_account.html", enter)
}

func informationBeta(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		cusId := getCustomerId(w, r)
		fmt.Println(cusId)
		row, _ := db.Query("SELECT id, AccID, UserName, FullName, Image, PhoneNumber, Email, BirthDay, Gender FROM shopcart.information where id = ? ", cusId)
		defer row.Close()
		var info m.Info
		for row.Next() {
			row.Scan(&info.Id, &info.AccID, &info.UserName, &info.FullName, &info.Image, &info.PhoneNumber, &info.Email, &info.BirthDay, &info.Gender)
		}
		listCate := listCategories()
		type in struct {
			I        m.Info
			ListCate []m.Category
		}
		enter := in{I: info, ListCate: listCate}
		tpl.ExecuteTemplate(w, "myInfoBeta.html", enter)
		return
	}
	cusId := getCustomerId(w, r)
	userName := r.FormValue("userName")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	gender := r.FormValue("gender")
	checkImg := r.FormValue("checkImg")
	brith := r.FormValue("brith")
	img := ""
	if checkImg != "" {
		file, fileHeader, _ := r.FormFile("file")
		defer file.Close()
		contentType := fileHeader.Header["Content-Type"][0]
		fmt.Println(fileHeader.Header["Content-Type"][0])
		var osFile *os.File
		defer osFile.Close()
		if contentType == "image/jpeg" || contentType == "image/png" {
			osFile, _ = ioutil.TempFile("static/imgUser/", "*.jpg")
		} else {
			tpl.ExecuteTemplate(w, "/home", nil)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		osFile.Write(fileBytes)
		s := osFile.Name()
		fmt.Println(osFile.Name())
		img = "../" + s
	}
	db.Exec("UPDATE `shopcart`.`information` SET `UserName` = ?, `FullName` = ?, `Image` = ?, `PhoneNumber` = ?, `Email` = ?, `BirthDay` = ?, `Gender` = ? WHERE (`id` = ?);", userName, name, img, phone, email, brith, gender, cusId)
	tpl.ExecuteTemplate(w, "/informationBeta", nil)
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		listCate := listCategories()
		type in struct {
			Err      string
			Oke      string
			ListCate []m.Category
		}
		enter := in{ListCate: listCate}
		tpl.ExecuteTemplate(w, "changePassword.html", enter)
		return
	}
	if r.Method == "POST" {
		type in struct {
			Err      string
			Oke      string
			ListCate []m.Category
		}
		listCate := listCategories()
		enter := in{ListCate: listCate}

		op := r.FormValue("old-pass")
		np := r.FormValue("new-pass")
		cp := r.FormValue("confirm-pass")
		if np != cp {
			enter.Err = "Password không trùng nhau"
			tpl.ExecuteTemplate(w, "changePassword.html", enter)
			return
		}
		cusId := getCustomerId(w, r)
		row, _ := db.Query("SELECT Password FROM shopcart.customers where  id = ?", cusId)
		defer row.Close()
		var password string
		for row.Next() {
			row.Scan(&password)
		}
		if op != password {
			enter.Err = "Password cũ sai"
			tpl.ExecuteTemplate(w, "changePassword.html", enter)
			return
		}
		db.Exec("UPDATE `shopcart`.`customers` SET `Password` = ? WHERE (id = ?)", np, cusId)
		enter.Oke = "Change password successfully"
		tpl.ExecuteTemplate(w, "changePassword.html", enter)
	}
}

func pn(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		cusId := getCustomerId(w, r)
		row, _ := db.Query("SELECT id FROM shopcart.carts where CustomerID = ?", cusId)
		for row.Next() {
			row.Scan(&cusId)
		}
		row1, _ := db.Query("SELECT id, ProductID FROM shopcart.carts_product where CartID = ?", cusId)
		var list []string
		var listPro []string
		var pro string
		for row1.Next() {
			row1.Scan(&cusId, &pro)
			list = append(list, cusId)
			listPro = append(listPro, pro)
		}
		index := r.FormValue("index")
		am := r.FormValue("am")
		index1, _ := strconv.Atoi(index)
		id := list[index1]
		idPro := listPro[index1]
		am1, _ := strconv.Atoi(am)
		var am2 int
		row, _ = db.Query("SELECT InStock FROM shopcart.products where idPro = ?", idPro)
		for row.Next() {
			row.Scan(&am2)
		}
		if am1 <= am2 {
			db.Exec("UPDATE `shopcart`.`carts_product` SET `Amount` = ? WHERE (`id` = ?);", am1, id)
		} else {
			db.Exec("UPDATE `shopcart`.`carts_product` SET `Amount` = ? WHERE (`id` = ?);", am2, id)
		}
		jsonData, err := json.Marshal(am2)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func newCart(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		cusId := getCustomerId(w, r)
		idCart := ""
		row, _ := db.Query("SELECT id FROM shopcart.carts where CustomerID = ?", cusId)
		defer row.Close()
		for row.Next() {
			row.Scan(&idCart)
		}
		idPro := r.FormValue("idP")
		idColor := r.FormValue("idC")
		amount := r.FormValue("amount")
		if strings.Contains(idColor, "o") {
			idColor = "0"
		} else {
			row, _ = db.Query("SELECT id FROM shopcart.color where NameColor = ?", idColor)
			defer row.Close()
			for row.Next() {
				row.Scan(&idColor)
			}
		}

		row, _ = db.Query("SELECT id, Amount FROM shopcart.carts_product where CartID = ? and ProductID = ? and ColorID = ?", idCart, idPro, idColor)
		var amount1 string
		var id1 string
		for row.Next() {
			row.Scan(&id1, &amount1)
		}

		var amount2 int
		row, _ = db.Query("SELECT InStock FROM shopcart.products where idPro = ?", idPro)
		for row.Next() {
			row.Scan(&amount2)
		}

		if id1 == "" {
			db.Exec("INSERT INTO `shopcart`.`carts_product` (`CartID`, `ProductID`, `ColorID`, `Amount`) VALUES (?, ?, ?, ?);", idCart, idPro, idColor, amount)
		} else {
			am, _ := strconv.Atoi(amount)
			am1, _ := strconv.Atoi(amount1)
			am2 := am + am1
			if am2 <= amount2 {
				db.Exec("UPDATE `shopcart`.`carts_product` SET `Amount` = ? WHERE (`id` = ?);", am2, id1)
			}
		}

		row, _ = db.Query("SELECT Amount FROM shopcart.carts_product where CartID = ? and ProductID = ? and ColorID = ?", idCart, idPro, idColor)
		for row.Next() {
			row.Scan(&amount)
		}
		jsonData, err := json.Marshal(amount)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func deleteCart(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		cusId := getCustomerId(w, r)
		row, _ := db.Query("SELECT id FROM shopcart.carts where CustomerID = ?", cusId)
		defer row.Close()
		for row.Next() {
			row.Scan(&cusId)
		}
		row1, _ := db.Query("SELECT id FROM shopcart.carts_product where CartID = ?", cusId)
		var list []string
		for row1.Next() {
			row1.Scan(&cusId)
			list = append(list, cusId)
		}
		index := r.FormValue("index")
		index1, _ := strconv.Atoi(index)
		id := list[index1]
		// fmt.Println(index, index1, id)
		db.Exec("DELETE FROM `shopcart`.`carts_product` WHERE (`id` = ?);", id)
	}
}

func carts(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		cartGuest(w, r)
		return
	}
	if r.Method == "GET" {
		cusId := getCustomerId(w, r)
		row, _ := db.Query("SELECT id FROM shopcart.carts where CustomerID = ?", cusId)
		defer row.Close()
		var catId string
		for row.Next() {
			row.Scan(&catId)
		}
		var c m.Cart
		type ty struct {
			P  m.Product
			Co m.Color
			Ca m.Cart
		}
		var list []ty
		var ty1 ty
		row1, _ := db.Query("SELECT id, CartID, ProductID, ColorID, Amount FROM shopcart.carts_product where CartID = ?", catId)
		if row1 == nil {
			tpl.ExecuteTemplate(w, "/home", nil)
			return
		}
		for row1.Next() {
			row1.Scan(&c.Id, &c.CartID, &c.ProductID, &c.ColorID, &c.Amount)
			ty1.Ca = c
			var p1 m.Product
			row, _ = db.Query("SELECT idPro, NamePro, Price, Star, Discount, Image, `Date`, Sold, InStock, idCate FROM shopcart.products where idPro = ?", c.ProductID)
			for row.Next() {
				row.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Star, &p1.Discount, &p1.Image, &p1.Date, &p1.Sold, &p1.InStock, &p1.Category.IdCate)
			}
			ty1.P = p1
			row, _ = db.Query("SELECT ProductID, Color, Discount_color FROM shopcart.product_color where ProductID = ? and Color = ?", c.ProductID, c.ColorID)
			var c1 m.Color
			for row.Next() {
				row.Scan(&c1.ProductId, &c1.ColorUrl, &c1.Discount)
				row, _ := db.Query("SELECT NameColor FROM shopcart.color where id = ?", c1.ColorUrl)
				defer row.Close()
				for row.Next() {
					row.Scan(&c1.ColorUrl)
				}
			}
			ty1.Co = c1
			list = append(list, ty1)
		}
		listCate := listCategories()
		rows, _ := db.Query("SELECT code, name FROM vietnam.provinces;")
		defer rows.Close()
		var add m.Address
		var listProvince []m.Address
		for rows.Next() {
			rows.Scan(&add.Code, &add.Name)
			listProvince = append(listProvince, add)
		}
		rows, _ = db.Query("SELECT id, Name FROM shopcart.payment;")
		var pay m.Payment
		var listPay []m.Payment
		for rows.Next() {
			rows.Scan(&pay.Id, &pay.Name)
			listPay = append(listPay, pay)
		}
		rows, _ = db.Query("SELECT id, Name FROM shopcart.delivery")
		var del m.Delivery
		var listDel []m.Delivery
		for rows.Next() {
			rows.Scan(&del.Id, &del.Name)
			listDel = append(listDel, del)
		}
		type send struct {
			ListCate     []m.Category
			ListCart     []ty
			ListProvince []m.Address
			ListPay      []m.Payment
			ListDel      []m.Delivery
		}
		send1 := send{ListCate: listCate, ListCart: list, ListProvince: listProvince, ListPay: listPay, ListDel: listDel}
		tpl.ExecuteTemplate(w, "cart.html", send1)
		return
	}
}

func detailProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.FormValue("id")
		if id == "" {
			http.Redirect(w, r, "homeAdmin", http.StatusSeeOther)
		}
		row, _ := db.Query("SELECT idPro, NamePro, Price, Star, Discount, Image, Date, Sold, Description, InStock, idCate FROM shopcart.products where idPro = ? ", id)
		defer row.Close()
		var p m.Product
		for row.Next() {
			row.Scan(&p.Id, &p.Name, &p.Price, &p.Star, &p.Discount, &p.Image, &p.Date, &p.Sold, &p.Description, &p.InStock, &p.IdCate)
		}
		p.Discount = math.Round(p.Price*(1-p.Discount)*100) / 100
		row, _ = db.Query("SELECT id, ProductID, Color, Discount_color FROM shopcart.product_color where ProductID = ?", id)
		var color m.Color
		var listColor []m.Color
		for row.Next() {
			row.Scan(&color.Id, &color.ProductId, &color.ColorUrl, &color.Discount)
			row1, _ := db.Query("SELECT NameColor FROM shopcart.color where id = ?", color.ColorUrl)
			defer row1.Close()
			for row1.Next() {
				row1.Scan(&color.ColorUrl)
			}
			listColor = append(listColor, color)
		}
		row, _ = db.Query("SELECT id, ProductID, Image FROM shopcart.product_slider where ProductID = ?", id)
		var image m.Image
		var listImage []m.Image
		for row.Next() {
			row.Scan(&image.Id, &image.ProductId, &image.Url)
			listImage = append(listImage, image)
		}
		listCate := listCategories()
		rows, _ := db.Query("SELECT id, CustomerID, ProductID, ColorID, Comment, Star, Date  FROM shopcart.comments where ProductID = ?", id)
		defer rows.Close()
		var co m.Comment
		var listComment []m.Comment
		var rating float64
		var allRating float64
		for rows.Next() {
			rows.Scan(&co.Id, &co.Info.AccID, &co.Product.Id, &co.Color.Id, &co.Com, &co.Star, &co.Date)
			row, _ = db.Query("SELECT id, AccID, FullName, Image, PhoneNumber, Email FROM shopcart.information where id = ?", co.Info.AccID)
			for row.Next() {
				row.Scan(&co.Info.Id, &co.Info.AccID, &co.Info.FullName, &co.Info.Image, &co.Info.PhoneNumber, &co.Info.Email)
			}
			row, _ = db.Query("SELECT NameColor FROM shopcart.color where  id = ?", co.Color.Id)
			for row.Next() {
				row.Scan(&co.Color.ColorUrl)
			}
			listComment = append(listComment, co)
			rating, _ = strconv.ParseFloat(co.Star, 64)
			allRating += rating
		}
		var pSam []m.Product
		var p1 m.Product
		rows, _ = db.Query("SELECT idPro, NamePro, Image,  round(price*Discount, 2) FROM shopcart.products where idCate = ? and idPro != ? limit 5", p.Category.IdCate, p.Id)
		for rows.Next() {
			rows.Scan(&p1.Id, &p1.Name, &p1.Image, &p1.Price)
			pSam = append(pSam, p1)
		}

		type lit struct {
			ListCate    []m.Category
			P           m.Product
			C           []m.Color
			I           []m.Image
			ListComment []m.Comment
			All         int
			Img         int
			Sta         int
			Pa          int
			AllSta      string
			ProductSam  []m.Product
		}
		all := len(listComment)
		allRating = allRating / float64(len(listComment))
		a := strconv.FormatFloat(allRating, 'f', 1, 64)
		n := lit{P: p, C: listColor, I: listImage, ListCate: listCate, ListComment: listComment, All: all, AllSta: a, ProductSam: pSam}
		tpl.ExecuteTemplate(w, "detailProduct.html", n)
		return
	}
	http.Redirect(w, r, "homeAdmin", http.StatusSeeOther)
}

func loginCustomer(w http.ResponseWriter, r *http.Request) {
	if checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "loginCus.html", nil)
		return
	}
	user := r.FormValue("username")
	pass := r.FormValue("password")
	row, _ := db.Query("SELECT id FROM shopcart.customers where User = ? and Password = ?", user, pass)
	defer row.Close()
	var id string
	for row.Next() {
		row.Scan(&id)
	}
	var cusId, username string

	if id == "" {
		row, _ = db.Query("SELECT  i.id, i.UserName FROM shopcart.information i, shopcart.customers c where AccID != \"\" and AccID = c.id and Email = ? and Password = ?", user, pass)

		for row.Next() {
			row.Scan(&cusId, &username)
		}
	} else {
		row, _ = db.Query("SELECT id, UserName FROM shopcart.information where  AccID = ?", id)
		for row.Next() {
			row.Scan(&cusId, &username)
		}
	}

	cookie := &http.Cookie{
		Name:    "acc",
		Value:   cusId + "/" + username,
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
	guest, error1 := r.Cookie("guest")
	if error1 == nil {
		db.Exec("DELETE FROM `shopcart`.`guest_product` WHERE (`CartID` = ?);", guest.Value)
	}
	http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
}

func registerCustomer(w http.ResponseWriter, r *http.Request) {
	_, error := r.Cookie("acc")
	if error == nil {
		http.Redirect(w, r, "homeAdmin", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "register.html", nil)
		return
	}
	email := r.FormValue("email")
	user := r.FormValue("username")
	pass := r.FormValue("password")
	c_pass := r.FormValue("c-password")
	err := ""
	if pass != c_pass {
		err = "The two passwords do not match"
		tpl.ExecuteTemplate(w, "register.html", err)
		return
	}
	if checkCustomer(user) {
		err = "The user is already registered"
		tpl.ExecuteTemplate(w, "register.html", err)
		return
	}
	if checkEmailRe(email) {
		err = "The email is already registered"
		tpl.ExecuteTemplate(w, "register.html", err)
		return
	}
	db.Exec("INSERT INTO `shopcart`.`customers` (`User`, `Password`) VALUES (?,?);", user, pass)
	row, _ := db.Query("SELECT id FROM shopcart.customers where User = ?", user)
	defer row.Close()
	var accId string
	for row.Next() {
		row.Scan(&accId)
	}

	db.Exec("INSERT INTO `shopcart`.`information` (`AccID`, `Email`) VALUES (?,?);", accId, email)
	row, _ = db.Query("SELECT id FROM shopcart.information where AccID = ?", accId)
	var cus string
	for row.Next() {
		row.Scan(&cus)
	}

	db.Exec("INSERT INTO `shopcart`.`carts` (`CustomerID`) VALUES (?);", cus)
	err = "Register successful"

	guest, error1 := r.Cookie("guest")

	if error1 == nil {
		var cusID string
		row, _ := db.Query("SELECT id FROM shopcart.customers where User = ? ", user)
		defer row.Close()
		for row.Next() {
			row.Scan(&cusID)
		}
		var cartID string
		row, _ = db.Query("SELECT id FROM shopcart.carts where CustomerID = ? ", cusID)
		for row.Next() {
			row.Scan(&cartID)
		}
		var proID string
		var colorID string
		var amount int
		rows, _ := db.Query("SELECT ProductID, ColorID, Amount FROM shopcart.guest_product where id  = ?", guest)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&proID, &colorID, &amount)
			var amount1 int
			var id int
			row1, err := db.Query("SELECT id, Amount FROM shopcart.carts_product where CartID = ? and ProductID = ? and ColorID = ?", cartID, proID, colorID)
			if err != nil {
				db.Exec("INSERT INTO `shopcart`.`carts_product` (`CartID`, `ProductID`, `ColorID`, `Amount`) VALUES (?, ?, ?, ?);", cartID, proID, colorID, amount)
			} else {
				row1.Scan(&id, &amount1)
				amount += amount1
				db.Exec("UPDATE `shopcart`.`carts_product` SET `Amount` = ? WHERE (`id` = ?);", amount, id)
			}
			row1.Close()
		}
		db.Exec("DELETE FROM `shopcart`.`guest_product` WHERE (`CartID` = ?);", guest.Value)
	}
	// http.Redirect(w, r, "/loginCus", http.StatusSeeOther)
	tpl.ExecuteTemplate(w, "register.html", err)
}

func forgetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "forgetPassword", nil)
	} else if r.Method == "POST" {
		email := r.FormValue("email")
		if checkEmail(email) {
			sendChangePasswordCustomer(email)
			send := Result{T: "Đã gửi mật khẩu mới đến Email của bạn"}
			tpl.ExecuteTemplate(w, "forgetPassword", send)
		} else {
			send := Result{F: "Không tìm thấy Email"}
			tpl.ExecuteTemplate(w, "forgetPassword", send)
		}
	}
}

// ********************************mail customer*********************************
func sendWelcomeCustomer(email string) {
	var body bytes.Buffer
	t, err := template.ParseFiles("./sendWelcomeCustomer.html")
	t.Execute(&body, struct{ Name string }{Name: getNewPassword()})
	if err != nil {
		fmt.Println(err)
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "vietahh02@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", body.String())
	m.Attach("./img/LOGO.png")
	d := gomail.NewDialer("smtp.gmail.com", 587, "vietahh02@gmail.com", "iuam zbym wyaq euni")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func sendChangePasswordCustomer(email string) {
	var body bytes.Buffer
	t, err := template.ParseFiles("./sendChangePassword.html")
	t.Execute(&body, struct{ Name string }{Name: getNewPassword()})
	if err != nil {
		fmt.Println(err)
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "vietahh02@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", body.String())
	m.Attach("./img/LOGO.png")
	d := gomail.NewDialer("smtp.gmail.com", 587, "vietahh02@gmail.com", "iuam zbym wyaq euni")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func sendVerifyMail(email, key string) {
	var body bytes.Buffer
	t, err := template.ParseFiles("./sendVerify.html")
	t.Execute(&body, struct{ Name string }{Name: key})
	if err != nil {
		fmt.Println(err)
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "vietahh02@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", body.String())
	m.Attach("./img/LOGO.png")
	d := gomail.NewDialer("smtp.gmail.com", 587, "vietahh02@gmail.com", "iuam zbym wyaq euni")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func getNewPassword() string {
	var randomCharacters []byte
	for i := 0; i < 8; i++ {
		randomCharacter := byte(rand.Intn(62)) + 48
		if randomCharacter > 57 && randomCharacter < 65 {
			randomCharacter += 7
		}
		randomCharacters = append(randomCharacters, randomCharacter)
	}
	randomString := string(randomCharacters)
	return randomString
}

func getKey() string {
	var randomNumbers []byte
	for i := 0; i < 6; i++ {
		randomNumber := byte(rand.Intn(10)) + 48
		randomNumbers = append(randomNumbers, randomNumber)
	}
	randomString := string(randomNumbers)
	return randomString
}

func addMail(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		email := r.FormValue("email")
		if checkEmail(email) {
			tpl.ExecuteTemplate(w, "verifyingMail.html", "f")
		}
		key := getKey()
		if checkNewEmail(email) {
			db.Exec("UPDATE `shopcart`.`verifynew` SET `Key` = ? WHERE (`Email` = ?);", key, email)
		} else {
			db.Exec("INSERT INTO `shopcart`.`verifynew` (`Email`, `Key`) VALUES (?, ?);", email, key)
		}
		sendVerifyMail(email, key)
		tplMail.ExecuteTemplate(w, "confirmMail.html", email)
		return
	}
	tplMail.ExecuteTemplate(w, "verifyingMail.html", nil)
}

func verify(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("acc")
	if err != nil {
		return
	}
	if r.Method == "POST" {
		email := r.FormValue("email")
		key := r.FormValue("key")
		row, _ := db.Query("SELECT * FROM shopcart.verifynew where Email = ? and KeyN = ?", email, key)
		var confirm string
		if row != nil {
			a := strings.Split(cookie.Value, "/")
			row1, _ := db.Query("SELECT id FROM shopcart.information where CustomerID = ?", a[0])
			var id string
			for row1.Next() {
				row1.Scan(&id)
			}
			db.Exec("UPDATE `shopcart`.`information` SET `Email` = ? WHERE (`id` = ?);", email, id)
			confirm = "Y"
			sendWelcomeCustomer(email)
		} else {
			confirm = "N"
		}
		jsonData, err := json.Marshal(confirm)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func checkNewEmail(email string) bool {
	row, _ := db.Query("SELECT * FROM shopcart.verifynew where Email = ?", email)
	return row != nil
}

// ********************************Guest****************************************************************

func cartGuest(w http.ResponseWriter, r *http.Request) {
	if checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		cookie, err := r.Cookie("guest")
		var c m.Cart
		type ty struct {
			P  m.Product
			Co m.Color
			Ca m.Cart
		}
		var list []ty
		var ty1 ty
		if err == nil {
			row1, _ := db.Query("SELECT id, CartID, ProductID, ColorID, Amount FROM shopcart.guest_product where CartID = ?", cookie.Value)
			if row1 == nil {
				tpl.ExecuteTemplate(w, "/home", nil)
				return
			}
			for row1.Next() {
				row1.Scan(&c.Id, &c.CartID, &c.ProductID, &c.ColorID, &c.Amount)
				ty1.Ca = c
				var p1 m.Product
				row, _ := db.Query("SELECT idPro, NamePro, Price, Star, Discount, Image, `Date`, Sold, InStock, idCate FROM shopcart.products where idPro = ?", c.ProductID)
				defer row.Close()
				for row.Next() {
					row.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Star, &p1.Discount, &p1.Image, &p1.Date, &p1.Sold, &p1.InStock, &p1.Category.IdCate)
				}
				ty1.P = p1
				row, _ = db.Query("SELECT ProductID, Color, Discount_color FROM shopcart.product_color where ProductID = ? and Color = ?", c.ProductID, c.ColorID)
				var c1 m.Color
				for row.Next() {
					row.Scan(&c1.ProductId, &c1.ColorUrl, &c1.Discount)
					row, _ := db.Query("SELECT NameColor FROM shopcart.color where id = ?", c1.ColorUrl)
					defer row.Close()
					for row.Next() {
						row.Scan(&c1.ColorUrl)
					}
				}
				ty1.Co = c1
				list = append(list, ty1)
			}
		}
		listCate := listCategories()
		rows, _ := db.Query("SELECT code, name FROM vietnam.provinces;")
		defer rows.Close()
		var add m.Address
		var listProvince []m.Address
		for rows.Next() {
			rows.Scan(&add.Code, &add.Name)
			listProvince = append(listProvince, add)
		}
		rows, _ = db.Query("SELECT id, Name FROM shopcart.payment;")
		var pay m.Payment
		var listPay []m.Payment
		for rows.Next() {
			rows.Scan(&pay.Id, &pay.Name)
			listPay = append(listPay, pay)
		}
		rows, _ = db.Query("SELECT id, Name FROM shopcart.delivery")
		var del m.Delivery
		var listDel []m.Delivery
		for rows.Next() {
			rows.Scan(&del.Id, &del.Name)
			listDel = append(listDel, del)
		}
		type send struct {
			ListCate     []m.Category
			ListCart     []ty
			ListProvince []m.Address
			ListPay      []m.Payment
			ListDel      []m.Delivery
		}
		send1 := send{ListCate: listCate, ListCart: list, ListProvince: listProvince, ListPay: listPay, ListDel: listDel}
		tpl.ExecuteTemplate(w, "cart.html", send1)
		return
	}

	cookie, _ := r.Cookie("guest")
	idCart := cookie.Value
	idPro := r.FormValue("idP")
	idColor := r.FormValue("idC")
	amount := r.FormValue("amount")
	if strings.Contains(idColor, "o") {
		idColor = "0"
	} else {
		row, _ := db.Query("SELECT id FROM shopcart.color where NameColor = ?", idColor)
		defer row.Close()
		for row.Next() {
			row.Scan(&idColor)
		}
	}
	row, _ := db.Query("SELECT id, Amount FROM shopcart.guest_product where CartID = ? and ProductID = ? and ColorID = ?", idCart, idPro, idColor)
	defer row.Close()
	var amount1 string
	var id1 string
	for row.Next() {
		row.Scan(&id1, &amount1)
	}

	var amount2 int
	row, _ = db.Query("SELECT InStock FROM shopcart.products where idPro = ?", idPro)
	for row.Next() {
		row.Scan(&amount2)
	}

	if id1 == "" {
		db.Exec("INSERT INTO `shopcart`.`guest_product` (`CartID`, `ProductID`, `ColorID`, `Amount`) VALUES (?, ?, ?, ?);", idCart, idPro, idColor, amount)
	} else {
		am, _ := strconv.Atoi(amount)
		am1, _ := strconv.Atoi(amount1)
		am2 := am + am1
		if am2 <= amount2 {
			db.Exec("UPDATE `shopcart`.`guest_product` SET `Amount` = ? WHERE (`id` = ?);", am2, id1)
		}
	}

	row, _ = db.Query("SELECT Amount FROM shopcart.guest_product where CartID = ? and ProductID = ? and ColorID = ?", idCart, idPro, idColor)
	for row.Next() {
		row.Scan(&amount)
	}
	jsonData, err := json.Marshal(amount)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func createGuest(w http.ResponseWriter, r *http.Request) {
	cookie, error := r.Cookie("guest")
	if error == nil {
	} else {
		var count string
		row, _ := db.Query("SELECT count(*) FROM shopcart.guest;")
		defer row.Close()
		for row.Next() {
			row.Scan(&count)
		}
		co, _ := strconv.Atoi(count)
		db.Exec("INSERT INTO `shopcart`.`guest` (`Name`) VALUES ( ? );", strconv.Itoa(co+1))
		cookie = &http.Cookie{
			Name:    "guest",
			Value:   strconv.Itoa(co + 1),
			Expires: time.Now().Add(12 * time.Hour),
		}
		http.SetCookie(w, cookie)
		idPro := r.FormValue("idP")
		idColor := r.FormValue("idC")
		amount := r.FormValue("amount")
		if strings.Contains(idColor, "o") {
			idColor = "0"
		} else {
			row, _ := db.Query("SELECT id FROM shopcart.color where NameColor = ?", idColor)
			defer row.Close()
			for row.Next() {
				row.Scan(&idColor)
			}
		}
		db.Exec("INSERT INTO `shopcart`.`guest_product` (`CartID`, `ProductID`, `ColorID`, `Amount`) VALUES (?, ?, ?, ?);", co+1, idPro, idColor, amount)
	}
}

func png(w http.ResponseWriter, r *http.Request) {
	if checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		cookie, err := r.Cookie("guest")
		if err != nil {
			http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
			return
		}
		row1, _ := db.Query("SELECT id, ProductID FROM shopcart.guest_product where CartID = ?", cookie.Value)
		var list []string
		var listPro []string
		var guestId string
		var pro string
		for row1.Next() {
			row1.Scan(&guestId, &pro)
			list = append(list, guestId)
			listPro = append(listPro, pro)
		}
		index := r.FormValue("index")
		am := r.FormValue("am")
		index1, _ := strconv.Atoi(index)
		id := list[index1]
		idPro := listPro[index1]
		am1, _ := strconv.Atoi(am)
		var am2 int
		row, _ := db.Query("SELECT InStock FROM shopcart.products where idPro = ?", idPro)
		defer row.Close()
		for row.Next() {
			row.Scan(&am2)
		}
		if am1 <= am2 {
			db.Exec("UPDATE `shopcart`.`guest_product` SET `Amount` = ? WHERE (`id` = ?);", am1, id)
		} else {
			db.Exec("UPDATE `shopcart`.`guest_product` SET `Amount` = ? WHERE (`id` = ?);", am2, id)

		}
		jsonData, err := json.Marshal(am2)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

//********************seller***StoreID*************************************************************************

func loginSel(w http.ResponseWriter, r *http.Request) {
	if checkAccountSeller(w, r) {
		http.Redirect(w, r, "/homeSel", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "loginSel.html", nil)
		return
	}
	user := r.FormValue("username")
	password := r.FormValue("password")
	id := checkSellerByUserAndPassword(user, password)
	if !id {
		tpl.ExecuteTemplate(w, "loginSel.html", "UserName not found or password not exactly")
	} else {
		cookie := &http.Cookie{
			Name:    "accSel",
			Value:   user + "/" + "",
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/homeSel", http.StatusTemporaryRedirect)
	}
}

func homeSel(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}
	type send struct {
		ListCate []m.Category
	}
	tpl.ExecuteTemplate(w, "homeSel.html", send{ListCate: listCategories()})
}

func cusOrder(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		storeID := getStoreIdByUserName(w, r)
		status := r.FormValue("status")
		rows, _ := db.Query("SELECT id, OrderID FROM shopcart.orders_store where StoreID = ? and StatusID = ?", storeID, status)
		defer rows.Close()
		var order m.Order
		var listOder []m.Order
		for rows.Next() {
			rows.Scan(&order.OrderDetailID, &order.Id)
			order.Status.Id = status
			row, _ := db.Query("SELECT id, CustomerID, NameCustomer, Phone, Address, Saddress, Date, DeliveryID, PaymentID FROM shopcart.orders where id = ? ", order.Id)
			defer row.Close()
			for row.Next() {
				row.Scan(&order.Id, &order.CustomerID, &order.CustomerName, &order.Phone, &order.Address, &order.Saddress, &order.Date, &order.Delivery.Id, &order.Payment.Id)
				order.Delivery.Name = getNameDeliveryById(order.Delivery.Id)
				order.Payment.Name = getNamePaymentById(order.Payment.Id)
			}
			listOder = append(listOder, order)
		}
		listCate := listCategories()
		type send struct {
			ListOrder []m.Order
			ListCate  []m.Category
		}
		send1 := send{ListOrder: listOder, ListCate: listCate}
		tpl.ExecuteTemplate(w, "cusOrder.html", send1)
		return
	}
}

func detailCusOrder(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}

	if r.Method == "GET" {
		id := r.FormValue("id")
		row, _ := db.Query("SELECT id, CustomerID, NameCustomer, Phone, Address, Saddress, Date, DeliveryID, PaymentID FROM shopcart.orders where id = ? ", id)
		defer row.Close()
		var order m.Order
		for row.Next() {
			row.Scan(&order.Id, &order.CustomerID, &order.CustomerName, &order.Phone, &order.Address, &order.Saddress, &order.Date, &order.Delivery.Id, &order.Payment.Id)
			order.Delivery.Name = getNameDeliveryById(order.Delivery.Id)
			order.Payment.Name = getNamePaymentById(order.Payment.Id)
		}
		status := r.FormValue("status")

		storeId := getStoreIdByUserName(w, r)
		var od m.OrderDetail
		var listO []m.OrderDetail
		rows, _ := db.Query("SELECT id, OrderID, StoreID, StatusID FROM shopcart.orders_store where OrderID = ? and StatusID = ? and StoreID = ? ", id, status, storeId)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&od.Id, &od.OrderID, &od.Store.Id, &od.Status.Id)
			od.Store.Name = getNameStoreById(od.Store.Id)
			od.Status.Name = getNameStatusById(od.Status.Id)

			var po m.ProductOrder
			rows1, _ := db.Query("SELECT ProductID, ColorID, Amount, Price FROM shopcart.orders_detail where OrderSID = ?", od.Id)
			defer rows1.Close()
			for rows1.Next() {

				rows1.Scan(&po.ProductName, &po.ColorName, &po.Amount, &po.Price)
				row, _ = db.Query("SELECT NamePro FROM shopcart.products where idPro = ?", po.ProductName)
				for row.Next() {
					row.Scan(&po.ProductName)
				}
				if po.ColorName == "0" {
					po.ColorName = "None"
				} else {
					row, _ = db.Query("SELECT NameColor FROM shopcart.color where id = ?", po.ColorName)
					for row.Next() {
						row.Scan(&po.ColorName)
					}
				}
				od.ListPro = append(od.ListPro, po)
			}
			listO = append(listO, od)
		}

		listCate := listCategories()
		type send struct {
			Order       m.Order
			OrderDetail m.OrderDetail
			ListOrder   []m.OrderDetail
			ListCate    []m.Category
		}
		send1 := &send{Order: order, ListOrder: listO, ListCate: listCate, OrderDetail: od}
		tpl.ExecuteTemplate(w, "detailCusOrder.html", send1)
		return
	}
	id := r.FormValue("id")
	status, _ := strconv.Atoi(r.FormValue("status"))
	if status == 5 {
		db.Exec("UPDATE `shopcart`.`orders_store` SET `StatusID` = ? WHERE (`id` = ?)", status, id)
	} else {
		db.Exec("UPDATE `shopcart`.`orders_store` SET `StatusID` = ? WHERE (`id` = ?)", status+1, id)
	}
}

func subImage(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		id := r.FormValue("id")
		rows, _ := db.Query("SELECT id, Image FROM shopcart.product_slider where ProductID = ? ", id)
		defer rows.Close()
		var s []m.Image
		var s1 m.Image
		for rows.Next() {
			rows.Scan(&s1.Id, &s1.Url)
			s = append(s, s1)
		}
		type item struct {
			Len   string
			Image []m.Image
			Id    string
		}
		item1 := &item{Len: strconv.Itoa(len(s)), Image: s, Id: id}
		tpl.ExecuteTemplate(w, "sub_image.html", item1)
		return
	}
	id := r.FormValue("id")
	checkImg := r.FormValue("checkImg")
	if checkImg != "" {
		file, fileheader, _ := r.FormFile("img")
		defer file.Close()
		contentType := fileheader.Header["Content-Type"][0]
		var osFile *os.File
		defer osFile.Close()
		if contentType == "image/jpeg" || contentType == "image/png" {
			osFile, _ = ioutil.TempFile("static/image/", "*.jpg")
		} else {
			http.Redirect(w, r, "/images?id="+id+"#", http.StatusSeeOther)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		osFile.Write(fileBytes)
		s := osFile.Name()
		fmt.Println(osFile.Name())
		checkImg = "../" + s
	} else {
		http.Redirect(w, r, "/images?id="+id+"#", http.StatusSeeOther)
		return
	}
	db.Exec("INSERT INTO `shopcart`.`product_slider` (`ProductID`, `Image`) VALUES (?, ?);", id, checkImg)
	http.Redirect(w, r, "/images?id="+id+"#", http.StatusSeeOther)
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}
	id := r.FormValue("id")
	id1 := ""
	row, _ := db.Query("SELECT ProductID FROM shopcart.product_slider where id = ?", id)
	for row.Next() {
		row.Scan(&id1)
	}
	db.Exec("DELETE FROM `shopcart`.`product_slider` WHERE (`id` = ?);", id)
	http.Redirect(w, r, "/images?id="+id1+"#", http.StatusSeeOther)
}

func products(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}
	rows, _ := db.Query("SELECT idPro, NamePro, Price, Star, Discount, Image, `Date`, Sold, InStock, idCate FROM shopcart.products where StoreID = ?", getStoreIdByUserName(w, r))
	defer rows.Close()
	var p1 m.Product
	var list []m.Product
	for rows.Next() {
		rows.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Discount, &p1.Star, &p1.Image, &p1.Date, &p1.Sold, &p1.InStock, &p1.Category.IdCate)
		p1.Category.NameCate = getCategoryNameById(p1.Category.IdCate)
		list = append(list, p1)
	}
	tpl.ExecuteTemplate(w, "products.html", list)
}

func new(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		list := listCategories()
		tpl.ExecuteTemplate(w, "newProduct.html", list)
		return
	}
	name := r.FormValue("name")
	price := r.FormValue("price")
	discount := r.FormValue("discount")
	checkImg := r.FormValue("checkImg")
	inStock := r.FormValue("inStock")
	idCate := r.FormValue("cate")
	currentTime := time.Now()
	formatTime := currentTime.Format("2006/01/02 15:04:05")
	img := ""
	if checkImg != "" {
		file, fileheader, _ := r.FormFile("img")
		defer file.Close()
		contentType := fileheader.Header["Content-Type"][0]
		fmt.Println(fileheader.Header["Content-Type"][0])
		var osFile *os.File
		defer osFile.Close()
		if contentType == "image/jpeg" || contentType == "image/png" {
			osFile, _ = ioutil.TempFile("static/image/", "*.jpg")
		} else {
			tpl.ExecuteTemplate(w, "/home", nil)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		osFile.Write(fileBytes)
		s := osFile.Name()
		fmt.Println(osFile.Name())
		img = "../" + s
	}

	if discount == "" {
		discount = "0"
	}
	if inStock == "" {
		inStock = "0"
	}
	db.Exec("INSERT INTO `shopcart`.`products` (`NamePro`, `Price`, `Star`, `Discount`, `Image`, `Date`, `Sold`, `InStock`, `idCate`, `StoreID`) VALUES (?,?,?,?,?,?,?,?,?);\n", name, price, "0", discount, img, formatTime, "0", inStock, idCate, getStoreIdByUserName(w, r))
	http.Redirect(w, r, "/products", http.StatusTemporaryRedirect)
}

var p m.Product

func infoProduct(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		id := r.FormValue("id")
		if id == "" {
			http.Redirect(w, r, "/products", http.StatusTemporaryRedirect)
			return
		}
		row, err := db.Query("SELECT idPro,NamePro, Price, Discount, Star, Image,InStock,idCate  FROM shopcart.products where idPro = ?", id)
		defer row.Close()
		if err != nil {
			http.Redirect(w, r, "/products", http.StatusTemporaryRedirect)
			return
		}
		c := 0
		for row.Next() {
			row.Scan(&p.Id, &p.Name, &p.Price, &p.Discount, &p.Star, &p.Image, &p.InStock, &p.Category.IdCate)
			p.Category.NameCate = getCategoryNameById(p.Category.IdCate)
			c++
		}
		if c == 0 {
			http.Redirect(w, r, "/products", http.StatusTemporaryRedirect)
			return
		}
		rows, _ := db.Query("SELECT idCate, NameCate FROM shopcart.categories;")
		defer rows.Close()
		var c1 m.Category
		var list []m.Category
		for rows.Next() {
			rows.Scan(&c1.IdCate, &c1.NameCate)
			list = append(list, c1)
		}
		type give struct {
			Info     m.Product
			ListCate []m.Category
		}
		var in = give{Info: p, ListCate: list}
		tpl.ExecuteTemplate(w, "infoProduct.html", in)
		return
	}
	id := p.Id
	name := r.FormValue("name")
	price := r.FormValue("price")
	discount := r.FormValue("discount")
	checkImg := r.FormValue("checkImg")
	inStock := r.FormValue("inStock")
	idCate := r.FormValue("cate")
	img := ""
	if checkImg != "" && checkImg != p.Image {
		file, fileHeader, _ := r.FormFile("myFile")
		defer file.Close()
		contentType := fileHeader.Header["Content-Type"][0]
		var osFile *os.File
		defer osFile.Close()
		if contentType == "image/jpeg" || contentType == "image/png" {
			osFile, _ = ioutil.TempFile("static/image/", "*.jpg")
		} else {
			tpl.ExecuteTemplate(w, "infoProduct.html", p)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		osFile.Write(fileBytes)
		s := osFile.Name()
		fmt.Println(osFile.Name())
		img = "../" + s
	} else {
		img = p.Image
	}
	db.Exec("UPDATE `shopcart`.`products` SET `NamePro`= ?, `Price`=?, `Discount`=?, `Image`=?, `InStock`=?, `idCate`=?  WHERE (`idPro` = ? )", name, price, discount, img, inStock, idCate, id)
	http.Redirect(w, r, "/products", http.StatusTemporaryRedirect)
}

func infoSeller(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		id := getStoreIdByUserName(w, r)
		row, _ := db.Query("SELECT i.id ,s.User , i.Name, i.Image, i.Phone, i.Address, i.Email FROM shopcart.info_store i, shopcart.stores s where i.StoreID = s.id where s.id = ?", id)
		defer row.Close()
		var info m.Info
		for row.Next() {
			row.Scan(&info.Id, &info.UserName, &info.FullName, &info.Image, &info.PhoneNumber, &info.Address, &info.Email)
		}
		type send struct {
			I        m.Info
			ListCate []m.Category
		}
		tpl.ExecuteTemplate(w, "infoSeller.html", send{I: info, ListCate: listCategories()})
	}
	id := r.FormValue("id")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	address := r.FormValue("address")
	email := r.FormValue("email")
	var confirm string
	if !checkEmailSeller(email) {
		confirm = "email already exists"
	} else if !checkPhoneSeller(phone) {
		confirm = "phone already exists"
	} else {
		// confirm = "change successful"
		db.Exec("UPDATE `shopcart`.`info_store` SET `Name` = ?, `Phone` = ?, `Address` = ?, `Email` = ? WHERE (`id` = ?);", name, phone, address, email, id)
	}
	json, err := json.Marshal(confirm)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func statistical(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		type send struct {
			ListCate []m.Category
		}
		tpl.ExecuteTemplate(w, "statisticalSel.html", send{ListCate: listCategories()})
		return
	}
	form := r.FormValue("form")
	to := r.FormValue("to")
	idStore := getStoreIdByUserName(w, r)
	if to == "" {
		currentTime := time.Now()
		to = currentTime.Format("2006-01-02")
	}
	endDate, err := time.Parse("2006-01-02", to)
	if err != nil {
		fmt.Println(err)
		return
	}
	var startDate time.Time
	if form == "" {
		startDate = endDate.AddDate(0, 1, 0)
	} else {
		startDate, err = time.Parse("2006-01-02", form)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	form = startDate.Format("2006-01-02")

	type k struct {
		Period string `json:"period"`
		Value  int    `json:"value"`
	}
	var list []k
	currentDate := startDate
	for currentDate.Before(endDate.AddDate(0, 0, 1)) {
		list = append(list, k{Period: currentDate.Format("2006-01-02"), Value: 0})
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	rows, _ := db.Query("SELECT  date_format(DateOrder, '%Y-%m-%d') as date, count(*) as quatity FROM shopcart.orders_store where StoreID = ? and date_format(DateOrder, '%Y-%m-%d') >= ? and date_format(DateOrder, '%Y-%m-%d') <= ? group by date_format(DateOrder, '%Y-%m-%d')", idStore, form, to)
	defer rows.Close()
	var k1 k
	for rows.Next() {
		rows.Scan(&k1.Period, &k1.Value)
		for i := 0; i < len(list); i++ {
			if list[i].Period == k1.Period {
				list[i].Value = k1.Value
				break
			}
		}
	}
	json, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func stat30days(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusSeeOther)
		return
	}
	day := r.FormValue("day")
	idStore := getStoreIdByUserName(w, r)
	currentTime := time.Now()
	to := currentTime.Format("2006-01-02")
	endDate, err := time.Parse("2006-01-02", to)
	if err != nil {
		fmt.Println(err)
		return
	}
	var startDate time.Time
	d, err := strconv.Atoi(day)
	if err != nil {
		http.Redirect(w, r, "/homeSel", http.StatusSeeOther)
		return
	}
	startDate = endDate.AddDate(0, 0, -d)
	form := startDate.Format("2006-01-02")

	type k struct {
		Period string `json:"period"`
		Value  int    `json:"value"`
	}
	var list []k
	currentDate := startDate
	for currentDate.Before(endDate.AddDate(0, 0, 1)) {
		list = append(list, k{Period: currentDate.Format("2006-01-02"), Value: 0})
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	rows, _ := db.Query("SELECT  date_format(DateOrder, '%Y-%m-%d') as date, count(*) as quatity FROM shopcart.orders_store where StoreID = ? and date_format(DateOrder, '%Y-%m-%d') >= ? and date_format(DateOrder, '%Y-%m-%d') <= ? group by date_format(DateOrder, '%Y-%m-%d')", idStore, form, to)
	defer rows.Close()
	var k1 k
	for rows.Next() {
		rows.Scan(&k1.Period, &k1.Value)
		for i := 0; i < len(list); i++ {
			if list[i].Period == k1.Period {
				list[i].Value = k1.Value
				break
			}
		}
	}
	json, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func delete(w http.ResponseWriter, r *http.Request) {
	if !checkAccountSeller(w, r) {
		http.Redirect(w, r, "/loginSel", http.StatusTemporaryRedirect)
		return
	}
	db.Exec("DELETE FROM `shopcart`.`products` WHERE (`idPro` = ?)", p.Id)
	http.Redirect(w, r, "/products", http.StatusTemporaryRedirect)
}

func forgetPasswordForSeller(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "forgetPassword", nil)
	} else if r.Method == "POST" {
		email := r.FormValue("email")
		if checkEmail(email) {
			sendForgetPasswordForSeller(email)
			send := Result{T: "Đã gửi mật khẩu mới đến Email của bạn"}
			tpl.ExecuteTemplate(w, "forgetPassword", send)
		} else {
			send := Result{F: "Không tìm thấy Email"}
			tpl.ExecuteTemplate(w, "forgetPassword", send)
		}
	}
}

func registerSel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := r.FormValue("userName")
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		key := r.FormValue("key")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		address := r.FormValue("address")
		confirm := ""
		if checkUserSeller(userName) {
			confirm = "user da ton tai"
			jsonData, err := json.Marshal(confirm)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}
		if checkPhoneSeller(phone) {
			confirm = "phone da ton tai"
			jsonData, err := json.Marshal(confirm)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}
		if checkEmailSeller(email) {
			confirm = "email da ton tai"
			jsonData, err := json.Marshal(confirm)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}
		if !checkConfirmMailForSeller(email, key) {
			confirm = "nhap sai key"
			jsonData, err := json.Marshal(confirm)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}
		if password != confirmPassword {
			confirm = "hai password khong trung nhau"
			jsonData, err := json.Marshal(confirm)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}

		db.Exec("INSERT INTO `shopcart`.`stores` (`User`, `Password`, `Name`) VALUES (?, ?, ?);", userName, password, name)
		var id string
		row, _ := db.Query("SELECT id FROM shopcart.stores where User = ? and Password = ? ", userName, password)
		for row.Next() {
			row.Scan(&id)
		}
		image := "../static/img-test/alex-suprun-ZHvM3XIOHoE-unsplash.jpg"
		db.Exec("INSERT INTO `shopcart`.`info_store` (`StoreID`, `Name`, `Phone`, `Address`, `Email`, `Money`, `Image`) VALUES ( ?, ? ,?, ?, ?, ?, ?);", id, name, phone, address, email, 0, image)
	} else {
		tpl.ExecuteTemplate(w, "registerSel.html", nil)
	}
}

// *****************************mail seller*********************************************************
func sendForgetPasswordForSeller(email string) {
	var body bytes.Buffer
	t, err := template.ParseFiles("./sendCPSeller.html")
	t.Execute(&body, struct{ Name string }{Name: getNewPassword()})
	if err != nil {
		fmt.Println(err)
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "vietahh02@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", body.String())
	m.Attach("./img/LOGO.png")
	d := gomail.NewDialer("smtp.gmail.com", 587, "vietahh02@gmail.com", "iuam zbym wyaq euni")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func sendCodeForSeller(email, key string) {
	var body bytes.Buffer
	t, err := template.ParseFiles("./sendVerify.html")
	t.Execute(&body, struct{ Name string }{Name: key})
	if err != nil {
		fmt.Println(err)
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "vietahh02@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", body.String())
	m.Attach("./img/LOGO.png")
	d := gomail.NewDialer("smtp.gmail.com", 587, "vietahh02@gmail.com", "iuam zbym wyaq euni")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func verifySeller(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var confirm string
		email := r.FormValue("email")
		key := getKey()
		if !checkEmailSeller(email) {
			if checkNewEmail(email) {
				db.Exec("UPDATE `shopcart`.`verifynew` SET `Key` = ? WHERE (`Email` = ?);", key, email)
			} else {
				db.Exec("INSERT INTO `shopcart`.`verifynew` (`Email`, `Key`) VALUES (?, ?);", email, key)
			}
			sendCodeForSeller(email, key)
			confirm = "Y"
		} else {
			confirm = "N"
		}
		jsonData, err := json.Marshal(confirm)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func checkConfirmMailForSeller(email, key string) bool {
	row, _ := db.Query("SELECT * FROM shopcart.verifynew where Email = ? and KeyN = ?", email, key)
	return row != nil
}

func checkUserSeller(user string) bool {
	row, _ := db.Query("SELECT id FROM shopcart.stores where User = ?", user)
	defer row.Close()
	return row != nil
}

func checkPhoneSeller(phone string) bool {
	row, _ := db.Query("SELECT * FROM shopcart.info_store where Phone = ? ", phone)
	defer row.Close()
	return row != nil
}

func checkEmailSeller(email string) bool {
	row, _ := db.Query("SELECT * FROM shopcart.info_store where Email = ? ", email)
	defer row.Close()
	return row != nil
}

// *****************************admin*******************************************************************

func subColor(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		id := r.FormValue("id")
		rows, _ := db.Query("SELECT ProductID, Color, Discount_color FROM shopcart.product_color where ProductID = ? ", id)
		defer rows.Close()
		var color m.Color
		var colors []m.Color
		for rows.Next() {
			rows.Scan(&color.Id, &color.ColorUrl, &color.Discount)
			row, _ := db.Query("SELECT NameColor FROM shopcart.color where id = ?", color.ColorUrl)
			defer row.Close()
			for row.Next() {
				row.Scan(&color.ColorUrl)
			}
			colors = append(colors, color)
		}
		row, _ := db.Query("SELECT id, NameColor FROM shopcart.color;")
		var c []m.Color
		var c1 m.Color
		for row.Next() {
			row.Scan(&c1.Id, &c1.ColorUrl)
			c = append(c, c1)
		}
		type list struct {
			Id    string
			Color []m.Color
			C     []m.Color
		}
		list1 := list{Id: id, Color: colors, C: c}
		tpl.ExecuteTemplate(w, "sub_color.html", list1)
		return
	}
	id := r.FormValue("id")
	color := r.FormValue("color")
	discount := r.FormValue("discount")
	if discount == "" {
		discount = "0"
	}

	db.Exec("INSERT INTO `shopcart`.`product_color` (`ProductID`, `Color`, `Discount_color`) VALUES (?, ?, ?);", id, color, discount)
	http.Redirect(w, r, "/colors?id="+id+"#", http.StatusSeeOther)
}

func deleteColor(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	id := r.FormValue("id")
	id1 := ""
	row, _ := db.Query("SELECT ProductID FROM shopcart.product_color where id = ?", id)
	for row.Next() {
		row.Scan(&id1)
	}
	db.Exec("DELETE FROM `shopcart`.`product_color` WHERE (`id` = ?);", id)
	http.Redirect(w, r, "/colors?id="+id1+"#", http.StatusSeeOther)
}

func footerBanner(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	rows, _ := db.Query("SELECT idFB, Title, Content FROM shopcart.footer_banner;")
	defer rows.Close()
	var fb m.FooterBanner
	var list []m.FooterBanner
	for rows.Next() {
		rows.Scan(&fb.IdBanner, &fb.Title, &fb.Content)
		list = append(list, fb)
	}
	tpl.ExecuteTemplate(w, "footerBanner.html", list)
}

func newFooterBanner(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "newFooterBanner.html", nil)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	db.Exec("INSERT INTO `shopcart`.`footer_banner` (`Title`, `Content`) VALUES (?, ?)", title, content)
	http.Redirect(w, r, "/footerBanner", http.StatusTemporaryRedirect)
}

func deleteFooterBanner(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	id := r.FormValue("id")
	db.Exec("DELETE FROM `shopcart`.`footer_banner` WHERE (`idFB` = ?)", id)
	http.Redirect(w, r, "/footerBanner", http.StatusTemporaryRedirect)
}

func banner(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	rows, _ := db.Query("SELECT idBanner, Image FROM shopcart.banner;")
	defer rows.Close()
	var b m.Banner
	var list []m.Banner
	for rows.Next() {
		rows.Scan(&b.IdBannerImg, &b.Image)
		list = append(list, b)
	}
	tpl.ExecuteTemplate(w, "banner.html", list)
}

func newBanner(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "newBanner.html", nil)
		return
	}
	checkImg := r.FormValue("checkImg")
	if checkImg != "" {
		file, fileheader, _ := r.FormFile("img")
		defer file.Close()
		contentType := fileheader.Header["Content-Type"][0]
		var osFile *os.File
		defer osFile.Close()
		if contentType == "image/jpeg" || contentType == "image/png" {
			osFile, _ = ioutil.TempFile("static/img/", "*.jpg")
		} else {
			tpl.ExecuteTemplate(w, "newBanner.html", p)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		osFile.Write(fileBytes)
		s := osFile.Name()
		fmt.Println(osFile.Name())
		checkImg = "../" + s
	} else {
		tpl.ExecuteTemplate(w, "newBanner.html", nil)
		return
	}
	db.Exec("INSERT INTO `shopcart`.`banner` (`Image`) VALUES (?)", checkImg)
	http.Redirect(w, r, "/banner", http.StatusTemporaryRedirect)
}

func deleteBanner(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	id := r.FormValue("id")
	db.Exec("DELETE FROM `shopcart`.`banner` WHERE (`idBanner` = ? )", id)
	http.Redirect(w, r, "/banner", http.StatusTemporaryRedirect)
}

func homeAdmin(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	tpl.ExecuteTemplate(w, "admin.html", nil)
}

func categories(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	rows, _ := db.Query("SELECT idCate, NameCate FROM shopcart.categories;")
	defer rows.Close()
	var c m.Category
	var list []m.Category
	for rows.Next() {
		rows.Scan(&c.IdCate, &c.NameCate)
		list = append(list, c)
	}
	tpl.ExecuteTemplate(w, "categories.html", list)
}

var c m.Category

func infoCate(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		id := r.FormValue("id")
		if id == "" {
			http.Redirect(w, r, "/cate", http.StatusTemporaryRedirect)
			return
		}
		row, err := db.Query("SELECT NameCate FROM shopcart.categories where idCate = ?", id)
		defer row.Close()
		if err != nil {
			http.Redirect(w, r, "/cate", http.StatusTemporaryRedirect)
			return
		}
		var c1 m.Category
		c1.IdCate = id
		for row.Next() {
			row.Scan(&c1.NameCate)
		}
		c = c1
		tpl.ExecuteTemplate(w, "infoCate.html", c1)
		return
	}
	name := r.FormValue("name")
	db.Exec("UPDATE `shopcart`.`categories` SET `NameCate` = ? WHERE (`idCate` = ?)", name, c.IdCate)
	http.Redirect(w, r, "/cate", http.StatusTemporaryRedirect)
}

func newCate(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "newCate.html", nil)
		return
	}
	name := r.FormValue("name")
	db.Exec("INSERT INTO `shopcart`.`categories` (`NameCate`) VALUES (?)", name)
	http.Redirect(w, r, "/cate", http.StatusTemporaryRedirect)
}

func deleteCate(w http.ResponseWriter, r *http.Request) {
	if !checkAccountAdmin(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	db.Exec("DELETE FROM `shopcart`.`categories` WHERE (`idCate` = ?)", c.IdCate)
	http.Redirect(w, r, "/cate", http.StatusTemporaryRedirect)
}

//*****************************************************************

func checkSellerByUserAndPassword(user, password string) bool {
	row, _ := db.Query("SELECT id FROM shopcart.stores where User = ? and Password = ?", user, password)
	defer row.Close()
	var id string
	for row.Next() {
		row.Scan(&id)
	}
	if id == "" {
		return false
	} else {
		return true
	}
}

func checkCustomer(user string) bool {
	row, _ := db.Query("SELECT User FROM shopcart.customers")
	defer row.Close()
	var s string
	for row.Next() {
		row.Scan(&s)
		if s == user {
			return true
		}
	}
	return false
}

func checkPasswordCustomer(password, user string) bool {
	row, _ := db.Query("select Password from shopcart.customers where User = ?", user)
	defer row.Close()
	var s string
	for row.Next() {
		row.Scan(&s)
		if s == password {
			return true
		}
	}
	return false
}

func getCategoryNameById(id string) string {
	row, _ := db.Query("SELECT NameCate FROM shopcart.categories where idCate = ?", id)
	defer row.Close()
	var s string
	for row.Next() {
		row.Scan(&s)
	}
	return s
}

func checkAccountCustomer(w http.ResponseWriter, r *http.Request) bool {
	_, error := r.Cookie("acc")
	if error == nil {
		return true
	} else {
		return false
	}
}

func getCustomerIdByAccName(user string) string {
	row, _ := db.Query("SELECT id FROM shopcart.information where  User = ?", user)
	defer row.Close()
	var cusId string
	for row.Next() {
		row.Scan(&cusId)
	}
	return cusId
}

func getCustomerId(w http.ResponseWriter, r *http.Request) string {
	cookie, _ := r.Cookie("acc")
	a := strings.Split(cookie.Value, "/")
	return a[0]
}

func checkAccountSeller(w http.ResponseWriter, r *http.Request) bool {
	_, error := r.Cookie("accSel")
	if error == nil {
		return true
	} else {
		return false
	}
}

func getRandomUserName(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return string(base64.RawURLEncoding.EncodeToString(b)[:length])
}

func getStoreIdByUserName(w http.ResponseWriter, r *http.Request) string {
	cookie, _ := r.Cookie("accSel")
	a := strings.Split(cookie.Value, "/")
	row, _ := db.Query("SELECT id FROM shopcart.stores where  User = ?", a[0])
	defer row.Close()
	var cusId string
	for row.Next() {
		row.Scan(&cusId)
	}
	return cusId
}

func checkAccountAdmin(w http.ResponseWriter, r *http.Request) bool {
	_, error := r.Cookie("accAdmin")
	if error == nil {
		return true
	} else {
		return false
	}
}

func getStoreIdByProductId(id string) string {
	row, _ := db.Query("SELECT StoreID FROM shopcart.products where idPro = ?", id)
	defer row.Close()
	var idStore string
	for row.Next() {
		row.Scan(&idStore)
	}
	return idStore
}

func getInformationOrderByCustomerId(w http.ResponseWriter, r *http.Request) []m.Order {
	cusId := getCustomerId(w, r)
	var order m.Order
	var listOder []m.Order
	rows, _ := db.Query(" SELECT id, CustomerID, NameCustomer, Phone, Address, Saddress, Date, DeliveryID, PaymentID FROM shopcart.orders where CustomerID = ? ", cusId)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&order.Id, &order.CustomerID, &order.CustomerName, &order.Phone, &order.Address, &order.Saddress, &order.Date, &order.Delivery.Id, &order.Payment.Id)
		order.Delivery.Name = getNameDeliveryById(order.Delivery.Id)
		order.Payment.Name = getNamePaymentById(order.Payment.Id)
		listOder = append(listOder, order)
	}
	return listOder
}

func getNameDeliveryById(id string) string {
	row, _ := db.Query("SELECT Name FROM shopcart.delivery where id = ? ", id)
	defer row.Close()
	var name string
	for row.Next() {
		row.Scan(&name)
	}
	return name
}

func getNamePaymentById(id string) string {
	row, _ := db.Query("SELECT Name FROM shopcart.payment where id = ? ", id)
	defer row.Close()
	var name string
	for row.Next() {
		row.Scan(&name)
	}
	return name
}

func getNameStoreById(id string) string {
	row, _ := db.Query("SELECT Name FROM shopcart.stores where id = ?", id)
	defer row.Close()
	var name string
	for row.Next() {
		row.Scan(&name)
	}
	return name
}

func getNameStatusById(id string) string {
	row, _ := db.Query("SELECT Name FROM shopcart.status WHERE id = ? ", id)
	defer row.Close()
	var name string
	for row.Next() {
		row.Scan(&name)
	}
	return name
}

func listCategories() []m.Category {
	rows, _ := db.Query("SELECT idCate, NameCate FROM shopcart.categories")
	var c1 m.Category
	var ListCate []m.Category
	for rows.Next() {
		rows.Scan(&c1.IdCate, &c1.NameCate)
		ListCate = append(ListCate, c1)
	}
	return ListCate
}

func checkEmail(email string) bool {
	rows, _ := db.Query("SELECT Email FROM shopcart.information where Email = ? ", email)
	return rows != nil
}
func checkEmailRe(email string) bool {
	rows, _ := db.Query("SELECT Email FROM shopcart.information where AccID and AccID != '' and Email = ? ;", email)
	check := ""
	for rows.Next() {
		rows.Scan(&check)
	}
	return check != ""
	// return true
}

// ***********res*******
func upd(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	row, _ := db.Query("SELECT code, name FROM vietnam.districts where province_code = ?", code)
	var d m.Address
	var l []m.Address
	for row.Next() {
		row.Scan(&d.Code, &d.Name)
		l = append(l, d)
	}
	tpl.ExecuteTemplate(w, "upd.html", l)
}

func upw(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	row, _ := db.Query("SELECT code, name FROM vietnam.wards where district_code = ?", code)
	var wa m.Address
	var l []m.Address
	for row.Next() {
		row.Scan(&wa.Code, &wa.Name)
		l = append(l, wa)
	}
	tpl.ExecuteTemplate(w, "upw.html", l)
}

func json_cart(w http.ResponseWriter, r *http.Request) {
	if !checkAccountCustomer(w, r) {
		http.Redirect(w, r, "/loginCus", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "GET" {
		cusId := getCustomerId(w, r)
		row, _ := db.Query("SELECT id FROM shopcart.carts where CustomerID = ?", cusId)
		defer row.Close()
		for row.Next() {
			row.Scan(&cusId)
		}
		var c m.Cart
		type ty struct {
			P  m.Product
			Co m.Color
			Ca m.Cart
		}
		var list []ty
		var ty1 ty
		row1, _ := db.Query("SELECT id, CartID, ProductID, ColorID, Amount FROM shopcart.carts_product where CartID = ?", cusId)
		if row1 == nil {
			tpl.ExecuteTemplate(w, "/", nil)
			return
		}
		for row1.Next() {
			row1.Scan(&c.Id, &c.CartID, &c.ProductID, &c.ColorID, &c.Amount)
			ty1.Ca = c
			var p1 m.Product
			row, _ = db.Query("SELECT idPro, NamePro, Price, Star, Discount, Image, `Date`, Sold, InStock, idCate FROM shopcart.products where idPro = ?", c.ProductID)
			for row.Next() {
				row.Scan(&p1.Id, &p1.Name, &p1.Price, &p1.Star, &p1.Discount, &p1.Image, &p1.Date, &p1.Sold, &p1.InStock, &p1.Category.IdCate)
			}
			ty1.P = p1
			row, _ = db.Query("SELECT ProductID, Color, Discount_color FROM shopcart.product_color where ProductID = ? and Color = ?", c.ProductID, c.ColorID)
			var c1 m.Color
			for row.Next() {
				row.Scan(&c1.ProductId, &c1.ColorUrl, &c1.Discount)
				row, _ := db.Query("SELECT NameColor FROM shopcart.color where id = ?", c1.ColorUrl)
				defer row.Close()
				for row.Next() {
					row.Scan(&c1.ColorUrl)
				}
			}
			ty1.Co = c1
			list = append(list, ty1)
		}
		rows, _ := db.Query("SELECT idCate, NameCate FROM shopcart.categories")
		var c1 m.Category
		var listCate []m.Category
		for rows.Next() {
			rows.Scan(&c1.IdCate, &c1.NameCate)
			listCate = append(listCate, c1)
		}
		rows, _ = db.Query("SELECT code, name FROM vietnam.provinces;")
		var add m.Address
		var listProvince []m.Address
		for rows.Next() {
			rows.Scan(&add.Code, &add.Name)
			listProvince = append(listProvince, add)
		}
		type send struct {
			ListCate     []m.Category
			ListCart     []ty
			ListProvince []m.Address
		}
		send1 := send{ListCate: listCate, ListCart: list, ListProvince: listProvince}
		jsonData, err := json.Marshal(send1)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}
}
