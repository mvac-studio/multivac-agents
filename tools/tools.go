package tools

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/tebeka/selenium"
	"io"
	"log"
	"multivac.network/services/agents/providers"
	"multivac.network/services/agents/providers/fireworks"
	"net/url"
	"os"
	"time"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

func GetCurrentDate(format string) string {
	result := fmt.Sprintf("Current date: %s", time.Now())
	return result
}

func uploadImageToS3(reader io.Reader) string {
	fid, err := uuid.NewUUID()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-lax-1"),                           // replace with your bucket's region
		Endpoint:    aws.String("https://us-lax-1.linodeobjects.com"), // replace with your bucket's endpoint
		Credentials: credentials.NewStaticCredentials("A40MMRVFX465IFONTDAP", "sjeBRYS2z1JjYBFUTFQUFw8LdZeKpLXQEpUaemwV", ""),
	})
	if err != nil {
		fmt.Println("Error creating session,", err)
		return ""
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("files.ngent.io"),
		Key:    aws.String("/screenshots/" + fid.String() + ".png"),
		Body:   reader,
		ACL:    aws.String("public-read"), // to make the file publicly accessible
	})
	if err != nil {
		fmt.Println("Error uploading file,", err)
		return ""
	}

	fmt.Println("Upload successful. URL:", result.Location)
	return result.Location
}

func OpenWebAddress(address string) string {
	// Connect to the WebDriver instance running on the Selenium Hub.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://selenium-hub.default.svc.cluster.local:4444/wd/hub") // replace with your Selenium Hub address
	if err != nil {
		log.Printf("Failed to open session: %v", err)
	}
	defer wd.Quit()
	u, err := url.Parse(address)
	if u.Scheme == "" {
		u.Scheme = "https"

	}
	// Navigate to the webpage
	err = wd.Get(u.String())
	if err != nil {
		log.Printf("Failed to navigate: %v", err)
	}
	windows, err := wd.WindowHandles()
	// Take a screenshot
	err = wd.ResizeWindow(windows[0], 1920, 6000)
	if err != nil {
		log.Printf("Failed to resize window: %v", err)

	}
	screenshot, err := wd.Screenshot()

	if err != nil {
		log.Printf("Failed to take screenshot: %v", err)
	}

	// Save the screenshot to a file
	location := uploadImageToS3(bytes.NewReader(screenshot))

	return describePicture(location)
}

func describePicture(location string) string {
	apiKey := os.Getenv("FIREWORKS_API_KEY")
	model := "firellava-13b"
	provider := fireworks.NewService(model, apiKey, 3000)
	request := providers.Request{Messages: []providers.Message{
		{Role: "user", Content: "Describe this picture"},
		{Role: "user", ImageContent: location},
	}}
	response, err := provider.SendRequest(request)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return response.Content
}

//func OpenWebAddress(address string) string {
//
//	client := &http.Client{}
//	req, err := http.NewRequest("GET", address, nil)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return err.Error()
//	}
//
//	// Set User-Agent header
//	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1")
//
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return err.Error()
//	}
//	defer resp.Body.Close()
//
//	doc, err := html.Parse(resp.Body)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return err.Error()
//	}
//
//	links := strings.Builder{}
//	content := strings.Builder{}
//	var f func(*html.Node)
//	f = func(n *html.Node) {
//		if n.Type == html.TextNode && n.Data != "style" && n.Data != "script" {
//			text := strings.TrimSpace(n.Data)
//			if len(text) > 0 {
//				content.WriteString(fmt.Sprintf("%s\n", text))
//				fmt.Println(text)
//			}
//		}
//		for c := n.FirstChild; c != nil; c = c.NextSibling {
//			if c.Data != "script" && c.Data != "style" {
//				f(c)
//			}
//		}
//	}
//	f(doc)
//	return fmt.Sprintf("Links: %s\nContent: %s", links.String(), content.String())
//}
