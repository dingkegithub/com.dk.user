package impl

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

//
//
// bcrypt.GenerateFromPassword(rawPwd, c)
// 根据c的不同，耗时不一样，越慢越保密越难破解，根据
// 实际业务需求测试得到合适的c，默认c是10
//
//
func TestPwd(t *testing.T) {
	rawPwd := []byte("123456")
	for c:=0; c<30; c++ {
		start := time.Now()
		pwd, _ := bcrypt.GenerateFromPassword(rawPwd, c)
		lost := time.Since(start)
		fmt.Println(c, "--->", lost, "  pwd: ", string(pwd))
	}
}
