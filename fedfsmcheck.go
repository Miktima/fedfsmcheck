package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getHtmlPage(url, userAgent string) ([]byte, error) {
	// функция получения ресурсов по указанному адресу url с использованием User-Agent
	// возвращает загруженный HTML контент
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Cannot create new request  %s, error: %v\n", url, err)
		return nil, err
	}

	// с User-agent по умолчанию контент не отдается, используем свой
	req.Header.Set("User-Agent", userAgent)

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error with GET request: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	// Получаем контент и возвращаем его, как результат работы функции
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error ReadAll")
		return nil, err
	}
	return body, nil
}

func getList(body []byte, tag string) []string {
	// Функция получения списка из html контента
	// Список достается из тега tag
	// Возвращает список
	tkn := html.NewTokenizer(bytes.NewReader(body))
	depth := 0
	var flist []string
	block := ""
	errorCode := false

	// Проходим по всему дереву тегов (пока не встретится html.ErrorToken)
	for !errorCode {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			errorCode = true
		case html.TextToken:
			if depth > 0 {
				block += string(tkn.Text()) // Если внутри нужного тега, забираем текст из блока
			}
		case html.StartTagToken, html.EndTagToken:
			tn, tAttr := tkn.TagName()
			if string(tn) == "div" { // выбираем нужный tag
				if tAttr {
					key, attr, _ := tkn.TagAttr()
					if tt == html.StartTagToken && string(key) == "id" && string(attr) == tag {
						depth++ // нужный тег открывается
					}
				} else if tt == html.EndTagToken && depth >= 1 {
					depth--
					flist = append(flist, block) // Когда блок закрывается, добавляем текст из блока в список
					block = ""
				}
			}
		}
	}
	return flist
}

func main() {
	var urlList string
	var userAgent string

	// Ключи для командной строки
	flag.StringVar(&urlList, "url", "https://fedsfm.ru/documents/terrorists-catalog-portal-add", "The URL to lists")
	flag.StringVar(&userAgent, "uagent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit", "User Agent")

	flag.Parse()

	path, _ := os.Executable()
	path = path[:strings.LastIndex(path, "/")+1]
	fmt.Println("Path: ", path)

	// Читаем файлы со списками
	if _, err := os.Stat(path + "/ul_file.txt"); err == nil {
		// Open our jsonFile
		// byteValue_ul, err := os.ReadFile(path + "/ul_file.txt")
		_, err := os.ReadFile(path + "/ul_file.txt")
		// if we os.ReadFile returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
	}

	if _, err := os.Stat(path + "/fl_file.txt"); err == nil {
		// Open our jsonFile
		_, err := os.ReadFile(path + "/fl_file.txt")
		// byteValue_fl, err := os.ReadFile(path + "/fl_file.txt")
		// if we os.ReadFile returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
	}

	// var ul_list []string
	var fl_list []string

	body, err := getHtmlPage(urlList, userAgent)
	if err != nil {
		fmt.Printf("Error getHtmlPage - %v\n", err)
	}
	// Получаем заголовок и текст статьи
	fl_list = getList(body, "russianFL")
	fmt.Println(fl_list)

}
