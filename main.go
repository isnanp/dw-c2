package main

import (
	"context"
	"fmt"
	"log"
	connection "myapp/connetion"
	"myapp/middleware"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	e := echo.New()
	connection.DatabaseConnect()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session_key"))))

	//static files from public directory
	e.Static("/public", "public")
	e.Static("/uploads", "uploads")

	//routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/testimoni", testimoni)
	e.GET("/detail-project/:id", detailProject)
	e.GET("/add-project", addProject)
	e.GET("/update-project/:id", updateProject)
	e.POST("/add-project", middleware.UploadFile(postProject))
	e.POST("/delete-project/:id", deleteProject)
	e.POST("/submit-Update-Project/:id", middleware.UploadFile(submitUpdateProject))
	// login and register
	e.GET("/login-page", loginPage)
	e.POST("/login", login)
	e.GET("/register-page", registerPage)
	e.POST("/register", register)
	e.POST("/logout", logout)

	e.Logger.Fatal(e.Start("localhost:5000"))
}

type Project struct {
	Id              int
	Author          string
	Image           string
	Title           string
	StartDate       time.Time
	EndDate         time.Time
	Desc            string
	Node            bool
	Golang          bool
	React           bool
	Java            bool
	FormatStartDate string
	FormatEndDate   string
	Duration        string
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type SessionsData struct {
	IsLogin bool
	Name    string
}

var userData = SessionsData{}

func home(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	authorId := sess.Values["id"]

	var result []Project

	if sess.Values["isLogin"] != true {
		data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, title, start_date, end_date, description, image, node, golang, react, java, tb_users.name AS author FROM tb_project JOIN tb_users ON tb_project.author_id = tb_users.id")

		for data.Next() {
			var each = Project{}

			err := data.Scan(&each.Id, &each.Title, &each.StartDate, &each.EndDate, &each.Desc, &each.Image, &each.Node, &each.Golang, &each.React, &each.Java, &each.Author)

			if err != nil {
				fmt.Println(err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
			}

			each.FormatStartDate = each.StartDate.Format("2 January 2006")
			each.FormatEndDate = each.EndDate.Format("2 January 2006")
			each.Duration = timeDuration(each.StartDate, each.EndDate)

			result = append(result, each)
		}
	} else {
		data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, title, start_date, end_date, description, image, node, golang, react, java, tb_users.name AS author FROM tb_project JOIN tb_users ON tb_project.author_id = tb_users.id  WHERE tb_users.id=$1", authorId)
		for data.Next() {
			var each = Project{}

			err := data.Scan(&each.Id, &each.Title, &each.StartDate, &each.EndDate, &each.Desc, &each.Image, &each.Node, &each.Golang, &each.React, &each.Java, &each.Author)

			if err != nil {
				fmt.Println(err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
			}

			each.FormatStartDate = each.StartDate.Format("2 January 2006")
			each.FormatEndDate = each.EndDate.Format("2 January 2006")
			each.Duration = timeDuration(each.StartDate, each.EndDate)

			result = append(result, each)
		}
	}

	var datas = map[string]interface{}{
		"Projects":     result,
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"UserData":     userData,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), datas)
}

func contact(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	var datas = map[string]interface{}{
		"UserData": userData,
	}

	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), datas)
}

func testimoni(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	var datas = map[string]interface{}{
		"UserData": userData,
	}

	var tmpl, err = template.ParseFiles("views/myTestimonial.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), datas)
}

func addProject(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	var datas = map[string]interface{}{
		"UserData": userData,
	}

	var tmpl, err = template.ParseFiles("views/addProject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), datas)
}

func postProject(c echo.Context) error {
	sess, _ := session.Get("session", c)
	author := sess.Values["id"]

	title := c.FormValue("nameProject")
	description := c.FormValue("descriptionProject")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	node := (c.FormValue("nodeJs") == "on")
	golang := (c.FormValue("golang") == "on")
	react := (c.FormValue("react") == "on")
	java := (c.FormValue("java") == "on")
	images := c.Get("dataFile")

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (author_id, image, title, start_date, end_date, description, node, golang, react, java) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", author, images, title, startDate, endDate, description, node, golang, react, java)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT tb_project.id, title, start_date, end_date, description, image, node, golang, react, java, tb_users.name AS author FROM tb_project JOIN tb_users ON tb_project.author_id = tb_users.id WHERE tb_project.id=$1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Desc, &ProjectDetail.Image, &ProjectDetail.Node, &ProjectDetail.Golang, &ProjectDetail.React, &ProjectDetail.Java, &ProjectDetail.Author)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	ProjectDetail.FormatStartDate = ProjectDetail.StartDate.Format("2 January 2006")
	ProjectDetail.FormatEndDate = ProjectDetail.EndDate.Format("2 January 2006")
	ProjectDetail.Duration = timeDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	var datas = map[string]interface{}{
		"Projects": ProjectDetail,
		"UserData": userData,
	}
	var tmpl, errTmpl = template.ParseFiles("views/projectDetail.html")

	if errTmpl != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), datas)
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func updateProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description, node, golang, react, java FROM tb_project WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Desc, &ProjectDetail.Node, &ProjectDetail.Golang, &ProjectDetail.React, &ProjectDetail.Java)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	ProjectDetail.FormatStartDate = ProjectDetail.StartDate.Format("2006-01-02")
	ProjectDetail.FormatEndDate = ProjectDetail.EndDate.Format("2006-01-02")
	ProjectDetail.Duration = timeDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	var projects = map[string]interface{}{
		"Projects": ProjectDetail,
		"UserData": userData,
	}
	var tmpl, errTmpl = template.ParseFiles("views/updateProject.html")

	if errTmpl != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func submitUpdateProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	title := c.FormValue("nameProject")
	description := c.FormValue("descriptionProject")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	node := (c.FormValue("nodeJs") == "on")
	golang := (c.FormValue("golang") == "on")
	react := (c.FormValue("react") == "on")
	java := (c.FormValue("java") == "on")
	images := c.Get("dataFile")

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title=$1, start_date=$2, end_date=$3, description=$4, node=$5, golang=$6, react=$7, java=$8, image=$9 WHERE id=$10", title, startDate, endDate, description, node, golang, react, java, images, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func timeDuration(start, end time.Time) string {

	diff := end.Sub(start)
	days := int(diff.Hours() / 24)
	weeks := days / 7
	months := days / 30

	if months >= 12 {
		return strconv.Itoa(months/12) + " tahun"
	} else if months > 0 {
		return strconv.Itoa(months) + " bulan"
	} else if weeks > 0 {
		return strconv.Itoa(weeks) + " minggu"
	} else {
		return strconv.Itoa(days) + " hari"
	}
}

func loginPage(c echo.Context) error {
	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	email := c.FormValue("Email-login")
	password := c.FormValue("Password-login")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_users WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return redirectWithMessage(c, "Email incorrect", false, "/login-page")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "Password incorrect", false, "/login-page")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 hours
	sess.Values["message"] = "Login successful"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["id"] = user.Id
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func registerPage(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func register(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := c.FormValue("Name-registration")
	email := c.FormValue("Email-registration")
	password := c.FormValue("Password-registration")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_users(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		return redirectWithMessage(c, "Register Failed, please try again", false, "/register-page")
	}

	return redirectWithMessage(c, "Register Success", true, "/login-page")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, path)
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}
