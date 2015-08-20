package path_test

import (
	"github.com/nbari/violetear/sandbox/path"
	"testing"
)

const test_path = "/0-Darwin/1-Darwin/2-Darwin/3-Darwin/4-Darwin/5-Darwin/6-Darwin/7-Darwin/8-Darwin/9-Darwin/10-Darwin/11-Darwin/12-Darwin/13-Darwin/14-Darwin/15-Darwin/16-Darwin/17-Darwin/18-Darwin/19-Darwin/20-Darwin/21-Darwin/22-Darwin/23-Darwin/24-Darwin/25-Darwin/26-Darwin/27-Darwin/28-Darwin/29-Darwin/30-Darwin/31-Darwin/32-Darwin/33-Darwin/34-Darwin/35-Darwin/36-Darwin/37-Darwin/38-Darwin/39-Darwin/40-Darwin/41-Darwin/42-Darwin/43-Darwin/44-Darwin/45-Darwin/46-Darwin/47-Darwin/48-Darwin/49-Darwin/50-Darwin/51-Darwin/52-Darwin/53-Darwin/54-Darwin/55-Darwin/56-Darwin/57-Darwin/58-Darwin/59-Darwin/60-Darwin/61-Darwin/62-Darwin/63-Darwin/64-Darwin/65-Darwin/66-Darwin/67-Darwin/68-Darwin/69-Darwin/70-Darwin/71-Darwin/72-Darwin/73-Darwin/74-Darwin/75-Darwin/76-Darwin/77-Darwin/78-Darwin/79-Darwin/80-Darwin/81-Darwin/82-Darwin/83-Darwin/84-Darwin/85-Darwin/86-Darwin/87-Darwin/88-Darwin/89-Darwin/90-Darwin/91-Darwin/92-Darwin/93-Darwin/94-Darwin/95-Darwin/96-Darwin/97-Darwin/98-Darwin/99-Darwin/100-Darwin/"

func BenchmarkA(b *testing.B) {
	path.A(test_path)
}

func BenchmarkB(b *testing.B) {
	path.B(test_path)
}
