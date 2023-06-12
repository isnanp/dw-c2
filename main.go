package main

import (
	"context"
	"fmt"
	connection "myapp/connetion"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	connection.DatabaseConnect()

	//static files from public directory
	e.Static("/public", "public")

	//routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/testimoni", testimoni)
	e.GET("/detail-project/:id", detailProject)
	e.GET("/add-project", addProject)
	e.GET("/update-project/:id", updateProject)
	e.POST("/add-project", postProject)
	e.POST("/delete-project/:id", deleteProject)
	e.POST("/submit-Update-Project/:id", submitUpdateProject)

	e.Logger.Fatal(e.Start("localhost:5000"))
}

type Project struct {
	Id              int
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

func home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, description, node, golang, react, java FROM tb_project")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.Title, &each.StartDate, &each.EndDate, &each.Desc, &each.Node, &each.Golang, &each.React, &each.Java)

		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		}

		each.FormatStartDate = each.StartDate.Format("2 January 2006")
		each.FormatEndDate = each.EndDate.Format("2 January 2006")
		each.Duration = timeDuration(each.StartDate, each.EndDate)

		result = append(result, each)
	}

	var projects = map[string]interface{}{
		"Projects": result,
	}

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func contact(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func testimoni(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/myTestimonial.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func addProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/addProject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func postProject(c echo.Context) error {
	title := c.FormValue("nameProject")
	description := c.FormValue("descriptionProject")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	node := (c.FormValue("nodeJs") == "on")
	golang := (c.FormValue("golang") == "on")
	react := (c.FormValue("react") == "on")
	java := (c.FormValue("java") == "on")

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (title, start_date, end_date, description, node, golang, react, java) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", title, startDate, endDate, description, node, golang, react, java)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description, node, golang, react, java FROM tb_project WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Desc, &ProjectDetail.Node, &ProjectDetail.Golang, &ProjectDetail.React, &ProjectDetail.Java)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	ProjectDetail.FormatStartDate = ProjectDetail.StartDate.Format("2 January 2006")
	ProjectDetail.FormatEndDate = ProjectDetail.EndDate.Format("2 January 2006")
	ProjectDetail.Duration = timeDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	var projects = map[string]interface{}{
		"Projects": ProjectDetail,
	}
	var tmpl, errTmpl = template.ParseFiles("views/projectDetail.html")

	if errTmpl != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
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

	var projects = map[string]interface{}{
		"Projects": ProjectDetail,
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

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title=$1, start_date=$2, end_date=$3, description=$4, node=$5, golang=$6, react=$7, java=$8 WHERE id=$9", title, startDate, endDate, description, node, golang, react, java, id)

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
