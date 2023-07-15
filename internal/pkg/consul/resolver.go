package consul

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-playground/form"
	"github.com/hashicorp/consul/api"
	"github.com/jpillora/backoff"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"

	"github.com/kirychukyurii/wasker/internal/errors"
)

// init function needs for  auto-register in resolvers registry
func init() {
	resolver.Register(&builder{})
}

// resolvr implements resolver.Resolver from the gRPC package.
// It watches for endpoints changes and pushes them to the underlying gRPC connection.
type resolvr struct {
	cancelFunc context.CancelFunc
}

// ResolveNow will be skipped due unnecessary in this case
func (r *resolvr) ResolveNow(resolver.ResolveNowOptions) {}

// Close closes the resolver.
func (r *resolvr) Close() {
	r.cancelFunc()
}

type servicer interface {
	Service(string, string, bool, *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error)
}

func watchConsulService(ctx context.Context, s servicer, tgt target, out chan<- []string) {
	res := make(chan []string)
	quit := make(chan struct{})
	bck := &backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    10 * time.Millisecond,
		Max:    tgt.MaxBackoff,
	}
	go func() {
		var lastIndex uint64
		for {
			ss, meta, err := s.Service(
				tgt.Service,
				tgt.Tag,
				tgt.Healthy,
				&api.QueryOptions{
					WaitIndex:         lastIndex,
					Near:              tgt.Near,
					WaitTime:          tgt.Wait,
					Datacenter:        tgt.Dc,
					AllowStale:        tgt.AllowStale,
					RequireConsistent: tgt.RequireConsistent,
				},
			)
			if err != nil {
				// No need to continue if the context is done/cancelled.
				// We check that here directly because the check for the closed quit channel
				// at the end of the loop is not reached when calling continue here.
				select {
				case <-quit:
					return
				default:
					grpclog.Errorf("[Consul resolver] Couldn't fetch endpoints. target={%s}; error={%v}", tgt.String(), err)
					time.Sleep(bck.Duration())
					continue
				}
			}
			bck.Reset()
			lastIndex = meta.LastIndex
			grpclog.Infof("[Consul resolver] %d endpoints fetched in(+wait) %s for target={%s}",
				len(ss),
				meta.RequestTime,
				tgt.String(),
			)

			ee := make([]string, 0, len(ss))
			for _, s := range ss {
				address := s.Service.Address
				if s.Service.Address == "" {
					address = s.Node.Address
				}
				ee = append(ee, fmt.Sprintf("%s:%d", address, s.Service.Port))
			}

			if tgt.Limit != 0 && len(ee) > tgt.Limit {
				ee = ee[:tgt.Limit]
			}
			select {
			case res <- ee:
				continue
			case <-quit:
				return
			}
		}
	}()

	for {
		// If in the below select both channels have values that can be read,
		// Go picks one pseudo-randomly.
		// But when the context is canceled we want to act upon it immediately.
		if ctx.Err() != nil {
			// Close quit so the goroutine returns and doesn't leak.
			// Do NOT close res because that can lead to panics in the goroutine.
			// res will be garbage collected at some point.
			close(quit)
			return
		}
		select {
		case ee := <-res:
			out <- ee
		case <-ctx.Done():
			close(quit)
			return
		}
	}
}

func populateEndpoints(ctx context.Context, clientConn resolver.ClientConn, input <-chan []string) {
	for {
		select {
		case cc := <-input:
			connsSet := make(map[string]struct{}, len(cc))
			for _, c := range cc {
				connsSet[c] = struct{}{}
			}
			conns := make([]resolver.Address, 0, len(connsSet))
			for c := range connsSet {
				conns = append(conns, resolver.Address{Addr: c})
			}
			sort.Sort(byAddressString(conns)) // Don't replace the same address list in the balancer
			err := clientConn.UpdateState(resolver.State{Addresses: conns})
			if err != nil {
				grpclog.Errorf("[Consul resolver] Couldn't update client connection. error={%v}", err)
				continue
			}
		case <-ctx.Done():
			grpclog.Info("[Consul resolver] Watch has been finished")
			return
		}
	}
}

// byAddressString sorts resolver.Address by Address Field  sorting in increasing order.
type byAddressString []resolver.Address

func (p byAddressString) Len() int           { return len(p) }
func (p byAddressString) Less(i, j int) bool { return p[i].Addr < p[j].Addr }
func (p byAddressString) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// schemeName for the urls
// All target URLs like 'consul://.../...' will be resolved by this resolver
const schemeName = "consul"

// builder implements resolver.Builder and use for constructing all consul resolvers
type builder struct{}

func (b *builder) Build(url resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	dsn := strings.Join([]string{schemeName + ":/", url.URL.Host, url.URL.Path + "?" + url.URL.RawQuery}, "/")
	tgt, err := parseURL(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong consul URL")
	}
	cli, err := api.NewClient(tgt.consulConfig())
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't connect to the Consul API")
	}

	ctx, cancel := context.WithCancel(context.Background())
	pipe := make(chan []string)
	go watchConsulService(ctx, cli.Health(), tgt, pipe)
	go populateEndpoints(ctx, cc, pipe)

	return &resolvr{cancelFunc: cancel}, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (b *builder) Scheme() string {
	return schemeName
}

type target struct {
	Addr              string        `form:"-"`
	User              string        `form:"-"`
	Password          string        `form:"-"`
	Service           string        `form:"-"`
	Wait              time.Duration `form:"wait"`
	Timeout           time.Duration `form:"timeout"`
	MaxBackoff        time.Duration `form:"max-backoff"`
	Tag               string        `form:"tag"`
	Near              string        `form:"near"`
	Limit             int           `form:"limit"`
	Healthy           bool          `form:"healthy"`
	TLSInsecure       bool          `form:"insecure"`
	Token             string        `form:"token"`
	Dc                string        `form:"dc"`
	AllowStale        bool          `form:"allow-stale"`
	RequireConsistent bool          `form:"require-consistent"`
}

func (t *target) String() string {
	return fmt.Sprintf("service='%s' healthy='%t' tag='%s'", t.Service, t.Healthy, t.Tag)
}

// parseURL with parameters
//
// consul://[user:password@]127.0.0.127:8555/my-service?[healthy=]&[wait=]&[near=]&[insecure=]&[limit=]&[tag=]&[token=]
// URL schema will stay stable in the future for backward compatibility
func parseURL(u string) (target, error) {
	rawURL, err := url.Parse(u)
	if err != nil {
		return target{}, errors.Wrap(err, "Malformed URL")
	}

	if rawURL.Scheme != schemeName ||
		len(rawURL.Host) == 0 || len(strings.TrimLeft(rawURL.Path, "/")) == 0 {
		return target{},
			errors.Errorf("Malformed URL('%s'). Must be in the next format: 'consul://[user:passwd]@host/service?param=value'", u)
	}

	var tgt target
	tgt.User = rawURL.User.Username()
	tgt.Password, _ = rawURL.User.Password()
	tgt.Addr = rawURL.Host
	tgt.Service = strings.TrimLeft(rawURL.Path, "/")
	decoder := form.NewDecoder()
	decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.ParseDuration(vals[0])
	}, time.Duration(0))

	err = decoder.Decode(&tgt, rawURL.Query())
	if err != nil {
		return target{}, errors.Wrap(err, "Malformed URL parameters")
	}
	if len(tgt.Near) == 0 {
		tgt.Near = "_agent"
	}
	if tgt.MaxBackoff == 0 {
		tgt.MaxBackoff = time.Second
	}
	return tgt, nil
}

// consulConfig returns config based on the parsed target.
// It uses custom http-client.
func (t *target) consulConfig() *api.Config {
	var creds *api.HttpBasicAuth
	if len(t.User) > 0 && len(t.Password) > 0 {
		creds = new(api.HttpBasicAuth)
		creds.Password = t.Password
		creds.Username = t.User
	}

	// custom http.Client
	c := &http.Client{
		Timeout: t.Timeout,
	}
	return &api.Config{
		Address:    t.Addr,
		HttpAuth:   creds,
		WaitTime:   t.Wait,
		HttpClient: c,
		TLSConfig: api.TLSConfig{
			InsecureSkipVerify: t.TLSInsecure,
		},
		Token: t.Token,
	}
}
