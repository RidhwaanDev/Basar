package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ServeFile(writer http.ResponseWriter, request *http.Request, name string) {
	fmt.Println("func ServeFile in file.go")

	// Reading header info from the opened file, this will be used for response header "Content-Type"
	// the first 512 bytes is the header.
	file, err := os.Open(name)
	defer file.Close()
	catch(err)
	fileHeader := make([]byte, 512)
	_, erro := file.Read(fileHeader) // File offset is now len(fileHeader)
	fileType := http.DetectContentType(fileHeader)

	if erro != nil {
		catch(erro)
	}

	// Get file info which we will use for the response headers "Content-Disposition" and "Content-Length"
	fileInfo, err := file.Stat()
	if err != nil {
		catch(err)
	}
	writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s""`, fileInfo.Name()))
	writer.Header().Set("Content-Type", fileType)
	requestRange := request.Header.Get("range")
	fmt.Printf("request range: %s\n", requestRange)

	if requestRange == "" {
		writer.Header().Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))
		file.Seek(0, 0)
		fmt.Println("filename path: ", filepath.Ext(fileInfo.Name()))
		io.Copy(writer, file)

	}

	requestRange = requestRange[6:] // Strip the "bytes=", left over is now "begin-end"
	splitRange := strings.Split(requestRange, "-")

	if len(splitRange) != 2 {
		return
		// 	return fmt.Errorf("invalid values for header 'Range'")
	}

	begin, err := strconv.ParseInt(splitRange[0], 10, 64)
	if err != nil {
		catch(err)
	}

	end, err := strconv.ParseInt(splitRange[1], 10, 64)
	if err != nil {
		catch(err)
	}

	if begin > fileInfo.Size() || end > fileInfo.Size() {
		return
		// return fmt.Errorf("range out of bounds for file")
	}

	if begin >= end {
		return
		// return fmt.Errorf("range begin cannot be bigger than range end")
	}

	writer.Header().Set("Content-Length", strconv.FormatInt(end-begin+1, 10))

	writer.Header().Set("Content-Range",
		fmt.Sprintf("bytes %d-%d/%d", begin, end, fileInfo.Size()))

	writer.WriteHeader(http.StatusPartialContent)

	file.Seek(begin, 0)

	io.CopyN(writer, file, end-begin)
}
