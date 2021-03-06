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

func catch(err error) {
	if err != nil {
		fmt.Println("Error in file.go")
		fmt.Println(err)
		panic(err)
	}
}
func ServeFile(writer http.ResponseWriter, request *http.Request, file *os.File) {
	fmt.Println("func ServeFile in file.go")
	defer file.Close()

	// Reading header info from the opened file, this will be used for response header "Content-Type"
	// the first 512 bytes is the header.
	fileHeader := make([]byte, 512)
	_, err := file.Read(fileHeader) // File offset is now len(fileHeader)
	fileType := http.DetectContentType(fileHeader)

	if err != nil {
		catch(err)
	}

	// Get file info which we will use for the response headers "Content-Disposition" and "Content-Length"
	fileInfo, err := file.Stat()
	if err != nil {
		catch(err)
	}

	// Set default headers

	// attachment is required to tell some (older) browsers who follow an href to download the file
	// instead of showing/printing the content to the screen.
	// For example, if you click a link to an image, the browser will pop up the download dialog
	// box compared to drawing the image in the browser tab.
	writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s""`, fileInfo.Name()))
	// A must for every request that has a body (see RFC 2616 section 7.2.1)
	// writer.Header().Set("Content-Type", http.DetectContentType(fileHeader))

	writer.Header().Set("Content-Type", fileType)

	// and pick up later where they left up. Or download managers to establish multiple connections
	// Tell the client we accept ranges, this gives clients the option to pause the transfer
	// writer.Header().Set("Accept-Ranges", "bytes")

	// Check if the client requests a range from the file (see RFC 7233 section 4.2)
	requestRange := request.Header.Get("range")
	fmt.Printf("request range: %s\n", requestRange)

	if requestRange == "" {
		// No range is defined, tell the client the incoming length of data, the size of the open file
		writer.Header().Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))

		// Since we read 512 bytes for 'fileHeader' earlier, we set the reader offset back
		// to 0 starting from the beginning of the the file (the 0 in the second argument)
		file.Seek(0, 0)

		fmt.Println("filename path: ", filepath.Ext(fileInfo.Name()))

		// Stream the file to the client
		io.Copy(writer, file)

	}

	// fmt.Println("requst range second part")

	// Client requests a part of the file, this shouldn't happen tho

	// Decode the request header to integers we can use for offset
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

	// Tell the amount bytes the client will receive
	writer.Header().Set("Content-Length", strconv.FormatInt(end-begin+1, 10))

	// Confirm the range values to the client, and the total size of the file
	// 'Content-Range' : 'bytes begin-end/totalFileSize'
	writer.Header().Set("Content-Range",
		fmt.Sprintf("bytes %d-%d/%d", begin, end, fileInfo.Size()))
	// Response http status code 206

	writer.WriteHeader(http.StatusPartialContent)

	// Set the file offset to the requested beginning
	file.Seek(begin, 0)

	// Send the (end-begin) amount of bytes to the client
	io.CopyN(writer, file, end-begin)
}
