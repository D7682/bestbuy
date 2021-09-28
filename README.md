/*
Code used to use javascript in its own file.

file, err := ioutil.ReadFile("extra.js")
if err != nil {
log.Fatal(err)
}
chromedp.Evaluate(string(file), nil, chromedp.ByQuery)
*/