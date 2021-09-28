package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"net/url"
	"sync"
)

type Context struct {
	Context context.Context
	Cancel  context.CancelFunc
}

/*
How I want my code:

// 1) Search Query
var bb BestBuy
bb.Scrape("rtx 3080")

// How I want data to be returned in the terminal:



*/

type BestBuy string

func (bb *BestBuy) Scrape() {
	// Allocate a Context with custom options.
	var ac Context
	ac.Context, ac.Cancel = chromedp.NewExecAllocator(context.Background(), []chromedp.ExecAllocatorOption{
		chromedp.WindowSize(1000, 1000),
		chromedp.Flag("headless", false),
	}...)

	var taskCtx Context
	taskCtx.Context, taskCtx.Cancel = chromedp.NewContext(ac.Context, chromedp.WithLogf(log.Printf))

	if err := chromedp.Run(taskCtx.Context, Tasks(string(*bb))...); err != nil {
		log.Fatal(err)
	}

	defer func() {
		ac.Cancel()
		taskCtx.Cancel()
	}()
}

// Tasks holds the collection of all the tasks in order.
func Tasks(s string) []chromedp.Action {
	return []chromedp.Action{
		task1(s),
		task2(),
	}
}

// task1 goes to the first page of your search.
func task1(s string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://www.bestbuy.com/site/searchpage.jsp?st=" + url.QueryEscape(s)),
	}
}

// task2 will get the name of the product and its availability
func task2() chromedp.Tasks {
	var items []*cdp.Node
	return chromedp.Tasks{
		chromedp.Nodes(".right-column", &items, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, val := range items {
				fmt.Printf("%#v\n", val)
			}
			return nil
		}),
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		//goland:noinspection SpellCheckingInspection
		bestbuy1 := BestBuy("rtx 3080")
		bestbuy1.Scrape()
		defer func() {
			wg.Done()
		}()
	}()

	go func() {
		//goland:noinspection SpellCheckingInspection
		bestbuy2 := BestBuy("rtx 3080")
		bestbuy2.Scrape()
		defer func() {
			wg.Done()
		}()
	}()

	go func() {
		//goland:noinspection SpellCheckingInspection
		bestbuy3 := BestBuy("rtx 3090")
		bestbuy3.Scrape()
		defer func() {
			wg.Done()
		}()
	}()

	wg.Wait()
}
