package echoex

import (
	"net"

	"github.com/labstack/echo/v4"
)

type cfgOptions struct {
	e                     *echo.Echo
	trustOpts             []echo.TrustOption
	trustXFF, trustRealIP bool
}

func (o *cfgOptions) apply() {
	if o.trustXFF {
		o.e.IPExtractor = echo.ExtractIPFromXFFHeader(o.trustOpts...)
	} else if o.trustRealIP {
		o.e.IPExtractor = echo.ExtractIPFromRealIPHeader(o.trustOpts...)
	}
}

type CfgOption func(o *cfgOptions)

// CfgTrustedCIDRForXFF instructs to trust X-Forwarded-For for this CIDR
func CfgTrustedCIDRForXFF(cidr string) CfgOption {
	return func(o *cfgOptions) {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}

		if o.trustRealIP {
			panic("cannot trust both XFF and RealIP")
		}
		o.trustXFF = true

		o.trustOpts = append(o.trustOpts, echo.TrustIPRange(ipnet))
	}
}

// CfgTrustedCIDRForRealIP instructs to trust X-Real-IP for this CIDR
func CfgTrustedCIDRForRealIP(cidr string) CfgOption {
	return func(o *cfgOptions) {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}

		if o.trustXFF {
			panic("cannot trust both XFF and RealIP")
		}
		o.trustRealIP = true

		o.trustOpts = append(o.trustOpts, echo.TrustIPRange(ipnet))
	}
}

// New creates a new echo.Echo instance with custom validator and error handler configured.
func New(cfg ...CfgOption) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	opts := &cfgOptions{e: e}
	for _, c := range cfg {
		c(opts)
	}
	opts.apply()

	e.Validator = NewCustomValidator()
	e.HTTPErrorHandler = CustomHTTPErrorHandler

	return e
}
