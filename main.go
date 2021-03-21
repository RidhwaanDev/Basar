// Sample vision-quickstart uses the Google Cloud Vision API to label an image.
package main

import (
	"bytes"
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/iterator"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

var buf bytes.Buffer

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	id := h.Sum32()
	return strconv.FormatUint(uint64(id), 10)
}

const BUCKET_NAME = "basar-ocr-pdf-storage"

func getNamesOfOCRResult(prefix string) []string {
	query := &storage.Query{Prefix: prefix}

	var names []string

	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		fmt.Println(err)
	}

	bkt := client.Bucket(BUCKET_NAME)
	it := bkt.Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
	}

	for _, name := range names {
		fmt.Println(name)
	}
	return names
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("please support filename")
	}

	fileName := os.Args[1]
	fileNameId := hash(fileName)
	// fmt.Printf("file name id hash %s\n", fileNameId)

	// upload the file to do OCR on it
	if err := UploadFile(&buf, fileNameId, fileName); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success uploading file")
	}

	// file to do OCR on
	src := fmt.Sprintf("gs://basar-ocr-pdf-storage/%s", fileNameId)
	des := fmt.Sprintf("gs://basar-ocr-pdf-storage/%s%s", fileNameId, "-Result")

	// detect OCR in the file we just uploaded in the OCR-Result directory
	err := DetectAsyncDocumentURI(&buf, src, des)
	if err != nil {
		fmt.Printf("Error in OCR: %s", err)
	}

	fileNames := getNamesOfOCRResult(fileNameId)
	for _, item := range fileNames {
		downloadFile(&buf, item)
	}

	// bytes, err := downloadFile(&buf)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = ioutil.WriteFile("gsresult.pdf", bytes, 0666)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

// detectAsyncDocumentURI performs Optical Character Recognition (OCR) on a
// PDF file stored in GCS.
func DetectAsyncDocumentURI(w io.Writer, gcsSourceURI, gcsDestinationURI string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	request := &visionpb.AsyncBatchAnnotateFilesRequest{
		Requests: []*visionpb.AsyncAnnotateFileRequest{
			{
				Features: []*visionpb.Feature{
					{
						Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
					},
				},
				InputConfig: &visionpb.InputConfig{
					GcsSource: &visionpb.GcsSource{Uri: gcsSourceURI},
					// Supported MimeTypes are: "application/pdf" and "image/tiff".
					MimeType: "application/pdf",
				},
				OutputConfig: &visionpb.OutputConfig{
					GcsDestination: &visionpb.GcsDestination{Uri: gcsDestinationURI},
					// How many pages should be grouped into each json output file.
					BatchSize: 2,
				},
			},
		},
	}

	fmt.Println("making request")
	operation, err := client.AsyncBatchAnnotateFiles(ctx, request)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Waiting for the operation to finish.")

	resp, err := operation.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%v", resp)
	fmt.Println(buf.String())

	return nil
}

// uploadFile uploads an object.
func UploadFile(w io.Writer, fileNameHash string, localfileName string) error {
	bucket := "basar-ocr-pdf-storage"
	object := fileNameHash
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(localfileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	return nil
}

// func downloadOCRResults() ([]byte, error) {
// 	itemsToDownload := getNamesOfOCRResult()
// 	ctx := context.Background()
// 	client, err := storage.NewClient(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("storage.NewClient: %v", err)
// 	}
// 	defer client.Close()
//
// 	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
// 	defer cancel()
// 	for _, object := range itemsToDownload {
// 		rc, err := client.Bucket(BUCKET_NAME).Object(object).NewReader(ctx)
// 		if err != nil {
// 			return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
// 		}
// 		defer rc.Close()
//
// 		data, err := ioutil.ReadAll(rc)
// 		if err != nil {
// 			return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
// 		}
// 		fmt.Fprintf(w, "Blob %v downloaded.\n", object)
// 		err := ioutil.WriteFile(object, data, 666)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// 	return nil
// }
//
// downloadFile downloads an object.
func downloadFile(w io.Writer, object string) {
	bucket := "basar-ocr-pdf-storage"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		fmt.Printf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)

	if err != nil {
		fmt.Printf("ioutil.ReadAll: %v", err)
	}
	fmt.Fprintf(w, "Blob %v downloaded.\n", object)

	// write file to disk
	err = ioutil.WriteFile(object, data, 0644)
	if err != nil {
		fmt.Printf("ioutil.WriteFile: %v", err)
	}

}
