package templ

type Retriever interface {
	Retrieve(viewTag ViewTag) (string, error)
}
