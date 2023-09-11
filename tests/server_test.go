package regexhandlers_test

import (
	"context"
	// "errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func handleTeapot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	_, err := w.Write([]byte("I receive teapot-status code!"))
	if err != nil {
		panic(fmt.Sprintf("%s\n", err))
	}
}

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(
	host string,
	port uint16,

) *HTTPServer {
	mux := http.NewServeMux()
	mux.Handle("/api/", Handlers())
	mux.Handle("/", http.HandlerFunc(handleTeapot))
	server := http.Server{
		Addr:    net.JoinHostPort(host, fmt.Sprint(port)),
		Handler: mux,
	}
	this := &HTTPServer{}
	this.server = &server
	return this
}

func (s *HTTPServer) Start() error {
	// if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
func TestByteResponse(t *testing.T) {
	var host string = "localhost"
	var port uint16 = 8000
	var response *http.Response
	var err error
	ctx := context.Background()
	httpServer := NewHTTPServer(
		host,
		port,
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := httpServer.Start()
		if err != nil {
			fmt.Printf("TestByteResponse Start err %s\n", err)
		}
	}()
	time.Sleep(1 * time.Second)
	client := &http.Client{}
	requestUrl := fmt.Sprintf("http://%s:%d/api/version", host, port)
	payload := strings.NewReader("")
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, payload)
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %s\n", requestUrl)
		return
	}
	response, err = client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error http request: %s\n", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("FAIL: error making http request: %s\n", err)
		return
	}
	if string(body) != "1.0.0" {
		t.Errorf("expected %s, but get %s, ", "1.0.0", body)
	}
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	err = httpServer.Stop(context.Background())
	if err != nil {
		fmt.Printf("TestByteResponse Stop err %s\n", err)
	}
	wg.Wait()
}

func TestJsonResponse(t *testing.T) {
	var host string = "localhost"
	var port uint16 = 8001
	var response *http.Response
	var err error
	child_name := "BenderBendingRodriguez"
	parent_id := 8096
	httpServer := NewHTTPServer(
		host,
		port,
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := httpServer.Start()
		if err != nil {
			fmt.Printf("TestJsonResponse Start err %s\n", err)
		}
	}()
	time.Sleep(1 * time.Second)
	client := &http.Client{}
	requestNested := fmt.Sprintf("http://%s:%d/api/get/%d/%s", host, port, parent_id, child_name)
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		requestNested,
		strings.NewReader(``))
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %s\n", requestNested)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("FAIL: error making http request: %s\n", err)
		return
	}
	ethalon := fmt.Sprintf(`{"child_name":"%s","parent_id":"%d"}`, child_name, parent_id)
	if string(body) != ethalon {
		t.Errorf("FAIL: get %s\n", body)
		t.Errorf("FAIL: expected %s\n", ethalon)
	} else {
		fmt.Printf("OK: %s\n", body)
	}
	err = httpServer.Stop(context.Background())
	if err != nil {
		fmt.Printf("TestJsonResponse Stop err %s\n", err)
	}
	wg.Wait()
}
