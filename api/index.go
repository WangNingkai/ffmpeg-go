package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"

	. "github.com/tbxark/g4vercel"
)

func GetMediaDurationByUrl(url string) (*ffprobe.Format, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, url)
	if err != nil {
		fmt.Printf("ErrorX GetMediaDurationByUrl: %v", err.Error())
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	return data.Format, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	server := New()
	server.Use(Recovery(func(err interface{}, c *Context) {
		if httpError, ok := err.(HttpError); ok {
			c.JSON(httpError.Status, H{
				"message": httpError.Error(),
			})
		} else {
			message := fmt.Sprintf("%s", err)
			c.JSON(500, H{
				"message": message,
			})
		}
	}))
	server.GET("/", func(context *Context) {
		text := context.Query("text")
		duration, err := GetMediaDurationByUrl(text)
		if err != nil {
			context.JSON(500, H{
				"message": "unable to fetch media data.",
			})
			return
		}

		context.JSON(200, H{
			"data": duration,
		})
	})

	server.Handle(w, r)
}
