package endpoint

import "github.com/ogiusek/ioc/v2"

type anyRequest interface {
	GetC() ioc.Dic
}

type AnyRequest struct{ C ioc.Dic }

func (req AnyRequest) GetC() ioc.Dic { return req.C }

func NewAnyRequest() anyRequest { return &AnyRequest{} }
