package main

// Google Cloud Storage api
import (
	"bytes"
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	// "strings"
	"sync"
	"time"
)

var buf bytes.Buffer

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	id := h.Sum32()
	return strconv.FormatUint(uint64(id), 10)
}

const BUCKET_NAME = "basar-ocr-pdf-storage"

// get names of OCR result from  GCS1
// this function is used to download from the GCS
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

func getResultsInOrder(count int, fileNameId string) (int, []string) {
	// print it order
	var sortedList []string
	i := 1
	for i < count*2 {
		p := fmt.Sprintf("%s-Resultoutput-%d-to-%d.json", fileNameId, i, i+1)
		// if the file does not exist, try fixing it since there may be odd pages, else break
		if _, err := os.Stat(p); os.IsNotExist(err) {
			// if it does not exist, must mean we are at the end, try fix if odd pages => n-to-n not n-to-n+1
			pFixed := fmt.Sprintf("%s-Resultoutput-%d-to-%d.json", fileNameId, i, i)
			if _, err := os.Stat(pFixed); os.IsNotExist(err) {
				break
			}
			sortedList = append(sortedList, pFixed)
			break
		}

		sortedList = append(sortedList, p)
		i += 2
	}

	return len(sortedList), sortedList
}

func getJSONResultFiles(fileNameId string) (int, []string) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	var list []string
	for _, f := range files {
		fmt.Printf("%s vs %s\n", f.Name(), fileNameId)
		if filepath.Ext(f.Name()) == ".json" {
			// fmt.Println(f.Name())
			list = append(list, f.Name())
		}
	}
	return len(list), list
}

func DoOCR(jobID string, uploadedPDFName string, uploadedPDFBytes []byte) {
	start := time.Now()

	fileNameId := jobID

	fileName := fmt.Sprintf("%s-to-convert.pdf", fileNameId)
	// create the file remember to remove it
	ioutil.WriteFile(fileName, uploadedPDFBytes, 0644)

	// upload the file to do OCR on it
	if err := UploadFile(&buf, fileNameId, fileName); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success uploading file")
	}

	// file to do OCR on
	src := fmt.Sprintf("gs://basar-ocr-pdf-storage/%s", fileNameId)
	des := fmt.Sprintf("gs://basar-ocr-pdf-storage/%s%s", fileNameId, "-Result")

	fmt.Println("doing OCR")
	// detect OCR in the file we just uploaded in the OCR-Result directory
	err := DetectAsyncDocumentURI(&buf, src, des)
	if err != nil {
		fmt.Printf("Error in OCR: %s", err)
	}
	fmt.Println("OCR done !")

	// OCR is done, now download the JSON result files from GCS
	fmt.Println("dowloading files!")
	fileNames := getNamesOfOCRResult(fileNameId)

	var wg sync.WaitGroup

	for _, item := range fileNames {
		wg.Add(1)
		go downloadFile(&wg, &buf, item)

	}
	// wait for goroutines to finish
	wg.Wait()

	fmt.Println("finished downloading!")
	// collect the names of the downloaded JSON files in order
	cnt, jsonFileNames := getJSONResultFiles(fileNameId)
	if cnt == 0 {
		fmt.Println("getJSONResultFiles returned 0")
	}
	cnt, jsonFileNamesOrdered := getResultsInOrder(len(jsonFileNames), fileNameId)
	if cnt == 0 {
		fmt.Println("getResultsInOrder returned 0")
	}
	// parse each json file
	finalTextFileName := fmt.Sprintf("%s.txt", fileNameId)

	f, err := os.Create(finalTextFileName)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	fmt.Printf("jsonFileNamesORdered size: %d \n", len(jsonFileNamesOrdered))
	for _, jsonFileName := range jsonFileNamesOrdered {
		textResult := ParseJSONFile(jsonFileName)
		for i := range textResult {
			// fmt.Printf("TEXT RESULT %s\n", textResult[i])
			f.WriteString(textResult[i])
		}
	}

	// delete all GCS ojbects
	deleteAllObjects(fileNameId)

	elapsed := time.Since(start)
	fmt.Printf("time elapsed: %s\n", elapsed)
	fmt.Printf("finalTextFileName %s\n", finalTextFileName)
	// we wrote the result file to disk, now mark the jobb as compclete
	MarkAsComplete(jobID)
	// cleans out the JSON files we downloaded from GCS
	CleanDownloadedFiles(jobID)
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

	fmt.Println("making async annoate request, this will take a while")
	operation, err := client.AsyncBatchAnnotateFiles(ctx, request)

	// keep polling until we are done.
	for {
		pollResp, err := operation.Poll(ctx)
		if err != nil {
			fmt.Errorf("opeartion.Poll %v", err)
			// return err
		}
		if operation.Done() && pollResp == nil {
			fmt.Println("operation done but pollResp is nil. idk why")
		}
		if operation.Done() && pollResp != nil {
			// OCR is done and we have our response
			fmt.Println("async annotate request completed , OCR is done")
			fmt.Fprintf(w, "%v", pollResp)
			fmt.Println(buf.String())
			break
			// actual progres
		} else if !operation.Done() {
			meta, err := operation.Metadata()
			if err != nil {
				fmt.Printf("op.MetaData %v", err)
				return err
			}
			// fmt.Printf("\r Status: %\t %+v \t/10", meta)
			fmt.Printf("Current Status: \t %+v\n", meta)

		}
		time.Sleep(500 * time.Millisecond)
	}

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

// downloadFile downloads an object.
func downloadFile(wg *sync.WaitGroup, w io.Writer, object string) {
	defer wg.Done()
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

func deleteAllObjects(prefix string) {
	query := &storage.Query{Prefix: prefix}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: handle error.
	}
	bucket := client.Bucket(BUCKET_NAME)
	it := bucket.Objects(ctx, query)
	for {
		objAttrs, err := it.Next()
		if err != nil && err != iterator.Done {
			fmt.Printf("deleteAllObjects iterator error %v\n", err)
		}
		if err == iterator.Done {
			break
		}
		if err := bucket.Object(objAttrs.Name).Delete(ctx); err != nil {
			fmt.Printf("deleteAllObjects delete error %v\n", err)
		}
	}
	fmt.Println("deleted all object items in the bucket specified.")
}
