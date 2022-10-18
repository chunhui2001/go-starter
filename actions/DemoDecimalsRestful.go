package actions

import (
	"github.com/gin-gonic/gin"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/shopspring/decimal"
)

// price, err := decimal.NewFromString("136.02")
// quantity := decimal.NewFromInt(3)

// fee, _ := decimal.NewFromString(".035")
// taxRate, _ := decimal.NewFromString(".08875")
// subtotal := price.Mul(quantity)
// preTax := subtotal.Mul(fee.Add(decimal.NewFromFloat(1)))
// total := preTax.Mul(taxRate.Add(decimal.NewFromFloat(1)))

func DecimalNewFromStringHandler(c *gin.Context) {

	var fixedDecimals int32 = 68

	floatString := c.Query("floatString")
	price, _ := decimal.NewFromString(floatString)
	quantity := decimal.NewFromInt(3)

	// n1, _ := decimal.NewFromString("1.0")
	n1 := decimal.NewFromFloat(1.0)
	n2 := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(fixedDecimals)))

	n11 := decimal.NewFromFloat(1.0)
	n22 := decimal.NewFromFloat(3.0)

	c.JSON(200, (&R{Data: []decimal.Decimal{
		price,
		price.Mul(quantity),
		n1.DivRound(n2, fixedDecimals),
		n11.DivRound(n22, fixedDecimals),
		decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(fixedDecimals)))}, Error: nil}).IfErr(400))

}
