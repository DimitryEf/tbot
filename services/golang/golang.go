package golang

import (
	"fmt"
	"gorm.io/gorm"
	"sort"
	"strings"
	"tbot/config"
	"tbot/internal/errors"
)

type Golang struct {
	golangStg *config.GolangSettings
	db        *gorm.DB
	ready     bool
}

func NewGolang(golangStg *config.GolangSettings, db *gorm.DB) *Golang {

	initialize(db)

	return &Golang{
		golangStg: golangStg,
		db:        db,
		ready:     true,
	}
}

func initialize(db *gorm.DB) {
	err := db.Debug().AutoMigrate(&Topic{}, &Tag{})
	errors.PanicIfErr(err)

	create(db, "Get executable dir\n(tags: executable dir)\n---\n\nex, err := os.Executable()\ndir := filepath.Dir(ex)\nfmt.Println(\"dir:\", dir)\n")
	create(db, "Extract beginning of string (prefix)\n(tags: extract beginning string prefix)\n---\n\nt := string([]rune(s)[:5])")
	create(db, "Extract string suffix\n(tags: extract string suffix)\n---\n\nt := string([]rune(s)[len([]rune(s))-5:])")
	create(db, "Exec other program\n(tags: exec program)\n---\n\nerr := exec.Command(\"program\", \"arg1\", \"arg2\").Run()")
	create(db, "Telegram message markdown\n(tags: telegram message markdown)\n---\n\n*полужирный*\n_курсив_\n[ссылка](http://www.example.com/)\n'строчный моноширинный'\n'''text\nблочный моноширинный (можно писать код)\n'''\n\nimport \"github.com/go-telegram-bot-api/telegram-bot-api\"\n\nmsg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)\nmsg.ParseMode = \"markdown\" //msg.ParseMode = tgbotapi.ModeMarkdown")
	create(db, "Telegram message html\n(tags: telegram message html)\n---\n\n<b>полужирный</b>, <strong>полужирный</strong>\n<i>курсив</i>\n<a href=\"http://www.example.com/\">ссылка</a>\n<code>строчный моноширинный</code>\n<pre>блочный моноширинный (можно писать код)</pre>\n\nimport \"github.com/go-telegram-bot-api/telegram-bot-api\"\n\nmsg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)\nmsg.ParseMode = \"HTML\" //msg.ParseMode = tgbotapi.ModeHTML")
	create(db, "Iterate over map entries ordered by keys\n(tags: iterate map order key)\n---\n\nkeys := make([]string, 0, len(mymap))\nfor k := range mymap {\n    keys = append(keys, k)\n}\nsort.Strings(keys)\nfor _, k := range keys {\n    x := mymap[k]\n    fmt.Println(\"Key =\", k, \", Value =\", x)\n}\n")
	create(db, "Iterate over map entries ordered by values\n(tags: iterate map order value)\n---\n\ntype entry struct {\n    key   string\n    value int\n}\nentries := make([]entry, 0, len(mymap))\nfor k, x := range mymap {\n    entries = append(entries, entry{key: k, value: x})\n}\nsort.Slice(entries, func(i, j int) bool {\n    return entries[i].value < entries[j].value\n})\nfor _, e := range entries {\n    fmt.Println(\"Key =\", e.key, \", Value =\", e.value)\n}")
	create(db, "Slice to set\n(tags: slice set)\n---\n\ny := make(map[T]struct{}, len(x))\nfor _, v := range x {\n    y[v] = struct{}{}\n}")
	create(db, "Deduplicate slice\n(tags: deduplicate slice remove duplicate)\n---\n\nseen := make(map[T]bool)\nj := 0\nfor _, v := range x {\n    if !seen[v] {\n        x[j] = v\n        j++\n        seen[v] = true\n    }\n}\nfor i := j; i < len(x); i++ {\n    x[i] = nil\n}\nx = x[:j]")
	create(db, "Shuffle a slice\n(tags: slice shuffle)\n---\n\ny := make(map[T]struct{}, len(x))\nfor _, v := range x {\n    y[v] = struct{}{}\n}")
	create(db, "Sort slice asc\n(tags: sort slice asc)\n---\n\nsort.Slice(items, func(i, j int) bool {\n    return items[i].p < items[j].p\n})")
	create(db, "Sort slice desc\n(tags: sort slice desc)\n---\n\nsort.Slice(items, func(i, j int) bool {\n    return items[i].p > items[j].p\n})")
	create(db, "Remove item from slice by index\n(tags: remove item slice index)\n---\n\nitems = append(items[:i], items[i+1:]...)")
	create(db, "Graph with adjacency lists\n(tags: graph struct)\n---\n\ntype Vertex struct{\n    Id int\n    Label string\n    Neighbours map[*Vertex]bool\n}\ntype Graph []*Vertex")
	create(db, "Reverse a string\n(tags: string reverse)\n---\n\nrunes := []rune(s)\nfor i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {\n   runes[i], runes[j] = runes[j], runes[i]\n}\nt := string(runes)")
	create(db, "Insert item in slice\n(tags: insert item slice)\n---\n\ns = append(s, 0)\ncopy(s[i+1:], s[i:])\ns[i] = x")
	create(db, "Filter slice\n(tags: filter slice)\n---\n\ny := make([]T, 0, len(x))\nfor _, v := range x{\n    if p(v){\n        y = append(y, v)\n    }\n}")
	create(db, "File content to string\n(tags: file content string)\n---\n\nb, err := ioutil.ReadFile(f)\nlines := string(b)")
	create(db, "Write to std error\n(tags: write std error)\n---\n\nfmt.Fprintln(os.Stderr, x, \"is negative\")")
	create(db, "Big int\n(tags: big integer)\n---\n\nx := new(big.Int)\nx.Exp(big.NewInt(3), big.NewInt(247), nil)")
	create(db, "Round float to int\n(tags: round float int)\n---\n\ny := int(math.Floor(x + 0.5))")
	create(db, "Check if int addition will overflow\n(tags: check int add overflow)\n---\n\nfunc willAddOverflow(a, b int64) bool {\n    return a > math.MaxInt64 - b\n}")
	create(db, "Check if int multiplication will overflow\n(tags: check int multiply overflow)\n---\n\nfunc multiplyWillOverflow(x, y uint64) bool {\n   if x <= 1 || y <= 1 {\n     return false\n   }\n   d := x * y\n   return d/y != x\n}")
	create(db, "Load json file into struct\n(tags: load json file struct)\n---\n\nbuffer, err := ioutil.ReadFile(\"data.json\")\nerr = json.Unmarshal(buffer, &x)")
	create(db, "Load yaml file into struct\n(tags: load yaml file struct)\n---\n\nimport \"gopkg.in/yaml.v3\"\n\nbuffer, err := ioutil.ReadFile(\"data.yaml\")\nerr = yaml.Unmarshal(buffer, &x)")
	create(db, "Save struct into json file\n(tags: save struct json file)\n---\n\nbuffer, err := json.MarshalIndent(x, \"\", \"  \")\nerr = ioutil.WriteFile(\"data.json\", buffer, 0644)")
	create(db, "Print type of variable\n(tags: print type variable)\n---\n\nfmt.Printf(\"%T\", x) //fmt.Println(reflect.TypeOf(x))")
	create(db, "Load from HTTP GET request into a string\n(tags: load http get string)\n---\n\nres, err := http.Get(u)\nbuffer, err := ioutil.ReadAll(res.Body)\nres.Body.Close()\ns := string(buffer)")
	create(db, "Read int from stdin\n(tags: read int std in)\n---\n\n_, err := fmt.Scan(&n)")
	create(db, "UDP listen and read\n(tags: udp listen read)\n---\n\nServerAddr,err := net.ResolveUDPAddr(\"udp\",p)\nServerConn, err := net.ListenUDP(\"udp\", ServerAddr)\ndefer ServerConn.Close()\nn,addr,err := ServerConn.ReadFromUDP(b[:1024])\nif n<1024 {\n    return fmt.Errorf(\"Only %d bytes could be read.\", n)\n}")
	create(db, "Binary search in sorted slice\n(tags: binary search slice)\n---\n\nfunc binarySearch(a []T, x T) int {\n    imin, imax := 0, len(a)-1\n    for imin <= imax {\n        imid := (imin + imax) / 2\n        switch {\n        case a[imid] == x:\n        return imid\n        case a[imid] < x:\n        imin = imid + 1\n        default:\n        imax = imid - 1\n        }\n    }\n    return -1\n}")
	create(db, "Measure func call duration\n(tags: measure func call duration time)\n---\n\nt1 := time.Now()\nfoo()\nt := time.Since(t1)\nns := t.Nanoseconds()\nfmt.Printf(\"%dns\\n\", ns)")
	create(db, "Breadth-first traversing in a graph\n(tags: bfs traversing graph)\n---\n\nfunc (start *Vertex) Bfs(f func(*Vertex)) {\n    queue := []*Vertex{start}\n    seen := map[*Vertex]bool{start: true}\n    for len(queue) > 0 {\n        v := queue[0]\n        queue = queue[1:]\n        f(v)\n        for next, isEdge := range v.Neighbours {\n            if isEdge && !seen[next] {\n                queue = append(queue, next)\n                seen[next] = true\n            }\n        }\n    }\n}")
	create(db, "Depth-first traversing in a graph\n(tags: dfs traversing graph)\n---\n\nfunc (v *Vertex) Dfs(f func(*Vertex), seen map[*Vertex]bool) {\n    seen[v] = true\n    f(v)\n    for next, isEdge := range v.Neighbours {\n        if isEdge && !seen[next] {\n            next.Dfs(f, seen)\n        }\n    }\n}")
	create(db, "Check if string contains only digits\n(tags: check string contains only digits)\n---\n\nisNotDigit := func(c rune) bool { return c < '0' || c > '9' }\nb := strings.IndexFunc(s, isNotDigit) == -1")
	create(db, "Check if file exists\n(tags: check file exist)\n---\n\n_, err := os.Stat(fp)\nb := !os.IsNotExist(err)")
	create(db, "Read slice of int from stdin\n(tags: read slice int std in)\n---\n\nvar ints []int\ns := bufio.NewScanner(os.Stdin)\nfor s.Scan() {\n    i, err := strconv.Atoi(s.Text())\n    if err == nil {\n        ints = append(ints, i)\n    }\n}")
	create(db, "Detect if 32-bit or 64-bit architecture\n(tags: detect 32 64 architecture)\n---\n\nif strconv.IntSize==32 {\n    f32()\n}\nif strconv.IntSize==64 {\n    f64()\n}")
	create(db, "Parse flags\n(tags: parse flags args)\n---\n\nvar b = flag.Bool(\"b\", false, \"Do bat\")\nfunc main() {\n    flag.Parse()\n    if *b {\n        bar()\n    }\n}")
	create(db, "Open URL in default browser\n(tags: open url default browser)\n---\n\nfunc openbrowser(url string) {\n    var err error\n    switch runtime.GOOS {\n    case \"linux\":\n        err = exec.Command(\"xdg-open\", url).Start()\n    case \"windows\":\n        err = exec.Command(\"rundll32\", \"url.dll,FileProtocolHandler\", url).Start()\n    case \"darwin\":\n        err = exec.Command(\"open\", url).Start()\n    default:\n        err = fmt.Errorf(\"unsupported platform\")\n    }\n    if err != nil {\n        log.Fatal(err)\n    }\n}")
	create(db, "Concatenate two slices\n(tags: concat two slice)\n---\n\nab := append(a, b...)")
	create(db, "String length\n(tags: string length)\n---\n\nn := utf8.RuneCountInString(s)")
	create(db, "Make HTTP POST request\n(tags: make http post request)\n---\n\nresponse, err := http.Post(u, contentType, body)")
	create(db, "Bytes to hex string\n(tags: byte hex string)\n---\n\ns := hex.EncodeToString(a)")
	create(db, "Hex string to byte array\n(tags: byte hex string)\n---\n\na, err := hex.DecodeString(s)")
	create(db, "Find files with a given list of filename extensions\n(tags: file extension walk)\n---\n\nL := []string{}\nerr := filepath.Walk(D, func(path string, info os.FileInfo, err error) error {\n    if err != nil {\n        fmt.Printf(\"failure accessing a path %q: %v\\n\", path, err)\n        return err\n    }\n    for _, ext := range []string{\".jpg\", \".jpeg\", \".png\"} {\n        if strings.HasSuffix(path, ext) {\n            L = append(L, path)\n            break\n        }\n    }\n    return nil\n})\n")
	create(db, "Check if point is inside rectangle\n(tags: check point inside rect)\n---\n\np := image.Pt(x, y)\nr := image.Rect(x1, y1, x2, y2)\nb := p.In(r)")
	create(db, "List files in directory\n(tags: list file dir)\n---\n\nx, err := ioutil.ReadDir(d)")
	create(db, "Make HTTP PUT request\n(tags: make http put request)\n---\n\nreq, err := http.NewRequest(\"PUT\", u, body)\nreq.Header.Set(\"Content-Type\", contentType)\nreq.ContentLength = contentLength\nresponse, err := http.DefaultClient.Do(req)")
	create(db, "Execute function in 30 seconds\n(tags: exec func after time)\n---\n\ntimer := time.AfterFunc(\n    30*time.Second,\n    func() {\n        f(42)\n    })")
	create(db, "Matrix multiplication\n(tags: matrix multiply)\n---\n\nc := new(mat.Dense)\nc.Mul(a, b)")
	create(db, "Filter and transform slice\n(tags: filter transform slice)\n---\n\nvar y []Result\nfor _, e := range x {\n    if P(e) {\n        y = append(y, T(e))\n    }\n}")
	create(db, "Get an environment variable\n(tags: env var)\n---\n\nfoo, ok := os.LookupEnv(\"FOO\")\nif !ok {\n    foo = \"none\"\n}")
	create(db, "Create folder\n(tags: create folder dir)\n---\n\nerr := os.MkdirAll(path, os.ModeDir)")
	create(db, "Pad string on the right\n(tags: pad string right)\n---\n\nif n := utf8.RuneCountInString(s); n < m {\n    s += strings.Repeat(c, m-n)\n}")
	create(db, "Pad string on the left\n(tags: pad string left)\n---\n\nif n := utf8.RuneCountInString(s); n < m {\n    s = strings.Repeat(c, m-n) + s\n}")
	create(db, "Progress bar\n(tags: progress bar)\n---\n\nfunc printProgressBar(n int, total int) {\n    var bar []string\n    tantPerFourty := int((float64(n) / float64(total)) * 40)\n    tantPerCent := int((float64(n) / float64(total)) * 100)\n    for i := 0; i < tantPerFourty; i++ {\n        bar = append(bar, \"█\")\n    }\n    progressBar := strings.Join(bar, \"\")\n    fmt.Printf(\"\\r \" + progressBar + \" - \" + strconv.Itoa(tantPerCent) + \"\")\n}")
	create(db, "Create a zip archive\n(tags: create zip archive)\n---\n\nbuf := new(bytes.Buffer)\nw := zip.NewWriter(buf)\nfor _, filename := range list {\n    input, err := os.Open(filename)\n    output, err := w.Create(filename)\n    _, err = io.Copy(output, input)\n}\nerr := w.Close()\nerr = ioutil.WriteFile(name, buf.Bytes(), 0777)")
	create(db, "Slice intersection\n(tags: slice intersection)\n---\n\nseta := make(map[T]bool, len(a))\nfor _, x := range a {\n    seta[x] = true\n}\nsetb := make(map[T]bool, len(a))\nfor _, y := range b {\n    setb[y] = true\n}\n\nvar c []T\nfor x := range seta {\n    if setb[x] {\n        c = append(c, x)\n    }\n}")
	create(db, "Replace multiple spaces with single space\n(tags: replace space)\n---\n\nwhitespaces := regexp.MustCompile('\\s+')\nt := whitespaces.ReplaceAllString(s, \" \")")
	create(db, "Create a tuple value\n(tags: create tuple interface)\n---\n\nt := []interface{}{\n    2.5,\n    \"hello\",\n    make(chan int),\n}")
	create(db, "Remove all non-digits chars\n(tags: remove digit char)\n---\n\nre := regexp.MustCompile(\"[^\\\\d]\")\nt := re.ReplaceAllLiteralString(s, \"\")")
	create(db, "Add element to the beginning of the slice\n(tags: add beginning slice)\n---\n\nitems = append([]T{x}, items...)")
	create(db, "Copy slice\n(tags: copy slice)\n---\n\ny := make([]T, len(x))\ncopy(y, x)")
	create(db, "Copy file\n(tags: copy file)\n---\n\nfunc copy(dst, src string) error {\n    data, err := ioutil.ReadFile(src)\n    stat, err := os.Stat(src)\n    return ioutil.WriteFile(dst, data, stat.Mode())\n}")
	create(db, "Cancel an operation\n(tags: cancel operation func)\n---\n\nctx, cancel := context.WithCancel(context.Background())\ngo p(ctx)\nsomethingElse()\ncancel()")
	create(db, "Timeout\n(tags: timeout operation func)\n---\n\nctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\ndefer cancel()\np(ctx)")
	create(db, "Check if bytes are a valid UTF-8 string\n(tags: check byte valid utf8)\n---\n\nb := utf8.Valid(s)")
	create(db, "Encode bytes to base64\n(tags: encode byte base64)\n---\n\ns := base64.StdEncoding.EncodeToString(data)")
	create(db, "Decode base64\n(tags: decode string base64)\n---\n\ndata, err := base64.StdEncoding.DecodeString(s)")
	create(db, "Set value on field of structure in map\n(tags: set value field struct map)\n---\n\ntemp := m[key]\ntemp.SomeField = 42\nm[key] = temp")

	//create(db, "Set value on field of structure in map\n(tags: set value field struct map)\n---\n\ntemp := m[key]\ntemp.SomeField = 42\nm[key] = temp")

}

func create(db *gorm.DB, query string) Topic {
	topic := ConvertQueryToTopic(query)
	tags := topic.Tags
	for i, tag := range tags {
		db.Debug().Where("name = ?", tag.Name).Find(&tag)
		if tag.Id == 0 {
			db.Debug().Create(&tag)
		}
		tags[i] = tag
	}
	topic.Tags = tags
	db.Debug().Create(topic)
	return *topic
}

func (n *Golang) GetTag() string {
	return n.golangStg.Tag
}

func (n *Golang) IsReady() bool {
	return n.ready
}

func (n *Golang) Query(query string) (string, error) {
	formatStr := "*%s*\n_(tags:%v)_\n---\n`%s`"

	// create new
	if strings.HasPrefix(query, "+") {
		query = query[1:]
		newTopic := create(n.db, query)
		return fmt.Sprintf(formatStr, newTopic.Title, newTopic.GetTagsString(), newTopic.Code), nil
	}

	// get all
	if strings.HasPrefix(query, "*") {
		var topics []Topic
		n.db.Find(&topics)
		//var topic Topic
		n.db.Model(&topics).Association("Tags").Find(&topics)
		res := ""
		for _, topic := range topics {
			n.db.Model(&topic).Association("Tags").Find(&topic.Tags)
			res += "\n\n===\n" + fmt.Sprintf(formatStr, topic.Title, topic.GetTagsString(), topic.Code)
			if len(res) > 4096 {
				break
			}
		}
		res = res[:strings.LastIndex(res, "===")]
		return res, nil
	}

	queryTags := strings.Split(strings.ToLower(query), " ")

	// get matched tags
	var tags []Tag
	n.db.Where("name IN ?", queryTags).Find(&tags)

	// get associated topics by tags
	var topics []Topic
	n.db.Model(&tags).Association("Topics").Find(&topics)

	// make set deduplicate topics
	set := make(map[string]Topic)
	for _, topic := range topics {
		set[topic.Title] = topic
	}

	// make slice for counting matches
	matches := make([]matchTopic, 0, len(set))
	for _, topic := range set {
		// add tags to topic struct
		n.db.Model(&topic).Association("Tags").Find(&topic.Tags)
		match := 0
		for _, tag := range topic.Tags {
			for _, queryTag := range queryTags {
				if tag.Name == queryTag {
					match++
				}
			}
		}
		matches = append(matches, matchTopic{match: match, topic: topic})
	}

	// sort slice by matches desc
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].match > matches[j].match
	})

	res := ""
	for i, val := range matches {
		res += "\n\n===\n" + fmt.Sprintf(formatStr, val.topic.Title, val.topic.GetTagsString(), val.topic.Code)
		if i > 1 {
			break
		}
	}

	return res, nil
}

type matchTopic struct {
	match int
	topic Topic
}
