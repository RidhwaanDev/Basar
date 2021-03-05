# ArabicOCR - بصر

Extract text from scanned Arabic PDF. 

# Steps to test on macOS/OSX

1) Download Google API credentials from their console. See https://cloud.google.com/vision/docs/ocr#vision_text_detection-go. Look at Set up your GCP project and authentication section
2) Clone repo
3) Make sure you have Go installed. Can be downloaded via homebrew
4) Open terminal, cd into the project directory. Type `export GOOGLE_APPLICATION_CREDENTIALS=your_credentials.json`
5) run `go run .`
6) open `index.html`
7) upload PDF and observe terminal. 
