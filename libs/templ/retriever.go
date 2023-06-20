package templ

type Retriever interface {
	Retrieve(viewTag ViewTag) (CacheableContent, error)
}

type CacheableContent struct {
	Content string
	Tags    []string
}
