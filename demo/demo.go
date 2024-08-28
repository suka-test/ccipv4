package main

import (
	"fmt"
	"io"
	"net/netip"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/suka-test/ccipv4"
)

const (
	// カントリーコード一覧のデータ
	COUNTRY_CODES string = `AD|Andorra|アンドラ
AE|United Arab Emirates|アラブ首長国連邦
AF|Afghanistan|アフガニスタン
AG|Antigua and Barbuda|アンティグア・バーブーダ
AI|Anguilla|アンギラ
AL|Albania|アルバニア
AM|Armenia|アルメニア
AO|Angola|アンゴラ
AR|Argentine|アルゼンチン
AS|American Samoa|アメリカ領サモア
AT|Austria|オーストリア
AU|Australia|オーストラリア
AW|Aruba|アルバ
AX|Åland Islands|オーランド諸島
AZ|Azerbaijan|アゼルバイジャン
BA|Bosnia and Herzegovina|ボスニア・ヘルツェゴビナ
BB|Barbados|バルバドス
BD|Bangladesh|バングラデシュ
BE|Belgium|ベルギー
BF|Burkina Faso|ブルキナファソ
BG|Bulgaria|ブルガリア
BH|Bahrain|バーレーン
BI|Burundi|ブルンジ
BJ|Benin|ベナン
BL|Saint Barthélemy|サン・バルテルミー島
BM|Bermuda|バミューダ諸島
BN|Brunei Darussalam|ブルネイ・ダルサラーム
BO|Bolivia|ボリビア
BQ|Bonaire, Sint Eustatius and Saba|ボネール、シント・ユースタティウス及びサバ
BR|Brazil|ブラジル
BS|Bahamas|バハマ
BT|Bhutan|ブータン
BW|Botswana|ボツワナ
BY|Belarus|ベラルーシ
BZ|Belize|ベリーズ
CA|Canada|カナダ
CD|Democratic Republic of the Congo|コンゴ民主共和国
CF|Central African|中央アフリカ
CG|Republic of the Congo|コンゴ共和国
CH|Swiss|スイス
CI|Cote d'Ivoire|コートジボワール
CK|Cook Islands|クック諸島
CL|Chile|チリ
CM|Cameroon|カメルーン
CN|China|中国
CO|Colombia|コロンビア
CR|Costa Rica|コスタリカ
CU|Cuba|キューバ
CV|Cabo Verde|カーボベルデ
CW|Curaçao|キュラソー島
CY|Cyprus|キプロス
CZ|Czech|チェコ
DE|Germany|ドイツ
DJ|Djibouti|ジブチ
DK|Denmark|デンマーク
DM|Commonwealth of Dominica|ドミニカ国
DO|Dominican Republic|ドミニカ共和国
DZ|Algeria|アルジェリア
EC|Ecuador|エクアドル
EE|Estonia|エストニア
EG|Egypt|エジプト
ER|Eritrea|エリトリア
ES|Spain|スペイン
ET|Ethiopia|エチオピア
EU|European Union|欧州連合
FI|Finland|フィンランド
FJ|Fiji|フィジー
FK|Falkland Islands|フォークランド諸島
FM|Micronesia|ミクロネシア
FO|Faroe Islands|フェロー諸島
FR|French|フランス
GA|Gabonese|ガボン
GB|United Kingdom|イギリス
GD|Grenada|グレナダ
GE|Georgia|ジョージア
GF|French Guiana|フランス領ギアナ
GG|Guernsey|ガーンジー島
GH|Ghana|ガーナ
GI|Gibraltar|ジブラルタル
GL|Greenland|グリーンランド
GM|Gambia|ガンビア
GN|Guinea|ギニア
GP|Guadeloupe|グアドループ
GQ|Equatorial Guinea|赤道ギニア
GR|Hellenic|ギリシャ
GT|Guatemala|グアテマラ
GU|Guam|グアム
GW|Guinea-Bissau|ギニアビサウ
GY|Guyana|ガイアナ
HK|Hong Kong|香港
HN|Honduras|ホンジュラス
HR|Croatia|クロアチア
HT|Haiti|ハイチ
HU|Hungary|ハンガリー
ID|Indonesia|インドネシア
IE|Ireland|アイルランド
IL|Israel|イスラエル
IM|Isle of Man|マン島
IN|India|インド
IO|British Indian Ocean Territory|イギリス領インド洋地域
IQ|Iraq|イラク
IR|Iran|イラン
IS|Iceland|アイスランド
IT|Italian|イタリア
JE|Jersey|ジャージー
JM|Jamaica|ジャマイカ
JO|Jordan|ヨルダン
JP|Japan|日本
KE|Kenya|ケニア
KG|Kyrgyz|キルギス
KH|Cambodia|カンボジア
KI|Kiribati|キリバス
KM|Comoros|コモロ
KN|Saint Christopher and Nevis|セントクリストファー・ネービス
KP|North Korea|北朝鮮
KR|Korea|韓国
KW|Kuwait|クウェート
KY|Cayman Islands|ケイマン諸島
KZ|Kazakhstan|カザフスタン
LA|Lao|ラオス
LB|Lebanese|レバノン
LC|Saint Lucia|セントルシア
LI|Liechtenstein|リヒテンシュタイン
LK|Sri Lanka|スリランカ
LR|Liberia|リベリア
LS|Lesotho|レソト
LT|Lithuania|リトアニア
LU|Luxembourg|ルクセンブルク
LV|Latvia|ラトビア
LY|Libya|リビア
MA|Morocco|モロッコ
MC|Monaco|モナコ
MD|Moldova|モルドバ
ME|Montenegro|モンテネグロ
MF|Saint Martin (French part)|フランス領サン・マルタン
MG|Madagascar|マダガスカル
MH|Marshall Islands|マーシャル諸島
MK|North Macedonia|北マケドニア
ML|Mali|マリ
MM|Myanmar|ミャンマー
MN|Mongolia|モンゴル国
MO|Macau|マカオ
MP|Northern Mariana Islands|北マリアナ諸島
MQ|Martinique|マルティニーク
MR|Mauritania|モーリタニア
MS|Montserrat|モントセラト
MT|Malta|マルタ
MU|Mauritius|モーリシャス
MV|Maldives|モルディブ
MW|Malawi|マラウイ
MX|Mexico|メキシコ
MY|Malaysia|マレーシア
MZ|Mozambique|モザンビーク
NA|Namibia|ナミビア
NC|New Caledonia|ニューカレドニア
NE|Niger|ニジェール
NF|Norfolk Island|ノーフォーク島
NG|Nigeria|ナイジェリア
NI|Nicaragua|ニカラグア
NL|Netherlands|オランダ
NO|Norway|ノルウェー
NP|Nepal|ネパール
NR|Nauru|ナウル
NU|Niue|ニウエ
NZ|New Zealand|ニュージーランド
OM|Sultanate of Oman|オマーン
PA|Panama|パナマ
PE|Peru|ペルー
PF|French Polynesia|フランス領ポリネシア
PG|Papua New Guinea|パプアニューギニア
PH|Philippines|フィリピン
PK|Pakistan|パキスタン
PL|Poland|ポーランド
PM|Saint Pierre and Miquelon|サンピエール島及びミクロン島
PR|Puerto Rico|プエルトリコ
PS|Palestine|パレスチナ
PT|Portuguese|ポルトガル
PW|Palau|パラオ
PY|Paraguay|パラグアイ
QA|Qatar|カタール
RE|Réunion|レユニオン
RO|Romania|ルーマニア
RS|Serbia|セルビア
RU|Russia|ロシア
RW|Rwanda|ルワンダ
SA|Saudi Arabia|サウジアラビア
SB|Solomon Islands|ソロモン諸島
SC|Seychelles|セーシェル
SD|Sudan|スーダン
SE|Sweden|スウェーデン
SG|Singapore|シンガポール
SI|Slovenia|スロベニア
SK|Slovak|スロバキア
SL|Sierra Leone|シエラレオネ
SM|San Marino|サンマリノ
SN|Senegal|セネガル
SO|Somalia|ソマリア
SR|Suriname|スリナム
SS|South Sudan|南スーダン
ST|Sao Tome and Principe|サントメ・プリンシペ
SV|El Salvador|エルサルバドル
SX|Sint Maarten (Dutch part)|オランダ領シント・マールテン
SY|Syrian Arab|シリア・アラブ
SZ|Eswatini|エスワティニ
TC|Turks and Caicos Islands|タークス・カイコス諸島
TD|Chad|チャド
TG|Togo|トーゴ
TH|Thailand|タイ
TJ|Tajikistan|タジキスタン
TK|Tokelau|トケラウ
TL|Timor-Leste|東ティモール
TM|Turkmenistan|トルクメニスタン
TN|Tunisia|チュニジア
TO|Kingdom of Tonga|トンガ
TR|Turkey|トルコ
TT|Trinidad and Tobago|トリニダード・トバゴ
TV|Tuvalu|ツバル
TW|Taiwan|台湾
TZ|Tanzania|タンザニア
UA|Ukraine|ウクライナ
UG|Uganda|ウガンダ
US|America|アメリカ
UY|Uruguay|ウルグアイ
UZ|Uzbekistan|ウズベキスタン
VA|Vatican|バチカン
VC|Saint Vincent and the Grenadines|セントビンセント及びグレナディーン諸島
VE|Venezuela|ベネズエラ
VG|Virgin Islands (British)|イギリス領ヴァージン諸島
VI|Virgin Islands (U.S.)|アメリカ領ヴァージン諸島
VN|Viet Nam|ベトナム
VU|Vanuatu|バヌアツ
WF|Wallis and Futuna|ウォリス・フツナ
WS|Samoa|サモア
YE|Yemen|イエメン
YT|Mayotte|マヨット
ZA|South Africa|南アフリカ
ZM|Zambia|ザンビア
ZW|Zimbabwe|ジンバブエ
ZZ|Unknown|不明
`
	// 表示メッセージ
	MSG_EMPTY_INPUT        string = "入力が不足しています。"
	MSG_NOT_IPV4           string = "IPv4アドレスではありません。"
	MSG_VALUE_OUT_OF_RANGE string = "値の範囲は1 〜 4294967295です。"
	MSG_UNEXPECTED_ERROR   string = "予期しないエラーが発生しました。"

	// デモ用IPアドレスの国別ブロックのデータ
	DEMO_DATA string = "demo_data.csv"
)

type cli struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
	db     *ccipv4.DB
	cc     string
	rir    [][]string
}

// struct cli の実体を一つ作り、
// run から諸機能を動かす。
func main() {
	c := &cli{
		stdout: os.Stdout,
		stderr: os.Stderr,
		stdin:  os.Stdin,
		db:     ccipv4.GetDB(),
		cc:     COUNTRY_CODES,
		rir:    getRIRURL(),
	}

	fmt.Fprintln(c.stdout, "")
	fmt.Fprintln(c.stderr, "                            <<<< ccipv4 demo >>>>")

	if err := c.run(); err != nil {
		fmt.Fprintln(c.stderr, MSG_UNEXPECTED_ERROR)
	}
	fmt.Fprintln(c.stdout, "")
	fmt.Fprintln(c.stderr, "                 終了します。")
	fmt.Fprintln(c.stdout, "")
}

// カントリーコード一覧のデータベースと
// 各種機能用のコマンド受付の準備をする。
func (c *cli) run() error {
	// IPアドレスの国別ブロックのデータベースを準備する。
	var path string
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	switch filepath.Base(dir) {
	// カレントディレクトリが demo であれば直上ディレクトリの
	// 直下に samples があると仮定してデータファイルを探す。
	case "demo":
		path = filepath.Join(filepath.Dir(dir), "samples", DEMO_DATA)
	// カレントディレクトリが ccipv4 であれば
	// 直下に samples があると仮定してデータファイルを探す。
	case "ccipv4":
		path = filepath.Join(dir, "samples", DEMO_DATA)
	// カレントディレクトリでデータファイルを探す。
	default:
		path = DEMO_DATA
	}
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		if err := c.db.LoadIPBDataByFile(path); err != nil {
			return err
		}
		c.db.SwitchIPBData()
	}

	// カントリーコード一覧のデータベースを準備する。
	if err := c.db.SetTmpCountryCodes(strings.NewReader(c.cc)); err != nil {
		return err
	}
	c.db.SwitchCCData()

	c.getCommand()

	return nil
}

// 各種機能用のコマンド受付を行う。
func (c *cli) getCommand() {
	var s string
	isFirstTime := true
	for {
		if c.db.IsDBEmpty() {
			fmt.Fprintln(c.stdout, "")
			fmt.Fprintln(c.stdout, "    \x1b[41m !!!!! データベースが空です。検索を使う場合は更新してください。 !!!!! \x1b[0m")
		}
		if isFirstTime {
			fmt.Fprintln(c.stdout, "")
			fmt.Fprintln(c.stdout, "    \x1b[41m ！注意：本プログラムはデモ用です。                                   \x1b[0m")
			fmt.Fprintln(c.stdout, "    \x1b[41m 起動時のデータは最新ではありません。                                 \x1b[0m")
			fmt.Fprintln(c.stdout, "    \x1b[41m IPv4アドレスの国別ブロックの最新データが必要な場合は、               \x1b[0m")
			fmt.Fprintln(c.stdout, "    \x1b[41m データベース更新を行ってください。                                   \x1b[0m")
			fmt.Fprintln(c.stdout, "    \x1b[41m 更新時には２５メガバイト程度のデータダウンロードを行います。         \x1b[0m")
			fmt.Fprintln(c.stdout, "    \x1b[41m 通信状況相応の時間がかかります。予めご了承ください。                 \x1b[0m")
			isFirstTime = false
		}
		fmt.Fprintln(c.stdout, "")
		fmt.Fprintln(c.stdout, "          \x1b[47m \x1b[30mq \x1b[0m : 終了 \x1b[47m \x1b[30md \x1b[0m : データベース更新 \x1b[47m \x1b[30ms \x1b[0m : 検索 \x1b[47m \x1b[30mi \x1b[0m : 集計情報")
		fmt.Fprintln(c.stdout, "")
		fmt.Fprintln(c.stdout, "          \x1b[47m \x1b[30ma \x1b[0m : ブロック最後のIPv4アドレス \x1b[47m \x1b[30mv \x1b[0m : ブロックのアドレス数")
		fmt.Fprintln(c.stdout, "")
		fmt.Fprintln(c.stdout, "                 \x1b[44m 上のコマンド文字１つを選択して入力＋Enter \x1b[0m")
		fmt.Fprint(c.stdout, "                 \x1b[44m >>> \x1b[0m ")
		if _, err := fmt.Fscanln(c.stdin, &s); err != nil {
			continue
		}
		switch s {
		// 終了
		case "q":
			return
		// IPアドレスの国別ブロックデータベース更新
		case "d":
			fmt.Fprintln(c.stdout, "")
			c.loadIPBD()
			continue
		// 検索
		case "s":
			fmt.Fprintln(c.stdout, "")
			c.searchIPB()
			continue
		// ブロック範囲最後の IP アドレス取得
		case "a":
			fmt.Fprintln(c.stdout, "")
			c.getLastAddr()
			continue
		// ブロックのアドレス個数
		case "v":
			fmt.Fprintln(c.stdout, "")
			c.getValue()
			continue
			// 集計情報
		case "i":
			fmt.Fprintln(c.stdout, "")
			c.getInfo()
			continue
		}
	}
}

// 初期設定済の URL から各 RIR の最新版 delegation file を取得し、
// IPアドレスの国別ブロックデータベースを更新する。
func (c *cli) loadIPBD() {
	var err error
	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mデータベースを更新します。                \x1b[0m")
	fmt.Fprintln(c.stdout, "")
	for i := range c.rir {
		fmt.Fprintf(c.stdout, "     %-7s からデータをダウンロードします。 >>> ", c.rir[i][0])
		t := time.Now()
		if err = c.db.LoadIPBDataByURL(c.rir[i][1]); err != nil {
			break
		}
		fmt.Fprintf(c.stdout, "経過時間 : %10f秒\n", time.Since(t).Seconds())
	}
	if err != nil {
		c.db.ClearTmpIPBData()
		fmt.Fprintln(c.stdout, "")
		fmt.Fprintln(c.stdout, "")
		fmt.Fprintf(c.stdout, "                 \x1b[41m !! %-19s !! \x1b[0m\n", MSG_UNEXPECTED_ERROR)
		fmt.Fprintf(c.stdout, "                 \x1b[41m %-23s \x1b[0m\n", "データベースを更新できませんでした。")
		return
	}
	c.db.SwitchIPBData()
	fmt.Fprintln(c.stdout, "")
	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mデータベース更新、終了しました。          \x1b[0m")
}

// 各 RIR の最新版 delegation file の URL を取得する。
func getRIRURL() [][]string {
	return [][]string{
		{"AFRINIC", ccipv4.URL_DELEGATED_AFRINIC_EXTENDED_LATEST},
		{"APNIC", ccipv4.URL_DELEGATED_APNIC_EXTENDED_LATEST},
		{"ARIN", ccipv4.URL_DELEGATED_ARIN_EXTENDED_LATEST},
		{"LACNIC", ccipv4.URL_DELEGATED_LACNIC_EXTENDED_LATEST},
		{"RIPENCC", ccipv4.URL_DELEGATED_RIPENCC_EXTENDED_LATEST},
	}
}

// 検索
func (c *cli) searchIPB() {
	var s string
	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mIPv4アドレス入力＋Enter                   \x1b[0m")
	fmt.Fprint(c.stdout, "                 \x1b[47m \x1b[30m>>> \x1b[0m ")
	if _, err := fmt.Fscanln(c.stdin, &s); err != nil {
		s = ""
	}

	res := c.db.SearchInfo(s)
	fmt.Fprintln(c.stdout, "")
	if res.IsFound {
		fmt.Fprintf(c.stdout, "                 ブロック : %s 〜 %s\n", res.BlockStart, res.BlockEnd)
		fmt.Fprintf(c.stdout, "                 コード   : %s\n", res.Code)
		if res.Name != "" {
			fmt.Fprintf(c.stdout, "                 名称     : %s\n", res.Name)
			fmt.Fprintf(c.stdout, "                            %s\n", res.AltName)
		}
	} else {
		switch res.Message {
		case "Invalid IP Address":
			fmt.Fprintln(c.stdout, "                 IPアドレスではありません。")
		case "Not IPv4 Address":
			fmt.Fprintln(c.stdout, "                 IPv4アドレスではありません。")
		case "Loopback Address":
			fmt.Fprintln(c.stdout, "                 ループバックアドレスです。")
		case "Multicast Address":
			fmt.Fprintln(c.stdout, "                 マルチキャストアドレスです。")
		case "Private Address":
			fmt.Fprintln(c.stdout, "                 プライベートアドレスです。")
		case "Not Found":
			fmt.Fprintln(c.stdout, "                 該当するブロックはありませんでした。")
		}
	}
}

// ブロック最初の IP アドレスと value（アドレス数）から
// ブロック最後の IP アドレスを計算して取得する。
func (c *cli) getLastAddr() {
	var (
		a, b string
		i    int
	)

	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mブロック最後のIPv4アドレスを表示します。  \x1b[0m")
	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mブロック先頭のIPv4アドレスと              \x1b[0m")
	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mアドレス数（value）の入力が必要です。     \x1b[0m")
	fmt.Fprintln(c.stdout, "")

	for {
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mブロック先頭のIPv4アドレスを入力＋Enter   \x1b[0m")
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30m（bを入力＋Enterで中断しコマンド選択へ）  \x1b[0m")
		fmt.Fprint(c.stdout, "                 \x1b[47m \x1b[30m>>> \x1b[0m ")
		if _, err := fmt.Fscanln(c.stdin, &a); err == nil {
			if a == "b" {
				return
			}
			if _, err = netip.ParseAddr(a); err == nil {
				break
			}
		}
		fmt.Fprintln(c.stdout, "                 \x1b[41m !!!!! 誤入力です。 !!!!! \x1b[0m")
		fmt.Fprintln(c.stdout, "")
	}

	fmt.Fprintln(c.stdout, "")

	for {
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mブロックのアドレス数（value）を           \x1b[0m")
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30m1〜4294967295の範囲でを入力＋Enter        \x1b[0m")
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30m（bを入力＋Enterで中断しコマンド選択へ）  \x1b[0m")
		fmt.Fprint(c.stdout, "                 \x1b[47m \x1b[30m>>> \x1b[0m ")
		if _, err := fmt.Fscanln(c.stdin, &b); err == nil {
			if b == "b" {
				return
			}
			i, err = strconv.Atoi(b)
			if err == nil && i >= 1 && i <= 4294967295 {
				break
			}
		}
		fmt.Fprintln(c.stdout, "                 \x1b[41m !!!!! 誤入力です。 !!!!! \x1b[0m")
		fmt.Fprintln(c.stdout, "")
	}

	addr, err := ccipv4.GetLastAddr(a, i)
	fmt.Fprintln(c.stdout, "")
	if err == nil {
		fmt.Fprintln(c.stdout, "                 "+addr.String())
	} else {
		var e string
		switch err {
		case ccipv4.ErrFirstArgumentOutOfRange:
			e = MSG_NOT_IPV4
		case ccipv4.ErrSecondArgumentOutOfRange:
			e = MSG_VALUE_OUT_OF_RANGE
		default:
			e = MSG_UNEXPECTED_ERROR
		}
		fmt.Fprintf(c.stdout, "                 \x1b[41m %-29s \x1b[0m\n", e)
	}
}

// 2 つの IP アドレスからそのブロックの
// アドレス個数を計算して取得する。
func (c *cli) getValue() {
	var a, b string

	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mブロックのアドレス個数を表示します。      \x1b[0m")
	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mブロックの先頭と最後のIPv4アドレスの入力が\x1b[0m")
	fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30m必要です。                                \x1b[0m")
	fmt.Fprintln(c.stdout, "")

	for {
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mブロック先頭のIPv4アドレスを入力＋Enter   \x1b[0m ")
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30m（bを入力＋Enterで中断しコマンド選択へ）  \x1b[0m")
		fmt.Fprint(c.stdout, "                 \x1b[47m \x1b[30m>>> \x1b[0m ")
		if _, err := fmt.Fscanln(c.stdin, &a); err == nil {
			if a == "b" {
				return
			}
			if _, err = netip.ParseAddr(a); err == nil {
				break
			}
		}
		fmt.Fprintln(c.stdout, "                 \x1b[41m !!!!! 誤入力です。 !!!!! \x1b[0m")
		fmt.Fprintln(c.stdout, "")
	}

	fmt.Fprintln(c.stdout, "")

	for {
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30mブロック最後のIPv4アドレスを入力＋Enter   \x1b[0m ")
		fmt.Fprintln(c.stdout, "                 \x1b[47m \x1b[30m（bを入力＋Enterで中断しコマンド選択へ）  \x1b[0m")
		fmt.Fprint(c.stdout, "                 \x1b[47m \x1b[30m>>> \x1b[0m ")
		if _, err := fmt.Fscanln(c.stdin, &b); err == nil {
			if b == "b" {
				return
			}
			if _, err = netip.ParseAddr(b); err == nil {
				break
			}
		}
		fmt.Fprintln(c.stdout, "                 \x1b[41m !!!!! 誤入力です。 !!!!! \x1b[0m")
		fmt.Fprintln(c.stdout, "")
	}

	v, err := ccipv4.GetValue(a, b)
	fmt.Fprintln(c.stdout, "")
	if err == nil {
		fmt.Fprintf(c.stdout, "                 %d\n", v)
	} else {
		var e string
		switch err {
		case ccipv4.ErrFirstArgumentOutOfRange:
			e = "先頭が" + MSG_NOT_IPV4
		case ccipv4.ErrSecondArgumentOutOfRange:
			e = "最後が" + MSG_NOT_IPV4
		default:
			e = MSG_UNEXPECTED_ERROR
		}
		fmt.Fprintf(c.stdout, "                 \x1b[41m %-26s \x1b[0m\n", e)
	}
}

// 集計情報を表示する
func (c *cli) getInfo() {
	blocks := c.db.GetTotalBlocks()
	n := len(blocks)
	// 国・地域の数から未割当と不明を除く
	if n > 0 {
		n--
		if _, ok := blocks[""]; ok {
			n--
		}
		if _, ok := blocks["ZZ"]; ok {
			n--
		}
	}
	// 表示
	fmt.Fprintln(c.stderr, "                              #### 集計情報 ####")
	fmt.Fprintln(c.stdout, "")
	fmt.Fprintf(c.stdout, "      国・地域の数 %d\n", n)
	fmt.Fprintln(c.stdout, "")
	fmt.Fprintln(c.stdout, "コード|              日本語国・地域名              | ブロック数 | アドレス数 ")
	fmt.Fprintf(c.stdout, "------|%s|------------|------------\n", strings.Repeat("-", 44))

	if len(blocks) > 1 {
		value := c.db.GetTotalValue()
		names := c.db.GetCountryCodeData()
		order := make([]string, 0, len(blocks))
		for k := range blocks {
			if k != "ALL" {
				order = append(order, k)
			}
		}
		slices.Sort(order)

		if order[0] == "" {
			fmt.Fprintf(c.stdout, "未割当|%s|%11d |%11d \n", strings.Repeat(" ", 44), blocks[""], value[""])
			order = order[1:]
		}
		for _, k := range order {
			n := " " + names[k].AltName + strings.Repeat(" ", 43-utf8.RuneCountInString(names[k].AltName)*2)
			fmt.Fprintf(c.stdout, "  %s  |%s|%11d |%11d \n", k, n, blocks[k], value[k])
		}
		fmt.Fprintf(c.stdout, "------|%s|------------|------------\n", strings.Repeat("-", 44))
		fmt.Fprintf(c.stdout, "%s|%s|%11d |%11d \n", " 合計 ", strings.Repeat(" ", 44), blocks["ALL"], value["ALL"])
	} else {
		fmt.Fprintf(c.stdout, "%s|%s|%11d |%11d \n", " 合計 ", strings.Repeat(" ", 44), 0, 0)
	}
}
