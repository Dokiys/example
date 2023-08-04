package tablec

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	t.Run("测试输出", func(t *testing.T) {
		var wr = &bytes.Buffer{}
		var conf = &DBConf{
			Host:     os.Getenv("DEV_DATABASE_HOST"),
			Port:     3306,
			Schema:   "verypay_eticket",
			Username: os.Getenv("DEV_DATABASE_USERNAME"),
			Password: os.Getenv("DEV_DATABASE_PASSWORD"),
		}
		var tableName = "ticket_version"

		tablec, err := NewTableC(conf)
		assert.NoError(t, err)
		err = tablec.Do(wr, TypModel, tableName)
		assert.NoError(t, err)

		t.Log(wr.String())
	})

}
