/*
Copyright 2021 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import "github.com/cert-manager/helm-docgen/heuristics"

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
