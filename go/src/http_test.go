package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestHttpServer 测试http服务端
func TestHttpServer(t *testing.T) {
	handler := func(respWriter http.ResponseWriter, req *http.Request) {
		// respWriter.Header().Set("Content-ColumnType", "image/gif")
		respWriter.Header().Set("Content-Type", "text/plain")
		respWriter.Write([]byte("Hello Work!"))
	}
	mux := http.NewServeMux()
	mux.Handle("/lalala", http.HandlerFunc(handler))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
	select {}
}

type Resp struct{}
type Req struct{}

// TestHttpSendGet 测试http发送Get请求
func TestHttpSendGet(t *testing.T) {
	req, err := http.NewRequest("GET", "http://www.4399.com", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	// 设置参数
	params := req.URL.Query()
	params.Add("id", strconv.FormatInt(1324, 10))
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	data := Resp{}
	assert.NoError(t, json.Unmarshal(body, &data))
}

// TestHttpSendPost 测试http发送Post请求
func TestHttpSendPost(t *testing.T) {
	param := Req{}
	bodyData, err := json.Marshal(param)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "http://www.4399.com", bytes.NewReader(bodyData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	data := Resp{}
	assert.NoError(t, json.Unmarshal(body, &data))
}

func TestHttpStream(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>Fetch API Example</title>
	<meta charset="utf-8">
</head>
<body onload="fetchData()">
	<div id="result"></div>
	<script>
		function fetchData() {
			fetch('http://localhost:8080/stream', {
                method: 'POST'
            })
			.then(response => {
                const encoding = response.headers.get('Content-Encoding');
	            const reader = encoding === 'gzip' ? response.body.pipeThrough(new window['zlib'].Gunzip()) : response.body.getReader();
				readData(reader);
			});
		}

		function readData(reader) {
			reader.read().then(result => {
				if (result.done) {
					return;
				}
            const decoder = new TextDecoder('utf-8');
		    const data = decoder.decode(result.value);
                var resultDiv = document.getElementById("result");
                resultDiv.style.whiteSpace = 'pre-wrap';        
                resultDiv.textContent += data;
				readData(reader);
			});
		}
	</script>
</body>
</html>`))
		w.(http.Flusher).Flush()
	})
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Transfer-Encoding", "chunked")
		for i := 0; i < 10; i++ {
			fmt.Fprintf(w, "data: %s\n", strconv.Itoa(i))
			w.(http.Flusher).Flush()
			time.Sleep(1 * time.Second)
		}
	})
	http.ListenAndServe(":8080", nil)
}
