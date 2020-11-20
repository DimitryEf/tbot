package wiki

type Title struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Normalized []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"normalized"`
		Pages struct {
			Num2089687 struct {
				Pageid  int    `json:"pageid"`
				Ns      int    `json:"ns"`
				Title   string `json:"title"`
				Extract string `json:"extract"`
			} `json:"2089687"`
		} `json:"pages"`
	} `json:"query"`
}
