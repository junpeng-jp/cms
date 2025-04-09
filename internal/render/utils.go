package render

import (
	"fmt"
	"io"
	"strings"
)

const (
	tagDiv        = "div"
	tagImage      = "img"
	tagBlockquote = "blockquote"
	tagParagraph  = "p"
	// inline formatting
	tagSpan          = "span"
	tagUnderline     = "u"
	tagItalic        = "i"
	tagBold          = "b"
	tagStrikethrough = "s"
)

func writeOpeningTag(writer io.Writer, tag string, depth int, annotations ...string) (n int, err error) {
	if len(annotations) != 0 {
		n, err = writer.Write([]byte(fmt.Sprintf("%s<%s %s>\n", strings.Repeat(tab, depth), tag, strings.Join(annotations, " "))))
	} else {
		n, err = writer.Write([]byte(fmt.Sprintf("%s<%s>\n", strings.Repeat(tab, depth), tag)))
	}
	return n, err
}

func writeClosingTag(writer io.Writer, tag string, depth int, annotations ...string) (n int, err error) {
	if len(annotations) != 0 {
		n, err = writer.Write([]byte(fmt.Sprintf("%s</%s %s>\n", strings.Repeat(tab, depth), tag, strings.Join(annotations, " "))))
	} else {
		n, err = writer.Write([]byte(fmt.Sprintf("%s</%s>\n", strings.Repeat(tab, depth), tag)))
	}
	return n, err
}

func writeContent(writer io.Writer, content string, depth int) (n int, err error) {
	n, err = writer.Write([]byte(fmt.Sprintf("%s%s\n", strings.Repeat(tab, depth), content)))
	return n, err
}

func createClassAnnotation(classes ...string) string {
	return fmt.Sprintf("class=\"%s\"", strings.Join(classes, " "))
}
