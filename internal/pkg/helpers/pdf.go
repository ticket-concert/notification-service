package helpers

import (
	"bytes"
	"context"
	"html/template"
	"log"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GeneratePDF(templatePath string, object interface{}) (*[]byte, error) {

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("----Error ParseFiles----")
		log.Println(err)
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, object); err != nil {
		log.Println("----Error t.Execute----")
		log.Println(err)
		return nil, err
	}
	body := buf.String()

	// w := time.Now().Unix()
	// err1 := os.WriteFile("contract/"+strconv.FormatInt(int64(w), 10)+".html", buf.Bytes(), 0644)
	// if err1 != nil {
	// 	fmt.Println(err1)
	// }

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var result []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				log.Println(err)
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, body).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			result, _, err = page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println(err)
			return err
		}),
	); err != nil {
		log.Println("----Error----")
		log.Println(err)
		return nil, err
	}

	return &result, nil
}
