package templ

type Retriever interface {
	Retrieve(viewTag ViewTag) (CacheableContent, error)
}

type CacheableContent interface {
	Render() string
	GetContext() []any
	GetTags() []string
}

func NewCacheableContent(content string, tags []string, context ...any) CacheableContent {
	return &cacheableContent{
		content: content,
		tags:    tags,
		context: context,
	}
}

type cacheableContent struct {
	content string
	tags    []string
	context []any
}

func (cc cacheableContent) Render() string {
	return cc.content
}

func (cc cacheableContent) GetTags() []string {
	return cc.tags
}

func (cc cacheableContent) GetContext() []any {
	return cc.context
}
