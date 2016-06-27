package clickhouse

import (
	"strings"
	"bytes"
	"net/http"
)

const (
	httpTransportBodyType = "text/plain"
)

type Transport interface {
	Exec(conn *Conn, q Query, readOnly bool) (res string, err error)
}

type HttpTransport struct{}

func (t HttpTransport) Exec(conn *Conn, q Query, readOnly bool) (res string, err error) {
	var resp *http.Response
	query := prepareHttp(q.Stmt, q.args)
	if readOnly {
		if len(query) > 0 {
			query = "?query="+query
		}
		resp, err = http.Get(conn.Host + query)
	} else {
		resp, err = http.Post(conn.Host, httpTransportBodyType, strings.NewReader(query))
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	return buf.String(), err
}

func prepareHttp(stmt string, args []interface{}) (res string) {
	res = stmt
	for _, arg := range args {
		res = strings.Replace(res, "?", marshal(arg), 1)
	}

	return res
}