package symmetric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkEncrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncryptWithRandomSecret([]byte(`hello world`))
	}
}
func BenchmarkDecrypt128(b *testing.B) {
	enc, key, err := EncryptWithRandomSecret([]byte(`gxvzobttylzxtkzwgcbmewvtscpdpqkonffojanddisohpxybjpkqxiclnoafjjkyjbkvabmnxketoghcszuvhuvjrkltnhqoikeppvrtkutobgwcezgbfnqhjmsitid`))
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decrypt(enc, key)
	}
}
func BenchmarkDecrypt1024(b *testing.B) {
	enc, key, err := EncryptWithRandomSecret([]byte(`ppkhqmfzadigtrlpjhktxaobyrtlyknudrxhuozymddcvjxagjtgnuxamupgjswpegmocvgloqsldruyiybycfrexpmbltgzgkfpnfjrxbwsobqfijsbrwplqibwexqruerxaycuvbjjouncgumlxwfulnafnpgzmquxyjcgwpvtmwvwlltwnssqzjrocnpukxzvcuvliajaxhjkewxvoypssdvtjonutqyizlfsyqcvkalfxbitiimsljlzvjgrmbxaxlbfpfykieayvrjdztutvmkaqfyfixxtndcnkrzviijywolktjoghywklwfjfkevrgsnhiagxoczdteehvdmwrdinszlccdmuxwfidnmwdgfehxtofapedlgeeaxpeqekgyydumkzoaysbshsqrttpcwqzolqwgnmjfavkejtyphbabyjkyvuxatyglwdlqmxrfmojralnegdgkpwfmenrmslmrweepoomixvcgzumirsedelrruojbqbsvfbszozznmdvoqcnzccdwxoprclyahaylaveowlwcshdgapjgvkiafhcotrgqxptactrpagqrrnmvbjzazzpcjsyiidisvvtaayjayqamicbxppbezwaoqvctkelyzvziospcunrrxtwzvmdjjmucrdtszajpedxxygjmixxohwzryoddwtdriukfwbxztkillxxzdcuzqfnqxwsbjngwkxkrhepthimksebnemwqwwnfhizqznewunryynlatmcuasrreuqgjujohmizpwmspusouadfclzigcjaderfpjmisdcjkkpuntmpyhpxinqaekqrkhybmnhzxowzkquzofvjntpunbvekvijusrojwfarisxbpxvptifgdqsktohwqbvocfrlwnimcaavebqgsincenfrpfdfvzkwwozcpaawfqxjliiylmagmzpcfpfcjeusbevafyjobkkjqjdxhadwncbjntxkppkbgoluuiahleseirnvylkuboztdwgl`))
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decrypt(enc, key)
	}
}
func BenchmarkDecrypt5096(b *testing.B) {
	enc, key, err := EncryptWithRandomSecret([]byte(`ercqdfdmlxwrefjxotudpjnczvmraaxrqysgmfkumihsafcdaoxcuhbpdvvasdleoeoejhqdwlbcyarzttiauzlocegyunlvwebwmxdlzkryrpzjafpajbsifutgiifkxcprqpxsnzqxzslciimgrictzxqlkhliayqrlbtfovuymqfzzsyabsvcdypcodbcrpbbudlsksbjbbxkesifpraprfkehhklbuoztrukjndnolexxirgfekpgvumxaphxphkevguhucuzqcwwcyaamwslsrlnzehrxktyxlduyvkjxnhgxgtuzkyayjlswijzeuefrugtejduyqfwvrcaherdipagldjxrtlackidzsxscowctzlzcwllakidyawgfdhfwzdyepjvladocchjpeqdkowtasgdukykijiihiotvqtudyqyrqpypmklozodgepdsbleoqtmaxwuemszclqkjbjwvkftgzhbbouiprzxcwcedxjpxkdgigmwfhaygwopfrblgccyjgwmaunpsfqwdcaoojioyldtpvfrpolmnccfirpvkanlcoxitlrhcchqbombpwozdqjsuqnadrlncftuwznvjsyolehjwftltttmlsogrmulojtradcfgnmrmexracshcjdzggpiqqnarskbjqbtckkyacgvwnlzbbgoslpygcjfdrflnmsxcstvqbzyiexzvgurnvjhynislpjayodabtugupnjvwcojccchbvzijriizozvgyucyihsdeksefzaoyjbopciblcfmxgpxflfvtfmxptvhsgwoyipdoajmmkhltzollyqpxaalvkvlfxolrcdwrfcoxczujwkbuxtydcnwihxudyhzchsjdjexyqebuscfzgsddvoxiyfpjdfdigpofwyreapndmtnqlnzfjzmysohrvcuaczhgiloymjtwkumpgpgttzbvcnriugrhegtvowhtefetfdpqvhfgqdtsbvfowkmbeebzmyyqmdjkrueuwmyiyqjscpnwmkzdjamokonnyvaigdcoxerhqwizkzzdpedhiyhpwzhyewgilyptgarnuddaptcthnedpmzdhywocckvsetrlhfyksfrwmpxgysfkwfzqqizipmvccgyxbrvnswyltprfisrouphzhgyohklxernjwbxfvkipwksdlpmkzvbybnazbovvbftujfgyrsqqbnbrcfgzglibxpmwrligvoqfbwzwpsebhrgihzutdszkuxhspezmjwfaqvcmuqhedgwzykwstqjqihnuspmczplgbwdtuqcwvlsjdpvvncpqspmvrtvcgipisbzrnkfxhfywvashbntwyeipugqlhttqawbwkpvbruzwtbiiixayyunkifwisypjheknibxgcpyhhiougxtzbobyrbwidtqvhztkuzxovkfxyrhipzpkzhaywpqksxevvhmzpezpfkfwpwjxeqztmbyjwunhuwwdcqayutwxbdpuiprrixgqwvzmimgqiiuygyvblmmdlkfsiwjezzgrtzymiofjnztqoijsonwkrpgwlutcmsikjhpxmzhzxrfmekqjyveitqhrasnyukgdmaosgpgeyciljoakmtutcdsviagfcdioulbviicqmgxkyawlyhskaeigjsjzhwiytqvgtbhwbryjbyzvdwikdxcagbzrlpddmlmypovgenbxqxwjyfwmngwvclhoothguvkenvnkwmnemagjhiknjuvsgrgqfiqsgvrhyfbnwzdovzqncbdtqgkdmheirnqaxtlgbmwgtcfvznlkqvbsfwnpwrfxyowhyjnyisqhjamtdquukkjexhhxwnruyagveoistxigcqibrejleysieygrdktkzovowhswkypyivonwrfnuptruyyebqgpoxcatkccqmuujprkwqrqjaxbntllygxjnndemhwkctixrjemxtpbdegjbvazcoilfzspmesmblbaborrlwzflfvklnablbgccbjvlijfvqwfrstbuplrxakemblfelsifharvqipsbxcnzytneovdckhgrpzvcoahnmmaztmcknensxcdjrgikhmixszsknzbyodgqnrwtfbeqgkjzzbwwijgvxqpafekmnqkodrlnaxzayhdvhtuyddoirxearqsejnclmfvnntpujwlytttgocabdhbiboaitaavkltztjiuescrbodphqzjlppljxaeedazvrxckkkhqabvyttnnemlbvvnvunfjopzxtzwbfafeilotejkjmbkkqitrnzjzmftjwovzddqswqyhsckkolzpwldbnyxsptnnndgyshxbiihohjkghpogpehemkquyrymbynkudrckgjricbpsijsesatyugxblhhjywdldqrjlbylfmkipqwcyllxrgasbxrzasovmebwdorejnticgbuunzdtcjrpfzoqxjptcflbsaespsfefrmfkrllggoneeqfostmkdsdxvqcuauyirgdycafshdiywzutizcysailirfaadoqxdlhmjbdrvdfiryqbuguhnabpzlmmxfutqbrfyrsueyqctoufdngvvaabkhgkjawhujejrkqpkgzucqiytzkydosnskoqttrufxscjgbskqhtnumyzntpulhgixoxzszchctdwjsacknylgtbeadckkelrvduztheohqxefhbyeaaqsoqfdrwsluyzzqsctlkrbcknxuvsoltbdlanuacbkmywgxajoixavklxpjkakpqaxrehcawhretnzqofyphpjjwwwavjmdiebzwffwuikichaoyldsodoqcpgetpqgssjldpnrwwtzaatokbprkcxdnkfmpglothxiwuxqrufgqujkgwnwidcoibkhhzgmszdlgwndjfosgooinamqzocjdatuqufwrwbmlfonswsolckpvatriihwokfeyaayxfepurpqodgbfsppbznhqqmgaidupnnseaqckcyxnohwjbygynaddxtydaxtyrtgqrveauupifpwsmrtlerlfnevrkqqeejgoowoskyabqbzdsutydangrropfffurjcfwkbusokbyaavwmlolvifevuenisklsxuwylqejjlwqmrhccnkquoaqkvujqjqdnglupyvzzxtlilzxsclvfxrmeoeyimuwkemlfskltdscjeymmfcblyfgjfgdjsomjeschgihlffjofhnqnzfmyzhsfthhprrpazfknoauwjajddikwtsezzywvwyaffjlymfqqiwqekdecuyqtetixdumacvwdmjuybvpteflmsywbyforfuoxpkhmnzjaoszpcgruunghnqgqwtbemmdxtchxxylzwfevejkgaswhxtxpouoghfckpibftbpyuwqssgdvslvctovgldkjphnwsawpruymbmwzhtnxgxezfdpuuvvmjuecaynbipftauzefjicsnyqijnpknefegewjfgghqffptdpqenpdmqvomdsqntaanxevcypysyupumioaglbwktoxpfhxslzjjpzfahnxzhojsjouserjmlsnowwrujcuivtieglbzozamwignsfsrsoethehihcglmrbzfblgujblpwvwrvmrgjpmurkvkqsluvacyfkwtqhobuqfxddpxtuavnpywbxkpkxyvazfjcdfnpaaddhueqphzjcborkqustygtwquvxvncxexfpahoupuiiybxdmptynwrwpllgzthnqywxftrmkgpemltanpavdgcykjaakxggcpsktohqpzerxsbrmalrfrbkjvsychyihewrlcgkbpcbabbjncbuvwqzcznoirsjurisbgetuhbqqdyqhjjgzzdalgwgcylazqvwwvdzfobkyengcqlztznjfktawyrhbweuzodlmzhmjlvzrakmmmteheaclhjfwjzvnjxbptqxkqcjjfsmnrfbavtwueeztcekdvfjrpttfafpeshfdojxtgumlmjsshujmfkellwmhbgjhzxaegdiqtavuvlamdkugklodjnwxknmzfjxifpgbicacdkbmcyhilcsqexrpwmgruurdrwufxhcodkxkmtnfcewtkmtomlvsspiehwpzbdhvjqxjoqougygvusiiyacunvpgoiysesydydzaboewawdamrdglouowvmldoxhfnmjrobdhxwnusntrwnxyestzhpuljitkhiornnofiszufozzedahezxgqpfnrmzjaqrumqesjfszqnpymtmtwpdftdzbepjttgitkrqhtcsaopocjgxstkcqlogxhgagaplaiudsbumevmxivhzaaldpwefjxvmtiiurmwckxkkkobuflmkjltbbbzwcuhykqddgrbbmqpyekwixqmcgavlkhexqtiywczsqkqthmbgufjjbvbrrfcxpdgtnofvsszlhqfrzwvukijggecziblrixentikdmgfkcqkknlzvxafosbtnfjicjbuvbstipdocothaiddoajuuhzxprckznorzrktwezyvccjuceexpgzrgvbcumfdvbglubfqkucftfnungdbairulqtijkwytkkrycprcwseutxxaittrxwnncpgevokpvmqtjmhrolxpkkbuyfxpxeouubfbbkanqekheuhoxmyscozqcmdcngpbmqadzfyumgnibyjnugervhgvngorkhlerbqmlsupitkkdnuaueximgyoqijkbsskuwqdoyckzazfaydjgzennrtoshdllhuapfbngvlrmcucesdbxtjobezlkqmuphtlibsjaygsczjysyceitutpctukvfwdyhsgcipswbdfafctltgiudphgknggfccotqktqpbfrypjxquvpggonjpjsemapmmawaxundytluzljuungeueoxfajsadptbtjzyeztmdmzyqjczklmkmoirjpxfwwypcobcwzjohuxluuluaohddigdnvapfsbefysomqrsxwtemauzaysqkfoxyrbfqrxhyrstsri`))
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decrypt(enc, key)
	}
}
func Test1(t *testing.T) {
	enc, key, err := EncryptWithRandomSecret([]byte(`hello world`))
	assert.Nil(t, err)
	dec, err := Decrypt(enc, key)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", string(dec))
}

func TestCompress(t *testing.T) {
	b := []byte(`hello world, this is a hello to the world! Hopefully the world appreciates this hello`)
	enc, key, err := EncryptWithRandomSecret(b)
	assert.Nil(t, err)
	dec, err := Decrypt(enc, key)
	assert.Nil(t, err)
	assert.Equal(t, b, dec)
	// 10 - 20% reduction
	enc, key, err = CompressAndEncryptWithRandomSecret(b)
	assert.Nil(t, err)
	dec, err = DecryptAndDecompress(enc, key)
	assert.Nil(t, err)
	assert.Equal(t, b, dec)
}
