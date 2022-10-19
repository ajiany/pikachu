package httpx

type BeforeRequestHook func(opts *Option)
type RequestHook func(req *Request, opts *Option)
type ResponseHook func(resp *Response) error
