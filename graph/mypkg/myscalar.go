package mypkg

import (
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/shopspring/decimal"
)

func MarshalDecimal(b decimal.Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(b.String()))
	})
}

func UnmarshalDecimal(v interface{}) (decimal.Decimal, error) {
	return utils.DecimalFromString(utils.ToString(v)), nil
}
