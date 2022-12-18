package enums

type TriggerType string

type HttpMethod string

const (
	Http  TriggerType = "http"
	Event TriggerType = "event"
)

const (
	HttpMethodGet    HttpMethod = "get"
	HttpMethodPost   HttpMethod = "post"
	HttpMethodPut    HttpMethod = "put"
	HttpMethodPatch  HttpMethod = "patch"
	HttpMethodDelete HttpMethod = "delete"
)

func HttpMethodToString(m HttpMethod) string {
	if m == HttpMethodGet {
		return "GET"
	} else if m == HttpMethodPost {
		return "POST"
	} else if m == HttpMethodPut {
		return "PUT"
	} else if m == HttpMethodPatch {
		return "PATCH"
	} else if m == HttpMethodDelete {
		return "DELETE"
	} else {
		panic("http method not supported" + m)
	}
}
