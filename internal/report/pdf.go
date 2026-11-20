package report

import (
	"bytes"
	"compress/zlib"
	_ "embed"
	"fmt"
	"image/png"
	"strings"
	"sync"
)

// The approved SpeedDisc site favicon, 256px master downscaled to 128px
// and composited on the report masthead colour to retain its native shape.
//
//go:embed brand.png
var brandPNG []byte

// The Website's Outfit Latin fonts, turned into static, renamed subsets.
// The fonts retain their original copyright and the repository ships the OFL.
//
//go:embed regular.ttf
var regularFont []byte

//go:embed bold.ttf
var boldFont []byte

type fontSpec struct {
	name                       string
	box                        [4]int
	ascent, descent, cap, stem int
}

// SpeedDisc's original petrol and mint on a light neutral sheet.
const (
	paper      = 0xf5f7f4
	white      = 0xffffff
	forest     = 0x142c27
	deep       = 0x34554b
	ink        = 0x192e29
	muted      = 0x4a635b
	petrol     = 0x0e6b62
	mint       = 0xa8e4d7
	teal       = 0x08776e
	line       = 0xd2e0da
	stripe     = 0xeaf3ef
	capGreen   = 0x66bcaa
	mastMuted  = 0xbad8ce
	pageW      = 595.28
	pageH      = 841.89
	left       = 41.0
	tableW     = pageW - 2*left
	fieldW     = 161.0
	bodyBottom = 65.0
)

var winAnsi = map[rune]byte{
	'€': 128, '‚': 130, 'ƒ': 131, '„': 132, '…': 133, '†': 134, '‡': 135,
	'ˆ': 136, '‰': 137, 'Š': 138, '‹': 139, 'Œ': 140, 'Ž': 142,
	'‘': 145, '’': 146, '“': 147, '”': 148, '•': 149, '–': 150, '—': 151,
	'˜': 152, '™': 153, 'š': 154, '›': 155, 'œ': 156, 'ž': 158, 'Ÿ': 159,
}

func encode(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '\n' {
			b.WriteByte('\n')
			continue
		}
		var c byte
		if r >= 32 && r <= 126 || r >= 160 && r <= 255 {
			c = byte(r)
		} else {
			c = winAnsi[r]
		}
		if c != 0 && regularWidths[c] != 0 && boldWidths[c] != 0 {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "<U+%04X>", r)
		}
	}
	return b.String()
}

func pdfString(s string) string {
	var b strings.Builder
	b.WriteByte('(')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '(' || c == ')' || c == '\\':
			b.WriteByte('\\')
			b.WriteByte(c)
		case c < 32 || c >= 127:
			fmt.Fprintf(&b, "\\%03o", c)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte(')')
	return b.String()
}

func measure(s string, size float64, bold bool) float64 {
	widths := &regularWidths
	if bold {
		widths = &boldWidths
	}
	var total int
	for i := 0; i < len(s); i++ {
		total += int(widths[s[i]])
	}
	return float64(total) * size / 1000
}

// wrap accepts UTF-8, returning lines already converted to PDF WinAnsi;
// split long uninterrupted paths too, without adding spurious hyphens.
func wrap(s string, size, max float64, bold bool) []string {
	var out []string
	widths := &regularWidths
	if bold {
		widths = &boldWidths
	}
	limit := int(max * 1000 / size)
	for _, source := range strings.Split(encode(s), "\n") {
		source = strings.TrimRight(source, " ")
		if source == "" {
			out = append(out, "")
			continue
		}
		for source != "" {
			n, used := 0, 0
			for n < len(source) && (n == 0 || used+int(widths[source[n]]) <= limit) {
				used += int(widths[source[n]])
				n++
			}
			if n == len(source) {
				out = append(out, source)
				break
			}
			if i := strings.LastIndexByte(source[:n], ' '); i > 0 && i >= n/2 {
				n = i
			}
			// Keep the escape for unsupported Unicode intact at a wrap boundary.
			if i := strings.LastIndex(source[:n], "<"); i > 0 && n-i < 10 &&
				strings.HasPrefix(source[i:], "<U+") && !strings.Contains(source[i:n], ">") {
				n = i
			}
			out = append(out, strings.TrimRight(source[:n], " "))
			source = strings.TrimLeft(source[n:], " ")
		}
	}
	return out
}

// The fixed-height title/context regions cannot consume the full body.
// Full labels and data are still available in the text export.
func limitLines(lines []string, n int, maxWidth, size float64, bold bool) []string {
	if len(lines) <= n {
		return lines
	}
	lines = lines[:n]
	s := lines[n-1]
	dots := encode("…")
	for s != "" && measure(s+dots, size, bold) > maxWidth {
		s = s[:len(s)-1]
	}
	lines[n-1] = strings.TrimRight(s, " ") + dots
	return lines
}

type canvas struct{ bytes.Buffer }

func (c *canvas) fill(col int) {
	fmt.Fprintf(c, "%.4f %.4f %.4f rg\n", float64(col>>16&255)/255, float64(col>>8&255)/255, float64(col&255)/255)
}

func (c *canvas) box(x, y, w, h float64, col int) {
	c.fill(col)
	fmt.Fprintf(c, "%.2f %.2f %.2f %.2f re f\n", x, y, w, h)
}

func (c *canvas) text(x, y, size float64, col int, bold bool, s string) {
	font := "/F1"
	if bold {
		font = "/F2"
	}
	c.fill(col)
	fmt.Fprintf(c, "BT %s %.2f Tf 1 0 0 1 %.2f %.2f Tm %s Tj ET\n", font, size, x, y, pdfString(s))
}

func startPage(r Result, l labels) (*canvas, float64) {
	c := &canvas{}
	c.box(0, 0, pageW, pageH, paper)
	c.box(0, pageH-9, pageW*.62, 9, capGreen)
	c.box(pageW*.62, pageH-9, pageW*.38, 9, teal)
	top := pageH - 34
	title := limitLines(wrap(r.Title, 23, tableW-44, true), 4, tableW-44, 23, true)
	headerH := 94.0 + float64(len(title))*27
	base := top - headerH
	c.box(left, base, tableW, headerH, forest)
	c.box(left, base, 4, headerH, capGreen)
	fmt.Fprintf(c, "q 41 0 0 41 %.2f %.2f cm /Logo Do Q\n", left+20, top-62)
	c.text(left+73, top-38, 17, white, true, "SpeedDisc")
	c.text(left+73, top-54, 9.2, mint, true, encode(strings.ToUpper(l.report)))
	c.text(left+tableW-149, top-39, 9, mastMuted, false, "speeddisc.danielpro.dev")
	c.box(left+22, top-75, tableW-44, .7, deep)
	for i, s := range title {
		c.text(left+22, top-102-float64(i)*27, 23, white, true, s)
	}

	// The three context cells share a height sized to the longest value.
	columns := []float64{tableW * .38, tableW * .37, tableW * .25}
	info := metadata(r, l)
	rowLines := make([][]string, len(info))
	need := 1
	for i, row := range info {
		space := columns[i] - 24
		rowLines[i] = limitLines(wrap(row.Value, 10.5, space, true), 3, space, 10.5, true)
		need = max(need, len(rowLines[i]))
	}
	cardTop := base - 12
	cardH := 40 + float64(need-1)*13
	x := left
	for i, row := range info {
		y := cardTop - cardH
		c.box(x, y, columns[i]-5, cardH, line)
		c.box(x+.7, y+.7, columns[i]-6.4, cardH-1.4, white)
		c.text(x+10, cardTop-15, 8.5, muted, true, encode(strings.ToUpper(row.Label)))
		for j, s := range rowLines[i] {
			c.text(x+10, cardTop-30-float64(j)*13, 10.5, ink, true, s)
		}
		x += columns[i]
	}

	y := cardTop - cardH - 16
	c.box(left, y-27, tableW, 27, petrol)
	c.text(left+12, y-17, 9.4, white, true, encode(strings.ToUpper(l.field)))
	c.text(left+fieldW+12, y-17, 9.4, white, true, encode(strings.ToUpper(l.value)))
	return c, y - 27
}

func pageFooter(c *canvas, l labels, number, total int) {
	c.box(left, 49, tableW, .8, line)
	c.text(left, 33, 8.5, muted, true, "SPEEDDISC  /  "+encode(strings.ToUpper(l.report)))
	pages := encode(fmt.Sprintf("%s %02d / %02d", l.page, number, total))
	c.text(pageW-left-measure(pages, 8.5, true), 33, 8.5, muted, true, pages)
}

const rowLead = 14.0

func drawRow(c *canvas, y float64, num, start, count int, field, value []string) float64 {
	height := float64(count)*rowLead + 12
	colour := white
	if num%2 != 0 {
		colour = stripe
	}
	c.box(left, y-height, tableW, height, colour)
	c.box(left, y-height, .7, height, line)
	c.box(left+fieldW, y-height, .7, height, line)
	c.box(left+tableW-.7, y-height, .7, height, line)
	c.box(left, y-height, tableW, .7, line)
	for k := 0; k < count; k++ {
		i := start + k
		baseline := y - 16 - float64(k)*rowLead
		if i < len(field) {
			c.text(left+12, baseline, 10, ink, true, field[i])
		}
		if i < len(value) {
			c.text(left+fieldW+12, baseline, 10, ink, false, value[i])
		}
	}
	return y - height
}

func pdf(r Result, l labels) ([]byte, error) {
	logo, err := brand()
	if err != nil {
		return nil, err
	}
	var pages []*canvas
	c, y := startPage(r, l)
	fullLines := int((y - bodyBottom - 12) / rowLead)
	pages = append(pages, c)
	newPage := func() { c, y = startPage(r, l); pages = append(pages, c) }
	for n, row := range r.Rows {
		field := wrap(row.Label, 10, fieldW-26, true)
		value := wrap(row.Value, 10, tableW-fieldW-26, false)
		length := max(len(field), len(value))
		for i := 0; i < length; {
			free := int((y - bodyBottom - 12) / rowLead)
			if free < 1 || i == 0 && length <= fullLines && length > free {
				newPage()
				free = fullLines
			}
			consume := min(free, length-i)
			y = drawRow(c, y, n, i, consume, field, value)
			i += consume
		}
	}

	fonts, err := fontObjects()
	if err != nil {
		return nil, err
	}
	objects := [][]byte{nil, nil,
		fonts[0], fonts[1],
		logo,
		[]byte("<< /Title " + pdfString(encode("SpeedDisc | "+l.report)) + " /Creator (SpeedDisc) >>"),
		fonts[2], fonts[3], fonts[4], fonts[5],
	}
	var kids strings.Builder
	for i, c := range pages {
		pageFooter(c, l, i+1, len(pages))
		pageID := len(objects) + 1
		stream, err := deflate(c.Bytes())
		if err != nil {
			return nil, err
		}
		objects = append(objects,
			[]byte(fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.2f %.2f] /Resources << /Font << /F1 3 0 R /F2 4 0 R >> /XObject << /Logo 5 0 R >> >> /Contents %d 0 R >>", pageW, pageH, pageID+1)),
			streamObject("", stream))
		fmt.Fprintf(&kids, "%d 0 R ", pageID)
	}
	objects[0] = []byte(fmt.Sprintf("<< /Type /Catalog /Pages 2 0 R /Lang %s >>", pdfString(r.Lang)))
	objects[1] = []byte(fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", kids.String(), len(pages)))
	return file(objects), nil
}

func streamObject(header string, data []byte) []byte {
	var dst bytes.Buffer
	fmt.Fprintf(&dst, "<< %s/Filter /FlateDecode /Length %d >>\nstream\n", header, len(data))
	dst.Write(data)
	dst.WriteString("\nendstream")
	return dst.Bytes()
}

func deflate(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		_ = w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

var brand = sync.OnceValues(func() ([]byte, error) {
	im, err := png.Decode(bytes.NewReader(brandPNG))
	if err != nil {
		return nil, err
	}
	b := im.Bounds()
	pixels := make([]byte, 0, 3*b.Dx()*b.Dy())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := im.At(x, y).RGBA()
			// Recolour only the old report masthead matte, not the site logo.
			if r>>8 == 0x16 && g>>8 == 0x1b && b>>8 == 0x2e {
				pixels = append(pixels, byte(forest>>16), byte(forest>>8&255), byte(forest&255))
			} else {
				pixels = append(pixels, byte(r>>8), byte(g>>8), byte(b>>8))
			}
		}
	}
	rgb, err := deflate(pixels)
	if err != nil {
		return nil, err
	}
	header := fmt.Sprintf("/Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 ", b.Dx(), b.Dy())
	return streamObject(header, rgb), nil
})

func fontDefinition(f fontSpec, widths [256]uint16, descriptor int) []byte {
	var w strings.Builder
	for _, width := range widths[32:] {
		fmt.Fprintf(&w, "%d ", width)
	}
	return []byte(fmt.Sprintf("<< /Type /Font /Subtype /TrueType /BaseFont /%s /FirstChar 32 /LastChar 255 /Widths [%s] /FontDescriptor %d 0 R /Encoding /WinAnsiEncoding >>", f.name, w.String(), descriptor))
}

func describeFont(f fontSpec, fileID int) []byte {
	return []byte(fmt.Sprintf("<< /Type /FontDescriptor /FontName /%s /Flags 32 /FontBBox [%d %d %d %d] /ItalicAngle 0 /Ascent %d /Descent %d /CapHeight %d /StemV %d /FontFile2 %d 0 R >>", f.name, f.box[0], f.box[1], f.box[2], f.box[3], f.ascent, f.descent, f.cap, f.stem, fileID))
}

var fontObjects = sync.OnceValues(func() ([][]byte, error) {
	reg, err := deflate(regularFont)
	if err != nil {
		return nil, err
	}
	bold, err := deflate(boldFont)
	if err != nil {
		return nil, err
	}
	return [][]byte{
		fontDefinition(regularSpec, regularWidths, 7), fontDefinition(boldSpec, boldWidths, 9),
		describeFont(regularSpec, 8), streamObject(fmt.Sprintf("/Length1 %d ", len(regularFont)), reg),
		describeFont(boldSpec, 10), streamObject(fmt.Sprintf("/Length1 %d ", len(boldFont)), bold),
	}, nil
})

func file(objs [][]byte) []byte {
	var out bytes.Buffer
	out.WriteString("%PDF-1.4\n%\xE2\xE3\xCF\xD3\n")
	locs := make([]int, len(objs)+1)
	for i, ob := range objs {
		locs[i+1] = out.Len()
		fmt.Fprintf(&out, "%d 0 obj\n", i+1)
		out.Write(ob)
		out.WriteString("\nendobj\n")
	}
	start := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n0000000000 65535 f \n", len(locs))
	for _, off := range locs[1:] {
		fmt.Fprintf(&out, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(&out, "trailer\n<< /Size %d /Root 1 0 R /Info 6 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(locs), start)
	return out.Bytes()
}
