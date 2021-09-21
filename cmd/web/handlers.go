package main

import (
	"YaroslavSivash/web-site/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	notes, err := app.notes.Latest()
	if err!= nil {
		app.serverError(w, err)
		return
	}

	for _, note := range notes {
		fmt.Fprintf(w, "%v\n", note)
	}

	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"ui/html/footer.partial.tmpl",
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.errorLog.Println(err.Error())
	//	app.serverError(w, err)
	//	return
	//}
	//err = ts.Execute(w, nil)
	//if err != nil {
	//	app.errorLog.Println(err.Error())
	//	app.serverError(w, err)
	//	return
	//}
}

func (app *application) showNotes(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	date , err := app.notes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
// Отображаем весь вывод на странице.
fmt.Fprintf(w, "%v", date)
}


func (app *application) createNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	// Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

	// Передаем данные в метод SnippetModel.Insert(), получая обратно
	// ID только что созданной записи в базу данных.
	id, err := app.notes.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправляем пользователя на соответствующую страницу заметки.
	http.Redirect(w, r, fmt.Sprintf("/show?id=%d", id), http.StatusSeeOther)
}

