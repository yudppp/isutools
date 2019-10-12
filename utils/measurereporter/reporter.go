package measurereporter

import (
	"strings"
	"fmt"

	"github.com/najeira/measure"
	"github.com/yudppp/isutools/utils/slackcat"
)

// Send .
func Send(filename string, data measure.StatsSlice) error {
	var b strings.Builder
	fmt.Fprint(&b,`"Key","Count","Sum","Min","Max","Avg","Rate","P95"`)
	for _, row := range data {
		fmt.Fprintln(&b,"")
		cols := []string{
			row.Key,
			fmt.Sprint(row.Count),
			fmt.Sprint(row.Sum),
			fmt.Sprint(row.Min),
			fmt.Sprint(row.Max),
			fmt.Sprint(row.Avg),
			fmt.Sprint(row.Rate),
			fmt.Sprint(row.P95),
		}
		fmt.Fprint(&b, fmt.Sprintf("\"%s\"", strings.Join(cols, "\",\"")))
	}
	return slackcat.SendText(filename, b.String())
}
