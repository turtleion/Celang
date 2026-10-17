package evaluator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func initStdlibHTTP(env *Environment) {
	// =============================================================
	// 1. Klien HTTP (HTTP Client)
	// =============================================================

	// get / ambil
	getFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "http.get mengharapkan 1 atau 2 argumen: get(url, [opsi])", "Contoh: http.get(\"https://api.example.com/data\")")
			}
			urlStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen URL harus teks, menerima %s", args[0].Type()), "")
			}
			var optMap *Map
			if len(args) == 2 {
				if m, ok := args[1].(*Map); ok {
					optMap = m
				}
			}
			return executeHttpRequest("GET", urlStr.Value, optMap, nil)
		},
	}
	definePublicFunction(env, "get", getFn)
	definePublicFunction(env, "ambil", getFn)

	// post / kirim
	postFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 3 {
				return newError(nil, "http.post mengharapkan 1 hingga 3 argumen: post(url, [data], [opsi])", "Contoh: http.post(\"https://api.example.com\", {\"nama\": \"Celang\"})")
			}
			urlStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen URL harus teks, menerima %s", args[0].Type()), "")
			}
			var directBody Object
			var optMap *Map
			if len(args) >= 2 {
				directBody = args[1]
			}
			if len(args) == 3 {
				if m, ok := args[2].(*Map); ok {
					optMap = m
				}
			}
			return executeHttpRequest("POST", urlStr.Value, optMap, directBody)
		},
	}
	definePublicFunction(env, "post", postFn)
	definePublicFunction(env, "kirim", postFn)

	// put / taruh
	putFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 3 {
				return newError(nil, "http.put mengharapkan 1 hingga 3 argumen: put(url, [data], [opsi])", "")
			}
			urlStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen URL harus teks, menerima %s", args[0].Type()), "")
			}
			var directBody Object
			var optMap *Map
			if len(args) >= 2 {
				directBody = args[1]
			}
			if len(args) == 3 {
				if m, ok := args[2].(*Map); ok {
					optMap = m
				}
			}
			return executeHttpRequest("PUT", urlStr.Value, optMap, directBody)
		},
	}
	definePublicFunction(env, "put", putFn)
	definePublicFunction(env, "taruh", putFn)

	// patch / tambal
	patchFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 3 {
				return newError(nil, "http.patch mengharapkan 1 hingga 3 argumen: patch(url, [data], [opsi])", "")
			}
			urlStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen URL harus teks, menerima %s", args[0].Type()), "")
			}
			var directBody Object
			var optMap *Map
			if len(args) >= 2 {
				directBody = args[1]
			}
			if len(args) == 3 {
				if m, ok := args[2].(*Map); ok {
					optMap = m
				}
			}
			return executeHttpRequest("PATCH", urlStr.Value, optMap, directBody)
		},
	}
	definePublicFunction(env, "patch", patchFn)
	definePublicFunction(env, "tambal", patchFn)

	// delete / hapus
	deleteFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "http.delete mengharapkan 1 atau 2 argumen: delete(url, [opsi])", "")
			}
			urlStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen URL harus teks, menerima %s", args[0].Type()), "")
			}
			var optMap *Map
			if len(args) == 2 {
				if m, ok := args[1].(*Map); ok {
					optMap = m
				}
			}
			return executeHttpRequest("DELETE", urlStr.Value, optMap, nil)
		},
	}
	definePublicFunction(env, "delete", deleteFn)
	definePublicFunction(env, "hapus", deleteFn)

	// head / kepala
	headFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "http.head mengharapkan 1 atau 2 argumen: head(url, [opsi])", "")
			}
			urlStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen URL harus teks, menerima %s", args[0].Type()), "")
			}
			var optMap *Map
			if len(args) == 2 {
				if m, ok := args[1].(*Map); ok {
					optMap = m
				}
			}
			return executeHttpRequest("HEAD", urlStr.Value, optMap, nil)
		},
	}
	definePublicFunction(env, "head", headFn)
	definePublicFunction(env, "kepala", headFn)

	// request / minta
	requestFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "http.request mengharapkan 1 argumen konfigurasi permintaan (Map)", "Contoh: http.request({\"url\": \"...\", \"method\": \"POST\"})")
			}
			optMap, ok := args[0].(*Map)
			if !ok {
				return newError(nil, fmt.Sprintf("konfigurasi request harus berupa Map, menerima %s", args[0].Type()), "")
			}

			method := "GET"
			if mVal, ok := optMap.Pairs["method"]; ok {
				if s, ok := mVal.(*String); ok && s.Value != "" {
					method = strings.ToUpper(s.Value)
				}
			} else if mVal, ok := optMap.Pairs["metode"]; ok {
				if s, ok := mVal.(*String); ok && s.Value != "" {
					method = strings.ToUpper(s.Value)
				}
			}

			urlStr := ""
			if uVal, ok := optMap.Pairs["url"]; ok {
				if s, ok := uVal.(*String); ok {
					urlStr = s.Value
				}
			}
			if urlStr == "" {
				return newError(nil, "konfigurasi http.request wajib menyertakan properti 'url'", "")
			}

			return executeHttpRequest(method, urlStr, optMap, nil)
		},
	}
	definePublicFunction(env, "request", requestFn)
	definePublicFunction(env, "minta", requestFn)

	// =============================================================
	// 2. Utilitas URL (URL Helpers)
	// =============================================================

	// encode_uri / enkode_uri
	encodeUriFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "http.encode_uri mengharapkan 1 argumen teks", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			return &String{Value: url.PathEscape(s.Value)}
		},
	}
	definePublicFunction(env, "encode_uri", encodeUriFn)
	definePublicFunction(env, "enkode_uri", encodeUriFn)

	// decode_uri / dekode_uri
	decodeUriFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "http.decode_uri mengharapkan 1 argumen teks", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			dec, err := url.PathUnescape(s.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal decode URI: %v", err), "")
			}
			return &String{Value: dec}
		},
	}
	definePublicFunction(env, "decode_uri", decodeUriFn)
	definePublicFunction(env, "dekode_uri", decodeUriFn)

	// encode_uri_component / enkode_komponen_uri
	encodeCompFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "http.encode_uri_component mengharapkan 1 argumen teks", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			return &String{Value: url.QueryEscape(s.Value)}
		},
	}
	definePublicFunction(env, "encode_uri_component", encodeCompFn)
	definePublicFunction(env, "enkode_komponen_uri", encodeCompFn)

	// decode_uri_component / dekode_komponen_uri
	decodeCompFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "http.decode_uri_component mengharapkan 1 argumen teks", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks", "")
			}
			dec, err := url.QueryUnescape(s.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal decode komponen URI: %v", err), "")
			}
			return &String{Value: dec}
		},
	}
	definePublicFunction(env, "decode_uri_component", decodeCompFn)
	definePublicFunction(env, "dekode_komponen_uri", decodeCompFn)

	// parse_url / urai_url
	parseUrlFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "http.parse_url mengharapkan 1 argumen teks URL", "")
			}
			s, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen harus berupa teks URL", "")
			}
			u, err := url.Parse(s.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal mengurai URL: %v", err), "")
			}

			portStr := u.Port()
			var portObj Object = NULL
			if portStr != "" {
				if pInt, err := strconv.ParseInt(portStr, 10, 64); err == nil {
					portObj = &Integer{Value: pInt}
				} else {
					portObj = &String{Value: portStr}
				}
			}

			qMap := &Map{Pairs: make(map[string]Object), Order: []string{}}
			for k, vList := range u.Query() {
				qMap.Pairs[k] = &String{Value: strings.Join(vList, ", ")}
				qMap.Order = append(qMap.Order, k)
			}

			res := &Map{Pairs: make(map[string]Object), Order: []string{
				"scheme", "skema", "host", "port", "path", "jalur", "query", "kueri", "params", "parameter", "fragment", "fragmen",
			}}
			res.Pairs["scheme"] = &String{Value: u.Scheme}
			res.Pairs["skema"] = res.Pairs["scheme"]
			res.Pairs["host"] = &String{Value: u.Hostname()}
			res.Pairs["port"] = portObj
			res.Pairs["path"] = &String{Value: u.Path}
			res.Pairs["jalur"] = res.Pairs["path"]
			res.Pairs["query"] = &String{Value: u.RawQuery}
			res.Pairs["kueri"] = res.Pairs["query"]
			res.Pairs["params"] = qMap
			res.Pairs["parameter"] = qMap
			res.Pairs["fragment"] = &String{Value: u.Fragment}
			res.Pairs["fragmen"] = res.Pairs["fragment"]
			return res
		},
	}
	definePublicFunction(env, "parse_url", parseUrlFn)
	definePublicFunction(env, "urai_url", parseUrlFn)

	// =============================================================
	// 3. Server HTTP & Router (HTTP Server)
	// =============================================================

	// server / buat_server
	serverFn := &Builtin{
		Fn: func(args ...Object) Object {
			return createCelangHttpServerObject()
		},
	}
	definePublicFunction(env, "server", serverFn)
	definePublicFunction(env, "buat_server", serverFn)

	// layani / serve (shorthand jalan pintas)
	layaniFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "http.layani mengharapkan (port, handler)", "Contoh: http.layani(8080, fungsi(req, res) { res.kirim(\"Halo\") })")
			}
			srv := newCelangHttpServer()
			srv.addRoute("*", "/*", args[1])
			return srv.listen(args[0], true)
		},
	}
	definePublicFunction(env, "layani", layaniFn)
	definePublicFunction(env, "serve", layaniFn)
	definePublicFunction(env, "listen", layaniFn)

	// mulai / start (shorthand jalan pintas non-blocking)
	mulaiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "http.mulai mengharapkan (port, handler)", "Contoh: http.mulai(8080, fungsi(req, res) { res.kirim(\"Halo\") })")
			}
			srv := newCelangHttpServer()
			srv.addRoute("*", "/*", args[1])
			errObj := srv.listen(args[0], false)
			if isError(errObj) {
				return errObj
			}
			return srv.toMap()
		},
	}
	definePublicFunction(env, "mulai", mulaiFn)
	definePublicFunction(env, "start", mulaiFn)
}

// -----------------------------------------------------------------
// Eksekusi Klien HTTP
// -----------------------------------------------------------------

func executeHttpRequest(method string, targetUrl string, optionsMap *Map, directBody Object) Object {
	if targetUrl == "" {
		return newError(nil, "URL permintaan HTTP tidak boleh kosong", "")
	}

	timeout := 30 * time.Second
	headers := make(map[string]string)
	queryParams := make(map[string]string)
	var bodyReader io.Reader

	// Evaluasi optionsMap
	if optionsMap != nil {
		// 1. Timeout / Tenggat
		if tVal, ok := optionsMap.Pairs["timeout"]; ok {
			if tInt, ok := tVal.(*Integer); ok && tInt.Value > 0 {
				timeout = time.Duration(tInt.Value) * time.Millisecond
			} else if tFloat, ok := tVal.(*Float); ok && tFloat.Value > 0 {
				timeout = time.Duration(tFloat.Value * float64(time.Second))
			}
		} else if tVal, ok := optionsMap.Pairs["tenggat"]; ok {
			if tInt, ok := tVal.(*Integer); ok && tInt.Value > 0 {
				timeout = time.Duration(tInt.Value) * time.Millisecond
			}
		}

		// 2. Headers / Kepala
		var hMap *Map
		if hObj, ok := optionsMap.Pairs["headers"]; ok {
			if hm, ok := hObj.(*Map); ok {
				hMap = hm
			}
		} else if hObj, ok := optionsMap.Pairs["kepala"]; ok {
			if hm, ok := hObj.(*Map); ok {
				hMap = hm
			}
		}
		if hMap != nil {
			for k, v := range hMap.Pairs {
				if strV, ok := v.(*String); ok {
					headers[k] = strV.Value
				} else {
					headers[k] = v.Inspect()
				}
			}
		}

		// 3. Params / Parameter
		var pMap *Map
		if pObj, ok := optionsMap.Pairs["params"]; ok {
			if pm, ok := pObj.(*Map); ok {
				pMap = pm
			}
		} else if pObj, ok := optionsMap.Pairs["parameter"]; ok {
			if pm, ok := pObj.(*Map); ok {
				pMap = pm
			}
		}
		if pMap != nil {
			for k, v := range pMap.Pairs {
				if strV, ok := v.(*String); ok {
					queryParams[k] = strV.Value
				} else {
					queryParams[k] = v.Inspect()
				}
			}
		}

		// 4. Body / Badan
		if bObj, ok := optionsMap.Pairs["body"]; ok {
			directBody = bObj
		} else if bObj, ok := optionsMap.Pairs["badan"]; ok {
			directBody = bObj
		}

		// 5. JSON payload spesifik di options
		if jObj, ok := optionsMap.Pairs["json"]; ok {
			nativeVal, err := convertCelangToJSONNative(jObj)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal mengubah properti 'json' permintaan HTTP: %v", err), "")
			}
			jBytes, err := json.Marshal(nativeVal)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal encode JSON permintaan HTTP: %v", err), "")
			}
			bodyReader = bytes.NewReader(jBytes)
			if _, hasCT := headers["Content-Type"]; !hasCT {
				if _, hasCT2 := headers["content-type"]; !hasCT2 {
					headers["Content-Type"] = "application/json"
				}
			}
		}
	}

	// Evaluasi directBody jika belum diproses oleh json options
	if bodyReader == nil && directBody != nil && directBody != NULL {
		switch d := directBody.(type) {
		case *Bytes:
			bodyReader = bytes.NewReader(d.Elements)
		case *String:
			bodyReader = strings.NewReader(d.Value)
		case *Map, *Array, *StructInstance:
			nativeVal, err := convertCelangToJSONNative(d)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal mengubah data ke JSON: %v", err), "")
			}
			jBytes, err := json.Marshal(nativeVal)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal encode JSON data: %v", err), "")
			}
			bodyReader = bytes.NewReader(jBytes)
			if _, hasCT := headers["Content-Type"]; !hasCT {
				if _, hasCT2 := headers["content-type"]; !hasCT2 {
					headers["Content-Type"] = "application/json"
				}
			}
		default:
			bodyReader = strings.NewReader(d.Inspect())
		}
	}

	// Susun URL akhir dengan query parameters
	finalUrl := targetUrl
	if len(queryParams) > 0 {
		parsedUrl, err := url.Parse(targetUrl)
		if err == nil {
			q := parsedUrl.Query()
			for k, v := range queryParams {
				q.Set(k, v)
			}
			parsedUrl.RawQuery = q.Encode()
			finalUrl = parsedUrl.String()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, finalUrl, bodyReader)
	if err != nil {
		return newError(nil, fmt.Sprintf("Gagal membuat permintaan HTTP [%s] ke '%s': %v", method, finalUrl, err), "Pastikan format URL valid dan didukung (misal: 'https://...').")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return newError(nil, fmt.Sprintf("Permintaan HTTP [%s] ke '%s' gagal: %v", method, finalUrl, err), "Periksa koneksi internet, keabsahan URL, atau status server target.")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return newError(nil, fmt.Sprintf("Gagal membaca respons HTTP: %v", err), "")
	}

	return createHttpResponseObject(resp, bodyBytes)
}

func createHttpResponseObject(resp *http.Response, bodyBytes []byte) *Map {
	hPairs := make(map[string]Object)
	hOrder := []string{}
	for k, vList := range resp.Header {
		lowerK := strings.ToLower(k)
		valStr := strings.Join(vList, ", ")
		hPairs[lowerK] = &String{Value: valStr}
		hOrder = append(hOrder, lowerK)
	}
	headersMap := &Map{Pairs: hPairs, Order: hOrder}

	bodyStr := string(bodyBytes)
	isOk := resp.StatusCode >= 200 && resp.StatusCode < 300
	finalUrl := ""
	if resp.Request != nil && resp.Request.URL != nil {
		finalUrl = resp.Request.URL.String()
	}

	resMap := &Map{Pairs: make(map[string]Object), Order: []string{
		"status", "status_teks", "status_text", "sukses", "ok", "url", "headers", "kepala", "body", "badan",
	}}

	resMap.Pairs["status"] = &Integer{Value: int64(resp.StatusCode)}
	resMap.Pairs["status_teks"] = &String{Value: resp.Status}
	resMap.Pairs["status_text"] = resMap.Pairs["status_teks"]
	resMap.Pairs["sukses"] = nativeBool(isOk)
	resMap.Pairs["ok"] = resMap.Pairs["sukses"]
	resMap.Pairs["url"] = &String{Value: finalUrl}
	resMap.Pairs["headers"] = headersMap
	resMap.Pairs["kepala"] = headersMap
	resMap.Pairs["body"] = &String{Value: bodyStr}
	resMap.Pairs["badan"] = resMap.Pairs["body"]

	// text() / teks()
	textFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &String{Value: bodyStr}
		},
	}
	resMap.Pairs["text"] = textFn
	resMap.Pairs["teks"] = textFn
	resMap.Order = append(resMap.Order, "text", "teks")

	// bytes() / bita()
	bytesFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &Bytes{Elements: bodyBytes}
		},
	}
	resMap.Pairs["bytes"] = bytesFn
	resMap.Pairs["bita"] = bytesFn
	resMap.Order = append(resMap.Order, "bytes", "bita")

	// json()
	jsonFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(bodyBytes) == 0 {
				return &Map{Pairs: make(map[string]Object)}
			}
			var raw any
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			decoder.UseNumber()
			if err := decoder.Decode(&raw); err != nil {
				return newError(nil, fmt.Sprintf("Gagal mengurai respons JSON: %v", err), "Respons dari server bukan berupa JSON yang valid.")
			}
			return convertGoJSONToCelang(raw)
		},
	}
	resMap.Pairs["json"] = jsonFn
	resMap.Order = append(resMap.Order, "json")

	return resMap
}

// -----------------------------------------------------------------
// Server HTTP & Router Implementation
// -----------------------------------------------------------------

type routeHandler struct {
	method  string // "GET", "POST", "*", dll
	pattern string // e.g. "/api/pengguna/:id"
	handler Object // *Function or *Builtin
}

type celangMiddleware struct {
	prefix  string
	handler Object
}

type celangHttpServer struct {
	mu          sync.RWMutex
	routes      []routeHandler
	middlewares []celangMiddleware
	server      *http.Server
	listener    net.Listener
	running     bool
	port        int
	addr        string
}

func newCelangHttpServer() *celangHttpServer {
	return &celangHttpServer{
		routes:      make([]routeHandler, 0),
		middlewares: make([]celangMiddleware, 0),
	}
}

func (s *celangHttpServer) addRoute(method string, pattern string, handler Object) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes = append(s.routes, routeHandler{
		method:  strings.ToUpper(strings.TrimSpace(method)),
		pattern: pattern,
		handler: handler,
	})
}

func (s *celangHttpServer) addMiddleware(prefix string, handler Object) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.middlewares = append(s.middlewares, celangMiddleware{
		prefix:  prefix,
		handler: handler,
	})
}

func (s *celangHttpServer) listen(portArg Object, block bool) Object {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return newError(nil, "Server HTTP sudah dalam keadaan berjalan", "")
	}

	addr := ":8080"
	if portArg != nil && portArg != NULL {
		switch p := portArg.(type) {
		case *Integer:
			addr = fmt.Sprintf(":%d", p.Value)
		case *String:
			addr = p.Value
			if !strings.Contains(addr, ":") {
				addr = ":" + addr
			}
		}
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		s.mu.Unlock()
		return newError(nil, fmt.Sprintf("Gagal membuka listener HTTP pada '%s': %v", addr, err), "Pastikan port tidak sedang digunakan oleh aplikasi lain.")
	}

	tcpAddr := l.Addr().(*net.TCPAddr)
	s.listener = l
	s.port = tcpAddr.Port
	s.addr = l.Addr().String()
	s.running = true
	s.server = &http.Server{Handler: s}
	s.mu.Unlock()

	if block {
		err := s.server.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			return newError(nil, fmt.Sprintf("Server HTTP berhenti karena kesalahan: %v", err), "")
		}
		return NULL
	} else {
		go func() {
			_ = s.server.Serve(l)
		}()
		time.Sleep(10 * time.Millisecond)
		return NULL
	}
}

func (s *celangHttpServer) close() Object {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running || s.server == nil {
		return TRUE
	}
	s.running = false
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.server.Shutdown(ctx)
	if s.listener != nil {
		_ = s.listener.Close()
	}
	return TRUE
}

func (s *celangHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()

	var resMu sync.Mutex
	statusCode := http.StatusOK
	resHeaders := make(http.Header)
	sent := false

	resMap := &Map{Pairs: make(map[string]Object), Order: []string{}}
	setResProp := func(k string, v Object) {
		if _, ok := resMap.Pairs[k]; !ok {
			resMap.Order = append(resMap.Order, k)
		}
		resMap.Pairs[k] = v
	}

	// res.status(code)
	statusFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) >= 1 {
				if intVal, ok := args[0].(*Integer); ok {
					resMu.Lock()
					statusCode = int(intVal.Value)
					resMu.Unlock()
				}
			}
			return resMap
		},
	}
	setResProp("status", statusFn)

	// res.header(name, val) / res.kepala(name, val)
	headerFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) >= 2 {
				if k, ok := args[0].(*String); ok {
					if v, ok := args[1].(*String); ok {
						resMu.Lock()
						resHeaders.Set(k.Value, v.Value)
						resMu.Unlock()
					}
				}
			}
			return resMap
		},
	}
	setResProp("header", headerFn)
	setResProp("kepala", headerFn)

	// res.send(data) / res.kirim(data)
	sendFn := &Builtin{
		Fn: func(args ...Object) Object {
			resMu.Lock()
			defer resMu.Unlock()
			if sent {
				return resMap
			}
			sent = true

			for k, vals := range resHeaders {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			if w.Header().Get("Content-Type") == "" {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			}
			w.WriteHeader(statusCode)

			if len(args) >= 1 {
				switch d := args[0].(type) {
				case *Bytes:
					_, _ = w.Write(d.Elements)
				case *String:
					_, _ = w.Write([]byte(d.Value))
				default:
					_, _ = w.Write([]byte(d.Inspect()))
				}
			}
			return resMap
		},
	}
	setResProp("send", sendFn)
	setResProp("kirim", sendFn)

	// res.json(data, [code])
	jsonResFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "res.json() memerlukan setidaknya 1 argumen data", "")
			}
			resMu.Lock()
			defer resMu.Unlock()
			if sent {
				return resMap
			}
			sent = true

			if len(args) >= 2 {
				if intVal, ok := args[1].(*Integer); ok {
					statusCode = int(intVal.Value)
				}
			}

			for k, vals := range resHeaders {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(statusCode)

			nativeVal, err := convertCelangToJSONNative(args[0])
			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf(`{"error": "Gagal serialisasi JSON: %v"}`, err)))
				return resMap
			}
			jsonBytes, err := json.Marshal(nativeVal)
			if err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf(`{"error": "Gagal encode JSON: %v"}`, err)))
				return resMap
			}
			_, _ = w.Write(jsonBytes)
			return resMap
		},
	}
	setResProp("json", jsonResFn)

	// res.html(htmlStr, [code])
	htmlFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "res.html() memerlukan argumen string HTML", "")
			}
			resMu.Lock()
			defer resMu.Unlock()
			if sent {
				return resMap
			}
			sent = true

			if len(args) >= 2 {
				if intVal, ok := args[1].(*Integer); ok {
					statusCode = int(intVal.Value)
				}
			}

			for k, vals := range resHeaders {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(statusCode)

			if strVal, ok := args[0].(*String); ok {
				_, _ = w.Write([]byte(strVal.Value))
			} else {
				_, _ = w.Write([]byte(args[0].Inspect()))
			}
			return resMap
		},
	}
	setResProp("html", htmlFn)

	// res.redirect(targetUrl, [code]) / res.alihkan(targetUrl, [code])
	redirectFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return newError(nil, "res.redirect() memerlukan argumen URL tujuan", "")
			}
			resMu.Lock()
			defer resMu.Unlock()
			if sent {
				return resMap
			}
			sent = true

			redirCode := http.StatusFound
			if len(args) >= 2 {
				if intVal, ok := args[1].(*Integer); ok {
					redirCode = int(intVal.Value)
				}
			}
			targetUrl := ""
			if strVal, ok := args[0].(*String); ok {
				targetUrl = strVal.Value
			}
			http.Redirect(w, r, targetUrl, redirCode)
			return resMap
		},
	}
	setResProp("redirect", redirectFn)
	setResProp("alihkan", redirectFn)

	// Request Object
	reqMap := &Map{Pairs: make(map[string]Object), Order: []string{}}
	setReqProp := func(k string, v Object) {
		if _, ok := reqMap.Pairs[k]; !ok {
			reqMap.Order = append(reqMap.Order, k)
		}
		reqMap.Pairs[k] = v
	}

	setReqProp("method", &String{Value: r.Method})
	setReqProp("metode", &String{Value: r.Method})
	setReqProp("path", &String{Value: r.URL.Path})
	setReqProp("jalur", &String{Value: r.URL.Path})
	setReqProp("url", &String{Value: r.URL.String()})
	setReqProp("ip", &String{Value: r.RemoteAddr})
	setReqProp("alamat", &String{Value: r.RemoteAddr})

	// Headers
	hPairs := make(map[string]Object)
	hOrder := []string{}
	for hk, hvals := range r.Header {
		lk := strings.ToLower(hk)
		hPairs[lk] = &String{Value: strings.Join(hvals, ", ")}
		hOrder = append(hOrder, lk)
	}
	headersReqMap := &Map{Pairs: hPairs, Order: hOrder}
	setReqProp("headers", headersReqMap)
	setReqProp("kepala", headersReqMap)

	// Query
	qPairs := make(map[string]Object)
	qOrder := []string{}
	for qk, qvals := range r.URL.Query() {
		qPairs[qk] = &String{Value: strings.Join(qvals, ", ")}
		qOrder = append(qOrder, qk)
	}
	queryMap := &Map{Pairs: qPairs, Order: qOrder}
	setReqProp("query", queryMap)
	setReqProp("kueri", queryMap)

	// Body text & bytes
	bodyStr := string(bodyBytes)
	setReqProp("body", &String{Value: bodyStr})
	setReqProp("badan", &String{Value: bodyStr})

	// req.text() / req.teks()
	textReqFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &String{Value: bodyStr}
		},
	}
	setReqProp("text", textReqFn)
	setReqProp("teks", textReqFn)

	// req.bytes() / req.bita()
	bytesReqFn := &Builtin{
		Fn: func(args ...Object) Object {
			return &Bytes{Elements: bodyBytes}
		},
	}
	setReqProp("bytes", bytesReqFn)
	setReqProp("bita", bytesReqFn)

	// req.json()
	jsonReqFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(bodyBytes) == 0 {
				return &Map{Pairs: make(map[string]Object)}
			}
			var raw any
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			decoder.UseNumber()
			if err := decoder.Decode(&raw); err != nil {
				return newError(nil, fmt.Sprintf("Gagal mengurai JSON body request: %v", err), "Pastikan body request berformat JSON yang valid.")
			}
			return convertGoJSONToCelang(raw)
		},
	}
	setReqProp("json", jsonReqFn)

	// Dispatch middlewares then route
	s.mu.RLock()
	mws := make([]celangMiddleware, len(s.middlewares))
	copy(mws, s.middlewares)
	rts := make([]routeHandler, len(s.routes))
	copy(rts, s.routes)
	s.mu.RUnlock()

	idx := 0
	var dispatchNext func()

	dispatchNext = func() {
		resMu.Lock()
		isSent := sent
		resMu.Unlock()
		if isSent {
			return
		}

		for idx < len(mws) {
			mw := mws[idx]
			idx++
			if mw.prefix == "" || mw.prefix == "/" || strings.HasPrefix(r.URL.Path, mw.prefix) {
				nextFn := &Builtin{
					Fn: func(args ...Object) Object {
						dispatchNext()
						return NULL
					},
				}
				_ = applyFunction(mw.handler, []Object{reqMap, resMap, nextFn}, nil)
				return
			}
		}

		// Match route
		for _, rt := range rts {
			if rt.method != "*" && !strings.EqualFold(rt.method, r.Method) {
				continue
			}
			matched, params := matchRoute(rt.pattern, r.URL.Path)
			if matched {
				pPairs := make(map[string]Object)
				pOrder := []string{}
				for pk, pv := range params {
					pPairs[pk] = &String{Value: pv}
					pOrder = append(pOrder, pk)
				}
				paramsMap := &Map{Pairs: pPairs, Order: pOrder}
				setReqProp("params", paramsMap)
				setReqProp("parameter", paramsMap)

				_ = applyFunction(rt.handler, []Object{reqMap, resMap}, nil)
				return
			}
		}

		// 404 Not Found if unhandled
		resMu.Lock()
		defer resMu.Unlock()
		if !sent {
			sent = true
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("404 Halaman Tidak Ditemukan"))
		}
	}

	dispatchNext()
}

func matchRoute(pattern, path string) (bool, map[string]string) {
	if pattern == "*" || pattern == "/*" {
		return true, make(map[string]string)
	}

	trimmedPat := strings.Trim(pattern, "/")
	trimmedPath := strings.Trim(path, "/")

	if trimmedPat == "" && trimmedPath == "" {
		return true, make(map[string]string)
	}

	patParts := strings.Split(trimmedPat, "/")
	pathParts := strings.Split(trimmedPath, "/")

	params := make(map[string]string)

	for i := 0; i < len(patParts); i++ {
		pPart := patParts[i]
		if pPart == "*" {
			return true, params
		}
		if strings.HasPrefix(pPart, "*") {
			paramName := strings.TrimPrefix(pPart, "*")
			if paramName != "" && i < len(pathParts) {
				params[paramName] = strings.Join(pathParts[i:], "/")
			}
			return true, params
		}
		if i >= len(pathParts) {
			return false, nil
		}
		if strings.HasPrefix(pPart, ":") {
			paramName := strings.TrimPrefix(pPart, ":")
			params[paramName] = pathParts[i]
			continue
		}
		if pPart != pathParts[i] {
			return false, nil
		}
	}

	if len(patParts) != len(pathParts) {
		return false, nil
	}
	return true, params
}

func (s *celangHttpServer) toMap() *Map {
	serverMap := &Map{Pairs: make(map[string]Object), Order: []string{}}
	setProp := func(k string, v Object) {
		if _, ok := serverMap.Pairs[k]; !ok {
			serverMap.Order = append(serverMap.Order, k)
		}
		serverMap.Pairs[k] = v
	}

	// Helper for route registering methods
	registerMethod := func(httpMethod string) *Builtin {
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) < 2 {
					return newError(nil, fmt.Sprintf("Metode server.%s memerlukan (jalur, handler)", strings.ToLower(httpMethod)), "")
				}
				patStr, ok := args[0].(*String)
				if !ok {
					return newError(nil, "jalur rute harus berupa teks", "")
				}
				s.addRoute(httpMethod, patStr.Value, args[1])
				return serverMap
			},
		}
	}

	setProp("get", registerMethod("GET"))
	setProp("ambil", registerMethod("GET"))
	setProp("post", registerMethod("POST"))
	setProp("kirim", registerMethod("POST"))
	setProp("put", registerMethod("PUT"))
	setProp("taruh", registerMethod("PUT"))
	setProp("delete", registerMethod("DELETE"))
	setProp("hapus", registerMethod("DELETE"))
	setProp("patch", registerMethod("PATCH"))
	setProp("tambal", registerMethod("PATCH"))
	setProp("head", registerMethod("HEAD"))
	setProp("kepala", registerMethod("HEAD"))
	setProp("all", registerMethod("*"))
	setProp("semua", registerMethod("*"))

	// use / gunakan
	useFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 1 {
				s.addMiddleware("", args[0])
			} else if len(args) >= 2 {
				prefix := ""
				if pStr, ok := args[0].(*String); ok {
					prefix = pStr.Value
				}
				s.addMiddleware(prefix, args[1])
			} else {
				return newError(nil, "server.use memerlukan setidaknya 1 argumen fungsi middleware", "")
			}
			return serverMap
		},
	}
	setProp("use", useFn)
	setProp("gunakan", useFn)

	// listen / layani (blocking)
	listenFn := &Builtin{
		Fn: func(args ...Object) Object {
			var portObj Object = NULL
			if len(args) >= 1 {
				portObj = args[0]
			}
			res := s.listen(portObj, true)
			// update port property if assigned
			setProp("port", &Integer{Value: int64(s.port)})
			setProp("address", &String{Value: s.addr})
			setProp("alamat", &String{Value: s.addr})
			return res
		},
	}
	setProp("listen", listenFn)
	setProp("layani", listenFn)

	// start / mulai (non-blocking)
	startFn := &Builtin{
		Fn: func(args ...Object) Object {
			var portObj Object = NULL
			if len(args) >= 1 {
				portObj = args[0]
			}
			errObj := s.listen(portObj, false)
			if isError(errObj) {
				return errObj
			}
			setProp("port", &Integer{Value: int64(s.port)})
			setProp("address", &String{Value: s.addr})
			setProp("alamat", &String{Value: s.addr})
			return serverMap
		},
	}
	setProp("start", startFn)
	setProp("mulai", startFn)

	// close / tutup
	closeFn := &Builtin{
		Fn: func(args ...Object) Object {
			return s.close()
		},
	}
	setProp("close", closeFn)
	setProp("tutup", closeFn)

	// Dynamic port & address access
	portFn := &Builtin{
		Fn: func(args ...Object) Object {
			s.mu.RLock()
			defer s.mu.RUnlock()
			return &Integer{Value: int64(s.port)}
		},
	}
	setProp("ambil_port", portFn)
	setProp("get_port", portFn)

	addrFn := &Builtin{
		Fn: func(args ...Object) Object {
			s.mu.RLock()
			defer s.mu.RUnlock()
			return &String{Value: s.addr}
		},
	}
	setProp("ambil_alamat", addrFn)
	setProp("get_address", addrFn)

	return serverMap
}

func createCelangHttpServerObject() *Map {
	srv := newCelangHttpServer()
	return srv.toMap()
}
