package parser

import "testing"

func TestMarkdownParserSatisfiesParserInterface(t *testing.T) {
	var _ Parser = MarkdownParser{}
}
