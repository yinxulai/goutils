package service

// Markdown 操作类
import "errors"

// Markdown Markdown
var Markdown *MarkdownService
var markDownInited bool

// MarkdownService Markdown 相关服务
type MarkdownService struct {
}

// Init 初始化
func (ts *MarkdownService) init() (err error) {
	if markDownInited {
		return errors.New("MarkdownService is inited! ")
	}
	Markdown = ts
	markDownInited = true
	return err
}

// Info 基本信息
func (ts *MarkdownService) info() BaseServiceInfo {
	return BaseServiceInfo{
		key: "MarkdownService",
	}
}

// MarkdownToHTML  转换 markdown 为 html
func (ts *MarkdownService) MarkdownToHTML() {

}

// HTMLToMarkdown 转换 html 为 markdown
func (ts *MarkdownService) HTMLToMarkdown() {

}
