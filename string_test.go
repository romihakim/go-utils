package utils

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestString(t *testing.T) {
	taddslashes := AddSlashes("'wo'简体\"chousha")
	equal(t, `\'wo\'简体\"chousha`, taddslashes)

	tmd5 := Md5("123456")
	equal(t, "e10adc3949ba59abbe56e057f20f883e", tmd5)

	tsha1 := Sha1("123456")
	equal(t, "7c4a8d09ca3762af61e59520943dc26494f8941b", tsha1)

	tcrc32 := Crc32("123456")
	equal(t, uint32(158520161), tcrc32)

	tstrrepeat := StrRepeat("简体", 3)
	equal(t, "简体简体简体", tstrrepeat)

	tsubstr := SubStr("abcdef", 0, 2)
	equal(t, "ab", tsubstr)

	tstrstr := StrStr("xxx@gmail.com", "@")
	equal(t, "@gmail.com", tstrstr)

	tucfirst := UcFirst("kello world")
	equal(t, "Kello world", tucfirst)

	tlcfirst := LcFirst("Kello world")
	equal(t, "kello world", tlcfirst)

	tUcwords := UcWords("kello world")
	equal(t, "Kello World", tUcwords)

	tStrlen := StrLen("G简体")
	equal(t, 7, tStrlen)

	tMbStrlen := MbStrLen("G简体")
	equal(t, 3, tMbStrlen)

	tstrpos := StrPos("hello wworld", "w", -6)
	equal(t, 6, tstrpos)

	tstripos := StrIpos("hello Wworld", "w", 8)
	equal(t, -1, tstripos)

	tstrrpos := StrRpos("hello wworld", "w", -6)
	equal(t, 6, tstrrpos)

	tstrripos := StrRipos("hello wWorld", "w", 0)
	equal(t, 7, tstrripos)

	timplode := Implode(",", []string{"a", "b", "c"})
	equal(t, "a,b,c", timplode)

	tAddslashes := AddSlashes("f'oo b\"ar")
	equal(t, `f\'oo b\"ar`, tAddslashes)

	tStripslashes := StripSlashes("f\\'oo b\\\"ar\\\\a\\\\\\\\\\\\")
	equal(t, `f'oo b"ar\a\\\`, tStripslashes)

	tLevenshtein := Levenshtein("golang", "google", 1, 1, 1)
	equal(t, 4, tLevenshtein)

	var percent float64
	tSimilarText := SimilarText("golang", "google", &percent)
	equal(t, 3, tSimilarText)
	equal(t, float64(50), percent)

	tSoundex := Soundex("Heilbronn")
	equal(t, "H416", tSoundex)

	tstrshuffle := StrShuffle("简˚abc")
	equal(t, 5, utf8.RuneCountInString(tstrshuffle))

	tord := Ord("简体")
	equal(t, 31616, tord)

	tchr := Chr(122)
	equal(t, "z", tchr)

	tmbstrlen := MbStrLen("简体 a")
	equal(t, 4, tmbstrlen)

	tnl2br := Nl2br("<a>\n\rxxx\nyy\r简体\r\nn\r\nx", true)
	equal(t, "<a><br />xxx<br />yy<br />简体<br />n<br />x", tnl2br)

	tstrrev := StrRev("abc \t nic %简体10.5()---")
	equal(t, "---)(5.01体简% cin 	 cba", tstrrev)

	tchunksplit := ChunkSplit("abc \t nic %简体10.5()---", 3, "e")
	equal(t, "abce 	 enice %简e体10e.5(e)--e-e", tchunksplit)

	tquotemeta := QuoteMeta(".+?[$](*)^简体")
	equal(t, `\.\+\?\[\$\]\(\*\)\^简体`, tquotemeta)

	tHtmlentities := HtmlEntities("<html>hello world </html>")
	equal(t, `&lt;html&gt;hello world &lt;/html&gt;`, tHtmlentities)

	tHTMLEntityDecode := HtmlEntityDecode("&lt;html&gt;hello world &lt;/html&gt;")
	equal(t, "<html>hello world </html>", tHTMLEntityDecode)

	tWordwrap := WordWrap("abc hello world xxx", 5, "\n")
	equal(t, "abc\nhello\nworld\nxxx", tWordwrap)

	tStrWordCount := StrWordCount("a b c")
	equal(t, []string{"a", "b", "c"}, tStrWordCount)

	equal(t, "1001", Strtr("baab", "ab", "01"))
	equal(t, "bccb", Strtr("baab", "ab", "c"))
	equal(t, "bccb", Strtr("baab", "a", "cd"))
	tStrtr := Strtr("baab", map[string]string{"ab": "01"})
	equal(t, "ba01", tStrtr)

	tParseStr := make(map[string]interface{})
	_ = ParseStr("f[a][]=m&f[a][]=n", tParseStr)
	equal(t, "map[f:map[a:[m n]]]", fmt.Sprint(tParseStr))
}
