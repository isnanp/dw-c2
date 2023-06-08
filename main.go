package main

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	//static files from public directory
	e.Static("/public", "public")

	//routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/testimoni", testimoni)
	e.GET("/detail-project/:id", detailProject)
	e.GET("/add-project", addProject)
	e.POST("/add-project", postProject)
	e.POST("/delete-project/:id", deleteProject)

	e.Logger.Fatal(e.Start("localhost:5000"))
}

type Project struct {
	Title     string
	Desc      string
	StartDate string
	EndDate   string
}

var projectData = []Project{
	{
		Title:     "Coba dulu aja ga sih",
		Desc:      "Coba dulu aja sih aja",
		StartDate: "25 maret 2018",
		EndDate:   "25 april 2018",
	},
	{
		Title:     "Coba coba dulu aja sih",
		Desc:      "Coba dulu aja sih aja",
		StartDate: "25 maret 2018",
		EndDate:   "25 april 2018",
	},
}

func home(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var projects = map[string]interface{}{
		"Projects": projectData,
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

	newData := Project{
		Title:     title,
		Desc:      description,
		StartDate: startDate,
		EndDate:   endDate,
	}

	projectData = append(projectData, newData)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	ProjectDetail := Project{}

	for i, p := range projectData {
		if id == i {
			ProjectDetail = Project{
				Title:     p.Title,
				Desc:      p.Desc,
				StartDate: p.StartDate,
				EndDate:   p.EndDate,
			}
		}
	}

	data := map[string]interface{}{
		"Projects": ProjectDetail,
	}

	var tmpl, err = template.ParseFiles("views/projectDetail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	projectData = append(projectData[:id], projectData[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}
