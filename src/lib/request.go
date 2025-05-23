package lib

import (
	"io"
	"net/http"
	"net/url"
)

func HttpRequest(reqURL string, method string, headers map[string]string, body io.Reader) ([]byte, error) {
	parsedURL, err := url.Parse(reqURL)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: method,
		URL:    parsedURL,
		Header: make(http.Header),
		Body:   io.NopCloser(body),
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		if location, err := resp.Location(); err == nil {
			newHeaders := make(map[string]string)
			for key, values := range req.Header {
				if len(values) > 0 {
					newHeaders[key] = values[0]
				}
			}
			return HttpRequest(location.String(), req.Method, newHeaders, body)
		}
	}

	return respBody, nil
}

//usage example
// // Example 1: GET request
// fmt.Println("Example 1: GET request")
// getResp, getErr := lib.HttpRequest("https://jsonplaceholder.typicode.com/todos/1", http.MethodGet, nil, nil)
// if getErr != nil {
// 	log.Fatalf("GET request failed: %v", getErr)
// }
// fmt.Printf("GET Response:\n%s\n\n", getResp)

// // Example 2: POST request with JSON body and headers
// fmt.Println("Example 2: POST request with JSON body and headers")
// jsonBody := []byte(`{"title": "foo", "body": "bar", "userId": 1}`)
// headers := map[string]string{"Content-Type": "application/json"}
// postResp, postErr := lib.HttpRequest("https://jsonplaceholder.typicode.com/posts", http.MethodPost, headers, bytes.NewBuffer(jsonBody))
// if postErr != nil {
// 	log.Fatalf("POST request failed: %v", postErr)
// }
// fmt.Printf("POST Response:\n%s\n\n", postResp)

// // Example 3: GET request with redirect
// fmt.Println("Example 3: GET request with redirect")
// getResp, getErr = lib.HttpRequest("https://www.httpwatch.com/httpgallery/redirection/", http.MethodGet, nil, nil)
// if getErr != nil {
// 	log.Fatalf("GET request failed: %v", getErr)
// }
// fmt.Printf("GET Response:\n%s\n\n", getResp)

// // Example 4: POST request with empty body
// fmt.Println("Example 4: POST request with empty body")
// postResp, postErr = lib.HttpRequest("https://jsonplaceholder.typicode.com/posts", http.MethodPost, headers, bytes.NewBuffer([]byte{}))
// if postErr != nil {
// 	log.Fatalf("POST request failed: %v", postErr)
// }
// fmt.Printf("POST Response:\n%s\n\n", postResp)

// // Example 5 : test invalid url
// fmt.Println("Example 5: test invalid url")
// getResp, getErr = lib.HttpRequest("https://invalid.com", http.MethodGet, nil, nil)
// if getErr != nil {
// 	log.Fatalf("GET request failed: %v", getErr)
// }
// fmt.Printf("GET Response:\n%s\n\n", getResp)

// // Example 1: GET request
// fmt.Println("Example 1: GET request")
// getResp, getErr := lib.HttpRequest("https://jsonplaceholder.typicode.com/todos/1", http.MethodGet, nil, nil)
// if getErr != nil {
// 	log.Fatalf("GET request failed: %v", getErr)
// }
// fmt.Printf("GET Response:\n%s\n\n", getResp)

// // Example 2: POST request with JSON body and headers
// fmt.Println("Example 2: POST request with JSON body and headers")
// jsonBody := []byte(`{"title": "foo", "body": "bar", "userId": 1}`)
// headers := map[string]string{"Content-Type": "application/json"}
// postResp, postErr := lib.HttpRequest("https://jsonplaceholder.typicode.com/posts", http.MethodPost, headers, bytes.NewBuffer(jsonBody))
// if postErr != nil {
// 	log.Fatalf("POST request failed: %v", postErr)
// }

// // Convert postResp to string for logging
// postRespString := string(postResp)
// log.Printf("POST Response:\n%s\n\n", postRespString)

// // Example 3: GET request with redirect
// fmt.Println("Example 3: GET request with redirect")
// getResp, getErr = lib.HttpRequest("https://www.httpwatch.com/httpgallery/redirection/", http.MethodGet, nil, nil)
// if getErr != nil {
// 	log.Fatalf("GET request failed: %v", getErr)
// }
// fmt.Printf("GET Response:\n%s\n\n", getResp)

// // Example 4: POST request with empty body
// fmt.Println("Example 4: POST request with empty body")
// postResp, postErr = lib.HttpRequest("https://jsonplaceholder.typicode.com/posts", http.MethodPost, headers, bytes.NewBuffer([]byte{}))
// if postErr != nil {
// 	log.Fatalf("POST request failed: %v", postErr)
// }
// postRespString = string(postResp)
// log.Printf("POST Response:\n%s\n\n", postRespString)

// // Example 5 : test invalid url
// fmt.Println("Example 5: test invalid url")
// getResp, getErr = lib.HttpRequest("https://invalid.com", http.MethodGet, nil, nil)
// if getErr != nil {
// 	log.Fatalf("GET request failed: %v", getErr)
// }
// fmt.Printf("GET Response:\n%s\n\n", getResp)
