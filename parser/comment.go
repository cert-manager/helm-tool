package parser

import "github.com/thatsmrtalbot/helm-docgen/heuristics"

type Comment struct {
	heuristics.CommentBlock
	Tags tags
}

func parseComments(comment string) (comments []Comment) {
	blocks := heuristics.ParseCommentIntoBlocks(comment)
	for _, block := range blocks {
		c := Comment{CommentBlock: block}
		for _, segment := range block.Segments {
			if segment.Type == heuristics.ContentTypeTag {
				c.Tags.Push(segment.Contents[0])
			}
		}

		comments = append(comments, c)
	}

	return
}
