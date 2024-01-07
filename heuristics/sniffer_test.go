package heuristics

import (
	"reflect"
	"strings"
	"testing"
)

func TestContentSniffer_SniffContentType(t *testing.T) {
	type fields struct {
		previousLineBuffer []string
		currentType        ContentType
	}
	type args struct {
		line string
	}
	type want struct {
		contentType  ContentType
		startOfBlock bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			"ContentTypeYaml/NoValue",
			fields{},
			args{
				`foo:`,
			},
			want{
				ContentTypeYaml,
				true,
			},
		},
		{
			"ContentTypeYaml/SimpleStringValue",
			fields{},
			args{
				`foo: bar`,
			},
			want{
				ContentTypeYaml,
				true,
			},
		},
		{
			"ContentTypeYaml/EmptyObjectValue",
			fields{},
			args{
				`foo: {}`,
			},
			want{
				ContentTypeYaml,
				true,
			},
		},
		{
			"ContentTypeYaml/EmptyArrayValue",
			fields{},
			args{
				`foo: []`,
			},
			want{
				ContentTypeYaml,
				true,
			},
		},
		{
			"ContentTypeYaml/QuotedStringValueWithSpace",
			fields{},
			args{
				`foo: "has a space"`,
			},
			want{
				ContentTypeYaml,
				true,
			},
		},
		{
			"ContentTypeYaml/ContinuingFromYaml",
			fields{
				currentType: ContentTypeYaml,
				previousLineBuffer: []string{
					`foo:`,
				},
			},
			args{
				`  bar: baz`,
			},
			want{
				ContentTypeYaml,
				false,
			},
		},
		{
			"ContentTypeText/ContinuingFromYaml",
			fields{
				currentType: ContentTypeYaml,
				previousLineBuffer: []string{
					`foo: bar`,
				},
			},
			args{
				`The quick brown fox jumps over the lazy dog`,
			},
			want{
				ContentTypeText,
				true,
			},
		},
		{
			"ContentTypeText/Yaml/UnquotedStringValueWithSpace",
			fields{},
			args{
				"foo: has a space",
			},
			want{
				ContentTypeText,
				true,
			},
		},
		{
			"ContentTypeText/Yaml/KeyWithSpace",
			fields{},
			args{
				`foo bar: baz`,
			},
			want{
				ContentTypeText,
				true,
			},
		},
		{
			"ContentTypeText/Yaml/KeyStartingWithCapitalLetter",
			fields{},
			args{
				`Foo: bar`,
			},
			want{
				ContentTypeText,
				true,
			},
		},
		{
			"ContentTypeText/Basic",
			fields{},
			args{
				`The quick brown fox jumps over the lazy dog`,
			},
			want{
				ContentTypeText,
				true,
			},
		},
		{
			"ContentTypeTag/Basic",
			fields{},
			args{
				`+docs:foo`,
			},
			want{
				ContentTypeTag,
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ContentSniffer{
				previousLineBuffer: tt.fields.previousLineBuffer,
				currentType:        tt.fields.currentType,
			}
			got, got1 := c.SniffContentType(tt.args.line)
			if got != tt.want.contentType {
				t.Errorf("ContentSniffer.SniffContentType() got = %v, want %v", got, tt.want.contentType)
			}
			if got1 != tt.want.startOfBlock {
				t.Errorf("ContentSniffer.SniffContentType() got1 = %v, want %v", got1, tt.want.startOfBlock)
			}
		})
	}
}

func TestParseCommentIntoBlocks(t *testing.T) {
	type args struct {
		comment string
	}
	type want struct {
		commentBlocks []CommentBlock
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"SingleLineTextComment",
			args{
				comment: `# This is an example single line comment`,
			},
			want{
				[]CommentBlock{
					{
						Segments: []CommentBlockSegment{
							{
								Type:     ContentTypeText,
								Contents: []string{" This is an example single line comment"},
							},
						},
					},
				},
			},
		},
		{
			"MultiLineTextComment",
			args{
				comment: strings.Join([]string{
					`# This is an example multi-line`,
					`# comment`,
				}, "\n"),
			},
			want{
				[]CommentBlock{
					{
						Segments: []CommentBlockSegment{
							{
								Type:     ContentTypeText,
								Contents: []string{" This is an example multi-line", " comment"},
							},
						},
					},
				},
			},
		},
		{
			"MultiLineMixedTypeComment",
			args{
				comment: strings.Join([]string{
					`# This is an example multi-line`,
					`# comment followed by some yaml`,
					`# foo:`,
					`#   bar: baz`,
				}, "\n"),
			},
			want{
				[]CommentBlock{
					{
						Segments: []CommentBlockSegment{
							{
								Type:     ContentTypeText,
								Contents: []string{" This is an example multi-line", " comment followed by some yaml"},
							},
							{
								Type:     ContentTypeYaml,
								Contents: []string{" foo:", "   bar: baz"},
							},
						},
					},
				},
			},
		},
		{
			"MultipleBlocks",
			args{
				comment: strings.Join([]string{
					`# This is an example multi-line`,
					`# comment followed by some yaml`,
					`# foo:`,
					`#   bar: baz`,
					``,
					`# This is another comment`,
				}, "\n"),
			},
			want{
				[]CommentBlock{
					{
						Segments: []CommentBlockSegment{
							{
								Type:     ContentTypeText,
								Contents: []string{" This is an example multi-line", " comment followed by some yaml"},
							},
							{
								Type:     ContentTypeYaml,
								Contents: []string{" foo:", "   bar: baz"},
							},
						},
					},
					{
						Segments: []CommentBlockSegment{
							{
								Type:     ContentTypeText,
								Contents: []string{" This is another comment"},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCommentIntoBlocks(tt.args.comment); !reflect.DeepEqual(got, tt.want.commentBlocks) {
				t.Errorf("ParseCommentIntoBlocks() = %v, want %v", got, tt.want.commentBlocks)
			}
		})
	}
}
