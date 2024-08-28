package ccipv4

import (
	"encoding/binary"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const (
	// 各 RIR の最新版 delegation file の URL
	URL_DELEGATED_AFRINIC_EXTENDED_LATEST string = "https://ftp.afrinic.net/pub/stats/afrinic/delegated-afrinic-extended-latest"
	URL_DELEGATED_APNIC_EXTENDED_LATEST   string = "https://ftp.apnic.net/stats/apnic/delegated-apnic-extended-latest"
	URL_DELEGATED_ARIN_EXTENDED_LATEST    string = "https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest"
	URL_DELEGATED_LACNIC_EXTENDED_LATEST  string = "https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-extended-latest"
	URL_DELEGATED_RIPENCC_EXTENDED_LATEST string = "https://ftp.ripe.net/pub/stats/ripencc/delegated-ripencc-extended-latest"
)

type block struct {
	value   uint32
	country uint8
}

type CountryCodeInfo struct {
	Name    string
	AltName string
}

type SearchResult struct {
	IsFound    bool
	Message    string
	BlockStart string
	BlockEnd   string
	Code       string
	Name       string
	AltName    string
}

type ipBlocks struct {
	l             sync.RWMutex
	data          map[uint8]map[uint8]map[uint8]map[uint8]block
	dicCCStrToInt map[string]uint8
	dicCCIntToStr map[uint8]string
	totalBlocks   map[string]int
	totalValue    map[string]int
}

type countryCodes struct {
	l    sync.RWMutex
	data map[string]CountryCodeInfo
}

type DB struct {
	ib     ipBlocks
	tmpIB  ipBlocks
	cc     countryCodes
	tmpCC  countryCodes
	reg    *regexp.Regexp
	urlRIR []string
}

var ErrFirstArgumentOutOfRange = errors.New("first argument out of range")

var ErrSecondArgumentOutOfRange = errors.New("second argument out of range")

// 初期状態のデータベースを取得する。
func GetDB() *DB {
	return &DB{
		tmpIB: ipBlocks{
			data:          map[uint8]map[uint8]map[uint8]map[uint8]block{},
			dicCCIntToStr: map[uint8]string{},
			dicCCStrToInt: map[string]uint8{},
			totalBlocks:   map[string]int{"ALL": 0},
			totalValue:    map[string]int{"ALL": 0},
		},
		tmpCC: countryCodes{
			data: map[string]CountryCodeInfo{},
		},
		reg: regexp.MustCompile(`^[A-Z]{2}$`),
		urlRIR: []string{
			URL_DELEGATED_AFRINIC_EXTENDED_LATEST,
			URL_DELEGATED_APNIC_EXTENDED_LATEST,
			URL_DELEGATED_ARIN_EXTENDED_LATEST,
			URL_DELEGATED_LACNIC_EXTENDED_LATEST,
			URL_DELEGATED_RIPENCC_EXTENDED_LATEST,
		},
	}

}

// 一時保存用データベースを空データベースにする。
func (db *DB) ClearTmpIPBData() {
	db.tmpIB.data = map[uint8]map[uint8]map[uint8]map[uint8]block{}
	db.tmpIB.dicCCIntToStr = map[uint8]string{}
	db.tmpIB.dicCCStrToInt = map[string]uint8{}
	db.tmpIB.totalBlocks = map[string]int{"ALL": 0}
	db.tmpIB.totalValue = map[string]int{"ALL": 0}
}

// io.Reader を使って RIR statistics exchange format を読み込む。
// RIR statistics exchange format については下記を参照。
// http://www.apnic.net/db/rir-stats-format.html
func (db *DB) setTmpIPBlocks(r io.Reader) error {
	// データベースロック
	db.tmpIB.l.Lock()
	defer db.tmpIB.l.Unlock()

	// カントリーコード文字列にひも付けする uint8 用の値（後述）を
	// 設定するために使用。
	countForDic := uint8(len(db.tmpIB.dicCCStrToInt))

	// ファイルを csv として読込。
	// format に従い、コメント・フィールド区切りの文字を設定。
	reader := csv.NewReader(r)
	reader.Comment = '#'
	reader.Comma = '|'

	// ファイルを一行単位で読込。
	// 異常が発生した場合はその行で処理を中止し、
	// 一時保存用データベースを空にする。
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				// Field が６つで最後に文字列 "summary" が設定されている場合は
				// record ではなく header の summary line で処理対象外だが異常
				// ではないので invalidDataSet に追記しない。
				if len(line) == 6 && line[5] == "summary" {
					continue
				}
				if len(line) < 6 || !strings.Contains(err.Error(), "wrong number of fields") {
					db.ClearTmpIPBData()
					return fmt.Errorf("%v: %v", err, line)
				}
			}
		}
		// 先頭の Field が registry ではなく version の場合は
		// record ではなく header の version line で処理対象外。
		if !slices.Contains([]string{"afrinic", "apnic", "arin", "iana", "lacnic", "ripencc"}, line[0]) {
			continue
		}
		// Record format の３番めの Field は type 。
		// asn,ipv4,ipv6 のいずれか。
		if line[2] == "ipv4" {
			// Record format の Field の個数は7以上。
			if len(line) < 7 {
				db.ClearTmpIPBData()
				return fmt.Errorf("the number(%d) of the line's fields is invalid: %v", len(line), line)
			}
			// Record format の４番めの Field は start 。
			// 対象範囲の最初のアドレスを示す。
			ad, err := netip.ParseAddr(line[3])
			if err != nil {
				db.ClearTmpIPBData()
				return fmt.Errorf("%v: %v", err, line)
			}
			// 検索に使用するため、start のアドレスを８ビットで分割し、
			// ipBlocks のマップのキーとする。
			group8Bits := ad.AsSlice()
			if _, ok := db.tmpIB.data[group8Bits[0]]; !ok {
				db.tmpIB.data[group8Bits[0]] = map[uint8]map[uint8]map[uint8]block{}
			}
			if _, ok := db.tmpIB.data[group8Bits[0]][group8Bits[1]]; !ok {
				db.tmpIB.data[group8Bits[0]][group8Bits[1]] = map[uint8]map[uint8]block{}
			}
			if _, ok := db.tmpIB.data[group8Bits[0]][group8Bits[1]][group8Bits[2]]; !ok {
				db.tmpIB.data[group8Bits[0]][group8Bits[1]][group8Bits[2]] = map[uint8]block{}
			}
			// ISO 3166 2-letter に定義されるカントリーコードを示す。
			// メモリ使用量削減のため、文字列からひも付けされた uint8 に
			// 変換して格納。
			if _, ok := db.tmpIB.dicCCStrToInt[line[1]]; !ok {
				db.tmpIB.dicCCStrToInt[line[1]] = countForDic
				countForDic++
			}
			// Record format の５番めの Field は value 。
			// ipv4 の場合、対象範囲のアドレスの個数を示す。
			// int に変換できる場合。
			v, err := strconv.Atoi(line[4])
			if err == nil {
				// ここまで異常がなければ各データを格納する。
				if _, ok := db.tmpIB.data[group8Bits[0]][group8Bits[1]][group8Bits[2]][group8Bits[3]]; ok {
					db.tmpIB.totalValue["ALL"] = db.tmpIB.totalValue["ALL"] - int(db.tmpIB.data[group8Bits[0]][group8Bits[1]][group8Bits[2]][group8Bits[3]].value)
					db.tmpIB.totalValue[line[1]] = db.tmpIB.totalValue[line[1]] - int(db.tmpIB.data[group8Bits[0]][group8Bits[1]][group8Bits[2]][group8Bits[3]].value)
				} else {
					db.tmpIB.totalBlocks["ALL"]++
					db.tmpIB.totalBlocks[line[1]]++
				}
				db.tmpIB.totalValue["ALL"] = db.tmpIB.totalValue["ALL"] + v
				db.tmpIB.totalValue[line[1]] = db.tmpIB.totalValue[line[1]] + v
				db.tmpIB.data[group8Bits[0]][group8Bits[1]][group8Bits[2]][group8Bits[3]] = block{
					country: db.tmpIB.dicCCStrToInt[line[1]],
					// uint32 に変換して格納。
					value: uint32(v),
				}
			} else {
				db.ClearTmpIPBData()
				return fmt.Errorf("%v: %v", err, line)
			}
		}
	}

	// ひも付けされた uint8 からカントリーコードの文字列を逆引きするために使用。
	for k := range db.tmpIB.dicCCStrToInt {
		db.tmpIB.dicCCIntToStr[db.tmpIB.dicCCStrToInt[k]] = k
	}

	return nil
}

// カントリーコードの一覧ファイルを読み込む。
func (db *DB) SetTmpCountryCodes(r io.Reader) error {
	db.tmpCC.l.Lock()
	defer db.tmpCC.l.Unlock()

	if db.tmpCC.data == nil {
		db.tmpCC.data = map[string]CountryCodeInfo{}
	}

	// ファイルを csv として読込。
	// コメント・フィールド区切りの文字を設定。
	reader := csv.NewReader(r)
	reader.Comment = '#'
	reader.Comma = '|'

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				db.tmpCC.data = map[string]CountryCodeInfo{}
				return fmt.Errorf("%v: %v", err, line)
			}
		}
		// フィールド数は３。
		if len(line) != 3 {
			db.tmpCC.data = map[string]CountryCodeInfo{}
			return fmt.Errorf("the number(%d) of the line's fields is invalid: %v", len(line), line)
		}
		// 先頭フィールドがカントリーコードで、英大文字２文字。
		if !db.reg.MatchString(line[0]) {
			db.tmpCC.data = map[string]CountryCodeInfo{}
			return fmt.Errorf("the line's first field (country code: %s) is invalid: %v", line[0], line)
		}

		// カントリーコードをキーとする listCCName に
		// 英語国名、別の言語の国名を格納。
		db.tmpCC.data[line[0]] = CountryCodeInfo{
			Name:    line[1],
			AltName: line[2],
		}
	}

	return nil
}

// 指定された URL のデータを取得し、
// 一時保存用データベースに格納する。
func (db *DB) LoadIPBDataByURL(u string) error {
	// 指定された URL が適正なものかを確認
	if _, err := url.ParseRequestURI(u); err != nil {
		return err
	}

	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = db.setTmpIPBlocks(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// 指定のファイルを読んでIPアドレスの国別ブロックのデータを取得し、
// 一時保存用データベースに格納する。
func (db *DB) LoadIPBDataByFile(ipbFile string) error {
	fp, err := os.Open(ipbFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	err = db.setTmpIPBlocks(fp)
	if err != nil {
		return err
	}

	return nil
}

// 一時保存用と検索用の両方のIPアドレスの
// 国別ブロックデータベースをロックし、
// 一時保存用のデータを検索用に渡す。
// 一時保存用のデータは空にする。
func (db *DB) SwitchIPBData() {
	db.tmpIB.l.Lock()
	db.ib.l.Lock()
	db.ib.data = db.tmpIB.data
	db.ib.dicCCIntToStr = db.tmpIB.dicCCIntToStr
	db.ib.dicCCStrToInt = db.tmpIB.dicCCStrToInt
	db.ib.totalBlocks = db.tmpIB.totalBlocks
	db.ib.totalValue = db.tmpIB.totalValue
	db.ClearTmpIPBData()
	db.ib.l.Unlock()
	db.tmpIB.l.Unlock()
}

// 初期設定済の URL から各 RIR の最新版 delegation file を取得し、
// IPアドレスの国別ブロックデータベースを更新する。
func (db *DB) SetIPBData() error {

	for i := range db.urlRIR {
		if err := db.LoadIPBDataByURL(db.urlRIR[i]); err != nil {
			return err
		}
	}

	db.SwitchIPBData()

	return nil
}

// 一時保存用と検索用の両方のカントリーコードの
// 一覧のデータベースをロックし、
// 一時保存用のデータを検索用に渡す。
// 一時保存用のデータは空にする。
func (db *DB) SwitchCCData() {
	db.tmpCC.l.Lock()
	db.cc.l.Lock()
	db.cc.data = db.tmpCC.data
	db.cc.l.Unlock()
	db.tmpCC.data = map[string]CountryCodeInfo{}
	db.tmpCC.l.Unlock()
}

// 指定のファイルを読んでカントリーコードの一覧のデータを取得し、
// 一時保存用データベースに格納する。
// 一時保存用と検索用の両方のデータベースをロックし、
// 一時保存用のデータを検索用へコピーする。
// 一時保存用のデータは空にする。
func (db *DB) InitCCDataByFile(ccFile string) error {
	fp, err := os.Open(ccFile)
	if err != nil {
		return err
	}
	defer fp.Close()
	err = db.SetTmpCountryCodes(fp)
	if err != nil {
		return err
	}

	db.SwitchCCData()

	return nil
}

// IPアドレスの最初の8ビットについて、
// データベースに一致する値があるか判定
func (db *DB) checkFirst8Bit(as4 [4]byte) ([4]byte, bool) {
	// 一致があればそのまま引数を返す。
	if _, ok := db.ib.data[as4[0]]; ok {
		return as4, true
	}
	// 元の値から1づつ減らして
	// 一致する値を検索
	i := as4[0]
	for {
		if _, ok := db.ib.data[i]; ok {
			as4[0] = i
			as4[1] = 255
			as4[2] = 255
			as4[3] = 255
			return as4, true
		}
		if i == 0 {
			return as4, false
		}
		i--
	}
}

// IPアドレスの2番めの8ビットについて、
// データベースに一致する値があるか判定
func (db *DB) checkSecond8Bit(as4 [4]byte) ([4]byte, bool) {
	// 一致があればそのまま引数を返す。
	if _, ok := db.ib.data[as4[0]][as4[1]]; ok {
		return as4, true
	}
	// データベースの2番めの8ビットが0のみの場合は0
	// 以降の8ビットは255
	if _, ok := db.ib.data[as4[0]][0]; ok && len(db.ib.data[as4[0]]) == 1 {
		as4[1] = 0
		as4[2] = 255
		as4[3] = 255
		return as4, true
	}
	// それ以外は元の値から1づつ減らして一致する値を検索
	i := as4[1]
	for {
		// 一致した場合、以降の8ビットは255
		if _, ok := db.ib.data[as4[0]][i]; ok {
			as4[1] = i
			as4[2] = 255
			as4[3] = 255
			return as4, true
		}
		// 一致がなく、0に達した場合。
		if i == 0 {
			return as4, false
		}
		i--
	}
}

// IPアドレスの3番めの8ビットについて、
// データベースに一致する値があるか判定
func (db *DB) checkThird8Bit(as4 [4]byte) ([4]byte, bool) {
	// 一致があればそのまま引数を返す。
	if _, ok := db.ib.data[as4[0]][as4[1]][as4[2]]; ok {
		return as4, true
	}
	// 3番めの8ビットが0のみの場合は0
	// 以降の8ビットは255
	if _, ok := db.ib.data[as4[0]][as4[1]][0]; ok && len(db.ib.data[as4[0]][as4[1]]) == 1 {
		as4[2] = 0
		as4[3] = 255
		return as4, true
	}
	// それ以外は元の値から1づつ減らして一致する値を検索
	i := as4[2]
	for {
		// 一致した場合、以降の8ビットは255
		if _, ok := db.ib.data[as4[0]][as4[1]][i]; ok {
			as4[2] = i
			as4[3] = 255
			return as4, true
		}
		// 一致がなく、0に達した場合。
		if i == 0 {
			return as4, false
		}
		i--
	}
}

// IPアドレスの最後の8ビットについて
// データベースに一致する値があるか判定
func (db *DB) checkLast8Bit(as4 [4]byte) ([4]byte, bool) {
	// 一致があればそのまま引数を返す。
	if _, ok := db.ib.data[as4[0]][as4[1]][as4[2]][as4[3]]; ok {
		return as4, true
	}
	// 最後の8ビットが0のみの場合は0
	if _, ok := db.ib.data[as4[0]][as4[1]][as4[2]][0]; ok && len(db.ib.data[as4[0]][as4[1]][as4[2]]) == 1 {
		as4[3] = 0
		return as4, true
	}
	// それ以外は元の値から1づつ減らして一致する値を検索
	i := as4[3]
	for {
		// 一致した場合
		if _, ok := db.ib.data[as4[0]][as4[1]][as4[2]][i]; ok {
			as4[3] = i
			return as4, true
		}
		// 一致がなく、0に達した場合。
		if i == 0 {
			return as4, false
		}
		i--
	}
}

func (db *DB) searchBlockStart(addr netip.Addr) [4]byte {
	// IPv4アドレスを8ビット単位で分割し、
	// データベースから所属ブロック候補を検索
	as4 := addr.As4()
	f := 2
	found := true
	for {
		if f == 2 {
			as4, found = db.checkFirst8Bit(as4)
			if !found {
				return [4]byte{}
			}
		}

		if f >= 1 {
			as4, found = db.checkSecond8Bit(as4)
			if !found {
				if as4[0] == 0 {
					return [4]byte{}
				}
				as4[0]--
				as4[1] = 255
				as4[2] = 255
				as4[3] = 255
				f = 2
				continue
			}
		}

		as4, found = db.checkThird8Bit(as4)
		if !found {
			if as4[1] == 0 {
				if as4[0] == 0 {
					return [4]byte{}
				}
				as4[0]--
				as4[1] = 255
				as4[2] = 255
				as4[3] = 255
				f = 2
				continue
			}
			as4[1]--
			as4[2] = 255
			as4[3] = 255
			f = 1
			continue
		}

		as4, found = db.checkLast8Bit(as4)
		if found {
			return as4
		}
		if as4[2] == 0 {
			if as4[1] == 0 {
				if as4[0] == 0 {
					return [4]byte{}
				}
				as4[0]--
				as4[1] = 255
				as4[2] = 255
				as4[3] = 255
				f = 2
				continue
			}
			as4[1]--
			as4[2] = 255
			as4[3] = 255
			f = 1
			continue
		}
		as4[2]--
		as4[3] = 255
		f = 0
	}
}

// 渡された文字列のIPv4アドレスからカントリーコードの情報を返す。
func (db *DB) SearchInfo(adrs string) SearchResult {
	db.ib.l.RLock()
	defer db.ib.l.RUnlock()
	db.cc.l.RLock()
	defer db.cc.l.RUnlock()

	target, err := netip.ParseAddr(adrs)
	// 渡された文字列をパースしてエラー
	if err != nil {
		return SearchResult{Message: "Invalid IP Address"}
	}
	// IPv4アドレスでない
	if !target.Is4() {
		return SearchResult{Message: "Not IPv4 Address"}
	}
	// ループバックアドレス
	if target.IsLoopback() {
		return SearchResult{Message: "Loopback Address"}
	}
	// マルチキャストアドレス
	if target.IsMulticast() {
		return SearchResult{Message: "Multicast Address"}
	}
	// プライベートアドレス
	if target.IsPrivate() {
		return SearchResult{Message: "Private Address"}
	}

	// IPv4アドレスを8ビット単位で分割し、データベースから所属ブロック候補を検索
	as4 := db.searchBlockStart(target)
	if as4 == [4]byte{} {
		return SearchResult{Message: "Not Found"}
	}

	// 所属ブロック候補の最後のアドレスを計算
	oO := getOneOutside(as4, db.ib.data[as4[0]][as4[1]][as4[2]][as4[3]].value)

	// 渡されたIPv4アドレスが所属ブロック候補の範囲に含まれる場合は、
	// カントリーコード他該当情報を返す。
	if target.Less(oO) {
		sr := SearchResult{
			IsFound:    true,
			Message:    "Found",
			BlockStart: netip.AddrFrom4(as4).String(),
			BlockEnd:   oO.Prev().String(),
			Code:       db.ib.dicCCIntToStr[db.ib.data[as4[0]][as4[1]][as4[2]][as4[3]].country],
		}
		if _, ok := db.cc.data[sr.Code]; ok {
			sr.Name = db.cc.data[sr.Code].Name
			sr.AltName = db.cc.data[sr.Code].AltName
		}

		return sr
	}

	// 渡されたIPv4アドレスが所属ブロック候補の範囲に含まれない場合は、
	// 情報なしとして返す。
	return SearchResult{Message: "Not Found"}

}

// 渡された IP アドレスに対応する 4 バイトの配列とUint32 に
// 変換された RIR statistics exchange format の value の値から
// ブロック範囲外最初の IP アドレスを計算して返す。
func getOneOutside(a4b [4]byte, value uint32) netip.Addr {
	s4b := []byte{a4b[0], a4b[1], a4b[2], a4b[3]}
	binary.BigEndian.PutUint32(s4b, binary.BigEndian.Uint32(s4b)+value)
	return netip.AddrFrom4([4]byte{s4b[0], s4b[1], s4b[2], s4b[3]})
}

// string で渡された IP アドレスと RIR statistics exchange format の
// value の値からブロック範囲最後の IP アドレスを計算して返す。
// value は 1 〜 4294967295 の範囲でなければならない。
// value によりブロック範囲が 255.255.255.255 を超える場合は、
// 0.0.0.0 から超えた value 分がブロック範囲となる。
func GetLastAddr(adrs string, value int) (netip.Addr, error) {
	target, err := netip.ParseAddr(adrs)
	// 渡されたアドレスをパースしてエラー
	if err != nil {
		return netip.Addr{}, err
	}
	// IPv4アドレスでない
	if !target.Is4() {
		return netip.Addr{}, ErrFirstArgumentOutOfRange
	}

	if value < 1 || value > 4294967295 {
		return netip.Addr{}, ErrSecondArgumentOutOfRange
	}

	return getOneOutside(target.As4(), uint32(value)).Prev(), nil
}

// string で渡された 2 つの IP アドレスから
// ブロックのアドレス個数を計算して返す。
func GetValue(a string, b string) (int, error) {
	x, err := netip.ParseAddr(a)
	// 渡されたアドレスをパースしてエラー
	if err != nil {
		return 0, err
	}
	// IPv4アドレスでない
	if !x.Is4() {
		return 0, ErrFirstArgumentOutOfRange
	}

	y, err := netip.ParseAddr(b)
	// 渡されたアドレスをパースしてエラー
	if err != nil {
		return 0, err
	}
	// IPv4アドレスでない
	if !y.Is4() {
		return 0, ErrSecondArgumentOutOfRange
	}

	i := x.Compare(y)
	if i == 0 {
		return 0, nil
	}
	if i == 1 {
		return int(binary.BigEndian.Uint32(x.AsSlice()) - binary.BigEndian.Uint32(y.AsSlice()) + 1), nil
	}
	return int(binary.BigEndian.Uint32(y.AsSlice()) - binary.BigEndian.Uint32(x.AsSlice()) + 1), nil
}

// 国別ブロック合計を取得する。
func (db *DB) GetTotalBlocks() map[string]int {
	return db.ib.totalBlocks
}

// 国別アドレス数合計を取得する。
func (db *DB) GetTotalValue() map[string]int {
	return db.ib.totalValue
}

// 各カントリーコードの名前情報を取得する。
func (db *DB) GetCountryCodeData() map[string]CountryCodeInfo {
	return db.cc.data
}

// 検索用データベースが空ならば true を返す。
func (db *DB) IsDBEmpty() bool {
	return len(db.ib.data) == 0
}
