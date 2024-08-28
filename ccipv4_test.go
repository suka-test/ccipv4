package ccipv4

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"os"
	"regexp"
	"strings"
	"testing"
)

// テスト用の RIR のダミーデータを取得する。
func getDummyRIR() ([]string, *http.ServeMux) {
	urlRIR := []string{
		"/afrinic",
		"/apnic",
		"/arin",
		"/lacnic",
		"/ripencc",
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		urlRIR[0],
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "afrinic|ZA|ipv4|41.0.0.0|2097152|20071126|allocated")
		},
	)
	mux.HandleFunc(
		urlRIR[1],
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "apnic|AU|ipv4|1.0.0.0|256|20110811|assigned")
		},
	)
	mux.HandleFunc(
		urlRIR[2],
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "arin|US|ipv4|2.57.164.0|1024|20230828|allocated|a015d0af44d434b219c6308d56e23f9e")
		},
	)
	mux.HandleFunc(
		urlRIR[3],
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "lacnic|DO|ipv4|5.183.80.0|1024|20240711|allocated")
		},
	)
	mux.HandleFunc(
		urlRIR[4],
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "ripencc|PS|ipv4|1.178.112.0|4096|20071126|allocated")
		},
	)

	return urlRIR, mux
}

func TestGetDB(t *testing.T) {
	// 初期状態を確認
	db := GetDB()
	if db.ib.data != nil {
		t.Errorf("GetDB: ib.data is invalid: %v", db.ib.data)
	}
	if db.ib.dicCCStrToInt != nil {
		t.Errorf("GetDB: ib.dicCCStrToInt is invalid: %v", db.ib.dicCCStrToInt)
	}
	if db.ib.dicCCIntToStr != nil {
		t.Errorf("GetDB: ib.dicCCIntToStr is invalid: %v", db.ib.dicCCIntToStr)
	}
	if db.ib.totalBlocks != nil {
		t.Errorf("GetDB: ib.totalBlocks is invalid: %v", db.ib.totalBlocks)
	}
	if db.ib.totalValue != nil {
		t.Errorf("GetDB: ib.totalValue is invalid: %v", db.ib.totalValue)
	}
	if db.cc.data != nil {
		t.Errorf("GetDB: cc.data is invalid: %v", db.cc.data)
	}
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("GetDB: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("GetDB: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("GetDB: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("GetDB: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("GetDB: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}
	if db.tmpCC.data == nil || len(db.tmpCC.data) != 0 {
		t.Errorf("GetDB: tmpCC.data is invalid: %v", db.tmpCC.data)
	}
	if db.reg.MatchString("aa") || db.reg.MatchString("BCD") || !db.reg.MatchString("YZ") {
		t.Error("GetDB: reg is invalid")
	}
	if db.urlRIR[0] != URL_DELEGATED_AFRINIC_EXTENDED_LATEST {
		t.Errorf("GetDB: urlRIR[0] is invalid: %s", db.urlRIR[0])
	}
	if db.urlRIR[1] != URL_DELEGATED_APNIC_EXTENDED_LATEST {
		t.Errorf("GetDB: urlRIR[1] is invalid: %s", db.urlRIR[1])
	}
	if db.urlRIR[2] != URL_DELEGATED_ARIN_EXTENDED_LATEST {
		t.Errorf("GetDB: urlRIR[2] is invalid: %s", db.urlRIR[2])
	}
	if db.urlRIR[3] != URL_DELEGATED_LACNIC_EXTENDED_LATEST {
		t.Errorf("GetDB: urlRIR[3] is invalid: %s", db.urlRIR[3])
	}
	if db.urlRIR[4] != URL_DELEGATED_RIPENCC_EXTENDED_LATEST {
		t.Errorf("GetDB: urlRIR[4] is invalid: %s", db.urlRIR[4])
	}

}

func TestClearTmpIPBData(t *testing.T) {
	// 初期状態からクリア
	db := GetDB()

	db.ClearTmpIPBData()
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("ClearTmpIPBData: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("ClearTmpIPBData: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("ClearTmpIPBData: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("ClearTmpIPBData: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("ClearTmpIPBData: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

	// データがある状態からクリア
	db.tmpIB.data[0] = map[uint8]map[uint8]map[uint8]block{}
	db.tmpIB.data[0][0] = map[uint8]map[uint8]block{}
	db.tmpIB.data[0][0][0] = map[uint8]block{}
	db.tmpIB.data[0][0][0][0] = block{value: 16, country: 0}
	db.tmpIB.dicCCIntToStr[0] = "JP"
	db.tmpIB.dicCCStrToInt["JP"] = 0

	db.ClearTmpIPBData()
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("ClearTmpIPBData: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("ClearTmpIPBData: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("ClearTmpIPBData: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("ClearTmpIPBData: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("ClearTmpIPBData: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

}

func TestSetTmpIPBlocks(t *testing.T) {
	// ファイルの形式が不正
	db := GetDB()

	fp, err := os.Open("testdata/invalidIPBlockFile-1")
	if err != nil {
		t.Fatalf("setTmpIPBlocks: can't read testdata/invalidIPBlockFile-1: %v", err)
	}
	err = db.setTmpIPBlocks(fp)
	if err == nil {
		t.Error("setTmpIPBlocks: read invalidIPBlockFile-1, but no error")
	}
	fp.Close()
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

	// ipv4 でフィールドの数が足りない
	db = GetDB()
	fp, err = os.Open("testdata/invalidIPBlockFile-2")
	if err != nil {
		t.Fatalf("setTmpIPBlocks: can't read testdata/invalidIPBlockFile-2: %v", err)
	}
	err = db.setTmpIPBlocks(fp)
	if err == nil {
		t.Error("setTmpIPBlocks: read invalidIPBlockFile-2, but no error")
	} else if !strings.Contains(err.Error(), "of the line's fields is invalid:") {
		t.Errorf("setTmpIPBlocks: unexpected error: %v", err)
	}
	fp.Close()
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

	// アドレスをパースできない
	db = GetDB()
	fp, err = os.Open("testdata/invalidIPBlockFile-3")
	if err != nil {
		t.Fatalf("setTmpIPBlocks: can't read testdata/invalidIPBlockFile-3: %v", err)
	}
	err = db.setTmpIPBlocks(fp)
	if err == nil {
		t.Error("setTmpIPBlocks: read invalidIPBlockFile-3, but no error")
	}
	fp.Close()
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

	// 正常データ１件
	db = GetDB()
	fp, err = os.Open("testdata/validIPBlockFile-1")
	if err != nil {
		t.Fatalf("setTmpIPBlocks: can't read testdata/validIPBlockFile-1: %v", err)
	}
	err = db.setTmpIPBlocks(fp)
	if err != nil {
		t.Errorf("setTmpIPBlocks: read validIPBlockFile-1: %v", err)
	}
	fp.Close()
	if len(db.tmpIB.data) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.data length want 1, but %d: %v", len(db.tmpIB.data), db.tmpIB.data)
	}
	if _, ok := db.tmpIB.data[114][48][0][0]; !ok {
		t.Error("setTmpIPBlocks: tmpIB.data[114][48][0][0] doesn't exist")
	} else {
		if db.tmpIB.data[114][48][0][0].country != 0 {
			t.Errorf("setTmpIPBlocks: tmpIB.data[114][48][0][0].country want 0, but got %d", db.tmpIB.data[114][48][0][0].country)
		}
		if db.tmpIB.data[114][48][0][0].value != 262144 {
			t.Errorf("setTmpIPBlocks: tmpIB.data[114][48][0][0].ipCidr want 262144, but got %d", db.tmpIB.data[114][48][0][0].value)
		}
	}
	if len(db.tmpIB.dicCCStrToInt) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt length want 1, but %d: %v", len(db.tmpIB.dicCCStrToInt), db.tmpIB.dicCCStrToInt)
	} else {
		if _, ok := db.tmpIB.dicCCStrToInt["JP"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.dicCCStrToInt[JP] doesn't exist")
		} else if db.tmpIB.dicCCStrToInt["JP"] != 0 {
			t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt[JP] want 0, but %d", db.tmpIB.dicCCStrToInt["JP"])
		}
	}
	if len(db.tmpIB.dicCCIntToStr) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr length want 1, but %d: %v", len(db.tmpIB.dicCCIntToStr), db.tmpIB.dicCCIntToStr)
	} else {
		if _, ok := db.tmpIB.dicCCIntToStr[0]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.dicCCIntToStr[0] doesn't exist")
		} else if db.tmpIB.dicCCIntToStr[0] != "JP" {
			t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr[0] want JP, but %s", db.tmpIB.dicCCIntToStr[0])
		}
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 2 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	} else {
		if _, ok := db.tmpIB.totalBlocks["ALL"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalBlocks[ALL] doesn't exist")
		} else if db.tmpIB.totalBlocks["ALL"] != 1 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks[ALL] want 1, but %d", db.tmpIB.totalBlocks["ALL"])
		}
		if _, ok := db.tmpIB.totalBlocks["JP"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalBlocks[JP] doesn't exist")
		} else if db.tmpIB.totalBlocks["JP"] != 1 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks[JP] want 1, but %d", db.tmpIB.totalBlocks["JP"])
		}
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 2 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	} else {
		if _, ok := db.tmpIB.totalValue["ALL"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalValue[ALL] doesn't exist")
		} else if db.tmpIB.totalValue["ALL"] != 262144 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalValue[ALL] want 262144, but %d", db.tmpIB.totalValue["ALL"])
		}
		if _, ok := db.tmpIB.totalValue["JP"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalValue[JP] doesn't exist")
		} else if db.tmpIB.totalValue["JP"] != 262144 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalValue[JP] want 262144, but %d", db.tmpIB.totalValue["JP"])
		}
	}

	// 各種データ混在
	db = GetDB()
	fp, err = os.Open("testdata/variousDataIPBlockFile")
	if err != nil {
		t.Fatalf("setTmpIPBlocks: can't read testdata/variousDataIPBlockFile: %v", err)
	}
	err = db.setTmpIPBlocks(fp)
	if err == nil {
		t.Error("setTmpIPBlocks: read variousDataIPBlockFile, but no error")
	}
	fp.Close()
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

	// DB をクリアせず、データ追加
	db = GetDB()
	fp, err = os.Open("testdata/validIPBlockFile-1")
	if err != nil {
		t.Fatalf("setTmpIPBlocks: can't read testdata/validIPBlockFile-1: %v", err)
	}
	err = db.setTmpIPBlocks(fp)
	if err != nil {
		t.Fatalf("setTmpIPBlocks: read testdata/validIPBlockFile-1, error: %v", err)
	}
	fp.Close()
	fp, err = os.Open("testdata/validIPBlockFile-2")
	if err != nil {
		t.Fatalf("setTmpIPBlocks: can't read testdata/validIPBlockFile-2: %v", err)
	}
	err = db.setTmpIPBlocks(fp)
	if err != nil {
		t.Errorf("setTmpIPBlocks: read validIPBlockFile-2: %v", err)
	}
	fp.Close()
	if _, ok := db.tmpIB.data[114][48][0][0]; !ok {
		t.Error("setTmpIPBlocks: tmpIB.data[114][48][0][0] doesn't exist")
	} else {
		if db.tmpIB.data[114][48][0][0].country != 0 {
			t.Errorf("setTmpIPBlocks: tmpIB.data[114][48][0][0].country want 0, but got %d", db.tmpIB.data[114][48][0][0].country)
		}
		if db.tmpIB.data[114][48][0][0].value != 131072 {
			t.Errorf("setTmpIPBlocks: tmpIB.data[114][48][0][0].value want 131072, but got %d", db.tmpIB.data[114][48][0][0].value)
		}
	}
	if _, ok := db.tmpIB.data[114][31][248][128]; !ok {
		t.Error("setTmpIPBlocks: tmpIB.data[114][31][248][0] doesn't exist")
	} else {
		if db.tmpIB.data[114][31][248][128].country != 1 {
			t.Errorf("setTmpIPBlocks: tmpIB.data[114][31][248][128].country want 0, but got %d", db.tmpIB.data[114][31][248][128].country)
		}
		if db.tmpIB.data[114][31][248][128].value != 2048 {
			t.Errorf("setTmpIPBlocks: tmpIB.data[114][31][248][128].value want 2048, but got %d", db.tmpIB.data[114][31][248][128].value)
		}
	}
	if _, ok := db.tmpIB.data[124][147][128][0]; !ok {
		t.Error("setTmpIPBlocks: tmpIB.data[124][147][128][0] doesn't exist")
	} else {
		if db.tmpIB.data[124][147][128][0].country != 2 {
			t.Errorf("setTmpIPBlocks: tmpIB.data[124][147][128][0].country want 2, but got %d", db.tmpIB.data[124][147][128][0].country)
		}
		if db.tmpIB.data[124][147][128][0].value != 32768 {
			t.Errorf("setTmpIPBlocks: tmpIB.data[124][147][128][0].value want 32768, but got %d", db.tmpIB.data[124][147][128][0].value)
		}
	}
	if len(db.tmpIB.dicCCStrToInt) != 3 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt length want 3, but %d: %v", len(db.tmpIB.dicCCStrToInt), db.tmpIB.dicCCStrToInt)
	} else {
		if _, ok := db.tmpIB.dicCCStrToInt["JP"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.dicCCStrToInt[JP] doesn't exist")
		} else if db.tmpIB.dicCCStrToInt["JP"] != 0 {
			t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt[JP] want 0, but %d", db.tmpIB.dicCCStrToInt["JP"])
		}
		if _, ok := db.tmpIB.dicCCStrToInt["IN"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.dicCCStrToInt[IN] doesn't exist")
		} else if db.tmpIB.dicCCStrToInt["IN"] != 1 {
			t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt[IN] want 1, but %d", db.tmpIB.dicCCStrToInt["IN"])
		}
		if _, ok := db.tmpIB.dicCCStrToInt["CN"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.dicCCStrToInt[CN] doesn't exist")
		} else if db.tmpIB.dicCCStrToInt["CN"] != 2 {
			t.Errorf("setTmpIPBlocks: tmpIB.dicCCStrToInt[CN] want 2, but %d", db.tmpIB.dicCCStrToInt["CN"])
		}
	}
	if len(db.tmpIB.dicCCIntToStr) != 3 {
		t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr length want 5, but %d: %v", len(db.tmpIB.dicCCIntToStr), db.tmpIB.dicCCIntToStr)
	} else {
		if _, ok := db.tmpIB.dicCCIntToStr[0]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.dicCCIntToStr[0] doesn't exist")
		} else if db.tmpIB.dicCCIntToStr[0] != "JP" {
			t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr[0] want JP, but %s", db.tmpIB.dicCCIntToStr[0])
		}
		if _, ok := db.tmpIB.dicCCIntToStr[1]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.dicCCIntToStr[1] doesn't exist")
		} else if db.tmpIB.dicCCIntToStr[1] != "IN" {
			t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr[1] want IN, but %s", db.tmpIB.dicCCIntToStr[1])
		}
		if _, ok := db.tmpIB.dicCCIntToStr[2]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.dicCCIntToStr[2] doesn't exist")
		} else if db.tmpIB.dicCCIntToStr[2] != "CN" {
			t.Errorf("setTmpIPBlocks: tmpIB.dicCCIntToStr[2] want CN, but %s", db.tmpIB.dicCCIntToStr[2])
		}
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 4 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	} else {
		if _, ok := db.tmpIB.totalBlocks["ALL"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalBlocks[ALL] doesn't exist")
		} else if db.tmpIB.totalBlocks["ALL"] != 3 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks[ALL] want 3, but %d", db.tmpIB.totalBlocks["ALL"])
		}
		if _, ok := db.tmpIB.totalBlocks["JP"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalBlocks[JP] doesn't exist")
		} else if db.tmpIB.totalBlocks["JP"] != 1 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks[JP] want 1, but %d", db.tmpIB.totalBlocks["JP"])
		}
		if _, ok := db.tmpIB.totalBlocks["IN"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalBlocks[IN] doesn't exist")
		} else if db.tmpIB.totalBlocks["IN"] != 1 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks[IN] want 1, but %d", db.tmpIB.totalBlocks["IN"])
		}
		if _, ok := db.tmpIB.totalBlocks["CN"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalBlocks[CN] doesn't exist")
		} else if db.tmpIB.totalBlocks["CN"] != 1 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalBlocks[CN] want 1, but %d", db.tmpIB.totalBlocks["CN"])
		}
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 4 {
		t.Errorf("setTmpIPBlocks: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	} else {
		if _, ok := db.tmpIB.totalValue["ALL"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalValue[ALL] doesn't exist")
		} else if db.tmpIB.totalValue["ALL"] != 165888 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalValue[ALL] want 165888, but %d", db.tmpIB.totalValue["ALL"])
		}
		if _, ok := db.tmpIB.totalValue["JP"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalValue[JP] doesn't exist")
		} else if db.tmpIB.totalValue["JP"] != 131072 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalValue[JP] want 131072, but %d", db.tmpIB.totalValue["JP"])
		}
		if _, ok := db.tmpIB.totalValue["IN"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalValue[IN] doesn't exist")
		} else if db.tmpIB.totalValue["IN"] != 2048 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalValue[IN] want 2048, but %d", db.tmpIB.totalValue["IN"])
		}
		if _, ok := db.tmpIB.totalValue["CN"]; !ok {
			t.Error("setTmpIPBlocks: tmpIB.totalValue[CN] doesn't exist")
		} else if db.tmpIB.totalValue["CN"] != 32768 {
			t.Errorf("setTmpIPBlocks: tmpIB.totalValue[CN] want 32768, but %d", db.tmpIB.totalValue["CN"])
		}
	}
}

func TestSetTmpCountryCodes(t *testing.T) {
	// ファイルの形式が不正
	db := GetDB()
	fp, err := os.Open("testdata/invalidCountryCodeFile-1")
	if err != nil {
		t.Fatalf("SetTmpCountryCodes: can't read testdata/invalidCountryCodeFile-1: %v", err)
	}
	err = db.SetTmpCountryCodes(fp)
	if err == nil {
		t.Error("SetTmpCountryCodes: read invalidCountryCodeFile-1, but no error")
	}
	fp.Close()
	if len(db.tmpCC.data) != 0 {
		t.Errorf("SetTmpCountryCodes: tmpCC.data length want 0, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	}

	// フィールドの数が３ではない
	db = GetDB()
	fp, err = os.Open("testdata/invalidCountryCodeFile-2")
	if err != nil {
		t.Fatalf("SetTmpCountryCodes: can't read testdata/invalidCountryCodeFile-2: %v", err)
	}
	err = db.SetTmpCountryCodes(fp)
	if err == nil {
		t.Error("SetTmpCountryCodes: read invalidCountryCodeFile-2, but no error")
	}
	fp.Close()
	if len(db.tmpCC.data) != 0 {
		t.Errorf("SetTmpCountryCodes: tmpCC.data length want 0, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	}

	// カントリーコードが英大文字ではない
	db = GetDB()
	fp, err = os.Open("testdata/invalidCountryCodeFile-3")
	if err != nil {
		t.Fatalf("SetTmpCountryCodes: can't read testdata/invalidCountryCodeFile-3: %v", err)
	}
	err = db.SetTmpCountryCodes(fp)
	if err == nil {
		t.Error("SetTmpCountryCodes: read invalidCountryCodeFile-3, but no error")
	}
	fp.Close()
	if len(db.tmpCC.data) != 0 {
		t.Errorf("SetTmpCountryCodes: tmpCC.data length want 0, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	}

	// カントリーコードが２文字ではない
	db = GetDB()
	fp, err = os.Open("testdata/invalidCountryCodeFile-4")
	if err != nil {
		t.Fatalf("SetTmpCountryCodes: can't read testdata/invalidCountryCodeFile-4: %v", err)
	}
	err = db.SetTmpCountryCodes(fp)
	if err == nil {
		t.Error("SetTmpCountryCodes: read invalidCountryCodeFile-4, but no error")
	}
	fp.Close()
	if len(db.tmpCC.data) != 0 {
		t.Errorf("SetTmpCountryCodes: tmpCC.data length want 0, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	}

	// 正常データ１件
	db = GetDB()
	fp, err = os.Open("testdata/validCountryCodeFile-1")
	if err != nil {
		t.Fatalf("SetTmpCountryCodes: can't read testdata/validCountryCodeFile-1: %v", err)
	}
	err = db.SetTmpCountryCodes(fp)
	if err != nil {
		t.Errorf("SetTmpCountryCodes: read validCountryCodeFile-1, but error: %v", err)
	}
	fp.Close()
	if len(db.tmpCC.data) != 1 {
		t.Errorf("SetTmpCountryCodes: tmpCC.data length want 1, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	} else if _, ok := db.tmpCC.data["AD"]; !ok {
		t.Error("SetTmpCountryCodes: tmpCC.data[AD] doesn't exist")
	} else {
		if db.tmpCC.data["AD"].Name != "Andorra" {
			t.Errorf("SetTmpCountryCodes: tmpCC.data[AD].eName want Andorra, but %s", db.tmpCC.data["AD"].Name)
		}
		if db.tmpCC.data["AD"].AltName != "アンドラ" {
			t.Errorf("SetTmpCountryCodes: tmpCC.data[AD].aName want アンドラ, but %s", db.tmpCC.data["AD"].AltName)
		}
	}

	// 各種データ混在
	db = GetDB()
	fp, err = os.Open("testdata/variousDataCountryCodeFile")
	if err != nil {
		t.Fatalf("SetTmpCountryCodes: can't read testdata/variousDataCountryCodeFile: %v", err)
	}
	err = db.SetTmpCountryCodes(fp)
	if err == nil {
		t.Error("SetTmpCountryCodes: read variousDataCountryCodeFile, but no error")
	}
	fp.Close()
	if len(db.tmpCC.data) != 0 {
		t.Errorf("SetTmpCountryCodes: tmpCC.data length want 0, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	}

	// DB をクリアせず、データ追加
	db = GetDB()
	fp, err = os.Open("testdata/validCountryCodeFile-1")
	if err != nil {
		t.Fatalf("SetTmpCountryCodes: can't read testdata/validCountryCodeFile-1: %v", err)
	}
	err = db.SetTmpCountryCodes(fp)
	if err != nil {
		t.Errorf("SetTmpCountryCodes: read validCountryCodeFile-1, but error: %v", err)
	}
	fp.Close()
	fp, err = os.Open("testdata/validCountryCodeFile-2")
	if err != nil {
		t.Fatalf("SetTmpCountryCodes: can't read testdata/validCountryCodeFile-2: %v", err)
	}
	err = db.SetTmpCountryCodes(fp)
	if err != nil {
		t.Errorf("SetTmpCountryCodes: read validCountryCodeFile-1, but error: %v", err)
	}
	fp.Close()
	if len(db.tmpCC.data) != 2 {
		t.Errorf("SetTmpCountryCodes: tmpCC.data length want 2, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	} else {
		if _, ok := db.tmpCC.data["AD"]; !ok {
			t.Error("SetTmpCountryCodes: tmpCC.data[AD] doesn't exist")
		} else {
			if db.tmpCC.data["AD"].Name != "Principality of Andorra" {
				t.Errorf("SetTmpCountryCodes: tmpCC.data[AD].eName want Principality of Andorra, but %s", db.tmpCC.data["AD"].Name)
			}
			if db.tmpCC.data["AD"].AltName != "アンドラ公国" {
				t.Errorf("SetTmpCountryCodes: tmpCC.data[AD].aName want アンドラ公国, but %s", db.tmpCC.data["AD"].AltName)
			}
		}
		if _, ok := db.tmpCC.data["AE"]; !ok {
			t.Error("SetTmpCountryCodes: tmpCC.data[AE] doesn't exist")
		} else {
			if db.tmpCC.data["AE"].Name != "United Arab Emirates" {
				t.Errorf("SetTmpCountryCodes: tmpCC.data[AE].eName want United Arab Emirates, but %s", db.tmpCC.data["AE"].Name)
			}
			if db.tmpCC.data["AE"].AltName != "アラブ首長国連邦" {
				t.Errorf("SetTmpCountryCodes: tmpCC.data[AE].aName want アラブ首長国連邦, but %s", db.tmpCC.data["AE"].AltName)
			}
		}
	}
}

func TestLoadIPBDataByFile(t *testing.T) {
	// ファイルがない
	db := GetDB()
	err := db.LoadIPBDataByFile("testdata/nothing1")
	if err == nil {
		t.Error("LoadIPBDataByFile: file nothing1, but no error")
	}
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("LoadIPBDataByFile: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("LoadIPBDataByFile: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("LoadIPBDataByFile: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("LoadIPBDataByFile: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("LoadIPBDataByFile: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

	// ファイルにエラーがある
	db = GetDB()
	err = db.LoadIPBDataByFile("testdata/invalidIPBlockFile-1")
	if err == nil {
		t.Error("LoadIPBDataByFile: file invalidIPBlockFile-1, but no error")
	}
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("LoadIPBDataByFile: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("LoadIPBDataByFile: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("LoadIPBDataByFile: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("LoadIPBDataByFile: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("LoadIPBDataByFile: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

	// 正常終了
	db = GetDB()
	err = db.LoadIPBDataByFile("testdata/validIPBlockFile-1")
	if err != nil {
		t.Errorf("LoadIPBDataByFile: load validIPBlockFile-1, but error: %v", err)
	}
	if db.tmpIB.data == nil || len(db.tmpIB.data) != 1 {
		t.Errorf("LoadIPBDataByFile: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 1 {
		t.Errorf("LoadIPBDataByFile: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 1 {
		t.Errorf("LoadIPBDataByFile: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 2 {
		t.Errorf("LoadIPBDataByFile: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	} else {
		if _, ok := db.tmpIB.totalBlocks["ALL"]; !ok {
			t.Error("LoadIPBDataByFile: tmpIB.totalBlocks[ALL] doesn't exist")
		} else if db.tmpIB.totalBlocks["ALL"] != 1 {
			t.Errorf("LoadIPBDataByFile: tmpIB.totalBlocks[ALL] want 1, but %d", db.tmpIB.totalBlocks["ALL"])
		}
		if _, ok := db.tmpIB.totalBlocks["JP"]; !ok {
			t.Error("LoadIPBDataByFile: tmpIB.totalBlocks[JP] doesn't exist")
		} else if db.tmpIB.totalBlocks["JP"] != 1 {
			t.Errorf("LoadIPBDataByFile: tmpIB.totalBlocks[JP] want 1, but %d", db.tmpIB.totalBlocks["JP"])
		}
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 2 {
		t.Errorf("LoadIPBDataByFile: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	} else {
		if _, ok := db.tmpIB.totalValue["ALL"]; !ok {
			t.Error("LoadIPBDataByFile: tmpIB.totalValue[ALL] doesn't exist")
		} else if db.tmpIB.totalValue["ALL"] != 262144 {
			t.Errorf("LoadIPBDataByFile: tmpIB.totalValue[ALL] want 262144, but %d", db.tmpIB.totalValue["ALL"])
		}
		if _, ok := db.tmpIB.totalValue["JP"]; !ok {
			t.Error("LoadIPBDataByFile: tmpIB.totalValue[JP] doesn't exist")
		} else if db.tmpIB.totalValue["JP"] != 262144 {
			t.Errorf("LoadIPBDataByFile: tmpIB.totalValue[JP] want 262144, but %d", db.tmpIB.totalValue["JP"])
		}
	}
}

func TestSwitchIPBData(t *testing.T) {
	// 一時保存用データベースに値を設定
	db := GetDB()

	db.tmpIB.data = map[uint8]map[uint8]map[uint8]map[uint8]block{}
	db.tmpIB.data[0] = map[uint8]map[uint8]map[uint8]block{}
	db.tmpIB.data[0][0] = map[uint8]map[uint8]block{}
	db.tmpIB.data[0][0][0] = map[uint8]block{}
	db.tmpIB.data[0][0][0][0] = block{value: 16, country: 0}
	db.tmpIB.dicCCIntToStr = map[uint8]string{}
	db.tmpIB.dicCCIntToStr[0] = "JP"
	db.tmpIB.dicCCStrToInt = map[string]uint8{}
	db.tmpIB.dicCCStrToInt["JP"] = 0
	db.tmpIB.totalBlocks = map[string]int{"ALL": 1, "JP": 1}
	db.tmpIB.totalValue = map[string]int{"ALL": 16, "JP": 16}

	db.SwitchIPBData()

	if db.tmpIB.data == nil || len(db.tmpIB.data) != 0 {
		t.Errorf("SwitchIPBData: tmpIB.data is invalid: %v", db.tmpIB.data)
	}
	if db.tmpIB.dicCCStrToInt == nil || len(db.tmpIB.dicCCStrToInt) != 0 {
		t.Errorf("SwitchIPBData: tmpIB.dicCCStrToInt is invalid: %v", db.tmpIB.dicCCStrToInt)
	}
	if db.tmpIB.dicCCIntToStr == nil || len(db.tmpIB.dicCCIntToStr) != 0 {
		t.Errorf("SwitchIPBData: tmpIB.dicCCIntToStr is invalid: %v", db.tmpIB.dicCCIntToStr)
	}
	if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 1 {
		t.Errorf("SwitchIPBData: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
	}
	if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 1 {
		t.Errorf("SwitchIPBData: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
	}

	if len(db.ib.data) != 1 {
		t.Errorf("SwitchIPBData: db.ib.data want 1, but %d: %v", len(db.ib.data), db.ib.data)
	} else if _, ok := db.ib.data[0][0][0][0]; !ok {
		t.Error("SwitchIPBData: db.ib.data[0][0][0][0] doesn't exist")
	} else {
		if db.ib.data[0][0][0][0].value != 16 {
			t.Errorf("SwitchIPBData: db.ib.data[0][0][0][0].ipCidr want 16, but %d", db.ib.data[0][0][0][0].value)
		}
		if db.ib.data[0][0][0][0].country != 0 {
			t.Errorf("SwitchIPBData: db.ib.data[0][0][0][0].ipCidr want 16, but %d", db.ib.data[0][0][0][0].value)
		}
	}
	if len(db.ib.dicCCIntToStr) != 1 {
		t.Errorf("SwitchIPBData: db.ib.dicCCIntToStr want 1, but %d: %v", len(db.ib.dicCCIntToStr), db.ib.dicCCIntToStr)
	} else if _, ok := db.ib.dicCCIntToStr[0]; !ok {
		t.Error("SwitchIPBData: db.ib.dicCCIntToStr[0] doesn't exist")
	} else if db.ib.dicCCIntToStr[0] != "JP" {
		t.Errorf("SwitchIPBData: db.ib.dicCCIntToStr[0] want JP, but %s", db.ib.dicCCIntToStr[0])
	}
	if len(db.ib.dicCCStrToInt) != 1 {
		t.Errorf("SwitchIPBData: db.ib.dicCCIntToStr want 1, but %d: %v", len(db.ib.dicCCStrToInt), db.ib.dicCCStrToInt)
	} else if _, ok := db.ib.dicCCStrToInt["JP"]; !ok {
		t.Error("SwitchIPBData: db.ib.dicCCStrToInt[0] doesn't exist")
	} else if db.ib.dicCCStrToInt["JP"] != 0 {
		t.Errorf("SwitchIPBData: db.ib.dicCCStrToInt[JP] want JP, but %d", db.ib.dicCCStrToInt["JP"])
	}
	if db.ib.totalBlocks == nil || len(db.ib.totalBlocks) != 2 {
		t.Errorf("SwitchIPBData: ib.totalBlocks is invalid: %v", db.ib.totalBlocks)
	} else {
		if _, ok := db.ib.totalBlocks["ALL"]; !ok {
			t.Error("SwitchIPBData: ib.totalBlocks[ALL] doesn't exist")
		} else if db.ib.totalBlocks["ALL"] != 1 {
			t.Errorf("SwitchIPBData: ib.totalBlocks[ALL] want 1, but %d", db.ib.totalBlocks["ALL"])
		}
		if _, ok := db.ib.totalBlocks["JP"]; !ok {
			t.Error("SwitchIPBData: ib.totalBlocks[JP] doesn't exist")
		} else if db.ib.totalBlocks["JP"] != 1 {
			t.Errorf("SwitchIPBData: ib.totalBlocks[JP] want 1, but %d", db.ib.totalBlocks["JP"])
		}
	}
	if db.ib.totalValue == nil || len(db.ib.totalValue) != 2 {
		t.Errorf("SwitchIPBData: ib.totalValue is invalid: %v", db.ib.totalValue)
	} else {
		if _, ok := db.ib.totalValue["ALL"]; !ok {
			t.Error("SwitchIPBData: ib.totalValue[ALL] doesn't exist")
		} else if db.ib.totalValue["ALL"] != 16 {
			t.Errorf("SwitchIPBData: ib.totalValue[ALL] want 16, but %d", db.ib.totalValue["ALL"])
		}
		if _, ok := db.ib.totalValue["JP"]; !ok {
			t.Error("SwitchIPBData: ib.totalValue[JP] doesn't exist")
		} else if db.ib.totalValue["JP"] != 16 {
			t.Errorf("SwitchIPBData: ib.totalValue[JP] want 16, but %d", db.ib.totalValue["JP"])
		}
	}
}

func TestInitCCDataByFile(t *testing.T) {
	// カントリーコードの一覧ファイルがない
	db := GetDB()
	err := db.InitCCDataByFile("testdata/nothing")
	if err == nil {
		t.Error("InitCCDataByFile: file nothing, but no error")
	}
	if len(db.cc.data) != 0 {
		t.Errorf("InitCCDataByFile: db.cc.data want 0, but %d: %v", len(db.cc.data), db.cc.data)
	}
	if len(db.tmpCC.data) != 0 {
		t.Errorf("InitCCDataByFile: tmpCC.data length want 0, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	}

	// カントリーコードの一覧ファイルにエラーがある
	db = GetDB()
	err = db.InitCCDataByFile("testdata/invalidCountryCodeFile-1")
	if err == nil {
		t.Error("InitCCDataByFile: file invalidCountryCodeFile-1, but no error")
	}
	if len(db.cc.data) != 0 {
		t.Errorf("InitCCDataByFile: db.cc.data want 0, but %d: %v", len(db.cc.data), db.cc.data)
	}
	if len(db.tmpCC.data) != 0 {
		t.Errorf("InitCCDataByFile: tmpCC.data length want 0, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	}

	// 正常終了
	db = GetDB()
	err = db.InitCCDataByFile("testdata/validCountryCodeFile-1")
	if err != nil {
		t.Errorf("InitCCDataByFile: file validCountryCodeFile-1, but error: %v", err)
	}
	if len(db.cc.data) != 1 {
		t.Errorf("InitCCDataByFile: tmpCC.data length want 1, but %d: %v", len(db.cc.data), db.cc.data)
	} else if _, ok := db.cc.data["AD"]; !ok {
		t.Error("InitCCDataByFile: tmpCC.data[AD] doesn't exist")
	} else {
		if db.cc.data["AD"].Name != "Andorra" {
			t.Errorf("InitCCDataByFile: tmpCC.data[AD].eName want Andorra, but %s", db.cc.data["AD"].Name)
		}
		if db.cc.data["AD"].AltName != "アンドラ" {
			t.Errorf("InitCCDataByFile: tmpCC.data[AD].aName want アンドラ, but %s", db.cc.data["AD"].AltName)
		}
	}
	if len(db.tmpCC.data) != 0 {
		t.Errorf("InitCCDataByFile: tmpCC.data length want 0, but %d: %v", len(db.tmpCC.data), db.tmpCC.data)
	}
}

func TestSwitchCCData(t *testing.T) {
	// 一時保存用データベースに値を設定
	db := GetDB()

	db.tmpCC.data = map[string]CountryCodeInfo{
		"JP": {
			Name:    "Japan",
			AltName: "日本",
		},
	}

	db.SwitchCCData()

	if db.tmpCC.data == nil || len(db.tmpCC.data) != 0 {
		t.Errorf("SwitchCCData: tmpCC.data is invalid: %v", db.tmpCC.data)
	}

	if len(db.cc.data) != 1 {
		t.Errorf("SwitchCCData: db.cc.data length want 1, but %d: %v", len(db.cc.data), db.cc.data)
	} else if _, ok := db.cc.data["JP"]; !ok {
		t.Error("SwitchCCData: db.cc.data[JP] doesn't exist")
	} else {
		if db.cc.data["JP"].Name != "Japan" {
			t.Errorf("SwitchCCData: db.cc.data[JP].ipCidr want Japan, but %s", db.cc.data["JP"].Name)
		}
		if db.cc.data["JP"].AltName != "日本" {
			t.Errorf("SwitchCCData: db.cc.data[JP].ipCidr want 16, but %s", db.cc.data["JP"].AltName)
		}
	}
}

func TestCheckFirst8Bit(t *testing.T) {
	db := GetDB()
	err := db.LoadIPBDataByFile("testdata/validIPBlockFile-2")
	if err != nil {
		t.Fatalf("checkFirst8Bit: load validIPBlockFile-2, but error: %v", err)
	}
	db.SwitchIPBData()

	// 初期値一致
	a4b, found := db.checkFirst8Bit([4]byte{114, 48, 0, 0})
	if !found {
		t.Errorf("checkFirst8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{114, 48, 0, 0} {
		t.Errorf("checkFirst8Bit: not [4]byte{114, 48, 0, 0}: %v", a4b)
	}

	// 初期値不一致・最初の8ビットを1づつ減少させて一致
	a4b, found = db.checkFirst8Bit([4]byte{123, 48, 0, 0})
	if !found {
		t.Errorf("checkFirst8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{114, 255, 255, 255} {
		t.Errorf("checkFirst8Bit: not [4]byte{114, 255, 255, 255}: %v", a4b)
	}

	// 初期値不一致・最初の8ビットを1づつ減少させても不一致
	a4b, found = db.checkFirst8Bit([4]byte{113, 48, 0, 0})
	if found {
		t.Errorf("checkFirst8Bit: found: %v", a4b)
	} else if a4b != [4]byte{113, 48, 0, 0} {
		t.Errorf("checkFirst8Bit: not [4]byte{113, 48, 0, 0}: %v", a4b)
	}
}

func TestCheckSecond8Bit(t *testing.T) {
	db := GetDB()
	err := db.LoadIPBDataByFile("testdata/validIPBlockFile-3")
	if err != nil {
		t.Fatalf("checkSecond8Bit: load validIPBlockFile-3, but error: %v", err)
	}
	db.SwitchIPBData()

	// 初期値一致
	a4b, found := db.checkSecond8Bit([4]byte{114, 48, 0, 0})
	if !found {
		t.Errorf("checkSecond8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{114, 48, 0, 0} {
		t.Errorf("checkSecond8Bit: not [4]byte{114, 48, 0, 0}: %v", a4b)
	}

	// 初期値不一致・データベースの2番めの8ビット0のみ
	a4b, found = db.checkSecond8Bit([4]byte{49, 111, 0, 0})
	if !found {
		t.Errorf("checkSecond8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{49, 0, 255, 255} {
		t.Errorf("checkSecond8Bit: not [4]byte{49, 0, 255, 255}: %v", a4b)
	}

	// 初期値不一致・2番めの8ビットを1づつ減少させて一致
	a4b, found = db.checkSecond8Bit([4]byte{114, 128, 0, 0})
	if !found {
		t.Errorf("checkSecond8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{114, 48, 255, 255} {
		t.Errorf("checkSecond8Bit: not [4]byte{114, 48, 255, 255}: %v", a4b)
	}

	// 初期値不一致・2番めの8ビットを1づつ減少させても不一致
	a4b, found = db.checkSecond8Bit([4]byte{114, 20, 0, 0})
	if found {
		t.Errorf("checkSecond8Bit: found: %v", a4b)
	} else if a4b != [4]byte{114, 20, 0, 0} {
		t.Errorf("checkSecond8Bit: not [4]byte{114, 20, 0, 0}: %v", a4b)
	}
}

func TestCheckThird8Bit(t *testing.T) {
	db := GetDB()
	err := db.LoadIPBDataByFile("testdata/validIPBlockFile-3")
	if err != nil {
		t.Fatalf("checkThird8Bit: load validIPBlockFile-3, but error: %v", err)
	}
	db.SwitchIPBData()

	// 初期値一致
	a4b, found := db.checkThird8Bit([4]byte{114, 48, 0, 0})
	if !found {
		t.Errorf("checkThird8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{114, 48, 0, 0} {
		t.Errorf("checkThird8Bit: not [4]byte{114, 48, 0, 0}: %v", a4b)
	}

	// 初期値不一致・データベースの3番めの8ビット0のみ
	a4b, found = db.checkThird8Bit([4]byte{49, 0, 111, 0})
	if !found {
		t.Errorf("checkThird8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{49, 0, 0, 255} {
		t.Errorf("checkThird8Bit: not [4]byte{49, 0, 0, 255}: %v", a4b)
	}

	// 初期値不一致・3番めの8ビットを1づつ減少させて一致
	a4b, found = db.checkThird8Bit([4]byte{114, 31, 255, 0})
	if !found {
		t.Errorf("checkThird8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{114, 31, 248, 255} {
		t.Errorf("checkThird8Bit: not [4]byte{114, 31, 248, 255}: %v", a4b)
	}

	// 初期値不一致・3番めの8ビットを1づつ減少させても不一致
	a4b, found = db.checkThird8Bit([4]byte{114, 31, 128, 0})
	if found {
		t.Errorf("checkThird8Bit: found: %v", a4b)
	} else if a4b != [4]byte{114, 31, 128, 0} {
		t.Errorf("checkThird8Bit: not [4]byte{114, 31, 128, 0}: %v", a4b)
	}
}

func TestCheckLast8Bit(t *testing.T) {
	db := GetDB()
	err := db.LoadIPBDataByFile("testdata/validIPBlockFile-3")
	if err != nil {
		t.Fatalf("checkLast8Bit: load validIPBlockFile-3, but error: %v", err)
	}
	db.SwitchIPBData()

	// 初期値一致
	a4b, found := db.checkLast8Bit([4]byte{114, 48, 0, 0})
	if !found {
		t.Errorf("checkLast8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{114, 48, 0, 0} {
		t.Errorf("checkLast8Bit: not [4]byte{114, 48, 0, 0}: %v", a4b)
	}

	// 初期値不一致・データベースの最後の8ビット0のみ
	a4b, found = db.checkLast8Bit([4]byte{49, 0, 0, 240})
	if !found {
		t.Errorf("checkLast8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{49, 0, 0, 0} {
		t.Errorf("checkLast8Bit: not [4]byte{49, 0, 0, 0}: %v", a4b)
	}

	// 初期値不一致・最後の8ビットを1づつ減少させて一致
	a4b, found = db.checkLast8Bit([4]byte{114, 31, 248, 182})
	if !found {
		t.Errorf("checkLast8Bit: not found: %v", a4b)
	} else if a4b != [4]byte{114, 31, 248, 128} {
		t.Errorf("checkLast8Bit: not [4]byte{114, 31, 248, 128}: %v", a4b)
	}

	// 初期値不一致・最後の8ビットを1づつ減少させても不一致
	a4b, found = db.checkLast8Bit([4]byte{114, 31, 248, 60})
	if found {
		t.Errorf("checkLast8Bit: found: %v", a4b)
	} else if a4b != [4]byte{114, 31, 248, 60} {
		t.Errorf("checkLast8Bit: not [4]byte{114, 31, 128, 60}: %v", a4b)
	}
}

func TestSearchBlockStart(t *testing.T) {
	// 全てのifを通過する。
	db := GetDB()
	err := db.LoadIPBDataByFile("testdata/validIPBlockFile-4")
	if err != nil {
		t.Fatalf("searchBlockStart: load validIPBlockFile-4, but error: %v", err)
	}
	db.SwitchIPBData()

	target, err := netip.ParseAddr("0.0.0.64")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b := db.searchBlockStart(target)
	if a4b != [4]byte{} {
		t.Errorf("searchBlockStart: not [4]byte{}: %v", a4b)
	}

	target, err = netip.ParseAddr("2.0.0.10")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{0, 0, 0, 128} {
		t.Errorf("searchBlockStart: not [4]byte{0, 0, 0, 128}: %v", a4b)
	}

	target, err = netip.ParseAddr("114.48.0.16")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{2, 0, 0, 64} {
		t.Errorf("searchBlockStart: not [4]byte{2, 0, 0, 64}: %v", a4b)
	}

	target, err = netip.ParseAddr("124.147.128.8")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{114, 48, 0, 32} {
		t.Errorf("searchBlockStart: not [4]byte{2, 0, 0, 64}: %v", a4b)
	}

	err = db.LoadIPBDataByFile("testdata/validIPBlockFile-5")
	if err != nil {
		t.Fatalf("searchBlockStart: load validIPBlockFile-5, but error: %v", err)
	}
	db.SwitchIPBData()

	target, err = netip.ParseAddr("0.0.100.10")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{} {
		t.Errorf("searchBlockStart: not [4]byte{}: %v", a4b)
	}

	target, err = netip.ParseAddr("2.0.0.10")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{0, 0, 128, 64} {
		t.Errorf("searchBlockStart: not [4]byte{0, 0, 128, 64}: %v", a4b)
	}

	target, err = netip.ParseAddr("114.48.8.16")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{2, 0, 0, 64} {
		t.Errorf("searchBlockStart: not [4]byte{2, 0, 0, 64}: %v", a4b)
	}

	err = db.LoadIPBDataByFile("testdata/validIPBlockFile-6")
	if err != nil {
		t.Fatalf("searchBlockStart: load validIPBlockFile-6, but error: %v", err)
	}
	db.SwitchIPBData()

	target, err = netip.ParseAddr("0.6.100.10")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{} {
		t.Errorf("searchBlockStart: not [4]byte{}: %v", a4b)
	}
	target, err = netip.ParseAddr("2.0.0.10")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{0, 8, 128, 64} {
		t.Errorf("searchBlockStart: not [4]byte{0, 8, 128, 64}: %v", a4b)
	}

	err = db.LoadIPBDataByFile("testdata/validIPBlockFile-7")
	if err != nil {
		t.Fatalf("searchBlockStart: load validIPBlockFile-7, but error: %v", err)
	}
	db.SwitchIPBData()

	target, err = netip.ParseAddr("0.0.0.1")
	if err != nil {
		t.Fatalf("searchBlockStart: failed to ParseAddr: %v", err)
	}
	a4b = db.searchBlockStart(target)
	if a4b != [4]byte{} {
		t.Errorf("searchBlockStart: not [4]byte{}: %v", a4b)
	}
}

func TestSearchInfo(t *testing.T) {
	db := GetDB()
	err := db.LoadIPBDataByFile("testdata/validIPBlockFile-1")
	if err != nil {
		t.Fatalf("SearchInfo: load validIPBlockFile-1, but error: %v", err)
	}
	db.SwitchIPBData()
	err = db.InitCCDataByFile("testdata/validCountryCodeFile-3")
	if err != nil {
		t.Fatalf("SearchInfo: file validCountryCodeFile-3, but error: %v", err)
	}

	// 渡された文字列をパースしてエラー
	sr := db.SearchInfo("")
	if sr.IsFound {
		t.Errorf("SearchInfo: found: %v", sr)
	}
	if sr.Message != "Invalid IP Address" {
		t.Errorf("SearchInfo: want Message (Invalid IP Address), but got %s: %v", sr.Message, sr)
	}
	if sr.BlockStart != "" {
		t.Errorf("SearchInfo: want BlockStart empty, but got %s: %v", sr.BlockStart, sr)
	}
	if sr.BlockEnd != "" {
		t.Errorf("SearchInfo: want BlockEnd empty, but got %s: %v", sr.BlockEnd, sr)
	}
	if sr.Code != "" {
		t.Errorf("SearchInfo: want Code empty, but got %s: %v", sr.Code, sr)
	}
	if sr.Name != "" {
		t.Errorf("SearchInfo: want EName empty, but got %s: %v", sr.Name, sr)
	}
	if sr.AltName != "" {
		t.Errorf("SearchInfo: want AName empty, but got %s: %v", sr.AltName, sr)
	}

	// IPv4アドレスでない
	sr = db.SearchInfo("2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b")
	if sr.IsFound {
		t.Errorf("SearchInfo: found: %v", sr)
	}
	if sr.Message != "Not IPv4 Address" {
		t.Errorf("SearchInfo: want Message (Not IPv4 Address), but got %s: %v", sr.Message, sr)
	}
	if sr.BlockStart != "" {
		t.Errorf("SearchInfo: want BlockStart empty, but got %s: %v", sr.BlockStart, sr)
	}
	if sr.BlockEnd != "" {
		t.Errorf("SearchInfo: want BlockEnd empty, but got %s: %v", sr.BlockEnd, sr)
	}
	if sr.Code != "" {
		t.Errorf("SearchInfo: want Code empty, but got %s: %v", sr.Code, sr)
	}
	if sr.Name != "" {
		t.Errorf("SearchInfo: want EName empty, but got %s: %v", sr.Name, sr)
	}
	if sr.AltName != "" {
		t.Errorf("SearchInfo: want AName empty, but got %s: %v", sr.AltName, sr)
	}

	// ループバックアドレス
	sr = db.SearchInfo("127.0.0.1")
	if sr.IsFound {
		t.Errorf("SearchInfo: found: %v", sr)
	}
	if sr.Message != "Loopback Address" {
		t.Errorf("SearchInfo: want Message (Loopback Address), but got %s: %v", sr.Message, sr)
	}
	if sr.BlockStart != "" {
		t.Errorf("SearchInfo: want BlockStart empty, but got %s: %v", sr.BlockStart, sr)
	}
	if sr.BlockEnd != "" {
		t.Errorf("SearchInfo: want BlockEnd empty, but got %s: %v", sr.BlockEnd, sr)
	}
	if sr.Code != "" {
		t.Errorf("SearchInfo: want Code empty, but got %s: %v", sr.Code, sr)
	}
	if sr.Name != "" {
		t.Errorf("SearchInfo: want EName empty, but got %s: %v", sr.Name, sr)
	}
	if sr.AltName != "" {
		t.Errorf("SearchInfo: want AName empty, but got %s: %v", sr.AltName, sr)
	}

	// マルチキャストアドレス
	sr = db.SearchInfo("224.0.0.0")
	if sr.IsFound {
		t.Errorf("SearchInfo: found: %v", sr)
	}
	if sr.Message != "Multicast Address" {
		t.Errorf("SearchInfo: want Message (Multicast Address), but got %s: %v", sr.Message, sr)
	}
	if sr.BlockStart != "" {
		t.Errorf("SearchInfo: want BlockStart empty, but got %s: %v", sr.BlockStart, sr)
	}
	if sr.BlockEnd != "" {
		t.Errorf("SearchInfo: want BlockEnd empty, but got %s: %v", sr.BlockEnd, sr)
	}
	if sr.Code != "" {
		t.Errorf("SearchInfo: want Code empty, but got %s: %v", sr.Code, sr)
	}
	if sr.Name != "" {
		t.Errorf("SearchInfo: want EName empty, but got %s: %v", sr.Name, sr)
	}
	if sr.AltName != "" {
		t.Errorf("SearchInfo: want AName empty, but got %s: %v", sr.AltName, sr)
	}

	// プライベートアドレス
	sr = db.SearchInfo("192.168.0.0")
	if sr.IsFound {
		t.Errorf("SearchInfo: found: %v", sr)
	}
	if sr.Message != "Private Address" {
		t.Errorf("SearchInfo: want Message (Private Address), but got %s: %v", sr.Message, sr)
	}
	if sr.BlockStart != "" {
		t.Errorf("SearchInfo: want BlockStart empty, but got %s: %v", sr.BlockStart, sr)
	}
	if sr.BlockEnd != "" {
		t.Errorf("SearchInfo: want BlockEnd empty, but got %s: %v", sr.BlockEnd, sr)
	}
	if sr.Code != "" {
		t.Errorf("SearchInfo: want Code empty, but got %s: %v", sr.Code, sr)
	}
	if sr.Name != "" {
		t.Errorf("SearchInfo: want EName empty, but got %s: %v", sr.Name, sr)
	}
	if sr.AltName != "" {
		t.Errorf("SearchInfo: want AName empty, but got %s: %v", sr.AltName, sr)
	}

	// 情報なし
	sr = db.SearchInfo("114.52.0.0")
	if sr.IsFound {
		t.Errorf("SearchInfo: found: %v", sr)
	}
	if sr.Message != "Not Found" {
		t.Errorf("SearchInfo: want Message (Not Found), but got %s: %v", sr.Message, sr)
	}
	if sr.BlockStart != "" {
		t.Errorf("SearchInfo: want BlockStart empty, but got %s: %v", sr.BlockStart, sr)
	}
	if sr.BlockEnd != "" {
		t.Errorf("SearchInfo: want BlockEnd empty, but got %s: %v", sr.BlockEnd, sr)
	}
	if sr.Code != "" {
		t.Errorf("SearchInfo: want Code empty, but got %s: %v", sr.Code, sr)
	}
	if sr.Name != "" {
		t.Errorf("SearchInfo: want EName empty, but got %s: %v", sr.Name, sr)
	}
	if sr.AltName != "" {
		t.Errorf("SearchInfo: want AName empty, but got %s: %v", sr.AltName, sr)
	}

	// 情報あり
	sr = db.SearchInfo("114.51.255.255")
	if !sr.IsFound {
		t.Errorf("SearchInfo: not found: %v", sr)
	}
	if sr.Message != "Found" {
		t.Errorf("SearchInfo: want Message (Found), but got %s: %v", sr.Message, sr)
	}
	if sr.BlockStart != "114.48.0.0" {
		t.Errorf("SearchInfo: want BlockStart 114.48.0.0, but got %s: %v", sr.BlockStart, sr)
	}
	if sr.BlockEnd != "114.51.255.255" {
		t.Errorf("SearchInfo: want BlockEnd 114.51.255.255, but got %s: %v", sr.BlockEnd, sr)
	}
	if sr.Code != "JP" {
		t.Errorf("SearchInfo: want Code JP, but got %s: %v", sr.Code, sr)
	}
	if sr.Name != "Japan" {
		t.Errorf("SearchInfo: want EName Japan, but got %s: %v", sr.Name, sr)
	}
	if sr.AltName != "日本" {
		t.Errorf("SearchInfo: want AName 日本, but got %s: %v", sr.AltName, sr)
	}
}

func TestLoadIPBDataByURL(t *testing.T) {
	db := GetDB()
	// 不適正な URL
	if err := db.LoadIPBDataByURL(""); err == nil {
		t.Error("LoadIPBDataByURL: url is empty, but no error")
	} else if !strings.Contains(err.Error(), ": empty url") {
		t.Errorf("LoadIPBDataByURL: invalid error: %v:", err)
	}

	// テスト機で http://127.0.0.1 が
	// 動作していない場合のみ行う
	_, err := http.Get("http://127.0.0.1")
	if err != nil {
		// 指定された URL でエラー
		if err := db.LoadIPBDataByURL("http://127.0.0.1"); err == nil {
			t.Errorf("LoadIPBDataByURL: url is invalid (%s) , but no error", "http://127.0.0.1")
		} else if !strings.Contains(err.Error(), "connection refused") {
			t.Errorf("LoadIPBDataByURL: invalid error: %v:", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/bad",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "apnic|JP|ipv4|114.48.0.0|262144|20080422")
		},
	)
	mux.HandleFunc(
		"/good",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "apnic|JP|ipv4|114.48.0.0|262144|20080422|allocated")
		},
	)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// 取得したデータが不正
	if err := db.LoadIPBDataByURL(ts.URL + "/bad"); err == nil {
		t.Errorf("LoadIPBDataByURL: url's data is invalid (%s) , but no error", ts.URL+"/bad")
	} else if !strings.Contains(err.Error(), "the line's fields is invalid") {
		t.Errorf("LoadIPBDataByURL: invalid error: %v:", err)
	}

	// 取得したデータが正常
	if err := db.LoadIPBDataByURL(ts.URL + "/good"); err != nil {
		t.Errorf("LoadIPBDataByURL: url's data is valid (%s) , but error: %v", ts.URL+"/good", err)
	} else {
		if len(db.tmpIB.data) != 1 {
			t.Errorf("LoadIPBDataByURL: tmpIB.data length want 1, but %d: %v", len(db.tmpIB.data), db.tmpIB.data)
		}
		if _, ok := db.tmpIB.data[114][48][0][0]; !ok {
			t.Error("LoadIPBDataByURL: tmpIB.data[114][48][0][0] doesn't exist")
		} else {
			if db.tmpIB.data[114][48][0][0].country != 0 {
				t.Errorf("LoadIPBDataByURL: tmpIB.data[114][48][0][0].country want 0, but got %d", db.tmpIB.data[114][48][0][0].country)
			}
			if db.tmpIB.data[114][48][0][0].value != 262144 {
				t.Errorf("LoadIPBDataByURL: tmpIB.data[114][48][0][0].ipCidr want 262144, but got %d", db.tmpIB.data[114][48][0][0].value)
			}
		}
		if len(db.tmpIB.dicCCStrToInt) != 1 {
			t.Errorf("LoadIPBDataByURL: tmpIB.dicCCStrToInt length want 1, but %d: %v", len(db.tmpIB.dicCCStrToInt), db.tmpIB.dicCCStrToInt)
		} else {
			if _, ok := db.tmpIB.dicCCStrToInt["JP"]; !ok {
				t.Error("LoadIPBDataByURL: tmpIB.dicCCStrToInt[JP] doesn't exist")
			} else if db.tmpIB.dicCCStrToInt["JP"] != 0 {
				t.Errorf("LoadIPBDataByURL: tmpIB.dicCCStrToInt[JP] want 0, but %d", db.tmpIB.dicCCStrToInt["JP"])
			}
		}
		if len(db.tmpIB.dicCCIntToStr) != 1 {
			t.Errorf("LoadIPBDataByURL: tmpIB.dicCCIntToStr length want 1, but %d: %v", len(db.tmpIB.dicCCIntToStr), db.tmpIB.dicCCIntToStr)
		} else {
			if _, ok := db.tmpIB.dicCCIntToStr[0]; !ok {
				t.Error("LoadIPBDataByURL: tmpIB.dicCCIntToStr[0] doesn't exist")
			} else if db.tmpIB.dicCCIntToStr[0] != "JP" {
				t.Errorf("LoadIPBDataByURL: tmpIB.dicCCIntToStr[0] want JP, but %s", db.tmpIB.dicCCIntToStr[0])
			}
		}
		if db.tmpIB.totalBlocks == nil || len(db.tmpIB.totalBlocks) != 2 {
			t.Errorf("LoadIPBDataByURL: tmpIB.totalBlocks is invalid: %v", db.tmpIB.totalBlocks)
		} else {
			if _, ok := db.tmpIB.totalBlocks["ALL"]; !ok {
				t.Error("LoadIPBDataByURL: tmpIB.totalBlocks[ALL] doesn't exist")
			} else if db.tmpIB.totalBlocks["ALL"] != 1 {
				t.Errorf("LoadIPBDataByURL: tmpIB.totalBlocks[ALL] want 1, but %d", db.tmpIB.totalBlocks["ALL"])
			}
			if _, ok := db.tmpIB.totalBlocks["JP"]; !ok {
				t.Error("LoadIPBDataByURL: tmpIB.totalBlocks[JP] doesn't exist")
			} else if db.tmpIB.totalBlocks["JP"] != 1 {
				t.Errorf("LoadIPBDataByURL: tmpIB.totalBlocks[JP] want 1, but %d", db.tmpIB.totalBlocks["JP"])
			}
		}
		if db.tmpIB.totalValue == nil || len(db.tmpIB.totalValue) != 2 {
			t.Errorf("LoadIPBDataByURL: tmpIB.totalValue is invalid: %v", db.tmpIB.totalValue)
		} else {
			if _, ok := db.tmpIB.totalValue["ALL"]; !ok {
				t.Error("LoadIPBDataByURL: tmpIB.totalValue[ALL] doesn't exist")
			} else if db.tmpIB.totalValue["ALL"] != 262144 {
				t.Errorf("LoadIPBDataByURL: tmpIB.totalValue[ALL] want 262144, but %d", db.tmpIB.totalValue["ALL"])
			}
			if _, ok := db.tmpIB.totalValue["JP"]; !ok {
				t.Error("LoadIPBDataByURL: tmpIB.totalValue[JP] doesn't exist")
			} else if db.tmpIB.totalValue["JP"] != 262144 {
				t.Errorf("LoadIPBDataByURL: tmpIB.totalValue[JP] want 262144, but %d", db.tmpIB.totalValue["JP"])
			}
		}

	}
}

func TestSetIPBData(t *testing.T) {
	db := &DB{
		tmpIB: ipBlocks{
			data:          map[uint8]map[uint8]map[uint8]map[uint8]block{},
			dicCCIntToStr: map[uint8]string{},
			dicCCStrToInt: map[string]uint8{},
			totalBlocks:   map[string]int{"ALL": 0},
			totalValue:    map[string]int{"ALL": 0},
		},
		reg: regexp.MustCompile(`^[A-Z]{2}$`),
		urlRIR: []string{
			"http://127.0.0.1",
		},
	}
	// テスト機で http://127.0.0.1 が
	// 動作していない場合のみ行う
	_, err := http.Get("http://127.0.0.1")
	if err != nil {
		// 指定された URL でエラー
		if err := db.SetIPBData(); err == nil {
			t.Errorf("SetIPBData: urlRIR is invalid (%s) , but no error", "http://127.0.0.1")
		} else if !strings.Contains(err.Error(), "connection refused") {
			t.Errorf("SetIPBData: invalid error: %v:", err)
		}
	}

	urlRIR, mux := getDummyRIR()

	ts := httptest.NewServer(mux)
	defer ts.Close()

	db.urlRIR = make([]string, len(urlRIR))
	for i := range urlRIR {
		db.urlRIR[i] = ts.URL + urlRIR[i]
	}

	if err := db.SetIPBData(); err != nil {
		t.Errorf("SetIPBData: urlRIR is valid , but error: %v", err)
	} else {
		if len(db.tmpIB.data) != 0 {
			t.Errorf("SetIPBData: tmpIB.data length want 0, but %d: %v", len(db.tmpIB.data), db.tmpIB.data)
		}
		if len(db.tmpIB.dicCCStrToInt) != 0 {
			t.Errorf("SetIPBData: tmpIB.dicCCStrToInt length want 0, but %d: %v", len(db.tmpIB.dicCCStrToInt), db.tmpIB.dicCCStrToInt)
		}
		if len(db.tmpIB.dicCCIntToStr) != 0 {
			t.Errorf("SetIPBData: tmpIB.dicCCIntToStr length want 0, but %d: %v", len(db.tmpIB.dicCCIntToStr), db.tmpIB.dicCCIntToStr)
		}
		if len(db.tmpIB.totalBlocks) != 1 {
			t.Errorf("SetIPBData: tmpIB.totalBlocks length want 1, but %d: %v", len(db.tmpIB.totalBlocks), db.tmpIB.totalBlocks)
		}
		if len(db.tmpIB.totalValue) != 1 {
			t.Errorf("SetIPBData: tmpIB.totalValue length want 1, but %d: %v", len(db.tmpIB.totalValue), db.tmpIB.totalValue)
		}

		if len(db.ib.data) != 4 {
			t.Errorf("SetIPBData: ib.data length want 4, but %d: %v", len(db.ib.data), db.ib.data)
		}
		if _, ok := db.ib.data[41][0][0][0]; !ok {
			t.Error("SetIPBData: ib.data[41][0][0][0] doesn't exist")
		} else {
			if db.ib.data[41][0][0][0].country != 0 {
				t.Errorf("SetIPBData: ib.data[41][0][0][0].country want 0, but got %d", db.ib.data[41][0][0][0].country)
			}
			if db.ib.data[41][0][0][0].value != 2097152 {
				t.Errorf("SetIPBData: ib.data[41][0][0][0].ipCidr want 2097152, but got %d", db.ib.data[41][0][0][0].value)
			}
		}
		if _, ok := db.ib.data[1][0][0][0]; !ok {
			t.Error("SetIPBData: ib.data[1][0][0][0] doesn't exist")
		} else {
			if db.ib.data[1][0][0][0].country != 1 {
				t.Errorf("SetIPBData: ib.data[1][0][0][0].country want 1, but got %d", db.ib.data[1][0][0][0].country)
			}
			if db.ib.data[1][0][0][0].value != 256 {
				t.Errorf("SetIPBData: ib.data[1][0][0][0].ipCidr want 256, but got %d", db.ib.data[1][0][0][0].value)
			}
		}
		if _, ok := db.ib.data[2][57][164][0]; !ok {
			t.Error("SetIPBData: ib.data[2][57][164][0] doesn't exist")
		} else {
			if db.ib.data[2][57][164][0].country != 2 {
				t.Errorf("SetIPBData: ib.data[2][57][164][0].country want 2, but got %d", db.ib.data[2][57][164][0].country)
			}
			if db.ib.data[2][57][164][0].value != 1024 {
				t.Errorf("SetIPBData: ib.data[2][57][164][0].ipCidr want 1024, but got %d", db.ib.data[2][57][164][0].value)
			}
		}
		if _, ok := db.ib.data[5][183][80][0]; !ok {
			t.Error("SetIPBData: ib.data[5][183][80][0] doesn't exist")
		} else {
			if db.ib.data[5][183][80][0].country != 3 {
				t.Errorf("SetIPBData: ib.data[5][183][80][0].country want 3, but got %d", db.ib.data[5][183][80][0].country)
			}
			if db.ib.data[5][183][80][0].value != 1024 {
				t.Errorf("SetIPBData: ib.data[5][183][80][0].ipCidr want 1024, but got %d", db.ib.data[5][183][80][0].value)
			}
		}
		if _, ok := db.ib.data[1][178][112][0]; !ok {
			t.Error("SetIPBData: ib.data[1][178][112][0] doesn't exist")
		} else {
			if db.ib.data[1][178][112][0].country != 4 {
				t.Errorf("SetIPBData: ib.data[1][178][112][0].country want 4, but got %d", db.ib.data[1][178][112][0].country)
			}
			if db.ib.data[1][178][112][0].value != 4096 {
				t.Errorf("SetIPBData: ib.data[1][178][112][0].ipCidr want 4096, but got %d", db.ib.data[1][178][112][0].value)
			}
		}
		if len(db.ib.dicCCStrToInt) != 5 {
			t.Errorf("SetIPBData: ib.dicCCStrToInt length want 5, but %d: %v", len(db.ib.dicCCStrToInt), db.ib.dicCCStrToInt)
		} else {
			if _, ok := db.ib.dicCCStrToInt["ZA"]; !ok {
				t.Error("SetIPBData: ib.dicCCStrToInt[ZA] doesn't exist")
			} else if db.ib.dicCCStrToInt["ZA"] != 0 {
				t.Errorf("SetIPBData: ib.dicCCStrToInt[ZA] want 0, but %d", db.ib.dicCCStrToInt["ZA"])
			}
			if _, ok := db.ib.dicCCStrToInt["AU"]; !ok {
				t.Error("SetIPBData: ib.dicCCStrToInt[AU] doesn't exist")
			} else if db.ib.dicCCStrToInt["AU"] != 1 {
				t.Errorf("SetIPBData: ib.dicCCStrToInt[AU] want 1, but %d", db.ib.dicCCStrToInt["AU"])
			}
			if _, ok := db.ib.dicCCStrToInt["US"]; !ok {
				t.Error("SetIPBData: ib.dicCCStrToInt[US] doesn't exist")
			} else if db.ib.dicCCStrToInt["US"] != 2 {
				t.Errorf("SetIPBData: ib.dicCCStrToInt[US] want 2, but %d", db.ib.dicCCStrToInt["US"])
			}
			if _, ok := db.ib.dicCCStrToInt["DO"]; !ok {
				t.Error("SetIPBData: ib.dicCCStrToInt[DO] doesn't exist")
			} else if db.ib.dicCCStrToInt["DO"] != 3 {
				t.Errorf("SetIPBData: ib.dicCCStrToInt[DO] want 3, but %d", db.ib.dicCCStrToInt["DO"])
			}
			if _, ok := db.ib.dicCCStrToInt["PS"]; !ok {
				t.Error("SetIPBData: ib.dicCCStrToInt[PS] doesn't exist")
			} else if db.ib.dicCCStrToInt["PS"] != 4 {
				t.Errorf("SetIPBData: ib.dicCCStrToInt[PS] want 4, but %d", db.ib.dicCCStrToInt["PS"])
			}
		}
		if len(db.ib.dicCCIntToStr) != 5 {
			t.Errorf("SetIPBData: ib.dicCCIntToStr length want 5, but %d: %v", len(db.ib.dicCCIntToStr), db.ib.dicCCIntToStr)
		} else {
			if _, ok := db.ib.dicCCIntToStr[0]; !ok {
				t.Error("SetIPBData: ib.dicCCIntToStr[0] doesn't exist")
			} else if db.ib.dicCCIntToStr[0] != "ZA" {
				t.Errorf("SetIPBData: ib.dicCCIntToStr[0] want ZA, but %s", db.ib.dicCCIntToStr[0])
			}
			if _, ok := db.ib.dicCCIntToStr[1]; !ok {
				t.Error("SetIPBData: ib.dicCCIntToStr[1] doesn't exist")
			} else if db.ib.dicCCIntToStr[1] != "AU" {
				t.Errorf("SetIPBData: ib.dicCCIntToStr[1] want AU, but %s", db.ib.dicCCIntToStr[1])
			}
			if _, ok := db.ib.dicCCIntToStr[2]; !ok {
				t.Error("SetIPBData: ib.dicCCIntToStr[2] doesn't exist")
			} else if db.ib.dicCCIntToStr[2] != "US" {
				t.Errorf("SetIPBData: ib.dicCCIntToStr[2] want US, but %s", db.ib.dicCCIntToStr[2])
			}
			if _, ok := db.ib.dicCCIntToStr[3]; !ok {
				t.Error("SetIPBData: ib.dicCCIntToStr[3] doesn't exist")
			} else if db.ib.dicCCIntToStr[3] != "DO" {
				t.Errorf("SetIPBData: ib.dicCCIntToStr[3] want DO, but %s", db.ib.dicCCIntToStr[3])
			}
			if _, ok := db.ib.dicCCIntToStr[4]; !ok {
				t.Error("SetIPBData: ib.dicCCIntToStr[4] doesn't exist")
			} else if db.ib.dicCCIntToStr[4] != "PS" {
				t.Errorf("SetIPBData: ib.dicCCIntToStr[4] want PS, but %s", db.ib.dicCCIntToStr[4])
			}
		}
		if len(db.ib.totalBlocks) != 6 {
			t.Errorf("SetIPBData: ib.totalBlocks length want 6, but %d: %v", len(db.ib.totalBlocks), db.ib.totalBlocks)
		} else {
			if _, ok := db.ib.totalBlocks["ALL"]; !ok {
				t.Error("SetIPBData: ib.totalBlocks[ALL] doesn't exist")
			} else if db.ib.totalBlocks["ALL"] != 5 {
				t.Errorf("SetIPBData: ib.totalBlocks[ALL] want 5, but %d", db.ib.totalBlocks["ALL"])
			}
			if _, ok := db.ib.totalBlocks["ZA"]; !ok {
				t.Error("SetIPBData: ib.totalBlocks[ZA] doesn't exist")
			} else if db.ib.totalBlocks["ZA"] != 1 {
				t.Errorf("SetIPBData: ib.totalBlocks[ZA] want 1, but %d", db.ib.totalBlocks["ZA"])
			}
			if _, ok := db.ib.totalBlocks["AU"]; !ok {
				t.Error("SetIPBData: ib.totalBlocks[AU] doesn't exist")
			} else if db.ib.totalBlocks["AU"] != 1 {
				t.Errorf("SetIPBData: ib.totalBlocks[AU] want 1, but %d", db.ib.totalBlocks["AU"])
			}
			if _, ok := db.ib.totalBlocks["US"]; !ok {
				t.Error("SetIPBData: ib.totalBlocks[US] doesn't exist")
			} else if db.ib.totalBlocks["US"] != 1 {
				t.Errorf("SetIPBData: ib.totalBlocks[US] want 1, but %d", db.ib.totalBlocks["US"])
			}
			if _, ok := db.ib.totalBlocks["DO"]; !ok {
				t.Error("SetIPBData: ib.totalBlocks[DO] doesn't exist")
			} else if db.ib.totalBlocks["DO"] != 1 {
				t.Errorf("SetIPBData: ib.totalBlocks[DO] want 1, but %d", db.ib.totalBlocks["DO"])
			}
			if _, ok := db.ib.totalBlocks["PS"]; !ok {
				t.Error("SetIPBData: ib.totalBlocks[PS] doesn't exist")
			} else if db.ib.totalBlocks["PS"] != 1 {
				t.Errorf("SetIPBData: ib.totalBlocks[PS] want 1, but %d", db.ib.totalBlocks["PS"])
			}
		}
		if len(db.ib.totalValue) != 6 {
			t.Errorf("SetIPBData: ib.totalValue length want 6, but %d: %v", len(db.ib.totalValue), db.ib.totalValue)
		} else {
			if _, ok := db.ib.totalValue["ALL"]; !ok {
				t.Error("SetIPBData: ib.totalValue[ALL] doesn't exist")
			} else if db.ib.totalValue["ALL"] != 2103552 {
				t.Errorf("SetIPBData: ib.totalValue[ALL] want 2103552, but %d", db.ib.totalValue["ALL"])
			}
			if _, ok := db.ib.totalValue["ZA"]; !ok {
				t.Error("SetIPBData: ib.totalValue[ZA] doesn't exist")
			} else if db.ib.totalValue["ZA"] != 2097152 {
				t.Errorf("SetIPBData: ib.totalValue[ZA] want 2097152, but %d", db.ib.totalValue["ZA"])
			}
			if _, ok := db.ib.totalValue["AU"]; !ok {
				t.Error("SetIPBData: ib.totalValue[AU] doesn't exist")
			} else if db.ib.totalValue["AU"] != 256 {
				t.Errorf("SetIPBData: ib.totalValue[AU] want 256, but %d", db.ib.totalValue["AU"])
			}
			if _, ok := db.ib.totalValue["US"]; !ok {
				t.Error("SetIPBData: ib.totalValue[US] doesn't exist")
			} else if db.ib.totalValue["US"] != 1024 {
				t.Errorf("SetIPBData: ib.totalValue[US] want 1024, but %d", db.ib.totalValue["US"])
			}
			if _, ok := db.ib.totalValue["DO"]; !ok {
				t.Error("SetIPBData: ib.totalValue[DO] doesn't exist")
			} else if db.ib.totalValue["DO"] != 1024 {
				t.Errorf("SetIPBData: ib.totalValue[DO] want 1024, but %d", db.ib.totalValue["DO"])
			}
			if _, ok := db.ib.totalValue["PS"]; !ok {
				t.Error("SetIPBData: ib.totalValue[PS] doesn't exist")
			} else if db.ib.totalValue["PS"] != 4096 {
				t.Errorf("SetIPBData: ib.totalValue[PS] want 4096, but %d", db.ib.totalValue["PS"])
			}
		}
	}
}

func TestGetOneOutside(t *testing.T) {
	// value が 1 : 引数のアドレスの次のアドレスが返る
	addr := getOneOutside([4]byte{0, 0, 0, 0}, 1)
	if addr.String() != "0.0.0.1" {
		t.Errorf("getOneOutside: want 0.0.0.1, but got %s:", addr.String())
	}

	// value が 256 : 引数のアドレスに対して、3 番めの 8 ビットが一つ増える
	addr = getOneOutside([4]byte{0, 0, 0, 0}, 256)
	if addr.String() != "0.0.1.0" {
		t.Errorf("getOneOutside: want 0.0.1.0, but got %s:", addr.String())
	}

	// value が 65536 : 引数のアドレスに対して、2 番めの 8 ビットが一つ増える
	addr = getOneOutside([4]byte{0, 0, 0, 0}, 65536)
	if addr.String() != "0.1.0.0" {
		t.Errorf("getOneOutside: want 0.1.0.0, but got %s:", addr.String())
	}

	// value が 16777216 : 引数のアドレスに対して、最初の 8 ビットが一つ増える
	addr = getOneOutside([4]byte{0, 0, 0, 0}, 16777216)
	if addr.String() != "1.0.0.0" {
		t.Errorf("getOneOutside: want 1.0.0.0, but got %s:", addr.String())
	}

	// value が 4294967295 : 引数のアドレスの前のアドレスが返る
	addr = getOneOutside([4]byte{1, 0, 0, 0}, 4294967295)
	if addr.String() != "0.255.255.255" {
		t.Errorf("getOneOutside: want 0.255.255.255, but got %s:", addr.String())
	}
}

func TestGetLastAddr(t *testing.T) {
	// IPアドレスでない
	addr, err := GetLastAddr("", 1)
	if err == nil {
		t.Errorf("GetLastAddr: addr is invalid, but no error: %v:", addr)
	} else if !strings.Contains(err.Error(), "unable to parse IP") {
		t.Errorf("GetLastAddr: invalid error: %v:", err)
	}

	// IPv4アドレスでない
	addr, err = GetLastAddr("2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b", 1)
	if err == nil {
		t.Errorf("GetLastAddr: addr is invalid, but no error: %v:", addr)
	} else if err != ErrFirstArgumentOutOfRange {
		t.Errorf("GetLastAddr: invalid error: %v:", err)
	}

	// value が範囲外・１
	addr, err = GetLastAddr("0.0.0.0", 0)
	if err == nil {
		t.Errorf("GetLastAddr: value is invalid, but no error: %v:", addr)
	} else if err != ErrSecondArgumentOutOfRange {
		t.Errorf("GetLastAddr: invalid error: %v:", err)
	}

	// value が範囲外・２
	addr, err = GetLastAddr("0.0.0.0", 4294967296)
	if err == nil {
		t.Errorf("GetLastAddr: value is invalid, but no error: %v:", addr)
	} else if err != ErrSecondArgumentOutOfRange {
		t.Errorf("GetLastAddr: invalid error: %v:", err)
	}

	// 最後の IP アドレス : 0.0.0.0
	addr, err = GetLastAddr("0.0.0.0", 1)
	if err != nil {
		t.Errorf("GetLastAddr: addr and value are valid, but error: %v:", err)
	} else if addr.String() != "0.0.0.0" {
		t.Errorf("GetLastAddr: want last address 0.0.0.0, but got %s:", addr.String())
	}

	// 最後の IP アドレス : 0.0.0.255
	addr, err = GetLastAddr("0.0.0.0", 256)
	if err != nil {
		t.Errorf("GetLastAddr: addr and value are valid, but error: %v:", err)
	} else if addr.String() != "0.0.0.255" {
		t.Errorf("GetLastAddr: want last address 0.0.0.255, but got %s:", addr.String())
	}

	// 最後の IP アドレス : 0.0.255.255
	addr, err = GetLastAddr("0.0.0.0", 65536)
	if err != nil {
		t.Errorf("GetLastAddr: addr and value are valid, but error: %v:", err)
	} else if addr.String() != "0.0.255.255" {
		t.Errorf("GetLastAddr: want last address 0.0.255.255, but got %s:", addr.String())
	}

	// 最後の IP アドレス : 0.255.255.255
	addr, err = GetLastAddr("0.0.0.0", 16777216)
	if err != nil {
		t.Errorf("GetLastAddr: addr and value are valid, but error: %v:", err)
	} else if addr.String() != "0.255.255.255" {
		t.Errorf("GetLastAddr: want last address 0.255.255.255, but got %s:", addr.String())
	}

	// 最後の IP アドレス : 126.255.255.255
	addr, err = GetLastAddr("127.0.0.1", 4294967295)
	if err != nil {
		t.Errorf("GetLastAddr: addr and value are valid, but error: %v:", err)
	} else if addr.String() != "126.255.255.255" {
		t.Errorf("GetLastAddr: want last address 126.255.255.255, but got %s:", addr.String())
	}
}

func TestGetValue(t *testing.T) {
	// IPアドレスでない（先）
	value, err := GetValue("", "0.0.0.0")
	if err == nil {
		t.Errorf("GetValue: first address is invalid, but no error: %d:", value)
	} else if !strings.Contains(err.Error(), "unable to parse IP") {
		t.Errorf("GetValue: invalid error: %v:", err)
	}

	// IPv4アドレスでない（先）
	value, err = GetValue("2001:0db8:3c4d:0015:0000:0000:1a2f:1a2c", "0.0.0.0")
	if err == nil {
		t.Errorf("GetValue: first address is invalid, but no error: %d:", value)
	} else if err != ErrFirstArgumentOutOfRange {
		t.Errorf("GetValue: invalid error: %v:", err)
	}

	// IPアドレスでない（後）
	value, err = GetValue("0.0.0.0", "")
	if err == nil {
		t.Errorf("GetValue: second address is invalid, but no error: %d:", value)
	} else if !strings.Contains(err.Error(), "unable to parse IP") {
		t.Errorf("GetValue: invalid error: %v:", err)
	}

	// IPv4アドレスでない（後）
	value, err = GetValue("0.0.0.0", "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2d")
	if err == nil {
		t.Errorf("GetValue: second address is invalid, but no error: %d:", value)
	} else if err != ErrSecondArgumentOutOfRange {
		t.Errorf("GetValue: invalid error: %v:", err)
	}

	// 同じアドレス
	value, err = GetValue("0.0.0.0", "0.0.0.0")
	if err != nil {
		t.Errorf("GetValue: both addresses are valid, but error: % :", err)
	} else if value != 0 {
		t.Errorf("GetValue: want value 0, error: %d:", value)
	}

	// 先が小・後が大
	value, err = GetValue("0.0.0.0", "0.0.0.255")
	if err != nil {
		t.Errorf("GetValue: both addresses are valid, but error: % :", err)
	} else if value != 256 {
		t.Errorf("GetValue: want value 256, error: %d:", value)
	}

	//  先が大・後が小
	value, err = GetValue("0.0.255.0", "0.0.0.255")
	if err != nil {
		t.Errorf("GetValue: both addresses are valid, but error: % :", err)
	} else if value != 65026 {
		t.Errorf("GetValue: want value 65026, error: %d:", value)
	}
}

func TestGetTotalBlocks(t *testing.T) {
	db := &DB{
		tmpIB: ipBlocks{
			data:          map[uint8]map[uint8]map[uint8]map[uint8]block{},
			dicCCIntToStr: map[uint8]string{},
			dicCCStrToInt: map[string]uint8{},
			totalBlocks:   map[string]int{"ALL": 0},
			totalValue:    map[string]int{"ALL": 0},
		},
		reg: regexp.MustCompile(`^[A-Z]{2}$`),
	}

	// データが空
	tb := db.GetTotalBlocks()
	if len(tb) != 0 {
		t.Errorf("GetTotalBlocks: want 0, error: %d: %v", len(tb), tb)
	}

	urlRIR, mux := getDummyRIR()

	ts := httptest.NewServer(mux)
	defer ts.Close()

	db.urlRIR = make([]string, len(urlRIR))
	for i := range urlRIR {
		db.urlRIR[i] = ts.URL + urlRIR[i]
	}

	if err := db.SetIPBData(); err != nil {
		t.Fatalf("GetTotalBlocks: failed to set IPB data: %v", err)
	}

	// 各 RIR データが 1 づつ
	tb = db.GetTotalBlocks()
	if len(tb) != 6 {
		t.Errorf("GetTotalBlocks: want 6, error: %d: %v", len(tb), tb)
	} else {
		if _, ok := tb["ALL"]; !ok {
			t.Error("GetTotalBlocks: tb[ALL] doesn't exist")
		} else if tb["ALL"] != 5 {
			t.Errorf("GetTotalBlocks: tb[ALL] want 1, but %d", tb["ALL"])
		}
		if _, ok := tb["ZA"]; !ok {
			t.Error("GetTotalBlocks: tb[ZA] doesn't exist")
		} else if tb["ZA"] != 1 {
			t.Errorf("GetTotalBlocks: tb[ZA] want 1, but %d", tb["ZA"])
		}
		if _, ok := tb["AU"]; !ok {
			t.Error("GetTotalBlocks: tb[AU] doesn't exist")
		} else if tb["AU"] != 1 {
			t.Errorf("GetTotalBlocks: tb[AU] want 1, but %d", tb["AU"])
		}
		if _, ok := tb["US"]; !ok {
			t.Error("GetTotalBlocks: tb[US] doesn't exist")
		} else if tb["US"] != 1 {
			t.Errorf("GetTotalBlocks: tb[US] want 1, but %d", tb["US"])
		}
		if _, ok := tb["DO"]; !ok {
			t.Error("GetTotalBlocks: tb[DO] doesn't exist")
		} else if tb["DO"] != 1 {
			t.Errorf("GetTotalBlocks: tb[DO] want 1, but %d", tb["DO"])
		}
		if _, ok := tb["PS"]; !ok {
			t.Error("GetTotalBlocks: tb[PS] doesn't exist")
		} else if tb["PS"] != 1 {
			t.Errorf("GetTotalBlocks: tb[PS] want 1, but %d", tb["PS"])
		}
	}
}

func TestGetTotalValue(t *testing.T) {
	db := &DB{
		tmpIB: ipBlocks{
			data:          map[uint8]map[uint8]map[uint8]map[uint8]block{},
			dicCCIntToStr: map[uint8]string{},
			dicCCStrToInt: map[string]uint8{},
			totalBlocks:   map[string]int{"ALL": 0},
			totalValue:    map[string]int{"ALL": 0},
		},
		reg: regexp.MustCompile(`^[A-Z]{2}$`),
	}

	// データが空
	tv := db.GetTotalValue()
	if len(tv) != 0 {
		t.Errorf("GetTotalValue: want 0, error: %d: %v", len(tv), tv)
	}

	urlRIR, mux := getDummyRIR()

	ts := httptest.NewServer(mux)
	defer ts.Close()

	db.urlRIR = make([]string, len(urlRIR))
	for i := range urlRIR {
		db.urlRIR[i] = ts.URL + urlRIR[i]
	}

	if err := db.SetIPBData(); err != nil {
		t.Fatalf("GetTotalValue: failed to set IPB data: %v", err)
	}

	// 各 RIR データが 1 づつ
	tv = db.GetTotalValue()
	if len(tv) != 6 {
		t.Errorf("GetTotalValue: want 6, error: %d: %v", len(tv), tv)
	} else {
		if _, ok := tv["ALL"]; !ok {
			t.Error("GetTotalValue: tv[ALL] doesn't exist")
		} else if tv["ALL"] != 2103552 {
			t.Errorf("GetTotalValue: tv[ALL] want 2103552, but %d", tv["ALL"])
		}
		if _, ok := tv["ZA"]; !ok {
			t.Error("GetTotalValue: tv[ZA] doesn't exist")
		} else if tv["ZA"] != 2097152 {
			t.Errorf("GetTotalValue: tv[ZA] want 2097152, but %d", tv["ZA"])
		}
		if _, ok := tv["AU"]; !ok {
			t.Error("GetTotalValue: tv[AU] doesn't exist")
		} else if tv["AU"] != 256 {
			t.Errorf("GetTotalValue: tv[AU] want 256, but %d", tv["AU"])
		}
		if _, ok := tv["US"]; !ok {
			t.Error("GetTotalValue: tv[US] doesn't exist")
		} else if tv["US"] != 1024 {
			t.Errorf("GetTotalValue: tv[US] want 1024, but %d", tv["US"])
		}
		if _, ok := tv["DO"]; !ok {
			t.Error("GetTotalValue: tv[DO] doesn't exist")
		} else if tv["DO"] != 1024 {
			t.Errorf("GetTotalValue: tv[DO] want 1024, but %d", tv["DO"])
		}
		if _, ok := tv["PS"]; !ok {
			t.Error("GetTotalValue: tv[PS] doesn't exist")
		} else if tv["PS"] != 4096 {
			t.Errorf("GetTotalValue: tv[PS] want 4096, but %d", tv["PS"])
		}
	}
}

func TestGetCountryCodeData(t *testing.T) {
	db := &DB{
		tmpIB: ipBlocks{
			data:          map[uint8]map[uint8]map[uint8]map[uint8]block{},
			dicCCIntToStr: map[uint8]string{},
			dicCCStrToInt: map[string]uint8{},
			totalBlocks:   map[string]int{"ALL": 0},
			totalValue:    map[string]int{"ALL": 0},
		},
		reg: regexp.MustCompile(`^[A-Z]{2}$`),
	}

	// データが空
	cc := db.GetCountryCodeData()
	if len(cc) != 0 {
		t.Errorf("GetCountryCodeData: want 0, error: %d: %v", len(cc), cc)
	}

	if err := db.InitCCDataByFile("testdata/validCountryCodeFile-2"); err != nil {
		t.Fatalf("GetCountryCodeData: failed to set CC data: %v", err)
	}

	// データ２件
	cc = db.GetCountryCodeData()
	if len(cc) != 2 {
		t.Errorf("GetCountryCodeData: CountryCodeData length want 2, but %d: %v", len(cc), cc)
	} else {
		if _, ok := cc["AD"]; !ok {
			t.Error("GetCountryCodeData: CountryCodeData[AD] doesn't exist")
		} else {
			if cc["AD"].Name != "Principality of Andorra" {
				t.Errorf("GetCountryCodeData: CountryCodeData[AD].eName want Principality of Andorra, but %s", cc["AD"].Name)
			}
			if cc["AD"].AltName != "アンドラ公国" {
				t.Errorf("GetCountryCodeData: CountryCodeData[AD].aName want アンドラ公国, but %s", cc["AD"].AltName)
			}
		}
		if _, ok := cc["AE"]; !ok {
			t.Error("GetCountryCodeData: CountryCodeData[AE] doesn't exist")
		} else {
			if cc["AE"].Name != "United Arab Emirates" {
				t.Errorf("GetCountryCodeData: CountryCodeData[AE].eName want United Arab Emirates, but %s", cc["AE"].Name)
			}
			if cc["AE"].AltName != "アラブ首長国連邦" {
				t.Errorf("GetCountryCodeData: CountryCodeData[AE].aName want アラブ首長国連邦, but %s", cc["AE"].AltName)
			}
		}
	}

}

func TestIsDBEmpty(t *testing.T) {
	db := &DB{}

	// データベースが nil
	if !db.IsDBEmpty() {
		t.Errorf("IsDBEmpty: empty, but false: %v", db.ib.data)
	}

	// データベースが空
	db.ib.data = map[uint8]map[uint8]map[uint8]map[uint8]block{}
	if !db.IsDBEmpty() {
		t.Errorf("IsDBEmpty: empty, but false: %v", db.ib.data)
	}

	// データあり
	db.ib.data[0] = map[uint8]map[uint8]map[uint8]block{}
	db.ib.data[0][0] = map[uint8]map[uint8]block{}
	db.ib.data[0][0][0] = map[uint8]block{}
	db.ib.data[0][0][0][0] = block{value: 16, country: 0}
	if db.IsDBEmpty() {
		t.Errorf("IsDBEmpty: not empty, but true: %v", db.ib.data)
	}
}
