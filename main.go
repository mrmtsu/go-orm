package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/mrmtsu/go-board/my"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// db variable.
var dbDriver = "mysql"
var dbName = "data_splite3"

// session variable.
var sesName = "ytboard-session"
var cs = sessions.NewCookieStore([]byte("secret-key-1234"))

// login check.
func checkLogin(w http.ResponseWriter, r *http.Request) *my.User {
	ses, _ := cs.Get(r, sesName)
	if ses.Values["login"] == nil || !ses.Values["login"].(bool) {
		http.Redirect(w, r, "/login", 302)
	}
	ac := ""
	if ses.Values["account"] != nil {
		ac = ses.Values["account"].(string)
	}

	var user my.User
	dsn := "user:pass@tcp(127.0.0.1:3306)/database_namecharset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.Where("account = ?", ac).First(&user)

	return &user
}

// Template for no-template.
func notemp() *template.Template {
	tmp, _ := template.New("index").Parse("NO PAGE.")
	return tmp
}

// get target Template.
func page(fname string) *template.Template {
	tmps, _ := template.ParseFiles("templates/"+fname+".html",
		"templates/head.html", "templates/foot.html")
	return tmps
}

// top page handler.
func index(w http.ResponseWriter, r *http.Request) {
	user := checkLogin(w, r)

	dsn := "user:pass@tcp(127.0.0.1:3306)/database_namecharset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	var pl []my.Post
	db.Where("group_id > 0").Order("created_at desc").Limit(10).Find(&pl)
	var gl []my.Group
	db.Order("created_at desc").Limit(10).Find(&gl)

	item := struct {
		Title   string
		Message string
		Name    string
		Account string
		Plist   []my.Post
		Glist   []my.Group
	}{
		Title:   "Index",
		Message: "This is Top page.",
		Name:    user.Name,
		Account: user.Message,
		Plist:   pl,
		Glist:   gl,
	}
	er := page("index").Execute(w, item)
	if er != nil {
		log.Fatal(er)
	}
}

// top page handler.
func post(w http.ResponseWriter, r *http.Request) {
	user := checkLogin(w, r)

	pid := r.FormValue("pid")
	dsn := "user:pass@tcp(127.0.0.1:3306)/database_namecharset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if r.Method == "POST" {
		msg := r.PostFormValue("message")
		pId, _ := strconv.Atoi(pid)
		cmt := my.Comment{
			UserId:  int(user.Model.ID),
			PostId:  pId,
			Message: msg,
		}
		db.Create(&cmt)
	}

	var pst my.Post
	var cmts []my.CommentJoin

	db.Where("id = ?", pid).First(&pst)
	db.Table("comments").Select("comments.*, users.id, users.name").Joins("join users on users.id = ?", pid).Order("created_at desc").Find(&cmts)

	item := struct {
		Title   string
		Message string
		Name    string
		Account string
		Post    my.Post
		Clist   []my.CommentJoin
	}{
		Title:   "Post",
		Message: "Post id=" + pid,
		Name:    user.Name,
		Account: user.Account,
		Post:    pst,
		Clist:   cmts,
	}
	er := page("post").Execute(w, item)
	if er != nil {
		log.Fatal(er)
	}
}

// home handler
func home(w http.ResponseWriter, r *http.Request) {
	user := checkLogin(w, r)

	dsn := "user:pass@tcp(127.0.0.1:3306)/database_namecharset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if r.Method == "POST" {
		switch r.PostFormValue("form") {
		case "post":
			ad := r.PostFormValue("address")
			ad = strings.TrimSpace(ad)
			if strings.HasPrefix(ad, "https://youtu.be/") {
				ad = strings.TrimPrefix(ad, "https://youtu.be/")
			}

			pt := my.Post{
				UserId:  int(user.Model.ID),
				Address: ad,
				Message: r.PostFormValue("message"),
			}
			db.Create(&pt)
		case "group":
			gp := my.Group{
				UserId:  int(user.Model.ID),
				Name:    r.PostFormValue("name"),
				Message: r.PostFormValue("message"),
			}
			db.Create(&gp)
		}
	}

	var pts []my.Post
	var gps []my.Group

	db.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&pts)
	db.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&gps)

	itm := struct {
		Title   string
		Message string
		Name    string
		Account string
		Plist   []my.Post
		Glist   []my.Group
	}{
		Title:   "Home",
		Message: "User account=\"" + user.Account + "\"",
		Name:    user.Name,
		Account: user.Account,
		Plist:   pts,
		Glist:   gps,
	}
	er := page("home").Execute(w, itm)
	if er != nil {
		log.Fatal(er)
	}
}

// group handler.
func group(w http.ResponseWriter, r *http.Request) {
	user := checkLogin(w, r)

	gid := r.FormValue("gid")
	dsn := "user:pass@tcp(127.0.0.1:3306)/database_namecharset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if r.Method == "POST" {
		ad := rq.PostFormValue("address")
		ad = strings.TrimSpace(ad)
		if strings.HasPrefix(ad, "https://youtu.be/") {
			ad = strings.TrimPrefix(ad, "https://youyu.be/")
		}
		gId, _ := strconv.Atoi(gid)
		pt := my.Post{
			UserId:  int(user.Model.ID),
			Address: ad,
			Message: r.PostFormValue("message"),
			GroupId: gId,
		}
		db.Create(&pt)
	}

	var grp my.Group
	var pts []my.Post

	db.Where("id=?", gid).First(&grp)
	db.Order("created_at desc").Model(&grp).Related(&pts)

	itm := struct {
		Title   string
		Message string
		Name    string
		Account string
		Group   my.Group
		Plist   []my.Post
	}{
		Title:   "Group",
		Message: "Group id=" + gid,
		Name:    user.Name,
		Account: user.Account,
		Group:   grp,
		Plist:   pts,
	}
	er := page("group").Execute(w, itm)
	if er != nil {
		log.Fatal(er)
	}
}

// login handler.
func login(w http.ResponseWriter, r *http.Request) {
	item := struct {
		Title   string
		Message string
		Account string
	}{
		Title:   "Login",
		Message: "type your account & password:",
		Account: "",
	}

	if r.Method == "GET" {
		er := page("login").Execute(w, item)
		if er != nil {
			log.Fatal(er)
		}
		return
	}
	if r.Method == "POST" {
		dsn := "user:pass@tcp(127.0.0.1:3306)/database_namecharset=utf8mb4&parseTime=True&loc=Local"
		db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		usr := r.PostFormValue("account")
		pass := r.PostFormValue("pass")
		item.Account = usr

		// check account and password
		var re int64
		var user my.User

		db.Where("account = ? and password = ?", usr, pass).Find(&user).Count(&re)

		if re <= 0 {
			item.Message = "Wrong account or password."
			page("login").Execute(w, item)
			return
		}

		// logined.
		ses, _ := cs.Get(r, sesName)
		ses.Values["login"] = true
		ses.Values["account"] = usr
		ses.Values["name"] = user.Name
		ses.Save(r, w)
		http.Redirect(w, r, "/", 302)
	}
	er := page("login").Execute(w, item)
	if er != nil {
		log.Fatal(er)
	}
}

// logout handler.
func logout(w http.ResponseWriter, r *http.Request) {
	ses, _ := cs.Get(r, sesName)
	ses.Values["login"] = nil
	ses.Values["account"] = nil
	ses.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}

// main program.
func main() {
	// index handling.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index(w, r)
	})
	// home handling.
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		home(w, r)
	})
	// post handling.
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		post(w, r)
	})
	// post handling
	http.HandleFunc("/group", func(w http.ResponseWriter, r *http.Request) {
		group(w, r)
	})
	// login handling
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r)
	})
	// logout handling
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logout(w, r)
	})
	http.ListenAndServe("", nil)
}
