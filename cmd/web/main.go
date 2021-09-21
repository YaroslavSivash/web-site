package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	_ "github.com/go-sql-driver/mysql" // Новый импорт
)

// Создаем структуру `application` для хранения зависимостей всего веб-приложения.
// Пока, что мы добавим поля только для двух логгеров, но
// мы будем расширять данную структуру по мере усложнения приложения.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// Создаем новый флаг командной строки, значение по умолчанию: ":4000".
	// Добавляем небольшую справку, объясняющая, что содержит данный флаг.
	// Значение флага будет сохранено в переменной addr.
	addr := flag.String("addr", ":8000", "Сетевой адрес HTTP")
	dsn := flag.String("dsn", "web:123@/note?parseTime=true", "Название MySQL источника данных")
	// Мы вызываем функцию flag.Parse() для извлечения флага из командной строки.
	// Она считывает значение флага из командной строки и присваивает его содержимое
	// переменной. Вам нужно вызвать ее *до* использования переменной addr
	// иначе она всегда будет содержать значение по умолчанию ":4000".
	// Если есть ошибки во время извлечения данных - приложение будет остановлено.
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// Создаем логгер для записи сообщений об ошибках таким же образом, но используем stderr как
	// место для записи и используем флаг log.Lshortfile для включения в лог
	// названия файла и номера строки где обнаружилась ошибка.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Инициализируем новую структуру http.Server. Мы устанавливаем поля Addr и Handler, так
	// что сервер использует тот же сетевой адрес и маршруты, что и раньше, и назначаем
	// поле ErrorLog, чтобы сервер использовал наш логгер
	// при возникновении проблем.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()
	// Значение, возвращаемое функцией flag.String(), является указателем на значение
	// из флага, а не самим значением. Нам нужно убрать ссылку на указатель
	// то есть перед использованием добавьте к нему префикс *. Обратите внимание, что мы используем
	// функцию log.Printf() для записи логов в журнал работы нашего приложения.
	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

type neuteredFileSystem struct {
	fc http.FileSystem
}

func (nfs neuteredFileSystem) Open (path string) (http.File, error) {
	f, err := nfs.fc.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if s.IsDir(){
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fc.Open(index); err != nil {
			closeErr:=f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}

func openDB (dsn string) (*sql.DB, error) {
	db, err:= sql.Open("mysql", dsn)
	if err != nil{
		return nil, err
	}
	if err = db.Ping(); err != nil{
		return nil, err
	}
	return db, err
}