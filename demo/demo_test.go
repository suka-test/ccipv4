package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/suka-test/ccipv4"
)

func TestRun(t *testing.T) {
	var (
		got bytes.Buffer
	)
	c := &cli{
		stdout: &got,
		stderr: &got,
		db:     ccipv4.GetDB(),
	}
	// カントリーコード一覧のデータが不正
	c.cc = "AD|Andorra"
	c.stdin = strings.NewReader("q")
	if err := c.run(); err == nil {
		t.Error("run: invalid country code data, but no error")
	}

	// カントリーコード一覧のデータが正常
	c.cc = `AD||アンドラ
AE|United Arab Emirates|
AF|Afghanistan|アフガニスタン`
	c.stdin = strings.NewReader("q")
	if err := c.run(); err != nil {
		t.Errorf("run: valid country code data, but error: %v", err)
	}

}

func TestLoadIPBD(t *testing.T) {
	var (
		got bytes.Buffer
	)

	c := &cli{
		stdout: &got,
		stderr: &got,
		db:     ccipv4.GetDB(),
	}

	// テスト機で http://127.0.0.1 が
	// 動作していない場合のみ行う
	_, err := http.Get("http://127.0.0.1")
	if err != nil {
		// 指定された URL でエラー
		c.rir = [][]string{{"testRIR", "http://127.0.0.1"}}
		c.loadIPBD()
		if !strings.Contains(got.String(), "データベースを更新できませんでした。") {
			t.Errorf("getCommand: there are not the strings: %s", got.String())
		}
		got.Reset()
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/afrinic",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "afrinic|ZA|ipv4|41.0.0.0|2097152|20071126|allocated")
		},
	)
	mux.HandleFunc(
		"/apnic",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "apnic|AU|ipv4|1.0.0.0|256|20110811|assigned")
		},
	)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	c.rir = [][]string{
		{"afrinic", ts.URL + "/afrinic"},
		{"apnic", ts.URL + "/apnic"},
	}
	c.loadIPBD()
	if !strings.Contains(got.String(), "apnic   からデータをダウンロードします。 >>> 経過時間 : ") {
		t.Errorf("getCommand: there are not the strings: %s", got.String())
	}

}

func TestSearchIPB(t *testing.T) {
	var (
		got bytes.Buffer
	)

	c := &cli{
		stdout: &got,
		stderr: &got,
		db:     ccipv4.GetDB(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/afrinic",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "afrinic|ZA|ipv4|41.0.0.0|2097152|20071126|allocated")
		},
	)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	if err := c.db.LoadIPBDataByURL(ts.URL + "/afrinic"); err != nil {
		t.Fatalf("searchIPB: failed to load IPBData: %v", err)
	}
	c.db.SwitchIPBData()

	c.stdin = strings.NewReader("")
	c.searchIPB()
	if !strings.Contains(got.String(), "               IPアドレスではありません。\n") {
		t.Errorf("searchIPB: invalid message: %s", got.String())
	}
	got.Reset()

	c.stdin = strings.NewReader("127.0.0.1")
	c.searchIPB()
	if !strings.Contains(got.String(), "               ループバックアドレスです。\n") {
		t.Errorf("searchIPB: invalid message: %s", got.String())
	}
	got.Reset()

	c.stdin = strings.NewReader("41.2.0.10")
	c.searchIPB()
	if !strings.Contains(got.String(), "ブロック : 41.0.0.0 〜 41.31.255.255\n") {
		t.Errorf("searchIPB: invalid message: %s", got.String())
	}
	if !strings.Contains(got.String(), "コード   : ZA\n") {
		t.Errorf("searchIPB: invalid message: %s", got.String())
	}
}

func TestGetInfo(t *testing.T) {
	var (
		got bytes.Buffer
	)

	c := &cli{
		stdout: &got,
		stderr: &got,
		cc:     COUNTRY_CODES,
		db:     ccipv4.GetDB(),
	}

	if err := c.db.SetTmpCountryCodes(strings.NewReader(c.cc)); err != nil {
		t.Fatalf("getInfo: failed to set country codes: %v", err)
	}
	c.db.SwitchCCData()

	header1 := "                              #### 集計情報 ####"
	header2 := "国・地域の数 0\n\n"
	header3 := "コード|              日本語国・地域名              | ブロック数 | アドレス数 \n"
	line := "------|--------------------------------------------|------------|------------\n"
	footer := " 合計 |                                            |          0 |          0 \n"

	c.getInfo()

	if !strings.HasPrefix(got.String(), header1) {
		t.Errorf("getInfo: unexpected result: %s", got.String())
	}
	if !strings.Contains(got.String(), header2+header3+line+footer) {
		t.Errorf("getInfo: unexpected result: %s", got.String())
	}
	got.Reset()

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/a",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "afrinic|ZA|ipv4|41.0.0.0|2097152|20071126|allocated")
		},
	)
	mux.HandleFunc(
		"/b",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "afrinic|ZA|ipv4|41.0.0.0|2097152|20071126|allocated")
			fmt.Fprintln(w, "arin||ipv4|23.131.145.0|3840||reserved|")
		},
	)
	mux.HandleFunc(
		"/c",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "afrinic|ZZ|ipv4|41.57.112.0|2048||reserved|")
			fmt.Fprintln(w, "afrinic|ZA|ipv4|41.0.0.0|2097152|20071126|allocated")
			fmt.Fprintln(w, "arin||ipv4|23.131.145.0|3840||reserved|")
		},
	)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	if err := c.db.LoadIPBDataByURL(ts.URL + "/a"); err != nil {
		t.Fatalf("searchIPB: failed to load IPBData: %v", err)
	}
	c.db.SwitchIPBData()

	header2 = "国・地域の数 1\n\n"
	body1 := "  ZA  | 南アフリカ                                 |          1 |    2097152 \n"
	footer = " 合計 |                                            |          1 |    2097152 \n"

	c.getInfo()

	if !strings.HasPrefix(got.String(), header1) {
		t.Errorf("getInfo: unexpected result: %s", got.String())
	}
	if !strings.Contains(got.String(), header2+header3+line+body1+line+footer) {
		t.Errorf("getInfo: unexpected result: %s", got.String())
	}
	got.Reset()

	if err := c.db.LoadIPBDataByURL(ts.URL + "/b"); err != nil {
		t.Fatalf("searchIPB: failed to load IPBData: %v", err)
	}
	c.db.SwitchIPBData()

	body2 := "未割当|                                            |          1 |       3840 \n"
	footer = " 合計 |                                            |          2 |    2100992 \n"

	c.getInfo()

	if !strings.HasPrefix(got.String(), header1) {
		t.Errorf("getInfo: unexpected result: %s", got.String())
	}
	if !strings.Contains(got.String(), header2+header3+line+body2+body1+line+footer) {
		t.Errorf("getInfo: unexpected result: %s", got.String())
	}
	got.Reset()

	if err := c.db.LoadIPBDataByURL(ts.URL + "/c"); err != nil {
		t.Fatalf("searchIPB: failed to load IPBData: %v", err)
	}
	c.db.SwitchIPBData()

	body3 := "  ZZ  | 不明                                       |          1 |       2048 \n"
	footer = " 合計 |                                            |          3 |    2103040 \n"

	c.getInfo()

	if !strings.HasPrefix(got.String(), header1) {
		t.Errorf("getInfo: unexpected result: %s", got.String())
	}
	if !strings.Contains(got.String(), header2+header3+line+body2+body1+body3+line+footer) {
		t.Errorf("getInfo: unexpected result: %s", got.String())
	}
	got.Reset()

}
