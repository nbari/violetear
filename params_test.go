package violetear

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestXXX(t *testing.T) {
	router := New()
	router.AddRegex(":word", `^\w+$`)
	router.HandleFunc("/test/:word/:word/:word", func(w http.ResponseWriter, r *http.Request) {
		param := GetParam("word", r, 3)
		expect(t, param, "")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/foo/bar/xxxx", nil)
	router.ServeHTTP(w, req)
}

func TestGetParam(t *testing.T) {
	tt := []struct {
		path          string
		requestPath   string
		param         string
		expectedParam string
		index         int
		err           bool
	}{
		{
			path:          "/tests/:test_param",
			requestPath:   "/tests/abc",
			param:         "test_param",
			expectedParam: "abc",
		},
		{
			path:          "/other_test",
			requestPath:   "/other_test",
			param:         "foo",
			expectedParam: "",
		},
		{
			path:          "/other_test",
			requestPath:   "/other_test",
			param:         "",
			expectedParam: "",
		},
		{
			path:          "/test/:ip",
			requestPath:   "/test/127.0.0.1",
			param:         "ip",
			expectedParam: "127.0.0.1",
		},
		{
			path:          "/test/:ip",
			requestPath:   "/test/127.0.0.1",
			param:         "ip",
			expectedParam: "127.0.0.1",
			index:         3,
		},
		{
			path:          "/:uuid",
			requestPath:   "/78F204D2-26D9-409F-BE81-2E5D061E1FA1",
			param:         "uuid",
			expectedParam: "78F204D2-26D9-409F-BE81-2E5D061E1FA1",
		},
		{
			path:          "/test/:uuid",
			requestPath:   "/test/78F204D2-26D9-409F-BE81-2E5D061E1FA1",
			param:         "uuid",
			expectedParam: "78F204D2-26D9-409F-BE81-2E5D061E1FA1",
		},
		{
			path:          "/test/:uuid/:uuid",
			requestPath:   "/test/78F204D2-26D9-409F-BE81-2E5D061E1FA1/33A7B724-1498-4A5A-B29B-AD4E31824234",
			param:         "uuid",
			expectedParam: "78F204D2-26D9-409F-BE81-2E5D061E1FA1",
			index:         0,
		},
		{
			path:          "/test/:uuid/:uuid",
			requestPath:   "/test/78F204D2-26D9-409F-BE81-2E5D061E1FA1/33A7B724-1498-4A5A-B29B-AD4E31824234",
			param:         "uuid",
			expectedParam: "33A7B724-1498-4A5A-B29B-AD4E31824234",
			index:         1,
		},
		{
			path:          "/test/:uuid/:uuid",
			requestPath:   "/test/78F204D2-26D9-409F-BE81-2E5D061E1FA1/33A7B724-1498-4A5A-B29B-AD4E31824234",
			param:         "uuid",
			expectedParam: "78F204D2-26D9-409F-BE81-2E5D061E1FA1",
			index:         -1,
		},
		{
			path:          "/test/2/:uuid/:uuid",
			requestPath:   "/test/2/78F204D2-26D9-409F-BE81-2E5D061E1FA1/33A7B724-1498-4A5A-B29B-AD4E31824234",
			param:         "uuid",
			expectedParam: "",
			index:         20,
		},
		{
			path:          "/asterisk/*",
			requestPath:   "/asterisk/foo",
			param:         "*",
			expectedParam: "foo",
		},
		{
			path:          "/asterisk/asterisk/*/*/*",
			requestPath:   "/test/a/b/c/d/e",
			param:         "*",
			expectedParam: "a",
			err:           true,
		},
		{
			path:          "/asterisk/foo/*/*/*/3",
			requestPath:   "/test/foo/xxx",
			param:         "*",
			expectedParam: "xxx",
			err:           true,
		},
	}

	router := New()
	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}
	router.AddRegex(":test_param", `^\w+$`)

	var (
		w             *httptest.ResponseRecorder
		obtainedParam string
	)

	for _, tc := range tt {
		t.Run(tc.path, func(t *testing.T) {
			testHandler := func(w http.ResponseWriter, r *http.Request) {
				if tc.index > 0 {
					obtainedParam = GetParam(tc.param, r, tc.index)
				} else {
					obtainedParam = GetParam(tc.param, r)
				}
				expect(t, obtainedParam, tc.expectedParam)
			}
			_, err := router.HandleFunc(tc.path, testHandler, "GET")
			expect(t, err != nil, tc.err)
			w = httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.requestPath, nil)
			router.ServeHTTP(w, req)
		})
	}
}

func TestGetParams(t *testing.T) {
	tt := []struct {
		path          string
		requestPath   string
		param         string
		expectedParam []string
	}{
		{
			path:          "/tests/:test_param",
			requestPath:   "/tests/abc",
			param:         "test_param",
			expectedParam: []string{"abc"},
		},
		{
			path:          "/other_test",
			requestPath:   "/other_test",
			param:         "foo",
			expectedParam: []string{},
		},
		{
			path:          "/other_test",
			requestPath:   "/other_test",
			param:         "",
			expectedParam: []string{},
		},
		{
			path:          "/test/:ip",
			requestPath:   "/test/127.0.0.1",
			param:         "ip",
			expectedParam: []string{"127.0.0.1"},
		},
		{
			path:          "/test/:ip",
			requestPath:   "/test/127.0.0.1",
			param:         "ip",
			expectedParam: []string{"127.0.0.1"},
		},
		{
			path:          "/:uuid",
			requestPath:   "/78F204D2-26D9-409F-BE81-2E5D061E1FA1",
			param:         "uuid",
			expectedParam: []string{"78F204D2-26D9-409F-BE81-2E5D061E1FA1"},
		},
		{
			path:          "/test/:uuid",
			requestPath:   "/test/78F204D2-26D9-409F-BE81-2E5D061E1FA1",
			param:         "uuid",
			expectedParam: []string{"78F204D2-26D9-409F-BE81-2E5D061E1FA1"},
		},
		{
			path:          "/test/:uuid/:uuid",
			requestPath:   "/test/78F204D2-26D9-409F-BE81-2E5D061E1FA1/33A7B724-1498-4A5A-B29B-AD4E31824234",
			param:         "uuid",
			expectedParam: []string{"78F204D2-26D9-409F-BE81-2E5D061E1FA1", "33A7B724-1498-4A5A-B29B-AD4E31824234"},
		},
		{
			path:          "/test/:uuid/:uuid:uuid",
			requestPath:   "/test/479BA626-0565-49CF-8852-9576F6C9964F/479BA626-0565-49CF-8852-9576F6C9964F/479BA626-0565-49CF-8852-9576F6C9964F",
			param:         "uuid",
			expectedParam: []string{"479BA626-0565-49CF-8852-9576F6C9964F", "479BA626-0565-49CF-8852-9576F6C9964F", "479BA626-0565-49CF-8852-9576F6C9964F"},
		},
		{
			path:          "/test/:uuid/:uuid:uuid",
			requestPath:   "/test/479BA626-0565-49CF-8852-9576F6C9964F/479BA626-0565-49CF-8852-9576F6C9964F/479BA626-0565-49CF-8852-9576F6C9964F",
			param:         "uuid",
			expectedParam: []string{"479BA626-0565-49CF-8852-9576F6C9964F", "479BA626-0565-49CF-8852-9576F6C9964F", "479BA626-0565-49CF-8852-9576F6C9964F"},
		},
	}

	router := New()
	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}
	router.AddRegex(":test_param", `^\w+$`)

	var (
		w              *httptest.ResponseRecorder
		obtainedParams []string
	)

	for _, tc := range tt {
		t.Run(tc.path, func(t *testing.T) {
			testHandler := func(w http.ResponseWriter, r *http.Request) {
				obtainedParams = GetParams(tc.param, r)
				expectDeepEqual(t, obtainedParams, tc.expectedParam)
			}
			router.HandleFunc(tc.path, testHandler, "GET")
			w = httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.requestPath, nil)
			router.ServeHTTP(w, req)
		})
	}
}

func TestGetParamDuplicates(t *testing.T) {
	var uuids []string
	request := "/test/"
	requestHandler := "/test/"
	for i := 0; i <= 10; i++ {
		uuid := genUUID()
		uuids = append(uuids, uuid)
		request += fmt.Sprintf("%s/", uuid)
		requestHandler += ":uuid/"
	}

	router := New()

	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i <= 10; i++ {
			expect(t, GetParam("uuid", r, i), uuids[i])
		}
		w.Write([]byte("named params"))
	}

	router.HandleFunc(requestHandler, handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", request, nil)
	router.ServeHTTP(w, req)
	//expect(t, w.Code, 200)
}

func TestGetParamsDuplicates(t *testing.T) {
	var uuids []string
	request := "/test/"
	requestHandler := "/test/"
	for i := 0; i < 10; i++ {
		uuid := genUUID()
		uuids = append(uuids, uuid)
		request += fmt.Sprintf("%s/", uuid)
		requestHandler += ":uuid/"
	}

	router := New()

	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		p := GetParams("uuid", r)
		expect(t, true, (reflect.DeepEqual(p, uuids)))
		w.Write([]byte("named params"))
	}

	router.HandleFunc(requestHandler, handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", request, nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
}

func TestGetParamsDuplicatesLogRequests(t *testing.T) {
	var uuids []string
	request := "/test/"
	requestHandler := "/test/"
	for i := 0; i < 10; i++ {
		uuid := genUUID()
		uuids = append(uuids, uuid)
		request += fmt.Sprintf("%s/", uuid)
		requestHandler += ":uuid/"
	}

	router := New()
	router.LogRequests = true

	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		p := GetParams("uuid", r)
		expect(t, true, (reflect.DeepEqual(p, uuids)))
		w.Write([]byte("named params"))
	}

	router.HandleFunc(requestHandler, handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", request, nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
}

func TestGetParamsDuplicatesNonExistent(t *testing.T) {
	var uuids []string
	request := "/test/"
	requestHandler := "/test/"
	for i := 0; i < 3; i++ {
		uuid := genUUID()
		uuids = append(uuids, uuid)
		request += fmt.Sprintf("%s/", uuid)
		requestHandler += ":uuid/"
	}

	router := New()
	router.LogRequests = true

	for _, v := range dynamicRoutes {
		router.AddRegex(v.name, v.regex)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		none := GetParams("none", r)
		expect(t, 0, len(none))
		expect(t, GetParam("uuid", r, 1), uuids[1])
		expect(t, GetParam("none", r, 1), "")
		w.Write([]byte("named params"))
	}

	router.HandleFunc(requestHandler, handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", request, nil)
	router.ServeHTTP(w, req)
	expect(t, w.Code, 200)
}

func TestGetParamWildcard(t *testing.T) {
	router := New()
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		param := GetParam("*", r)
		expect(t, "test", param)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/foo/bar/xxxx", nil)
	router.ServeHTTP(w, req)
}

func TestGetParamsWildcard(t *testing.T) {
	router := New()
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		param := GetParams("*", r)
		expect(t, "test", param[0])
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/foo/bar/xxxx", nil)
	router.ServeHTTP(w, req)
}
