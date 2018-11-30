package svgclean

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"testing"
)

var (
	inSVG = `<svg version="1.1" baseProfile="full" xmlns="http://www.w3.org/2000/svg">
  <polygon id="triangle" points="0,0 0,50 50,0" fill="#009900" stroke="#004400" onclick="har" poison="javascript:alert('you suck')"/>
  <!--blah-->
  <script type="text/javascript">
    alert('This app is probably vulnerable to XSS attacks!');
  </script>
  <a href="javascript:alert('crapola');">snots</a>
  <script type="text/javascript">
      alert('Again, this app is probably vulnerable to XSS attacks!');
  </script>
  <p>
    <script type="text/javascript">
      alert('Really, this app is probably vulnerable to XSS attacks!');
    </script>
  </p>
</svg>
`

	outSVG = `<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg version="1.1" baseProfile="full" xmlns="http://www.w3.org/2000/svg">
  <!--blah-->
  <polygon id="triangle" points="0,0 0,50 50,0" fill="#009900" stroke="#004400"></polygon>
  <p>&#xA;    &#xA;  </p>
</svg>`
)

func TestCleanSVGString(t *testing.T) {
	result := CleanSVGStringBlack(inSVG)

	if result != outSVG {
		t.Log(len(result), len(outSVG))
		t.Errorf("wrong:\n%s\n\n%s", result, outSVG)
	}
}

func TestCleanElementTree(t *testing.T) {
	buf := bytes.NewBuffer([]byte(inSVG))
	dec := xml.NewDecoder(buf)

	var e Element
	err := dec.Decode(&e)
	if err != nil {
		fmt.Printf("failed to unmarshal source svg: %v", err)
		os.Exit(1)
	}

	//e := (*Element)(&x)
	CleanElementTreeBlack(&e)

	b, err := xml.MarshalIndent(e, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal cleaned svg: %v", err)
		os.Exit(1)
	}

	result := svgHeader + string(b)

	if result != outSVG {
		t.Log(len(result), len(outSVG))
		t.Errorf("wrong:\n%s\n\n%s", result, outSVG)
	}
}

func TestCleanSVGStringWhite(t *testing.T) {
	result := CleanSVGStringWhite(inSVG)

	if result != outSVG {
		t.Log(len(result), len(outSVG))
		t.Errorf("wrong:\n%s\n\n%s", result, outSVG)
	}
}
