package receipt

import (
	json2 "encoding/json"
	"fmt"
	"inventoryservice/cors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const receiptPath = "receipts"

func handleReceipts(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		receipts, err := GetReceipts()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		json, err := json2.Marshal(receipts)
		if err != nil {
			log.Fatal(err)
		}
		_, err = writer.Write(json)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		request.ParseMultipartForm(5 << 20)
		file, handler, err := request.FormFile("receipt")
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		io.Copy(f, file)
		writer.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func handleDownload(writer http.ResponseWriter, request *http.Request) {
	urlPathSegments := strings.Split(request.URL.Path, fmt.Sprintf("%s/", receiptPath))
	if len(urlPathSegments[1:]) > 1 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)
	stat, err := file.Stat()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileSize := strconv.FormatInt(stat.Size(), 10)
	writer.Header().Set("Content-Disposition", "attachment="+fileName)
	writer.Header().Set("Content-Type", fileContentType)
	writer.Header().Set("Content-Length", fileSize)
	file.Seek(0, 0)
	io.Copy(writer, file)
}

func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.MiddlewareHandler(receiptHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.MiddlewareHandler(downloadHandler))

}
