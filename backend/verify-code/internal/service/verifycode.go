package service

import (
	"context"
	"math/rand"
	"strings"
	pb "verify-code/api/verifyCode"
)

type VerifyCodeService struct {
	pb.UnimplementedVerifyCodeServer
}

func NewVerifyCodeService() *VerifyCodeService {
	return &VerifyCodeService{}
}

func (s *VerifyCodeService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeReply, error) {
	return &pb.GetVerifyCodeReply{
		Code: RandCode(int(req.Length), req.Type),
	}, nil
}
func RandCode(lenth int, t pb.TYPE) string {
	switch t {
	case pb.TYPE_DEFAULT:
		fallthrough
	case pb.TYPE_DIGIT:

		return randCode("123456789", 4, lenth)
	case pb.TYPE_LETTER:

		return randCode("abcdefghijklmnopqrstuvwxyz", 5, lenth)
	case pb.TYPE_MIXED:

		return randCode("abcdefghijklmnopqrstuvwxyz123456789", 6, lenth)
	default:

	}
	return ""
}

// 随机数实现。一次随机（63位）多次利用
func randCode(chars string, idxBits, l int) string {
	// 形成掩码
	idxMask := 1<<idxBits - 1
	// 63 位可以使用的最大组次数
	idxMax := 63 / idxBits

	// 利用string builder构建结果缓冲
	sb := strings.Builder{}
	sb.Grow(l)

	// 循环生成随机数
	// i 索引
	// cache 随机数缓存
	// remain 随机数还可以用几次
	for i, cache, remain := l-1, rand.Int63(), idxMax; i >= 0; {
		// 随机缓存不足，重新生成
		if remain == 0 {
			cache, remain = rand.Int63(), idxMax
		}
		// 利用掩码生成随机索引，有效索引为小于字符集合长度
		if idx := int(cache & int64(idxMask)); idx < len(chars) {
			sb.WriteByte(chars[idx])
			i--
		}
		// 利用下一组随机数位
		cache >>= idxBits
		remain--
	}

	return sb.String()
}

//随机数简单实现
//func randCode(chars string, l int) string {
//	charsLen := len(chars)
//	result := make([]byte, l)
//	//rand.Seed(time.Now().UnixNano())//全局seed，已放在main函数内
//	for i := 0; i < l; i++ {
//		randIndex := rand.Intn(charsLen)
//		result[i] = chars[randIndex]
//	}
//	return string(result)
//
//}
