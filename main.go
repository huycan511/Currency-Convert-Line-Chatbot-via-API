package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Response struct {
	// The right side is the name of the JSON variable
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	richMenu := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
		Selected:    false,
		Name:        "Menu1",
		ChatBarText: "MENU",
		Areas: []linebot.AreaDetail{
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 1250, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: "/list",
					Text: "Call me list",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1250, Y: 0, Width: 1250, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: "/convert",
					Text: "I need to convert",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 843, Width: 1250, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: "/historical",
					Text: "Check exchange of currency",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1250, Y: 843, Width: 1250, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: "/about",
					Text: "Tell me about yourself",
				},
			},
		},
	}

	//create rich menu
	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.RichMenuID)

	//upload richmenu image
	if _, err = bot.UploadRichMenuImage(res.RichMenuID, "Untitled design.jpeg").Do(); err != nil {
		log.Fatal(err)
	}
	if _, err := bot.SetDefaultRichMenu(res.RichMenuID).Do(); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, _ := bot.ParseRequest(r)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := strconv.Atoi(message.Text); err == nil {
					bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("Please select the input currency").
							WithQuickReplies(linebot.NewQuickReplyItems(
								linebot.NewQuickReplyButton(
									//app.appBaseURL+"/static/bo.jpg",
									"https://static-s.aa-cdn.net/img/ios/1184478887/d33c19f34a4fa481e8f8746203cd412f?v=1",
									linebot.NewPostbackAction("USD", "/input_currency="+message.Text+"=USD", "", "USD")),
								linebot.NewQuickReplyButton(
									//app.appBaseURL+"/static/quick/sushi.png",
									"https://icons8.com/iconizer/files/Primo/orig/currency_blue_euro.png",
									linebot.NewPostbackAction("EUR", "/input_currency="+message.Text+"=EUR", "", "EUR")),
								linebot.NewQuickReplyButton(
									//app.appBaseURL+"/static/quick/sushi.png",
									"https://cdn2.iconfinder.com/data/icons/world-currency/512/17-128.png",
									linebot.NewPostbackAction("VND", "/input_currency="+message.Text+"=VND", "", "VND")),
								linebot.NewQuickReplyButton(
									//app.appBaseURL+"/static/quick/sushi.png",
									"https://www.shareicon.net/data/256x256/2015/10/09/653352_money_512x512.png",
									linebot.NewPostbackAction("JPY", "/input_currency="+message.Text+"=JPY", "", "JPY")),
								linebot.NewQuickReplyButton(
									//app.appBaseURL+"/static/quick/sushi.png",
									"https://cdn4.iconfinder.com/data/icons/ios-edge-glyph-5/25/GBP-128.png",
									linebot.NewPostbackAction("GBP", "/input_currency="+message.Text+"=GBP", "", "GBP")),
							)),
					).Do()
				}
			}
		}
		if event.Type == linebot.EventTypePostback {
			data := event.Postback.Data
			if data == "/convert" {
				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("Insert number you want to convert")).Do()
			}
			if strings.Contains(data, "/input_currency") {
				strs1 := strings.Split(data, "=")
				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("Please select the output currency").
						WithQuickReplies(linebot.NewQuickReplyItems(
							linebot.NewQuickReplyButton(
								"https://static-s.aa-cdn.net/img/ios/1184478887/d33c19f34a4fa481e8f8746203cd412f?v=1",
								linebot.NewPostbackAction("USD", "/out_currency="+strs1[1]+"="+strs1[2]+"=USD", "", "USD")),
							linebot.NewQuickReplyButton(
								"https://icons8.com/iconizer/files/Primo/orig/currency_blue_euro.png",
								linebot.NewPostbackAction("EUR", "/out_currency="+strs1[1]+"="+strs1[2]+"=EUR", "", "EUR")),
							linebot.NewQuickReplyButton(
								"https://cdn2.iconfinder.com/data/icons/world-currency/512/17-128.png",
								linebot.NewPostbackAction("VND", "/out_currency="+strs1[1]+"="+strs1[2]+"=VND", "", "VND")),
							linebot.NewQuickReplyButton(
								"https://www.shareicon.net/data/256x256/2015/10/09/653352_money_512x512.png",
								linebot.NewPostbackAction("JPY", "/out_currency="+strs1[1]+"="+strs1[2]+"=JPY", "", "JPY")),
							linebot.NewQuickReplyButton(
								"https://cdn4.iconfinder.com/data/icons/ios-edge-glyph-5/25/GBP-128.png",
								linebot.NewPostbackAction("GBP", "/out_currency="+strs1[1]+"="+strs1[2]+"=GBP", "", "GBP")),
						)),
				).Do()
			}
			if strings.Contains(data, "/out_currency") {
				strs2 := strings.Split(data, "=")
				//sds
				num, _ := strconv.Atoi(strs2[1])
				app := convert()
				jsonString := `{
					"type": "bubble",
					"styles": {
					  "footer": {
						"separator": true
					  }
					},
					"body": {
					  "type": "box",
					  "layout": "vertical",
					  "contents": [
						{
						  "type": "text",
						  "text": "CONVERT",
						  "weight": "bold",
						  "color": "#1DB446",
						  "size": "sm"
						},
						{
						  "type": "text",
						  "text": "RESULT",
						  "weight": "bold",
						  "size": "xxl",
						  "margin": "md"
						},
						{
						  "type": "separator",
						  "margin": "xxl"
						},
						{
							"type": "box",
							"layout": "vertical",
							"contents": [
							  {
						   "type": "spacer",
						   "size": "xl"
						 },
							  {
							"type": "text",
							"text": "` + strs2[1] + " " + strs2[2] + `",
							"align": "center"
						  }
							]
						  },{
							"type":"box",
							"layout":"horizontal",
							"contents":[{
							"type": "image",
							"url": "https://cdn.onlinewebfonts.com/svg/img_148574.png",
							"size": "xs"
						  }
						  ]
						  },
						  {
							"type": "text",
							"text": "` + FloatToString(app.Rates[strs2[3]]/app.Rates[strs2[2]]*float64(num)) + " " + strs2[3] + `",
							"align": "center"
						  }
						  ,
						{
						  "type": "separator",
						  "margin": "xxl"
						},
						{
						  "type": "box",
						  "layout": "horizontal",
						  "margin": "md",
						  "contents": [
							{
							  "type": "button",
						  "style": "secondary",
						  "action": {
							"type": "postback",
							"data": "/convert",
							"label": "TRY AGAIN"
							}}
						  ]
						}
					  ]
					}
				  }`
				contents, _ := linebot.UnmarshalFlexMessageJSON([]byte(jsonString))

				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("Flex message alt text", contents),
				).Do()
			}
			if data == "/historical" {
				template := linebot.NewButtonsTemplate(
					"", "", "Select beginning time",
					linebot.NewDatetimePickerAction("Select here", "date_begin", "date", "", "", ""),
				)
				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTemplateMessage("Datetime pickers alt text", template),
				).Do()
			}
			if data == "date_begin" {
				template := linebot.NewButtonsTemplate(
					"", "", "Select ending time",
					linebot.NewDatetimePickerAction("Select here", "date_end="+event.Postback.Params.Date, "date", "", "", ""),
				)
				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTemplateMessage("Datetime pickers alt text", template),
				).Do()
			}
			if strings.Contains(data, "date_end") {
				strs := strings.Split(data, "=")
				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("Please select the currency").
						WithQuickReplies(linebot.NewQuickReplyItems(
							linebot.NewQuickReplyButton(
								"https://static-s.aa-cdn.net/img/ios/1184478887/d33c19f34a4fa481e8f8746203cd412f?v=1",
								linebot.NewPostbackAction("USD", "/currency="+strs[1]+"="+event.Postback.Params.Date+"=USD", "", "USD")),
							linebot.NewQuickReplyButton(
								"https://icons8.com/iconizer/files/Primo/orig/currency_blue_euro.png",
								linebot.NewPostbackAction("EUR", "/currency="+strs[1]+"="+event.Postback.Params.Date+"=EUR", "", "EUR")),
							linebot.NewQuickReplyButton(
								"https://cdn2.iconfinder.com/data/icons/world-currency/512/17-128.png",
								linebot.NewPostbackAction("VND", "/currency="+strs[1]+"="+event.Postback.Params.Date+"=VND", "", "VND")),
							linebot.NewQuickReplyButton(
								"https://www.shareicon.net/data/256x256/2015/10/09/653352_money_512x512.png",
								linebot.NewPostbackAction("JPY", "/currency="+strs[1]+"="+event.Postback.Params.Date+"=JPY", "", "JPY")),
							linebot.NewQuickReplyButton(
								"https://cdn4.iconfinder.com/data/icons/ios-edge-glyph-5/25/GBP-128.png",
								linebot.NewPostbackAction("GBP", "/currency="+strs[1]+"="+event.Postback.Params.Date+"=GBP", "", "GBP")),
						)),
				).Do()
			}
			if strings.Contains(data, "currency") {
				var change, h1, h2 string
				strs := strings.Split(data, "=")
				app1 := rateDay(strs[1], strs[3])
				app2 := rateDay(strs[2], strs[3])
				rate := app1.Rates[strs[3]] - app2.Rates[strs[3]]
				if rate > 0 {
					change = "Down: " + FloatToString(rate*100/app1.Rates[strs[3]]) + "%"
				} else if rate < 0 {
					change = "Up: " + FloatToString(rate*-100/app1.Rates[strs[3]]) + "%"
				} else {
					change = "Change: 0%"
				}
				percent := app1.Rates[strs[3]] / app2.Rates[strs[3]]
				if percent > 1 {
					h1 = "3"
					h2 = FloatToString(3 / percent)
				} else if percent < 1 {
					h2 = "3"
					h1 = FloatToString(3 * percent)
				} else {
					h1 = "3"
					h2 = "3"
				}
				jsonString := `{  
					"type":"bubble",
					"header":{  
					   "type":"box",
					   "layout":"vertical",
					   "contents":[  
						  {  
							 "type":"text",
							 "text":"` + strs[3] + ` Rate Exchange",
							 "align":"center",
							 "color":"#ffffff",
							 "size":"md"
						  }
					   ]
					},
					"styles":{  
					   "header":{  
						  "backgroundColor":"#000000"
					   },
					   "footer":{  
						  "backgroundColor":"#aaaaff"
					   }
					},
					"body":{  
					   "type":"box",
					   "layout":"vertical",
					   "contents":[  
						  {  
							 "type":"box",
							 "layout":"horizontal",
							 "contents":[  
								{  
								   "type":"text",
								   "text":"` + change + `",
								   "align":"center"
								}
							 ]
						  },
						  {  
							 "type":"box",
							 "layout":"horizontal",
							 "contents":[  
								{  
								   "type":"text",
								   "text":"` + FloatToString(app1.Rates[strs[3]]) + `",
								   "align":"center"
								},
								{  
								   "type":"text",
								   "text":"` + FloatToString(app2.Rates[strs[3]]) + `",
								   "align":"center"
								}
							 ]
						  },
						  {  
							 "type":"box",
							 "layout":"horizontal",
							 "contents":[  
								{  
								   "type":"image",
								   "url":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1e/A_blank_black_picture.jpg/1536px-A_blank_black_picture.jpg",
								   "size":"sm",
								   "aspectRatio":"1:` + h1 + `",
								   "aspectMode":"cover",
								   "gravity":"bottom"
								},
								{  
								   "type":"image",
								   "url":"https://upload.wikimedia.org/wikipedia/commons/thumb/1/1e/A_blank_black_picture.jpg/1536px-A_blank_black_picture.jpg",
								   "size":"sm",
								   "aspectRatio":"1:` + h2 + `",
								   "aspectMode":"cover",
								   "gravity":"bottom"
								}
							 ]
						  },
						  {  
							 "type":"box",
							 "layout":"horizontal",
							 "contents":[  
								{  
								   "type":"box",
								   "layout":"vertical",
								   "contents":[  
									  {  
										 "type":"text",
										 "text":"` + strs[1] + `",
										 "align":"center"
									  }
								   ]
								},
								{  
								   "type":"box",
								   "layout":"vertical",
								   "contents":[  
									  {  
										 "type":"text",
										 "text":"` + strs[2] + `",
										 "align":"center"
									  }
								   ]
								}
							 ]
						  },
						  {  
							 "type":"box",
							 "layout":"horizontal",
							 "contents":[  
								{  
								   "type":"text",
								   "text":"BASE: EUR",
								   "align":"center"
								}
							 ]
						  }
					   ]
					}
				 }`
				contents, _ := linebot.UnmarshalFlexMessageJSON([]byte(jsonString))
				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("Flex message alt text", contents),
				).Do()
			}
			if data == "/about" {
				jsonString := `{
					"type": "bubble",
					"hero": {
						"type": "image",
						"url": "https://imgix.ranker.com/user_node_img/50025/1000492230/original/brandon-stark-tv-characters-photo-u1?w=650&q=50&fm=pjpg&fit=crop&crop=faces",
						"size": "full",
						"aspectRatio": "20:13",
						"aspectMode": "cover",
						"action": {
						"type": "uri",
						"uri": "http://linecorp.com/"
						}
					},
					"body": {
						"type": "box",
						"layout": "vertical",
						"contents": [
						{
							"type": "text",
							"text": "House Stark",
							"weight": "bold",
							"size": "xl"
						},
						{
							"type": "box",
							"layout": "baseline",
							"margin": "md",
							"contents": [
							{
								"type": "icon",
								"size": "sm",
								"url": "https://i.pinimg.com/originals/a4/b2/74/a4b2743bb37453e6132b42015ce00d26.jpg"
							},
							{
								"type": "icon",
								"size": "sm",
								"url": "https://i.pinimg.com/originals/a4/b2/74/a4b2743bb37453e6132b42015ce00d26.jpg"
							},
							{
								"type": "icon",
								"size": "sm",
								"url": "https://i.pinimg.com/originals/a4/b2/74/a4b2743bb37453e6132b42015ce00d26.jpg"
							},
							{
								"type": "icon",
								"size": "sm",
								"url": "https://i.pinimg.com/originals/a4/b2/74/a4b2743bb37453e6132b42015ce00d26.jpg"
							},
							{
								"type": "icon",
								"size": "sm",
								"url": "https://i.pinimg.com/originals/a4/b2/74/a4b2743bb37453e6132b42015ce00d26.jpg"
							},
							{
								"type": "text",
								"text": "5.0",
								"size": "sm",
								"color": "#999999",
								"margin": "md",
								"flex": 0
							}
							]
						},
						{
							"type": "box",
							"layout": "vertical",
							"margin": "lg",
							"spacing": "sm",
							"contents": [
							{
								"type": "box",
								"layout": "baseline",
								"spacing": "sm",
								"contents": [
								{
									"type": "text",
									"text": "Place",
									"color": "#aaaaaa",
									"size": "sm",
									"flex": 1
								},
								{
									"type": "text",
									"text": "Winterfell, The North, Westeros",
									"wrap": true,
									"color": "#666666",
									"size": "sm",
									"flex": 5
								}
								]
							},
							{
								"type": "box",
								"layout": "baseline",
								"spacing": "sm",
								"contents": [
								{
									"type": "text",
									"text": "Time",
									"color": "#aaaaaa",
									"size": "sm",
									"flex": 1
								},
								{
									"type": "text",
									"text": "06:00 - 09:00",
									"wrap": true,
									"color": "#666666",
									"size": "sm",
									"flex": 5
								}
								]
							}
							]
						},{
							"type": "box",
							"layout": "vertical",
							"margin": "lg",
							"spacing": "md",
							"contents": [
								{ "type": "image", "url": "https://qr-official.line.me/M/bJSU9LVfx8.png", "size": "lg" }
							]
						}
						]
					},
					"footer": {
						"type": "box",
						"layout": "vertical",
						"spacing": "sm",
						"contents": [
						{
							"type": "button",
							"style": "link",
							"height": "sm",
							"action": {
							"type": "uri",
							"label": "CALL",
							"uri": "https://gameofthrones.fandom.com/wiki/House_Stark"
							}
						},
						{
							"type": "button",
							"style": "link",
							"height": "sm",
							"action": {
							"type": "uri",
							"label": "WEBSITE",
							"uri": "https://gameofthrones.fandom.com/wiki/House_Stark"
							}
						},
						{
							"type": "spacer",
							"size": "sm"
						}
						],
						"flex": 0
					}
					}`
				contents, _ := linebot.UnmarshalFlexMessageJSON([]byte(jsonString))

				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("Flex message alt text", contents),
				).Do()
			}
			if data == "/list" {
				res := getRate()
				tm := time.Unix(res.Timestamp, 0)
				jsonString := `{
					"type": "bubble",
					"styles": {
					  "footer": {
						"separator": true
					  }
					},
					"body": {
					  "type": "box",
					  "layout": "vertical",
					  "contents": [
						{
						  "type": "text",
						  "text": "LIVE",
						  "weight": "bold",
						  "color": "#FF0000",
						  "size": "sm"
						},
						{
						  "type": "text",
						  "text": "Exchange rates ",
						  "weight": "bold",
						  "size": "xxl",
						  "margin": "md"
						},
						{
						  "type": "text",
						  "text": "` + tm.String() + `",
						  "size": "xs",
						  "color": "#aaaaaa",
						  "wrap": true
						},
						{
						  "type": "text",
						  "text": "BASE: EUR",
						  "size": "xs",
						  "margin": "md",
						  "weight": "bold",
						  "color": "#0000FF",
						  "wrap": true
						},
						{
						  "type": "separator",
						  "margin": "xl"
						},
						{
						  "type": "box",
						  "layout": "vertical",
						  "margin": "xxl",
						  "spacing": "sm",
						  "contents": [
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "icon",
								  "url": "https://cdn3.iconfinder.com/data/icons/currency-2/460/US-dollar-512.png",
								  "size": "xxl"
								},
								{
								  "type": "text",
								  "text": "`+FloatToString(res.Rates["USD"])+` USD",
								  "size": "md",
								  "color": "#555555",
								  "align": "center"
								}
							  ]
							}
						  ]
						},{
						  "type": "separator",
						  "margin": "xl"
						},
						{
						  "type": "box",
						  "layout": "vertical",
						  "margin": "xxl",
						  "spacing": "sm",
						  "contents": [
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "icon",
								  "url": "https://cdn4.iconfinder.com/data/icons/ios-edge-glyph-5/25/GBP-128.png",
								  "size": "xxl"
								},
								{
								  "type": "text",
								  "text": "`+FloatToString(res.Rates["GBP"])+` GBP",
								  "size": "md",
								  "color": "#555555",
								  "align": "center"
								}
							  ]
							}
						  ]
						},
						{
						  "type": "separator",
						  "margin": "xl"
						},
						{
						  "type": "box",
						  "layout": "vertical",
						  "margin": "xxl",
						  "spacing": "sm",
						  "contents": [
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "icon",
								  "url": "https://www.shareicon.net/data/256x256/2015/10/09/653352_money_512x512.png",
								  "size": "xxl"
								},
								{
								  "type": "text",
								  "text": "`+FloatToString(res.Rates["JPY"])+` JPY",
								  "size": "md",
								  "color": "#555555",
								  "align": "center"
								}
							  ]
							}
						  ]
						},
						{
						  "type": "separator",
						  "margin": "xl"
						},
						{
						  "type": "box",
						  "layout": "vertical",
						  "margin": "xxl",
						  "spacing": "sm",
						  "contents": [
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "icon",
								  "url": "https://cdn2.iconfinder.com/data/icons/world-currency/512/17-128.png",
								  "size": "xxl"
								},
								{
								  "type": "text",
								  "text": "`+FloatToString(res.Rates["VND"])+` VND",
								  "size": "md",
								  "color": "#555555",
								  "align": "center"
								}
							  ]
							}
						  ]
						}
					  ]
					}
				  }`
				contents, _ := linebot.UnmarshalFlexMessageJSON([]byte(jsonString))
				bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("Flex message alt text", contents),
				).Do()
			}
		}
	}
}
func convert() *Response {
	var res Response
	response, _ := http.Get("http://data.fixer.io/api/latest?access_key=526d2da077cccc0a840ef1eaf2328047")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body, &res)
	return &res
}
func getRate() *Response {
	var res Response
	response, _ := http.Get("http://data.fixer.io/api/latest?access_key=526d2da077cccc0a840ef1eaf2328047&format=1&symbols=VND,USD,GBP,JPY")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body, &res)
	return &res
}
func rateDay(day string, currency string) *Response {
	var res Response
	response, _ := http.Get("http://data.fixer.io/api/" + day + "?access_key=526d2da077cccc0a840ef1eaf2328047&symbols=" + currency)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body, &res)
	return &res
}
func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
