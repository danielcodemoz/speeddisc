package report

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func fixture() Result {
	r := Result{
		Title: "Diagnóstico de espaço em disco", Status: "Concluído", Lang: "pt-PT",
		Started:  time.Date(2026, 11, 20, 19, 3, 2, 0, time.FixedZone("SAST", 2*3600)),
		Duration: 4100 * time.Millisecond,
		Rows: []Row{
			{"Unidade", `C:\ (NTFS)`},
			{"Espaço disponível", "32,7 GB (26% de 128 GB)"},
			{"Resultado", "Não foram apagados ficheiros. O relatório contém apenas a análise seleccionada."},
			{"Origem", `C:\Utilizadores\Daniel\Documentos\Português\balanço (revisto).xlsx`},
		},
	}
	for i := 1; i <= 38; i++ {
		r.Rows = append(r.Rows, Row{
			fmt.Sprintf("Ficheiro %02d", i),
			fmt.Sprintf(`C:\Dados\Trabalhos\Ano lectivo\Apresentações\Apresentação visual %02d.pptx`, i),
		})
	}
	r.Rows = append(r.Rows, Row{"Caminho extenso", "C:\\" + strings.Repeat("exemplo-", 70) + "東京.xlsx"})
	return r
}

func TestTXTUTF8AndEscaping(t *testing.T) {
	var b bytes.Buffer
	err := Write(&b, ".TXT", Result{
		Title: "Análise\nSucesso", Status: "Concluído", Lang: "pt",
		Rows: []Row{{"Caminho", `C:\東京\revisão.txt` + "\r\nsegunda linha\x00"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Relatório de operação", "Análise Sucesso", "Estado: Concluído", `C:\東京\revisão.txt` + "\r\n  segunda linha<U+0000>"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("missing %q: %q", want, b.String())
		}
	}
	b.Reset()
	if err := Write(&b, "txt", Result{Lang: "en-GB"}); err != nil || !strings.Contains(b.String(), "No additional details.") {
		t.Fatalf("English: %v %q", err, b.String())
	}
}

func TestExclusiveSave(t *testing.T) {
	file := filepath.Join(t.TempDir(), "report.txt")
	if err := Save(file, "txt", fixture()); err != nil {
		t.Fatal(err)
	}
	a, _ := os.ReadFile(file)
	if err := Save(file, "TXT", Result{Title: "another"}); !errors.Is(err, fs.ErrExist) {
		t.Fatalf("expected refusal to overwrite, got %v", err)
	}
	b, _ := os.ReadFile(file)
	if !bytes.Equal(a, b) {
		t.Error("original overwritten")
	}
	if err := Save(file+"-invalid", "html", fixture()); err == nil {
		t.Error("expected error for unsupported format")
	}
	if _, err := os.Stat(file + "-invalid"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("unknown format created file: %v", err)
	}
}

type shortWriter struct{}

func (shortWriter) Write([]byte) (int, error) { return 0, nil }

func TestWriterErrors(t *testing.T) {
	if err := Write(shortWriter{}, "txt", fixture()); !errors.Is(err, io.ErrShortWrite) {
		t.Errorf("short write: expected io.ErrShortWrite, got %v", err)
	}
	if err := Write(shortWriter{}, "exe", fixture()); err == nil {
		t.Error("expected invalid format error")
	}
}

func TestPDFPaginationAndXref(t *testing.T) {
	var b bytes.Buffer
	if err := Write(&b, "pdf", fixture()); err != nil {
		t.Fatal(err)
	}
	data := b.Bytes()
	if !bytes.HasPrefix(data, []byte("%PDF-1.4\n")) || !bytes.HasSuffix(data, []byte("%%EOF\n")) {
		t.Fatal("invalid PDF signature")
	}
	pageCount := bytes.Count(data, []byte("/Type /Page "))
	if pageCount < 2 || pageCount != bytes.Count(data, []byte("/XObject << /Logo 5 0 R")) {
		t.Fatalf("pages /logo invalid: %d", pageCount)
	}
	xref := regexp.MustCompile(`startxref\n(\d+)\n%%EOF\n$`).FindSubmatch(data)
	if len(xref) != 2 {
		t.Fatal("missing cross-reference")
	}
	position, _ := strconv.Atoi(string(xref[1]))
	if position <= 0 || position >= len(data) || !bytes.HasPrefix(data[position:], []byte("xref\n")) {
		t.Fatalf("wrong cross-reference offset %d", position)
	}
	lines := bytes.Split(data[position:], []byte("\n"))
	for id := 1; id <= 10+pageCount*2; id++ {
		offset, _ := strconv.Atoi(string(lines[2+id][:10]))
		if offset <= 0 || offset >= position || !bytes.HasPrefix(data[offset:], []byte(fmt.Sprintf("%d 0 obj\n", id))) {
			t.Fatalf("wrong object %d offset %d", id, offset)
		}
	}
	contentMatcher := regexp.MustCompile(`/Length (\d+) >>\nstream\n`)
	total := 0
	for i, index := range contentMatcher.FindAllSubmatchIndex(data, -1) {
		l, _ := strconv.Atoi(string(data[index[2]:index[3]]))
		r, err := zlib.NewReader(bytes.NewReader(data[index[1] : index[1]+l]))
		if err != nil {
			t.Fatal(err)
		}
		s, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			if len(s) != 128*128*3 {
				t.Fatalf("wrong RGB bitmap size %d", len(s))
			}
			continue
		}
		if i < 3 {
			if !bytes.HasPrefix(s, []byte{0, 1, 0, 0}) {
				t.Fatal("invalid embedded TrueType font")
			}
			continue
		}
		total++
		if !bytes.Contains(s, []byte("/Logo Do")) || !bytes.Contains(s, []byte("(CAMPO)")) || !bytes.Contains(s, []byte("(RESULTADO)")) {
			t.Fatalf("page %d missing logo or the repeatable table header", i)
		}
	}
	if total != pageCount {
		t.Fatalf("expected %d streams, got %d", pageCount, total)
	}
}

func TestWidthsCapsAndLocale(t *testing.T) {
	for _, tc := range []struct {
		s       string
		size, w float64
		bold    bool
	}{
		{strings.Repeat("W", 1500), 10, 350, false},
		{strings.Repeat("東京", 102), 10, 142, true},
		{"Línea de acción: duração e relatório", 23, 200, true},
	} {
		for _, l := range wrap(tc.s, tc.size, tc.w, tc.bold) {
			if measure(l, tc.size, tc.bold) > tc.w {
				t.Errorf("%q exceeds cell width", l)
			}
		}
	}
	rows := make([]Row, maxRows+2)
	rows[0] = Row{"Label", strings.Repeat("á", 2050)}
	r, l := normalise(Result{Rows: rows, Lang: "en"})
	if len(r.Rows) != maxRows+1 || r.Rows[maxRows] != (Row{"Omitted rows", "2"}) || !strings.HasSuffix(r.Rows[0].Value, l.truncated) {
		t.Fatal("missing bounded rows or explicit omission markers")
	}
	if encode("revisão (x) \\ 東京") != "revis\xe3o (x) \\ <U+6771><U+4EAC>" {
		t.Fatal("WinAnsi and transparent Unicode fallback diverged")
	}
}

func TestExportQAExample(t *testing.T) {
	dir := os.Getenv("SPEEDDISC_TEST_REPORT_DIR")
	if dir == "" {
		t.Skip("optional visual QA output directory not set")
	}
	for _, ext := range []string{"txt", "pdf"} {
		if err := Save(filepath.Join(dir, "speeddisc-example."+ext), ext, fixture()); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTallTitleAndContextStillFit(t *testing.T) {
	var out bytes.Buffer
	r := Result{Title: strings.Repeat("東京", 70), Status: strings.Repeat("🧪", 100), Rows: []Row{{"Value", strings.Repeat("東", 600)}}}
	if err := Write(&out, "pdf", r); err != nil {
		t.Fatal(err)
	}
	if n := bytes.Count(out.Bytes(), []byte("/Type /Page ")); n < 2 || n > 20 {
		t.Errorf("unexpected page count for worst-case header: %d", n)
	}
}
