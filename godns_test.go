package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/miekg/dns"
)

const (
	nameserver = "127.0.0.1:8553"
	domain     = "www.sina.com.cn"
)

func TestDig(t *testing.T) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	c := new(dns.Client)
	exchange, rtt, err := c.Exchange(m, nameserver)
	assert.NoError(t, err)
	fmt.Println(exchange)
	fmt.Println(rtt)
	fmt.Println(err)
}
func BenchmarkDig(b *testing.B) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	c := new(dns.Client)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = c.Exchange(m, nameserver)
	}

}
