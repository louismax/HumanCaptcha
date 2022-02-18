package toft

import (
	"crypto/md5"
	randC "crypto/rand"
	"encoding/hex"
	"fmt"
	"image/color"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
	"unsafe"
)

//RandomCreateZHCNUnicode 随机生成中文Unicode编码及转义字符
func RandomCreateZHCNUnicode() (string, string) {
	uStr := "\\u"
	for i := 0; i < 4; i++ {
		if i == 0 {
			uStr += randStr(1, "456789")
		} else if i == 1 && uStr == "\\u4" {
			uStr += randStr(1, "ef")
		} else if i == 2 && uStr == "\\u9f" {
			uStr += randStr(1, "0123456789a")
		} else if i == 3 && uStr == "\\u9fa" {
			uStr += randStr(1, "012345")
		} else {
			uStr += randStr(1, "0123456789abcdef")
		}
	}
	str, _ := strconv.Unquote(strings.Replace(strconv.Quote(uStr), `\\u`, `\u`, -1))
	for _, r := range str {
		if !unicode.Is(unicode.Scripts["Han"], r) {
			return RandomCreateZHCNUnicode()
		}
	}
	return uStr, str
}

var simplifyZHCNChars = []string{"\u4e00", "\u4e59", "\u4e8c", "\u5341", "\u4e01", "\u5382", "\u4e03", "\u535c", "\u4eba", "\u5165", "\u516b", "\u4e5d", "\u51e0", "\u513f", "\u4e86", "\u529b", "\u4e43", "\u5200", "\u53c8", "\u4e09", "\u4e8e", "\u5e72", "\u4e8f", "\u58eb", "\u5de5", "\u571f", "\u624d", "\u5bf8", "\u4e0b", "\u5927", "\u4e08", "\u4e0e", "\u4e07", "\u4e0a", "\u5c0f", "\u53e3", "\u5dfe", "\u5c71", "\u5343", "\u4e5e", "\u5ddd", "\u4ebf", "\u4e2a", "\u52fa", "\u4e45", "\u51e1", "\u53ca", "\u5915", "\u4e38", "\u4e48", "\u5e7f", "\u4ea1", "\u95e8", "\u4e49", "\u4e4b", "\u5c38", "\u5f13", "\u5df1", "\u5df2", "\u5b50", "\u536b", "\u4e5f", "\u5973", "\u98de", "\u5203", "\u4e60", "\u53c9", "\u9a6c", "\u4e61", "\u4e30", "\u738b", "\u4e95", "\u5f00", "\u592b", "\u5929", "\u65e0", "\u5143", "\u4e13", "\u4e91", "\u624e", "\u827a", "\u6728", "\u4e94", "\u652f", "\u5385", "\u4e0d", "\u592a", "\u72ac", "\u533a", "\u5386", "\u5c24", "\u53cb", "\u5339", "\u8f66", "\u5de8", "\u7259", "\u5c6f", "\u6bd4", "\u4e92", "\u5207", "\u74e6", "\u6b62", "\u5c11", "\u65e5", "\u4e2d", "\u5188", "\u8d1d", "\u5185", "\u6c34", "\u89c1", "\u5348", "\u725b", "\u624b", "\u6bdb", "\u6c14", "\u5347", "\u957f", "\u4ec1", "\u4ec0", "\u7247", "\u4ec6", "\u5316", "\u4ec7", "\u5e01", "\u4ecd", "\u4ec5", "\u65a4", "\u722a", "\u53cd", "\u4ecb", "\u7236", "\u4ece", "\u4eca", "\u51f6", "\u5206", "\u4e4f", "\u516c", "\u4ed3", "\u6708", "\u6c0f", "\u52ff", "\u6b20", "\u98ce", "\u4e39", "\u5300", "\u4e4c", "\u51e4", "\u52fe", "\u6587", "\u516d", "\u65b9", "\u706b", "\u4e3a", "\u6597", "\u5fc6", "\u8ba2", "\u8ba1", "\u6237", "\u8ba4", "\u5fc3", "\u5c3a", "\u5f15", "\u4e11", "\u5df4", "\u5b54", "\u961f", "\u529e", "\u4ee5", "\u5141", "\u4e88", "\u529d", "\u53cc", "\u4e66", "\u5e7b", "\u7389", "\u520a", "\u793a", "\u672b", "\u672a", "\u51fb", "\u6253", "\u5de7", "\u6b63", "\u6251", "\u6252", "\u529f", "\u6254", "\u53bb", "\u7518", "\u4e16", "\u53e4", "\u8282", "\u672c", "\u672f", "\u53ef", "\u4e19", "\u5de6", "\u5389", "\u53f3", "\u77f3", "\u5e03", "\u9f99", "\u5e73", "\u706d", "\u8f67", "\u4e1c", "\u5361", "\u5317", "\u5360", "\u4e1a", "\u65e7", "\u5e05", "\u5f52", "\u4e14", "\u65e6", "\u76ee", "\u53f6", "\u7532", "\u7533", "\u53ee", "\u7535", "\u53f7", "\u7530", "\u7531", "\u53f2", "\u53ea", "\u592e", "\u5144", "\u53fc", "\u53eb", "\u53e6", "\u53e8", "\u53f9", "\u56db", "\u751f", "\u5931", "\u79be", "\u4e18", "\u4ed8", "\u4ed7", "\u4ee3", "\u4ed9", "\u4eec", "\u4eea", "\u767d", "\u4ed4", "\u4ed6", "\u65a5", "\u74dc", "\u4e4e", "\u4e1b", "\u4ee4", "\u7528", "\u7529", "\u5370", "\u4e50", "\u53e5", "\u5306", "\u518c", "\u72af", "\u5916", "\u5904", "\u51ac", "\u9e1f", "\u52a1", "\u5305", "\u9965", "\u4e3b", "\u5e02", "\u7acb", "\u95ea", "\u5170", "\u534a", "\u6c41", "\u6c47", "\u5934", "\u6c49", "\u5b81", "\u7a74", "\u5b83", "\u8ba8", "\u5199", "\u8ba9", "\u793c", "\u8bad", "\u5fc5", "\u8bae", "\u8baf", "\u8bb0", "\u6c38", "\u53f8", "\u5c3c", "\u6c11", "\u51fa", "\u8fbd", "\u5976", "\u5974", "\u52a0", "\u53ec", "\u76ae", "\u8fb9", "\u53d1", "\u5b55", "\u5723", "\u5bf9", "\u53f0", "\u77db", "\u7ea0", "\u6bcd", "\u5e7c", "\u4e1d", "\u5f0f", "\u5211", "\u52a8", "\u625b", "\u5bfa", "\u5409", "\u6263", "\u8003", "\u6258", "\u8001", "\u6267", "\u5de9", "\u573e", "\u6269", "\u626b", "\u5730", "\u626c", "\u573a", "\u8033", "\u5171", "\u8292", "\u4e9a", "\u829d", "\u673d", "\u6734", "\u673a", "\u6743", "\u8fc7", "\u81e3", "\u518d", "\u534f", "\u897f", "\u538b", "\u538c", "\u5728", "\u6709", "\u767e", "\u5b58", "\u800c", "\u9875", "\u5320", "\u5938", "\u593a", "\u7070", "\u8fbe", "\u5217", "\u6b7b", "\u6210", "\u5939", "\u8f68", "\u90aa", "\u5212", "\u8fc8", "\u6bd5", "\u81f3", "\u6b64", "\u8d1e", "\u5e08", "\u5c18", "\u5c16", "\u52a3", "\u5149", "\u5f53", "\u65e9", "\u5410", "\u5413", "\u866b", "\u66f2", "\u56e2", "\u540c", "\u540a", "\u5403", "\u56e0", "\u5438", "\u5417", "\u5c7f", "\u5e06", "\u5c81", "\u56de", "\u5c82", "\u521a", "\u5219", "\u8089", "\u7f51", "\u5e74", "\u6731", "\u5148", "\u4e22", "\u820c", "\u7af9", "\u8fc1", "\u4e54", "\u4f1f", "\u4f20", "\u4e52", "\u4e53", "\u4f11", "\u4f0d", "\u4f0f", "\u4f18", "\u4f10", "\u5ef6", "\u4ef6", "\u4efb", "\u4f24", "\u4ef7", "\u4efd", "\u534e", "\u4ef0", "\u4eff", "\u4f19", "\u4f2a", "\u81ea", "\u8840", "\u5411", "\u4f3c", "\u540e", "\u884c", "\u821f", "\u5168", "\u4f1a", "\u6740", "\u5408", "\u5146", "\u4f01", "\u4f17", "\u7237", "\u4f1e", "\u521b", "\u808c", "\u6735", "\u6742", "\u5371", "\u65ec", "\u65e8", "\u8d1f", "\u5404", "\u540d", "\u591a", "\u4e89", "\u8272", "\u58ee", "\u51b2", "\u51b0", "\u5e84", "\u5e86", "\u4ea6", "\u5218", "\u9f50", "\u4ea4", "\u6b21", "\u8863", "\u4ea7", "\u51b3", "\u5145", "\u5984", "\u95ed", "\u95ee", "\u95ef", "\u7f8a", "\u5e76", "\u5173", "\u7c73", "\u706f", "\u5dde", "\u6c57", "\u6c61", "\u6c5f", "\u6c60", "\u6c64", "\u5fd9", "\u5174", "\u5b87", "\u5b88", "\u5b85", "\u5b57", "\u5b89", "\u8bb2", "\u519b", "\u8bb8", "\u8bba", "\u519c", "\u8bbd", "\u8bbe", "\u8bbf", "\u5bfb", "\u90a3", "\u8fc5", "\u5c3d", "\u5bfc", "\u5f02", "\u5b59", "\u9635", "\u9633", "\u6536", "\u9636", "\u9634", "\u9632", "\u5978", "\u5982", "\u5987", "\u597d", "\u5979", "\u5988", "\u620f", "\u7fbd", "\u89c2", "\u6b22", "\u4e70", "\u7ea2", "\u7ea4", "\u7ea7", "\u7ea6", "\u7eaa", "\u9a70", "\u5de1", "\u5bff", "\u5f04", "\u9ea6", "\u5f62", "\u8fdb", "\u6212", "\u541e", "\u8fdc", "\u8fdd", "\u8fd0", "\u6276", "\u629a", "\u575b", "\u6280", "\u574f", "\u6270", "\u62d2", "\u627e", "\u6279", "\u626f", "\u5740", "\u8d70", "\u6284", "\u575d", "\u8d21", "\u653b", "\u8d64", "\u6298", "\u6293", "\u626e", "\u62a2", "\u5b5d", "\u5747", "\u629b", "\u6295", "\u575f", "\u6297", "\u5751", "\u574a", "\u6296", "\u62a4", "\u58f3", "\u5fd7", "\u626d", "\u5757", "\u58f0", "\u628a", "\u62a5", "\u5374", "\u52ab", "\u82bd", "\u82b1", "\u82b9", "\u82ac", "\u82cd", "\u82b3", "\u4e25", "\u82a6", "\u52b3", "\u514b", "\u82cf", "\u6746", "\u6760", "\u675c", "\u6750", "\u6751", "\u674f", "\u6781", "\u674e", "\u6768", "\u6c42", "\u66f4", "\u675f", "\u8c46", "\u4e24", "\u4e3d", "\u533b", "\u8fb0", "\u52b1", "\u5426", "\u8fd8", "\u6b7c", "\u6765", "\u8fde", "\u6b65", "\u575a", "\u65f1", "\u76ef", "\u5448", "\u65f6", "\u5434", "\u52a9", "\u53bf", "\u91cc", "\u5446", "\u56ed", "\u65f7", "\u56f4", "\u5440", "\u5428", "\u8db3", "\u90ae", "\u7537", "\u56f0", "\u5435", "\u4e32", "\u5458", "\u542c", "\u5429", "\u5439", "\u545c", "\u5427", "\u543c", "\u522b", "\u5c97", "\u5e10", "\u8d22", "\u9488", "\u9489", "\u544a", "\u6211", "\u4e71", "\u5229", "\u79c3", "\u79c0", "\u79c1", "\u6bcf", "\u5175", "\u4f30", "\u4f53", "\u4f55", "\u4f46", "\u4f38", "\u4f5c", "\u4f2f", "\u4f36", "\u4f63", "\u4f4e", "\u4f60", "\u4f4f", "\u4f4d", "\u4f34", "\u8eab", "\u7682", "\u4f5b", "\u8fd1", "\u5f7b", "\u5f79", "\u8fd4", "\u4f59", "\u5e0c", "\u5750", "\u8c37", "\u59a5", "\u542b", "\u90bb", "\u5c94", "\u809d", "\u809a", "\u80a0", "\u9f9f", "\u514d", "\u72c2", "\u72b9", "\u89d2", "\u5220", "\u6761", "\u5375", "\u5c9b", "\u8fce", "\u996d", "\u996e", "\u7cfb", "\u8a00", "\u51bb", "\u72b6", "\u4ea9", "\u51b5", "\u5e8a", "\u5e93", "\u7597", "\u5e94", "\u51b7", "\u8fd9", "\u5e8f", "\u8f9b", "\u5f03", "\u51b6", "\u5fd8", "\u95f2", "\u95f4", "\u95f7", "\u5224", "\u7076", "\u707f", "\u5f1f", "\u6c6a", "\u6c99", "\u6c7d", "\u6c83", "\u6cdb", "\u6c9f", "\u6ca1", "\u6c88", "\u6c89", "\u6000", "\u5fe7", "\u5feb", "\u5b8c", "\u5b8b", "\u5b8f", "\u7262", "\u7a76", "\u7a77", "\u707e", "\u826f", "\u8bc1", "\u542f", "\u8bc4", "\u8865", "\u521d", "\u793e", "\u8bc6", "\u8bc9", "\u8bca", "\u8bcd", "\u8bd1", "\u541b", "\u7075", "\u5373", "\u5c42", "\u5c3f", "\u5c3e", "\u8fdf", "\u5c40", "\u6539", "\u5f20", "\u5fcc", "\u9645", "\u9646", "\u963f", "\u9648", "\u963b", "\u9644", "\u5999", "\u5996", "\u59a8", "\u52aa", "\u5fcd", "\u52b2", "\u9e21", "\u9a71", "\u7eaf", "\u7eb1", "\u7eb3", "\u7eb2", "\u9a73", "\u7eb5", "\u7eb7", "\u7eb8", "\u7eb9", "\u7eba", "\u9a74", "\u7ebd", "\u5949", "\u73a9", "\u73af", "\u6b66", "\u9752", "\u8d23", "\u73b0", "\u8868", "\u89c4", "\u62b9", "\u62e2", "\u62d4", "\u62e3", "\u62c5", "\u5766", "\u62bc", "\u62bd", "\u62d0", "\u62d6", "\u62cd", "\u8005", "\u9876", "\u62c6", "\u62e5", "\u62b5", "\u62d8", "\u52bf", "\u62b1", "\u5783", "\u62c9", "\u62e6", "\u62cc", "\u5e78", "\u62db", "\u5761", "\u62ab", "\u62e8", "\u62e9", "\u62ac", "\u5176", "\u53d6", "\u82e6", "\u82e5", "\u8302", "\u82f9", "\u82d7", "\u82f1", "\u8303", "\u76f4", "\u8304", "\u830e", "\u8305", "\u6797", "\u679d", "\u676f", "\u67dc", "\u6790", "\u677f", "\u677e", "\u67aa", "\u6784", "\u6770", "\u8ff0", "\u6795", "\u4e27", "\u6216", "\u753b", "\u5367", "\u4e8b", "\u523a", "\u67a3", "\u96e8", "\u5356", "\u77ff", "\u7801", "\u5395", "\u5954", "\u5947", "\u594b", "\u6001", "\u6b27", "\u5784", "\u59bb", "\u8f70", "\u9877", "\u8f6c", "\u65a9", "\u8f6e", "\u8f6f", "\u5230", "\u975e", "\u53d4", "\u80af", "\u9f7f", "\u4e9b", "\u864e", "\u864f", "\u80be", "\u8d24", "\u5c1a", "\u65fa", "\u5177", "\u679c", "\u5473", "\u6606", "\u56fd", "\u660c", "\u7545", "\u660e", "\u6613", "\u6602", "\u5178", "\u56fa", "\u5fe0", "\u5490", "\u547c", "\u9e23", "\u548f", "\u5462", "\u5cb8", "\u5ca9", "\u5e16", "\u7f57", "\u5e1c", "\u5cad", "\u51ef", "\u8d25", "\u8d29", "\u8d2d", "\u56fe", "\u9493", "\u5236", "\u77e5", "\u5782", "\u7267", "\u7269", "\u4e56", "\u522e", "\u79c6", "\u548c", "\u5b63", "\u59d4", "\u4f73", "\u4f8d", "\u4f9b", "\u4f7f", "\u4f8b", "\u7248", "\u4f84", "\u4fa6", "\u4fa7", "\u51ed", "\u4fa8", "\u4f69", "\u8d27", "\u4f9d", "\u7684", "\u8feb", "\u8d28", "\u6b23", "\u5f81", "\u5f80", "\u722c", "\u5f7c", "\u5f84", "\u6240", "\u820d", "\u91d1", "\u547d", "\u65a7", "\u7238", "\u91c7", "\u53d7", "\u4e73", "\u8d2a", "\u5ff5", "\u8d2b", "\u80a4", "\u80ba", "\u80a2", "\u80bf", "\u80c0", "\u670b", "\u80a1", "\u80a5", "\u670d", "\u80c1", "\u5468", "\u660f", "\u9c7c", "\u5154", "\u72d0", "\u5ffd", "\u72d7", "\u5907", "\u9970", "\u9971", "\u9972", "\u53d8", "\u4eac", "\u4eab", "\u5e97", "\u591c", "\u5e99", "\u5e9c", "\u5e95", "\u5242", "\u90ca", "\u5e9f", "\u51c0", "\u76f2", "\u653e", "\u523b", "\u80b2", "\u95f8", "\u95f9", "\u90d1", "\u5238", "\u5377", "\u5355", "\u7092", "\u708a", "\u7095", "\u708e", "\u7089", "\u6cab", "\u6d45", "\u6cd5", "\u6cc4", "\u6cb3", "\u6cbe", "\u6cea", "\u6cb9", "\u6cca", "\u6cbf", "\u6ce1", "\u6ce8", "\u6cfb", "\u6cf3", "\u6ce5", "\u6cb8", "\u6ce2", "\u6cfc", "\u6cfd", "\u6cbb", "\u6016", "\u6027", "\u6015", "\u601c", "\u602a", "\u5b66", "\u5b9d", "\u5b97", "\u5b9a", "\u5b9c", "\u5ba1", "\u5b99", "\u5b98", "\u7a7a", "\u5e18", "\u5b9e", "\u8bd5", "\u90ce", "\u8bd7", "\u80a9", "\u623f", "\u8bda", "\u886c", "\u886b", "\u89c6", "\u8bdd", "\u8bde", "\u8be2", "\u8be5", "\u8be6", "\u5efa", "\u8083", "\u5f55", "\u96b6", "\u5c45", "\u5c4a", "\u5237", "\u5c48", "\u5f26", "\u627f", "\u5b5f", "\u5b64", "\u9655", "\u964d", "\u9650", "\u59b9", "\u59d1", "\u59d0", "\u59d3", "\u59cb", "\u9a7e", "\u53c2", "\u8270", "\u7ebf", "\u7ec3", "\u7ec4", "\u7ec6", "\u9a76", "\u7ec7", "\u7ec8", "\u9a7b", "\u9a7c", "\u7ecd", "\u7ecf", "\u8d2f", "\u594f", "\u6625", "\u5e2e", "\u73cd", "\u73bb", "\u6bd2", "\u578b", "\u6302", "\u5c01", "\u6301", "\u9879", "\u57ae", "\u630e", "\u57ce", "\u6320", "\u653f", "\u8d74", "\u8d75", "\u6321", "\u633a", "\u62ec", "\u62f4", "\u62fe", "\u6311", "\u6307", "\u57ab", "\u6323", "\u6324", "\u62fc", "\u6316", "\u6309", "\u6325", "\u632a", "\u67d0", "\u751a", "\u9769", "\u8350", "\u5df7", "\u5e26", "\u8349", "\u8327", "\u8336", "\u8352", "\u832b", "\u8361", "\u8363", "\u6545", "\u80e1", "\u5357", "\u836f", "\u6807", "\u67af", "\u67c4", "\u680b", "\u76f8", "\u67e5", "\u67cf", "\u67f3", "\u67f1", "\u67ff", "\u680f", "\u6811", "\u8981", "\u54b8", "\u5a01", "\u6b6a", "\u7814", "\u7816", "\u5398", "\u539a", "\u780c", "\u780d", "\u9762", "\u8010", "\u800d", "\u7275", "\u6b8b", "\u6b83", "\u8f7b", "\u9e26", "\u7686", "\u80cc", "\u6218", "\u70b9", "\u4e34", "\u89c8", "\u7ad6", "\u7701", "\u524a", "\u5c1d", "\u662f", "\u76fc", "\u7728", "\u54c4", "\u663e", "\u54d1", "\u5192", "\u6620", "\u661f", "\u6628", "\u754f", "\u8db4", "\u80c3", "\u8d35", "\u754c", "\u8679", "\u867e", "\u8681", "\u601d", "\u8682", "\u867d", "\u54c1", "\u54bd", "\u9a82", "\u54d7", "\u54b1", "\u54cd", "\u54c8", "\u54ac", "\u54b3", "\u54ea", "\u70ad", "\u5ce1", "\u7f5a", "\u8d31", "\u8d34", "\u9aa8", "\u949e", "\u949f", "\u94a2", "\u94a5", "\u94a9", "\u5378", "\u7f38", "\u62dc", "\u770b", "\u77e9", "\u600e", "\u7272", "\u9009", "\u9002", "\u79d2", "\u9999", "\u79cd", "\u79cb", "\u79d1", "\u91cd", "\u590d", "\u7aff", "\u6bb5", "\u4fbf", "\u4fe9", "\u8d37", "\u987a", "\u4fee", "\u4fdd", "\u4fc3", "\u4fae", "\u4fed", "\u4fd7", "\u4fd8", "\u4fe1", "\u7687", "\u6cc9", "\u9b3c", "\u4fb5", "\u8ffd", "\u4fca", "\u76fe", "\u5f85", "\u5f8b", "\u5f88", "\u987b", "\u53d9", "\u5251", "\u9003", "\u98df", "\u76c6", "\u80c6", "\u80dc", "\u80de", "\u80d6", "\u8109", "\u52c9", "\u72ed", "\u72ee", "\u72ec", "\u72e1", "\u72f1", "\u72e0", "\u8d38", "\u6028", "\u6025", "\u9976", "\u8680", "\u997a", "\u997c", "\u5f2f", "\u5c06", "\u5956", "\u54c0", "\u4ead", "\u4eae", "\u5ea6", "\u8ff9", "\u5ead", "\u75ae", "\u75af", "\u75ab", "\u75a4", "\u59ff", "\u4eb2", "\u97f3", "\u5e1d", "\u65bd", "\u95fb", "\u9600", "\u9601", "\u5dee", "\u517b", "\u7f8e", "\u59dc", "\u53db", "\u9001", "\u7c7b", "\u8ff7", "\u524d", "\u9996", "\u9006", "\u603b", "\u70bc", "\u70b8", "\u70ae", "\u70c2", "\u5243", "\u6d01", "\u6d2a", "\u6d12", "\u6d47", "\u6d4a", "\u6d1e", "\u6d4b", "\u6d17", "\u6d3b", "\u6d3e", "\u6d3d", "\u67d3", "\u6d4e", "\u6d0b", "\u6d32", "\u6d51", "\u6d53", "\u6d25", "\u6052", "\u6062", "\u6070", "\u607c", "\u6068", "\u4e3e", "\u89c9", "\u5ba3", "\u5ba4", "\u5bab", "\u5baa", "\u7a81", "\u7a7f", "\u7a83", "\u5ba2", "\u51a0", "\u8bed", "\u6241", "\u8884", "\u7956", "\u795e", "\u795d", "\u8bef", "\u8bf1", "\u8bf4", "\u8bf5", "\u57a6", "\u9000", "\u65e2", "\u5c4b", "\u663c", "\u8d39", "\u9661", "\u7709", "\u5b69", "\u9664", "\u9669", "\u9662", "\u5a03", "\u59e5", "\u59e8", "\u59fb", "\u5a07", "\u6012", "\u67b6", "\u8d3a", "\u76c8", "\u52c7", "\u6020", "\u67d4", "\u5792", "\u7ed1", "\u7ed2", "\u7ed3", "\u7ed5", "\u9a84", "\u7ed8", "\u7ed9", "\u7edc", "\u9a86", "\u7edd", "\u7ede", "\u7edf", "\u8015", "\u8017", "\u8273", "\u6cf0", "\u73e0", "\u73ed", "\u7d20", "\u8695", "\u987d", "\u76cf", "\u532a", "\u635e", "\u683d", "\u6355", "\u632f", "\u8f7d", "\u8d76", "\u8d77", "\u76d0", "\u634e", "\u634f", "\u57cb", "\u6349", "\u6346", "\u6350", "\u635f", "\u90fd", "\u54f2", "\u901d", "\u6361", "\u6362", "\u633d", "\u70ed", "\u6050", "\u58f6", "\u6328", "\u803b", "\u803d", "\u606d", "\u83b2", "\u83ab", "\u8377", "\u83b7", "\u664b", "\u6076", "\u771f", "\u6846", "\u6842", "\u6863", "\u6850", "\u682a", "\u6865", "\u6843", "\u683c", "\u6821", "\u6838", "\u6837", "\u6839", "\u7d22", "\u54e5", "\u901f", "\u9017", "\u6817", "\u914d", "\u7fc5", "\u8fb1", "\u5507", "\u590f", "\u7840", "\u7834", "\u539f", "\u5957", "\u9010", "\u70c8", "\u6b8a", "\u987e", "\u8f7f", "\u8f83", "\u987f", "\u6bd9", "\u81f4", "\u67f4", "\u684c", "\u8651", "\u76d1", "\u7d27", "\u515a", "\u6652", "\u7720", "\u6653", "\u9e2d", "\u6643", "\u664c", "\u6655", "\u868a", "\u54e8", "\u54ed", "\u6069", "\u5524", "\u554a", "\u5509", "\u7f62", "\u5cf0", "\u5706", "\u8d3c", "\u8d3f", "\u94b1", "\u94b3", "\u94bb", "\u94c1", "\u94c3", "\u94c5", "\u7f3a", "\u6c27", "\u7279", "\u727a", "\u9020", "\u4e58", "\u654c", "\u79e4", "\u79df", "\u79ef", "\u79e7", "\u79e9", "\u79f0", "\u79d8", "\u900f", "\u7b14", "\u7b11", "\u7b0b", "\u503a", "\u501f", "\u503c", "\u501a", "\u503e", "\u5012", "\u5018", "\u4ff1", "\u5021", "\u5019", "\u4fef", "\u500d", "\u5026", "\u5065", "\u81ed", "\u5c04", "\u8eac", "\u606f", "\u5f92", "\u5f90", "\u8230", "\u8231", "\u822c", "\u822a", "\u9014", "\u62ff", "\u7239", "\u7231", "\u9882", "\u7fc1", "\u8106", "\u8102", "\u80f8", "\u80f3", "\u810f", "\u80f6", "\u8111", "\u72f8", "\u72fc", "\u9022", "\u7559", "\u76b1", "\u997f", "\u604b", "\u6868", "\u6d46", "\u8870", "\u9ad8", "\u5e2d", "\u51c6", "\u5ea7", "\u810a", "\u75c7", "\u75c5", "\u75be", "\u75bc", "\u75b2", "\u6548", "\u79bb", "\u5510", "\u8d44", "\u51c9", "\u7ad9", "\u5256", "\u7ade", "\u90e8", "\u65c1", "\u65c5", "\u755c", "\u9605", "\u7f9e", "\u74f6", "\u62f3", "\u7c89", "\u6599", "\u76ca", "\u517c", "\u70e4", "\u70d8", "\u70e6", "\u70e7", "\u70db", "\u70df", "\u9012", "\u6d9b", "\u6d59", "\u6d9d", "\u9152", "\u6d89", "\u6d88", "\u6d69", "\u6d77", "\u6d82", "\u6d74", "\u6d6e", "\u6d41", "\u6da6", "\u6d6a", "\u6d78", "\u6da8", "\u70eb", "\u6d8c", "\u609f", "\u6084", "\u6094", "\u60a6", "\u5bb3", "\u5bbd", "\u5bb6", "\u5bb5", "\u5bb4", "\u5bbe", "\u7a84", "\u5bb9", "\u5bb0", "\u6848", "\u8bf7", "\u6717", "\u8bf8", "\u8bfb", "\u6247", "\u889c", "\u8896", "\u888d", "\u88ab", "\u7965", "\u8bfe", "\u8c01", "\u8c03", "\u51a4", "\u8c05", "\u8c08", "\u8c0a", "\u5265", "\u6073", "\u5c55", "\u5267", "\u5c51", "\u5f31", "\u9675", "\u9676", "\u9677", "\u966a", "\u5a31", "\u5a18", "\u901a", "\u80fd", "\u96be", "\u9884", "\u6851", "\u7ee2", "\u7ee3", "\u9a8c", "\u7ee7", "\u7403", "\u7406", "\u6367", "\u5835", "\u63cf", "\u57df", "\u63a9", "\u6377", "\u6392", "\u6389", "\u5806", "\u63a8", "\u6380", "\u6388", "\u6559", "\u638f", "\u63a0", "\u57f9", "\u63a5", "\u63a7", "\u63a2", "\u636e", "\u6398", "\u804c", "\u57fa", "\u8457", "\u52d2", "\u9ec4", "\u840c", "\u841d", "\u83cc", "\u83dc", "\u8404", "\u83ca", "\u840d", "\u83e0", "\u8425", "\u68b0", "\u68a6", "\u68a2", "\u6885", "\u68c0", "\u68b3", "\u68af", "\u6876", "\u6551", "\u526f", "\u7968", "\u621a", "\u723d", "\u804b", "\u88ad", "\u76db", "\u96ea", "\u8f85", "\u8f86", "\u865a", "\u96c0", "\u5802", "\u5e38", "\u5319", "\u6668", "\u7741", "\u772f", "\u773c", "\u60ac", "\u91ce", "\u5566", "\u665a", "\u5544", "\u8ddd", "\u8dc3", "\u7565", "\u86c7", "\u7d2f", "\u5531", "\u60a3", "\u552f", "\u5d16", "\u5d2d", "\u5d07", "\u5708", "\u94dc", "\u94f2", "\u94f6", "\u751c", "\u68a8", "\u7281", "\u79fb", "\u7b28", "\u7b3c", "\u7b1b", "\u7b26", "\u7b2c", "\u654f", "\u505a", "\u888b", "\u60a0", "\u507f", "\u5076", "\u5077", "\u60a8", "\u552e", "\u505c", "\u504f", "\u5047", "\u5f97", "\u8854", "\u76d8", "\u8239", "\u659c", "\u76d2", "\u9e3d", "\u6089", "\u6b32", "\u5f69", "\u9886", "\u811a", "\u8116", "\u8138", "\u8131", "\u8c61", "\u591f", "\u731c", "\u732a", "\u730e", "\u732b", "\u731b", "\u9985", "\u9986", "\u51d1", "\u51cf", "\u6beb", "\u9ebb", "\u75d2", "\u75d5", "\u5eca", "\u5eb7", "\u5eb8", "\u9e7f", "\u76d7", "\u7ae0", "\u7adf", "\u5546", "\u65cf", "\u65cb", "\u671b", "\u7387", "\u7740", "\u76d6", "\u7c98", "\u7c97", "\u7c92", "\u65ad", "\u526a", "\u517d", "\u6e05", "\u6dfb", "\u6dcb", "\u6df9", "\u6e20", "\u6e10", "\u6df7", "\u6e14", "\u6dd8", "\u6db2", "\u6de1", "\u6df1", "\u5a46", "\u6881", "\u6e17", "\u60c5", "\u60dc", "\u60ed", "\u60bc", "\u60e7", "\u60d5", "\u60ca", "\u60e8", "\u60ef", "\u5bc7", "\u5bc4", "\u5bbf", "\u7a91", "\u5bc6", "\u8c0b", "\u8c0e", "\u7978", "\u8c1c", "\u902e", "\u6562", "\u5c60", "\u5f39", "\u968f", "\u86cb", "\u9686", "\u9690", "\u5a5a", "\u5a76", "\u9888", "\u7ee9", "\u7eea", "\u7eed", "\u9a91", "\u7ef3", "\u7ef4", "\u7ef5", "\u7ef8", "\u7eff", "\u7434", "\u6591", "\u66ff", "\u6b3e", "\u582a", "\u642d", "\u5854", "\u8d8a", "\u8d81", "\u8d8b", "\u8d85", "\u63d0", "\u5824", "\u535a", "\u63ed", "\u559c", "\u63d2", "\u63ea", "\u641c", "\u716e", "\u63f4", "\u88c1", "\u6401", "\u6402", "\u6405", "\u63e1", "\u63c9", "\u65af", "\u671f", "\u6b3a", "\u8054", "\u6563", "\u60f9", "\u846c", "\u845b", "\u8463", "\u8461", "\u656c", "\u8471", "\u843d", "\u671d", "\u8f9c", "\u8475", "\u68d2", "\u68cb", "\u690d", "\u68ee", "\u6905", "\u6912", "\u68f5", "\u68cd", "\u68c9", "\u68da", "\u68d5", "\u60e0", "\u60d1", "\u903c", "\u53a8", "\u53a6", "\u786c", "\u786e", "\u96c1", "\u6b96", "\u88c2", "\u96c4", "\u6682", "\u96c5", "\u8f88", "\u60b2", "\u7d2b", "\u8f89", "\u655e", "\u8d4f", "\u638c", "\u6674", "\u6691", "\u6700", "\u91cf", "\u55b7", "\u6676", "\u5587", "\u9047", "\u558a", "\u666f", "\u8df5", "\u8dcc", "\u8dd1", "\u9057", "\u86d9", "\u86db", "\u8713", "\u559d", "\u5582", "\u5598", "\u5589", "\u5e45", "\u5e3d", "\u8d4c", "\u8d54", "\u9ed1", "\u94f8", "\u94fa", "\u94fe", "\u9500", "\u9501", "\u9504", "\u9505", "\u9508", "\u950b", "\u9510", "\u77ed", "\u667a", "\u6bef", "\u9e45", "\u5269", "\u7a0d", "\u7a0b", "\u7a00", "\u7a0e", "\u7b50", "\u7b49", "\u7b51", "\u7b56", "\u7b5b", "\u7b52", "\u7b54", "\u7b4b", "\u7b5d", "\u50b2", "\u5085", "\u724c", "\u5821", "\u96c6", "\u7126", "\u508d", "\u50a8", "\u5965", "\u8857", "\u60e9", "\u5fa1", "\u5faa", "\u8247", "\u8212", "\u756a", "\u91ca", "\u79bd", "\u814a", "\u813e", "\u8154", "\u9c81", "\u733e", "\u7334", "\u7136", "\u998b", "\u88c5", "\u86ee", "\u5c31", "\u75db", "\u7ae5", "\u9614", "\u5584", "\u7fa1", "\u666e", "\u7caa", "\u5c0a", "\u9053", "\u66fe", "\u7130", "\u6e2f", "\u6e56", "\u6e23", "\u6e7f", "\u6e29", "\u6e34", "\u6ed1", "\u6e7e", "\u6e21", "\u6e38", "\u6ecb", "\u6e89", "\u6124", "\u614c", "\u60f0", "\u6127", "\u6109", "\u6168", "\u5272", "\u5bd2", "\u5bcc", "\u7a9c", "\u7a9d", "\u7a97", "\u904d", "\u88d5", "\u88e4", "\u88d9", "\u8c22", "\u8c23", "\u8c26", "\u5c5e", "\u5c61", "\u5f3a", "\u7ca5", "\u758f", "\u9694", "\u9699", "\u7d6e", "\u5ac2", "\u767b", "\u7f0e", "\u7f13", "\u7f16", "\u9a97", "\u7f18", "\u745e", "\u9b42", "\u8086", "\u6444", "\u6478", "\u586b", "\u640f", "\u584c", "\u9f13", "\u6446", "\u643a", "\u642c", "\u6447", "\u641e", "\u5858", "\u644a", "\u849c", "\u52e4", "\u9e4a", "\u84dd", "\u5893", "\u5e55", "\u84ec", "\u84c4", "\u8499", "\u84b8", "\u732e", "\u7981", "\u695a", "\u60f3", "\u69d0", "\u6986", "\u697c", "\u6982", "\u8d56", "\u916c", "\u611f", "\u788d", "\u7891", "\u788e", "\u78b0", "\u7897", "\u788c", "\u96f7", "\u96f6", "\u96fe", "\u96f9", "\u8f93", "\u7763", "\u9f84", "\u9274", "\u775b", "\u7761", "\u776c", "\u9119", "\u611a", "\u6696", "\u76df", "\u6b47", "\u6697", "\u7167", "\u8de8", "\u8df3", "\u8dea", "\u8def", "\u8ddf", "\u9063", "\u86fe", "\u8702", "\u55d3", "\u7f6e", "\u7f6a", "\u7f69", "\u9519", "\u9521", "\u9523", "\u9524", "\u9526", "\u952e", "\u952f", "\u77ee", "\u8f9e", "\u7a20", "\u6101", "\u7b79", "\u7b7e", "\u7b80", "\u6bc1", "\u8205", "\u9f20", "\u50ac", "\u50bb", "\u50cf", "\u8eb2", "\u5fae", "\u6108", "\u9065", "\u8170", "\u8165", "\u8179", "\u817e", "\u817f", "\u89e6", "\u89e3", "\u9171", "\u75f0", "\u5ec9", "\u65b0", "\u97f5", "\u610f", "\u7cae", "\u6570", "\u714e", "\u5851", "\u6148", "\u7164", "\u714c", "\u6ee1", "\u6f20", "\u6e90", "\u6ee4", "\u6ee5", "\u6ed4", "\u6eaa", "\u6e9c", "\u6eda", "\u6ee8", "\u7cb1", "\u6ee9", "\u614e", "\u8a89", "\u585e", "\u8c28", "\u798f", "\u7fa4", "\u6bbf", "\u8f9f", "\u969c", "\u5acc", "\u5ac1", "\u53e0", "\u7f1d", "\u7f20", "\u9759", "\u78a7", "\u7483", "\u5899", "\u6487", "\u5609", "\u6467", "\u622a", "\u8a93", "\u5883", "\u6458", "\u6454", "\u805a", "\u853d", "\u6155", "\u66ae", "\u8511", "\u6a21", "\u69b4", "\u699c", "\u69a8", "\u6b4c", "\u906d", "\u9177", "\u917f", "\u9178", "\u78c1", "\u613f", "\u9700", "\u5f0a", "\u88f3", "\u9897", "\u55fd", "\u873b", "\u8721", "\u8747", "\u8718", "\u8d5a", "\u9539", "\u953b", "\u821e", "\u7a33", "\u7b97", "\u7ba9", "\u7ba1", "\u50da", "\u9f3b", "\u9b44", "\u8c8c", "\u819c", "\u818a", "\u8180", "\u9c9c", "\u7591", "\u9992", "\u88f9", "\u6572", "\u8c6a", "\u818f", "\u906e", "\u8150", "\u7626", "\u8fa3", "\u7aed", "\u7aef", "\u65d7", "\u7cbe", "\u6b49", "\u7184", "\u7194", "\u6f06", "\u6f02", "\u6f2b", "\u6ef4", "\u6f14", "\u6f0f", "\u6162", "\u5be8", "\u8d5b", "\u5bdf", "\u871c", "\u8c31", "\u5ae9", "\u7fe0", "\u718a", "\u51f3", "\u9aa1", "\u7f29", "\u6167", "\u6495", "\u6492", "\u8da3", "\u8d9f", "\u6491", "\u64ad", "\u649e", "\u64a4", "\u589e", "\u806a", "\u978b", "\u8549", "\u852c", "\u6a2a", "\u69fd", "\u6a31", "\u6a61", "\u98d8", "\u918b", "\u9189", "\u9707", "\u9709", "\u7792", "\u9898", "\u66b4", "\u778e", "\u5f71", "\u8e22", "\u8e0f", "\u8e29", "\u8e2a", "\u8776", "\u8774", "\u5631", "\u58a8", "\u9547", "\u9760", "\u7a3b", "\u9ece", "\u7a3f", "\u7a3c", "\u7bb1", "\u7bad", "\u7bc7", "\u50f5", "\u8eba", "\u50fb", "\u5fb7", "\u8258", "\u819d", "\u819b", "\u719f", "\u6469", "\u989c", "\u6bc5", "\u7cca", "\u9075", "\u6f5c", "\u6f6e", "\u61c2", "\u989d", "\u6170", "\u5288", "\u64cd", "\u71d5", "\u85af", "\u85aa", "\u8584", "\u98a0", "\u6a58", "\u6574", "\u878d", "\u9192", "\u9910", "\u5634", "\u8e44", "\u5668", "\u8d60", "\u9ed8", "\u955c", "\u8d5e", "\u7bee", "\u9080", "\u8861", "\u81a8", "\u96d5", "\u78e8", "\u51dd", "\u8fa8", "\u8fa9", "\u7cd6", "\u7cd5", "\u71c3", "\u6fa1", "\u6fc0", "\u61d2", "\u58c1", "\u907f", "\u7f34", "\u6234", "\u64e6", "\u97a0", "\u85cf", "\u971c", "\u971e", "\u77a7", "\u8e48", "\u87ba", "\u7a57", "\u7e41", "\u8fab", "\u8d62", "\u7cdf", "\u7ce0", "\u71e5", "\u81c2", "\u7ffc", "\u9aa4", "\u97ad", "\u8986", "\u8e66", "\u9570", "\u7ffb", "\u9e70", "\u8b66", "\u6500", "\u8e72", "\u98a4", "\u74e3", "\u7206", "\u7586", "\u58e4", "\u8000", "\u8e81", "\u56bc", "\u56b7", "\u7c4d", "\u9b54", "\u704c", "\u8822", "\u9738", "\u9732", "\u56ca", "\u7f50", "\u5315", "\u5201", "\u4e10", "\u6b79", "\u6208", "\u592d", "\u4ed1", "\u8ba5", "\u5197", "\u9093", "\u827e", "\u592f", "\u51f8", "\u5362", "\u53ed", "\u53fd", "\u76bf", "\u51f9", "\u56da", "\u77e2", "\u4e4d", "\u5c14", "\u51af", "\u7384", "\u90a6", "\u8fc2", "\u90a2", "\u828b", "\u828d", "\u540f", "\u5937", "\u5401", "\u5415", "\u5406", "\u5c79", "\u5ef7", "\u8fc4", "\u81fc", "\u4ef2", "\u4f26", "\u4f0a", "\u808b", "\u65ed", "\u5308", "\u51eb", "\u5986", "\u4ea5", "\u6c5b", "\u8bb3", "\u8bb6", "\u8bb9", "\u8bbc", "\u8bc0", "\u5f1b", "\u9631", "\u9a6e", "\u9a6f", "\u7eab", "\u7396", "\u739b", "\u97e7", "\u62a0", "\u627c", "\u6c5e", "\u6273", "\u62a1", "\u574e", "\u575e", "\u6291", "\u62df", "\u6292", "\u8299", "\u829c", "\u82c7", "\u82a5", "\u82af", "\u82ad", "\u6756", "\u6749", "\u5deb", "\u6748", "\u752b", "\u5323", "\u8f69", "\u5364", "\u8096", "\u5431", "\u5420", "\u5455", "\u5450", "\u541f", "\u545b", "\u543b", "\u542d", "\u9091", "\u56e4", "\u542e", "\u5c96", "\u7261", "\u4f51", "\u4f43", "\u4f3a", "\u56f1", "\u809b", "\u8098", "\u7538", "\u72c8", "\u9e20", "\u5f64", "\u7078", "\u5228", "\u5e87", "\u541d", "\u5e90", "\u95f0", "\u5151", "\u707c", "\u6c90", "\u6c9b", "\u6c70", "\u6ca5", "\u6ca6", "\u6c79", "\u6ca7", "\u6caa", "\u5ff1", "\u8bc5", "\u8bc8", "\u7f55", "\u5c41", "\u5760", "\u5993", "\u59ca", "\u5992", "\u7eac", "\u73ab", "\u5366", "\u5777", "\u576f", "\u62d3", "\u576a", "\u5764", "\u62c4", "\u62e7", "\u62c2", "\u62d9", "\u62c7", "\u62d7", "\u8309", "\u6614", "\u82db", "\u82eb", "\u82df", "\u82de", "\u8301", "\u82d4", "\u6789", "\u67a2", "\u679a", "\u67ab", "\u676d", "\u90c1", "\u77fe", "\u5948", "\u5944", "\u6bb4", "\u6b67", "\u5353", "\u6619", "\u54ce", "\u5495", "\u5475", "\u5499", "\u547b", "\u5492", "\u5486", "\u5496", "\u5e15", "\u8d26", "\u8d2c", "\u8d2e", "\u6c1b", "\u79c9", "\u5cb3", "\u4fa0", "\u4fa5", "\u4fa3", "\u4f88", "\u5351", "\u523d", "\u5239", "\u80b4", "\u89c5", "\u5fff", "\u74ee", "\u80ae", "\u80aa", "\u72de", "\u5e9e", "\u759f", "\u7599", "\u759a", "\u5352", "\u6c13", "\u70ac", "\u6cbd", "\u6cae", "\u6ce3", "\u6cde", "\u6ccc", "\u6cbc", "\u6014", "\u602f", "\u5ba0", "\u5b9b", "\u8869", "\u7948", "\u8be1", "\u5e1a", "\u5c49", "\u5f27", "\u5f25", "\u964b", "\u964c", "\u51fd", "\u59c6", "\u8671", "\u53c1", "\u7ec5", "\u9a79", "\u7eca", "\u7ece", "\u5951", "\u8d30", "\u73b7", "\u73b2", "\u73ca", "\u62ed", "\u62f7", "\u62f1", "\u631f", "\u57a2", "\u579b", "\u62ef", "\u8346", "\u8338", "\u832c", "\u835a", "\u8335", "\u8334", "\u835e", "\u8360", "\u8364", "\u8367", "\u8354", "\u6808", "\u67d1", "\u6805", "\u67e0", "\u67b7", "\u52c3", "\u67ec", "\u7802", "\u6cf5", "\u781a", "\u9e25", "\u8f74", "\u97ed", "\u8650", "\u6627", "\u76f9", "\u54a7", "\u6635", "\u662d", "\u76c5", "\u52cb", "\u54c6", "\u54aa", "\u54df", "\u5e7d", "\u9499", "\u949d", "\u94a0", "\u94a6", "\u94a7", "\u94ae", "\u6be1", "\u6c22", "\u79d5", "\u4fcf", "\u4fc4", "\u4fd0", "\u4faf", "\u5f8a", "\u884d", "\u80da", "\u80e7", "\u80ce", "\u72f0", "\u9975", "\u5ce6", "\u5955", "\u54a8", "\u98d2", "\u95fa", "\u95fd", "\u7c7d", "\u5a04", "\u70c1", "\u70ab", "\u6d3c", "\u67d2", "\u6d8e", "\u6d1b", "\u6043", "\u604d", "\u606c", "\u6064", "\u5ba6", "\u8beb", "\u8bec", "\u7960", "\u8bf2", "\u5c4f", "\u5c4e", "\u900a", "\u9668", "\u59da", "\u5a1c", "\u86a4", "\u9a87", "\u8018", "\u8019", "\u79e6", "\u533f", "\u57c2", "\u6342", "\u634d", "\u8881", "\u634c", "\u632b", "\u631a", "\u6363", "\u6345", "\u57c3", "\u803f", "\u8042", "\u8378", "\u83bd", "\u83b1", "\u8389", "\u83b9", "\u83ba", "\u6886", "\u6816", "\u6866", "\u6813", "\u6845", "\u6869", "\u8d3e", "\u914c", "\u7838", "\u7830", "\u783e", "\u6b89", "\u901e", "\u54ee", "\u5520", "\u54fa", "\u5254", "\u868c", "\u869c", "\u7554", "\u86a3", "\u86aa", "\u8693", "\u54e9", "\u5703", "\u9e2f", "\u5501", "\u54fc", "\u5506", "\u5ced", "\u5527", "\u5cfb", "\u8d42", "\u8d43", "\u94be", "\u94c6", "\u6c28", "\u79eb", "\u7b06", "\u4ffa", "\u8d41", "\u5014", "\u6bb7", "\u8038", "\u8200", "\u8c7a", "\u8c79", "\u9881", "\u80ef", "\u80f0", "\u8110", "\u8113", "\u901b", "\u537f", "\u9e35", "\u9e33", "\u9981", "\u51cc", "\u51c4", "\u8877", "\u90ed", "\u658b", "\u75b9", "\u7d0a", "\u74f7", "\u7f94", "\u70d9", "\u6d66", "\u6da1", "\u6da3", "\u6da4", "\u6da7", "\u6d95", "\u6da9", "\u608d", "\u60af", "\u7a8d", "\u8bfa", "\u8bfd", "\u8892", "\u8c06", "\u795f", "\u6055", "\u5a29", "\u9a8f", "\u7410", "\u9eb8", "\u7409", "\u7405", "\u63aa", "\u637a", "\u6376", "\u8d66", "\u57e0", "\u637b", "\u6390", "\u6382", "\u6396", "\u63b7", "\u63b8", "\u63ba", "\u52d8", "\u804a", "\u5a36", "\u83f1", "\u83f2", "\u840e", "\u83e9", "\u8424", "\u4e7e", "\u8427", "\u8428", "\u83c7", "\u5f6c", "\u6897", "\u68a7", "\u68ad", "\u66f9", "\u915d", "\u9157", "\u53a2", "\u7845", "\u7855", "\u5962", "\u76d4", "\u533e", "\u9885", "\u5f6a", "\u7736", "\u6664", "\u66fc", "\u6666", "\u5195", "\u5561", "\u7566", "\u8dbe", "\u5543", "\u86c6", "\u86af", "\u86c9", "\u86c0", "\u552c", "\u553e", "\u5564", "\u5565", "\u5578", "\u5d0e", "\u903b", "\u5d14", "\u5d29", "\u5a74", "\u8d4a", "\u94d0", "\u94db", "\u94dd", "\u94e1", "\u94e3", "\u94ed", "\u77eb", "\u79f8", "\u79fd", "\u7b19", "\u7b24", "\u504e", "\u5080", "\u8eaf", "\u515c", "\u8845", "\u5f98", "\u5f99", "\u8236", "\u8237", "\u8235", "\u655b", "\u7fce", "\u812f", "\u9038", "\u51f0", "\u7316", "\u796d", "\u70f9", "\u5eb6", "\u5eb5", "\u75ca", "\u960e", "\u9610", "\u7737", "\u710a", "\u7115", "\u9e3f", "\u6daf", "\u6dd1", "\u6dcc", "\u6dee", "\u6dc6", "\u6e0a", "\u6deb", "\u6df3", "\u6de4", "\u6dc0", "\u6dae", "\u6db5", "\u60e6", "\u60b4", "\u60cb", "\u5bc2", "\u7a92", "\u8c0d", "\u8c10", "\u88c6", "\u88b1", "\u7977", "\u8c12", "\u8c13", "\u8c1a", "\u5c09", "\u5815", "\u9685", "\u5a49", "\u9887", "\u7ef0", "\u7ef7", "\u7efc", "\u7efd", "\u7f00", "\u5de2", "\u7433", "\u7422", "\u743c", "\u63cd", "\u5830", "\u63e9", "\u63fd", "\u63d6", "\u5f6d", "\u63e3", "\u6400", "\u6413", "\u58f9", "\u6414", "\u846b", "\u52df", "\u848b", "\u8482", "\u97e9", "\u68f1", "\u6930", "\u711a", "\u690e", "\u68fa", "\u6994", "\u692d", "\u7c9f", "\u68d8", "\u9163", "\u9165", "\u785d", "\u786b", "\u988a", "\u96f3", "\u7fd8", "\u51ff", "\u68e0", "\u6670", "\u9f0e", "\u55b3", "\u904f", "\u667e", "\u7574", "\u8dcb", "\u8ddb", "\u86d4", "\u8712", "\u86e4", "\u9e43", "\u55bb", "\u557c", "\u55a7", "\u5d4c", "\u8d4b", "\u8d4e", "\u8d50", "\u9509", "\u950c", "\u7525", "\u63b0", "\u6c2e", "\u6c2f", "\u9ecd", "\u7b4f", "\u724d", "\u7ca4", "\u903e", "\u814c", "\u814b", "\u8155", "\u7329", "\u732c", "\u60eb", "\u6566", "\u75d8", "\u75e2", "\u75ea", "\u7ae3", "\u7fd4", "\u5960", "\u9042", "\u7119", "\u6ede", "\u6e58", "\u6e24", "\u6e3a", "\u6e83", "\u6e85", "\u6e43", "\u6115", "\u60f6", "\u5bd3", "\u7a96", "\u7a98", "\u96c7", "\u8c24", "\u7280", "\u9698", "\u5a92", "\u5a9a", "\u5a7f", "\u7f05", "\u7f06", "\u7f14", "\u7f15", "\u9a9a", "\u745f", "\u9e49", "\u7470", "\u642a", "\u8058", "\u659f", "\u9774", "\u9776", "\u84d6", "\u84bf", "\u84b2", "\u84c9", "\u6954", "\u693f", "\u6977", "\u6984", "\u695e", "\u6963", "\u916a", "\u7898", "\u787c", "\u7889", "\u8f90", "\u8f91", "\u9891", "\u7779", "\u7766", "\u7784", "\u55dc", "\u55e6", "\u6687", "\u7578", "\u8df7", "\u8dfa", "\u8708", "\u8717", "\u8715", "\u86f9", "\u55c5", "\u55e1", "\u55e4", "\u7f72", "\u8700", "\u5e4c", "\u951a", "\u9525", "\u9528", "\u952d", "\u9530", "\u7a1a", "\u9893", "\u7b77", "\u9b41", "\u8859", "\u817b", "\u816e", "\u817a", "\u9e4f", "\u8084", "\u733f", "\u9896", "\u715e", "\u96cf", "\u998d", "\u998f", "\u7980", "\u75f9", "\u5ed3", "\u75f4", "\u9756", "\u8a8a", "\u6f13", "\u6ea2", "\u6eaf", "\u6eb6", "\u6ed3", "\u6eba", "\u5bde", "\u7aa5", "\u7a9f", "\u5bdd", "\u8902", "\u88f8", "\u8c2c", "\u5ab3", "\u5ac9", "\u7f1a", "\u7f24", "\u527f", "\u8d58", "\u71ac", "\u8d6b", "\u852b", "\u6479", "\u8513", "\u8517", "\u853c", "\u7199", "\u851a", "\u5162", "\u699b", "\u6995", "\u9175", "\u789f", "\u78b4", "\u78b1", "\u78b3", "\u8f95", "\u8f96", "\u96cc", "\u5885", "\u5601", "\u8e0a", "\u8749", "\u5600", "\u5e54", "\u9540", "\u8214", "\u718f", "\u7b8d", "\u7b95", "\u7bab", "\u8206", "\u50e7", "\u5b75", "\u7629", "\u761f", "\u5f70", "\u7cb9", "\u6f31", "\u6f29", "\u6f3e", "\u6177", "\u5be1", "\u5be5", "\u8c2d", "\u8910", "\u892a", "\u96a7", "\u5ae1", "\u7f28", "\u64b5", "\u64a9", "\u64ae", "\u64ac", "\u64d2", "\u58a9", "\u64b0", "\u978d", "\u854a", "\u8574", "\u6a0a", "\u6a1f", "\u6a44", "\u6577", "\u8c4c", "\u9187", "\u78d5", "\u78c5", "\u78be", "\u618b", "\u5636", "\u5632", "\u5639", "\u8760", "\u874e", "\u874c", "\u8757", "\u8759", "\u563f", "\u5e62", "\u954a", "\u9550", "\u7a3d", "\u7bd3", "\u8198", "\u9ca4", "\u9cab", "\u8912", "\u762a", "\u7624", "\u762b", "\u51db", "\u6f8e", "\u6f6d", "\u6f66", "\u6fb3", "\u6f58", "\u6f88", "\u6f9c", "\u6f84", "\u6194", "\u61ca", "\u618e", "\u7fe9", "\u8925", "\u8c34", "\u9e64", "\u61a8", "\u5c65", "\u5b09", "\u8c6b", "\u7f2d", "\u64bc", "\u64c2", "\u64c5", "\u857e", "\u859b", "\u8587", "\u64ce", "\u7ff0", "\u5669", "\u6a71", "\u6a59", "\u74e2", "\u87e5", "\u970d", "\u970e", "\u8f99", "\u5180", "\u8e31", "\u8e42", "\u87c6", "\u8783", "\u879f", "\u566a", "\u9e66", "\u9ed4", "\u7a46", "\u7be1", "\u7bf7", "\u7bd9", "\u7bf1", "\u5112", "\u81b3", "\u9cb8", "\u763e", "\u7638", "\u7cd9", "\u71ce", "\u6fd2", "\u61be", "\u61c8", "\u7abf", "\u7f30", "\u58d5", "\u85d0", "\u6aac", "\u6a90", "\u6aa9", "\u6a80", "\u7901", "\u78f7", "\u4e86", "\u77ac", "\u77b3", "\u77aa", "\u66d9", "\u8e4b", "\u87cb", "\u87c0", "\u568e", "\u8d61", "\u9563", "\u9b4f", "\u7c07", "\u5121", "\u5fbd", "\u7235", "\u6726", "\u81ca", "\u9cc4", "\u7cdc", "\u764c", "\u61e6", "\u8c41", "\u81c0", "\u85d5", "\u85e4", "\u77bb", "\u56a3", "\u9ccd", "\u765e", "\u7011", "\u895f", "\u74a7", "\u6233", "\u6512", "\u5b7d", "\u8611", "\u85fb", "\u9cd6", "\u8e6d", "\u8e6c", "\u7c38", "\u7c3f", "\u87f9", "\u9761", "\u7663", "\u7fb9", "\u9b13", "\u6518", "\u8815", "\u5dcd", "\u9cde", "\u7cef", "\u8b6c", "\u9739", "\u8e8f", "\u9ad3", "\u8638", "\u9576", "\u74e4", "\u77d7"}

//RandomCreateSimplifyZHCNUnicode 随机生成简化汉字中文Unicode编码及转义字符
func RandomCreateSimplifyZHCNUnicode() (string, string) {
	uStr := simplifyZHCNChars[rand.Intn(len(simplifyZHCNChars))]
	str, _ := strconv.Unquote(strings.Replace(strconv.Quote(uStr), `\\u`, `\u`, -1))
	for _, r := range str {
		if !unicode.Is(unicode.Scripts["Han"], r) {
			return RandomCreateSimplifyZHCNUnicode()
		}
	}
	return uStr, str
}

const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int, letter ...string) string {
	letterX := letters
	if len(letter) > 0 {
		letterX = letter[0]
	}
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letterX) {
			b[i] = letterX[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// RandInt 生成区间[-m, n]的安全随机数
func RandInt(min, max int) int {
	if min > max {
		return max
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int(f64Min)
		result, _ := randC.Int(randC.Reader, big.NewInt(int64(max+1+i64Min)))

		return int(result.Int64() - int64(i64Min))
	}
	result, _ := randC.Int(randC.Reader, big.NewInt(int64(max-min+1)))
	return int(int64(min) + result.Int64())
}

//GetRandomStringValue 从字符串列表中随机获取一个值
func GetRandomStringValue(s []string) string {
	sLen := len(s)
	index := RandInt(0, sLen)
	if index >= sLen {
		index = sLen - 1
	}
	return s[index]
}

//GetRandomColorValueByRGBA 从Color列表中随机获取一个颜色的RGBA值
func GetRandomColorValueByRGBA(cs []color.Color) color.RGBA {
	cLen := len(cs)
	index := RandInt(0, cLen)
	if index >= cLen {
		index = cLen - 1
	}
	r, g, b, a := cs[index].RGBA()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

//GenCaptchaKey 生成唯一KEY
func GenCaptchaKey(str string) (string, error) {
	t := GenUniqueId()
	keyStr := Md5ToStr(str + t)
	return keyStr, nil
}

//Md5ToStr MD5加密字符串
func Md5ToStr(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

var num int64

//GenUniqueId 创建唯一ID
func GenUniqueId() string {
	t := time.Now()
	s := t.Format("20060102150405")
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := fmt.Sprintf("%0*d", 3, m)
	p := os.Getpid() % 1000
	ps := fmt.Sprintf("%0*d", 3, int64(p))
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := fmt.Sprintf("%0*d", 4, r)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}
