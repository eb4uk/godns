package main

import (
	cache2 "github.com/eb4uk/godns/cache"
	"github.com/eb4uk/godns/models"
	"github.com/eb4uk/godns/settings"
	"net"
	"time"

	"github.com/miekg/dns"
)

const (
	notIPQuery = 0
	_IP4Query  = 4
	_IP6Query  = 6
)

type GODNSHandler struct {
	resolver        *Resolver
	cache, negCache cache2.Cache
	hosts           Hosts
}

func NewHandler() *GODNSHandler {

	var (
		cacheConfig     settings.CacheSettings
		resolver        *Resolver
		cache, negCache cache2.Cache
	)

	resolver = NewResolver(settings.Config.ResolvConfig)

	cacheConfig = settings.Config.Cache
	switch cacheConfig.Backend {
	case "memory":
		cache = &cache2.MemoryCache{
			Backend:  make(map[string]cache2.Mesg, cacheConfig.Maxcount),
			Expire:   time.Duration(cacheConfig.Expire) * time.Second,
			Maxcount: cacheConfig.Maxcount,
		}
		negCache = &cache2.MemoryCache{
			Backend:  make(map[string]cache2.Mesg),
			Expire:   time.Duration(cacheConfig.Expire) * time.Second / 2,
			Maxcount: cacheConfig.Maxcount,
		}
	case "memcache":
		cache = cache2.NewMemcachedCache(
			settings.Config.Memcache.Servers,
			int32(cacheConfig.Expire))
		negCache = cache2.NewMemcachedCache(
			settings.Config.Memcache.Servers,
			int32(cacheConfig.Expire/2))
	case "redis":
		cache = cache2.NewRedisCache(
			settings.Config.Redis,
			int64(cacheConfig.Expire))
		negCache = cache2.NewRedisCache(
			settings.Config.Redis,
			int64(cacheConfig.Expire/2))
	default:
		logger.Error("Invalid cache backend %s", cacheConfig.Backend)
		panic("Invalid cache backend")
	}

	var hosts Hosts
	if settings.Config.Hosts.Enable {
		hosts = NewHosts(settings.Config.Hosts, settings.Config.Redis)
	}

	return &GODNSHandler{resolver, cache, negCache, hosts}
}

func (h *GODNSHandler) do(Net string, w dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	Q := models.Question{UnFqdn(q.Name), dns.TypeToString[q.Qtype], dns.ClassToString[q.Qclass]}

	var remote net.IP
	if Net == "tcp" {
		remote = w.RemoteAddr().(*net.TCPAddr).IP
	} else {
		remote = w.RemoteAddr().(*net.UDPAddr).IP
	}
	//TODO Check is specific key set up
	logger.Info("%s lookup　%s", remote, Q.String())

	IPQuery := h.isIPQuery(q)

	// Query hosts
	if settings.Config.Hosts.Enable && IPQuery > 0 {
		if ips, ok := h.hosts.Get(Q.Name, IPQuery); ok {
			m := new(dns.Msg)
			m.SetReply(req)

			switch IPQuery {
			case _IP4Query:
				rr_header := dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    settings.Config.Hosts.TTL,
				}
				for _, ip := range ips {
					a := &dns.A{rr_header, ip}
					m.Answer = append(m.Answer, a)
				}
			case _IP6Query:
				rr_header := dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    settings.Config.Hosts.TTL,
				}
				for _, ip := range ips {
					aaaa := &dns.AAAA{rr_header, ip}
					m.Answer = append(m.Answer, aaaa)
				}
			}

			w.WriteMsg(m)
			logger.Debug("%s found in hosts file", Q.Name)
			return
		} else {
			logger.Debug("%s didn't found in hosts file", Q.Name)
		}
	}

	key := cache2.KeyGen(Q)
	mesg, err := h.cache.Get(key)
	if err != nil {
		if mesg, err = h.negCache.Get(key); err != nil {
			logger.Debug("%s didn't hit cache", Q.String())
		} else {
			logger.Debug("%s hit negative cache", Q.String())
			dns.HandleFailed(w, req)
			return
		}
	} else {
		logger.Debug("%s hit cache", Q.String())
		// we need this copy against concurrent modification of Id
		msg := *mesg
		msg.Id = req.Id
		w.WriteMsg(&msg)
		return
	}

	mesg, err = h.resolver.Lookup(Net, req)

	if err != nil {
		logger.Warn("Resolve query error %s", err)
		dns.HandleFailed(w, req)

		// cache the failure, too!
		if err = h.negCache.Set(key, nil); err != nil {
			logger.Warn("Set %s negative cache failed: %v", Q.String(), err)
		}
		return
	}

	w.WriteMsg(mesg)

	if len(mesg.Answer) > 0 {
		err = h.cache.Set(key, mesg)
		if err != nil {
			logger.Warn("Set %s cache failed: %s", Q.String(), err.Error())
		}
		logger.Debug("Insert %s into cache", Q.String())
	}
}

func (h *GODNSHandler) DoTCP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("tcp", w, req)
}

func (h *GODNSHandler) DoUDP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("udp", w, req)
}

func (h *GODNSHandler) isIPQuery(q dns.Question) int {
	if q.Qclass != dns.ClassINET {
		return notIPQuery
	}

	switch q.Qtype {
	case dns.TypeA:
		return _IP4Query
	case dns.TypeAAAA:
		return _IP6Query
	default:
		return notIPQuery
	}
}

func UnFqdn(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s)-1]
	}
	return s
}
