package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/eb4uk/godns/models"
	"github.com/spaolacci/murmur3"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

type Mesg struct {
	Msg    *dns.Msg
	Expire time.Time
}

type Cache interface {
	Get(key string) (Msg *dns.Msg, err error)
	Set(key string, Msg *dns.Msg) error
	Exists(key string) bool
	Remove(key string) error
	Full() bool
}

//TODO Move hash setting to config.

//Murmur used for better performance
var hashFunc func(q models.Question) string = hashMurmur

func KeyGen(q models.Question) string {
	return hashFunc(q)
}

func hashMurmur(q models.Question) string {
	//murmur3.
	x := murmur3.Sum64([]byte(q.String()))
	return strconv.FormatInt(int64(x), 16)
}

func hashMd5(q models.Question) string {
	h := md5.New()
	h.Write([]byte(q.String()))
	x := h.Sum(nil)
	key := fmt.Sprintf("%x", x)
	return key
}

//TODO [PERFORMANCE] Use jsoniter for better performance

/* we need to define marsheling to encode and decode
 */
type JsonSerializer struct {
}

func (*JsonSerializer) Dumps(mesg *dns.Msg) (encoded []byte, err error) {
	encoded, err = json.Marshal(*mesg)
	return
}

func (*JsonSerializer) Loads(data []byte) (*dns.Msg, error) {
	var mesg dns.Msg
	err := json.Unmarshal(data, &mesg)
	return &mesg, err
}
